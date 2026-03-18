# Phase 4: Validation & Runtime Governance - Context

**Gathered:** 2026-03-18
**Status:** Ready for planning

<domain>
## Phase Boundary

本阶段仅交付 activation -> validation -> runtime 的制度化闭环：让校验触发、结果判定、证据关联与 runtime 回写具备稳定规则与可追踪性，避免依赖会话约定。

不包含新增业务能力或新主域扩展；新增能力应进入后续独立 phase。

</domain>

<decisions>
## Implementation Decisions

### 验证触发与执行节奏
- 默认采用混合触发：计划内关键里程碑（如计划完成、关键规则变更、validator 清单变化）必须自动触发验证。
- 允许手动补跑，但必须保留历史记录，且只标记最新一次为当前有效结论，不覆盖历史。
- 失败处理采用按等级分流：`failed(error)` 阻断当前计划推进；`warned(warn)` 可继续但必须留痕处理意见；`passed` 正常放行。
- 阻断边界默认限定在当前计划（plan）级别，避免跨 phase 误阻断；跨计划阻断需人工显式升级。

### 验证结果判定与可读性
- 状态分级固定为三态：`passed` / `warned` / `failed`。
- 对下游的最小输出采用可审计契约：结论状态、错误/规则码、触发原因、修复建议、证据引用。
- 结果展示采用双视图：机器可读结构化视图（供 MCP/自动化消费）+ 人类摘要视图（供评审与交接）。
- `warned` 语义固定为“可继续但必须留痕”，不得静默通过。

### Artifacts/Handoff 关联粒度
- 证据关联主键采用 `run_id`（单次验证执行实例），并绑定 phase/plan 作为上下文维度。
- 单次验证最小证据集合采用完整链路：验证结论、命中规则、输入摘要、对应计划/summary/handoff/runtime 引用。
- 引用强度采用混合策略：关键节点（validation result、handoff、runtime ledger）使用强 ID 串联；说明性文本允许弱引用。
- 历史留存采用“全量保留 + 当前有效标记”，禁止只保留最新而丢失审计轨迹。

### Runtime 账本回写规则
- 回写触发采用混合模式：关键事件即时回写（规则变更、失败结论、阻断决策），其余允许在计划收尾批量回写。
- runtime 强制字段至少包含：`trace_id/run_id`、记录类型（assumption/decision/change/validation）、时间戳、来源引用、当前状态。
- 允许延后补记，但必须记录延后原因与补记截止点；超过窗口需升级为显式风险条目。
- 同一事项多次回写采用“版本追加”而非覆盖，保持时间序审计链；当前有效版本通过状态字段标记。

### Claude's Discretion
- 触发“关键里程碑”的具体判定阈值（在不违背上述强约束前提下）可由 Claude 在研究/规划阶段细化。
- 双视图中的字段排布与展示顺序可由 Claude 按调用方消费成本优化。
- 延后补记窗口时长可由 Claude 结合现有流程负载给出默认值。

</decisions>

<specifics>
## Specific Ideas

- 保持与既有 Phase 2/3 的规则表达一致：沿用 `passed/warned/failed` 与 rule trace 语义，避免新增并行语义体系。
- 继续坚持“可执行输出”原则：每条非通过结果必须包含可执行下一步，而不仅是状态描述。

</specifics>

<deferred>
## Deferred Ideas

None - discussion stayed within phase scope.

</deferred>

---

*Phase: 04-validation-runtime-governance*
*Context gathered: 2026-03-18*
