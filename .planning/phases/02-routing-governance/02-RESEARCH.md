# Phase 02: Routing Governance - Research

**Researched:** 2026-03-17
**Domain:** Go 路由治理（candidate-space-first、attach-only、可解释路由、canonical 缺失显式失败）
**Confidence:** HIGH

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

### 路由解释信息粒度（面向 MCP 调用方开发 Agent）
- 默认反馈采用“极简必要信息”策略，优先提供可被 agent 稳定消费的最小字段。
- 采用双层模型：默认返回极简摘要；需要时可扩展返回结构化明细。
- 返回风格为“可操作”：失败/partial 场景提供简短原因 + 下一步建议。
- 响应中预留 `docs_ref` 扩展位（当前可为空），为后续官方文档链接能力留接口。

### 无 canonical 映射时的失败语义
- 默认严格模式：无 canonical 映射时不允许隐式回退，直接显式失败（hard fail）。
- 使用稳定机器错误码（示例：`ROUTE_CANONICAL_MISSING`）+ 简短 message。
- 默认不回传完整候选列表，避免暴露过多内部治理细节；仅返回必要上下文。
- 返回固定短建议（检查 registry canonical 映射）并支持 `docs_ref` 指向后续文档。

### capability 附挂可见性与可审计性
- 必须返回“为何附挂/为何未附挂”的极简解释（原因码 + 短说明）。
- 返回规则标识（如 `BR-02`、`BR-03`）作为依据，并可选 `docs_ref`。
- 默认输出“最终附挂列表 + 每项极简原因”，不默认输出冗长全链路。
- capability 被拒绝附挂时，返回拒绝原因码 + 1 条下一步建议。

### 候选冲突与稳定决策策略
- 同分冲突采用稳定 tie-break，保证同输入可复现（禁止随机选择）。
- tie-break 默认优先级：canonical 命中优先 > 明确主域匹配 > 规则优先级 > 名称字典序。
- 若无法稳定判定，返回 explicit error/partial，不做猜测性回退。
- 返回简短可复现标记（如 `decision_trace_id` 或 `decision_basis`）。

### Claude's Discretion
- 极简响应与详细响应的具体字段命名（在不违背“默认极简”前提下）。
- `decision_trace_id` 与 `decision_basis` 的具体数据结构和编码形式。
- 规则标识与错误码在响应中的嵌套层级设计。

### Deferred Ideas (OUT OF SCOPE)
- 官方文档链接体系的完整建设（文档站点结构、链接路由策略、版本化说明）作为后续增强事项记录；当前阶段仅预留 `docs_ref` 扩展位。
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|-----------------|
| ROUT-01 | Maintainer can enforce candidate-space-first routing (scope/mode filter before scoring) | 把 `RouteQuery` 重构为“先候选空间过滤、后评分排序”；增加测试断言 attach-only/跨域节点不会进入 primary 候选池。 |
| ROUT-02 | Maintainer can guarantee capability is attach-only after primary domain selection | 明确两阶段决策：先主域与主包，再 capability attach；输出 attach/deny 原因码并写入结果。 |
| ROUT-03 | Maintainer can explain why main domain/main package/capabilities were selected | 定义最小可解释结果结构（主决策 + capability 附挂理由 + 规则依据 + 可复现标记）。 |
| ROUT-04 | Maintainer can return explicit error/partial when no canonical registry mapping exists (no silent fallback) | 删除 target_pack 分支中的隐式 fallback；canonical 缺失返回稳定错误码与短建议，不泄露全量候选细节。 |
</phase_requirements>

## Summary

当前路由主链路位于 `internal/query/query.go`，已具备 attach-only 的部分护栏，但与本 phase 的锁定决策仍有关键差距：第一，当前实现是“先评分再按候选空间归类”，不满足 ROUT-01 的 candidate-space-first 严格语义；第二，`target_pack` 命中 registry 但 canonical node 缺失/不可读时，当前会构造 `registry fallback` 候选并继续，违反 ROUT-04 的“无 canonical 映射必须显式失败”；第三，解释信息目前仅是 `reason []string`，缺少规则标识、原因码、decision trace、attach 拒绝理由，无法满足 ROUT-03。

从现有代码与测试可见，attach-only 基线已存在（`domainNodeAllowedInGlobal`、`workflowNodeAllowedInDomain`、`capabilityAttachAllowed`），且 `tests/m1_minimal_test.go` 已覆盖部分行为。但输出契约仍偏“内部调试字段”，不是“面向 MCP 调用方 agent 的稳定最小可解释协议”。Phase 2 的核心不在加功能，而在“固化行为+固化输出语义+固化失败语义”。

建议本 phase 采用“路由治理内核 + 输出适配层 + 回归测试矩阵”三件套：先在 query/activation 层落地 deterministic 两阶段决策与显式错误，再在模型层定义极简可解释输出（含可扩展 detail），最后用正/负例测试锁死回退与冲突边界。

**Primary recommendation:** 先重构 `RouteQuery` 为严格两阶段 candidate-space-first 决策并移除 canonical fallback，再统一 route/activate 输出为“默认极简、可扩展明细”的可解释契约并用测试锁死。

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Go | 1.25 | 实现路由治理逻辑 | 仓库主语言与既有测试基线 |
| `database/sql` + `modernc.org/sqlite` | std + v1.38.2 | 节点与元数据查询 | 现有 blueprint index 查询栈 |
| `internal/registry` | in-repo | canonical/package 身份真相源 | BR-04/BR-10 已冻结为必须依赖 |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| `sort` | std | 候选排序与 tie-break | 同分冲突稳定判定 |
| Go `testing` | std | 路由行为与错误语义回归 | 锁定 ROUT-01~04 行为 |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| registry canonical 强约束 | 运行时 heuristic fallback | 会造成不可解释漂移，直接违背 ROUT-04 |
| 最小+可扩展解释结构 | 仅返回自由文本 reason | 可读但不稳定，MCP agent 难做程序化消费 |
| 两阶段主域->capability | capability 与 workflow 同池竞争 | 破坏 attach-only 治理边界 |

**Installation:**
```bash
# 本 phase 无新增依赖；保持现有 go.mod
go test ./...
```

## Architecture Patterns

### Recommended Project Structure
```text
internal/
├── query/
│   └── query.go                    # 路由治理核心（两阶段决策）
├── activation/
│   └── activation.go               # route 结果消费与 status 归并
├── model/
│   └── model.go                    # RouteResult/错误与解释结构
└── registry/
    └── registry.go                 # canonical 真相校验
tests/
└── m1_minimal_test.go              # 路由治理回归（可拆分 m2 路由专项）
```

### Pattern 1: Candidate-Space-First Pipeline
**What:** 先按 `level + visibility_scope + activation_mode + domain` 过滤合法候选，再对合法候选评分。
**When to use:** 所有 route_query 路径（含 target_domain 场景）。
**Example:**
```go
// Source: internal/query/query.go (phase target pattern)
allowed := filterByCandidateSpace(nodes, level, activeDomain)
scored := scoreAllowedCandidates(task, allowed, meta, evidence)
sorted := stableTieBreak(scored)
```

### Pattern 2: Two-Stage Decision (Primary then Attach)
**What:** 阶段 1 只选主域/主包；阶段 2 仅在 primary 已确认后评估 capability attach。
**When to use:** L1 workflow 路由与 activation 主流程。
**Example:**
```go
primary := pickPrimaryWorkflow(workflowCandidates)
capability := evaluateAttachOnlyCaps(capCandidates, primary.Domain)
```

### Pattern 3: Explainable Minimal Contract
**What:** 默认返回极简机器可读字段（原因码/规则码/短建议），可选 detail 扩展。
**When to use:** route_query 与 activation 对外输出。
**Example:**
```json
{
  "status": "partial",
  "decision_basis": "canonical>domain>rule>lexicographic",
  "decision_trace_id": "rt-20260317-abc123",
  "main": {"pack": "wxt-manifest", "reason_code": "PRIMARY_DOMAIN_MATCH", "rule_ref": "BR-02"},
  "capabilities": [
    {"pack": "security-permissions", "attach": true, "reason_code": "ATTACH_ALLOWED", "rule_ref": "BR-03"}
  ],
  "next_action": "检查 registry canonical 映射",
  "docs_ref": ""
}
```

### Pattern 4: Canonical Missing = Explicit Error
**What:** `target_pack` 已给定但 registry canonical 无法映射到可路由节点时，直接失败，不返回 fallback candidate。
**When to use:** `RouteQuery` 的 target_pack 快路径。
**Example:**
```go
if !canonicalUsable(entry) {
    return RouteResult{Status: "failed", ErrorCode: "ROUTE_CANONICAL_MISSING"}, nil
}
```

### Anti-Patterns to Avoid
- **先评分再过滤候选空间：** 评分阶段已引入非法候选，解释链路会污染。
- **registry fallback 伪成功：** canonical 不可用却返回候选，后续 activate 行为不可预测。
- **只给自由文本 reason：** 缺少稳定错误码/规则码，无法写强断言测试与调用方策略。
- **随机或隐式 tie-break：** 同输入不同结果会破坏可复现与审计。

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| package 身份识别 | 从 node id/文件名猜 canonical | `registry.FindByName/FindByNode` | registry 是唯一身份真相源 |
| 错误表达 | 拼接临时字符串错误 | 稳定 `error_code + message + next_action` | 便于 agent 稳定消费与测试断言 |
| attach 解释 | 手写散乱 reason 文本 | 规则码驱动（BR-02/BR-03）+短说明 | 保证跨版本可读且可审计 |
| 冲突处理 | “分数一样随便选” | 固定 tie-break 链与 trace 标记 | 同输入可复现 |

**Key insight:** 路由治理难点不是“选出一个候选”，而是“可复现地选、可解释地选、失败时显式地停”。

## Common Pitfalls

### Pitfall 1: Candidate-space-first 被实现顺序悄悄破坏
**What goes wrong:** 非法候选先参与评分，最终虽然没入选，但 reason/排序已受污染。
**Why it happens:** 过滤与评分混在单循环中。
**How to avoid:** 拆成“过滤函数 + 评分函数”两个阶段并单测。
**Warning signs:** 结果里出现 `attach-only` 候选的竞争痕迹或 reason。

### Pitfall 2: target_pack 分支残留 fallback
**What goes wrong:** registry canonical 缺失时返回“看似成功”的候选。
**Why it happens:** 为兼容旧行为保留 fallback。
**How to avoid:** 统一 hard fail + `ROUTE_CANONICAL_MISSING`。
**Warning signs:** 响应 reason 出现 `registry fallback`。

### Pitfall 3: 解释字段不稳定
**What goes wrong:** 上游 agent 无法依赖输出做自动决策。
**Why it happens:** 解释字段只面向人看，不面向程序消费。
**How to avoid:** 固定最小字段（reason_code/rule_ref/next_action/docs_ref）。
**Warning signs:** 相同错误场景返回文本措辞频繁变化。

### Pitfall 4: tie-break 不可审计
**What goes wrong:** 同分候选结果不可复现，引发 flaky 测试。
**Why it happens:** 仅按 score 排序或依赖 map/遍历顺序。
**How to avoid:** 明确优先级链并在结果中回传 `decision_basis`。
**Warning signs:** 同输入偶发命中不同 candidate。

## Code Examples

Verified patterns from repo/current implementation and official APIs:

### attach-only 不进入 primary 竞争
```go
// Source: internal/query/query.go
func workflowNodeAllowedInDomain(candidate nodeRecord, activeDomain string) bool {
    if candidate.Domain != activeDomain {
        return false
    }
    if candidate.ActivationMode == "attach-only" {
        return false
    }
    return candidate.VisibilityScope == "domain-scoped" || candidate.VisibilityScope == "global"
}
```

### 当前存在的 canonical fallback（Phase 2 需移除）
```go
// Source: internal/query/query.go
return model.RouteResult{
    Candidates: []model.RouteCandidate{{
        ID: entry.CanonicalBlueprintNode,
        Reason: []string{"target_pack match", "registry fallback"},
    }},
}, nil
```

### 稳定排序基础（同分按 ID 字典序）
```go
// Source: internal/query/query.go
sort.Slice(workflowCandidates, func(i, j int) bool {
    if workflowCandidates[i].Score == workflowCandidates[j].Score {
        return workflowCandidates[i].ID < workflowCandidates[j].ID
    }
    return workflowCandidates[i].Score > workflowCandidates[j].Score
})
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| target_pack canonical 不可用时容忍 fallback | canonical 缺失应 explicit fail/partial | Phase 2 目标 | 消除隐式回退，提升治理可信度 |
| reason 文本片段式解释 | 原因码+规则码+短建议+trace 的结构化解释 | Phase 2 目标 | MCP agent 可稳定消费 |
| score 与过滤耦合 | candidate-space-first 再评分 | Phase 2 目标 | 满足 ROUT-01 可验证语义 |

**Deprecated/outdated:**
- `registry fallback` reason 路径：与 ROUT-04 和 CONTEXT 锁定决策冲突。

## Open Questions

1. **route_query 的 explicit error 与 partial 如何分界？**
   - What we know: canonical 缺失必须 hard fail；domain 已指定但证据不足可 partial。
   - What's unclear: MCP 层是否统一返回 200+payload，还是部分场景走 MCP error。
   - Recommendation: 先在 payload 固化 `status + error_code`，MCP transport 维持兼容。

2. **`decision_trace_id` 生成策略是否需要跨进程唯一？**
   - What we know: 需要可复现标记用于审计与排障。
   - What's unclear: 仅请求内唯一还是全局唯一。
   - Recommendation: phase 内先做请求内稳定（基于 request_id + deterministic hash），避免引入新基础设施。

3. **解释 detail 的开关入口放在 route_query 还是 activate？**
   - What we know: 默认极简，按需扩展。
   - What's unclear: 是 `explain_level=minimal|detailed` 还是 `include_debug=true`。
   - Recommendation: 采用 `explain_level`，语义更清晰且兼容未来分级。

## Sources

### Primary (HIGH confidence)
- `D:/dev/projects/the-agent-packs/.planning/phases/02-routing-governance/02-CONTEXT.md` - 锁定决策与阶段边界
- `D:/dev/projects/the-agent-packs/.planning/REQUIREMENTS.md` - ROUT-01~04 定义
- `D:/dev/projects/the-agent-packs/internal/query/query.go` - 路由核心现状与差距
- `D:/dev/projects/the-agent-packs/internal/activation/activation.go` - activate 状态语义与输出现状
- `D:/dev/projects/the-agent-packs/internal/model/model.go` - RouteResult/ActivationResult 数据结构基线
- `D:/dev/projects/the-agent-packs/internal/registry/registry.go` - canonical 与 attach-only 注册表约束
- `D:/dev/projects/the-agent-packs/tests/m1_minimal_test.go` - 路由主行为测试现状
- `D:/dev/projects/the-agent-packs/tests/m2_registry_test.go` - registry/attach-only 验证现状
- `D:/dev/projects/the-agent-packs/.planning/ROADMAP.md` - Phase 2 成功标准

### Secondary (MEDIUM confidence)
- https://pkg.go.dev/sort - 排序与可复现 tie-break 参考
- https://pkg.go.dev/encoding/json - 输出结构化响应契约参考

### Tertiary (LOW confidence)
- None

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - 完全基于仓库既有栈与依赖。
- Architecture: HIGH - 直接映射 ROUT-01~04 + CONTEXT 锁定决策到代码差距。
- Pitfalls: HIGH - 已由现有实现中的可见路径（如 registry fallback）直接验证。

**Research date:** 2026-03-17
**Valid until:** 2026-04-16
