# Domain Pitfalls

**Domain:** Agent Pack 增强包生产与上下文交付系统  
**Researched:** 2026-03-16

## Critical Pitfalls

Mistakes that cause rewrites or major issues.

### Pitfall 1: 语义真相源与实现真相源漂移
**What goes wrong:** 代码已经改动，但 `docs/AIDP` 与 runtime 账本未同步，导致后续 agent 按旧语义执行。  
**Why it happens:** 把 `.planning` 或对话当成真相源，忽略 AIDP 回写。  
**Consequences:** 需求-路由-验证链条断裂，阶段成果不可复用。  
**Prevention:** 每次关键变更同步更新 `runtime/02,03,06`，必要时更新 `core/`。  
**Detection:** 同一问题在不同 agent 会话出现相互矛盾结论。

### Pitfall 2: capability 越级成为全局入口
**What goes wrong:** 横线能力（security/release/auth/payment）在第一轮路由与主域竞争。  
**Why it happens:** 追求“命中率”而放宽候选空间约束。  
**Consequences:** 主域识别失真、上下文串线、交付包不可解释。  
**Prevention:** 强制执行 `global -> domain -> capability`，capability 仅 attach。  
**Detection:** 路由日志中 capability 在主域未确认前被选为 primary。

### Pitfall 3: 手工解析导致静默错误
**What goes wrong:** `package.yaml/frontmatter` 格式稍变后被误解析或字段丢失。  
**Why it happens:** 使用字符串 split 替代标准解析器。  
**Consequences:** registry 映射错误、validator 丢配、激活异常。  
**Prevention:** 引入 `yaml.v3` + KnownFields/schema 校验；解析失败即阻断。  
**Detection:** 同一文件人工可读正确，但系统行为异常且无显式报错。

### Pitfall 4: 索引重建非事务化
**What goes wrong:** rebuild 过程中失败，留下半成品索引。  
**Why it happens:** DROP/DELETE + 分步写入但无事务/原子替换。  
**Consequences:** route_query 结果不稳定，激活随机失败。  
**Prevention:** 事务包裹全流程或临时 DB 原子替换。  
**Detection:** 编译失败后路由结果异常波动。

## Moderate Pitfalls

### Pitfall 1: 硬编码默认主包回退
**What goes wrong:** 未匹配时默认回退固定包（当前审计为 wxt-manifest）。  
**Prevention:** 改为显式错误/partial，并要求 registry 配齐。

### Pitfall 2: 运行目录敏感导致 registry 加载失败
**What goes wrong:** 在不同 cwd 启动时无法定位项目根与 registry。  
**Prevention:** 增加显式 root 参数与更稳健的路径解析策略。

### Pitfall 3: 过早扩域
**What goes wrong:** 第二主域接入先于基础治理修复。  
**Prevention:** 先完成解析稳定化、路由解释性、契约验证，再扩域。

## Minor Pitfalls

### Pitfall 1: CLI 错误处理不可组合
**What goes wrong:** 参数错误直接 `os.Exit`，上层难以复用。  
**Prevention:** 对核心逻辑返回 error，入口层再决定退出策略。

### Pitfall 2: MCP 输入缺少体积限制
**What goes wrong:** 大输入可能引起资源消耗异常。  
**Prevention:** 使用 `io.LimitReader` 并做字段级校验。

## Phase-Specific Warnings

| Phase Topic | Likely Pitfall | Mitigation |
|-------------|---------------|------------|
| Phase 1: 解析与索引稳定化 | 局部修复后仍留手工解析路径 | 一次性替换并加回归测试（成功/失败/边界） |
| Phase 2: 路由治理 | capability 越级竞争复发 | 增加候选空间断言测试与 explain 输出 |
| Phase 3: 消费契约验证 | “最小”与“完整”只验证其一 | 同时校验包含必要项与排除无关项 |
| Phase 4: 多域扩展试点 | 新域沿用旧硬编码分支 | 先完成配置化域推断，再接新域 |

## Sources

- .planning/codebase/CONCERNS.md（HIGH，已发现风险）  
- docs/AIDP/core/05-范围与边界.md（HIGH，反向定义 anti-feature）  
- docs/AIDP/core/06-业务规则与关键对象.md（HIGH，路由与交付契约）  
- docs/AIDP/protocol/04-默认假设协议.md（HIGH，假设治理）  
- docs/AIDP/protocol/10-增强开发与迭代协议.md（HIGH，变更流程）  
- MCP spec security section（MEDIUM，工具与数据暴露安全）: https://modelcontextprotocol.io/specification/2025-06-18
