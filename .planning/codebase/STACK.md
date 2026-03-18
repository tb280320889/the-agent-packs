# Technology Stack

**Analysis Date:** 2026-03-16

## Languages

**Primary:**
- Go 1.25 - 代码实现集中于 `cmd/agent-pack-mcp/main.go`、`internal/` 与 `tests/` 下的 `.go` 文件

**Secondary:**
- Markdown - Blueprint 与文档资产存放在 `blueprint/`、`docs/`、`workflow-packages/**/README.md`
- JSON - 注册表与契约数据：`workflow-packages/registry.json`、`workflow-packages/wxt-manifest/contracts/*.json`、`fixtures/*.json`
- YAML - 包清单与自动化配置：`workflow-packages/**/package.yaml`、`.github/workflows/release-agent-pack-mcp.yml`

## Runtime

**Environment:**
- Go toolchain 1.25（由 `go.mod` 约束）

**Package Manager:**
- Go Modules
- Lockfile: `go.sum`（存在）

## Frameworks

**Core:**
- 无传统 Web 框架；以 Go 标准库 + 自定义模块构建 CLI/MCP 服务（入口 `cmd/agent-pack-mcp/main.go`）

**Testing:**
- Go `testing` 标准库（测试文件位于 `tests/*.go`）

**Build/Dev:**
- GitHub Actions 构建/发布：`.github/workflows/release-agent-pack-mcp.yml`

## Key Dependencies

**Critical:**
- `modernc.org/sqlite` v1.38.2 - SQLite 驱动，用于 Blueprint 索引读写（`internal/compiler/compiler.go`、`internal/query/query.go`）

**Infrastructure:**
- `github.com/google/uuid` v1.6.0（间接）- 依赖链引入（见 `go.mod`）
- `golang.org/x/exp` 等（间接）- 依赖链引入（见 `go.mod`）

## Configuration

**Environment:**
- 未发现 `.env` 或运行时环境变量读取逻辑（`internal/` 与 `cmd/` 代码中未检出 `os.Getenv`/`LookupEnv`）
- 运行参数通过 CLI flags 提供（`cmd/agent-pack-mcp/main.go`）

**Build:**
- Go 模块声明：`go.mod`
- 发布构建矩阵：`.github/workflows/release-agent-pack-mcp.yml`
- Blueprint 索引产物默认输出至 `blueprint/index/`（被 `.gitignore` 排除）

## Platform Requirements

**Development:**
- Go 1.25
- 本地 Blueprint 索引路径默认 `blueprint/index/blueprint.db`（由 CLI 传参）

**Production:**
- 发布为多平台二进制（Windows/Linux/macOS），由 GitHub Actions 构建并上传 release 资产（`.github/workflows/release-agent-pack-mcp.yml`）

---

*Stack analysis: 2026-03-16*
