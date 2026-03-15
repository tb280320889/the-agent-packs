# wxt-manifest 指令

## 你要做什么
- 审查 WXT manifest 配置是否满足最小权限与商店审核要求。
- 输出 `manifest-review.md`，结构必须符合模板。

## 核心检查项
- manifest 生成规则是否清晰且可追溯。
- permissions 与 host_permissions 是否最小化。
- browser-specific overrides 是否覆盖目标浏览器差异。
- CSP 与 store-facing 风险是否被明确标注。

## 输出要求
- 使用 `templates/docs/manifest-review.md`。
- 必须包含：任务摘要、主/必带包、风险、审查结论、下一步建议。
