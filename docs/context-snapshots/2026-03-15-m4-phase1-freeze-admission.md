# Context Snapshot: M4 Phase1 冻结与扩域准入

## 1. 当前阶段
- 所属里程碑：M4
- 关联 bead：the-agent-packs-3w0
- 当前状态：in_progress

## 2. 当前事实
- 当前要解决的问题：将 M0~M3 已完成链路收口为 Phase1 冻结面，并给出第二领域线准入与复用执行骨架。
- 当前已完成内容：已增强 M4 三份主文档，补齐冻结对象可追溯清单、non-breaking/breaking 判定与流程、准入检查清单、复用步骤输入输出、复用失败信号与下游交接包定义。
- 当前尚未完成内容：第二领域线尚未实际落地完整 pack；当前阶段只完成准入规则与复用手册。

## 3. 已冻结对象
- 系统主链路：`Blueprint -> compiler -> SQLite index -> Blueprint Query MCP -> minimal context bundle -> workflow package -> validator -> activation result / handoff`。
- 顶层 envelope 对象：Activation Request、Routing Result、Context Bundle、Artifact、Validation Plan、Validator Result、Activation Result、Handoff Bundle。
- route 优先级：`target_pack > target_domain > triggers/anti_triggers > selected_files/fragments > context_hints`。
- 状态语义：`completed/partial/handoff/failed` 与 `failed > handoff > partial > completed` 裁决基线。
- 结构模板：workflow package 标准模板与 tests 基本形状。

## 4. 当前输入
- 上游交付物：M3 验证闭环实现与回归基线。
- 依赖文档：`docs/50-M4_Phase1冻结与扩域准入_开发指导.md`、`docs/51-M4_上下文_冻结策略_修改纪律_BreakingChange.md`、`docs/52-M4_上下文_第二领域线准入与可复用骨架.md`。
- 依赖协作文档：`docs/04-多Agent接力开发与bd协作规则.md`、`docs/05-统一Handoff_ContextSnapshot_共享文档模板.md`。

## 5. 当前输出
- 已更新文件：
  - `docs/50-M4_Phase1冻结与扩域准入_开发指导.md`
  - `docs/51-M4_上下文_冻结策略_修改纪律_BreakingChange.md`
  - `docs/52-M4_上下文_第二领域线准入与可复用骨架.md`
- 已产出文件：
  - `docs/context-snapshots/2026-03-15-m4-phase1-freeze-admission.md`
  - `docs/handoffs/the-agent-packs-3w0-handoff.md`

## 6. 风险与阻塞
- 风险：第二领域线候选若依赖重宿主运行时，可能频繁触发 breaking 讨论并拖慢推进。
- 风险：若不严格执行准入清单，可能出现“口头准入”导致边界漂移。
- 阻塞：无。
- 是否需要人工决策：否（当前可继续按清单推进）。

## 7. 下一步建议
- 建议下一个 bead：the-agent-packs-3w0（继续，进入准入评审与候选领域选择）。
- 建议先执行的命令：`go test ./...`、`bd status --json`。
- 建议先阅读的文档：`docs/51-M4_上下文_冻结策略_修改纪律_BreakingChange.md`、`docs/52-M4_上下文_第二领域线准入与可复用骨架.md`。
