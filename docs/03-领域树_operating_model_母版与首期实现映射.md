# 03 领域树 operating model 母版与首期实现映射

## 文档类型
系统层主文档

## 作用域
保存整套长期 operating model 母版，以及它与当前首期实现切片之间的映射关系。  
它的作用是防止两种常见错误：

1. 因为首期切片很小，就误以为长期 operating model 也很小  
2. 因为长期 operating model 很大，就在首期实现里一次铺满

## 本文解决什么
1. workflow package 的最小正确单位是什么。
2. 长期领域树母版包含哪些层次。
3. 横向能力线为什么必须显式存在。
4. 为什么长期优先级与当前首期切片选择看起来不同，但并不冲突。
5. 如何把“首个完整 pack”理解为方法验证，不是领域战略排序。

---

## 一、workflow package 的最小正确单位

一个 workflow package 不是一段 prompt。  
它最小应由以下对象组成：

- 一个编排层 skill
- 一组 MCP tools / resources / prompts
- 一组 validator
- 一套跨包交接规则
- 一套产物模板
- 一套退出条件

这套定义必须长期稳定。  
它保证包的边界清楚，且能从“建议”提升为“约束”。

---

## 二、长期 operating model 母版

### 顶层：总编排层
全系统入口是 global orchestrator。  
它负责：
- 识别产品形态
- 识别宿主环境
- 识别数据位置
- 识别身份体系
- 识别分发方式
- 选择一个主 workflow
- 挂接必要横线
- 生成验证计划

### 一级产品线
长期母版至少包含：
- Tauri
- OSS
- WXT
- Web3（Base / TON-TMA / Solana）
- 小程序
- HarmonyOS
- Axum
- SQLite / SurrealDB

### 横向能力线
长期母版必须显式存在：
- Agent-native tooling（CLI -> MCP）
- Design System / UIUX
- Identity / Auth
- Payment / Transaction
- Security / Privacy / Compliance
- DevOps / CI/CD / Release Engineering
- Observability / Supportability

### 为什么横线必须显式存在
因为以下前置协议必须长期成立：
- 安全前置协议
- 身份前置协议
- 发布前置协议
- 交易前置协议
- 可观测性回填协议

如果横线不存在，领域树不会闭合。

---

## 三、统一包模板与命名母版

### 标准内部结构
每个 workflow package 长期都应使用统一模板：
- README.md
- package.yaml
- skill/
- mcp/
- validators/
- contracts/
- templates/
- fixtures/
- tests/
- CHANGELOG.md

### 命名规范
统一使用：
- 编排包：`<domain>-orchestrator`
- 一级/二级能力包：`<domain>-<subdomain>`
- 横向能力包：`<capability-line>-<subdomain>`
- validator：`validator-<scope>-<name>`
- MCP server：`mcp-<surface>-server`

### 禁止
- misc
- common-stuff
- utils-everything
- tauri-all-in-one
- web3-core
- miniapp-tools
- platform-service

---

## 四、长期实施优先级母版

长期视角下，更稳的推进顺序是：

### Phase 1：主骨架
- Global orchestrator
- Tauri
- Axum
- SQLite
- OSS
- Agent-native tooling
- Security / Identity / Release 三条横线

### Phase 2：第二产品面
- WXT
- SurrealDB
- 小程序
- Design System
- Observability
- Payment

### Phase 3：宿主与链扩展
- Base Mini App
- TON / Telegram Mini Apps
- Solana
- HarmonyOS

---

## 五、为什么当前首期实现选择 wxt-manifest

看起来这和长期 Phase 1 不同。  
实际上不冲突。

### 长期优先级回答的问题
“长期 operating model 需要先保证哪些主干线具备战略完整性？”

### 当前首期切片回答的问题
“在当前仓库里，哪一条最适合作为 **第一条可验证闭环**？”

### 选择 wxt-manifest 作为首个完整 pack 的原因
- 边界清楚
- 输入形状明确
- artifact 容易定义
- validator 好写
- route/bundle 行为可控
- 易挂接 security / release 横线
- 不要求先把大宿主系统跑到很深

### 正确理解
`wxt-manifest` 不是“长期第一战略领域”。  
它是“当前仓库首个完整、可控、可验证的垂直切片”。

---

## 六、当前仓库的双层推进观

### 长期母版层
由 `03` 保存，不允许在首期切片中丢失。

### 首期实现层
由 `10~52` 驱动，要求：
- 面尽量小
- contract 尽量硬
- 闭环尽量早
- 可复制性尽量高

### 两者关系
首期实现是长期母版的一次 **受控投影**。  
它不是长期母版的删减版，更不是替代版。

---

## 七、所有后续里程碑必须保持的解释框架

1. 长期领域树负责“系统会长成什么”
2. 当前里程碑负责“当前只做哪一小块”
3. 首个完整 pack 负责“证明方法”
4. 第二领域线准入负责“证明可复制”
5. 任何时候都不允许因为首期实现而损坏长期母版
