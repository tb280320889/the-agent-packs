# security-permissions

`security-permissions` 是横线 capability package，用于在主域 workflow 确认后附挂权限最小化与敏感面审查。

## 目标
- 审查权限范围是否过宽。
- 识别敏感权限与风险暴露面。
- 为主包提供权限最小化建议。

## 边界
- 不作为第一轮主竞争入口。
- 默认由主域包或 orchestrator 通过 attach 方式带出。
