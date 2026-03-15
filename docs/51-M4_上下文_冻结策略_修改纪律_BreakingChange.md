# 51 M4 上下文：冻结策略、修改纪律、Breaking Change

## 文档类型
单里程碑所需上下文文档

## 作用域
定义 Phase 1 冻结面、修改纪律与 breaking change 治理。  
它是防止第二领域线接入时把第一阶段壳层改烂的核心文档。

## 本文必须回答的问题
1. 哪些对象已经冻结。
2. 哪些改动属于非 breaking。
3. 哪些改动属于 breaking。
4. 谁有权发起 breaking change。
5. 第二领域线能在什么边界内扩展。

---

## 一、Phase 1 冻结对象

以下对象默认进入冻结态：
- 系统层总设计判断
- 核心对象名
- 顶层 envelope shape
- route 优先级
- partial / handoff / failed 语义
- Blueprint L0/L1/L2/L3 职责
- Blueprint Query MCP surface 名称
- workflow package 标准模板
- Validation Plan / Validator Result / Activation Result 的顶层语义
- 固定主任务 regression 基线

---

## 二、允许修改的范围（非 breaking）

### 可以改
- 新增可选字段
- 新增 domain-specific contract
- 新增包内模板
- 新增 validator
- 新增 artifact 子类型
- 在不改 envelope 语义前提下增加更细的 findings 字段

### 不算 breaking 的前提
- 不删除旧字段
- 不改旧字段语义
- 不改旧状态枚举
- 不改 route 基线顺序

---

## 三、breaking change

以下都视为 breaking：
- 改系统层主判断
- 改 envelope 顶层语义
- 删除或重命名顶层关键字段
- 改 route 优先级顺序
- 改 bundle 最小结构
- 改 MCP surface 名称
- 改 workflow package 的最小模板组成
- 改 Validation Plan / Activation Result 的状态逻辑

---

## 四、breaking change 的处理流程

1. 先写变更提案
2. 说明为什么不能用加法扩展解决
3. 评估对既有 package / validator / tests 的影响
4. 给迁移路径
5. 只有在冻结面被显式解冻后才允许落地

### 禁止
- 在某个 pack 的开发过程中顺手改全局 envelope
- 先实现再补变更说明
- 只在口头上说“影响不大”

---

## 五、第二领域线的修改纪律

第二领域线接入时：
- 只能新增 domain / subdomain / pack / validator / templates
- 只能在局部 contract 做加法
- 不得修改系统层母版
- 不得修改首期切片已通过 regression 的公共壳层

---

## 六、冻结的真正目的
冻结不是阻止演进。  
冻结是为了保证：

1. 第二条领域线是在 **复用骨架**
2. 而不是借着扩域 **重造系统**
