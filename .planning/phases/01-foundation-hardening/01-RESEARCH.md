# Phase 01: Foundation Hardening - Research

**Researched:** 2026-03-16
**Domain:** Go 解析/索引编译（YAML/frontmatter + SQLite 事务化重建）
**Confidence:** MEDIUM

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

### Claude's Discretion
- 解析错误策略：按类似产品与工程最佳实践决定严格性、失败与继续规则。
- frontmatter 支持范围：按常见语法优先且工程可维护性优先来裁定覆盖边界。
- 索引重建回滚策略：按可靠性与可恢复性优先的工程实践来决定回滚与状态标记方案。
- 结构化错误输出：按可调试、可追踪的最佳实践决定字段与粒度。

### Deferred Ideas (OUT OF SCOPE)
None — discussion stayed within phase scope.
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|-----------------|
| PARS-01 | Maintainer can parse `package.yaml` with standard YAML parser and fail fast on invalid fields | 使用 `gopkg.in/yaml.v3` 的 Decoder + KnownFields 严格模式；将 `package.yaml` 映射到结构体并校验未知字段。 |
| PARS-02 | Maintainer can parse Blueprint frontmatter with robust YAML semantics (not string-split) | 读取 frontmatter 块后用 YAML v3 解码，覆盖列表、引号、多行、null 等语法；保留字段行号用于结构化错误。 |
| INDX-01 | Maintainer can rebuild index transactionally so failed compile does not leave partial DB | 采用“临时 DB 构建 + 原子替换/回滚”或 DB 事务 + 完整写入后 commit 的模式，确保失败不破坏旧索引。 |
| INDX-02 | Maintainer can detect index build/report write failure with explicit error outcome | 编译流程分阶段返回结构化错误（解析/索引/报告写入），失败时显式失败状态并可测试覆盖。 |
</phase_requirements>

## Summary

当前实现中 `package.yaml` 与 frontmatter 解析都采用手写字符串拆分，难以支持引号、多行与复杂列表语法，也无法严格拒绝未知字段；同时索引重建直接在目标 SQLite 文件上执行删除/插入，没有事务或文件级回滚机制，一旦中途失败会留下部分写入状态。报告输出也仅在本地 `WriteFile` 失败时返回错误，但整体编译结果没有明确的结构化失败态。

建议 Phase 1 直接引入标准 YAML 解析（`gopkg.in/yaml.v3`），以 `Decoder.KnownFields(true)` 做严格模式，分别用于 `package.yaml` 和 frontmatter 的解码。索引重建建议采用“临时 DB 写入 + 原子替换 + 失败回滚”的文件级事务策略，或在单连接事务中完成 schema + 数据写入并在失败时回滚，确保失败不会覆盖旧索引。同时将编译流程拆分为解析、索引写入、报告写入三个阶段，统一返回结构化错误（包含阶段、路径、错误码、消息、可选行列号），满足 PARS/INDX 的显式失败语义。

**Primary recommendation:** 统一使用 `yaml.v3` 严格解码 + 事务化索引重建（临时 DB/回滚）+ 结构化错误输出，作为 Phase 1 的硬化主线。

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Go | 1.25 | 主语言 | 项目既定技术约束 | 
| modernc.org/sqlite | v1.38.2 | SQLite 驱动 | 纯 Go SQLite 驱动，项目既定依赖 | 
| gopkg.in/yaml.v3 | v3.0.1 | YAML 解析 | Go 生态标准 YAML 解析器，支持严格字段检查（KnownFields）与 YAML 1.2 语义 | 

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| database/sql (std) | go1.25 | 事务与连接管理 | 索引重建事务、失败回滚 |
| os (std) | go1.25 | 原子替换/重命名 | 临时 DB/报告文件写入后替换 |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| yaml.v3 严格解码 | 手写解析 | 维护成本高，无法覆盖多行/引号/列表语法，错误定位弱 |
| 临时 DB + 原子替换 | 直接在目标 DB 上写入 | 失败时容易留下部分索引或损坏状态 |

**Installation:**
```bash
go get gopkg.in/yaml.v3
```

## Architecture Patterns

### Recommended Project Structure
```
internal/
├── compiler/        # 编译与索引重建
├── registry/        # registry.json + package.yaml 解析
├── model/           # 数据模型
└── query/           # SQLite 查询
```

### Pattern 1: YAML 严格解码（KnownFields）
**What:** 使用 yaml.v3 的 Decoder 开启 KnownFields，拒绝未知字段并返回结构化错误。
**When to use:** 解析 `package.yaml` 与 blueprint frontmatter。
**Example:**
```go
// Source: https://pkg.go.dev/gopkg.in/yaml.v3#Decoder.KnownFields
dec := yaml.NewDecoder(reader)
dec.KnownFields(true)
if err := dec.Decode(&manifest); err != nil {
    return err
}
```

### Pattern 2: 临时 DB 构建 + 原子替换
**What:** 在同目录创建临时 DB 文件，完整写入后再替换目标索引；失败时保留旧索引或回滚到备份。
**When to use:** blueprint 索引重建（INDX-01）。
**Example:**
```go
// Source: https://pkg.go.dev/os#Rename
tmpPath := dbPath + ".tmp"
// 写入临时 DB（成功后）
if err := os.Rename(tmpPath, dbPath); err != nil {
    return err
}
```
> 注意：跨平台替换语义可能不同，必要时采用“旧文件 -> 备份 -> 新文件 -> 目标”的回滚路径。

### Pattern 3: 结构化错误输出
**What:** 解析/索引/报告阶段统一返回 `[]CompilerError`，包含 `phase`、`path`、`code`、`message`、`line/column`。
**When to use:** 编译失败、解析失败、报告写入失败（INDX-02）。
**Example:**
```go
type CompilerError struct {
    Phase   string `json:"phase"`
    Path    string `json:"path"`
    Code    string `json:"code"`
    Message string `json:"message"`
    Line    int    `json:"line,omitempty"`
    Column  int    `json:"column,omitempty"`
}
```

### Anti-Patterns to Avoid
- **手写 YAML 解析**：导致无法覆盖多行/引号/列表，错误定位弱，难以满足 PARS-01/02。
- **直接覆盖目标 DB**：失败时留下部分数据，违反 INDX-01 可恢复要求。
- **只返回 string error**：难以覆盖测试与定位，违反 INDX-02 明确失败态。

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| YAML 解析 | 字符串 split/trim 解析 | `gopkg.in/yaml.v3` | 支持标准 YAML 语义、错误信息更可靠 |
| 事务化重建 | 手动删除 + 逐条 insert | SQLite 事务或临时 DB 原子替换 | 失败回滚更可靠，不损坏旧索引 |

**Key insight:** YAML 与索引重建都存在隐性复杂度，标准库/生态工具更可靠且可测试。

## Common Pitfalls

### Pitfall 1: KnownFields 未开启导致未知字段被忽略
**What goes wrong:** `package.yaml`/frontmatter 中新增字段被静默忽略，无法及时发现配置错误。
**Why it happens:** 使用 `yaml.Unmarshal` 默认允许未知字段。
**How to avoid:** 使用 `yaml.Decoder` + `KnownFields(true)`。
**Warning signs:** 解析成功但运行态字段缺失，错误难以定位。

### Pitfall 2: 解析错误与编译错误混在一起
**What goes wrong:** 报告写入失败被误认为解析失败，导致修复方向错误。
**Why it happens:** 没有结构化错误阶段字段。
**How to avoid:** 统一 `phase` 字段并分阶段汇总。
**Warning signs:** 错误信息只有“compile failed”无具体阶段。

### Pitfall 3: DB 重建失败导致索引损坏
**What goes wrong:** 重建中断后索引库只写入部分数据。
**Why it happens:** 没有事务/文件级原子替换。
**How to avoid:** 临时 DB + 原子替换或事务回滚。
**Warning signs:** 失败后查询报缺失节点或 schema 为空。

## Code Examples

Verified patterns from official sources:

### YAML 严格解码
```go
// Source: https://pkg.go.dev/gopkg.in/yaml.v3#Decoder.KnownFields
dec := yaml.NewDecoder(reader)
dec.KnownFields(true)
if err := dec.Decode(&payload); err != nil {
    return err
}
```

### SQL 事务
```go
// Source: https://pkg.go.dev/database/sql#DB.BeginTx
tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
if err != nil {
    return err
}
if _, err := tx.ExecContext(ctx, stmt); err != nil {
    _ = tx.Rollback()
    return err
}
if err := tx.Commit(); err != nil {
    return err
}
```

### 文件级替换
```go
// Source: https://pkg.go.dev/os#Rename
if err := os.Rename(tmpPath, dbPath); err != nil {
    return err
}
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| 手写 frontmatter/YAML 解析 | yaml.v3 严格解码 + KnownFields | 2026-03 (Phase 1) | 解析更稳定、错误可追踪 |
| 原地 DB 删除/插入 | 临时 DB 构建 + 原子替换 | 2026-03 (Phase 1) | 失败不破坏旧索引 |

**Deprecated/outdated:**
- 手写 YAML 解析：无法满足 PARS-01/02 要求，容易静默忽略字段。

## Open Questions

1. **结构化错误字段的最小集合？**
   - What we know: 需要区分解析/索引/报告阶段，并包含路径与错误码。
   - What's unclear: 是否要求行号/列号、是否需要错误分级（warn/error）。
   - Recommendation: 先落地 `phase/path/code/message`，行列号可选。

2. **索引重建的回滚策略选型？**
   - What we know: 需要避免失败时覆盖旧索引。
   - What's unclear: 是否允许短暂无索引窗口（rename 过程）以及 Windows 替换语义要求。
   - Recommendation: 优先临时 DB + 备份/恢复策略；必要时用事务 + 单文件重建。

## Sources

### Primary (HIGH confidence)
- https://pkg.go.dev/gopkg.in/yaml.v3 — Decoder.KnownFields, Unmarshal 行为
- https://pkg.go.dev/database/sql — BeginTx/Tx 事务语义
- https://pkg.go.dev/os#Rename — 文件替换/重命名

### Secondary (MEDIUM confidence)
- 本仓库实现：`internal/compiler/compiler.go`、`internal/registry/registry.go`

### Tertiary (LOW confidence)
- 无

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - 依赖与官方文档明确
- Architecture: MEDIUM - 需结合当前实现迁移路径
- Pitfalls: MEDIUM - 基于现有代码与常见失败模式

**Research date:** 2026-03-16
**Valid until:** 2026-04-15
