# 50 M4 Phase 1 冻结与扩域准入 开发指导

## 作用域
本里程碑只负责：
- 判断 Phase 1 是否已经足够稳定
- 冻结对象与修改纪律
- 形成第二领域线准入门槛
- 明确后续是“补 WXT 周边”还是“开第二领域线”

本里程碑不负责：
- 直接开第二领域线
- 大改 schema
- 重构基础骨架
- 回到 M0 重新讨论设计哲学

## 本里程碑唯一目标
把第一阶段从“已跑通”推进到“已冻结、可复用、可作为后续领域模板”。

## 必读文档
- `01-通用开发方法论与里程碑解耦规则.md`
- `51-M4_上下文_修改纪律_冻结策略_BreakingChange.md`
- `52-M4_上下文_第二领域线准入与可复用骨架.md`

## 输入前提
- M3 已交付主任务闭环
- 主 pack / 主 validators / 主 artifact 已稳定
- regression matrix 已存在

## 本里程碑交付物
1. Phase 1 冻结清单
2. breaking change 执行规则
3. 允许修改 / 禁止修改对象表
4. 第二领域线准入判定
5. 复制式扩域模板
6. WXT 后续补齐优先级建议

## 实施步骤

### Step 1：确认已冻结对象
至少确认：
- 主 pack：wxt-manifest
- 主任务
- 主 validators
- activation / bundle / artifact / handoff 顶层对象
- activation result / validation plan 正式定义
- frontmatter / package.yaml 必填字段
- artifact 类型表

### Step 2：确认可修改对象
例如：
- 局部模板内容
- pack 内 checks 细节
- 非破坏性例子与说明

### Step 3：确认 breaking change 流程
明确：
- 谁能提改动
- 何时需要版本号提升
- 何时必须先证明跨领域复用价值

### Step 4：决定下一步优先级
优先判断：
1. 是否先补 WXT 周边 packs
2. 是否已满足第二领域线准入条件

### Step 5：输出复用骨架
必须形成：
- 新领域线复制清单
- 可替换项
- 禁止重写项

## 不要做的事
- 不要在 M4 才开始重写 M1 / M2 基础对象
- 不要因为单个 pack 不方便就放宽全局规则
- 不要未冻结就同时开多个领域
- 不要把“可复用”理解成复制 WXT 术语

## 完成定义
只有同时满足以下条件，本里程碑才算完成：
1. Phase 1 冻结对象明确
2. breaking change 规则可执行
3. 第二领域线准入条件明确
4. 已能说明下一步是“补 WXT 周边”还是“开第二领域线”
5. 后续扩展可在不重写骨架的前提下进行

## 对后续的 handoff
后续阶段只能在两条路径里选一条：
- **Path A：补 WXT 周边 packs**
- **Path B：进入第二领域线**

无论选哪条，都不得重写：
- activation model
- bundle model
- artifact model
- handoff model
- validator 分层
