# Handoff: 改造计划 v1 M4 迁移实施、兼容验证与准入演练

## 1. 交接对象
- 来源 bead：the-agent-packs-2ji
- 下一 bead：无强制后续 bead；若进入实现迁移，需新建 bead 承接
- 来源里程碑：M4
- 目标角色：项目内部维护 agent / 迭代开发子 agent / Verifier

## 2. 已完成什么
- 为 M4 开发指导补充了输入边界、设计顺序、必须回答问题、分阶段实施要求、非目标、完成标准、最小冻结对象与 Verifier/Handoff 要求。
- 为 M4 上下文文档补充了迁移实施阶段表、迁移差异最小检查项、分层回滚策略、兼容验证矩阵、第二领域线准入演练清单、最终收口标准与验收结论类型。
- 新增 `docs/改造计划v1/52-M4_验收结果与收口结论.md`，逐项确认 M4 已达到验收收口标准。
- 明确 M4 当前产物属于“可执行迁移计划已冻结并已验收收口”，而不是“实现迁移已完成”。
- 新增了 M4 阶段快照：`docs/改造计划v1/context-snapshots/2026-03-15-m4-migration-readiness.md`。

## 3. 下一位 agent 可直接依赖什么
- M4 开发指导：`docs/改造计划v1/50-M4_迁移实施_兼容验证与准入演练_开发指导.md`
- M4 上下文：`docs/改造计划v1/51-M4_上下文_迁移步骤_回滚策略与验收清单.md`
- M4 验收结论：`docs/改造计划v1/52-M4_验收结果与收口结论.md`
- M4 阶段快照：`docs/改造计划v1/context-snapshots/2026-03-15-m4-migration-readiness.md`
- M3 阶段快照：`docs/改造计划v1/context-snapshots/2026-03-15-m3-knowledge-ingestion.md`
- M2/M3 真相源与实现：`workflow-packages/registry.json`、`internal/registry/registry.go`

## 4. 下一位 agent 必须先做什么
- 先 claim：若要继续做真实实现迁移，必须新建并认领新的 bead，不得继续沿用已收口的 M4 bead。
- 先阅读：`docs/改造计划v1/52-M4_验收结果与收口结论.md`、`docs/改造计划v1/50-M4_迁移实施_兼容验证与准入演练_开发指导.md`、`docs/改造计划v1/51-M4_上下文_迁移步骤_回滚策略与验收清单.md`、`docs/改造计划v1/context-snapshots/2026-03-15-m4-migration-readiness.md`。
- 先验证：`go test ./...`、`bd show the-agent-packs-2ji --json`。

## 5. 不要做什么
- 不要把 M4 文档冻结误标成“正式实现迁移已完成”。
- 不要为了加快演练而临时放宽裸名禁令、attach-only、候选空间裁剪或资产落位闸口。
- 不要在没有新增 bead 与 snapshot 的情况下直接展开第二领域线正式接入。
- 不要再修改已收口的 M4 结论去承载实现执行细节；实现问题应进入新 bead。

## 6. 风险与未决项
- 若后续实现层跳过 M4 的负例校验，第二领域线正式接入时仍可能重新引入越界路由与未归属资产问题。
- 当前未决项已不属于 M4 收口范围，而属于后续实现阶段是否按 M4 计划执行的问题。

## 7. 推荐下一动作
- M4 已收口完成；若进入下一阶段，应基于 `docs/改造计划v1/52-M4_验收结果与收口结论.md` 新建实现执行 bead，并把 M4 当作已冻结输入而非继续编辑对象。
