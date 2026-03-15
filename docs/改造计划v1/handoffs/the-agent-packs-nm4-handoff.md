# Handoff: 改造计划 v1 M1 分层索引与总编排草案

## 1. 交接对象
- 来源 bead：the-agent-packs-nm4
- 下一 bead：待创建（建议 M1 实现前设计 bead 或 M2 注册表设计 bead）
- 来源里程碑：M1
- 目标角色：项目内部维护 agent / 迭代开发子 agent

## 2. 已完成什么
- 为 M1 开发指导补充了共享产物落点与给 M2 的最小输入要求。
- 为节点分类文档补充了 `node_kind / visibility_scope / activation_mode` 到注册表字段的映射建议。
- 为 Routing 分层文档补充了“当前实现 -> 目标模型”的差异摘要，以及实现层最小规则摘要。
- 新增了 M1 阶段快照：`docs/改造计划v1/context-snapshots/2026-03-15-m1-layered-routing.md`。

## 3. 下一位 agent 可直接依赖什么
- M1 开发指导：`docs/改造计划v1/20-M1_分层索引与总编排骨架_开发指导.md`
- 节点分类与字段映射：`docs/改造计划v1/21-M1_上下文_节点分类_作用域与可见性.md`
- 路由分层规则与实现差异摘要：`docs/改造计划v1/22-M1_上下文_Routing分层化与候选集裁剪.md`
- 阶段事实与风险：`docs/改造计划v1/context-snapshots/2026-03-15-m1-layered-routing.md`

## 4. 下一位 agent 必须先做什么
- 先 claim：`the-agent-packs-nm4` 或新建其后续 bead。
- 先阅读：`docs/改造计划v1/20-M1_分层索引与总编排骨架_开发指导.md`、`docs/改造计划v1/21-M1_上下文_节点分类_作用域与可见性.md`、`docs/改造计划v1/22-M1_上下文_Routing分层化与候选集裁剪.md`。
- 先验证：`go test ./...`。

## 5. 不要做什么
- 不要回头使用 `docs/handoffs/` 或 `docs/context-snapshots/` 作为 v1 的主协作入口。
- 不要在尚未明确实现 bead 的情况下直接改正式主线的顶层 envelope 或状态语义。
- 不要让 capability 节点重新回到第一轮主竞争空间。

## 6. 风险与未决项
- 当前只完成了 M1 设计文档补强，尚未进入实现型 bead。
- 若后续不及时衔接 M2，注册表字段需求仍可能只停留在文档层。

## 7. 推荐下一动作
- 若继续做 v1 M1，优先创建“分层 route 实现前差异清单” bead，对 `internal/query/query.go`、Blueprint frontmatter、schema 的改动面做受控拆分。
- 若切换到 v1 M2，直接消费 M1 给出的字段映射与候选空间规则，设计注册表最小字段与冲突裁决。
