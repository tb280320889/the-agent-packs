---
id: L1.demo.multiline
level: L1
domain: demo
subdomain: multiline
capability: null
node_kind: workflow-entry
visibility_scope: domain-scoped
activation_mode: direct
title: "Demo: Multi-line"
summary: |
  这是多行摘要。
  包含 "引号" 与 : 冒号。
aliases:
  - "demo alias"
  - 'multi line'
triggers:
  - manifest
  - "permissions"
anti_triggers: []
required_with: []
may_include: []
children: []
entry_conditions:
  - entry_ok
stop_conditions:
  - stop_ok
---

body
