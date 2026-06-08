package tools

import "os"

// ReadTool 读取指定路径的文件内容
type ReadTool struct{}

func (r *ReadTool) Name() string        { return "read_file" }
func (r *ReadTool) Description() string { return "读取指定路径的文件内容" }

// Execute 读取 input 指定的文件路径，返回文件内容
// 输入: 文件路径（如 "agent/llm.go"）
// 输出: 文件完整内容
func (r *ReadTool) Execute(input string) (string, error) {
	// input 是 LLM 给的文件路径，直接读
	content, err := os.ReadFile(input)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
