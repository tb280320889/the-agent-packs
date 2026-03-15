# 21 M1 上下文：Blueprint、命名、目录

## 作用域
本文只回答：
1. 命名怎么固定
2. 目录怎么固定
3. Blueprint 节点怎么写
4. 首批节点集合是什么

## 一、命名规则

### 1. Blueprint Node ID
固定格式：
`L<level>.<domain>[.<subdomain>][.<capability>]`

示例：
- `L0.wxt`
- `L1.wxt.manifest`
- `L1.security.permissions`
- `L2.wxt.content-script.page-injection`

规则：
- 小写
- 多词短横线
- 层级用 `.`
- 禁止 `general` / `misc` / `default`

### 2. Pack 命名
- 领域 pack：`<domain>-<subdomain>`
- 横线 pack：`<capability-line>-<subdomain>`

示例：
- `wxt-manifest`
- `wxt-content-script`
- `security-permissions`
- `release-engineering`

### 3. Validator 命名
固定格式：
`validator-<scope>-<target>`

示例：
- `validator-core-output`
- `validator-domain-wxt-manifest`

### 4. Artifact 命名
推荐：
- `<purpose>.md`
- 必要时 `<domain>-<purpose>.md`

首批固定示例：
- `manifest-review.md`
- `permission-audit.md`
- `content-script-compat-report.md`
- `store-release-checklist.md`

## 二、目录结构

### 顶层目录
```text
agent-pack/
├─ README.md
├─ AGENT_PACK.md
├─ docs/
├─ blueprint/
├─ tools/
├─ packages/
├─ skills/
├─ mcp/
├─ resources/
├─ fixtures/
└─ .github/
```

### docs/
- `concepts/`
- `protocols/`
- `standards/`
- `runbooks/`

### blueprint/
- `L0/`
- `L1/`
- `L2/`
- `L3/`

### tools/
- `activation-entry/`
- `blueprint-compiler/`
- `blueprint-query/`
- `validator-runner/`
- `artifact-generator/`

### packages/
- `orchestrators/`
- `domain-packs/`
- `cross-cutting/`
- `validators/`
- `contracts/`
- `shared/`

## 三、Frontmatter 固定字段

### 必填字段
- `id`
- `level`
- `domain`
- `subdomain`
- `title`
- `summary`
- `triggers`
- `required_with`
- `entry_conditions`
- `stop_conditions`

### 推荐字段
- `aliases`
- `anti_triggers`
- `may_include`
- `children`
- `excludes`
- `recommended_artifacts`
- `recommended_validators`

### 语义要求
每个字段都必须服务于：
- route
- bundle
- validation
- handoff

如果某个字段无法被这四类行为消费，就不应加入。

## 四、首批 blueprint 节点

### 第一阶段固定集合
- `L0.wxt`
- `L1.wxt.manifest`
- `L1.wxt.content-script`
- `L1.security.permissions`
- `L1.release.browser-store`
- `L2.wxt.content-script.page-injection`

### 层级职责
- L0：领域入口
- L1：子领域编排
- L2：执行策略 / 分支展开
- L3：边缘情况 / 升级分析

第一阶段不写满 L3。

## 五、首批节点写作要求

### L0.wxt
只回答：
- 什么任务属于 WXT
- 常见子问题
- 常见横线
- 可能 route 到哪些 L1

### L1.wxt.manifest
必须写清：
- 何时进入 manifest 视角
- 何时挂 security
- 何时挂 release
- 何时停止
- 何时交接

### L1.wxt.content-script
必须写清：
- 何时进入 content-script 视角
- 何时下钻 L2.page-injection
- 何时 handoff 到 permissions / runtime

### L1.security.permissions
必须写清：
- 何时权限成为主风险
- 何时作为 required_with
- 何时作为 handoff 目标

### L1.release.browser-store
必须写清：
- 何时 submission / packaging / store readiness 成为主问题
- 何时 checklist 是主输出
- 何时不该继续停留在 wxt-manifest

### L2.wxt.content-script.page-injection
只解决一个问题：
- 内容脚本的页面注入边界与页面上下文假设

## 六、最小 frontmatter 示例
```yaml
id: L1.wxt.manifest
level: L1
domain: wxt
subdomain: manifest
title: WXT Manifest Entry
summary: 处理 manifest-centered 任务，并在需要时挂接 permissions 与 release。
triggers:
  - manifest
  - permissions
  - browser store
required_with:
  - L1.security.permissions
entry_conditions:
  - task 明确涉及 manifest 或 browser permissions
stop_conditions:
  - 已明确主要 manifest 风险
  - 已可生成至少一个 review 或 handoff
recommended_artifacts:
  - manifest-review.md
recommended_validators:
  - validator-core-output
  - validator-domain-wxt-manifest
```

## 非目标
本文不定义：
- route 优先级细则
- bundle 输出结构
- activation entry 行为
这些由 22 文档负责。
