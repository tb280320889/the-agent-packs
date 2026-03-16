# 30 M2 package 注册表与命名空间治理 开发指导

## 文档类型
开发里程碑开发指导文档

## 作用域
M2 负责设计 workflow package 注册表、命名空间与冲突治理规则。

## 本里程碑目标
1. 定义 package 注册表为身份真相源
2. 冻结 package 命名规则
3. 冻结保留名与裸名禁止规则
4. 定义冲突裁决与索引归属判断

## 本里程碑输入边界
- 只消费共享文档 `00~05`、M0 冻结面、M1 分层 route 输出、当前 M2 文档。
- 直接承接 M1 已冻结的最小结构判断：`global / domain / capability` 候选空间、`visibility_scope`、`activation_mode`、`canonical_blueprint_node`。
- 不为兼容单个现有 package 放宽全局规则；单包样本只能用于校验，不得反向定义系统边界。

## 必须产出
1. 注册表最小字段
2. package 命名规则
3. 保留名清单
4. 冲突裁决流程

## 本里程碑建议直接落地的共享产物
- M2 阶段快照：`docs/改造计划v1/context-snapshots/<date>-m2-package-registry.md`
- M2 交接文档：`docs/改造计划v1/handoffs/<task-id>-handoff.md`
- M2 输出摘要：在文档中明确“字段语义 / 命名空间 / 裁决顺序 / 对 route 的消费方式”，避免 M3/M4 再回头猜。

## 设计顺序
1. 先冻结注册表最小字段与字段语义。
2. 再冻结 canonical name、alias、reserved name 的边界。
3. 再定义主域 package 与 capability package 的归属判断。
4. 最后定义新增 package 的冲突裁决与准入流程。

## M2 必须回答的问题
- 一个 package 的身份真相由哪些字段构成，哪些字段只是辅助检索。
- 一个名称是否可注册，是否必须改名，是否只能作为 alias 或 reserved 词。
- 一个 package 属于 `orchestrator / workflow / capability` 中哪一类，以及它出现在哪一层候选空间。
- 一个 `attach-only` capability 如何在注册表里被表达，并被 route 与后续迁移同时消费。

## 明确不产出
- 顶层 envelope 或状态枚举的重定义。
- 绕过注册表直接新增裸名 package 的例外规则。
- M3 文档吸收管道与资产落位细节。

## 完成标准
1. 后续新增 package 时可以自动判断是否冲突
2. `release/security` 一类高复用词不会再以裸名出现
3. package 的分类与路由空间可由注册表直接判断
4. M3/M4 可以直接依赖注册表判断 package 归属、可见性与 attach-only 语义

## 本阶段已落地实现对齐点
- 注册表真相源文件：`workflow-packages/registry.json`
- capability 样本包：`workflow-packages/security-permissions/package.yaml`、`workflow-packages/release-store-review/package.yaml`
- 注册表加载与校验实现：`internal/registry/registry.go`
- route / bundle 对注册表的消费：`internal/query/query.go`
- 回归测试：`tests/m2_registry_test.go`

## 给后续里程碑的最小输出要求
- 对 M3：外部资料纳入后，必须能挂回某个已注册的 `domain / capability / package` 命名空间，不允许生成游离资产。
- 对 M4：迁移与准入演练必须先经过注册表校验，再进入 route / validator 兼容验证。
