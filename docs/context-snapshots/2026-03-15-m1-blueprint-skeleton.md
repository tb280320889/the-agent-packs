# Context Snapshot: M1 Blueprint 知识骨架与最薄入口

## 1. 当前阶段
- 所属里程碑：M1
- 关联 bead：the-agent-packs-695
- 当前状态：completed

## 2. 当前事实
- 当前要解决的问题：落地 Blueprint 首批节点与最薄 compiler/query/entry 骨架，保证 M2 可直接复用。
- 当前已完成内容：创建 blueprint 目录与最小节点集；新增 compiler/query/entry 说明与 fixtures 样例；补最薄 compiler/query/entry 的可运行脚手架并完成 smoke case 验证；完成 P0/P1 加固（route 优先级与回退策略修正、最小 TDD 测试集）；新增 route 冲突与 frontmatter 边界测试。
- 当前尚未完成内容：无（已完成最薄链路与增强 TDD 验证）。

## 3. 已冻结对象
- M0 协议与对象冻结：见 `docs/12-M0_上下文_协议契约与缺口补齐.md`
- MCP surface 冻结：见 `docs/13-M0_上下文_四层系统闭合与Blueprint_Query_MCP.md`

## 4. 当前输入
- 上游交付物：M0 冻结文档与 handoff
- 依赖文档：`docs/20-23`，`docs/12-13`
- 依赖 schema / 模板 / fixtures：`blueprint/frontmatter-examples.md`

## 5. 当前输出
- 已产出文件：
  - `blueprint/README.md`
  - `blueprint/frontmatter-examples.md`
  - `blueprint/schema.md`
  - `blueprint/L0/wxt/overview.md`
  - `blueprint/L0/security/overview.md`
  - `blueprint/L0/release/overview.md`
  - `blueprint/L1/wxt/manifest.md`
  - `blueprint/L1/security/permissions.md`
  - `blueprint/L1/release/store-review.md`
  - `blueprint/L2/wxt/manifest/permissions-review.md`
  - `blueprint/L2/wxt/manifest/browser-overrides.md`
  - `blueprint/L2/security/permissions/minimization.md`
  - `blueprint/L2/release/store-review/browser-extension-checklist.md`
  - `blueprint/L3/wxt/manifest/edge-cases.md`
  - `blueprint/L3/release/store-review/browser-specific-edge-cases.md`
  - `tools/compiler/README.md`
  - `tools/query-mcp/README.md`
  - `tools/activation-entry/README.md`
  - `tools/compiler/compiler.py`
  - `tools/query-mcp/query_mcp.py`
  - `tools/activation-entry/activation_entry.py`
  - `tests/test_m1_minimal.py`
  - `fixtures/README.md`
  - `fixtures/activation-request.sample.json`
  - `fixtures/context-bundle.sample.json`
  - `fixtures/route-result.sample.json`
  - `fixtures/activation-result.sample.json`
  - `fixtures/smoke-run.md`
  - `blueprint/index/validation-report.json`
  - `blueprint/index/missing-reference-report.json`

## 6. 风险与阻塞
- 风险：当前仍未覆盖完整 YAML 特性（如多行字符串、复杂嵌套），后续若升级 parser 需补兼容回归。
- 阻塞：无
- 是否需要人工决策：否

## 7. 下一步建议
- 建议下一个 bead：进入 M2（首个完整 workflow package：wxt-manifest）
- 建议先执行的命令：`python -m unittest tests/test_m1_minimal.py`
- 建议先阅读的文档：`docs/30-M2_首个完整Pack_wxt_manifest_开发指导.md`
