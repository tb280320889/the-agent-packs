# Handoff: M3 入口（承接 M2 注册表与命名空间治理）

## 1. 交接对象
- 来源 GSD 任务项：the-agent-packs-bsj / the-agent-packs-3eq
- 下一 GSD 任务项：the-agent-packs-s91
- 来源里程碑：M2
- 目标角色：项目内部维护 agent / 文档型 agent / 迭代开发子 agent

## 2. 已完成什么
- M2 已冻结并实现 package 注册表真相源：`workflow-packages/registry.json`。
- 已建立注册表准入校验：`internal/registry/registry.go`。
- 已把 route / bundle 的 package 映射改为消费注册表：`internal/query/query.go`。
- 已补齐 capability 样本包：`workflow-packages/security-permissions/`、`workflow-packages/release-store-review/`。
- 已通过回归测试：`go test ./...`。

## 3. 下一位 agent 可直接依赖什么
- M2 开发指导：`docs/改造计划v1/30-M2_package注册表与命名空间治理_开发指导.md`
- M2 命名与冲突规则：`docs/改造计划v1/31-M2_上下文_命名规则_注册表字段与冲突裁决.md`
- M2 阶段快照：`docs/改造计划v1/context-snapshots/2026-03-15-m2-package-registry.md`
- M2 正式交接：`docs/改造计划v1/handoffs/the-agent-packs-bsj-handoff.md`
- 注册表真相源与实现：`workflow-packages/registry.json`、`internal/registry/registry.go`

## 4. 下一位 agent 必须先做什么
- 先在 GSD 任务记录中认领任务项：创建并认领 M3 任务项。
- 先阅读：`docs/改造计划v1/40-M3_文档吸收管道与知识资产化_开发指导.md`、`docs/改造计划v1/41-M3_上下文_外部资料纳入_语义映射与资产落位.md`、`docs/改造计划v1/31-M2_上下文_命名规则_注册表字段与冲突裁决.md`。
- 先验证：`go test ./...`。

## 5. 不要做什么
- 不要把外部资料直接复制进主文档或 workflow package 目录而不做语义映射。
- 不要跳过注册表归属判断，生成无法挂回 `domain / capability / package` 命名空间的资产。
- 不要重新发明 package 命名规则、保留名或 attach-only 语义。

## 6. 风险与未决项
- M3 的主要风险不是缺资料，而是资料进入仓库后失去归属与可消费性。
- 若 M3 不显式消费注册表，会重新出现“资料有了，但不知道属于哪个 package / capability”的问题。

## 7. 推荐下一动作
- 在 M3 中把“候选收集 -> 语义归类 -> 文档映射 -> 冻结审查 -> 资产落位”明确落成管道，并要求所有落位结果绑定到 M2 注册表中的已注册命名空间。
