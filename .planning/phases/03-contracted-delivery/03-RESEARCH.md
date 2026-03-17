# Phase 03: Contracted Delivery - Research

**Researched:** 2026-03-17
**Domain:** 上下文交付契约化（最小且完整、include/exclude rationale、可重复检查）
**Confidence:** HIGH

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

### 最小且完整判定口径
- 采用“双条件门禁”：**任务可完成性 + 规则可解释性** 同时满足才算“required”。
- “完整”最低要求为：覆盖主流程、常见失败场景、关键约束依据（不是只覆盖 happy path）。
- “最小”采用平衡策略：必要项 + 少量防误用上下文；禁止因“保险”混入跨域信息。
- 当“最小”与“完整”冲突时，默认优先“完整”，但必须在 rationale 中显式记录为何放宽最小化。

### Include / Exclude 理由表达
- rationale 必须同时支持机器检查与人工审阅：
  - machine-readable：稳定字段（reason_code、source_rule、scope、decision_basis）
  - human-readable：简短说明（为何包含/排除，若排除是否有风险）
- include 与 exclude 都必须给出依据，不允许仅记录 include。
- 每条 rationale 必须可回溯到规则来源（requirements/roadmap/core rules）。

### 负例契约边界
- 负例定义：出现无关域节点、越界 capability、缺失关键 required 节点、或 rationale 缺失/不可追溯。
- 结果语义采用分级：
  - 阻断级（P0）：跨域混入、required 缺失、规则不可追溯 -> 检查失败（hard fail）
  - 非阻断级：说明文本不清晰但结构可追溯 -> 警告（warning）
- Phase 3 至少包含 1 个明确负例并保证可重复触发与复验。

### 契约检查执行形态
- 检查应纳入可重复执行命令（测试流可直接调用），避免依赖会话人工判断。
- 输出应包含：通过/失败结论、失败项列表、对应规则映射、最小修复建议。
- 检查结果需可用于后续 phase（Validation & Runtime Governance）接力，而非一次性日志。

### Claude's Discretion
- 在不破坏以上锁定决策前提下，默认由 Claude 基于工程最佳实践自主决策具体实现细节。
- 仅当出现以下条件才升级为用户决策：
  1) P0 风险；
  2) 上下文不足导致无法安全决策；
  3) 该决策将对后续 phase 产生重大不可逆影响。

### Deferred Ideas (OUT OF SCOPE)
None — discussion stayed within phase scope.
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|-----------------|
| CONT-01 | Consumer agent can receive context bundle containing required domain knowledge only | 在 `BuildContextBundle` 结果上新增 contract 结构与域边界检查：required/main/must_include 仅允许目标域与合法 attach-only capability，禁止跨域节点泄露。 |
| CONT-02 | Consumer agent can inspect include/exclude rationale for delivered context | 为每个 include/exclude 节点生成 machine-readable + human-readable rationale，并强制 `source_rule` 可追溯到 BR/Requirement/Roadmap。 |
| CONT-03 | Maintainer can verify “minimal yet complete” delivery with repeatable checks | 落地可重复执行的契约测试矩阵（WXT 正例 + 至少 1 个负例），输出 hard fail / warning、规则映射和最小修复建议。 |
</phase_requirements>

## Summary

当前仓库在 Phase 2 已具备“可解释路由契约”基础（`reason_code/rule_ref/decision_basis`、canonical 缺失 hard-fail），但 Context Bundle 仍只有结构化数据（`main/required/execution_children/deferred`），尚未把“为什么包含、为什么排除、排除风险”固化为正式交付契约。换句话说，路由可解释性已经有了，交付可解释性还没闭环。

实现 Phase 3 的关键不是新增大量能力，而是把现有 Route + Bundle + Validator 三段串成一个“可验证交付合同”：一是定义 Contract Envelope（include/exclude rationale 与规则追溯）；二是定义检查器（最小且完整双条件门禁 + 分级失败语义）；三是把检查器接入可重复执行命令（Go tests/CLI），并沉淀正负例 fixture。这样才能真正满足 CONT-01/02/03。

从代码现状看，最优路径是“在现有模型上增量扩展，不重写主链路”：`internal/model` 增合同字段，`internal/query.BuildContextBundle` 产出 rationale 原材料，`internal/validator` 增 `validator-contract-delivery` 执行规则校验，`tests/` 新增 phase3 合同回归。该路径风险小、与 Phase 4 可自然衔接。

**Primary recommendation:** 先冻结“Contract Envelope + 规则映射表 + 失败分级语义”，再以 validator + tests 方式把 WXT 正例与负例一次性跑通，避免先写实现后补契约导致反复返工。

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Go | 1.25 | 契约构建与检查逻辑实现 | 仓库主语言，现有 activation/query/validator 均为 Go |
| `database/sql` + `modernc.org/sqlite` | std + v1.38.2 | 从 index 读取节点/边用于 bundle 构建 | 当前编译/查询链路已稳定使用 |
| in-repo `internal/query` / `internal/validator` | N/A | 交付对象构建 + 合同检查执行 | 现成扩展点，最小改动可达成 Phase 3 目标 |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| Go `testing` | std | 可重复契约检查（正例/负例） | 作为 CONT-03 的默认执行入口 |
| JSON 编码（std `encoding/json`） | std | 对外输出可机读契约对象 | CLI/MCP 输出 contract details 时使用 |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| 在现有 bundle 上增量加 contract 字段 | 独立新建 bundle 生成器 | 会复制逻辑并增加漂移风险 |
| validator 驱动契约检查 | 在 route/query 内嵌大量 if 校验 | 规则与执行耦合，难复验与扩展 |
| Go 测试矩阵（fixture） | 仅文档/手工 checklist | 不可重复，无法满足 CONT-03 |

**Installation:**
```bash
go test ./...
```

## Architecture Patterns

### Recommended Project Structure
```text
internal/
├── model/
│   └── model.go                     # ContextBundle/Contract 结构扩展
├── query/
│   └── query.go                     # bundle 构建 + include/exclude 原始依据
├── validator/
│   ├── registry.go                  # 注册 contract validator
│   └── contract_delivery.go         # Phase 3 契约检查器（新增）
tests/
├── m3_validation_closure_test.go    # 现有基线
└── m3_contracted_delivery_test.go   # Phase 3 契约正负例（建议新增）
```

### Pattern 1: Contract Envelope First
**What:** 先定义交付契约对象，再落实现逻辑。对象至少包含：`included[]`、`excluded[]`、`reason_code`、`source_rule`、`scope`、`decision_basis`、`human_note`。
**When to use:** 所有向消费侧输出 context bundle 的路径。
**Example:**
```go
type ContractDecision struct {
    NodeID         string `json:"node_id"`
    Action         string `json:"action"` // include|exclude
    ReasonCode     string `json:"reason_code"`
    SourceRule     string `json:"source_rule"`
    Scope          string `json:"scope"`
    DecisionBasis  string `json:"decision_basis"`
    HumanNote      string `json:"human_note"`
}
```

### Pattern 2: 双条件门禁（最小且完整）
**What:** required 判定必须同时满足“任务可完成性”与“规则可解释性”。
**When to use:** include 集合裁剪与 validator 检查。
**Example:**
```go
required := taskCompletable(node, task) && ruleExplainable(node, rationale)
```

### Pattern 3: 分级失败语义
**What:** P0（hard fail）与 warning 分离，且输出最小修复建议。
**When to use:** contract validator 输出与测试断言。
**Example:**
```json
{
  "status": "failed",
  "failures": [{"code": "CONTRACT_CROSS_DOMAIN_INCLUDED", "severity": "error", "rule_ref": "BR-05A"}],
  "suggestions": ["移除非目标域节点并补充 exclude rationale"]
}
```

### Pattern 4: 规则映射显式化
**What:** reason_code 不能裸用，必须映射到规则来源（BR/REQ/ROADMAP）。
**When to use:** 产出 rationale 与检查报告。
**Example:**
```go
var ruleMap = map[string]string{
    "INCLUDE_REQUIRED_WITH": "BR-05B",
    "EXCLUDE_OUT_OF_DOMAIN": "CONT-01",
}
```

### Anti-Patterns to Avoid
- **只记录 include 不记录 exclude：** 无法证明“最小性”，直接违背 CONT-02。
- **把所有 may_include 都当 required：** 会导致跨域污染与过量交付。
- **检查器只输出 pass/fail：** 缺少规则映射和修复建议，不可用于 Phase 4 接力。
- **将文档样例当自动验证替代：** `runtime/10-WXT-消费契约验证样例.md` 明确当前仅文档级，不足以满足 CONT-03。

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| 规则追溯 | 在字符串里硬编码“看起来像规则”的描述 | 统一 `reason_code -> source_rule` 映射表 | 防止理由漂移，便于机器校验 |
| 合同检查执行 | 临时脚本或人工 checklist | Go `testing` + validator runner | 可重复执行、可回归 |
| 节点域判断 | 依赖 node id 字符串猜测 | 使用 nodes 表字段 + registry 映射 | 降低误判跨域风险 |
| 负例构造 | 在线修改真实资产试错 | 独立 fixture + deterministic 测试 | 可稳定复现与复验 |

**Key insight:** Phase 3 的产物不是“更多上下文”，而是“更可证明的上下文”。

## Common Pitfalls

### Pitfall 1: 把 Route 可解释性误当 Delivery 可解释性
**What goes wrong:** 只复用 route 的 `reason_code`，却没有 include/exclude 级别理由。
**Why it happens:** 误以为 Phase 2 已覆盖 Phase 3。
**How to avoid:** 在 bundle 层新增 contract decisions，而不是只看 route result。
**Warning signs:** 输出里有主包原因，但没有排除项依据。

### Pitfall 2: “最小”被实现成“更少”而不是“够用”
**What goes wrong:** 删除关键 required 节点，任务无法完成。
**Why it happens:** 只做体积最小化，没有完整性门禁。
**How to avoid:** required 判定必须通过双条件门禁并加负例验证。
**Warning signs:** 测试通过率依赖 happy path，失败场景缺失。

### Pitfall 3: rationale 可读但不可机检
**What goes wrong:** 文本说明很好看，但缺少稳定字段。
**Why it happens:** 只按人工审阅写说明。
**How to avoid:** 强制 machine-readable 字段必填，human note 可选增强。
**Warning signs:** 同一场景不同运行生成不同措辞，无法稳定断言。

### Pitfall 4: 负例语义混乱
**What goes wrong:** P0 问题被当 warning 放行。
**Why it happens:** 缺少 severity policy 与规则分级表。
**How to avoid:** 在 validator 固化 error/warn 对应码与拦截策略。
**Warning signs:** 测试里只断言 `status != completed`，不校验错误码。

## Code Examples

Verified patterns from current repository:

### ContextBundle 当前结构（Phase 3 扩展点）
```go
// Source: internal/model/model.go
type ContextBundle struct {
    Main              *NodeSummary  `json:"main"`
    Required          []NodeSummary `json:"required"`
    ExecutionChildren []NodeSummary `json:"execution_children"`
    Deferred          []NodeSummary `json:"deferred"`
}
```

### Bundle 构建入口（应在此注入 include/exclude rationale）
```go
// Source: internal/query/query.go
func BuildContextBundle(db *sql.DB, mainNode string, includeRequired, includeMayInclude, includeChildren bool) (model.ContextBundle, error)
```

### Validator 执行框架（Phase 3 检查器挂载点）
```go
// Source: internal/validator/runner.go
func Run(plan model.ValidationPlan, input ExecutionInput) []model.ValidatorResult
```

### 现有文档级口径提示（需升级为自动化）
```text
// Source: docs/AIDP/runtime/10-WXT-消费契约验证样例.md
本样例当前是文档级验证口径，不等于代码层自动化测试已经全部具备
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| 文档主张“最小且完整” | 路由已可解释，但交付 contract 未结构化 | Phase 2 完成后 | 具备进入 Phase 3 基线，但 CONT-01/02/03 尚未闭环 |
| 交付验证偏人工阅读 | 已有 Go 测试框架与 validator runner 可复用 | 现状 | 可快速落地重复执行检查 |
| 负例边界未正式分级 | CONTEXT 已冻结 P0/warning 语义 | 本 phase 前置决策 | 只需实现，不再讨论标准 |

**Deprecated/outdated:**
- 仅依赖 `docs/AIDP/runtime/10-WXT-消费契约验证样例.md` 的文档检查路径，不足以满足 CONT-03。

## Open Questions

1. **Contract 数据是挂在 `ContextBundle` 还是 `ActivationResult` 双写？**
   - What we know: 消费侧需要 inspect（CONT-02），维护侧需要检查（CONT-03）。
   - What's unclear: 单一来源放在 bundle，还是在 activation 透传以便调用方直接取用。
   - Recommendation: 以 bundle 为真相源，activation 只透传摘要与检查结果，减少重复状态。

2. **exclude 集合的“候选全集”边界如何定义？**
   - What we know: 必须说明为什么不包含其他知识。
   - What's unclear: “其他知识”是同域候选，还是全库候选。
   - Recommendation: Phase 3 先限定为“同主域 + attach-only capability 候选空间”内的排除项，避免全库枚举噪音。

3. **负例 fixture 放在 tests 还是 workflow-packages 内部？**
   - What we know: 需要可重复触发与复验。
   - What's unclear: 最佳落位。
   - Recommendation: 规则级负例放 `tests/fixtures`；包语义相关负例放 `workflow-packages/wxt-manifest/fixtures/fail`，两者并存。

## Sources

### Primary (HIGH confidence)
- `D:/dev/projects/the-agent-packs/.planning/phases/03-contracted-delivery/03-CONTEXT.md` - 锁定决策、失败分级、输出要求
- `D:/dev/projects/the-agent-packs/.planning/REQUIREMENTS.md` - CONT-01/02/03 的正式定义
- `D:/dev/projects/the-agent-packs/.planning/ROADMAP.md` - Phase 3 成功标准
- `D:/dev/projects/the-agent-packs/internal/model/model.go` - 当前 ContextBundle/Route/Activation 数据模型
- `D:/dev/projects/the-agent-packs/internal/query/query.go` - BuildContextBundle 与 must_include 现状
- `D:/dev/projects/the-agent-packs/internal/activation/activation.go` - activation 输出链路与 validation plan 组装
- `D:/dev/projects/the-agent-packs/internal/validator/runner.go` - validator 执行框架
- `D:/dev/projects/the-agent-packs/internal/validator/registry.go` - validator 注册扩展点
- `D:/dev/projects/the-agent-packs/tests/m1_minimal_test.go` - route/bundle 基线测试
- `D:/dev/projects/the-agent-packs/tests/m3_validation_closure_test.go` - validation 闭环测试基线
- `D:/dev/projects/the-agent-packs/docs/AIDP/runtime/10-WXT-消费契约验证样例.md` - 当前文档级契约样例与限制

### Secondary (MEDIUM confidence)
- `https://pkg.go.dev/testing` - Go 可重复测试执行约定
- `https://pkg.go.dev/encoding/json` - 机读契约输出约定

### Tertiary (LOW confidence)
- None

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - 完全基于仓库现有技术栈与代码结构。
- Architecture: HIGH - 直接映射 CONTEXT 锁定决策到现有扩展点。
- Pitfalls: HIGH - 均可在当前实现与文档限制中定位到具体风险。

**Research date:** 2026-03-17
**Valid until:** 2026-04-16
