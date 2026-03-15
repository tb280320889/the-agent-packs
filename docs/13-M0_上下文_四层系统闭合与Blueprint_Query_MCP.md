# 13 M0 上下文：四层系统闭合与 Blueprint Query MCP

## 文档类型
单里程碑所需上下文文档

## 作用域
本文件把系统层母版中最容易在实现期被压扁的部分单独提出来冻结：  
**四层系统闭合 + Blueprint Query MCP 的最小接口面。**

## 本文必须回答的问题
1. 四层系统如何闭合。
2. Blueprint Query MCP 为什么是一等对象。
3. 最小 resources / tools / prompts 该长什么样。
4. route / expand / read / bundle 的职责边界是什么。
5. M1 能做到多薄，底线在哪里。

---

## 一、四层系统闭合的最低要求

四层系统必须至少形成下面这条链：

`Blueprint Markdown -> Compiler -> SQLite Index -> Blueprint Query MCP -> Agent`

### 如果缺任何一层，会发生什么
- 缺 Blueprint：知识无结构
- 缺 Compiler / Index：查询无稳定图谱
- 缺 MCP：agent 只能自由搜文档
- 缺 minimal bundle：agent 上下文失控

---

## 二、Blueprint Query MCP：首期固定接口面

### A. Resources
至少保留以下 URI 面：

1. `blueprint://node/{id}`
   - 返回节点摘要或正文

2. `blueprint://children/{id}`
   - 返回 children 节点列表

3. `blueprint://required/{id}`
   - 返回 required_with 节点列表

4. `blueprint://bundle/{bundle_id}`
   - 返回一次已构建好的最小上下文包

### B. Tools
首期至少实现：

1. `route_query`
2. `expand_node`
3. `read_node`
4. `build_context_bundle`

二阶段可补：
5. `validate_blueprint_graph`
6. `rebuild_index`

### C. Prompts
首期至少固定以下 prompt 名称：
1. `route-task`
2. `expand-subdomain`
3. `debug-validator-failure`

---

## 三、每个 tool 的职责边界

### 1. route_query
职责：
- 只做 L0/L1 路由
- 返回候选、得分、理由、must_include

不负责：
- 自由全文搜索
- 进入执行态细节
- 拼装大上下文

### 2. expand_node
职责：
- 从当前节点按关系展开
- 受 mode / level_limit / max_children 约束

不负责：
- 自由遍历图谱
- 无边界扩展深层节点

### 3. read_node
职责：
- 读取指定节点的摘要或正文指定 section

不负责：
- 任意模糊检索
- 代替 route

### 4. build_context_bundle
职责：
- 构造当前步骤所需的最小上下文包
- 输出 main / required / execution_children / deferred

不负责：
- 打包未来可能用到的所有节点
- 为了“完整”而塞入大量正文

### 5. validate_blueprint_graph
职责：
- 校验 Blueprint 文档系统自身的完整性

### 6. rebuild_index
职责：
- 从 Git 工作树重新 parse 并更新 SQLite 索引

---

## 四、agent 的固定使用顺序

### 总编排阶段
允许：
- `route_query(L0)`
- `route_query(L1)`
- `build_context_bundle`

### 子领域编排阶段
允许：
- 读主节点 L1 摘要
- 读必带横线 L1 摘要
- 视情况带少量 L2

### 执行态阶段
允许：
- 只补必要 L2
- 项目态事实改走代码 / 配置 / 运行时资源

### validator 阶段
允许：
- 先依赖检查
- 再 `expand_node`
- 最后才开 L3

---

## 五、M1 的“最薄实现”底线

M1 可以很薄，但不能薄到以下四点缺任意一个：

1. 能 parse frontmatter
2. 能把节点与边关系写入最小索引
3. 能按 route 选中主节点
4. 能构造最小 context bundle

### M1 不必一开始就做
- 复杂全文检索
- 高级 ranker
- 完整 graph validator
- 增量重建优化
- 大量缓存
- 可视化面板

---

## 六、M0 在这里必须冻结什么
- MCP surface 名称
- resource URI 语义
- tool 名称
- prompt 名称
- route / expand / read / bundle 的职责边界
- agent 的固定使用顺序

这些对象如果在 M1 里一边实现一边改，就会导致系统层母版再次松动。

---

## 七、冻结清单（可直接引用）

### 1. MCP surface 名称
- Blueprint Query MCP

### 2. Resource URI 语义
- `blueprint://node/{id}`: 节点摘要或正文
- `blueprint://children/{id}`: children 列表
- `blueprint://required/{id}`: required_with 列表
- `blueprint://bundle/{bundle_id}`: 已构建的最小上下文包

### 3. Tool 名称与职责边界
- `route_query`: 只做 L0/L1 路由
- `expand_node`: 受限展开关系
- `read_node`: 读指定节点摘要或正文
- `build_context_bundle`: 产出 main/required/execution_children/deferred
- `validate_blueprint_graph`: 校验图谱完整性（第二阶段）
- `rebuild_index`: 重建 SQLite 索引（第二阶段）

### 4. Prompt 名称
- `route-task`
- `expand-subdomain`
- `debug-validator-failure`
