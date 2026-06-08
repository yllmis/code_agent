# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

CLI 代码助手 Agent，Go 后端实现。核心模式：用户提问 → LLM 思考 → 调用本地工具（搜索/读文件/列目录）→ 观察结果 → 继续循环直到回答。

## 常用命令

```bash
go build ./...          # 编译检查
go run .                # 运行 CLI（需在项目根目录，读取 config.yaml）
go vet ./...            # 静态检查
```

## 架构

- `config/config.go` — 配置结构体 + YAML 加载
- `svc/servicecontext.go` — 全局依赖注入（ServiceContext 模式），持有 Config 等全局状态
- `agent/llm.go` — Mimo API 调用（OpenAI 兼容格式）
- `agent/agent.go` — Agent 核心循环（思考→工具→观察）
- `tools/` — 具体工具实现（search, read_file, list_files）

依赖流向：`main.go` → `svc.NewServiceContext()` → `agent.Chat(svc, messages)`

## LLM API

使用 Mimo token plan API，非 Anthropic SDK：
- Base URL: `https://token-plan-cn.xiaomimimo.com/v1/chat/completions`
- Auth Header: `api-key`（不是 `Authorization: Bearer`）
- 配置在 `config.yaml`（已 gitignore，不提交）

## 开发者

Go 后端开发者，在学习 Agent 模式。当用户没有指出需要设计或修改时，希望 Claude 提供指导和关键改动点，不要直接给完整实现代码。

## Git 提交规范

- 提交信息用中文
- 格式：`类型: 简要描述改动`
- 只写一句话总结整体改动，不要逐文件列出
- 类型：feat / fix / refactor / docs / test
- 示例：`feat: 实现 Chat 函数调用 Mimo API`
- 示例：`fix: 修复 config.yaml 读取路径问题`
