# Context Snapshot: 改造计划 v1 M1 分层索引与总编排草案

## 1. 当前阶段
- 所属里程碑：M1
- 关联 bead：the-agent-packs-nm4
- 当前状态：completed

## 2. 当前事实
- 当前要解决的问题：把当前平铺式 route 抽象成“总编排 -> 主域编排 -> 横线挂接”的两段式候选裁剪模型，并给 M2 留下注册表字段需求。
- 当前已完成内容：已补充 M1 文档中的共享产物要求、给 M2 的最小输入要求、节点分类到注册表字段的映射建议，以及从当前实现到目标模型的差异摘要；M1 现已显式依赖 M0 冻结摘要；并已在 Blueprint frontmatter、compiler schema、query candidate 裁剪与测试中落地最小结构实现。
- 当前尚未完成内容：尚未拆出更细的 M1 后续实现 bead；当前阶段已完成最小实现，但尚未进入更大范围 route 重构。

## 3. 已冻结对象
- M0 冻结对象不变：Activation Request、Routing Result、Context Bundle、Artifact、Handoff Bundle、Validation Plan、Validator Result、Activation Result。
- M1 结构判断：第一轮先做主域筛选，第二轮再做域内 workflow 选择，横线能力默认 attach-only。

## 4. 当前输入
- 上游交付物：`the-agent-packs-j03` 初始化文档包、`the-agent-packs-5oq` 的 M0 冻结交付。
- 依赖文档：`docs/改造计划v1/20-M1_分层索引与总编排骨架_开发指导.md`、`docs/改造计划v1/21-M1_上下文_节点分类_作用域与可见性.md`、`docs/改造计划v1/22-M1_上下文_Routing分层化与候选集裁剪.md`。
- 依赖实现：`internal/query/query.go` 当前 route 实现、`internal/compiler/compiler.go` 索引编译实现、`blueprint/schema.md` 当前索引 schema。

## 5. 当前输出
- 已更新文件：
  - `docs/改造计划v1/20-M1_分层索引与总编排骨架_开发指导.md`
  - `docs/改造计划v1/21-M1_上下文_节点分类_作用域与可见性.md`
  - `docs/改造计划v1/22-M1_上下文_Routing分层化与候选集裁剪.md`
  - `blueprint/schema.md`
  - `blueprint/frontmatter-examples.md`
  - `blueprint/L0/wxt/overview.md`
  - `blueprint/L0/security/overview.md`
  - `blueprint/L0/release/overview.md`
  - `blueprint/L1/wxt/manifest.md`
  - `blueprint/L1/security/permissions.md`
  - `blueprint/L1/release/store-review.md`
  - `internal/model/model.go`
  - `internal/compiler/compiler.go`
  - `internal/query/query.go`
  - `tests/m1_minimal_test.go`
  - `tests/m3_validation_closure_test.go`
- 已产出文件：
  - `docs/改造计划v1/context-snapshots/2026-03-15-m1-layered-routing.md`
- 已创建 bead：`the-agent-packs-nm4`（已完成并关闭）

## 6. 风险与阻塞
- 风险：若没有实现型 bead 继续承接，M1 仍可能停留在文档层设计，无法验证与现有 query 实现的映射精度。
- 风险：若 M2 不直接消费 M1 的字段映射，注册表字段仍可能重复设计。
- 风险：若后续 agent 跳过 `the-agent-packs-5oq` 的 M0 交接，仍可能误把结构设计问题升级为角色或协议重定义。
- 阻塞：无。
- 是否需要 breaking 评估：否（当前仍是结构设计与交接补强，不改正式顶层协议）。

## 7. 下一步建议
- 建议下一个 bead：v1 M2 package 注册表与命名空间治理 bead。
- 建议先阅读的文档：`docs/改造计划v1/30-M2_package注册表与命名空间治理_开发指导.md`、`docs/改造计划v1/31-M2_上下文_命名规则_注册表字段与冲突裁决.md`。
- 建议先验证的命令：`go test ./...`、`bd show the-agent-packs-nm4 --json`。
