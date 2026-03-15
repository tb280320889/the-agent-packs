# Handoff: 改造计划 v1 M2 package 注册表与命名空间治理

## 1. 交接对象
- 来源 bead：the-agent-packs-bsj
- 下一 bead：待创建（建议 M2 实现层注册表准入 bead）
- 来源里程碑：M2
- 目标角色：项目内部维护 agent / 迭代开发子 agent / Verifier

## 2. 已完成什么
- 为 M2 开发指导补充了输入边界、共享产物落点、设计顺序、必须回答的问题、非目标与给后续里程碑的最小输出要求。
- 为 M2 上下文文档补充了 canonical name、alias、主域/capability/orchestrator 命名空间判断、保留名解释、字段语义、冲突裁决顺序与结果类型。
- 基于 `wxt-manifest` 与 `security-permissions` / `release-store-review` 的现实样本，明确了主域 package 与 capability package 的命名区分方式。
- 新增 `workflow-packages/registry.json` 作为注册表真相源，并补齐 `security-permissions`、`release-store-review` capability 包样本目录。
- 新增 `internal/registry/registry.go` 实现注册表加载、字段校验、保留名约束、alias 冲突检查与 package manifest 对齐验证。
- 已把 `internal/query/query.go` 改为消费注册表，不再依赖平铺硬编码映射；并新增 `tests/m2_registry_test.go` 做 M2 回归覆盖。
- 新增了 M2 阶段快照：`docs/改造计划v1/context-snapshots/2026-03-15-m2-package-registry.md`。

## 3. 下一位 agent 可直接依赖什么
- M2 开发指导：`docs/改造计划v1/30-M2_package注册表与命名空间治理_开发指导.md`
- M2 命名与冲突规则：`docs/改造计划v1/31-M2_上下文_命名规则_注册表字段与冲突裁决.md`
- M2 注册表真相源：`workflow-packages/registry.json`
- M2 注册表实现：`internal/registry/registry.go`
- M1 字段映射与候选空间规则：`docs/改造计划v1/21-M1_上下文_节点分类_作用域与可见性.md`、`docs/改造计划v1/22-M1_上下文_Routing分层化与候选集裁剪.md`
- 阶段事实与风险：`docs/改造计划v1/context-snapshots/2026-03-15-m2-package-registry.md`

## 4. 下一位 agent 必须先做什么
- 先 claim：新建并认领 M3 bead。
- 先阅读：`docs/改造计划v1/40-M3_文档吸收管道与知识资产化_开发指导.md`、`docs/改造计划v1/41-M3_上下文_外部资料纳入_语义映射与资产落位.md`、`docs/改造计划v1/30-M2_package注册表与命名空间治理_开发指导.md`、`docs/改造计划v1/31-M2_上下文_命名规则_注册表字段与冲突裁决.md`。
- 先验证：`go test ./...`。

## 5. 不要做什么
- 不要为了兼容单个历史 package 放宽 canonical name、裸名禁令或 attach-only 规则。
- 不要让 alias 成为身份真相源，也不要让 alias 绕过候选空间裁剪。
- 不要绕过 `workflow-packages/registry.json` 直接新增 package 或生成未归属资产。

## 6. 风险与未决项
- 当前 M2 已完成规则冻结与最小实现落地；未决项已转移为 M3 如何让外部资料稳定挂回已注册命名空间。
- 若 M3 直接复制外部资料而不先做语义映射和注册表归属判断，会破坏 M2 已冻结的命名治理。

## 7. 推荐下一动作
- 直接进入 M3：围绕外部资料纳入、语义映射、资产落位设计吸收管道，并要求所有新资产先绑定已注册 package / capability / domain 命名空间。
