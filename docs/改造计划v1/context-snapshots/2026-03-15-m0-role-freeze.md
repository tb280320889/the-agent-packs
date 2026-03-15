# Context Snapshot: 改造计划 v1 M0 角色冻结与边界校正

## 1. 当前阶段
- 所属里程碑：M0
- 关联 bead：the-agent-packs-5oq
- 当前状态：in_progress

## 2. 当前事实
- 当前要解决的问题：把 v1 初始化阶段已经写出的角色模型、冻结面、兼容面与禁止事项，收口成可被 M1 直接消费的正式 M0 交付。
- 当前已完成内容：已为 M0 开发指导补充共享产物落点与给 M1 的最小摘要；已为角色模型文档补充角色冻结摘要；已为冻结面文档补充给 M1 的冻结摘要与停止条件。
- 当前尚未完成内容：尚未关闭 M0 bead；尚需把 M1 与 M0 的依赖关系在交接文档中写实。

## 3. 已冻结对象
- 角色冻结：项目内部维护 agent、迭代开发子 agent、文档型 agent、用户、用户侧外部 agent。
- 顶层冻结面：Activation Request、Routing Result、Context Bundle、Artifact、Handoff Bundle、Validation Plan、Validator Result、Activation Result、状态枚举与优先级、workflow package 基本模板。
- 兼容演化面：Blueprint 节点分类字段、route 候选集筛选逻辑、package 注册表、package 命名与保留名规则、文档吸收治理流程。

## 4. 当前输入
- 上游交付物：`the-agent-packs-j03` 初始化文档包与初始化快照/交接。
- 依赖文档：`docs/改造计划v1/10-M0_角色冻结与边界校正_开发指导.md`、`docs/改造计划v1/11-M0_上下文_角色模型_职责与非目标.md`、`docs/改造计划v1/12-M0_上下文_冻结面_兼容面与禁止事项.md`。
- 依赖协作规则：`docs/改造计划v1/04-多Agent接力开发与bd协作规则_改造版.md`、`docs/改造计划v1/05-统一Handoff_ContextSnapshot_共享文档模板_改造版.md`。

## 5. 当前输出
- 已更新文件：
  - `docs/改造计划v1/10-M0_角色冻结与边界校正_开发指导.md`
  - `docs/改造计划v1/11-M0_上下文_角色模型_职责与非目标.md`
  - `docs/改造计划v1/12-M0_上下文_冻结面_兼容面与禁止事项.md`
- 已产出文件：
  - `docs/改造计划v1/context-snapshots/2026-03-15-m0-role-freeze.md`
- 已创建 bead：`the-agent-packs-5oq`

## 6. 风险与阻塞
- 风险：如果 M1 不明确引用 M0 的冻结摘要，后续仍可能把角色讨论和结构讨论混在一起。
- 风险：若新增实现型 bead 没有依赖 M0，可能再次出现“先做 M1、后补 M0”的顺序漂移。
- 阻塞：无。
- 是否需要 breaking 评估：否（当前为文档冻结与交接补齐，不改正式主线协议）。

## 7. 下一步建议
- 建议下一个 bead：`the-agent-packs-nm4`（在已补回 M0 后继续推进 M1）。
- 建议先阅读的文档：`docs/改造计划v1/20-M1_分层索引与总编排骨架_开发指导.md`、`docs/改造计划v1/21-M1_上下文_节点分类_作用域与可见性.md`、`docs/改造计划v1/22-M1_上下文_Routing分层化与候选集裁剪.md`。
- 建议先验证的命令：`go test ./...`、`bd show the-agent-packs-5oq --json`。
