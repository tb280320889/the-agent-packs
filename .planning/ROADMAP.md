# Roadmap: the-agent-packs

**Created:** 2026-03-16  
**Input Sources:** `.planning/PROJECT.md`, `.planning/REQUIREMENTS.md`, `.planning/research/*`, `docs/AIDP/*`

## Overview

**5 phases** | **16 v1 requirements mapped** | **All v1 requirements covered ✓**

| # | Phase | Goal | Requirements | Success Criteria |
|---|-------|------|--------------|------------------|
| 1 | Foundation Hardening | Complete    | 2026-03-17 | 4 |
| 2 | Routing Governance | 固化路由边界与可解释性，移除隐式回退 | ROUT-01, ROUT-02, ROUT-03, ROUT-04 | 4 |
| 3 | Contracted Delivery | 让“最小且完整”上下文交付可验证 | CONT-01, CONT-02, CONT-03 | 3 |
| 4 | Validation & Runtime Governance | 校验闭环与运行态回写制度化 | VALD-01, VALD-02, GOVR-01 | 3 |
| 5 | Domain Expansion Pilot | 在护栏内接入第二主域并验证不破坏 WXT | DOMN-01, DOMN-02 | 3 |

## Phase Details

### Phase 1: Foundation Hardening
**Goal:** 替换易碎解析路径并实现事务化索引重建，确保基础能力可重复、可恢复。

**Requirements:** PARS-01, PARS-02, INDX-01, INDX-02

**Plans:** 4/4 plans complete

Plans:
- [ ] 01-01-PLAN.md — package.yaml 与 frontmatter 严格 YAML 解析硬化
- [ ] 01-02-PLAN.md — 事务化索引重建与结构化编译错误输出
- [ ] 01-03-PLAN.md — 解析回归测试与固定 fixture 基线
- [ ] 01-04-PLAN.md — 关闭验证缺口：报告失败场景下旧索引保留与回归验证

**Success Criteria:**
1. `package.yaml` 解析对未知/非法字段可稳定报错（非静默忽略）。
2. Blueprint frontmatter 解析覆盖列表、引号、多行等常见语法。
3. 索引重建失败后，系统仍保留上一个可用索引或明确回滚状态。
4. 编译与报告写入失败均有结构化错误输出并可被测试覆盖。

---

### Phase 01.1: 当前项目是brownfield , 核心业务文档逻辑有过改动,我需要进行一次当前代码的大检查和大重构, 一些目录都需要根据语义 领域 等等 设计 和 架构好,而不是散乱在各个地方和根目录,最后再进行一次 git commit 提交 (INSERTED)

**Goal:** 完成 brownfield 仓库的语义化结构重构与分批提交治理：先冻结目标目录与迁移映射，再执行迁移与路径修复，最后在“先归档后删除（需确认）”约束下完成可审阅的多批次提交与运行态回写。
**Depends on:** Phase 1
**Plans:** 4/4 plans complete

Plans:
- [ ] 01.1-01-PLAN.md — 冻结目标目录模型、迁移映射与归档清单
- [ ] 01.1-02-PLAN.md — 执行结构迁移与路径修复，产出分批提交切片
- [ ] 01.1-03-PLAN.md — 执行分批提交、删除决策门禁与运行态回写

### Phase 2: Routing Governance
**Goal:** 把主域优先与 capability attach-only 规则固化为实现与测试可验证行为。

**Requirements:** ROUT-01, ROUT-02, ROUT-03, ROUT-04

**Plans:** 2 plans

Plans:
- [ ] 02-01-PLAN.md — 重构 route 治理内核：candidate-space-first 与 attach-only 两阶段决策
- [ ] 02-02-PLAN.md — 固化可解释输出与 canonical 缺失显式失败语义

**Success Criteria:**
1. 所有路由路径先执行候选空间过滤再评分。
2. capability 不会在主域未确认前成为 primary 结果。
3. route/activate 输出包含“为何选中/为何附挂”的解释字段。
4. 无 registry canonical 映射时返回 explicit error/partial，不再回退硬编码包。

---

### Phase 3: Contracted Delivery
**Goal:** 把“最小且完整”的消费契约从文档主张落到可执行校验。

**Requirements:** CONT-01, CONT-02, CONT-03

**Success Criteria:**
1. context bundle 仅包含目标域任务必需节点，不混入无关域内容。
2. 交付对象中包含 include/exclude rationale，且可追溯到规则依据。
3. 建立可重复执行的契约检查用例（至少覆盖 WXT 样板与一个负例）。

---

### Phase 4: Validation & Runtime Governance
**Goal:** 让 activation→validation→runtime 回写形成稳定制度，而非会话约定。

**Requirements:** VALD-01, VALD-02, GOVR-01

**Success Criteria:**
1. validator 计划可从 registry 一致生成并执行（core + domain）。
2. validation result 与 artifacts/handoff 具备可追踪关联。
3. 关键变更能同步回写 runtime 账本（assumption/decision/change/validation）。

---

### Phase 5: Domain Expansion Pilot
**Goal:** 以受控试点方式引入第二主域，验证体系可扩展且不破坏既有样板。

**Requirements:** DOMN-01, DOMN-02

**Success Criteria:**
1. 第二主域可按同一 registry/routing 契约完成准入与激活。
2. WXT 样板路径在新增主域后仍保持行为一致与验证通过。
3. 形成可复用的“新主域准入清单”（命名、路由、验证、交付契约）。

## Dependency & Ordering Rationale

- **Phase 1 → Phase 2:** 路由治理依赖可靠解析与稳定索引。
- **Phase 2 → Phase 3:** 契约化交付依赖清晰路由边界与解释能力。
- **Phase 3 → Phase 4:** 校验制度化需要先明确交付契约目标。
- **Phase 4 → Phase 5:** 扩域前先确保验证与运行态治理闭环可持续。

## Risk Watch

- 若 Phase 1 未完成，后续所有 phase 结论都可能建立在不稳定数据面上。
- 若 Phase 2 放宽 attach-only，将直接破坏消费契约可信度。
- 若 Phase 4 回写未制度化，项目会再次退化到“会话真相源”。

---
*Roadmap created: 2026-03-16*
