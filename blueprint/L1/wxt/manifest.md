---
id: L1.wxt.manifest
level: L1
domain: wxt
subdomain: manifest
capability: null
title: WXT Manifest
summary: Review WXT manifest generation rules, browser overrides, permissions, and store-facing risks.
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
  - L2.wxt.manifest.browser-overrides
children:
  - L2.wxt.manifest.permissions-review
  - L2.wxt.manifest.browser-overrides
entry_conditions:
  - browser_extension_host_confirmed
stop_conditions:
  - execution_plan_possible
---

## 子领域目标
- 明确 manifest 生成与浏览器覆盖规则。

## 输入输出
- 输入：manifest 配置片段、目标浏览器、权限声明。
- 输出：权限审查结论与必要调整建议。

## 依赖
- 必带：L1.security.permissions、L1.release.store-review。

## 常见分支
- 权限过宽需要最小化。
- 浏览器差异需要 overrides。

## 升级条件
- 进入执行态时补 L2 执行策略。

## 停止条件
- 已形成可执行的审查与调整计划。
