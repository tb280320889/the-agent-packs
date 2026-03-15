---
id: L2.security.permissions.minimization
level: L2
domain: security
subdomain: permissions
capability: null
title: Permissions Minimization
summary: Provide concrete steps to minimize permissions while preserving required functionality.
aliases:
  - permission minimization
triggers:
  - minimize permissions
anti_triggers: []
required_with: []
may_include: []
children: []
entry_conditions:
  - execution_needed
stop_conditions:
  - minimization_plan_ready
---

## 执行步骤
1. 列出高风险权限与对应功能。
2. 评估替代方案与降级路径。
3. 形成最小权限清单。

## 配置点
- 目标功能与安全策略约束。

## 工具调用点
- 无强制工具，使用规则清单即可。

## validator 前置
- 需要已识别的权限需求清单。
