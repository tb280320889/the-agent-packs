# Blueprint Query MCP（M1）

本目录定义 M1 最薄 Query MCP 的接口与职责边界。

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
- `validate_blueprint_graph`（二阶段）
- `rebuild_index`（二阶段）

### Prompts
- `route-task`
- `expand-subdomain`
- `debug-validator-failure`

## 最薄实现边界
- route_query 只做 L0/L1 路由
- build_context_bundle 输出 main/required/execution_children/deferred
- 默认只输出摘要
