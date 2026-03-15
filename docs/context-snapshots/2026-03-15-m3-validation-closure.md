# Context Snapshot: M3 validators 与主任务闭环

## 1. 当前阶段
- 所属里程碑：M3
- 关联 bead：the-agent-packs-2dh
- 当前状态：in_progress

## 2. 当前事实
- 当前要解决的问题：把 Validation Plan -> Validators -> Activation Result 从占位结构升级为可执行闭环，并补齐固定主任务回归矩阵。
- 当前已完成内容：已将 ActivationResult 与验证对象改为强类型；新增 `internal/validator/` 执行框架；实现 `validator-core-output` 与 `validator-domain-wxt-manifest`；重构 activation 状态计算为策略函数；新增 M3 闭环测试并通过。
- 当前尚未完成内容：M4 冻结阶段的策略收口与扩域准入不在本 bead 范围。

## 3. 已冻结对象
- `model.ValidationPlan`：作为验证计划一等对象，固定 `validators`、`severity_policy`、`artifacts_under_validation`。
- `model.ValidatorResult`：固定 `status/findings/repair_suggestions/validated_artifacts` 结构。
- `model.ValidationEnvelope`：固定 activation 中 `validation_results` 的统一容器形状。
- `activation.deriveStatus`：固定状态裁决优先级 `failed > handoff > partial > completed`。

## 4. 当前输入
- 上游交付物：M2 `wxt-manifest` 包与 M2 回归基线。
- 依赖文档：`docs/40-M3_Validators与主任务闭环_开发指导.md`、`docs/41-M3_上下文_ValidationPlan_ValidatorResult_ActivationResult.md`、`docs/42-M3_上下文_测试矩阵与固定主任务闭环.md`。
- 依赖 schema / 模板 / fixtures：`fixtures/activation-request.sample.json`、`fixtures/activation-result.sample.json`。

## 5. 当前输出
- 已产出文件：`internal/validator/types.go`、`internal/validator/registry.go`、`internal/validator/runner.go`、`internal/validator/core_output.go`、`internal/validator/domain_wxt_manifest.go`、`tests/m3_validation_closure_test.go`。
- 已更新文件：`internal/model/model.go`、`internal/activation/activation.go`、`tests/m2_wxt_manifest_test.go`。
- 已创建 bead：`the-agent-packs-2dh`。

## 6. 风险与阻塞
- 风险：测试使用共享 `blueprint/index/blueprint.db`，并行执行可能触发 SQLite busy；需保持串行执行。
- 阻塞：无。
- 是否需要人工决策：否。

## 7. 下一步建议
- 建议下一个 bead：M4 冻结与扩域准入准备。
- 建议先执行的命令：`go test ./...`、`bd status --json`。
- 建议先阅读的文档：`docs/50-M4_Phase1冻结与扩域准入_开发指导.md`、`docs/51-M4_上下文_冻结策略_修改纪律_BreakingChange.md`。
