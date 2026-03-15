# release-store-review

`release-store-review` 是横线 capability package，用于在主域 workflow 确认后附挂商店审核与发布前检查。

## 目标
- 对齐目标商店的审核要求。
- 汇总发布前检查清单。
- 补充 store-facing 风险说明。

## 边界
- 不作为第一轮主竞争入口。
- 默认由主域包或 orchestrator 通过 attach 方式带出。
