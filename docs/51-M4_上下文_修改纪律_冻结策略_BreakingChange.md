# 51 M4 上下文：修改纪律、冻结策略、Breaking Change

## 作用域
本文只回答：
1. 什么对象已经冻结
2. 何时允许改总纲 / schema / artifact 类型
3. 什么改动算 breaking change
4. 冻结期如何治理修改

## 一、权威层级
优先级从高到低：
1. 总纲 / 主设计基点
2. 协议规范
3. 标准规范（frontmatter / package.yaml / naming）
4. 单个 pack 文档
5. 实现细节
6. 临时说明

局部文档与上层冲突时，应优先修正局部文档，而不是轻易改上层规则。

## 二、冻结对象清单
进入 M4 后，以下对象原则上视为冻结：
- frontmatter 必填字段集合
- package.yaml 必填字段集合
- artifact 类型表
- activation result 状态枚举
- validation plan 的顶层字段
- 主任务
- 主 pack
- 主 validators
- 首批 blueprint node id

## 三、允许修改总纲的条件
只有同时满足以下条件，才允许修改总纲：
1. 当前总纲明确阻碍闭环或扩展
2. 问题无法通过局部 pack / protocol / template 调整解决
3. 修改后的规则对多个领域都成立，而不是单个领域特例

只要不满足三条中的任意一条，就不应改总纲。

## 四、允许修改 schema 的条件
必须同时满足：
1. 当前 schema 不足以表达必要对象
2. 新字段可跨领域复用
3. 不会让现有 activation / bundle / artifact / handoff 失控

否则不应改全局 schema。

## 五、允许新增 artifact 类型的条件
必须同时满足：
1. 现有类型确实无法覆盖
2. 新类型可跨多个领域复用
3. 调用 agent 能明确消费
4. validator 模型不需要大面积重写

否则不应新增。

## 六、Breaking Change 分级

### 级别 0：说明性修改
- 改写文案
- 加例子
- 加注释
无需版本提升。

### 级别 1：兼容性扩展
- 新增可选字段
- 新增 recommended_validators
- 新增非强制模板块
需要记录变更，但不强制迁移。

### 级别 2：破坏性协议变更
- 改必填字段
- 改顶层状态
- 改字段语义
- 改 artifact 类型含义
必须提升版本并提供迁移说明。

### 级别 3：冻结期禁止变更
- 主任务
- 主 pack
- 主 validators
- frontmatter 必填字段集合
- package.yaml 必填字段集合
- activation result 状态枚举

这些对象在 Phase 1 结束前禁止改。

## 七、冻结期工作原则
冻结期允许做：
- 例子更清楚
- tests 更完整
- validator checks 更稳
- 模板正文更可读

冻结期不允许做：
- 为单个 pack 特例扩协议
- 为临时需求改全局 schema
- 改 route 优先级
- 改 pack 命名体系

## 八、修改提案的最小模板
任何结构性改动都必须回答：
1. 当前问题是什么
2. 为什么不能在局部解决
3. 改动影响哪些对象
4. 是否跨领域成立
5. 是否会破坏已有 pack / validator / artifacts
6. 是否需要版本提升
7. 是否需要迁移说明

## 九、本文的目的
不是阻止演进，而是阻止：
- 规则漂移
- 单 pack 绑架全局
- 过早扩张
- 在未冻结前进入多领域并行
