---
phase: 01-foundation-hardening
plan: 02
subsystem: compiler
tags: [index, sqlite, compiler, cli, mcp]

# Dependency graph
requires:
  - phase: 01-01
    provides: 解析硬化的 YAML 严格解码
provides:
  - 事务化索引重建与原子替换
  - 编译结构化错误输出（phase/path/code/message）
  - CLI/MCP 输出可被消费的编译结果
affects: [compiler, cli, tests]

# Tech tracking
tech-stack:
  added: []
  patterns: [临时 DB 写入 + 原子替换, 编译阶段结构化错误输出]

key-files:
  created:
    - internal/compiler/compiler_errors.go
    - tests/m4_compile_errors_test.go
    - tests/m4_index_rebuild_test.go
    - tests/m4_cli_compile_output_test.go
  modified:
    - internal/compiler/compiler.go
    - internal/model/model.go
    - cmd/agent-pack-mcp/main.go
    - tests/m1_minimal_test.go
    - tests/m4_parsing_hardening_test.go

key-decisions:
  - "编译结果统一返回 CompileResult(errors) 以便 CLI/MCP 与测试消费"
  - "索引重建先写临时 DB，再在报告成功后原子替换目标索引"

patterns-established:
  - "编译阶段错误统一为 phase/path/code/message 的结构化输出"
  - "索引重建采用 tmp + bak 的原子替换回滚策略"

requirements-completed: [INDX-01, INDX-02]

# Metrics
duration: 3 min
completed: 2026-03-16
---

# Phase 01 Plan 02: Foundation Hardening Summary

**索引重建以临时 DB 原子替换保障失败不破坏旧索引，并提供可序列化的编译结构化错误输出**

## Performance

- **Duration:** 3 min
- **Started:** 2026-03-16T13:28:51Z
- **Completed:** 2026-03-16T13:32:33Z
- **Tasks:** 3
- **Files modified:** 9

## Accomplishments
- 引入编译阶段结构化错误与 CompileResult，覆盖 parse/index/report 失败路径
- 索引重建改为临时 DB 写入 + 原子替换，报告失败不触发替换
- CLI/MCP 输出统一返回结构化编译结果并补齐测试覆盖

## Task Commits

Each task was committed atomically:

1. **Task 1: 定义结构化编译错误并贯穿编译阶段** - `d0ce69c` (feat)
2. **Task 2: 实现索引重建的临时 DB + 原子替换回滚** - `081f292` (feat)
3. **Task 3: 调整 CLI/MCP 输出以返回结构化错误** - `bda383c` (feat)

**Plan metadata:** _pending_

_Note: TDD tasks may have multiple commits (test → feat → refactor)_

## Files Created/Modified
- `internal/compiler/compiler_errors.go` - 编译阶段错误与结果结构定义
- `internal/compiler/compiler.go` - 编译流程结构化错误、临时 DB 重建与原子替换
- `internal/model/model.go` - 编译错误与结果序列化结构
- `cmd/agent-pack-mcp/main.go` - CLI/MCP 输出结构化编译结果
- `tests/m4_compile_errors_test.go` - 结构化错误输出测试
- `tests/m4_index_rebuild_test.go` - 索引重建事务性测试
- `tests/m4_cli_compile_output_test.go` - CLI 输出 JSON 结果测试
- `tests/m1_minimal_test.go` - YAML 变体测试修正
- `tests/m4_parsing_hardening_test.go` - 结构化错误断言调整

## Decisions Made
- 统一使用 CompileResult(errors) 作为编译结果输出，便于 CLI/MCP 与测试消费
- 在报告写入成功后再原子替换索引，避免失败覆盖旧索引

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] 修复 YAML 冒号解析测试失败**
- **Found during:** Task 3 (CLI 输出调整后的全量测试)
- **Issue:** YAML 摘要含冒号未加引号导致解析失败
- **Fix:** 将测试用例中的 summary 改为带引号字符串
- **Files modified:** tests/m1_minimal_test.go
- **Verification:** go test ./...
- **Committed in:** bda383c (Task 3 commit)

---

**Total deviations:** 1 auto-fixed (1 bug)
**Impact on plan:** 修复测试用例以匹配严格 YAML 解析要求，无范围扩张。

## Issues Encountered
None

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness
- 索引与编译失败态已结构化，可进入路由治理的后续计划
- Phase 01 仍需完成 01-03 计划

---
## Self-Check: PASSED

*Phase: 01-foundation-hardening*
*Completed: 2026-03-16*
