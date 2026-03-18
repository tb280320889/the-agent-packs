# External Integrations

**Analysis Date:** 2026-03-16

## APIs & External Services

**MCP (Model Context Protocol):**
- 本地 MCP server（stdio）实现于 `cmd/agent-pack-mcp/main.go`
  - 工具：`route_query`、`read_node`、`build_context_bundle`、`expand_node`、`rebuild_index`
  - 资源：`blueprint://node/{id}` 等（见 `cmd/agent-pack-mcp/main.go` 与 `tools/query-mcp/README.md`）
  - Auth: 未检测到（未读取环境变量，未见外部认证配置）

## Data Storage

**Databases:**
- SQLite（本地文件）
  - Connection: CLI 参数 `--db`（默认 `blueprint/index/blueprint.db`）
  - Client: `modernc.org/sqlite`（`internal/compiler/compiler.go`、`internal/query/query.go`）

**File Storage:**
- 本地文件系统（Blueprint Markdown 与包资源均在仓库内）

**Caching:**
- None

## Authentication & Identity

**Auth Provider:**
- Not detected（未见 OAuth/第三方身份服务或本地认证逻辑）
  - Implementation: 无

## Monitoring & Observability

**Error Tracking:**
- None

**Logs:**
- 标准输出/错误输出（CLI 与 MCP 响应直接打印；`cmd/agent-pack-mcp/main.go`）

## CI/CD & Deployment

**Hosting:**
- GitHub Releases（构建产物上传 release 资产；`.github/workflows/release-agent-pack-mcp.yml`）

**CI Pipeline:**
- GitHub Actions（`.github/workflows/release-agent-pack-mcp.yml`）

## Environment Configuration

**Required env vars:**
- Not detected

**Secrets location:**
- Not detected（未发现 `.env*`，未读取环境变量）

## Webhooks & Callbacks

**Incoming:**
- None

**Outgoing:**
- None

---

*Integration audit: 2026-03-16*
