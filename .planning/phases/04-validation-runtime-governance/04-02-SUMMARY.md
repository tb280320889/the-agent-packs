---
phase: 04-validation-runtime-governance
plan: 02
subsystem: validation
tags: [go, validation, registry, governance, trigger]

requires:
  - phase: 04-validation-runtime-governance
    provides: validation trace contract and run-level evidence chain baseline
provides:
  - Registry-driven validation plan generation (core + domain)
  - Hybrid trigger governance with manual rerun and auto trigger kinds
  - Plan-scoped blocking/status mapping regression coverage for VALD-01
affects: [05-domain-expansion-pilot, validator-consumers, runtime-ledger]

tech-stack:
  added: []
  patterns: [registry-defined validation manifest, trigger-kind normalization, machine-human status split]

key-files:
  created: []
  modified: [internal/activation/activation.go, internal/query/query.go, internal/model/model.go, tests/m4_validation_governance_test.go]

key-decisions:
  - "Validation plan由registry聚合main_pack与required_packs声明，强制core优先且稳定排序。"
  - "ValidationMachineView固定三态(passed/warned/failed)，ActivationResult保留流程态(completed/partial/failed)并建立映射。"
  - "validation_manual_rerun触发manual_rerun，warned路径必须输出run_id留痕动作。"

patterns-established:
  - "Registry Plan Stability: validator-core-output首位 + 其余字典序排序 + signature摘要"
  - "Plan-Scoped Blocking: failed只阻断当前activation，warned可继续但必须留痕"

requirements-completed: [VALD-01]

duration: 5 min
completed: 2026-03-18
---

# Phase 4 Plan 2: Validation Runtime Governance Summary

**实现了由registry稳定生成的core+domain校验计划，并把自动/手动混合触发与plan级阻断分流固化为可回归验证行为。**

## Performance

- **Duration:** 5 min
- **Started:** 2026-03-18T05:54:59Z
- **Completed:** 2026-03-18T06:00:37Z
- **Tasks:** 3
- **Files modified:** 4

## Accomplishments
- 在 `buildValidationPlan` 中落地 registry 聚合（main_pack + required_packs），输出稳定排序且包含 core 首位策略。
- 新增并固化 4 类 TriggerKind 与 manual rerun 入口，明确 machine/human 分层语义及 status 映射。
- 补齐 M4 回归测试，覆盖计划生成、触发类型、warned/failed 分流及 run_id 关联约束。

## Task Commits

Each task was committed atomically:

1. **Task 1: 用 registry 统一生成 core+domain validation 计划** - `5d8ab50` (feat)
2. **Task 2: 落地混合触发判定与 plan 级阻断分流** - `7a24424` (feat)
3. **Task 3: 补齐 M4 回归测试（计划生成 + 触发 + 分流）** - `89dc96c` (test)

**Plan metadata:** 待本计划文档提交生成

## Files Created/Modified
- `internal/query/query.go` - 新增 `RecommendedValidators` 与 validator manifest signature 生成。
- `internal/activation/activation.go` - 重构 validation plan 组装、trigger 解析、machine/activation status 分流与 warned 留痕动作。
- `internal/model/model.go` - 新增 validation trigger/status 常量，固定三态与流程态字面值来源。
- `tests/m4_validation_governance_test.go` - 新增 VALD-01 回归测试，校验 registry 计划生成与触发分流映射。

## Decisions Made
- validation plan 不再依赖单一 bundle 列表，而是以 registry 为真相源聚合生成并输出稳定顺序。
- failed/warned/passed 与 completed/partial/failed 分层表达：machine view 用三态，activation 保持流程态。
- `validation_manual_rerun=true` 优先覆盖 trigger kind 为 `manual_rerun`，并要求 warned 路径输出 run_id 留痕动作。

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered
- 并行执行两条 `go test` 验证命令时偶发 `SQLITE_BUSY`（`database is locked`）；改为顺序执行后通过，未改动业务代码。

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- VALD-01 已形成实现与回归保护，可继续推进 04-03 对 GOVR-01 的 runtime 账本回写制度化。

---
*Phase: 04-validation-runtime-governance*
*Completed: 2026-03-18*

## Self-Check: PASSED

- FOUND: `.planning/phases/04-validation-runtime-governance/04-02-SUMMARY.md`
- FOUND: `internal/query/query.go`
- FOUND: `internal/activation/activation.go`
- FOUND: `internal/model/model.go`
- FOUND: `tests/m4_validation_governance_test.go`
- FOUND commit: `5d8ab50`
- FOUND commit: `7a24424`
- FOUND commit: `89dc96c`
