# Phase 5: Domain Expansion Pilot - Context

**Gathered:** 2026-03-18
**Status:** Ready for planning

<domain>
## Phase Boundary

在既有护栏内引入第二主域试点，并验证不会破坏 WXT 样板契约：
- 第二主域需按既有 registry/routing/validation/runtime 契约准入
- 既有 WXT 样板路径必须保持行为一致并通过回归
- 产出可复用的新主域准入清单（命名、路由、验证、交付契约、runtime 证据）

本阶段仅澄清并落地“如何准入与如何验收”，不扩展新产品能力。

</domain>

<decisions>
## Implementation Decisions

### 主域试点对象
- 第二主域试点方向锁定为 `monorepo-oss-governance`（以 monorepo 架构与 OSS 规范治理运营为语义主轴）。
- 试点边界采用“单子域闭环”，只要求 1 条可执行主链路打通（route -> bundle -> validate -> runtime），不做多子域并行扩展。
- 与既有能力线关系采用“同名能力线保留”：保留当前 attach-only 能力线语义，不做一次性替换迁移。
- 失败回退策略采用“特性开关回退”，确保可快速恢复到仅 WXT 主域稳定路径。

### 准入验收闸门（DOMN-01）
- DOMN-01 采用“契约四闸门”通过标准：命名治理、路由契约、交付契约、验证/回写制度必须同时通过。
- P0 阻断采用“四类全阻断”：主域竞争越界、attach-only 破坏、WXT 回归失败、runtime ledger/run_id 证据链断裂。
- 非阻断问题（warn）允许继续，但必须带整改项与 deadline，并写入 runtime ledger 留痕。
- 闸门通过判定要求 machine-readable 与 human-readable 结果同时可审计。

### WXT 回归口径（DOMN-02）
- 回归范围锁定“主链路三件套”：`route_query`、`build_context_bundle`、`activate`。
- 行为不变判定采用“关键字段锁定”，至少覆盖 status、main pack、required_packs、decision_basis、contract decisions。
- 新主域引入后，冲突场景纳入 P0（抗干扰优先）：新增主域触发词不得抢占原本应落入 WXT 的请求。
- 放行策略锁定为“WXT 失败即阻断”：任一 WXT 主链路阻断失败均不允许 Phase 5 通过。

### 准入清单产物
- 在 phase 内采用“主清单”模式：`05-CONTEXT.md` 锁定决策，同时补充一份可执行引用的准入清单文档。
- 清单结构采用“五段式”：命名治理、路由契约、交付契约、验证回归、runtime 证据链。
- 通过证据最小集合为“命令 + 产物 + ledger”：必须可回溯执行命令、关键输出摘要、`run_id` 与 ledger refs。
- 非 P0 自主决策按“决策项逐条回写”，禁止只写抽象口号。

### Claude's Discretion
- 非 P0 决策默认授权 Claude 按当前阶段现状与工程最佳实践做出，并在上下文文档中逐条回写。
- 仅在以下条件升级为用户决策：P0 风险、不可逆影响、与既有锁定决策冲突。
- warn 类 deadline 的默认值、检查命令编排细节、清单执行顺序可由 Claude 在规划阶段细化，但不得违背已锁定闸门与阻断规则。

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Phase scope and requirements
- `.planning/ROADMAP.md` — Phase 5 目标、DOMN-01/DOMN-02 与成功标准边界。
- `.planning/REQUIREMENTS.md` — DOMN-01/DOMN-02 的需求映射与状态口径。
- `.planning/PROJECT.md` — 扩域仍需遵循“最小且完整”交付与 bounded-context 核心约束。

### Governance and contracts
- `.planning/phases/02-routing-governance/02-CONTEXT.md` — candidate-space-first、attach-only、可解释路由等已锁定规则。
- `.planning/phases/03-contracted-delivery/03-CONTEXT.md` — include/exclude rationale 与“最小且完整”交付契约。
- `.planning/phases/04-validation-runtime-governance/04-CONTEXT.md` — validation/run_id/runtime ledger 制度化约束。

### Existing implementation anchors
- `workflow-packages/registry.json` — package 身份、domain/category、activation_mode 与 required_packs 真相源。
- `internal/query/query.go` — 路由主流程、candidate-space 过滤、decision_basis/capability decisions 产出。
- `internal/activation/activation.go` — validation plan/run_id/runtime ledger 写入模式与 evidence 链路。
- `tests/m1_minimal_test.go` — 路由治理与 attach-only 关键回归样例。
- `tests/m2_wxt_manifest_test.go` — WXT 样板链路与 activation 结果回归样例。

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `workflow-packages/registry.json`: 已有 1 条主域线（`wxt-manifest`）+ 2 条 capability 线（`security-permissions`、`release-store-review`），可直接复用为“主域扩展前后对照基线”。
- `internal/query/query.go`: 已具备 candidate-space-first、attach-only 隔离、decision trace 输出，可复用为 DOMN-01 路由闸门实现基础。
- `internal/activation/activation.go`: 已具备 validation plan 聚合、`run_id`、runtime ledger immediate/batch_finalize，可复用 DOMN-01/DOMN-02 的证据链要求。
- `tests/m1_minimal_test.go` 与 `tests/m2_wxt_manifest_test.go`: 已有 WXT 与 capability 相关回归，可扩展为 Phase 5 抗干扰与不破坏验证矩阵。

### Established Patterns
- 路由治理模式：先候选空间过滤，再评分与稳定 tie-break；attach-only 不进入主域竞争。
- 契约交付模式：context bundle 以 include/exclude rationale 为 machine-readable 契约。
- 验证治理模式：`passed/warned/failed` 三态 + run 级 evidence refs + runtime ledger append-only。
- 文档治理模式：phase 决策先锁定到 CONTEXT，再由 research/planner 转为可执行方案。

### Integration Points
- 新主域试点主要接入点：`blueprint/L0|L1` 新增域节点、`workflow-packages/*` 新 package、`workflow-packages/registry.json` 注册扩展。
- 行为验证接入点：`query.RouteQuery` / `query.BuildContextBundle` / `activation.Execute` 组合回归。
- 运行态回写接入点：`ActivationResult.ValidationRunID` 与 `RuntimeLedgerEntries`（trace/run 关联）。

</code_context>

<specifics>
## Specific Ideas

- 用户明确要求：结合业务需求与当前阶段现状，由 Claude 按开发领域工程师最佳实践自主决定非 P0 项。
- 用户明确锁定：本阶段聚焦“怎么准入与怎么验证不破坏”，不讨论新增能力扩展。
- 第二主域语义方向由用户指定为“monorepo 架构 + OSS 规范治理运营横切”。

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 05-domain-expansion-pilot*
*Context gathered: 2026-03-18*
