# wxt-manifest

`wxt-manifest` 是 M2 首个完整 workflow package，用于审查与修复 WXT 扩展的 manifest 相关风险。

## 目标
- 聚焦 manifest 生成规则、permissions、host permissions、browser overrides、store-facing 风险。
- 产出结构化主 artifact：`manifest-review.md`。
- 为 M3 提供稳定 validator 输入面与回归基线。

## 不在本包范围
- content script 注入架构设计。
- background runtime 运行时编排。
- 全量 browser API 适配与 UI/UX 设计。

## 进入条件
- 任务显式提到 manifest、permissions、host permissions 或 store submission。
- 路由命中 `L1.wxt.manifest`。

## 退出条件
- manifest 核心风险已结构化说明。
- 权限最小化建议与 store-facing 风险已给出。
- 当前包职责已完成，必要时 handoff 给相邻包。

## 交接边界
- handoff 到 `security-permissions`：当核心问题是权限最小化和敏感面审查。
- handoff 到 `release-store-review`：当核心问题是商店审核与发布前检查。
