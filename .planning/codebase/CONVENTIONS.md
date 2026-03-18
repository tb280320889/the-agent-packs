# Coding Conventions

**Analysis Date:** 2026-03-16

## Naming Patterns

**Files:**
- Go 源文件使用小写下划线或小写单词：`internal/validator/domain_wxt_manifest.go`、`tests/m1_minimal_test.go`、`internal/activation/activation.go`。

**Functions:**
- 导出函数使用 PascalCase：`internal/activation/activation.go` 的 `Execute`、`internal/query/query.go` 的 `RouteQuery`。
- 包内函数使用 camelCase：`internal/query/query.go` 的 `fetchNodes`、`scoreCandidate`。

**Variables:**
- 使用 camelCase：`internal/query/query.go` 中的 `targetDomainPtr`、`selectedFiles`。
- 常量切片使用小写：`internal/compiler/compiler.go` 的 `requiredKeys`。

**Types:**
- 结构体/类型使用 PascalCase：`internal/model/model.go` 的 `ActivationResult`、`RouteCandidate`。
- 仅包内使用的结构体也遵循 PascalCase：`internal/query/query.go` 的 `packProfile`、`nodeRecord`。

## Code Style

**Formatting:**
- Go 标准 `gofmt` 风格（缩进、对齐、import 分组），无显式配置文件。

**Linting:**
- 未检测到 lint 配置文件（无 `.golangci.*`）。

## Import Organization

**Order:**
1. 标准库
2. 空行分隔
3. 内部模块

示例：`cmd/agent-pack-mcp/main.go`。

**Path Aliases:**
- Go module 路径直接使用 `the-agent-packs/...`：`internal/query/query.go`、`tests/m1_minimal_test.go`。

## Error Handling

**Patterns:**
- 返回 `error` 并立即上抛：`internal/query/query.go` 的 `OpenDB`、`BuildContextBundle`。
- 使用 `errors.New`/`fmt.Errorf` 生成错误：`internal/registry/registry.go`、`internal/compiler/compiler.go`。
- CLI 入口打印错误并 `os.Exit(1)`：`cmd/agent-pack-mcp/main.go`。
- 测试中使用 `t.Fatalf` 失败并输出上下文：`tests/m1_minimal_test.go`。

## Logging

**Framework:**
- 标准库 `fmt`/`os`，无日志框架。

**Patterns:**
- CLI 命令输出状态：`cmd/agent-pack-mcp/main.go` 的 `fmt.Println`、`fmt.Fprintln`。

## Comments

**When to Comment:**
- 代码中几乎不使用注释，依靠函数/变量命名表达语义：`internal/activation/activation.go`、`internal/query/query.go`。

**JSDoc/TSDoc:**
- 不适用（Go 代码库）。

## Function Design

**Size:**
- 大函数承担流程编排，内部拆分小助手函数：`cmd/agent-pack-mcp/main.go` 的 `cmdMCP`，配合 `parseParam*` 系列。

**Parameters:**
- 显式传递上下文与输入结构体：`internal/validator/runner.go` 的 `Run(plan, input)`。

**Return Values:**
- 结构体 + `error` 风格：`internal/query/query.go` 的 `RouteQuery`、`BuildContextBundle`。

## Module Design

**Exports:**
- 包导出以功能边界命名（activation/query/registry/validator）：`internal/activation/activation.go`、`internal/query/query.go`。

**Barrel Files:**
- 不使用（Go 包按目录组织）。

---

*Convention analysis: 2026-03-16*
