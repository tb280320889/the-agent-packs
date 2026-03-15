# Handoff: 改造计划 v1 M5 实现迁移执行与差异整改

## 1. 交接对象
- 来源 bead：the-agent-packs-1lg
- 下一 bead：the-agent-packs-18k
- 来源里程碑：M5
- 目标角色：项目内部维护 agent / 迭代开发子 agent / Verifier

## 2. 已完成什么
- 已完成注册表与 `package.yaml` 的关键声明一致性校验，包括 `depends_on`、`validators`、`artifacts` 与 `required_packs` 的实现层对齐。
- 已完成 route 对显式 `target_pack` 的注册表优先回退，并把 `required_packs` 注入 `must_include`。
- 已补齐 capability 包的推荐 validator / artifact，打通资产归属与 activation 输出链路。
- 已在 validator 层校验 handoff `to_packs` 与注册表 `required_packs` 一致。
- 已新增 M5 收口结论：`docs/改造计划v1/62-M5_整改结果与收口结论.md`
- 已新增 M5 收口快照：`docs/改造计划v1/context-snapshots/2026-03-15-m5-execution-alignment.md`

## 3. 下一位 agent 可直接依赖什么
- M5 收口结论：`docs/改造计划v1/62-M5_整改结果与收口结论.md`
- M5 收口快照：`docs/改造计划v1/context-snapshots/2026-03-15-m5-execution-alignment.md`
- 已对齐实现：`workflow-packages/registry.json`、`internal/registry/registry.go`、`internal/query/query.go`、`internal/activation/activation.go`、`internal/validator/core_output.go`
- 已补齐测试：`tests/m1_minimal_test.go`、`tests/m2_registry_test.go`、`tests/m2_wxt_manifest_test.go`、`tests/m3_validation_closure_test.go`
- bead：`the-agent-packs-18k`

## 4. 下一位 agent 必须先做什么
- 先 claim：`bd update the-agent-packs-18k --status=in_progress --json`
- 先阅读：`docs/改造计划v1/62-M5_整改结果与收口结论.md`、`docs/改造计划v1/70-M6_运行验证_第二领域线准入实施与持续治理_开发指导.md`、`docs/改造计划v1/71-M6_上下文_验证矩阵_准入执行与治理节奏.md`
- 先验证：`go test ./...`、`bd show the-agent-packs-18k --json`

## 5. 不要做什么
- 不要重新打开 M5 已通过的 registry / route 对齐结论去承载 M6 的运行策略讨论。
- 不要绕过已收口的 `required_packs -> must_include -> handoff to_packs` 约束链路。
- 不要在第二领域线准入时重新引入空 asset、空 validator 的 capability 注册项。

## 6. 风险与未决项
- 当前无 M5 范围内未决整改项。
- M6 需继续验证第二领域线是否同样满足注册表真相源、attach-only 与最小上下文消费约束。
- lint 自动检测工具当前存在工具侧异常，本次未作为 M5 阻塞项处理。

## 7. 推荐下一动作
- 直接认领 `the-agent-packs-18k`，基于 M5 已收口输出执行运行验证矩阵、第二领域线准入与持续治理闭环。
