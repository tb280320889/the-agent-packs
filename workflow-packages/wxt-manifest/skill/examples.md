# wxt-manifest 示例

## 示例任务
`review WXT manifest permissions for browser store submission`

## 预期主结果
- main pack: `wxt-manifest`
- required packs: `security-permissions`, `release-store-review`
- artifact: `manifest-review.md`

## partial 示例
- 输入缺 `manifest.permissions` 片段，仅能输出局部结论并请求补充上下文。

## handoff 示例
- 权限风险需深入审查：handoff 到 `security-permissions`。
- 商店提交约束需细化检查：handoff 到 `release-store-review`。
