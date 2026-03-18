---
phase: 04-validation-runtime-governance
plan: 04
subsystem: runtime-governance
tags: [runtime-ledger, record-type-mapping, gap-closure, governance]

# Dependency graph
requires:
  - phase: 04-validation-runtime-governance
    provides: runtime ledger append-only governance baseline (04-03)
provides:
  - key change -> record_type 映射规则（assumption/decision/change/validation）
  - activation 路径多类型 runtime ledger 写入与稳定顺序输出
  - 四类 record_type 触发矩阵与按类型版本追加回归保护
affects: [activation, model, tests, verification]

# Tech tracking
tech-stack:
  added: []
  patterns: [deterministic-type-mapping, append-only-by-trace-and-type, tdd-regression-matrix]

key-files:
  created: []
  modified:
    - internal/activation/activation.go
    - internal/model/model.go
    - tests/m4_validation_governance_test.go

key-decisions:
  - "record_type 由关键事件规则动态映射，validation 永远写入，其余类型按触发条件增量追加。"
  - "多类型输出固定按 assumption -> decision -> change -> validation 顺序，确保同输入稳定可回归。"
  - "append-only 粒度维持 TraceID+RecordType，禁止跨类型共享版本号或覆盖历史。"

requirements-completed: [GOVR-01]

# Metrics
duration: 29 min
completed: 2026-03-18
---

# Phase 04 Plan 04: Runtime Ledger Record-Type Gap Closure Summary

**已将 runtime ledger 从“仅 validation 单类型写入”修复为“按关键变更同构写入 assumption/decision/change/validation 四类 record_type”，并补齐可回归测试矩阵，闭合 GOVR-01 阻断缺口。**

## Performance

- **Tasks:** 2
- **Task commits:** 3（TDD RED + GREEN + 回归矩阵）
- **Files modified:** 3

## Accomplishments

- 在 `internal/activation/activation.go` 重构 `BuildRuntimeLedgerEntries`：
  - 移除固定 `RecordType=validation` 单写入路径。
  - 新增 `determineRuntimeLedgerRecordTypes`，按 plan 规则映射四类 record_type。
  - 多类型写入执行去重 + 固定顺序输出。
  - 延续并适配 `TraceID+RecordType` append-only 版本逻辑。
- 在 `internal/model/model.go` 新增 `IsRuntimeLedgerRecordType` 以集中约束 record_type 判定。
- 在 `tests/m4_validation_governance_test.go` 新增并强化回归用例：
  - `TestM4RuntimeLedgerRecordTypesFromKeyChanges`
  - `TestM4RuntimeLedgerRecordTypeVersionAppendOnly`
  - 同步修正既有 append-only/deferred 用例断言，匹配多类型输出后的真实行为。

## Task Commits

1. **Task 1 (RED):** `ed9f271` — `test(04-04): add failing test for runtime ledger record type mapping`
2. **Task 1 (GREEN):** `019f239` — `feat(04-04): map key changes to runtime ledger record types`
3. **Task 2:** `af67642` — `test(04-04): add record-type regression matrix and version checks`

## Verification

- `go test ./... -run "TestM4RuntimeLedger(RecordTypesFromKeyChanges|RecordTypeVersionAppendOnly|VersionAppendOnly|DeferredWindowAndEscalation)" -count=1` ✅
- `go test ./... -run "TestM4Validation.*" -count=1` ✅

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] 修复既有 runtime ledger 测试对“单类型条目数”的过时断言**
- **Found during:** Task 2 verification
- **Issue:** 多类型写入启用后，旧测试仍按单类型假设 `len(entries)`，导致误报失败。
- **Fix:** 更新 `TestM4RuntimeLedgerVersionAppendOnly` 与 `TestM4RuntimeLedgerDeferredWindowAndEscalation` 的断言口径，改为按 `record_type` 维度验证版本链与 current 指针。
- **Files modified:** `tests/m4_validation_governance_test.go`
- **Commit:** `af67642`

---

**Total deviations:** 1 auto-fixed (Rule 1)

## Issues Encountered

None.

## Next Phase Readiness

- Phase 04 的 GOVR-01 gap 已按代码+测试闭环修复。
- 可进入后续 re-verify，确认 04-VERIFICATION.md 的阻断项转为 satisfied。

## Self-Check: PASSED

- FOUND: `.planning/phases/04-validation-runtime-governance/04-04-SUMMARY.md`
- FOUND commits: `ed9f271`, `019f239`, `af67642`
