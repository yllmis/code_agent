package tools

import (
	"fmt"
	"os"
)

// ReadParams 读文件工具的参数
type ReadParams struct {
	Path string `json:"path" description:"文件路径"`
}

// ReadTool 读取指定路径的文件内容
type ReadTool struct{}

func (r *ReadTool) Name() string        { return "read_file" }
func (r *ReadTool) Description() string { return "读取指定路径的文件内容" }

func (r *ReadTool) Schema() ToolSchema {
	return StructToSchema("read_file", r.Description(), ReadParams{})
}

// Execute 读取指定路径的文件，返回文件内容
// 输入: ReadParams
// 输出: 文件完整内容
func (r *ReadTool) Execute(args any) (string, error) {
	params, ok := args.(ReadParams)
	if !ok {
		return "", fmt.Errorf("参数类型错误，期望 ReadParams")
	}
	if params.Path == "" {
		return "", fmt.Errorf("path 参数不能为空")
	}
	content, err := os.ReadFile(params.Path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
