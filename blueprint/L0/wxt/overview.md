---
id: L0.wxt
level: L0
domain: wxt
subdomain: null
capability: null
title: WXT
summary: Enter when the task is about browser extension workflows using WXT; avoid when the host is not a web extension project.
aliases:
  - web extension
  - browser extension
triggers:
  - wxt
  - web extension
  - browser extension
anti_triggers:
  - tauri
  - telegram miniapp
required_with: []
may_include:
  - L1.wxt.manifest
children:
  - L1.wxt.manifest
entry_conditions:
  - extension_task_confirmed
stop_conditions:
  - subdomain_selected
---

## 领域定义
- WXT 是面向浏览器扩展的开发与配置体系。

## 何时进入
- 任务明确涉及浏览器扩展或 WXT 配置。

## 何时不要进入
- 任务与浏览器扩展无关或被 anti_triggers 排除。

## 必带横线
- 安全与权限相关子域（security）。
- 应用发布与商店审核相关子域（release）。

## 推荐下一跳
- L1.wxt.manifest
