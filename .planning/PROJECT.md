# the-agent-packs

## What This Is

这是一个 **agent-packs 增强包生产者系统**：把 Blueprint/AIDP 等知识资产转成可路由、可验证、可交付的上下文包，并通过 MCP/CLI 提供激活与消费能力。它服务的核心对象是开发与维护增强包生产链路的团队，以及需要按目标域获取最小完整上下文的消费侧 agents。

## Core Value

在不泄露全仓语义的前提下，让消费侧 agent 始终拿到“**目标域相关、最小但不遗漏、可验证**”的上下文交付结果。

## Requirements

### Validated

- ✓ Blueprint 节点可编译为 SQLite 索引并用于路由查询 — existing
- ✓ 路由与上下文包构建（route_query/build_context_bundle）可通过 CLI/MCP 执行 — existing
- ✓ Activation 可汇聚路由结果并输出 validation/result/handoff 结构 — existing
- ✓ WXT 主域样板（wxt-manifest）已形成最薄闭环验证路径 — existing

### Active

- [ ] 解析层稳定化：替换手工 YAML/frontmatter 解析，减少静默错误
- [ ] 索引构建事务化：避免编译中断导致半成品索引
- [ ] 路由治理硬化：坚持 global→domain→capability 与 attach-only 规则
- [ ] 消费契约验证化：可自动验证“最小且完整”的上下文交付
- [ ] 多域准入机制化：在护栏内支持第二主域样板接入

### Out of Scope

- 全量重写当前 Go 运行时 — 偏离当前里程碑目标，成本高风险大
- 让消费侧默认读取完整 docs/AIDP — 违背渐进披露与最小上下文契约
- capability 与主域同层全局竞争 — 会破坏路由边界和可解释性

## Context

- 当前仓库是 brownfield，已有可运行最薄闭环与测试资产（`.planning/codebase/*` 已完成扫描）
- `docs/AIDP/` 已是新的项目语义入口，`docs/改造计划v1/` 作为历史来源
- 当前首要任务不是“扩功能”，而是“稳定生产系统语义 + 交付契约 + 运行态维护”
- `.planning/research/` 已完成 stack/features/architecture/pitfalls 研究，建议先基础硬化再扩域

## Constraints

- **Architecture**: 必须保持 activation-first 与 bounded-context — 这是当前系统核心护栏
- **Routing**: capability 默认 attach-only，不可抢主域入口 — 防止跨域串线
- **Truth Source**: registry 是包身份真相源，AIDP 是语义真相源 — 禁止各自漂移
- **Runtime**: 关键变化需回写 runtime 工件 — 保证多-agent 接力可追溯
- **Technology**: 维持 Go + SQLite 主栈 — 当前阶段不做换栈迁移

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| docs/AIDP 作为默认语义入口 | 统一项目语义真相源，避免旧文档分叉 | ✓ Good |
| 保持 Go + SQLite + MCP 组合 | 与现有实现与交付契约匹配，最小扰动 | — Pending |
| 先修解析/索引稳定性再扩第二主域 | 扩域前先消除系统性脆弱点，降低返工 | — Pending |

---
*Last updated: 2026-03-16 after initialization*
