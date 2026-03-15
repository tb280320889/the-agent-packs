# Handoff: M2 wxt-manifest 首个完整包（完成）

## 1. 交接对象
- 来源 bead：the-agent-packs-dm5
- 下一 bead：M3 kickoff（待创建）
- 来源里程碑：M2
- 目标角色：Verifier / Handoff

## 2. 已完成什么
- 在 `workflow-packages/wxt-manifest/` 下落地首个完整 workflow package 标准结构。
- 固定 `manifest-review.md` 为主 artifact，并补齐 contracts/templates/fixtures/tests 子目录。
- 在 `internal/query/query.go` 增加按主节点注入推荐 validators/artifacts 的能力。
- 在 `internal/activation/activation.go` 增加 validation plan / validator results / handoff payload 结构化输出。
- 补充 M2 测试 `tests/m2_wxt_manifest_test.go`，并通过 `go test ./...`。
- 固化包路径约定：后续统一使用 `workflow-packages/<package>/`，并更新文档约束。

## 3. 下一位 agent 可直接依赖什么
- 包根路径约定：`workflow-packages/<package>/`。
- 首包实现：`workflow-packages/wxt-manifest/`。
- M2 测试基线：`tests/m2_wxt_manifest_test.go`。
- 结构化输出样例：`fixtures/activation-result.sample.json`。

## 4. 下一位 agent 必须先做什么
- 先 claim：the-agent-packs-dm5（若继续收口）或新建 M3 bead。
- 先阅读：`docs/31-M2_上下文_workflow_package标准模板与跨包边界.md`、`docs/32-M2_上下文_wxt_manifest_Pack规格_artifact_handoff.md`、`docs/41-M3_上下文_ValidationPlan_ValidatorResult_ActivationResult.md`。
- 先验证：`go test ./...`。

## 5. 不要做什么
- 不要把 `wxt-manifest` 扩展成全量 WXT 体系。
- 不要修改 M0/M1 冻结对象语义。
- 不要绕开 `workflow-packages/` 根目录约定创建分散 package。

## 6. 风险与未决项
- 目前 validation 结果为结构化占位语义，真实 validator 逻辑由 M3 实现与收紧。
- 需要在 M3 明确 severity policy 与状态计算最终裁决边界。

## 7. 推荐下一动作
- 进入 M3：基于当前 `wxt-manifest` 输出面实现 `validator-core-output` 与 `validator-domain-wxt-manifest` 的真实执行闭环。
