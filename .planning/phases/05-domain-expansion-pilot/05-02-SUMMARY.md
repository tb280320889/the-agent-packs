---
phase: 05-domain-expansion-pilot
plan: 02
subsystem: testing
tags: [wxt, non-regression, onboarding-checklist, runtime-ledger]

requires:
  - phase: 05-domain-expansion-pilot-01
    provides: 第二主域准入资产与 feature switch 候选空间控制
provides:
  - DOMN-02 的 WXT 非回归矩阵（route/bundle/activate/conflict）
  - 新主域五段式准入清单与最小证据集合模板
affects: [future-domain-onboarding, DOMN-02, validation-governance]

tech-stack:
  added: []
  patterns: [wxt anti-steal regression, handoff-safe status compatibility]

key-files:
  created:
    - tests/m5_domain_expansion_regression_test.go
    - .planning/phases/05-domain-expansion-pilot/05-DOMAIN-ONBOARDING-CHECKLIST.md
  modified:
    - internal/activation/activation.go
    - internal/model/model.go

key-decisions:
  - "DOMN-02 回归以 WXT 三件套（route/build_context_bundle/activate）+ conflict anti-steal 为固定矩阵。"
  - "handoff 状态兼容策略限定在无 phase/plan 标记请求，避免破坏 Phase 4 的 phase-tagged completed 语义。"

patterns-established:
  - "扩域测试必须同时覆盖正向准入与既有样板抗抢占回归。"
  - "新主域准入清单统一包含 P0 阻断、warn deadline 与 runtime ledger 留痕要求。"

requirements-completed: [DOMN-01, DOMN-02]

duration: 6min
completed: 2026-03-18
---

# Phase 5 Plan 02: WXT Regression & Onboarding Checklist Summary

**固化 WXT 非回归矩阵并沉淀五段式准入清单，使第二主域接入可验证且可复制。**

## Performance

- **Duration:** 6 min
- **Started:** 2026-03-18T08:34:40Z
- **Completed:** 2026-03-18T08:38:43Z
- **Tasks:** 2
- **Files modified:** 4

## Accomplishments
- 新增 DOMN-02 回归测试矩阵，覆盖 WXT 路由不破坏、bundle 契约不破坏、activation trace 不破坏、冲突词防抢占。
- 输出 `05-DOMAIN-ONBOARDING-CHECKLIST.md`，固定命名治理/路由契约/交付契约/验证回归/runtime 证据链五段式模板。
- 自动修复 handoff 状态兼容性回归，确保既有 handoff 测试与 phase-tagged validation 测试同时通过。

## Task Commits

1. **Task 1: 增加 DOMN-02 回归矩阵** - `7a0783b` (fix)
2. **Task 2: 产出五段式新主域准入清单** - `e360c85` (docs)

## Files Created/Modified
- `tests/m5_domain_expansion_regression_test.go` - DOMN-02 WXT 非回归矩阵与冲突抗干扰断言。
- `.planning/phases/05-domain-expansion-pilot/05-DOMAIN-ONBOARDING-CHECKLIST.md` - 五段式准入模板与最小证据集合。
- `internal/activation/activation.go` - handoff 状态兼容逻辑（仅 legacy 请求触发）。
- `internal/model/model.go` - 增加 `ActivationStatusHandoff` 常量。

## Decisions Made
- 将冲突语义抗抢占纳入固定回归场景，防止后续主域扩展时静默侵入 WXT 主路径。
- 将 handoff 状态保持为兼容行为，但限制触发条件为无 phase/plan 标记的请求，避免影响治理态验证语义。

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] 修复 handoff 状态回归（M2 测试失败）**
- **Found during:** Task 1（执行 `go test ./... -run "TestM5DomainExpansion.*"` 后联动验证）
- **Issue:** 现有 `TestM2ActivationHandoffCarriesCarryContext` 期望 `handoff`，但结果为 `partial`。
- **Fix:** 新增 `ActivationStatusHandoff`，并在 `activation.Execute` 中对“无 phase/plan 标记且无阻断错误”的 handoff 请求恢复 `handoff` 状态。
- **Files modified:** `internal/model/model.go`, `internal/activation/activation.go`
- **Verification:** `go test ./... -run "TestM2ActivationHandoffCarriesCarryContext|TestM3HandoffContainsCarryContext|TestM6CapabilityPackActivationProducesRegisteredHandoff|TestM4ValidationTraceLinksArtifactsAndHandoff|TestM5DomainExpansion.*" -count=1`
- **Committed in:** `7a0783b`

---

**Total deviations:** 1 auto-fixed（Rule 1）
**Impact on plan:** 修复与新增回归直接相关的兼容性缺陷，无架构级变更。

## Issues Encountered
- 全量测试中发现 Phase 4 用例依赖 `completed` 状态语义；通过 phase/plan 条件约束 handoff 兼容逻辑避免互相破坏。

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- DOMN-02 回归矩阵与模板已就绪，后续新主域可复用相同准入与验收口径。
- 建议下一轮在真实第二主域任务中持续复用 `TestM5DomainExpansion.*` 与 checklist 的最小证据集合。

## Self-Check: PASSED

- FOUND: `.planning/phases/05-domain-expansion-pilot/05-02-SUMMARY.md`
- FOUND: commit `7a0783b`
- FOUND: commit `e360c85`
