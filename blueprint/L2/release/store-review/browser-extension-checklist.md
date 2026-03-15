---
id: L2.release.store-review.browser-extension-checklist
level: L2
domain: release
subdomain: store-review
capability: null
title: Browser Extension Store Checklist
summary: Build a checklist covering manifest, permissions, and policy alignment for store submission.
aliases:
  - store checklist
triggers:
  - checklist
  - store review
anti_triggers: []
required_with: []
may_include: []
children: []
entry_conditions:
  - execution_needed
stop_conditions:
  - checklist_ready
---

## 执行步骤
1. 列出商店审核要求与禁止项。
2. 对照 manifest 与权限声明逐条检查。
3. 输出可执行的提交前清单。

## 配置点
- 目标商店的最新政策。

## 工具调用点
- 无强制工具，使用清单模板。

## validator 前置
- 需要已收集的 manifest 与权限信息。
