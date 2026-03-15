# Blueprint 最小索引 schema（M1）

本文件定义 M1 最薄 compiler 的 SQLite schema 目标，用于校验与查询骨架。

## nodes

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| id | text | 节点唯一 id（与路径一致） |
| level | text | L0/L1/L2/L3 |
| domain | text | 一级域 |
| subdomain | text | 子域 |
| capability | text | 能力线或 null |
| node_kind | text | domain-root/domain-orchestrator/capability-root/workflow-entry/execution-strategy |
| visibility_scope | text | global/domain-scoped/capability-scoped |
| activation_mode | text | direct/attach-only |
| title | text | 节点标题 |
| summary | text | 最小摘要 |
| path | text | 源文件路径 |
| parent_id | text | 父节点 id |
| body_md | text | 正文（可选） |
| entry_conditions_json | text | JSON 数组 |
| stop_conditions_json | text | JSON 数组 |
| checksum | text | 内容校验 |
| updated_at | text | 更新时间 |

## node_meta

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| node_id | text | 节点 id |
| aliases | text | JSON 数组 |
| triggers | text | JSON 数组 |
| anti_triggers | text | JSON 数组 |
| tags | text | JSON 数组 |

## edges

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| source_id | text | 源节点 |
| target_id | text | 目标节点 |
| edge_type | text | child/required_with/may_include/excludes |

> 注意：schema 只是 M1 目标形状，实际实现可在 M2/M3 细化。
