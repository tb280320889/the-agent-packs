# Phase 05: Domain Expansion Pilot - Research

**Researched:** 2026-03-18
**Domain:** 第二主域受控准入（monorepo-oss-governance）
**Confidence:** HIGH

## Summary

Phase 5 的最优实现路径是 **“增量扩域，不改治理骨架”**：

1. 在 blueprint 新增 L0/L1 的 monorepo 主域节点（不改既有 WXT 节点）
2. 新增 `workflow-packages/monorepo-oss-governance/package.yaml`
3. 在 `workflow-packages/registry.json` 注册新 workflow 包（保留既有 attach-only capability 线不迁移）
4. 用现有 `RouteQuery -> BuildContextBundle -> Execute` 链路打通试点
5. 增加一组 M5 回归：新域可路由 + WXT 不回归 + 冲突词不抢占 + feature flag 可回退

这条路径满足 DOMN-01/DOMN-02，且符合 Phase 5 锁定决策（单子域闭环、四闸门、WXT 失败即阻断、feature switch 回退）。

## Locked Decisions Mapping

- **主域试点对象锁定：** `monorepo-oss-governance`
- **闭环范围锁定：** 仅 1 条主链路（route -> bundle -> validate -> runtime）
- **能力线策略锁定：** 保留当前 attach-only 能力线语义
- **P0 阻断锁定：** 主域越界、attach-only 破坏、WXT 回归失败、run_id/ledger 证据链断裂
- **回退策略锁定：** feature switch 回退

## Implementation Recommendations

### A. Domain modeling

- `domain` 采用 `monorepo`
- `subdomain` 采用 `oss-governance`
- canonical package 名称固定：`monorepo-oss-governance`

原因：符合当前 registry 的 `workflow = domain-subdomain` 规则，且保留用户指定语义方向。

### B. Feature switch rollback

- 在 `internal/query/query.go` 增加开关函数（例如读取 `DOMAIN_MONOREPO_ENABLED`）
- 开关关闭时，L0/L1 的 monorepo 节点不进入候选空间
- 默认开启（便于 Phase 5 验收），关闭可快速回退至仅 WXT 稳定路径

### C. Regression strategy (DOMN-02)

新增 `tests/m5_domain_expansion_test.go`，最少覆盖：

1. monorepo 任务可选中新主域主包
2. 新域进入后 WXT 样板主链路（route/build_context_bundle/activate）行为保持
3. 冲突词场景（同时包含 monorepo + extension 线索）不会错误抢占 WXT 既有路径
4. 关闭 feature switch 后新域不可路由，WXT 仍正常

## Validation Architecture

### Gate 1 — 命名治理

- `registry.Validate` 可通过
- 新 package 满足 `workflow` 命名规则（`domain-subdomain`）
- `canonical_blueprint_node` 非空且可读

### Gate 2 — 路由契约

- 路由保持 candidate-space-first
- attach-only 不进入 primary candidates
- 冲突场景保留稳定 tie-break，可输出 `decision_basis`

### Gate 3 — 交付契约

- `BuildContextBundle` 继续输出 include/exclude rationale
- 不混入无关域节点

### Gate 4 — 验证/回写制度

- `activation.Execute` 输出 validation run_id
- `RuntimeLedger` 维持 append-only（含 record_type/version/is_current）

## Risks & Mitigations

1. **风险：** 新域触发词污染 WXT 既有路径
   - **缓解：** 增加冲突回归测试，WXT 失败即阻断

2. **风险：** 仅加 registry，不加 blueprint，造成 canonical 空映射
   - **缓解：** 计划中把 blueprint + package + registry 作为同一任务原子落地

3. **风险：** 开关策略默认关闭导致 Phase 5 无法验收
   - **缓解：** 默认开启，保留显式关闭路径用于回退

## Recommended Verification Commands

- `go test ./... -run "TestM2Registry.*" -count=1`
- `go test ./... -run "TestM2Activation.*" -count=1`
- `go test ./... -run "TestM5DomainExpansion.*" -count=1`
- `go test ./... -run "TestM4Validation.*" -count=1`

## Output

研究结论支持直接进入 Phase 5 planning，无需新增外部依赖或账户配置。
