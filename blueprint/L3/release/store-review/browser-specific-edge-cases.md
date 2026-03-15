---
id: L3.release.store-review.browser-specific-edge-cases
level: L3
domain: release
subdomain: store-review
capability: null
title: Store Review Browser Edge Cases
summary: Capture browser-specific store review edge cases for extension submissions.
aliases:
  - store edge cases
triggers:
  - store rejection
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
- 商店政策更新导致已通过的配置失效。
- 描述与权限不匹配触发审核失败。

## 平台差异
- 不同商店对同类权限的接受度不同。

## 深度排障
- 逐条比对拒绝原因与提交内容。

## 回退路径
- 简化权限并补齐审核说明。
