---
id: L3.wxt.manifest.edge-cases
level: L3
domain: wxt
subdomain: manifest
capability: null
title: Manifest Edge Cases
summary: Capture edge cases for WXT manifest handling across browsers and store rules.
aliases:
  - manifest edge cases
triggers:
  - edge cases
anti_triggers: []
required_with: []
may_include: []
children: []
entry_conditions:
  - validator_failure
stop_conditions:
  - edge_case_resolved
---

## edge cases
- 浏览器版本差异导致字段无效。
- 权限声明与实际功能不一致。

## 平台差异
- 不同浏览器对同一权限的限制不同。

## 深度排障
- 对照商店拒绝原因回溯 manifest 字段。

## 回退路径
- 回退到最小权限与最小字段集合。
