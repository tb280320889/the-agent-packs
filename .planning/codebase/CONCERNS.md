# Codebase Concerns

**Analysis Date:** 2026-03-16

## Tech Debt

**YAML 解析实现过于手工且易碎：**
- Issue: `package.yaml` 的解析通过自定义字符串拆分完成，缺少 YAML 语法覆盖与错误提示。
- Files: `internal/registry/registry.go`
- Impact: YAML 新增字段/缩进风格变化可能被静默忽略，导致注册表与包清单不一致。
- Fix approach: 使用标准 YAML 解析库并引入结构化 schema 校验，覆盖 `depends_on/validators/artifacts` 等列表字段。

**前置校验/解析器对 frontmatter 处理过于简化：**
- Issue: `parseFrontmatter` 仅支持简单 `key: value` 与 `- item` 列表，缺少引号/多行/嵌套处理。
- Files: `internal/compiler/compiler.go`
- Impact: Blueprint 文档格式稍有变化即编译失败或误解析，降低扩展性。
- Fix approach: 使用成熟 frontmatter/YAML 解析库，确保列表/字符串/空值/多行语法正确解析。

**SQLite 索引重建全量删除，无事务保护：**
- Issue: `writeIndex` 先 `DROP`/`DELETE` 再逐条写入，未使用事务，失败时可能留下空索引。
- Files: `internal/compiler/compiler.go`
- Impact: 编译过程中断将导致索引不完整，影响路由与激活逻辑。
- Fix approach: 使用事务包裹 schema 变更与写入；采用临时 DB 文件 + 原子替换。

**默认主包回退硬编码：**
- Issue: 当路由节点无法映射包名时默认回退到 `wxt-manifest`。
- Files: `internal/activation/activation.go`
- Impact: 非 WXT 场景会被错误映射到固定包，造成错误输出与验证结果。
- Fix approach: 取消硬编码回退，要求注册表必须提供映射，否则返回明确错误或 `partial`。

## Known Bugs

**命令行在缺少 args 时直接退出：**
- Symptoms: CLI 调用缺少参数即 `os.Exit(1)`，无法以库形式复用或在上层捕获错误。
- Files: `cmd/agent-pack-mcp/main.go`
- Trigger: 未传入子命令或子命令错误。
- Workaround: 上层调用前先做参数完整性校验。

## Security Considerations

**外部输入路径直接读取：**
- Risk: `activate` 直接读取请求文件路径，缺少路径白名单或大小限制。
- Files: `internal/activation/activation.go`, `cmd/agent-pack-mcp/main.go`
- Current mitigation: 无显式限制。
- Recommendations: 添加路径限制（仅允许工作目录或固定目录），并检查文件大小/格式。

**MCP JSON 输入缺少体积限制与强校验：**
- Risk: 标准输入 JSON 解析没有最大长度限制，潜在 DoS 风险。
- Files: `cmd/agent-pack-mcp/main.go`
- Current mitigation: 依赖 Go 标准库 JSON 解码器。
- Recommendations: 使用 `io.LimitReader` 限制输入大小，并在解析后做字段存在性校验。

## Performance Bottlenecks

**编译器全量遍历与全量重建：**
- Problem: 每次编译都遍历 `blueprint/` 下所有 Markdown 并重建 SQLite。
- Files: `internal/compiler/compiler.go`
- Cause: 未实现增量更新或校验缓存。
- Improvement path: 引入 checksum 索引与增量更新策略，仅重建变化节点。

**路由评分策略以字符串 contains 为主：**
- Problem: 评分依赖 `strings.Contains`，在大型节点库中成本较高。
- Files: `internal/query/query.go`
- Cause: 无索引/预计算的 token 结构。
- Improvement path: 预计算 token 集合或使用 SQLite FTS 作为候选过滤。

## Fragile Areas

**注册表默认加载依赖工作目录：**
- Files: `internal/registry/registry.go`
- Why fragile: `findProjectRoot` 依赖 `go.mod` 与 `workflow-packages/registry.json` 的相对路径；运行目录变化即失效。
- Safe modification: 对 `Default()` 增加可配置根路径或显式参数，避免隐式 `os.Getwd()` 依赖。
- Test coverage: `tests/m2_registry_test.go` 覆盖正常加载路径，但未覆盖不同工作目录。

**Validator 注册表静态写死：**
- Files: `internal/validator/registry.go`, `internal/validator/runner.go`
- Why fragile: 新增 validator 需要改代码发布，无法由包配置动态扩展。
- Safe modification: 将注册表与 `workflow-packages/registry.json` 或插件机制联动。
- Test coverage: `tests/m3_validation_closure_test.go` 覆盖当前两个 validator，但无扩展场景。

**路由域推断逻辑过于特化：**
- Files: `internal/query/query.go`
- Why fragile: `inferMainDomain` 只识别 `wxt`/`manifest`/`browser extension`，对新域无通用规则。
- Safe modification: 将 domain 触发词配置化（来自 blueprint 或 registry），避免硬编码。
- Test coverage: `tests/m1_minimal_test.go` 覆盖 WXT 场景，不覆盖多域扩展。

## Scaling Limits

**节点/边关系无数据库约束：**
- Current capacity: SQLite 无外键与索引约束。
- Limit: 大规模节点时边关系完整性依赖编译器，易出现脏数据。
- Scaling path: 为 `nodes/edges` 增加索引与外键，或使用更强的校验管道。

## Dependencies at Risk

**SQLite 驱动固定为 `modernc.org/sqlite`：**
- Risk: 纯 Go SQLite 依赖可能在某些平台性能较弱。
- Impact: 大规模索引读写延迟升高。
- Migration plan: 预留可切换驱动（`mattn/go-sqlite3`）或抽象 DB 接口。

## Missing Critical Features

**缺少索引版本/迁移机制：**
- Problem: `blueprint/index/blueprint.db` 结构升级缺少版本标记与迁移流程。
- Blocks: 未来 schema 调整无法安全升级。

## Test Coverage Gaps

**编译器错误输出未覆盖异常路径：**
- What's not tested: `writeIndex` 失败、`writeReports` 写入失败等错误路径。
- Files: `internal/compiler/compiler.go`, `tests/m1_minimal_test.go`
- Risk: 失败时可能生成半成品索引或空报告。
- Priority: Medium

**MCP 输入解析边界缺少测试：**
- What's not tested: `parseParam*` 对异常 JSON 的处理与类型错误路径。
- Files: `cmd/agent-pack-mcp/main.go`
- Risk: 运行时错误难以定位，异常输入处理不稳定。
- Priority: Low

---

*Concerns audit: 2026-03-16*
