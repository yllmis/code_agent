package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-agents/svc"
	"io"
	"net/http"
)

// ChatMessage 对应 OpenAI 格式的消息
// Role: "system" 设定角色, "user" 用户输入, "assistant" LLM 回复
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 发给 大模型 的请求体
// POST https://token-plan-cn.xiaomimimo.com/v1/chat/completions
type ChatRequest struct {
	Model               string        `json:"model"`                 // 模型名，用 "mimo-v2.5-pro"
	Messages            []ChatMessage `json:"messages"`              // 完整对话历史
	MaxCompletionTokens int           `json:"max_completion_tokens"` // 最大生成 token 数
}

// ChatResponse Mimo API 的返回
type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"` // LLM 的回复
	} `json:"choices"`
}

// Chat 核心函数：把消息发给 Mimo，拿回 LLM 的文本回复
// 实现时注意：
// - 从config MIMO_API_KEY 读 key
// - Header 用 api-key，不是 Authorization: Bearer
// - 返回 choices[0].Message.Content

/*
curl --location --request POST 'BASE_URL/chat/completions' \
--header "api-key: $MIMO_API_KEY" \
--header "Content-Type: application/json" \

	--data-raw '{
	    "model": "mimo-v2.5-pro",
	    "messages": [
	        {
	            "role": "system",
	            "content": "You are MiMo, an AI assistant developed by Xiaomi. Today is date: Tuesday, December 16, 2025. Your knowledge cutoff date is December 2024."
	        },
	        {
	            "role": "user",
	            "content": "please introduce yourself"
	        }
	    ],
	    "max_completion_tokens": 1024
	}'
*/
func Chat(svc *svc.ServiceContext, messages []ChatMessage) (string, error) {
	req := ChatRequest{
		Model:               svc.Config.Model,
		Messages:            messages,
		MaxCompletionTokens: svc.Config.MaxCompletionTokens,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	// base url
	url := svc.Config.BaseURL + "/chat/completions"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	// header
	httpReq.Header.Set("api-key", svc.Config.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// 请求
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", err
	}

	return chatResp.Choices[0].Message.Content, nil
}
