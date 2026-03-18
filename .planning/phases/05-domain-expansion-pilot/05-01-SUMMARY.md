---
phase: 05-domain-expansion-pilot
plan: 01
subsystem: routing
tags: [monorepo, registry, feature-switch, domain-onboarding]

requires:
  - phase: 04-validation-runtime-governance
    provides: validation trace、runtime ledger 与 attach-only 路由治理基线
provides:
  - 第二主域 monorepo-oss-governance 的 blueprint/package/registry 准入资产
  - DOMAIN_MONOREPO_ENABLED 特性开关与 route->activate 准入回归测试
affects: [05-02, DOMN-01, domain-onboarding]

tech-stack:
  added: []
  patterns: [candidate-space-first feature switch, registry canonical onboarding]

key-files:
  created:
    - blueprint/L0/monorepo/overview.md
    - blueprint/L1/monorepo/oss-governance.md
    - workflow-packages/monorepo-oss-governance/package.yaml
    - tests/m5_domain_expansion_registry_test.go
    - tests/m5_domain_expansion_onboarding_test.go
  modified:
    - workflow-packages/registry.json
    - internal/query/query.go
    - tests/m2_registry_test.go

key-decisions:
  - "第二主域采用 monorepo-oss-governance 单子域闭环接入，保持 attach-only 能力线不迁移。"
  - "通过 DOMAIN_MONOREPO_ENABLED 控制 monorepo 域是否进入候选空间，默认开启以支持验收，关闭即回退。"

patterns-established:
  - "新主域准入先落 blueprint+package+registry，再补 route/activation 自动化回归。"
  - "主域开关应作用于 candidate-space 过滤层，而非评分层。"

requirements-completed: [DOMN-01]

duration: 6min
completed: 2026-03-18
---

# Phase 5 Plan 01: Domain Onboarding Summary

**以不改治理骨架的方式接入 monorepo-oss-governance 第二主域，并通过开关化候选空间控制实现可回退准入。**

## Performance

- **Duration:** 6 min
- **Started:** 2026-03-18T08:31:00Z
- **Completed:** 2026-03-18T08:34:38Z
- **Tasks:** 2
- **Files modified:** 8

## Accomplishments
- 新增 monorepo 域 L0/L1 blueprint 节点与 workflow package 定义。
- registry 注册 `monorepo-oss-governance`，canonical 映射到 `L1.monorepo.oss-governance`，并保持 capability attach-only 线不变。
- 在 `RouteQuery` 增加 `DOMAIN_MONOREPO_ENABLED` 开关，并补齐准入/回退自动化测试。

## Task Commits

1. **Task 1: 落地第二主域资产与注册表准入** - `4ce0af6` (feat)
2. **Task 2: 实现新域 feature switch 与 DOMN-01 准入回归** - `0d4d9cb` (feat)

## Files Created/Modified
- `blueprint/L0/monorepo/overview.md` - 第二主域 L0 根节点定义。
- `blueprint/L1/monorepo/oss-governance.md` - monorepo 子域 workflow 入口节点。
- `workflow-packages/monorepo-oss-governance/package.yaml` - 新 workflow package 契约。
- `workflow-packages/registry.json` - 注册新 workflow package canonical 映射。
- `internal/query/query.go` - 新增 `DOMAIN_MONOREPO_ENABLED` 候选空间开关。
- `tests/m5_domain_expansion_registry_test.go` - 新域 registry 准入与 capability 不回归测试。
- `tests/m5_domain_expansion_onboarding_test.go` - route/activate 正例与回退用例。
- `tests/m2_registry_test.go` - 同步 registry 包数量基线到 4。

## Decisions Made
- 使用独立 monorepo 主域（非 capability 线）作为第二主域试点，以验证主域扩展机制。
- feature switch 放在路由候选空间过滤阶段，确保不开启时不会产生“隐性竞争”。

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] 更新 registry 包数量回归基线**
- **Found during:** Task 1（registry 校验回归）
- **Issue:** `TestM2RegistryLoadsAndValidates` 仍固定期望 3 个包，新增主域后失败。
- **Fix:** 将期望更新为 4，保持基线与 registry 真实状态一致。
- **Files modified:** `tests/m2_registry_test.go`
- **Verification:** `go test ./... -run "TestM2RegistryLoadsAndValidates|TestM5DomainExpansionRegistryOnboarding" -count=1`
- **Committed in:** `4ce0af6`

---

**Total deviations:** 1 auto-fixed（Rule 1）
**Impact on plan:** 仅修复新增主域导致的基线漂移，无额外范围扩张。

## Issues Encountered
- 回退用例初次验证未生效，补充 `DOMAIN_MONOREPO_ENABLED` 候选空间过滤后恢复预期行为。

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- DOMN-01 已具备可执行准入闭环。
- 可进入 05-02 执行 WXT 非回归矩阵与准入清单沉淀。

## Self-Check: PASSED

- FOUND: `.planning/phases/05-domain-expansion-pilot/05-01-SUMMARY.md`
- FOUND: commit `4ce0af6`
- FOUND: commit `0d4d9cb`
