# Feature Landscape

**Domain:** Agent Pack 增强包生产系统（Producer-side）  
**Researched:** 2026-03-16

## Table Stakes

Features users expect. Missing = product feels incomplete.

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Blueprint 编译与索引构建 | 生产系统必须把知识源转为可查询结构 | Med | 当前已具备，但需补强事务与迁移能力。 |
| 受约束路由（global → domain → capability） | 避免能力串线、越权激活 | Med | 是本项目 identity 的核心，不可降级为“自由检索”。 |
| Context Bundle 最小且完整交付 | 消费侧需要“够用且不冗余”的上下文 | High | 直接决定消费体验与正确率。 |
| Registry 作为包真相源 | package 身份、依赖、validator 映射必须可追踪 | Med | 若退化为文件名推断，会导致系统级不一致。 |
| Activation + Validation 闭环 | 必须有“可执行 + 可验收”结果 | Med | 无验证闭环会让增强包只剩静态文档价值。 |
| CLI/MCP 双入口 | 本地调试与外部工具消费都需要 | Low | 当前已有主入口，重点是契约稳定性。 |

## Differentiators

Features that set product apart. Not expected, but valued.

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| 渐进披露上下文契约（producer 明确 not leak all） | 显著降低 token 噪声与跨域误导 | High | 相比“全仓灌入”是实质差异化能力。 |
| 三角色语义隔离（生产者/消费侧/开发侧） | 减少系统演进中的角色混淆与错误决策 | Med | 文档协议化后可持续复用到多里程碑。 |
| capability attach-only 治理 | 横线能力不与主域抢入口，减少错误匹配 | Med | 对多域扩展稳定性非常关键。 |
| 运行态工件驱动协作（assumption/decision/validation ledger） | 多 agent 接力时可追溯、可审计 | Med | 降低“只在对话里成立”的隐性上下文风险。 |
| WXT 样板消费契约可验证示例 | 从抽象规则落到可测试样板 | Low | 便于后续复制到第二主域。 |

## Anti-Features

Features to explicitly NOT build.

| Anti-Feature | Why Avoid | What to Do Instead |
|--------------|-----------|-------------------|
| 让消费侧默认读取完整 docs/AIDP | 破坏最小上下文契约，增加跨域串线 | 只交付目标域 + 目标任务上下文切片 |
| capability 作为第一轮全局竞争入口 | 会削弱主域边界，造成误激活 | 坚持主域确认后 capability attach |
| 在当前阶段重写 Go 运行时/全面换栈 | 风险与成本高，偏离里程碑目标 | 先做协议一致性与脆弱点修复 |
| 用“更多上下文”掩盖语义不清 | 只会提高噪声，不会提高正确性 | 提升路由可解释性与交付契约清晰度 |

## Feature Dependencies

```text
Blueprint 编译/索引
  → Registry 真相源对齐
  → 路由候选空间裁剪
  → Context Bundle 组装
  → Activation 执行
  → Validation 闭环
  → Handoff / 消费侧交付

角色边界清晰
  → 渐进披露契约可执行
  → 多域扩展可控
```

## MVP Recommendation

Prioritize:
1. **索引与解析稳定化**（YAML/frontmatter 标准解析 + 事务化索引重建）
2. **路由与包身份治理硬化**（消除硬编码回退、强化 registry 约束）
3. **消费契约可验证化**（最小且完整交付的自动化验证基线）

Defer: **第二主域大规模接入**：在基础协议与校验稳定前直接扩域，会放大技术债并污染契约边界。

## Sources

- docs/AIDP/core/05-范围与边界.md（HIGH，项目边界）  
- docs/AIDP/core/06-业务规则与关键对象.md（HIGH，关键规则与对象）  
- docs/AIDP/core/10-验收标准.md（HIGH，验收定义）  
- .planning/codebase/ARCHITECTURE.md（MEDIUM，当前实现闭环）  
- .planning/codebase/CONCERNS.md（MEDIUM，当前脆弱点）  
- MCP Spec 2025-06-18（MEDIUM，协议能力与安全原则）: https://modelcontextprotocol.io/specification/2025-06-18
