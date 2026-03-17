---
phase: 01-foundation-hardening
plan: 04
subsystem: compiler
tags: [index, sqlite, compiler, recoverability, testing]

# Dependency graph
requires:
  - phase: 01-02
    provides: 事务化索引替换基础流程与结构化编译错误模型
provides:
  - 修复 Compile 时序，确保报告成功前不替换正式索引
  - 报告失败时旧索引可查询保留、成功路径完成索引替换
  - 增强 report_write 结构化错误与可恢复性回归断言
affects: [compiler, tests, phase-01-verification]

# Tech tracking
tech-stack:
  added: []
  patterns: [tmp->report->backup->swap 事务替换时序, 索引可查询真值回归测试]

key-files:
  created: []
  modified:
    - internal/compiler/compiler.go
    - tests/m4_index_rebuild_test.go
    - tests/m4_compile_errors_test.go

key-decisions:
  - "Compile 仅在 writeReports 成功后执行 dbPath 原子替换，失败路径不得触碰旧索引"
  - "回归测试以旧索引可查询性与内容稳定性作为失败路径验收真值"

patterns-established:
  - "先生成 tmp 索引，再写报告，最后做 backup/swap，避免 report 失败破坏可用索引"
  - "错误断言不仅检查 err!=nil，还验证 phase/code/path 的结构化输出"

requirements-completed: [INDX-01]

# Metrics
duration: 8 min
completed: 2026-03-16
---

# Phase 01 Plan 04: Gap Closure Summary

**修复编译时序缺口：报告失败不再提前覆盖旧索引，并以“旧索引仍可查询”建立可恢复性回归证据。**

## Performance

- **Duration:** 8 min
- **Started:** 2026-03-16T23:51:30Z
- **Completed:** 2026-03-16T23:59:30Z
- **Tasks:** 2
- **Files modified:** 3

## Accomplishments
- 移除 `Compile` 中报告前的提前 `writeIndex(dbPath, ...)` 覆写路径。
- 保留并收敛为 `tmp -> writeReports -> backup -> swap` 的事务化替换时序。
- 增强回归测试：失败场景验证旧索引可查询且内容不变，成功场景验证新索引替换生效。

## Task Commits

Each task was committed atomically:

1. **Task 1: 修复 Compile 中索引替换时序，消除报告失败前写入** - `2cb523d` (fix)
2. **Task 2: 补强索引可恢复性回归测试并执行全量回归** - `7616be6` (test)

**Plan metadata:** _pending_

## Files Created/Modified
- `internal/compiler/compiler.go` - 删除报告前提前写入正式索引的路径，避免失败覆盖旧索引。
- `tests/m4_index_rebuild_test.go` - 新增旧索引可查询保留与成功替换双路径断言。
- `tests/m4_compile_errors_test.go` - 新增 `report_write` 结构化错误断言。

## Decisions Made
- 将索引正式替换严格后移到报告成功后执行，保证失败路径可恢复。
- 将“旧索引可查询”作为用户可观察真值写入回归测试，避免仅停留在 `err != nil` 级别断言。

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] 修复回归测试中的错误蓝图路径**
- **Found during:** Task 2（补强回归测试）
- **Issue:** 初版测试将被修改节点路径写成 `L0/wxt.md`，实际文件为 `L0/wxt/overview.md`，导致测试误失败。
- **Fix:** 调整测试目标路径为真实蓝图文件路径，并保留断言逻辑不变。
- **Files modified:** `tests/m4_index_rebuild_test.go`
- **Verification:** `go test ./... -run TestM4IndexRebuild`
- **Committed in:** `7616be6` (Task 2 commit)

---

**Total deviations:** 1 auto-fixed (1 bug)
**Impact on plan:** 偏差仅为测试路径修正，不改变实现目标或范围。

## Issues Encountered
None

## Authentication Gates
None

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness
- Phase 01 的 INDX-01 缺口已闭合，可进入 Phase 02 路由治理。
- 验证口径已包含失败可恢复性真值与结构化错误可解释性。

---
## Self-Check: PASSED

*Phase: 01-foundation-hardening*
*Completed: 2026-03-16*
