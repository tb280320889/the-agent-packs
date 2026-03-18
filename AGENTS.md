# Agent Instructions

## 语言硬规则

**记住：所有 agent 的回答、说明、推理、备注、handoff、issue 描述补充，默认都必须使用中文。**。

## 项目定位

本仓库当前是 **以以项目专属 AIDP 为默认入口** 的开发与协作仓库。

- 当前主目标不是堆功能，而是在现有最薄闭环基础上，把仓库推进为支持多主域、多横线能力、多 workflow package 接入的稳定开发系统。
- 当前第一条固定领域线是 `WXT`，首个完整样板 workflow package 是 `wxt-manifest`。
- 所有开发都必须遵循 `activation-first`、`bounded-context`、`milestone-decoupled`、`progressive-disclosure` 原则，禁止越界扩张。
- 新的项目语义入口、agent 阅读入口、运行态维护入口统一位于 `docs/AIDP/`。

## 启动顺序

任何 agent 进入仓库后，默认按下面顺序启动：

1. 先阅读本文件 `AGENTS.md`
2. 再阅读 `docs/AIDP/README.md`
3. 再阅读 `docs/AIDP/AGENTS.md`
4. 再阅读 `docs/AIDP/AIDP-MANIFEST.yaml`
5. 再阅读 `docs/AIDP/VERSION.md`
6. 再阅读 `docs/AIDP/COMPATIBILITY.md`
7. 再阅读 `docs/AIDP/protocol/01-输入输出定义.md`
8. 再阅读 `docs/AIDP/protocol/04-默认假设协议.md`
9. 再阅读 `docs/AIDP/protocol/05-输出生成协议.md`
10. 再阅读 `docs/AIDP/protocol/10-增强开发与迭代协议.md`
11. 再阅读 `docs/AIDP/core/01-项目总览.md`
12. 再阅读 `docs/AIDP/core/05-范围与边界.md`
13. 再阅读 `docs/AIDP/core/06-业务规则与关键对象.md`
14. 再阅读 `docs/AIDP/core/08-技术约束与工程约定.md`
15. 再阅读 `docs/AIDP/core/10-验收标准.md`
16. 再阅读 `docs/AIDP/core/12-默认假设与待确认问题.md`
17. 再读取 GSD 进度追踪记录，以及 `docs/AIDP/runtime/` 下最新工件
18. 确认本次要执行的 GSD 任务项
19. 如需追溯历史来源，再按需补读 `docs/改造计划v1/` 对应文档

## 文档阅读边界

- 默认先读 `docs/AIDP/`，不要把 `docs/改造计划v1/` 当新的主入口。
- 如果发现 AIDP 中缺少关键来源事实，不要脑补，必须：记录 issue、标记阻塞，并回写到 AIDP runtime 工件。
- 若当前任务属于新里程碑、结构性增强或版本升级，必须先检查 `docs/AIDP/VERSION.md`、`docs/AIDP/COMPATIBILITY.md` 与 `docs/AIDP/protocol/10-增强开发与迭代协议.md`。

统一规则：

- 任务状态放在 GSD 的文档上下文(.planning目录)系统记录
- 系统与项目语义约束放在 `docs/AIDP/`
- 不允许把只存在于对话中的关键上下文当成“已共享上下文”
