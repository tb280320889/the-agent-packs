# Context Snapshot: M1 Go模式B MCP迁移

> 弃用说明：本文件属于正式主线 `docs/00~52` 的历史快照，仅保留用于追溯；改造计划 v1 请改读 `docs/改造计划v1/context-snapshots/`。

## 1. 当前阶段
- 所属里程碑：M1
- 关联 bead：the-agent-packs-ld9
- 当前状态：completed

## 2. 当前事实
- 当前要解决的问题：将 M1 的 Python 脚本实现替换为 Go 的模式B二进制 MCP 运行面，确保下游不再依赖 Python/Node 运行时。
- 当前已完成内容：
  - 新增统一入口二进制 `agent-pack-mcp`，包含 `compile`、`route_query`、`read_node`、`build_context_bundle`、`expand_node`、`activate`、`mcp` 子命令。
  - 新增 Go 核心实现：`internal/compiler`、`internal/query`、`internal/activation`。
  - 删除旧 Python 实现与 Python 测试。
  - 迁移测试到 Go：`tests/m1_minimal_test.go`。
  - 更新依赖与 smoke 文档，确保命令与运行面一致。
- 当前尚未完成内容：无。

## 3. 已冻结对象
- 协议形状保持不变：route result、context bundle、activation result 的字段语义保持与 M1 冻结契约一致。
- 运行面冻结：M1 默认运行面为 Go 模式B（二进制 runner + MCP server），不再以 Python 作为主路径。

## 4. 当前输入
- 上游交付物：M1 最薄骨架（Blueprint + compiler/query/entry 语义）。
- 依赖文档：`docs/20-M1_Blueprint知识骨架与最薄入口_开发指导.md`、`docs/22-M1_上下文_Routing_Bundle_ActivationEntry.md`、`docs/23-M1_上下文_Compiler_SQLite_QueryMCP骨架.md`。
- 依赖 fixtures：`fixtures/activation-request.sample.json`、`fixtures/route-result.sample.json`、`fixtures/context-bundle.sample.json`、`fixtures/activation-result.sample.json`。

## 5. 当前输出
- 已新增文件：
  - `go.mod`
  - `go.sum`
  - `cmd/agent-pack-mcp/main.go`
  - `internal/model/model.go`
  - `internal/compiler/compiler.go`
  - `internal/query/query.go`
  - `internal/activation/activation.go`
  - `tests/m1_minimal_test.go`
  - `.github/workflows/release-agent-pack-mcp.yml`
  - `docs/context-snapshots/2026-03-15-m1-go-modeb-mcp-migration.md`
- 已删除文件：
  - `tools/compiler/compiler.py`
  - `tools/query-mcp/query_mcp.py`
  - `tools/activation-entry/activation_entry.py`
  - `tests/test_m1_minimal.py`
- 已更新文件：
  - `docs/运行依赖与降级策略.md`
  - `fixtures/smoke-run.md`
  - `fixtures/README.md`
  - `fixtures/context-bundle.sample.json`
  - `tools/compiler/README.md`
  - `tools/query-mcp/README.md`
  - `tools/activation-entry/README.md`
  - `docs/context-snapshots/2026-03-15-m1-blueprint-skeleton.md`
  - `docs/handoffs/the-agent-packs-695-handoff.md`
  - `.gitignore`

## 6. 风险与阻塞
- 风险：MCP server 目前为最薄 JSON-RPC/stdio 形态，后续若接入更多 MCP client，需要补更严格的协议兼容测试。
- 阻塞：无。
- 是否需要人工决策：否。

## 7. 下一步建议
- 建议下一个 bead：M2 pack 实现继续直接依赖 Go MCP 入口，不再写 Python 回退路径。
- 建议先执行的命令：`go test ./...`。
- 建议先阅读的文档：`docs/运行依赖与降级策略.md`、`fixtures/smoke-run.md`。
