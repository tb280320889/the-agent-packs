# wxt-manifest 回退策略

## partial 条件
- 域已明确（WXT）但缺少关键上下文（如 manifest 权限片段）。

## handoff 条件
- manifest 审查完成，剩余风险需 `security-permissions` 深审。
- manifest 审查完成，剩余风险需 `release-store-review` 发布前校验。

## failed 条件
- activation request 关键字段缺失。
- 无法形成可执行 route 且无法进入 partial/handoff。
