---
id: L1.release.store-review
level: L1
domain: release
subdomain: store-review
capability: null
node_kind: workflow-entry
visibility_scope: capability-scoped
activation_mode: attach-only
title: Store Review
summary: Prepare browser extension submissions by aligning manifest, permissions, and policy checks.
aliases:
  - store review
  - submission review
triggers:
  - store review
  - submission
  - release
anti_triggers: []
required_with: []
may_include:
  - L2.release.store-review.browser-extension-checklist
children:
  - L2.release.store-review.browser-extension-checklist
entry_conditions:
  - store_submission_present
stop_conditions:
  - checklist_ready
---

## 子领域目标
- 对齐商店审核要求并形成检查清单。

## 输入输出
- 输入：目标商店、manifest 与权限信息。
- 输出：审核清单与风险点。

## 依赖
- 无强依赖。

## 常见分支
- 不同商店政策差异。
- 权限或描述不匹配。

## 升级条件
- 进入执行态时补 L2 审核策略。

## 停止条件
- 可执行的审核清单已形成。
