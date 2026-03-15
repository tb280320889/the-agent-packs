# Handoff: 改造计划 v1 M6 运行验证与 V1 收口

## 1. 交接对象
- 来源 bead：the-agent-packs-18k
- 下一 bead：待新建
- 来源里程碑：M6
- 目标角色：项目内部维护 agent / 迭代开发子 agent / Verifier

## 2. 已完成什么
- 已完成 M6 运行验证，确认 registry、route、asset、validator、handoff 的执行链路均已形成结构化验证闭环。
- 已把 activation 的 validator / handoff 生成逻辑改为优先消费注册表与 context bundle 输出，而不是只依赖 `wxt-manifest` 硬编码分支。
- 已补齐 capability 包 validator / artifact 验证与 handoff 断言测试。
- 已新增 M6 验收结论：`docs/改造计划v1/72-M6_验收结果与V1收口结论.md`
- 已新增 M6 收口快照：`docs/改造计划v1/context-snapshots/2026-03-15-m6-runtime-governance.md`

## 3. 下一位 agent 可直接依赖什么
- V1 最终收口结论：`docs/改造计划v1/72-M6_验收结果与V1收口结论.md`
- M6 收口快照：`docs/改造计划v1/context-snapshots/2026-03-15-m6-runtime-governance.md`
- 已对齐实现：`internal/activation/activation.go`、`internal/query/query.go`、`internal/registry/registry.go`、`internal/validator/core_output.go`
- 已补齐测试：`tests/m2_registry_test.go`、`tests/m2_wxt_manifest_test.go`、`tests/m3_validation_closure_test.go`

## 4. 下一位 agent 必须先做什么
- 先确认是否已新建后续 bead；M6 已收口，不应继续混写。
- 先阅读：`docs/改造计划v1/72-M6_验收结果与V1收口结论.md`
- 先验证：`go test ./...`

## 5. 不要做什么
- 不要把“继续扩域开发”伪装成 M6 的未完成尾巴。
- 不要重新引入按 pack 名硬编码 validator / handoff 的逻辑分支。
- 不要在无 bead、snapshot、handoff 的情况下继续宣称新增领域线已经准入。

## 6. 风险与未决项
- v1 已完成，但真实第二主域业务样板尚未在仓库落地；这不是 M6 阻塞项，而是后续 roadmap 项。
- 后续若新增新的 domain validator，必须继续保持注册表与 bundle 驱动，避免逻辑漂移。

## 7. 推荐下一动作
- 为“真实第二主域样板接入”或“v1 后续扩域 roadmap”新建 bead，并在 v1 已冻结边界内推进。
