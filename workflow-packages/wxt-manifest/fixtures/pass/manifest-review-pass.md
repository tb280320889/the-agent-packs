# Manifest Review

## 1. 当前任务摘要
- review WXT manifest permissions for browser store submission

## 2. main pack / required packs
- main pack: wxt-manifest
- required packs:
  - security-permissions
  - release-store-review

## 3. manifest 相关风险摘要
- 权限范围过宽会导致商店审核失败。

## 4. permissions / host permissions 审查
- 建议仅保留必要 host 权限并按域名最小化。

## 5. browser-specific overrides 提示
- Firefox 对某些字段要求更严格，需单独覆盖。

## 6. store-facing 风险与建议
- 提交前补齐权限用途说明与最小化理由。

## 7. validator 预期项
- validator-core-output
- validator-domain-wxt-manifest

## 8. 下一步建议 / 可能 handoff
- 如需细化权限面，handoff 到 security-permissions。
