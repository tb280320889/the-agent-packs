---
phase: 03-contracted-delivery
plan: 02
subsystem: validation
tags: [contract-delivery, validator, wxt, go-test]

# Dependency graph
requires:
  - phase: 03-01
    provides: ContextBundle contract decisions（included/excluded）
provides:
  - validator-contract-delivery 契约检查器（P0 hard-fail + warning）
  - contract-delivery validator 注册与 runner 可执行链路
  - WXT 正/负/warning 可重复检查矩阵
affects: [phase-04-validation-runtime-governance, activation-validation-flow]

# Tech tracking
tech-stack:
  added: []
  patterns: [contract-finding-with-rule-trace, p0-hard-fail-policy, deterministic-contract-test-matrix]

key-files:
  created: [internal/validator/contract_delivery.go]
  modified: [internal/model/model.go, internal/validator/registry.go, internal/validator/types.go, internal/activation/activation.go, internal/query/query.go, tests/m3_contract_bundle_test.go]

key-decisions:
  - "将 contract delivery 检查器作为独立 validator（validator-contract-delivery）接入现有 runner，而非内嵌到 route/query 逻辑。"
  - "P0 规则统一以 finding.severity=error + status=failed 表达，非阻断语义以 warned 表达。"
  - "validator 输入以 ContextBundle 为契约真相源，Activation 层仅透传 ContractBundle。"

patterns-established:
  - "Pattern 1: finding 同时携带 code/severity/rule_ref/source_rule，支持规则追溯与后续治理接力"
  - "Pattern 2: 合同检查输出必须包含最小修复建议，不允许仅 pass/fail 文本"

requirements-completed: [CONT-03]

# Metrics
duration: 7 min
completed: 2026-03-17
---

# Phase 3 Plan 2: Contract Delivery Validator Summary

**交付契约检查已升级为可执行 validator：可对跨域混入、required 缺失、规则不可追溯执行 hard-fail，并输出结构化修复建议。**

## Performance

- **Duration:** 7 min
- **Started:** 2026-03-17T04:14:45Z
- **Completed:** 2026-03-17T04:21:38Z
- **Tasks:** 3
- **Files modified:** 7

## Accomplishments
- 新增 `validator-contract-delivery`，实现 P0/warning 分级与结构化检查输出。
- 将契约检查器接入现有 registry/runner，并保持主流程测试可通过。
- 扩展 `m3_contract_bundle_test.go` 为正例 + 负例 + warning 的可重复检查矩阵。

## Task Commits

Each task was committed atomically:

1. **Task 1: 实现 contract delivery validator 与分级失败策略** - `985510e` (feat)
2. **Task 2: 注册 validator 并接入现有 runner 执行链路** - `4c88fc2` (feat)
3. **Task 3: 扩展正负例测试矩阵覆盖 hard-fail 与 warning 语义** - `ebe7606` (feat)

## Files Created/Modified
- `internal/validator/contract_delivery.go` - 契约检查器实现（P0/warning、规则追溯、修复建议）
- `internal/model/model.go` - Finding 新增 `rule_ref/source_rule` 字段
- `internal/validator/registry.go` - 注册 `validator-contract-delivery`
- `internal/validator/types.go` - ExecutionInput 增加 `ContractBundle`
- `internal/activation/activation.go` - 将 bundle 契约数据透传给 validator
- `internal/query/query.go` - 修正 attach-only required 节点的 include 决策 scope/依据
- `tests/m3_contract_bundle_test.go` - 新增/扩展 contract validator 正负例与 warning 断言

## Decisions Made
- 采用独立 validator 承载契约检查，保持 route/query 构建与验证职责分离。
- P0 违规固定为 `failed`（error），说明弱但结构完整固定为 `warned`（warn）。
- contract 规则追溯统一落在 finding 结构字段中，供 Phase 4 直接消费。

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] 补齐 validator 输入中的 ContractBundle 透传**
- **Found during:** Task 3 (扩展正负例测试矩阵覆盖 hard-fail 与 warning 语义)
- **Issue:** 新增契约检查器若拿不到 bundle contract 数据，将无法执行 required/source_rule/cross-domain 规则验证。
- **Fix:** 在 `ExecutionInput` 增加 `ContractBundle`，并在 activation 执行链路透传 bundle。
- **Files modified:** `internal/validator/types.go`, `internal/activation/activation.go`
- **Verification:** `go test ./... -run "TestContractDeliveryValidator.*"` 与 `go test ./... -run "TestM3GoldenCompleted"` 通过
- **Committed in:** `ebe7606`

**2. [Rule 1 - Bug] 修复 attach-only required 节点被误判为跨域混入**
- **Found during:** Task 3 (扩展正负例测试矩阵覆盖 hard-fail 与 warning 语义)
- **Issue:** required_with 引入的 attach-only 节点沿用 `target_domain` scope，导致 validator 误报 `CONTRACT_CROSS_DOMAIN_INCLUDED`。
- **Fix:** 在 bundle 构建时，attach-only required 节点改为 `attach_only_capability` scope，并记录合法依赖决策依据。
- **Files modified:** `internal/query/query.go`
- **Verification:** `go test ./... -run "TestContractDeliveryValidator.*"` 通过且 warning/negative 语义稳定
- **Committed in:** `ebe7606`

---

**Total deviations:** 2 auto-fixed (1 blocking, 1 bug)
**Impact on plan:** 均为完成 CONT-03 所必需修复，无范围外扩。

## Issues Encountered
- `TestContractDeliveryValidatorPositive` 在并发跑测时偶发 `SQLITE_BUSY`，改为顺序执行验证命令后稳定通过。

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- Contract delivery 检查语义已结构化，可直接为 Phase 4 的 validation/runtime 治理复用。
- 已具备规则追溯字段与修复建议输出，可作为 runtime 回写输入。

---
*Phase: 03-contracted-delivery*
*Completed: 2026-03-17*

## Self-Check: PASSED
