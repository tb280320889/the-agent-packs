# Phase 3: Contracted Delivery - Context

**Gathered:** 2026-03-17
**Status:** Ready for planning

<domain>
## Phase Boundary

将“最小且完整”的消费侧上下文交付契约落地为可执行、可复验的检查机制：
- 交付内容仅包含目标域任务必需知识（不混入无关域）
- 交付对象必须携带 include/exclude rationale 且可追溯规则依据
- 建立可重复执行的契约检查用例（至少覆盖 WXT 样板与一个负例）

不包含新能力扩展（如新增主域、新增产品能力），仅在现有契约范围内固化交付与校验行为。

</domain>

<decisions>
## Implementation Decisions

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
  - 阻断级（P0）：跨域混入、required 缺失、规则不可追溯 → 检查失败（hard fail）
  - 非阻断级：说明文本不清晰但结构可追溯 → 警告（warning）
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

</decisions>

<specifics>
## Specific Ideas

- 用户授权：由 Claude 依据当前项目业务需求与工程最佳实践先行做实现决策。
- 决策升级规则：仅在 P0 + 上下文不足 + 后续影响重大时，再向用户提问。

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 03-contracted-delivery*
*Context gathered: 2026-03-17*
