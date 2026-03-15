# Agent Instructions

## 语言硬规则

**记住：所有 agent 的回答、说明、推理、备注、handoff、issue 描述补充，默认都必须使用中文。**

- 只有在以下情况才允许保留英文原文：协议字段名、命令、文件路径、schema key、代码、外部工具原始输出。
- 即使引用英文对象名，也必须用中文解释其作用。
- 如果新增文档、issue、注释或交接说明，默认先写中文，再按需要附英文标识。

## 项目定位

本仓库当前是 **Agent Pack 改造计划 v1 的文档与协作规范仓库**。

- 当前主目标不是堆功能，而是把改造计划 v1 的 `M0 -> M1 -> M2 -> M3 -> M4` 这条里程碑链做成可被多个 agent 接力执行的稳定开发系统。
- 当前第一条固定领域线是 `WXT`，首个完整样板 workflow package 是 `wxt-manifest`。
- 所有开发都必须遵循 `activation-first`、`bounded-context`、`milestone-decoupled`、`progressive-disclosure` 原则，禁止越界扩张。

## 启动顺序

任何 agent 进入仓库后，默认按下面顺序启动：

1. 先阅读本文件 `AGENTS.md`
2. 再阅读 `docs/改造计划v1/00-总索引与使用说明.md`
3. 再阅读 `docs/改造计划v1/01-通用改造方法论与增量迁移规则.md`
4. 再阅读 `docs/改造计划v1/02-目标系统模型_总编排_子编排与分层索引.md`
5. 再阅读 `docs/改造计划v1/03-角色体系与协作模型.md`
6. 再阅读 `docs/改造计划v1/04-多Agent接力开发与bd协作规则_改造版.md`
7. 再阅读 `docs/改造计划v1/05-统一Handoff_ContextSnapshot_共享文档模板_改造版.md`
8. 运行 `bd prime`
9. 运行 `bd ready --json`
10. 运行 `bd status --json`
11. 找到本次要执行的 bead，先 `claim`，再按所属里程碑补读最小阅读集

如果 issue 已明确属于某个里程碑，则继续阅读：

- `M0`：`docs/改造计划v1/10-M0_角色冻结与边界校正_开发指导.md` + `docs/改造计划v1/11-M0_上下文_角色模型_职责与非目标.md` + `docs/改造计划v1/12-M0_上下文_冻结面_兼容面与禁止事项.md`
- `M1`：`docs/改造计划v1/20-M1_分层索引与总编排骨架_开发指导.md` + `docs/改造计划v1/21-M1_上下文_节点分类_作用域与可见性.md` + `docs/改造计划v1/22-M1_上下文_Routing分层化与候选集裁剪.md`
- `M2`：`docs/改造计划v1/30-M2_package注册表与命名空间治理_开发指导.md` + `docs/改造计划v1/31-M2_上下文_命名规则_注册表字段与冲突裁决.md`
- `M3`：`docs/改造计划v1/40-M3_文档吸收管道与知识资产化_开发指导.md` + `docs/改造计划v1/41-M3_上下文_外部资料纳入_语义映射与资产落位.md`
- `M4`：`docs/改造计划v1/50-M4_迁移实施_兼容验证与准入演练_开发指导.md` + `docs/改造计划v1/51-M4_上下文_迁移步骤_回滚策略与验收清单.md`

## 文档阅读边界

- 每个里程碑 agent 默认只读自己的最小阅读集，加上 `docs/改造计划v1/00~05` 六份共享文档。
- 除非 issue 明确要求，不要回退到“全量总纲式阅读”。
- 当前里程碑只能消费前一里程碑的固定交付物、冻结对象、风险清单与禁止事项。
- `docs/改造计划v1/02-目标系统模型_总编排_子编排与分层索引.md` 是当前系统层母版，任何里程碑都不允许绕过。
- `docs/改造计划v1/03-角色体系与协作模型.md` 保存长期协作骨架，不代表当前就要全做。
- `docs/改造计划v1/04-多Agent接力开发与bd协作规则_改造版.md` 与 `docs/改造计划v1/05-统一Handoff_ContextSnapshot_共享文档模板_改造版.md` 是共享协作层，不是可选材料。
- 如果发现前置交付物缺失，不要脑补，必须：记录 issue、标记阻塞、补 handoff 或 context snapshot。

## 系统层判断

当前系统骨架固定为：

`Git truth layer -> Markdown + frontmatter -> SQLite index -> Blueprint Query MCP -> minimal context bundle -> workflow package -> validator -> activation result / handoff`

必须始终记住：
- Git 管源，SQLite 管查
- Blueprint 是图谱，不是自由文档库
- Blueprint Query MCP 是唯一受控查询入口层
- Agent 默认只消费最小上下文包，不自由搜全文文档
- L3 是升级层，不是默认阅读层

## Issue Tracking with bd

本项目使用 **bd (beads)** 作为唯一任务追踪系统。不要使用 markdown TODO、临时清单或外部任务系统替代。

### 启动必跑命令

```bash
bd prime
bd ready --json
bd status --json
```

### 日常核心命令

```bash
bd show <id> --json
bd update <id> --claim --json
bd create "标题" --description "上下文" -t task -p 1 --json
bd dep add <blocked> <blocking> --json
bd close <id> --reason "Completed" --json
bd dolt push
```

### 任务设计规则

- `epic`：一条需要多个 agent 接力的大目标
- `task`：一个 agent 可以独立完成的工作单元
- `bug`：发现实际缺口、冲突、错误、阻塞
- `chore`：初始化、同步、治理、维护类工作

### 依赖关系规则

- `parent-child`：用于 epic 与子任务
- `blocks`：用于严格前后置依赖
- `related`：用于相关但不阻塞
- `discovered-from`：执行过程中发现的新工作，必须挂回来源 issue

### 多-agent 接力规则

- 开工前先 `bd update <id> --claim --json`
- 一个 agent 同一时刻只应主负责一个 bead
- 发现新工作，不要口头留下，立即创建 bead 并建立依赖
- 如果被外部条件阻塞，用 `bd gate` 表达等待条件
- 如果并行分支需要避免合并打架，用 `bd merge-slot acquire` / `release`
- 如果需要多工作目录并行开发，用 `bd worktree create <name>`
- 关键决策、偏航原因、人工介入点可写入 `bd audit record`
- 长时运行或多 agent 编排时，可用 `bd agent state` 与 `bd agent heartbeat` 标记状态

## 共享上下文系统

本项目默认使用“**bd issue + Context Snapshot + Handoff**”三层共享上下文：

1. `bd issue`：记录任务归属、状态、依赖、阻塞
2. `Context Snapshot`：记录当前阶段事实、冻结对象、输入输出、未决项
3. `Handoff`：记录从当前 agent 到下一个 agent 的可执行交接信息

统一规则：

- 任务状态放在 `bd`
- 系统与里程碑约束放在 `docs/改造计划v1/`
- 阶段性交接使用 `docs/改造计划v1/05-统一Handoff_ContextSnapshot_共享文档模板_改造版.md` 中的模板
- 不允许把只存在于对话中的关键上下文当成“已共享上下文”

## 多-agent 工作流

推荐采用下面的接力模型：

1. `Planner Agent`：拆 epic、建依赖、定义里程碑 bead
2. `Executor Agent`：领取单个 bead，按最小阅读集实施
3. `Verifier Agent`：验证交付物、测试、schema、闭环状态
4. `Handoff Agent`：补齐 snapshot、更新风险、关闭或转交 bead

每次接力必须至少完成四件事：

1. 更新 bead 状态
2. 产出或更新 context snapshot
3. 明确冻结对象与未决项
4. 给下一个 agent 明确“可依赖什么 / 不要做什么”

## 交付要求

每完成一个阶段性工作，至少要留下：

- 对应 bead 的状态更新
- 一份简短中文 handoff
- 一份可追溯的 Context Snapshot 或其更新
- 明确的风险/阻塞说明
- 若有结构性修改，说明影响范围与是否属于 breaking change

## Non-Interactive Shell Commands

**始终使用非交互参数**，避免命令卡住等待确认。

```bash
cp -f source dest
mv -f source dest
rm -f file
rm -rf directory
cp -rf source dest
```

其他可能交互的命令：

- `scp` 使用 `-o BatchMode=yes`
- `ssh` 使用 `-o BatchMode=yes`
- `apt-get` 使用 `-y`
- `brew` 使用 `HOMEBREW_NO_AUTO_UPDATE=1`

## 会话结束前必须做的事

1. 为剩余工作创建或更新 bead
2. 运行必要验证：测试、lint、构建、文档校验
3. 更新 handoff / context snapshot
4. 关闭已完成 bead，或把 bead 状态改成真实状态
5. 执行同步：

```bash
git pull --rebase
bd dolt push
git push
git status
```

6. 确认没有“只存在本地、未交接、未入 issue”的上下文残留

## 关键禁令

- 禁止用英文输出替代中文说明
- 禁止绕过 `bd` 私下记录任务
- 禁止跳过 claim 直接并行改同一工作单元
- 禁止在无 handoff 的情况下把问题口头转交给下一个 agent
- 禁止把未冻结对象当成稳定协议传播
- 禁止因为单个 workflow package 方便就修改全局边界
- 禁止让 agent 直接自由搜索 Blueprint Markdown 作为主路由方式
- 禁止在 route 阶段默认深读 L2/L3
