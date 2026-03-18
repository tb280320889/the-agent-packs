# Architecture Patterns

**Domain:** Agent Pack 增强包生产系统  
**Researched:** 2026-03-16

## Recommended Architecture

采用 **Compile-Route-Bundle-Activate-Validate** 五段式生产架构，并以 `registry` + `AIDP 协议` 形成双真相源：

- **实现真相源（runtime truth）**：registry + SQLite index + activation artifacts
- **语义真相源（semantic truth）**：docs/AIDP（角色边界、交付契约、验收规则）

两者必须同步演进，禁止单独前进。

### Component Boundaries

| Component | Responsibility | Communicates With |
|-----------|---------------|-------------------|
| Blueprint Source | 维护 L0-L3 节点及 frontmatter | Compiler |
| Compiler | 解析节点、校验结构、构建 SQLite 索引 | SQLite Index, Validation Report |
| Registry | package 身份/依赖/validator/artifact 映射真相源 | Query, Activation, Validator |
| Query Router | 候选空间裁剪、评分、主域与 capability 选择 | SQLite Index, Registry, Bundle Builder |
| Bundle Builder | 生成最小且完整 context bundle（含 deferred） | Query Router, Activation |
| Activation Engine | 编排执行、组装 artifacts/handoff/result | Query Router, Bundle Builder, Validator Runner |
| Validator Runner | 执行 core/domain validators，回传结果 | Activation Engine, Registry |
| MCP/CLI Adapter | 对外暴露工具与命令 | Activation Engine, Query Router, Compiler |
| AIDP Runtime Artifacts | 记录假设/决策/验证/变更摘要 | 人与 agent 协作流程 |

### Data Flow

1. `Blueprint Markdown` 经 Compiler 解析并写入 `SQLite index`。  
2. 请求进入 Query Router，先裁候选空间（主域优先，capability attach-only）。  
3. Bundle Builder 生成 Context Bundle（main/required/children/deferred）。  
4. Activation Engine 结合 registry 构造 validation plan 与 artifacts。  
5. Validator Runner 执行校验，输出 validator result。  
6. Activation Result + Handoff Bundle 输出到调用方；运行态记录同步更新。

## Patterns to Follow

### Pattern 1: Candidate-Space-First Routing
**What:** 先按 `visibility_scope + activation_mode` 裁候选，再做打分。  
**When:** 所有 route_query/activate 路径。  
**Example:**
```typescript
// conceptual pseudocode
const candidates = filterByScopeAndMode(allNodes, request)
const scored = scoreWithinCandidates(candidates, request.task)
const primary = pickPrimaryDomain(scored)
const attached = attachCapabilities(primary, scored)
```

### Pattern 2: Registry-as-Source-of-Truth
**What:** 包身份、依赖、validator/artifact 映射仅信 registry，不信目录名猜测。  
**When:** 激活主包解析、validator 计划生成、artifact 归档。  
**Example:**
```typescript
// conceptual pseudocode
const pkg = registry.lookupByCanonicalName(node.package)
if (!pkg) throw new Error("unmapped package")
buildValidationPlan(pkg.validators)
```

### Pattern 3: Transactional Index Rebuild
**What:** 索引重建必须事务化或临时库原子替换。  
**When:** compile/rebuild_index。  
**Example:**
```typescript
beginTransaction()
rebuildNodesAndEdges()
writeReports()
commit()
```

### Pattern 4: Contracted Context Delivery
**What:** 交付给消费侧的是“最小且完整”的 bundle，并可解释包含/排除理由。  
**When:** build_context_bundle、activate。  
**Example:**
```typescript
bundle = {
  main,
  required,
  deferred,
  rationale: { included: [...], excluded: [...] }
}
```

## Anti-Patterns to Avoid

### Anti-Pattern 1: Hardcoded Main Pack Fallback
**What:** 无映射时默认回退到固定包（如 wxt-manifest）。  
**Why bad:** 会在新域场景输出错误结果且不易察觉。  
**Instead:** 显式报错或 partial，并要求 registry 完整映射。

### Anti-Pattern 2: Hand-Rolled YAML/Frontmatter Parsing
**What:** 用字符串拆分模拟 YAML/frontmatter 解析。  
**Why bad:** 格式稍变即误解析，导致静默数据错配。  
**Instead:** 使用标准解析库 + schema 校验。

### Anti-Pattern 3: Semantic/Runtime Drift
**What:** 代码实现变化但 AIDP/runtime 工件不更新。  
**Why bad:** 多 agent 协作时真相源分裂，错误决策累积。  
**Instead:** 变更同步更新决策日志、变更摘要、验证记录。

## Scalability Considerations

| Concern | At 100 users | At 10K users | At 1M users |
|---------|--------------|--------------|-------------|
| Route latency | 单库查询足够 | 需索引优化与 token 预计算 | 需分片（按域/租户）与缓存层 |
| Write concurrency | 单 writer 可接受 | 编译与查询需解耦（读写分离时序） | 必须分片/队列化，可能升级存储架构 |
| Registry growth | 手工维护可控 | 需要 schema lint + CI 校验 | 需要自动治理与版本化策略 |
| Validator extensibility | 静态注册可用 | 需要配置驱动注册 | 需要插件化与沙箱机制 |
| Context quality | 人工 spot check | 引入契约测试 | 需要系统化评测与回归基线 |

## Sources

- .planning/codebase/ARCHITECTURE.md（HIGH，当前实现分层）  
- .planning/codebase/CONCERNS.md（HIGH，反模式与脆弱点）  
- docs/AIDP/core/06-业务规则与关键对象.md（HIGH，架构规则）  
- docs/AIDP/core/08-技术约束与工程约定.md（HIGH，工程约束）  
- docs/AIDP/protocol/10-增强开发与迭代协议.md（HIGH，演进流程）  
- SQLite use-cases（MEDIUM，存储边界）: https://www.sqlite.org/whentouse.html
