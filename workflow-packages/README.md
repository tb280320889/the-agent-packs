# Workflow Packages

本目录是 workflow package 的统一根目录。

## 目录约定
- 每个 package 固定放在 `workflow-packages/<package>/`。
- 当前首个完整包：`workflow-packages/wxt-manifest/`。
- M2 起统一注册表位于 `workflow-packages/registry.json`，作为 package 身份、命名空间与准入判断的真相源。

## 一致性规则
- 后续里程碑新增 package 时，保持同一层级，不在仓库根目录分散创建。
- 所有 handoff、snapshot、测试与文档引用均使用该根路径。
- 新增 package 必须先通过注册表校验，再进入 Blueprint、validator 与 handoff 设计。
