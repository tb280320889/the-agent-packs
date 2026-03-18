# Architecture

**Analysis Date:** 2026-03-16

## Pattern Overview

**Overall:** Blueprint 索引驱动的编排管线（Compile -> Route -> Bundle -> Activate -> Validate）

**Key Characteristics:**
- Markdown frontmatter 作为路由知识源，编译为 SQLite 索引（`blueprint/` -> `blueprint/index/blueprint.db`）
- 路由与上下文裁剪基于索引查询（`internal/query/query.go`）
- Activation 汇聚路由结果、上下文包与校验结果（`internal/activation/activation.go`）

## Layers

**Blueprint 内容层：**
- Purpose: 定义 L0/L1/L2/L3 路由节点与执行策略
- Location: `blueprint/`
- Contains: Markdown + frontmatter 节点（例如 `blueprint/L1/wxt/manifest.md`）
- Depends on: 无运行时依赖（静态源数据）
- Used by: Compiler（`internal/compiler/compiler.go`）

**索引编译层：**
- Purpose: 校验 frontmatter 并构建 SQLite 索引与报告
- Location: `internal/compiler/compiler.go`
- Contains: frontmatter 解析、路径/ID一致性校验、节点/边写入
- Depends on: Blueprint 源文件（`blueprint/`）
- Used by: Query 层（`internal/query/query.go`）

**注册表层：**
- Purpose: Workflow package 注册表与命名治理真相源
- Location: `internal/registry/registry.go`, `workflow-packages/registry.json`
- Contains: 包名规则、依赖与 validator/artifact 映射
- Depends on: 各包 manifest（`workflow-packages/*/package.yaml`）
- Used by: Query 层与 Activation 构建 bundle/validator 计划（`internal/query/query.go`, `internal/activation/activation.go`）

**路由查询层：**
- Purpose: 在索引中筛选候选并形成 must_include/attach 列表
- Location: `internal/query/query.go`
- Contains: 评分、候选筛选、可见性/activation_mode 门槛
- Depends on: SQLite 索引（`blueprint/index/blueprint.db`）与 Registry（`internal/registry/registry.go`）
- Used by: Activation 与 MCP/CLI（`internal/activation/activation.go`, `cmd/agent-pack-mcp/main.go`）

**Activation 编排层：**
- Purpose: 执行 route -> bundle -> validation -> result
- Location: `internal/activation/activation.go`
- Contains: 请求解析、候选路由、context bundle、validation plan、handoff
- Depends on: Query 层与 Validator 层
- Used by: CLI/MCP activate 命令（`cmd/agent-pack-mcp/main.go`）

**Validation 层：**
- Purpose: 对 artifacts 与包依赖做校验与修复建议
- Location: `internal/validator/`
- Contains: core validator、domain validator、执行器
- Depends on: Activation 的 ValidationPlan 与 ExecutionInput（`internal/activation/activation.go`, `internal/validator/types.go`）
- Used by: Activation（`internal/activation/activation.go`）

## Data Flow

**Blueprint 编译与路由闭环：**

1. 编译 Blueprint：`cmd/agent-pack-mcp/main.go` 调用 `compiler.Compile`（`internal/compiler/compiler.go`）
2. 生成索引：SQLite 写入 `blueprint/index/blueprint.db`（`internal/compiler/compiler.go`）
3. 路由查询：`query.RouteQuery` 读取索引筛选候选（`internal/query/query.go`）
4. 构建 bundle：`query.BuildContextBundle` 汇总 required/may_include/children（`internal/query/query.go`）
5. Activation 汇聚：`activation.Execute` 组装 artifacts、validators、handoff（`internal/activation/activation.go`）
6. 运行校验：`validator.Run` 调用注册 validators（`internal/validator/runner.go`, `internal/validator/registry.go`）
7. 返回结果：`model.ActivationResult`（`internal/model/model.go`）

**State Management:**
- SQLite 索引是路由状态单一真相源（`blueprint/index/blueprint.db`）
- 注册表是 package 关系真相源（`workflow-packages/registry.json`）

## Key Abstractions

**Blueprint Node（frontmatter）：**
- Purpose: 定义路由节点与可见性/激活模式
- Examples: `blueprint/L1/wxt/manifest.md`, `blueprint/L1/security/permissions.md`
- Pattern: 路径与 `id` 严格一致（`internal/compiler/compiler.go`）

**Registry PackageEntry：**
- Purpose: 统一治理包名、领域、validators 与 artifacts
- Examples: `workflow-packages/registry.json`, `internal/registry/registry.go`
- Pattern: workflow/capability 命名规则验证（`internal/registry/registry.go`）

**ContextBundle：**
- Purpose: 产出最小可执行上下文（main/required/children/deferred）
- Examples: `internal/model/model.go`, `internal/query/query.go`
- Pattern: L3 节点进入 Deferred（`internal/query/query.go`）

**ValidationPlan/Result：**
- Purpose: 统一校验计划与结果结构
- Examples: `internal/model/model.go`, `internal/activation/activation.go`, `internal/validator/*.go`
- Pattern: core + domain validators 组合（`internal/validator/registry.go`）

## Entry Points

**CLI/MCP Binary:**
- Location: `cmd/agent-pack-mcp/main.go`
- Triggers: `compile`, `route_query`, `read_node`, `build_context_bundle`, `expand_node`, `activate`, `mcp`
- Responsibilities: 统一对外入口、调用 compiler/query/activation 组件

**MCP Server:**
- Location: `cmd/agent-pack-mcp/main.go`
- Triggers: `mcp` 子命令（stdio）
- Responsibilities: 暴露 route/build/read/expand/rebuild 能力（对应 `internal/query/*` 与 `internal/compiler/compiler.go`）

## Error Handling

**Strategy:** 返回结构化错误或降级结果，Activation 以状态表达失败/部分完成。

**Patterns:**
- 编译错误写入报告文件（`internal/compiler/compiler.go` -> `blueprint/index/validation-report.json`）
- Activation 请求无效直接返回 `failed`（`internal/activation/activation.go`）

## Cross-Cutting Concerns

**Logging:** 以 CLI 输出为主（`cmd/agent-pack-mcp/main.go`）
**Validation:** 编译期验证 frontmatter，运行期验证 artifacts（`internal/compiler/compiler.go`, `internal/validator/*`）
**Authentication:** Not applicable（无鉴权层）

---

*Architecture analysis: 2026-03-16*
