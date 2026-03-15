---
id: L1.security.permissions
level: L1
domain: security
subdomain: permissions
capability: null
node_kind: workflow-entry
visibility_scope: capability-scoped
activation_mode: attach-only
title: Permissions Review
summary: Assess permission scope, minimize sensitive capabilities, and surface security risks for store review.
aliases:
  - permissions review
  - security permissions
triggers:
  - permissions
  - host permissions
  - sensitive access
anti_triggers: []
required_with: []
may_include:
  - L2.security.permissions.minimization
children:
  - L2.security.permissions.minimization
entry_conditions:
  - permission_scope_present
stop_conditions:
  - review_plan_ready
---

## 子领域目标
- 评估权限范围并给出最小化策略。

## 输入输出
- 输入：权限声明、访问范围、目标功能描述。
- 输出：权限缩减建议与风险提示。

## 依赖
- 无强依赖。

## 常见分支
- 高风险权限需强解释。
- 权限不足导致功能受限。

## 升级条件
- 进入执行态时补 L2 细化步骤。

## 停止条件
- 权限最小化结论已具备。
