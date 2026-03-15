# 20 M1 Blueprint 知识骨架与最薄入口 开发指导

## 文档类型
开发里程碑开发指导文档

## 作用域
M1 负责把 **Blueprint -> compiler/query -> context bundle -> activation entry** 的最薄骨架搭起来。  
它不做完整 pack，只做能支撑首个 pack 落地的基础设施最薄版。

## 本里程碑的目标
1. 建立 Blueprint 目录与命名规则
2. 固化 frontmatter 与节点模板
3. 落首批必要节点
4. 实现最薄 compiler
5. 实现最薄 route / read / bundle 查询层
6. 实现最薄 activation entry
7. 让 M2 能在稳定骨架上做首个完整 pack

## 必读文档
- 01
- 02
- 20
- 21
- 22
- 23

## 参考输入
- M0 冻结对象与 route / contract 规则

## 默认不读
- `30+`  
理由：M1 只做骨架，不做完整 pack 细节。

## 必须产出
1. `blueprint/` 基础目录
2. 首批节点文件
3. `frontmatter` 规范与合法 / 非法示例
4. 最薄 `compiler`
5. 最薄 `query`
6. 最薄 `context bundle builder`
7. 最薄 `activation entry`
8. 基础 fixtures 与 smoke cases

## 明确不产出
- 完整 WXT pack
- 复杂 validator
- 全量 FTS 排序能力
- 大量领域树节点

---

## 执行顺序

### 第一步：先固目录与 frontmatter
先让 Blueprint 能被写成稳定节点，而不是先写查询代码。

### 第二步：只落首批必要节点
只写首个闭环需要的节点。  
不要发散扩树。

### 第三步：做最薄 compiler
先把：
- 解析 frontmatter
- 校验 id/path
- 生成 nodes / node_meta / edges
跑通。

### 第四步：做最薄 query MCP
先做到：
- route
- read
- bundle

expand 与 graph validator 可以很薄，但接口名要固定。

### 第五步：做最薄 activation entry
先只做到：
- 收 activation request
- 识别 target_pack / target_domain
- 调 route / bundle
- 返回最基本 activation result 或 route result

---

## 完成标准
1. 首批节点有稳定知识源
2. route 可以选出主节点
3. bundle 不会默认塞大量节点
4. activation entry 已能消费 M0 的 request shape
5. M2 无需重造 compiler/query 壳层

## handoff 给 M2 的内容
- Blueprint 首批节点清单
- frontmatter 冻结版本
- compiler/query/entry 的最薄接口
- 已知限制，但不阻塞首个 pack 开发

## 本里程碑绝对不要做的事
- 以“顺手”为名扩很多节点
- 把 compiler 做成复杂平台
- 在 M1 就堆控制台或管理面
- 因为觉得后面有用就让 bundle 变大
