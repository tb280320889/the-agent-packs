# 04 多Agent接力开发与bd协作规则

## 作用域
本文只回答：
1. 多个 agent 如何在同一项目中接力开发
2. `bd` 在任务拆分、依赖、阻塞、并行中的使用规则
3. 什么信息放在 issue，什么信息放在文档，什么信息必须 handoff

## 一、核心原则

### 1. issue 是任务真相源
所有任务状态、依赖关系、阻塞原因，统一以 `bd` 为准。

禁止：
- 用 markdown TODO 代替 issue
- 只在对话里说明“后面要做什么”
- 只在本地笔记里留待办

### 2. 文档是约束真相源
设计边界、固定对象、冻结规则、模板，统一放在 `docs/`。

### 3. handoff 是执行真相源
当前 agent 做到哪里、下一位 agent 可依赖什么、不要做什么，必须通过 handoff 明确写出。

## 二、推荐的多-agent 角色分工

### Planner Agent
负责：
- 拆分 epic
- 建立 parent-child / blocks 依赖
- 指定当前工作属于哪个里程碑
- 明确冻结对象与边界

### Executor Agent
负责：
- 领取 bead
- 按最小阅读集执行改动
- 补齐测试、文档、实现
- 发现新工作时创建 bead 并挂依赖

### Verifier Agent
负责：
- 检查交付物是否满足完成定义
- 运行测试、lint、结构校验
- 判断是否可关闭 bead 或必须退回

### Handoff Agent
负责：
- 产出 Context Snapshot
- 写 Handoff
- 标记风险、阻塞、冻结对象
- 明确下一位 agent 的输入边界

同一个 agent 在小项目中可以兼任多个角色，但输出物不能省略。

## 三、bd 的基础工作流

### 启动必跑
```bash
bd prime
bd ready --json
bd status --json
```

### 开工前
```bash
bd show <id> --json
bd update <id> --claim --json
```

### 发现新工作
```bash
bd create "标题" --description "上下文" -t task -p 1 --json
bd dep add <new-id> <current-id> --type discovered-from --json
```

### 存在前后置关系
```bash
bd dep add <blocked-id> <blocking-id> --type blocks --json
```

### 收尾
```bash
bd close <id> --reason "Completed" --json
bd dolt push
```

## 四、推荐 issue 设计法

### epic
适用于：
- 跨多个里程碑
- 需要多个 agent 接力
- 有多个可独立关闭的子任务

### task
适用于：
- 一个 agent 可独立完成
- 有明确交付物
- 范围边界清晰

### bug
适用于：
- 协议冲突
- 文档缺口
- 实现错误
- 闭环阻塞

### chore
适用于：
- 初始化
- 治理
- 同步
- 工具配置

## 五、依赖关系怎么用

### parent-child
用于：
- epic 与子任务
- 主任务与子交付

### blocks
用于：
- 明确前置完成后，后置才能开工
- 例如 M0 未冻结，M1 不能正式开工

### related
用于：
- 相关但不构成阻塞

### discovered-from
用于：
- 执行中发现的新问题、新任务、新风险
- 必须挂回来源 bead，避免孤儿任务

## 六、多-agent 接力的最小流程

### 流程 A：标准接力
1. Planner 创建 epic 和子任务
2. Executor claim 一个 task
3. Executor 完成后补 Context Snapshot
4. Verifier 验证完成定义
5. Handoff Agent 写 handoff，决定 close 或转下一任务

### 流程 B：执行中发现新工作
1. 当前 Executor 先继续收敛本任务主目标
2. 新发现工作立即创建 bead
3. 使用 `discovered-from` 挂回来源 bead
4. 若构成阻塞，再补 `blocks`

### 流程 C：被外部条件阻塞
适用场景：
- 等待人工确认
- 等待上游 issue 完成
- 等待跨仓依赖

可使用：
```bash
bd gate list
bd gate resolve <id>
```

规则：
- 阻塞必须进入 issue 或 gate
- 不允许只在对话里说“先等等”

## 七、并行开发规则

### 工作分叉
当两个 agent 需要并行开发时：
- 优先拆成不同 bead
- 避免共同改同一文件集合
- 如必须并行且存在冲突风险，使用独立 worktree

```bash
bd worktree create <name>
bd worktree list
```

### 合并冲突串行化
当多个分支需要争抢同一合并窗口时，使用 merge slot：

```bash
bd merge-slot acquire
bd merge-slot release
```

适用场景：
- 同时修改核心协议
- 同时修改 AGENTS.md
- 同时修改共享模板

## 八、agent 状态与审计

### agent state
长时任务可记录 agent 运行状态：

```bash
bd agent state <agent-id> working
bd agent heartbeat <agent-id>
```

适用场景：
- 多 agent 编排
- 长跑验证
- 需要监控某 agent 是否卡住

### audit record
重要决策建议写入审计：

```bash
bd audit record
```

适合记录：
- 为什么拆出新 bead
- 为什么修改边界
- 为什么需要人工介入
- 为什么从 completed 降回 partial

## 九、信息落位规则

### 放在 bd issue 的内容
- 任务标题
- 当前状态
- 阻塞关系
- 所属里程碑
- 新发现工作

### 放在 Context Snapshot 的内容
- 当前阶段事实
- 已冻结对象
- 当前输入输出
- 风险、阻塞、缺口

### 放在 Handoff 的内容
- 下一位 agent 要做什么
- 下一位 agent 可依赖什么
- 下一位 agent 不要做什么
- 推荐的下一个 bead / 命令 / 阅读集

## 十、一个科学的多-agent 项目追踪模型

推荐采用三层结构：

1. `bd epic/task/bug/chore`：负责任务追踪
2. `docs/`：负责规则、模板、协议、边界
3. `Context Snapshot + Handoff`：负责接力上下文共享

这三层缺一不可：
- 只有 issue，没有文档，会导致约束漂移
- 只有文档，没有 issue，会导致状态不可追踪
- 只有实现，没有 handoff，会导致交接失败

## 十一、完成定义
多-agent 接力机制只有在以下条件下才算成立：
- 新工作都能进入 bd
- 每个 bead 都有清晰负责人或可 claim 状态
- 关键边界都有文档承载
- 每次接力都有快照和 handoff
- 后续 agent 不需要依赖私聊上下文才能继续工作

## 非目标
本文不定义：
- 具体某个 pack 的实现细节
- 某个里程碑的业务内容
- 外部平台集成流程

本文只定义多-agent 如何协作、如何追踪、如何共享上下文。
