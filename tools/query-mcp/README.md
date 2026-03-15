# Blueprint Query MCP（M1）

本目录定义 M1 最薄 Query MCP 的接口与职责边界。

## 运行方式（Go 模式B）
- 本能力由统一二进制 `agent-pack-mcp` 提供。
- CLI 子命令：
  - `go run ./cmd/agent-pack-mcp route_query ...`
  - `go run ./cmd/agent-pack-mcp read_node ...`
  - `go run ./cmd/agent-pack-mcp build_context_bundle ...`
  - `go run ./cmd/agent-pack-mcp expand_node ...`
- MCP server（stdio）：`go run ./cmd/agent-pack-mcp mcp --db blueprint/index/blueprint.db`

## 固定接口

### Resources
- `blueprint://node/{id}`
- `blueprint://children/{id}`
- `blueprint://required/{id}`
- `blueprint://bundle/{bundle_id}`

### Tools
- `route_query`
- `expand_node`
- `read_node`
- `build_context_bundle`
- `rebuild_index`
- `validate_blueprint_graph`（二阶段）

### Prompts
- `route-task`
- `expand-subdomain`
- `debug-validator-failure`

## 最薄实现边界
- route_query 只做 L0/L1 路由
- build_context_bundle 输出 main/required/execution_children/deferred
- 默认只输出摘要
