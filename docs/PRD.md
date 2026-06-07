# CLI 代码助手 Agent — 需求文档

## 1. 项目目标

构建一个命令行代码助手 Agent，用户通过自然语言提问，Agent 自动搜索和阅读本地代码库，给出解答。

核心能力：让开发者在终端里用中文问代码问题，Agent 像一个熟悉项目的人一样帮你查代码、找定义、解释逻辑。

## 2. 技术架构

```
用户输入 → Agent 循环（ReAct）→ LLM 思考 → 调用本地工具 → 返回结果给 LLM → 继续循环 → 输出回答
```

- 语言：Go
- LLM：Mimo API（OpenAI 兼容格式）
- 配置：YAML 文件 + ServiceContext 依赖注入

## 3. 功能需求

### 3.1 LLM 对话（已实现）

- 调用 Mimo token plan API 进行对话
- 支持多轮对话（保留 messages 历史）
- 配置项：API Key、Model、Base URL

### 3.2 Agent 循环（待实现）

- ReAct 模式：思考 → 调工具 → 观察 → 再思考
- 最大循环次数限制（防止死循环，建议 10 轮）
- 工具调用协议：LLM 返回特定格式文本（如 `<tool>name:input</tool>`），Agent 解析执行
- 将工具描述注入 system prompt，让 LLM 知道有哪些工具可用

### 3.3 本地工具（已实现接口）

| 工具 | 功能 | 输入 | 输出 |
|------|------|------|------|
| search | 搜索文件中包含关键词的行 | 关键词 | 文件路径:行号: 匹配内容 |
| read_file | 读取指定文件内容 | 文件路径 | 文件完整内容 |
| list_files | 列出目录下文件 | 目录路径 | 文件名列表（目录加 /） |

工具接口：
```go
type Tool interface {
    Name() string
    Description() string
    Run(args string) (string, error)  // 注意：当前实现写的是 Execute，需统一
}
```

### 3.4 CLI 交互（待实现）

- 启动时读取 config.yaml，初始化 ServiceContext
- 交互式输入：用户输入问题，Agent 回答，循环
- 退出命令：`exit` 或 `quit`
- 欢迎信息：提示用户可以问什么

## 4. 配置

`config.yaml`（已 gitignore）：

```yaml
api_key: "your-mimo-api-key"
model: "mimo-v2.5-pro"
base_url: "https://token-plan-cn.xiaomimimo.com/v1"
```

## 5. 项目结构

```
code_agent/
├── main.go                 # CLI 入口
├── config.yaml             # 配置文件（不提交）
├── config/config.go        # 配置加载
├── svc/servicecontext.go   # 全局依赖注入
├── agent/
│   ├── llm.go              # Mimo API 调用（已实现）
│   └── agent.go            # Agent 核心循环（待实现）
└── tools/
    ├── tool.go             # Tool 接口定义
    ├── search.go           # 搜索工具
    ├── read.go             # 文件读取工具
    └── list.go             # 目录列表工具
```

## 6. 已知问题

- `tool.go` 接口方法名 `Run` 与工具实现的 `Execute` 不一致，需统一
- `main.go` 为空，需要接入所有模块

## 7. 后续扩展（非本期）

- 支持多文件并行搜索
- 支持正则搜索
- 支持 git blame / git log 工具
- TS 前端 / VS Code 插件集成
