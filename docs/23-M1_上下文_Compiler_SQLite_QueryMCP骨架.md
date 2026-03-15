# 23 M1 上下文：Compiler、SQLite、Query MCP 骨架

## 文档类型
单里程碑所需上下文文档

## 作用域
定义 M1 最薄编译链与最薄索引 / 查询骨架。  
目标是把系统层母版落成可运行最小链路，而不是做复杂平台。

## 本文必须回答的问题
1. compiler 的输入、处理、输出是什么。
2. SQLite 最小 schema 要长什么样。
3. query 层至少要支持哪些查询。
4. 哪些能力可以留到后面增强。
5. M1 交给 M2 的实现边界在哪里。

---

## 一、compiler 的最小职责

### 输入
- Git 工作树中的 Blueprint Markdown 文件

### 处理
1. parse frontmatter
2. 校验 id 与路径一致
3. 生成 `nodes`
4. 生成 `node_meta`
5. 生成 `edges`
6. 生成轻量 summary 索引
7. 输出报告

### 输出
- 最新 SQLite 索引
- 校验报告
- 失效引用报告

---

## 二、SQLite 最小 schema

首期至少要有：

### `nodes`
用于存主节点信息：
- `id`
- `level`
- `domain`
- `subdomain`
- `capability`
- `title`
- `summary`
- `path`
- `parent_id`
- `body_md`
- `entry_conditions_json`
- `stop_conditions_json`
- `checksum`
- `updated_at`

### `node_meta`
用于存多值元数据：
- `aliases`
- `triggers`
- `anti_triggers`
- `tags`

### `edges`
用于存图关系：
- `child`
- `required_with`
- `may_include`
- `excludes`

### `fts_nodes`
首期可选。  
若做，只索引轻量字段：
- title
- summary
- aliases
- triggers

> 不把正文全文作为 route 主索引。

---

## 三、最薄 Query MCP 骨架

### route_query
输入：
- level
- task
- host
- constraints
- max_results

输出：
- candidates
- must_include

### expand_node
输入：
- node_id
- mode
- max_children
- level_limit

### read_node
输入：
- node_id
- section

### build_context_bundle
输入：
- main_node
- include_required
- include_may_include
- include_children
- body_mode

---

## 四、首期不急着做的增强项
- 高级 FTS rank
- 增量重建优化
- bundle cache
- orphan node / dead edge 全量治理
- CI 增量重建
- 图可视化
- 大规模领域覆盖率报告

---

## 五、M1 的实现边界

### 必须完成
- 路径 / id 校验
- 最小索引写入
- 最小 route
- 最小 read
- 最小 bundle

### 可以留薄
- graph validator
- rebuild_index 工具化
- FTS
- anti_trigger 高级裁剪策略

---

## 六、M1 向 M2 的 handoff 内容
M2 应拿到：
- 可用的 Blueprint 目录与节点
- 可用的最薄索引
- 可用的 route / read / bundle 接口
- 已知限制清单

M2 不应再自己重写 compiler/query 的 envelope。
