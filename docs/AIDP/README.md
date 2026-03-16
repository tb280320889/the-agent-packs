# the-agent-packs Local AIDP

这是 `the-agent-packs` 仓库的 **local project AIDP**（项目实例级 AIDP）文档集。

它不是 AIDP system 模板包的复制品，而是基于本仓库现状、改造计划 v1 文档、`.planning/codebase/` Brownfield 扫描结果，以及上游 system AIDP 规范实例化后的 **项目语义真实源**。

首先要明确：`the-agent-packs` 的顶层身份，不是泛化的 agent 编排系统，而是 `agent-packs` 增强包的生产者系统。

本项目的核心职责，是把经过治理、分域、分层、可验证的知识与 workflow 能力，生产为可被消费侧 agents 通过 MCP 等方式按需消费的增强 pack；让消费侧 agent 在渐进披露模式下，仅获得目标域相关、最精简但不遗漏的上下文知识。

## 它在整体系统中的位置

本地工作链路应理解为：

`system AIDP -> local project AIDP + AGENTS.md -> GSD 生成/维护 .planning -> 开发 agent 在 .planning + AGENTS.md + 必要 local AIDP 指引下工作`

这意味着：
- 上游 system AIDP 提供协议模板与治理规则
- `docs/AIDP/` 提供本项目作为“增强包生产者系统”的业务语义真实源
- GSD 继续负责自己的任务流与 `.planning/` 工件体系
- 后续迭代时，应先升级 local AIDP，再继续 GSD 工作流

## 文档目标

本包首先用于让开发本项目的内部 agent、维护者与协作者，稳定理解以下事实：

1. 这个仓库当前在生产什么
2. 消费侧 agent 将如何消费本项目产出的增强 pack
3. 开发 agent 与消费侧 agent 的身份边界是什么
4. 哪些边界已经冻结，哪些仍属于待确认或后续阶段
5. 继续开发时应该读取哪些文档、维护哪些运行态工件
6. 如何把项目语义映射到 GSD 规划与执行工件

它不是消费侧 agent 的默认全文入口。消费侧真正应拿到的，是经过治理、裁剪和验证后的目标域上下文切片，而不是整个 `docs/AIDP/`。

## 当前项目类型

- `brownfield`：已有代码、已有最薄闭环、已有测试与文档资产
- `milestone-extension`：在既有实现基础上推进 Agent Pack 改造计划 v1 及其后续执行

## 当前北极星

把现有以 `WXT -> wxt-manifest` 为样板的最薄闭环，升级为一个稳定的增强包生产系统：既能支持多主域、多横线能力、多 workflow package 持续接入，又能保证消费侧 agents 通过渐进披露模式只获得目标域相关、最小完整、可验证的上下文包。

## 三类身份视角

### 生产者角色
本项目作为增强包生产者系统，负责治理知识、组织 package、定义边界、维护验证，并生产可被消费的增强 pack。

### 消费侧 agent 角色
消费侧 agent 通过 MCP 等方式消费特定主域或 capability 对应的上下文切片，而不是默认理解整个仓库治理语义。

### 开发本项目的开发 agent 角色
开发 agent 负责建设和维护这个生产系统本身，必须理解项目全局语义、规则和运行态，但不能把自己默认代入最终消费侧角色。

## 推荐阅读顺序

### 给 agent
1. `AGENTS.md`
2. `AIDP-MANIFEST.yaml`
3. `core/00-身份视角与项目定位.md`
4. `protocol/01-输入输出定义.md`
5. `protocol/04-默认假设协议.md`
6. `protocol/05-输出生成协议.md`
7. `core/01-项目总览.md`
8. `core/05-范围与边界.md`
9. `core/06-业务规则与关键对象.md`
10. `core/08-技术约束与工程约定.md`
11. `core/10-验收标准.md`
12. `core/12-默认假设与待确认问题.md`
13. `runtime/01-默认假设账本.md`
14. `runtime/02-决策日志.md`
15. `adapters/gsd/00-协同原则.md`
16. `adapters/gsd/01-全局产物映射.md`

### 给人类维护者
1. 先看本文件
2. 再看 `core/01-项目总览.md`
3. 再看 `core/05-范围与边界.md`
4. 若要推进当前版本实施，再看 `core/10-验收标准.md` 与 `adapters/gsd/`

## 与旧文档的关系

- `docs/改造计划v1/` 仍保留为历史改造设计与阶段性结论来源
- 但新的 agent 入口、项目语义总入口、运行态维护入口，统一以 `docs/AIDP/` 为准
- 如果 `docs/改造计划v1/` 与 `docs/AIDP/` 出现表达差异，应以 `docs/AIDP/` 中最新的项目语义、假设账本、决策日志与变更摘要为准

## 与 system AIDP 的关系

- 上游模板版本见 `VERSION.md` 中的 `Base System AIDP Version`
- 上游负责模板协议与治理演进，本地负责项目实例真相
- 不应把上游模板升级直接当成本地项目语义升级；必须先做兼容性与思想对齐判断

## 版本治理与持续迭代

本地 AIDP 已引入版本与迭代机制，相关文件包括：
- `VERSION.md`
- `CHANGELOG.md`
- `COMPATIBILITY.md`
- `GOVERNANCE.md`
- `CONTRIBUTING.md`
- `RELEASE-CHECKLIST.md`
- `protocol/10-增强开发与迭代协议.md`
- `runtime/08-增强开发任务启动单.md`
- `runtime/09-迭代会话交接.md`
- `adapters/gsd/06-增强开发与变更落地.md`

## 当前最小闭环

当前仓库的最小闭环是：

`Blueprint Markdown -> SQLite index -> route query -> context bundle -> workflow package activation -> validator -> activation result / handoff`

这是内部生产闭环中的技术主链路，不等于消费侧默认看到的全部语义。消费侧默认应看到的是与其目标主域和能力相关的增强 pack / context bundle。

## 运行态维护要求

每次重要开发、规划、验证或修正后，至少要更新：

1. `runtime/02-决策日志.md`（如有关键决策）
2. `runtime/03-变更摘要.md`
3. `runtime/04-phase-context.md`
4. `runtime/06-验证记录.md`

## 一句话结论

本 local AIDP 文档集是本仓库新的默认项目入口，负责把旧改造计划、现有代码事实、上游 system AIDP 规则与 GSD 执行方式统一收敛成“增强包生产者系统”的业务语义真实源，并明确消费侧 agent 应如何按渐进披露模式获得目标域上下文。
