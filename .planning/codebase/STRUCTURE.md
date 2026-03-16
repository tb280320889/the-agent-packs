# Codebase Structure

**Analysis Date:** 2026-03-16

## Directory Layout

```
[project-root]/
├── cmd/                  # Go CLI/MCP 入口
├── internal/             # 核心运行时模块
├── blueprint/            # Blueprint 节点与索引
├── workflow-packages/    # workflow package 目录与注册表
├── tools/                # 工具边界说明文档
├── fixtures/             # 示例请求/结果与 smoke 数据
├── tests/                # Go 测试
├── docs/                 # 文档与改造计划包
├── .planning/            # 规划文档输出目录
├── go.mod                # Go module 定义
└── go.sum                # 依赖锁定
```

## Directory Purposes

**cmd/**
- Purpose: 统一 CLI/MCP 入口
- Contains: `agent-pack-mcp` 主程序
- Key files: `cmd/agent-pack-mcp/main.go`

**internal/**
- Purpose: 核心业务模块（activation/compiler/query/registry/validator）
- Contains: Go 业务逻辑与模型
- Key files: `internal/activation/activation.go`, `internal/compiler/compiler.go`, `internal/query/query.go`, `internal/registry/registry.go`, `internal/validator/runner.go`

**blueprint/**
- Purpose: 路由知识库与索引
- Contains: L0/L1/L2/L3 节点 Markdown、索引与 schema
- Key files: `blueprint/README.md`, `blueprint/schema.md`, `blueprint/index/blueprint.db`

**workflow-packages/**
- Purpose: workflow package 根目录与注册表
- Contains: `registry.json`、各包 `package.yaml`
- Key files: `workflow-packages/registry.json`, `workflow-packages/wxt-manifest/package.yaml`

**tools/**
- Purpose: 各工具边界说明
- Contains: README（activation/compiler/query-mcp）
- Key files: `tools/activation-entry/README.md`, `tools/compiler/README.md`, `tools/query-mcp/README.md`

**fixtures/**
- Purpose: 示例请求/结果与 smoke 测试输入
- Contains: activation 请求、route 结果样例
- Key files: `fixtures/activation-request.sample.json`, `fixtures/route-result.sample.json`

**tests/**
- Purpose: Go 单元/集成测试
- Contains: M1/M2/M3 最薄验证用例
- Key files: `tests/m1_minimal_test.go`, `tests/m2_registry_test.go`

## Key File Locations

**Entry Points:**
- `cmd/agent-pack-mcp/main.go`: CLI/MCP 入口

**Configuration:**
- `workflow-packages/registry.json`: package 注册表
- `blueprint/schema.md`: Blueprint 索引 schema 说明

**Core Logic:**
- `internal/compiler/compiler.go`: Blueprint 编译为 SQLite
- `internal/query/query.go`: 路由查询与 bundle 构建
- `internal/activation/activation.go`: activation 主流程
- `internal/registry/registry.go`: registry 解析与校验
- `internal/validator/*.go`: validator 逻辑

**Testing:**
- `tests/*.go`: Go 测试文件

## Naming Conventions

**Files:**
- Blueprint 节点文件：`blueprint/L{0|1|2|3}/{domain}/.../*.md`（如 `blueprint/L1/wxt/manifest.md`）
- Workflow package：`workflow-packages/<package>/package.yaml`

**Directories:**
- Blueprint 分层：`L0/`、`L1/`、`L2/`、`L3/`（`blueprint/README.md`）
- Workflow packages 统一根：`workflow-packages/`（`workflow-packages/README.md`）

## Where to Add New Code

**New Feature:**
- Primary code: `internal/<module>/`（例如新增路由能力放 `internal/query/`）
- Tests: `tests/<feature>_test.go`

**New Component/Module:**
- Implementation: `internal/<new-module>/` 并在 `cmd/agent-pack-mcp/main.go` 连接入口

**Utilities:**
- Shared helpers: 贴近使用模块（例如 `internal/query/` 中的局部 helper）

## Special Directories

**blueprint/index/**
- Purpose: SQLite 索引与编译报告输出
- Generated: Yes（由 `internal/compiler/compiler.go` 生成）
- Committed: Yes（当前仓库存在 `blueprint/index` 目录引用）

---

*Structure analysis: 2026-03-16*
