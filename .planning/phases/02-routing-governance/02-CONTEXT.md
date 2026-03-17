# Phase 2: Routing Governance - Context

**Gathered:** 2026-03-17
**Status:** Ready for planning

<domain>
## Phase Boundary

本阶段仅固化路由治理行为：主域优先、capability attach-only、可解释路由结果、无 canonical 映射时的显式失败/partial 语义；不新增独立产品能力，不扩展到其他 phase 功能。

</domain>

<decisions>
## Implementation Decisions

### 路由解释信息粒度（面向 MCP 调用方开发 Agent）
- 默认反馈采用“极简必要信息”策略，优先提供可被 agent 稳定消费的最小字段。
- 采用双层模型：默认返回极简摘要；需要时可扩展返回结构化明细。
- 返回风格为“可操作”：失败/partial 场景提供简短原因 + 下一步建议。
- 响应中预留 `docs_ref` 扩展位（当前可为空），为后续官方文档链接能力留接口。

### 无 canonical 映射时的失败语义
- 默认严格模式：无 canonical 映射时不允许隐式回退，直接显式失败（hard fail）。
- 使用稳定机器错误码（示例：`ROUTE_CANONICAL_MISSING`）+ 简短 message。
- 默认不回传完整候选列表，避免暴露过多内部治理细节；仅返回必要上下文。
- 返回固定短建议（检查 registry canonical 映射）并支持 `docs_ref` 指向后续文档。

### capability 附挂可见性与可审计性
- 必须返回“为何附挂/为何未附挂”的极简解释（原因码 + 短说明）。
- 返回规则标识（如 `BR-02`、`BR-03`）作为依据，并可选 `docs_ref`。
- 默认输出“最终附挂列表 + 每项极简原因”，不默认输出冗长全链路。
- capability 被拒绝附挂时，返回拒绝原因码 + 1 条下一步建议。

### 候选冲突与稳定决策策略
- 同分冲突采用稳定 tie-break，保证同输入可复现（禁止随机选择）。
- tie-break 默认优先级：canonical 命中优先 > 明确主域匹配 > 规则优先级 > 名称字典序。
- 若无法稳定判定，返回 explicit error/partial，不做猜测性回退。
- 返回简短可复现标记（如 `decision_trace_id` 或 `decision_basis`）。

### Claude's Discretion
- 极简响应与详细响应的具体字段命名（在不违背“默认极简”前提下）。
- `decision_trace_id` 与 `decision_basis` 的具体数据结构和编码形式。
- 规则标识与错误码在响应中的嵌套层级设计。

</decisions>

<specifics>
## Specific Ideas

- 反馈对象是“调用方开发 Agent 使用 MCP 时的响应”，不是终端用户 UI。
- 核心诉求是默认极简必要信息，避免噪声，优先稳定可解析。
- 解释信息未来应可衔接官方文档链接，以支持快速排障与自助理解。

</specifics>

<deferred>
## Deferred Ideas

- 官方文档链接体系的完整建设（文档站点结构、链接路由策略、版本化说明）作为后续增强事项记录；当前阶段仅预留 `docs_ref` 扩展位。

</deferred>

---

*Phase: 02-routing-governance*
*Context gathered: 2026-03-17*
