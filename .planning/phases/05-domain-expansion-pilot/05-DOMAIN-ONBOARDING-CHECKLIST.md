# Phase 05 新主域准入清单（五段式）

本清单用于后续任意新主域接入时的统一执行模板，确保遵循 `candidate-space-first`、`attach-only`、`最小且完整` 与 `run_id/runtime-ledger` 证据链约束。

## 命名治理

- 必备文件
  - `blueprint/L0/<domain>/overview.md`
  - `blueprint/L1/<domain>/<subdomain>.md`
  - `workflow-packages/<domain>-<subdomain>/package.yaml`
  - `workflow-packages/registry.json`
- 通过标准
  - package 名称符合 `domain-subdomain`
  - `category=workflow` 且 `activation_mode=direct`
  - `canonical_blueprint_node` 指向有效 L1 节点
- 失败阻断标准
  - registry 校验失败即阻断
  - 命名不一致（目录名/package.yaml/registry name）即阻断
- 推荐命令
  - `go test ./... -run "TestM2RegistryLoadsAndValidates|TestM5DomainExpansionRegistryOnboarding" -count=1`

## 路由契约

- 必备检查
  - 新主域可在开关开启时进入 L0/L1 候选空间
  - 开关关闭时不进入主候选
  - attach-only 仍不得进入 primary candidate
- 通过标准
  - `RouteQuery` 能选中新主域主包
  - 开关关闭返回 `ROUTE_NO_PRIMARY_CANDIDATE`
  - `DecisionBasis` 与 `CapabilityDecisions` 可解释
- 失败阻断标准
  - **P0 阻断：主域竞争越界**
  - **P0 阻断：attach-only 破坏**
- 推荐命令
  - `go test ./... -run "TestM5DomainExpansion(OnboardMonorepoRouteAndActivation|FeatureSwitchRollback)" -count=1`

## 交付契约

- 必备检查
  - `BuildContextBundle` 的 include/exclude decisions 非空
  - 目标域 required packs 与推荐 artifacts/validators 可追踪
  - 无关域节点不混入默认交付
- 通过标准
  - include 决策可解释“为何包含”
  - exclude 决策可解释“为何不包含”
  - 关键 required packs 不遗漏
- 失败阻断标准
  - 关键 required packs 丢失即阻断
  - include/exclude 决策缺失即阻断
- 推荐命令
  - `go test ./... -run "TestM5DomainExpansionWXTBundleContractNonRegression" -count=1`

## 验证回归

- 必备检查
  - WXT 主链路三件套回归通过：`route_query` / `build_context_bundle` / `activate`
  - 冲突语义下新主域不抢占 WXT
- 通过标准
  - `wxt-manifest` 仍是 WXT 目标请求主候选
  - `security-permissions` 与 `release-store-review` 仍在 WXT required packs
  - 冲突词场景在 `target_domain=wxt` 下仍落入 WXT
- 失败阻断标准
  - **P0 阻断：WXT 回归失败**
- 推荐命令
  - `go test ./... -run "TestM5DomainExpansion(WXTNonRegression|WXTBundleContractNonRegression|WXTActivationTraceNonRegression|ConflictDoesNotStealWXT)" -count=1`
  - `go test ./... -run "TestM2Activation.*" -count=1`

## runtime 证据链

- 必备检查
  - `ActivationResult.CurrentValidationRunID` 非空
  - `RuntimeLedger` 非空且存在 `record_type=validation` 当前项
  - warned 项必须记录 deadline 与 runtime ledger 留痕
- 通过标准
  - run_id 可关联 validation results
  - runtime ledger 可关联 trace_id/source refs
  - overdue 时可识别 `runtime-ledger-overdue`
- 失败阻断标准
  - **P0 阻断：run_id/ledger 证据链断裂**
- warn 处理要求
  - 必须写明整改 `deadline`
  - 必须写入 `runtime ledger` 留痕后方可放行

## 最小证据集合

每次新主域准入至少交付以下证据集合：

1. `command`：执行过的测试命令列表（含 go test 过滤表达式）
2. `artifact output`：关键结果摘要（主候选、required packs、status）
3. `run_id`：至少一个可追踪 validation run_id
4. `runtime-ledger refs`：对应 runtime-ledger 条目引用（record_type/version/is_current）
