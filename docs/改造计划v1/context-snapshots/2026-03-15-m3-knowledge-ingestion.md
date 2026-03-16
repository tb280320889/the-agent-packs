# Context Snapshot: 改造计划 v1 M3 文档吸收管道与知识资产化

## 1. 当前阶段
- 所属里程碑：M3
- 关联 GSD 任务项：the-agent-packs-s91
- 当前状态：completed

## 2. 当前事实
- 当前要解决的问题：把外部资料吸收从“收集一些文档”收口为可执行治理管道，明确候选池、元数据、语义映射、审查闸口与资产落位规则，避免资料进入仓库后失去归属与可消费性。
- 当前已完成内容：已补充 M3 开发指导中的输入边界、共享产物、设计顺序、必须回答问题、非目标、最小冻结对象、文档型 agent 工作流要求与对 M4 的输出要求；已补充 M3 上下文文档中的候选资料最小元数据、元数据约束、标准吸收流程、语义映射规则、目录与对象落位约束、文档型 agent 职责边界、审查闸口与执行层约束。
- 当前尚未完成内容：无；M3 文档层规则已收口，后续工作转入 M4 迁移实施、兼容验证与准入演练。

## 3. 已冻结对象
- M2 注册表与命名治理延续：所有正式资产都必须绑定已注册的 `domain / capability / package` 命名空间。
- M3 候选资料最小元数据：`source_type`、`source_uri`、`captured_at`、`owner`、`trust_level`、`target_scope`、`target_ref`、`disposition`。
- M3 处理状态最小集合：`candidate`、`mapped`、`frozen`、`backlog`、`rejected`。
- M3 审查闸口：命名空间校验、语义去重、非目标检查、落位检查、可消费性检查。

## 4. 当前输入
- 上游交付物：`docs/改造计划v1/context-snapshots/2026-03-15-m2-package-registry.md`、`docs/改造计划v1/handoffs/M3-ENTRY-from-M2.md`。
- 依赖文档：`docs/改造计划v1/40-M3_文档吸收管道与知识资产化_开发指导.md`、`docs/改造计划v1/41-M3_上下文_外部资料纳入_语义映射与资产落位.md`、`docs/改造计划v1/31-M2_上下文_命名规则_注册表字段与冲突裁决.md`。
- 依赖实现：`workflow-packages/registry.json`、`internal/registry/registry.go`。

## 5. 当前输出
- 已更新文件：
  - `docs/改造计划v1/40-M3_文档吸收管道与知识资产化_开发指导.md`
  - `docs/改造计划v1/41-M3_上下文_外部资料纳入_语义映射与资产落位.md`
- 已产出文件：
  - `docs/改造计划v1/context-snapshots/2026-03-15-m3-knowledge-ingestion.md`
  - `docs/改造计划v1/handoffs/the-agent-packs-s91-handoff.md`
- 已创建 GSD 任务项：`the-agent-packs-s91`

## 6. 风险与阻塞
- 风险：若后续执行层只保留流程名称，不把候选元数据与审查闸口做成硬性准入点，M3 仍可能退化为“资料堆放规范”。
- 风险：若新增主域或 capability line 时先吸收资料后补注册表，会破坏 M2 的身份真相源。
- 风险：若 backlog 与正式资产边界不清，后续迁移时会出现不可追溯资产混入正式知识层。
- 阻塞：无。
- 是否需要 breaking 评估：否（当前为治理层规则冻结，不改顶层协议）。

## 7. 下一步建议
- 建议下一个 GSD 任务项：M4 迁移实施、兼容验证与准入演练任务项。
- 建议先阅读的文档：`docs/改造计划v1/50-M4_迁移实施_兼容验证与准入演练_开发指导.md`、`docs/改造计划v1/51-M4_上下文_迁移步骤_回滚策略与验收清单.md`、`docs/改造计划v1/40-M3_文档吸收管道与知识资产化_开发指导.md`。
- 建议先验证的命令：`go test ./...`、确认 GSD 任务记录状态。
