# Handoff: M3 验证闭环实现

> 弃用说明：本文件属于正式主线 `docs/00~52` 的历史交接，仅保留用于追溯；改造计划 v1 请改读 `docs/改造计划v1/handoffs/`。

## 1. 交接对象
- 来源 bead：the-agent-packs-2dh
- 下一 bead：M4 kickoff（待创建）
- 来源里程碑：M3
- 目标角色：Verifier / Handoff

## 2. 已完成什么
- 将 `ActivationResult` 从 `any` 聚合改为强类型对象，固定 `artifacts`、`validation_results`、`handoff` 的结构。
- 新增 `internal/validator/` 执行器与注册机制，实现 `validator-core-output`、`validator-domain-wxt-manifest` 两个可执行 validator。
- 重构 `internal/activation/activation.go`：固定执行顺序 route -> bundle -> artifact -> plan -> validators -> result，并统一状态裁决。
- 新增 `tests/m3_validation_closure_test.go`，覆盖 golden/negative/partial/handoff 关键场景。
- 执行 `go test ./tests -v` 与 `go test ./...`，结果通过。

## 3. 下一位 agent 可直接依赖什么
- 验证执行入口：`validator.Run(plan, input)`。
- 状态裁决函数：`deriveStatus(...)`（位于 `internal/activation/activation.go`）。
- M3 回归测试集：`tests/m3_validation_closure_test.go`。
- 强类型模型：`internal/model/model.go` 中 Validation 相关对象。

## 4. 下一位 agent 必须先做什么
- 先 claim：M4 bead（创建后领取）。
- 先阅读：`docs/50-M4_Phase1冻结与扩域准入_开发指导.md`、`docs/51-M4_上下文_冻结策略_修改纪律_BreakingChange.md`、`docs/52-M4_上下文_第二领域线准入与可复用骨架.md`。
- 先验证：`go test ./...`。

## 5. 不要做什么
- 不要回退 M3 已冻结验证对象为松散 `any` 结构。
- 不要把 `validator-domain-wxt-manifest` 扩张为跨领域通用规则引擎。
- 不要在 route 阶段默认深读 L2/L3 来替代最小上下文策略。

## 6. 风险与未决项
- SQLite 测试数据库存在并发锁风险，测试建议串行运行。
- 第二领域线 validator 复用策略尚未定义，属于 M4 范围。

## 7. 推荐下一动作
- 进入 M4：冻结 M3 产物接口，定义 second-line 准入模板，验证不破坏 M0~M3 冻结对象。
