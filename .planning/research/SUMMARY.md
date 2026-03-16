# Research Summary: the-agent-packs

**Domain:** Agent Pack 增强包生产者系统（MCP 上下文交付）  
**Researched:** 2026-03-16  
**Overall confidence:** HIGH

## Executive Summary

当前生态判断非常明确：这个项目不是“再造一个通用编排器”，而是一个**生产侧语义治理 + 运行时交付契约**系统。它的价值不在功能堆叠，而在于把领域知识稳定生产成可被消费侧精准调用的增强包，并确保交付上下文“最小但完整”。

技术上，现有 Go + SQLite + MCP 组合是对位的，不建议换栈。Go 与现有实现/测试/发布链深度耦合；SQLite 对本地索引型场景适配度高；MCP 则提供了协议化能力协商与工具调用边界。真正风险不在选型，而在工程脆弱点：手工 YAML/frontmatter 解析、非事务索引重建、路由硬编码回退。

特性上，table stakes 已具备雏形（compile/route/bundle/activate/validate），但 differentiator 还未完全产品化：特别是“渐进披露契约的可验证性”与“三角色语义隔离的执行闭环”。这两项是后续 roadmap 的主轴，优先级高于第二主域扩展。

架构上应坚持 Candidate-space-first（先裁候选再评分）与 Registry-as-truth（注册表真相源）两条硬规则。只要这两条被破坏，多域扩展会快速退化为不可解释的上下文拼接系统，最终导致返工。

## Key Findings

**Stack:** 继续 Go + SQLite + MCP，先修解析与事务一致性，再扩功能。  
**Architecture:** Compile → Route → Bundle → Activate → Validate 五段式必须保持，并强化可解释性输出。  
**Critical pitfall:** 语义真相源（AIDP）与实现真相源（registry/index/runtime）漂移，会直接破坏多 agent 协作质量。

## Implications for Roadmap

Based on research, suggested phase structure:

1. **Foundation Hardening（基础硬化）** - 先消除脆弱点，防止后续扩域放大风险  
   - Addresses: 索引构建、YAML/frontmatter 解析、registry 一致性  
   - Avoids: 手工解析误配、非事务重建导致的不稳定

2. **Routing Governance（路由治理）** - 固化主域优先 + capability attach-only 规则  
   - Addresses: 候选空间裁剪、主包选择可解释性  
   - Avoids: capability 越级、硬编码回退

3. **Contracted Delivery（契约化交付）** - 让“最小且完整”可被自动验证  
   - Addresses: context bundle inclusion/exclusion rationale、validation 闭环  
   - Avoids: 全量灌入、关键约束遗漏

4. **Domain Expansion Pilot（扩域试点）** - 在护栏内接入第二主域样板  
   - Addresses: 多域准入、命名治理、跨域隔离验证  
   - Avoids: 过早扩域引发系统性返工

5. **Operationalization（运行态制度化）** - 固化回写、验收、交接机制  
   - Addresses: runtime 工件维护、phase 验证、交接质量  
   - Avoids: 结论只停留在会话上下文

**Phase ordering rationale:**
- 先修基础可靠性（解析/索引）是所有上层能力的前置依赖。  
- 路由治理在契约交付前完成，才能确保交付边界可解释。  
- 契约验证稳定后再扩域，避免“带病扩展”。  
- 最后把运行态制度化，确保可持续迭代。

**Research flags for phases:**
- Phase 1: 标准工程问题，研究深度需求低（已知路径明确）  
- Phase 2: 需中等深度研究（尤其可解释路由输出格式）  
- Phase 3: 需高深度研究（“最小且完整”自动化评测口径）  
- Phase 4: 需高深度研究（新域准入策略与跨域隔离验证）

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | HIGH | 官方文档 + 现有代码事实一致，方向清晰 |
| Features | HIGH | AIDP 边界定义明确，table stakes/differentiators 可落地 |
| Architecture | HIGH | 当前实现已有闭环，问题主要在硬化而非重构 |
| Pitfalls | HIGH | 本地代码审计已给出具体脆弱点，且与 AIDP 规则一致 |

## Gaps to Address

- `modernc.org/sqlite` 从 v1.38.2 升级到最新稳定线的兼容性验证仍需专项测试（含 `modernc.org/libc` 配套）。  
- “最小但不遗漏”需要建立可量化验证指标（当前更多是规则声明，自动化口径不足）。  
- 第二主域候选与准入标准尚未冻结，Phase 4 前需补充里程碑级研究。
