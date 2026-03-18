---
id: L0.monorepo
level: L0
domain: monorepo
subdomain: null
capability: null
node_kind: domain-root
visibility_scope: global
activation_mode: direct
title: Monorepo
summary: Enter when the task is about monorepo governance and OSS contribution policy; avoid when task is extension-specific.
aliases:
  - monorepo governance
  - oss governance
triggers:
  - monorepo
  - oss governance
  - contribution policy
anti_triggers:
  - browser extension
  - wxt
required_with: []
may_include:
  - L1.monorepo.oss-governance
children:
  - L1.monorepo.oss-governance
entry_conditions:
  - monorepo_scope_confirmed
stop_conditions:
  - subdomain_selected
---

## 领域定义
- Monorepo 是面向多仓治理统一的结构、贡献与合规约束体系。

## 何时进入
- 任务明确涉及 monorepo 治理、OSS 规范、贡献流程或仓库策略。

## 何时不要进入
- 任务属于浏览器扩展实现路径或 WXT 配置本身。

## 推荐下一跳
- L1.monorepo.oss-governance
