package tools

type Tool interface {
	Name() string
	Description() string // 描述工具的功能和用途
	Execute(args string) (string, error)
}
