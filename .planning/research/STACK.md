# Technology Stack

**Project:** the-agent-packs（增强包生产者系统）  
**Researched:** 2026-03-16

## Recommended Stack

### Core Framework
| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| Go | 1.25.x（建议跟进 1.26 评估） | 核心编排、CLI/MCP 服务、validator 运行时 | 现有代码与测试已全面使用 Go；单二进制分发、静态部署成本低，适合“本地可执行 + CI 发布”模式。 |
| MCP Protocol | 2025-06-18 规范 | Host/Client/Server 能力协商、工具暴露、上下文交付协议化 | 项目目标是“受控交付上下文切片”，MCP 本身支持 capability 协商与工具化调用，天然契合。 |

### Database
| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| SQLite | 3.51.x（通过 modernc 驱动内置） | Blueprint 索引、路由查询、上下文裁剪中间层 | 本项目核心是“本地知识索引 + 快速查询 + 单文件工件”，SQLite 在单机/嵌入式场景成熟，运维负担低。 |
| modernc.org/sqlite | 现状 v1.38.2（建议升级到最新稳定线） | Go 无 CGO SQLite 驱动 | 现有代码已绑定；跨平台构建便利。需注意其官方强调 `modernc.org/libc` 版本配套约束。 |

### Infrastructure
| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| GitHub Actions | 当前仓库工作流 | 多平台构建与发布 | 仓库已有 release workflow；与 Go 二进制分发路径一致。 |
| Markdown + Frontmatter | 当前实现 | 知识节点源格式 | 面向知识维护者可读可审阅，便于 AIDP/Blueprint 协同。 |

### Supporting Libraries
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| gopkg.in/yaml.v3 | v3.x | 替换手工 YAML 解析（registry/package.yaml/frontmatter） | 只要涉及 YAML/frontmatter 解析与校验，就应优先使用，避免字符串拆分导致脆弱性。 |
| database/sql + SQLite pragma 管理 | 标准库 + 驱动参数 | 事务化索引重建、连接级约束（如 `_pragma=foreign_keys(1)`） | 索引重建、数据一致性保障、后续 schema 演化时必须启用。 |

## Alternatives Considered

| Category | Recommended | Alternative | Why Not |
|----------|-------------|-------------|---------|
| Runtime language | Go | Node.js/TypeScript | 本仓库已深度 Go 化；迁移成本高且对当前“单文件二进制”目标不增益。 |
| Database | SQLite | PostgreSQL | 当前主链路是本地索引与单实例路由，不需要分布式并发写；引入 PG 会显著增加运维复杂度。 |
| SQLite driver | modernc.org/sqlite（短期） | mattn/go-sqlite3 | 后者依赖 CGO，跨平台发布链更重；可作为“性能/兼容性兜底方案”保留评估。 |
| YAML parser | gopkg.in/yaml.v3 | 继续手工解析 | 手工解析已被现状审计明确为脆弱点，且扩展字段时易静默错误。 |

## Installation

```bash
# Core
go mod tidy

# Add robust YAML parsing (if not yet added)
go get gopkg.in/yaml.v3
```

## Decision Notes (Prescriptive)

1. **继续 Go + SQLite，不换栈。** 这是最小扰动且最符合“增强包生产系统”定位的路径。  
2. **先修解析与事务，再谈多域扩展。** 若 `package.yaml/frontmatter` 仍手工解析、索引重建无事务，多域接入会放大不确定性。  
3. **短期不引入 client/server 数据库。** 当前问题是语义治理与交付契约，不是数据库吞吐瓶颈。

## Sources

- Go Release History（go1.25/go1.26 发布时间与支持节奏，HIGH）: https://go.dev/doc/devel/release  
- MCP Specification 2025-06-18（协议能力与安全原则，HIGH）: https://modelcontextprotocol.io/specification/2025-06-18  
- SQLite “Appropriate Uses” （单机/嵌入式适配边界，HIGH）: https://www.sqlite.org/whentouse.html  
- modernc.org/sqlite 包文档（版本、平台支持、`libc` 约束，MEDIUM）: https://pkg.go.dev/modernc.org/sqlite  
- go-yaml v3 包文档（Decoder/KnownFields 等能力，MEDIUM）: https://pkg.go.dev/gopkg.in/yaml.v3
