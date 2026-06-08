package tools

import (
	"fmt"
	"os"
)

// ListTool 列出指定目录下的文件和子目录
type ListTool struct{}

func (l *ListTool) Name() string        { return "list_files" }
func (l *ListTool) Description() string { return "列出指定目录下的文件和子目录" }

// Execute 列出 input 目录下的文件
// 输入: 目录路径（如 "." 或 "agent"）
// 输出: 每行一个文件名，目录后面加 /
func (l *ListTool) Execute(input string) (string, error) {
	// 空输入默认当前目录
	if input == "" {
		input = "."
	}

	entries, err := os.ReadDir(input)
	if err != nil {
		return "", err
	}

	var results []string
	for _, entry := range entries {
		// 目录名后面加 / 方便 LLM 区分
		if entry.IsDir() {
			results = append(results, entry.Name()+"/")
		} else {
			results = append(results, entry.Name())
		}
	}

	if len(results) == 0 {
		return "目录为空", nil
	}
	return fmt.Sprintf("%s", results), nil
}
