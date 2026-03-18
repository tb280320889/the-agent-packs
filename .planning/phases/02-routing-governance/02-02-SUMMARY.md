---
phase: 02-routing-governance
plan: 02
subsystem: api
tags: [routing, activation, mcp, contract, explainability]

# Dependency graph
requires:
  - phase: 02-routing-governance
    provides: candidate-space-first 与 attach-only 两阶段路由基线
provides:
  - route_query/activate 的可解释极简契约（status/error_code/next_action/decision_trace_id/docs_ref）
  - canonical 缺失显式失败语义（ROUTE_CANONICAL_MISSING）
  - capability 附挂/拒绝的 machine-readable 原因码与规则码
affects: [phase-03-contracted-delivery, mcp-consumers, route-contract-tests]

# Tech tracking
tech-stack:
  added: []
  patterns: [minimal-explainable-contract, explicit-canonical-failure, stable-decision-trace]

key-files:
  created: [.planning/phases/02-routing-governance/02-02-SUMMARY.md]
  modified:
    - internal/model/model.go
    - internal/query/query.go
    - internal/activation/activation.go
    - cmd/agent-pack-mcp/main.go
    - tests/m1_minimal_test.go
    - tests/m3_validation_closure_test.go

key-decisions:
  - "RouteResult 默认输出极简 machine-readable 字段并预留 details/docs_ref 扩展位。"
  - "target_pack canonical 缺失或不可路由统一 hard-fail，错误码固定为 ROUTE_CANONICAL_MISSING。"
  - "ActivationResult 透传 route 语义字段，避免接入层吞掉失败原因。"

patterns-established:
  - "Pattern: 路由失败必须返回可执行 next_action，禁止静默 fallback。"
  - "Pattern: capability 附挂输出 reason_code/rule_ref 以支持调用方程序化消费。"

requirements-completed: [ROUT-03, ROUT-04]

# Metrics
duration: 8 min
completed: 2026-03-17
---

# Phase 02 Plan 02: 可解释路由契约与 canonical 显式失败 Summary

**route/activate 现已输出可解析的极简解释契约，并在 canonical 映射缺失时稳定返回 ROUTE_CANONICAL_MISSING 与下一步建议。**

## Performance

- **Duration:** 8 min
- **Started:** 2026-03-17T03:03:49Z
- **Completed:** 2026-03-17T03:11:58Z
- **Tasks:** 3
- **Files modified:** 6

## Accomplishments
- RouteResult 新增 status/error_code/message/next_action/decision_trace_id/docs_ref 与 capability_decisions 字段，默认保持极简。
- 删除 canonical 缺失场景的静默回退路径，固化 ROUTE_CANONICAL_MISSING 显式失败语义。
- activation 与 MCP 层对 route 语义完成透传，并通过回归测试锁死机器可读字段。

## Task Commits

1. **Task 1: 定义可解释路由结果契约并实现默认极简输出** - `51d3932` (feat)
2. **Task 2: 删除 canonical fallback 并固化 ROUTE_CANONICAL_MISSING 显式失败语义** - `5cea7e8` (feat)
3. **Task 3: 增加解释契约与失败语义回归矩阵** - `9298221` (test)

## Files Created/Modified
- `internal/model/model.go` - 扩展 RouteResult / ActivationResult 可解释契约字段
- `internal/query/query.go` - canonical missing hard-fail 与最小解释输出实现
- `internal/activation/activation.go` - route 语义透传到 activation 输出
- `cmd/agent-pack-mcp/main.go` - MCP tools 增加 activate 并保持新字段可见
- `tests/m1_minimal_test.go` - 路由与 activation 解释字段断言增强
- `tests/m3_validation_closure_test.go` - canonical 缺失与 route 语义回归断言

## Decisions Made
- 采用“双层模型”：默认极简字段 + details 预留扩展位，满足调用方稳定解析与后续扩展并存。
- canonical 缺失统一 hard-fail，不再返回 fallback candidate，避免伪成功语义污染。
- activation summary 在无候选场景优先复用 route message，确保错误语义一致。

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered
- 测试初次断言将“上下文不足导致 activation partial”误绑定为 route partial，已修正为 route completed + activation partial（由验证策略降级）。

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness
- ROUT-03/ROUT-04 已闭合，Phase 2 全部计划执行完成。
- 已具备进入 Phase 3（Contracted Delivery）所需的稳定 route/activation 解释契约。

## Self-Check: PASSED
- FOUND: `.planning/phases/02-routing-governance/02-02-SUMMARY.md`
- FOUND: commit `51d3932`
- FOUND: commit `5cea7e8`
- FOUND: commit `9298221`
