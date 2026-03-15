---
id: L0.security
level: L0
domain: security
subdomain: null
capability: null
node_kind: capability-root
visibility_scope: capability-scoped
activation_mode: attach-only
title: Security
summary: Enter when a task涉及权限、敏感数据或合规风险；作为 cross-cutting line 参与其他主域。
aliases:
  - security
  - permissions
triggers:
  - permission
  - permissions
  - security
anti_triggers: []
required_with: []
may_include:
  - L1.security.permissions
children:
  - L1.security.permissions
entry_conditions:
  - security_risk_present
stop_conditions:
  - subdomain_selected
---

## 领域定义
- 安全与权限相关的横线能力。

## 何时进入
- 主域涉及权限申请、敏感能力或合规风险。

## 何时不要进入
- 任务与安全/权限无关。

## 必带横线
- 无（本身为横线）。

## 推荐下一跳
- L1.security.permissions
