# Context Snapshot: 改造计划 v1 M2 package 注册表与命名空间治理

## 1. 当前阶段
- 所属里程碑：M2
- 关联 bead：the-agent-packs-bsj
- 当前状态：completed

## 2. 当前事实
- 当前要解决的问题：把 workflow package 的身份判断从目录约定提升为正式注册表规则，冻结最小字段、命名空间、保留名与冲突裁决流程，并落地为机器可校验的实现。
- 当前已完成内容：已补充 M2 开发指导中的输入边界、共享产物要求、设计顺序、非目标与向后续里程碑输出的最小要求；已补充 M2 上下文文档中的 canonical name 规则、alias 边界、主域/capability/orchestrator 命名空间判断、字段语义、裁决顺序、结果类型与样例判断；已新增 `workflow-packages/registry.json` 作为注册表真相源，补齐 capability 样本包 `security-permissions`、`release-store-review`，并在 `internal/registry/registry.go` 与 `internal/query/query.go` 中落地注册表加载、准入校验与 route/bundle 消费逻辑，同时补充 `tests/m2_registry_test.go` 回归测试。
- 当前尚未完成内容：无；M2 文档层与实现层已收口，后续工作转入 M3 文档吸收治理。

## 3. 已冻结对象
- M0 冻结对象不变：Activation Request、Routing Result、Context Bundle、Artifact、Handoff Bundle、Validation Plan、Validator Result、Activation Result、状态枚举与优先级、workflow package 基本模板。
- M1 结构判断延续：候选空间至少分为 `global / domain / capability`，横线能力默认 `attach-only`。
- M2 注册表判断：package 身份由 canonical name 与注册表最小字段共同定义，alias 不能替代身份真相。

## 4. 当前输入
- 上游交付物：`the-agent-packs-nm4` 的 M1 分层索引文档、快照与 handoff。
- 依赖文档：`docs/改造计划v1/30-M2_package注册表与命名空间治理_开发指导.md`、`docs/改造计划v1/31-M2_上下文_命名规则_注册表字段与冲突裁决.md`、`docs/改造计划v1/21-M1_上下文_节点分类_作用域与可见性.md`、`docs/改造计划v1/22-M1_上下文_Routing分层化与候选集裁剪.md`。
- 依赖实现：`workflow-packages/wxt-manifest/package.yaml`、`workflow-packages/wxt-manifest/README.md`、`workflow-packages/README.md`。

## 5. 当前输出
- 已更新文件：
  - `docs/改造计划v1/30-M2_package注册表与命名空间治理_开发指导.md`
  - `docs/改造计划v1/31-M2_上下文_命名规则_注册表字段与冲突裁决.md`
  - `workflow-packages/README.md`
  - `internal/query/query.go`
- 已产出文件：
  - `docs/改造计划v1/context-snapshots/2026-03-15-m2-package-registry.md`
  - `workflow-packages/registry.json`
  - `workflow-packages/security-permissions/package.yaml`
  - `workflow-packages/security-permissions/README.md`
  - `workflow-packages/release-store-review/package.yaml`
  - `workflow-packages/release-store-review/README.md`
  - `internal/registry/registry.go`
  - `tests/m2_registry_test.go`
- 已创建 bead：`the-agent-packs-bsj`

## 6. 风险与阻塞
- 风险：若 M3 在吸收外部资料时绕过注册表直接生成新包或游离资产，会重新引入命名与归属漂移。
- 风险：若后续实现层再引入硬编码 package 映射而不消费注册表，会造成双真相源。
- 风险：当前注册表校验聚焦现有样本与最小字段，后续新增包仍需持续补充样本覆盖。
- 阻塞：无。
- 是否需要 breaking 评估：否（当前仍为结构层规则冻结，不改顶层协议）。

## 7. 下一步建议
- 建议下一个 bead：M3 文档吸收管道与知识资产化 bead。
- 建议先阅读的文档：`docs/改造计划v1/40-M3_文档吸收管道与知识资产化_开发指导.md`、`docs/改造计划v1/41-M3_上下文_外部资料纳入_语义映射与资产落位.md`、`docs/改造计划v1/31-M2_上下文_命名规则_注册表字段与冲突裁决.md`。
- 建议先验证的命令：`go test ./...`、`bd show the-agent-packs-bsj --json`。
