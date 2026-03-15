---
id: L2.wxt.manifest.permissions-review
level: L2
domain: wxt
subdomain: manifest
capability: null
node_kind: execution-strategy
visibility_scope: domain-scoped
activation_mode: direct
title: Manifest Permissions Review
summary: Execute permission review steps for WXT manifest and align with store policies.
aliases:
  - manifest permission review
triggers:
  - permissions review
anti_triggers: []
required_with:
  - L1.security.permissions
may_include: []
children: []
entry_conditions:
  - execution_needed
stop_conditions:
  - permission_review_done
---

## 执行步骤
1. 收集 manifest 中声明的 permissions 与 host_permissions。
2. 对照功能需求解释每个权限必要性。
3. 标记可移除或可替代权限。

## 配置点
- 目标浏览器与商店政策。

## 工具调用点
- 仅需要静态规则比对与人工检查。

## validator 前置
- 需要已有权限最小化策略草案。
