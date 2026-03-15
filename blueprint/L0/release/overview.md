---
id: L0.release
level: L0
domain: release
subdomain: null
capability: null
node_kind: capability-root
visibility_scope: capability-scoped
activation_mode: attach-only
title: Release
summary: Enter when the task涉及商店审核、发布与合规检查；作为 cross-cutting line 参与其他主域。
aliases:
  - release
  - store review
triggers:
  - store review
  - submission
  - release
anti_triggers: []
required_with: []
may_include:
  - L1.release.store-review
children:
  - L1.release.store-review
entry_conditions:
  - store_submission_present
stop_conditions:
  - subdomain_selected
---

## 领域定义
- 发布与商店审核相关的横线能力。

## 何时进入
- 主域涉及浏览器扩展发布或商店审核。

## 何时不要进入
- 任务与发布/审核无关。

## 必带横线
- 无（本身为横线）。

## 推荐下一跳
- L1.release.store-review
