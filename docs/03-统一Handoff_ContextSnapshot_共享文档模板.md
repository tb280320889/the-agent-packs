# 03 统一 Handoff / Context Snapshot 共享文档模板

## 作用域
本文只回答：
1. Context Snapshot 应该记录什么
2. Handoff 应该记录什么
3. 哪些内容必须写入共享文档，不能只留在对话中

## 一、为什么需要统一模板

在多-agent 接力开发中，最常见的失败原因不是代码写不出来，而是：
- 当前状态只存在于某个 agent 的上下文窗口里
- 下一个 agent 不知道哪些对象已经冻结
- issue 有状态，但没有说明为什么变成这个状态
- 实现做完了，但没有留下“下一步该怎么接”

所以本项目要求所有阶段性交接统一使用：
- `Context Snapshot`
- `Handoff`

## 二、Context Snapshot 的作用

`Context Snapshot` 用来回答：
- 现在项目的事实状态是什么
- 当前阶段已固定什么
- 当前输入和输出是什么
- 风险、阻塞、缺口是什么

它更像“阶段快照”，不是任务分配单。

## 三、Handoff 的作用

`Handoff` 用来回答：
- 下一个 agent 该做什么
- 哪些东西可以直接依赖
- 哪些东西绝对不能动
- 如果继续做，应该先读什么、跑什么、看什么 bead

它是“交接说明”，不是阶段总结散文。

## 四、Context Snapshot 模板

建议文件名：
- `docs/context-snapshots/<milestone>-<topic>.md`
- 或按日期：`docs/context-snapshots/2026-03-15-m2-wxt-manifest.md`

模板如下：

```md
# Context Snapshot: <标题>

## 1. 当前阶段
- 所属里程碑：M0 / M1 / M2 / M3 / M4
- 关联 bead：<id>
- 当前状态：open / in_progress / blocked / completed

## 2. 当前事实
- 当前要解决的问题：
- 当前已完成内容：
- 当前尚未完成内容：

## 3. 已冻结对象
- <对象名>：<冻结原因 / 说明>
- <对象名>：<冻结原因 / 说明>

## 4. 当前输入
- 上游交付物：
- 依赖文档：
- 依赖 schema / 模板 / fixtures：

## 5. 当前输出
- 已产出文件：
- 已更新文件：
- 已创建 bead：

## 6. 风险与阻塞
- 风险：
- 阻塞：
- 是否需要人工决策：是 / 否

## 7. 下一步建议
- 建议下一个 bead：
- 建议先执行的命令：
- 建议先阅读的文档：
```

## 五、Handoff 模板

建议文件名：
- `docs/handoffs/<bead-id>-handoff.md`
- 或 `docs/handoffs/<milestone>-to-<milestone>.md`

模板如下：

```md
# Handoff: <标题>

## 1. 交接对象
- 来源 bead：<id>
- 下一 bead：<id>
- 来源里程碑：M0 / M1 / M2 / M3 / M4
- 目标角色：Planner / Executor / Verifier / Handoff

## 2. 已完成什么
- 

## 3. 下一位 agent 可直接依赖什么
- 

## 4. 下一位 agent 必须先做什么
- 先 claim：<id>
- 先阅读：<文档列表>
- 先验证：<命令或检查项>

## 5. 不要做什么
- 

## 6. 风险与未决项
- 

## 7. 推荐下一动作
- 
```

## 六、最小必填要求

### Context Snapshot 最少必须有
- 当前阶段
- 当前事实
- 已冻结对象
- 风险与阻塞

### Handoff 最少必须有
- 已完成什么
- 下一位 agent 可依赖什么
- 不要做什么
- 风险与未决项

## 七、什么时候必须更新 Snapshot / Handoff

以下情况必须更新至少其中之一：

### 必须更新 Context Snapshot
- 完成一个里程碑阶段性目标后
- 发现前置假设失效后
- 结构性对象被冻结后
- issue 状态从 open/in_progress 变成 blocked/completed 后

### 必须写 Handoff
- 准备把任务交给下一个 agent
- 当前 bead 完成，但后续工作要继续
- 当前 bead 到达边界，需要 handoff 到别的 pack / 别的里程碑

## 八、与 bd 的协同规则

统一要求：
- `bd` 记录任务状态
- `Context Snapshot` 记录事实状态
- `Handoff` 记录执行交接

三者之间要互相可追溯：
- Snapshot 内写 bead id
- Handoff 内写来源 bead 与下一 bead
- bead 描述里可引用对应文档路径

## 九、推荐的文档目录

```text
docs/
├─ context-snapshots/
│  └─ <snapshot files>
└─ handoffs/
   └─ <handoff files>
```

如果项目规模还小，可以先不创建大量实例文件，但模板必须先冻结。

## 十、质量判断标准

一个合格的 Snapshot / Handoff 必须满足：
- 不依赖写作者的记忆才能理解
- 不需要回看聊天记录才能继续工作
- 能明确区分“已冻结”“未决”“禁止修改”
- 能让下一个 agent 直接开始，而不是重新探索全部上下文

## 非目标
本文不定义：
- 各里程碑业务细节
- pack 内部模板正文
- validator 规则细节

本文只定义共享上下文文档应该如何写、何时写、写到什么程度。
