# Handoff: M4 Phase1 冻结与扩域准入（规则固化阶段）

> 弃用说明：本文件属于正式主线 `docs/00~52` 的历史交接，仅保留用于追溯；改造计划 v1 请改读 `docs/改造计划v1/handoffs/`。

## 1. 交接对象
- 来源 bead：the-agent-packs-3w0
- 下一 bead：the-agent-packs-3w0（同 bead 继续执行后续子任务）
- 来源里程碑：M4
- 目标角色：Executor / Verifier / Handoff

## 2. 已完成什么
- 在 `docs/50-M4_Phase1冻结与扩域准入_开发指导.md` 增补了 M4 子任务拆分、交付落点、对外能力边界说明与 DoD。
- 在 `docs/51-M4_上下文_冻结策略_修改纪律_BreakingChange.md` 落地了冻结对象可追溯表、non-breaking 判定表、breaking 影响评估项与处理流程模板。
- 在 `docs/52-M4_上下文_第二领域线准入与可复用骨架.md` 落地了准入检查清单、接入步骤输入输出定义、复用失败信号与下游最小交接包。
- 在 `docs/52-M4_上下文_第二领域线准入与可复用骨架.md` 补充了候选领域评审输出模板与优先筛选提示，便于正式评估第二领域线。
- 新增 M4 阶段快照：`docs/context-snapshots/2026-03-15-m4-phase1-freeze-admission.md`。

## 3. 下一位 agent 可直接依赖什么
- 冻结对象与红线边界：`docs/51-M4_上下文_冻结策略_修改纪律_BreakingChange.md`。
- 第二领域线准入模板与复用执行步骤：`docs/52-M4_上下文_第二领域线准入与可复用骨架.md`。
- 候选领域评审记录模板：`docs/52-M4_上下文_第二领域线准入与可复用骨架.md` 中“候选领域评审输出模板”小节。
- 阶段事实与风险基线：`docs/context-snapshots/2026-03-15-m4-phase1-freeze-admission.md`。

## 4. 下一位 agent 必须先做什么
- 先 claim：`the-agent-packs-3w0`（如果发生重新分配或状态变化）。
- 先阅读：`docs/50-M4_Phase1冻结与扩域准入_开发指导.md`、`docs/51-M4_上下文_冻结策略_修改纪律_BreakingChange.md`、`docs/52-M4_上下文_第二领域线准入与可复用骨架.md`。
- 先验证：`go test ./...`（确认文档更新未影响既有实现与测试基线）。

## 5. 不要做什么
- 不要以“扩域效率”为由改系统层主判断或顶层 envelope 语义。
- 不要在未走 breaking 流程时调整 route 优先级或最小 bundle 结构。
- 不要跳过准入检查清单直接启动第二领域线实现。

## 6. 风险与未决项
- 第二领域线候选尚未最终确定；需基于准入清单与评审输出模板做一次正式评审。
- 若候选领域对宿主运行时依赖重，需提前设定“准入失败回退”方案。

## 7. 推荐下一动作
- 以 `docs/52-M4_上下文_第二领域线准入与可复用骨架.md` 的检查清单与评审模板组织一次候选领域评估，优先比较 `security-permissions` 与 `release-store-review`。
- 将评估结果回填到 `bd` 与新的 Context Snapshot，再决定是否创建“第二领域线接入” bead。
