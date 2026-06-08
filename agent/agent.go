package agent

import (
	"fmt"
	"strings"
	"sync"

	"go-agents/svc"
	"go-agents/tools"
)

// ToolCall 表示一个工具调用请求
type ToolCall struct {
	Name  string
	Input string
}

// Run Agent 核心循环：用户提问 → LLM 思考 → 工具调用 → 观察 → 继续循环
func Run(svc *svc.ServiceContext, toolList []tools.Tool, userMessage string) (string, error) {
	// 1. 构造 system prompt，注入工具列表和调用格式
	systemPrompt := buildSystemPrompt(toolList)

	// 2. 初始化消息历史
	messages := []ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userMessage},
	}

	// 3. 循环，最多 max_rounds 轮
	for round := 0; round < svc.Config.MaxRounds; round++ {
		// 调 LLM
		reply, err := Chat(svc, messages)
		if err != nil {
			return "", fmt.Errorf("LLM call failed: %w", err)
		}

		// 解析工具调用（可能有多个）
		toolCalls := parseToolCalls(reply)
		if len(toolCalls) == 0 {
			// 没有工具调用，LLM 直接回答
			// 检查回复是否为空
			if strings.TrimSpace(reply) == "" {
				return "", fmt.Errorf("LLM 返回空回复")
			}
			return reply, nil
		}

		// 并发执行所有工具调用
		results := make([]string, len(toolCalls))
		var wg sync.WaitGroup
		for i, tc := range toolCalls {
			wg.Add(1)
			go func(i int, tc ToolCall) {
				defer wg.Done()
				tool := findTool(toolList, tc.Name)
				if tool == nil {
					results[i] = fmt.Sprintf("工具 %q 不存在", tc.Name)
					return
				}
				result, err := tool.Execute(tc.Input)
				if err != nil {
					results[i] = fmt.Sprintf("工具 %s 执行出错: %v", tc.Name, err)
					return
				}
				results[i] = result
			}(i, tc)
		}
		wg.Wait()

		// 把所有工具结果加入消息历史
		var resultText string
		for i, tc := range toolCalls {
			resultText += fmt.Sprintf("工具 %s 的执行结果:\n%s\n\n", tc.Name, results[i])
		}
		messages = append(messages,
			ChatMessage{Role: "assistant", Content: reply},
			ChatMessage{Role: "user", Content: strings.TrimSpace(resultText)},
		)
	}

	return "", fmt.Errorf("超过最大轮次 %d，未得到最终回答", svc.Config.MaxRounds)
}

// buildSystemPrompt 构造系统提示词，告诉 LLM 有哪些工具可用
func buildSystemPrompt(toolList []tools.Tool) string {
	var sb strings.Builder
	sb.WriteString("你是一个代码助手。你可以使用以下工具：\n\n")

	for _, t := range toolList {
		sb.WriteString(fmt.Sprintf("%s: %s\n", t.Name(), t.Description()))
	}

	sb.WriteString("\n当你需要使用工具时，用这个格式：\n")
	sb.WriteString("<tool>工具名:参数</tool>\n")
	sb.WriteString("\n可以一次调用多个工具，每个工具用独立的 <tool> 标签包裹。")
	sb.WriteString("\n请一步步思考，使用工具来帮助回答用户的问题。")
	return sb.String()
}

// parseToolCalls 从 LLM 回复中解析所有工具调用
// 格式：<tool>工具名:参数</tool>，可以有多个
func parseToolCalls(reply string) []ToolCall {
	var calls []ToolCall
	remaining := reply

	for {
		start := strings.Index(remaining, "<tool>")
		end := strings.Index(remaining, "</tool>")
		if start == -1 || end == -1 || start >= end {
			break
		}

		content := remaining[start+len("<tool>") : end]
		colonIdx := strings.Index(content, ":")
		if colonIdx == -1 {
			remaining = remaining[end+len("</tool>"):]
			continue
		}

		name := strings.TrimSpace(content[:colonIdx])
		input := strings.TrimSpace(content[colonIdx+1:])
		calls = append(calls, ToolCall{Name: name, Input: input})

		remaining = remaining[end+len("</tool>"):]
	}

	return calls
}

// findTool 根据工具名查找工具
func findTool(toolList []tools.Tool, name string) tools.Tool {
	for _, t := range toolList {
		if t.Name() == name {
			return t
		}
	}
	return nil
}
