package tools

import "reflect"

// Parameter 描述工具的单个参数
type Parameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`        // "string", "int" 等
	Description string `json:"description"` // 告诉 LLM 这个参数是什么
	Required    bool   `json:"required"`    // 是否必填
}

// ToolSchema 工具的参数说明书
type ToolSchema struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  []Parameter `json:"parameters"`
}

// Tool 工具接口
type Tool interface {
	Name() string
	Description() string
	Schema() ToolSchema
	Execute(args any) (string, error)
}

// StructToSchema 从结构体 tag 自动生成 ToolSchema
// 读取 json tag 获取参数名，description tag 获取描述
// omitempty 表示可选参数
func StructToSchema(name, description string, paramType any) ToolSchema {
	t := reflect.TypeOf(paramType)
	params := make([]Parameter, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		descTag := field.Tag.Get("description")

		// 解析 json tag，处理 omitempty
		paramName := jsonTag
		required := true
		if idx := len(jsonTag); idx > 0 {
			for j, c := range jsonTag {
				if c == ',' {
					paramName = jsonTag[:j]
					if jsonTag[j+1:] == "omitempty" {
						required = false
					}
					break
				}
			}
		}

		params = append(params, Parameter{
			Name:        paramName,
			Type:        field.Type.String(),
			Description: descTag,
			Required:    required,
		})
	}

	return ToolSchema{
		Name:        name,
		Description: description,
		Parameters:  params,
	}
}
