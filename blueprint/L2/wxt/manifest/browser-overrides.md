---
id: L2.wxt.manifest.browser-overrides
level: L2
domain: wxt
subdomain: manifest
capability: null
node_kind: execution-strategy
visibility_scope: domain-scoped
activation_mode: direct
title: Manifest Browser Overrides
summary: Adjust manifest overrides per browser differences and store requirements.
aliases:
  - browser overrides
triggers:
  - overrides
  - browser differences
anti_triggers: []
required_with: []
may_include: []
children: []
entry_conditions:
  - execution_needed
stop_conditions:
  - overrides_ready
---

## 执行步骤
1. 确认目标浏览器列表与版本。
2. 识别 manifest 中需要覆盖的字段。
3. 输出差异化配置建议。

## 配置点
- 浏览器兼容目标与商店限制。

## 工具调用点
- 无强制工具，优先人工核对。

## validator 前置
- 需要 L1 层的商店审核结论。
