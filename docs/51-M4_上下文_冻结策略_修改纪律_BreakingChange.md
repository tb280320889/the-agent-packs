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

### 冻结对象可追溯清单（M4 版）
| 冻结对象 | 来源里程碑 | 约束级别 | 说明 |
|---|---|---|---|
| 系统层主链路（Git->Markdown->SQLite->MCP->bundle->pack->validator->result） | M0/M1/M2/M3 | 不可变 | 第二领域线只能复用，不可改写链路语义 |
| Activation Request / Routing Result / Context Bundle / Artifact / Handoff Bundle / Validation Plan / Validator Result / Activation Result | M0 | 不可变 | 顶层语义冻结，仅允许局部 contract 加法扩展 |
| route 优先级（target_pack > target_domain > triggers/anti_triggers > selected_files/fragments > context_hints） | M0 | 不可变 | 不允许在第二领域线调整基线顺序 |
| partial / handoff / failed 状态语义 | M0/M3 | 不可变 | 仅允许补充 findings 细节，不允许改状态解释 |
| Blueprint L0/L1/L2/L3 职责划分 | M0/M1 | 不可变 | route 阶段禁止默认深读 L2/L3 |
| Blueprint Query MCP surface（resource/tool/prompt 命名） | M0/M1 | 不可变 | 不允许重命名已冻结接口面 |
| workflow package 标准模板目录形状 | M2 | 不可变 | 第二领域线必须按模板复制 |
| 固定主任务 regression 基线（golden/negative/partial/handoff） | M3 | 不可变 | 新领域需增补回归，不得替换原基线 |

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

### non-breaking 判定检查表
同时满足以下条件，才可判定 non-breaking：
1. 仅新增字段，且新增字段为可选。
2. 顶层 envelope 语义与状态枚举不变。
3. route 优先级与 bundle 最小结构不变。
4. 既有 pack/validator/tests 无需强制迁移即可通过。
5. 对外消费方可忽略新增字段而不影响行为。

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

### breaking 影响评估最小项
发起 breaking 变更前必须给出以下评估：
1. 受影响 envelope 列表与字段映射。
2. 受影响 workflow package 列表。
3. 受影响 validator 与 regression 用例列表。
4. 迁移窗口、兼容策略与回滚策略。
5. 是否需要跨里程碑解冻决定。

---

## 四、breaking change 的处理流程

1. 先写变更提案
2. 说明为什么不能用加法扩展解决
3. 评估对既有 package / validator / tests 的影响
4. 给迁移路径
5. 只有在冻结面被显式解冻后才允许落地

### 处理流程模板（建议落地到 bead 说明）
1. 提案：说明目标、现状缺陷、为何不能用加法扩展。
2. 影响分析：列出受影响对象、文件、测试、下游 agent 工作面。
3. 迁移方案：给出兼容期、迁移步骤、回滚步骤。
4. 审核决策：明确谁批准、何时解冻、解冻范围。
5. 实施与验证：更新文档 + 代码 + 回归 + snapshot/handoff。
6. 再冻结：变更完成后重新冻结并记录新基线。

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

### 第二领域线接入“红线”
以下任一触发即必须停止实施并回到 breaking 流程：
1. 需要改 envelope 顶层字段含义。
2. 需要调整 route 优先级。
3. 需要改 `main/required/execution_children/deferred` 最小结构。
4. 需要重命名 Blueprint Query MCP surface。
5. 需要改 workflow package 最小模板组成。

---

## 六、冻结的真正目的
冻结不是阻止演进。  
冻结是为了保证：

1. 第二条领域线是在 **复用骨架**
2. 而不是借着扩域 **重造系统**

## 七、M4 执行者最小动作清单
1. 在 `bd` 中确认 M4 bead 状态与依赖关系真实。
2. 产出或更新 Context Snapshot，记录冻结对象与风险。
3. 产出 Handoff，明确下一位 agent 可依赖项与禁止事项。
4. 在提交前验证回归基线不受破坏。
