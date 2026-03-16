# Handoff: 改造计划 v1 M3 文档吸收管道与知识资产化

## 1. 交接对象
- 来源 GSD 任务项：the-agent-packs-s91
- 下一 GSD 任务项：the-agent-packs-2ji
- 来源里程碑：M3
- 目标角色：项目内部维护 agent / 文档型 agent / Verifier

## 2. 已完成什么
- 为 M3 开发指导补充了输入边界、共享产物、设计顺序、必须回答的问题、非目标、完成标准、最小冻结对象与工作流要求。
- 为 M3 上下文文档补充了候选资料最小元数据、元数据约束、标准吸收流程、语义映射规则、目录与对象落位约束、职责边界与审查闸口。
- 明确所有正式资产都必须绑定到已注册的 `system / domain / capability / package` 命名空间，无法映射者只能进入 backlog 或 reject。
- 新增了 M3 阶段快照：`docs/改造计划v1/context-snapshots/2026-03-15-m3-knowledge-ingestion.md`。

## 3. 下一位 agent 可直接依赖什么
- M3 开发指导：`docs/改造计划v1/40-M3_文档吸收管道与知识资产化_开发指导.md`
- M3 上下文：`docs/改造计划v1/41-M3_上下文_外部资料纳入_语义映射与资产落位.md`
- M3 阶段快照：`docs/改造计划v1/context-snapshots/2026-03-15-m3-knowledge-ingestion.md`
- M2 注册表真相源：`workflow-packages/registry.json`
- M2 注册表实现：`internal/registry/registry.go`

## 4. 下一位 agent 必须先做什么
- 先在 GSD 任务记录中认领任务项：新建并认领 M4 任务项。
- 先阅读：`docs/改造计划v1/50-M4_迁移实施_兼容验证与准入演练_开发指导.md`、`docs/改造计划v1/51-M4_上下文_迁移步骤_回滚策略与验收清单.md`、`docs/改造计划v1/40-M3_文档吸收管道与知识资产化_开发指导.md`、`docs/改造计划v1/41-M3_上下文_外部资料纳入_语义映射与资产落位.md`。
- 先验证：`go test ./...`。

## 5. 不要做什么
- 不要把 backlog 候选直接当正式资产落位。
- 不要让外部资料绕过注册表新增主域、capability line 或 package。
- 不要让 package README 承担系统层或 capability 层真相源。

## 6. 风险与未决项
- 若 M4 只验证“文件存在”而不验证 `source_uri -> target_scope -> target_ref -> 落位` 链路，M3 的治理价值会被削弱。
- 若后续领域扩展先收资料再补注册表，会重新引入命名漂移与归属漂移。

## 7. 推荐下一动作
- 在 M4 中把 M2 注册表校验与 M3 资产归属校验合并进迁移实施与准入演练，形成“先身份、后资料、再兼容验证”的闭环。
