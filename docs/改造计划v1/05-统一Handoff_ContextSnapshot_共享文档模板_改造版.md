# 05 统一 Handoff / Context Snapshot 共享文档模板（改造版）

## 作用域
定义改造计划 v1 的共享上下文模板。

## 一、Context Snapshot 模板

建议文件名：
- `docs/改造计划v1/context-snapshots/<date>-<topic>.md`

模板：

```md
# Context Snapshot: <标题>

## 1. 当前阶段
- 所属里程碑：M0 / M1 / M2 / M3 / M4
- 关联 GSD 任务项：<id>
- 当前状态：open / in_progress / blocked / completed

## 2. 当前事实
- 当前要解决的问题：
- 当前已完成内容：
- 当前尚未完成内容：

## 3. 已冻结对象
- <对象名>：<说明>

## 4. 当前输入
- 上游交付物：
- 依赖文档：
- 依赖实现：

## 5. 当前输出
- 已产出文件：
- 已更新文件：
- 已创建 GSD 任务项：

## 6. 风险与阻塞
- 风险：
- 阻塞：
- 是否需要 breaking 评估：是 / 否

## 7. 下一步建议
- 建议下一个 GSD 任务项：
- 建议先阅读的文档：
- 建议先验证的命令：
```

## 二、Handoff 模板

建议文件名：
- `docs/改造计划v1/handoffs/<task-id>-handoff.md`

模板：

```md
# Handoff: <标题>

## 1. 交接对象
- 来源 GSD 任务项：<id>
- 下一 GSD 任务项：<id>
- 来源里程碑：M0 / M1 / M2 / M3 / M4
- 目标角色：项目内部维护 agent / 迭代开发子 agent / 文档型 agent / Verifier

## 2. 已完成什么
-

## 3. 下一位 agent 可直接依赖什么
-

## 4. 下一位 agent 必须先做什么
- 先在 GSD 任务记录中认领任务项：
- 先阅读：
- 先验证：

## 5. 不要做什么
-

## 6. 风险与未决项
-

## 7. 推荐下一动作
-
```
