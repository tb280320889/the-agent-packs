---
phase: 02-routing-governance
plan: 01
subsystem: routing
tags: [go, routing, governance, candidate-space-first, attach-only]

requires:
  - phase: 01-foundation-hardening
    provides: 可靠 YAML 解析与可恢复索引重建基线
provides:
  - RouteQuery 候选空间先过滤后评分的治理管线
  - target_pack canonical 不可用时的显式 hard-fail 语义
  - 路由稳定决策依据 decision_basis 与回归测试锁定
affects: [phase-03-contracted-delivery, activation-routing, mcp-route-query]

tech-stack:
  added: []
  patterns: [candidate-space-first, two-stage-primary-attach, stable-tie-break]

key-files:
  created: [.planning/phases/02-routing-governance/02-01-SUMMARY.md]
  modified: [internal/query/query.go, internal/model/model.go, tests/m1_minimal_test.go, tests/m3_validation_closure_test.go]

key-decisions:
  - "RouteQuery 统一为候选空间过滤后再评分，避免非法候选污染主竞争。"
  - "target_pack canonical 不可用时直接 hard-fail，不再返回 registry fallback 候选。"
  - "路由输出增加 decision_basis 作为可复验稳定决策依据。"

patterns-established:
  - "Pattern 1: L0/L1 先做候选空间白名单过滤，再进入统一评分与稳定 tie-break。"
  - "Pattern 2: attach-only capability 仅进入 must_include 附挂集合，不参与 primary candidates。"

requirements-completed: [ROUT-01, ROUT-02]

duration: 5 min
completed: 2026-03-17
---

# Phase 2 Plan 1: Routing Governance Core Summary

**RouteQuery 已落地 candidate-space-first + attach-only 两阶段治理，并通过 decision_basis 与回归测试矩阵保证冲突决策稳定可复验。**

## Performance

- **Duration:** 5 min
- **Started:** 2026-03-17T02:52:29Z
- **Completed:** 2026-03-17T02:58:27Z
- **Tasks:** 3
- **Files modified:** 4

## Accomplishments
- 重构 `RouteQuery` 为“候选空间过滤 → 评分 → 稳定排序”管线，阻断非法候选先评分的问题。
- 固化 `target_pack` 的 canonical 不可用 hard-fail 语义，并移除 `registry fallback` 伪成功路径。
- 输出 `decision_basis` 机器可读字段并补齐回归断言，锁死 attach-only 边界与 must_include 稳定顺序。

## Task Commits

Each task was committed atomically:

1. **Task 1: 重构 RouteQuery 为 candidate-space-first 管线** - `acb568e` (feat)
2. **Task 2: 固化稳定冲突决策链并保留可复现依据** - `b0f8f31` (feat)
3. **Task 3: 扩展治理回归用例锁死候选污染与附挂边界** - `889030c` (test)

## Files Created/Modified
- `internal/query/query.go` - 两阶段候选池过滤、稳定排序、target_pack hard-fail 与 must_include 稳定化
- `internal/model/model.go` - 新增 `RouteResult.decision_basis` 机器可读解释字段
- `tests/m1_minimal_test.go` - 增加 decision_basis 与 must_include 稳定顺序断言
- `tests/m3_validation_closure_test.go` - 增加 target_pack 快路径 decision_basis 与 level 不匹配 hard-fail 回归

## Decisions Made
- 采用统一稳定 tie-break 输出 `score>canonical>domain>rule>lexicographic`，确保同输入可复现。
- `target_pack` 路径不再容忍 canonical fallback，避免错误映射污染主决策。
- 测试断言优先 machine-readable 字段（`decision_basis`、`must_include`），而非仅文案片段。

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered
None

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness
- `ROUT-01/ROUT-02` 已闭合，RouteQuery 治理内核进入可验证稳定状态。
- 已具备进入 `02-02-PLAN.md` 的输入基础（可解释输出细化与 canonical 缺失语义延展）。

## Self-Check: PASSED

- FOUND: `.planning/phases/02-routing-governance/02-01-SUMMARY.md`
- FOUND: `acb568e`
- FOUND: `b0f8f31`
- FOUND: `889030c`

---
*Phase: 02-routing-governance*
*Completed: 2026-03-17*
