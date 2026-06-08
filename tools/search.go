package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// isBinaryFile 根据扩展名判断是否为二进制文件
func isBinaryFile(path string) bool {
	binaryExts := map[string]bool{
		".exe": true, ".bin": true, ".so": true, ".dll": true,
		".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".ico": true, ".svg": true,
		".mp3": true, ".mp4": true, ".avi": true, ".mov": true,
		".zip": true, ".tar": true, ".gz": true, ".rar": true,
		".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
		".pyc": true, ".class": true, ".o": true, ".a": true,
	}
	ext := strings.ToLower(filepath.Ext(path))
	return binaryExts[ext]
}

type SearchTool struct{}

func (s *SearchTool) Name() string        { return "search" }
func (s *SearchTool) Description() string { return "在当前目录搜索包含关键词的文件，返回文件名、行号和匹配内容" }

// Execute 递归搜索当前目录下所有文件，返回包含关键词的行
// 输入: 关键词（如 "func main"）
// 输出: 每行格式为 "文件路径:行号: 匹配内容"
func (s *SearchTool) Execute(input string) (string, error) {
	var results []string

	// filepath.Walk 会递归遍历目录树
	// 回调函数对每个文件/目录执行一次
	// 返回值控制遍历行为：
	//   nil = 继续
	//   filepath.SkipDir = 跳过当前目录（不是文件）
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		// 遍历出错（如权限不足）跳过这个文件，不中断整个搜索
		if err != nil {
			return nil
		}
		// 跳过 .git 目录，避免搜索 git 内部文件
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		// 目录不是文件，跳过（Walk 会自动进入子目录）
		if info.IsDir() {
			return nil
		}

		// 跳过二进制文件和常见非文本文件
		if isBinaryFile(path) {
			return nil
		}

		// 读取文件全部内容
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		// 按行拆分，逐行检查是否包含关键词
		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			if strings.Contains(line, input) {
				results = append(results, fmt.Sprintf("%s:%d: %s", path, i+1, strings.TrimSpace(line)))
			}
		}

		return nil
	})

	if len(results) == 0 {
		return "未找到匹配结果", nil
	}
	return strings.Join(results, "\n"), nil
}
