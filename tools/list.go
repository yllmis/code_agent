package tools

import (
	"fmt"
	"os"
)

// ListParams 列目录工具的参数
type ListParams struct {
	Path string `json:"path,omitempty" description:"目录路径，默认当前目录"`
}

// ListTool 列出指定目录下的文件和子目录
type ListTool struct{}

func (l *ListTool) Name() string        { return "list_files" }
func (l *ListTool) Description() string { return "列出指定目录下的文件和子目录" }

func (l *ListTool) Schema() ToolSchema {
	return StructToSchema("list_files", l.Description(), ListParams{})
}

// Execute 列出指定目录下的文件
// 输入: ListParams
// 输出: 每行一个文件名，目录后面加 /
func (l *ListTool) Execute(args any) (string, error) {
	params, ok := args.(ListParams)
	if !ok {
		return "", fmt.Errorf("参数类型错误，期望 ListParams")
	}
	dir := params.Path
	if dir == "" {
		dir = "."
	}

	entries, err := os.ReadDir(dir)
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
