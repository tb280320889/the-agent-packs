# 22 M1 上下文：Routing、Bundle、最薄入口

## 作用域
本文只回答：
1. activation entry 的最小行为是什么
2. route 如何选主 pack / 主节点
3. context bundle 最少要带什么
4. 上下文不足时怎么诚实返回

## 一、最薄 activation entry

### 输入
Activation Request 最少要求：
- `request_id`
- `task`

常用增强字段：
- `target_pack`
- `target_domain`
- `selected_files`
- `config_fragments`
- `context_hints`
- `allowed_operations`
- `desired_outputs`

### 输出
最薄 activation entry 只允许返回两类东西：
- `route result`
- `activation result`

### 行为边界
它只做：
- 解析 request
- 触发 route
- 触发最小 bundle
- 在明显上下文不足时返回 partial / failed

它不做：
- 深分析
- artifact 深生成
- validator 深运行

## 二、Routing 输入与输出

### 输入源
- task
- target_pack
- target_domain
- selected_files
- config_fragments
- context_hints
- blueprint triggers / anti_triggers
- package capability declarations

### 输出
最少返回：
```json
{
  "main_pack": "wxt-manifest",
  "main_blueprint_node": "L1.wxt.manifest",
  "required_packs": ["security-permissions", "release-engineering"],
  "required_blueprint_nodes": [
    "L1.security.permissions",
    "L1.release.browser-store"
  ],
  "route_reason": "manifest task 且存在 permissions / store submission 信号",
  "recommended_validators": [
    "validator-core-output",
    "validator-domain-wxt-manifest"
  ],
  "recommended_artifacts": ["manifest-review.md"]
}
```

## 三、Routing 最小行为

### 规则 1：显式 target_pack 绝对优先
只要 `target_pack` 存在且合法，模糊 hints 不能覆盖它。

### 规则 2：target_domain 只在没有 target_pack 时生效
它能缩小到某个领域，但不能直接省略子问题判断。

### 规则 3：task + triggers 是主路由信号
只有当 task 足够具体，才允许进入 L1 / pack。

### 规则 4：selected_files / fragments 只做校正，不做扩权
它们用于帮助判断，不允许变成默认扫描起点。

### 规则 5：上下文不足时必须退回 L0 或 partial
不能强行猜测宿主结构。

## 四、Bundle 最小结构

### 目标
Bundle 的目标不是“多带资料”，而是“只带完成当前 task 需要的最小增强上下文”。

### 最小字段
```json
{
  "bundle_id": "cb-001",
  "main_pack": "wxt-manifest",
  "main_blueprint_node": "L1.wxt.manifest",
  "required_packs": ["security-permissions", "release-engineering"],
  "required_blueprint_nodes": [
    "L1.security.permissions",
    "L1.release.browser-store"
  ],
  "execution_children": [],
  "deferred_nodes": [],
  "recommended_validators": [
    "validator-core-output",
    "validator-domain-wxt-manifest"
  ],
  "recommended_artifacts": ["manifest-review.md"]
}
```

### Bundle 行为原则
- Summary first
- Required first
- Deferred by default
- Minimality
- Traceability

## 五、典型流

### 流 A：Manifest 主任务
输入：
- task：`review WXT manifest permissions for browser store submission`

应得到：
- main_pack：`wxt-manifest`
- required：`security-permissions` + `release-engineering`
- main node：`L1.wxt.manifest`
- artifacts：`manifest-review.md`
- validators：core + domain

### 流 B：Content Script 任务
输入：
- task：`assess risks in current content script page injection approach`

应得到：
- main_pack：`wxt-content-script`
- required：`wxt-manifest`（如需要 manifest 关联）与 `security-permissions`
- execution child：`L2.wxt.content-script.page-injection`

### 流 C：泛任务 / 上下文不足
输入：
- task：`help me with my WXT extension`

应得到：
- 只 route 到 `L0.wxt` 或返回 partial
- 推荐补充：
  - selected_files
  - config_fragments
  - 更具体 task
- 不允许自动扫描仓库

## 六、最薄入口的完成条件
M1 里该模块只有在以下条件满足时才算完成：
1. 可以接受 activation request
2. 可以产出 route result
3. 可以产出最小 bundle
4. 可以诚实返回上下文不足
5. route 不会被模糊 hint 误导

## 非目标
本文不定义：
- 完整 pack 目录骨架
- artifact 正文模板
- validator 具体 checks
这些由后续文档负责。
