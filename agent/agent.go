package agent

import (
	"fmt"
	"strings"

	"go-agents/svc"
	"go-agents/tools"
)

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

		// 解析工具调用
		toolName, toolInput, found := parseToolCall(reply)
		if !found {
			// 没有工具调用，LLM 直接回答，返回结果
			return reply, nil
		}

		// 找到对应工具
		tool := findTool(toolList, toolName)
		if tool == nil {
			// 工具不存在，告诉 LLM 重试
			messages = append(messages,
				ChatMessage{Role: "assistant", Content: reply},
				ChatMessage{Role: "user", Content: fmt.Sprintf("工具 %q 不存在，请使用可用的工具", toolName)},
			)
			continue
		}

		// 执行工具
		result, err := tool.Execute(toolInput)
		if err != nil {
			result = fmt.Sprintf("工具执行出错: %v", err)
		}

		// 把 LLM 回复和工具结果加入消息历史，继续循环
		messages = append(messages,
			ChatMessage{Role: "assistant", Content: reply},
			ChatMessage{Role: "user", Content: fmt.Sprintf("工具 %s 的执行结果:\n%s", toolName, result)},
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
	sb.WriteString("\n请一步步思考，使用工具来帮助回答用户的问题。")
	return sb.String()
}

// parseToolCall 从 LLM 回复中解析工具调用
// 格式：<tool>工具名:参数</tool>
func parseToolCall(reply string) (name, input string, found bool) {
	start := strings.Index(reply, "<tool>")
	end := strings.Index(reply, "</tool>")
	if start == -1 || end == -1 || start >= end {
		return "", "", false
	}

	// 提取 <tool> 和 </tool> 之间的内容
	content := reply[start+len("<tool>") : end]

	// 按第一个 : 分割，前面是工具名，后面是参数
	colonIdx := strings.Index(content, ":")
	if colonIdx == -1 {
		return "", "", false
	}

	name = strings.TrimSpace(content[:colonIdx])
	input = strings.TrimSpace(content[colonIdx+1:])
	return name, input, true
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
