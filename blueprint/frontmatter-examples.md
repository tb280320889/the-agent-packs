# Blueprint Frontmatter 示例

本文件提供合法/非法示例，辅助校验 M1 的最小 frontmatter 规范。

## 合法示例

```yaml
id: L1.wxt.manifest
level: L1
domain: wxt
subdomain: manifest
capability: null
node_kind: workflow-entry
visibility_scope: domain-scoped
activation_mode: direct
title: WXT Manifest
summary: Review manifest generation rules, browser overrides, permissions and store-facing risks.
aliases:
  - web extension manifest
  - wxt config
triggers:
  - manifest
  - permissions
  - host permissions
anti_triggers:
  - tauri
  - telegram miniapp
required_with:
  - L1.security.permissions
  - L1.release.store-review
may_include:
  - L2.wxt.manifest.permissions-review
children:
  - L2.wxt.manifest.permissions-review
entry_conditions:
  - browser_extension_host_confirmed
stop_conditions:
  - execution_plan_possible
```

## 非法示例 A（路径与 id 不一致）

```yaml
id: L1.browser.manifest
level: L1
domain: wxt
subdomain: manifest
```

## 非法示例 B（字段语义漂移）

```yaml
required_with:
  - maybe security
```

## 非法示例 C（summary 过长）

summary 只能是可直接进入 bundle 的摘要，不可替代正文。
