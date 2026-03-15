# 40 M3 Validators 与主任务闭环 开发指导

## 作用域
本里程碑只负责：
- 让输出可验证
- 让 activation result 与 validation plan 真正进入运行闭环
- 跑通固定主任务
- 补齐回归测试矩阵

本里程碑不负责：
- 全量 validator 体系
- 第二领域线
- 大量新 pack
- 自动 patch 生成

## 本里程碑唯一目标
把第一阶段从“能运行”推进到“能验证、能验收、能确认闭环成立”。

## 必读文档
- `01-通用开发方法论与里程碑解耦规则.md`
- `41-M3_上下文_ValidationPlan_ActivationResult_ValidatorResult.md`
- `42-M3_上下文_测试矩阵与闭环验收.md`

## 输入前提
- M2 已交付完整 wxt-manifest
- 主 artifact 模板已稳定
- route / bundle / handoff 行为已稳定

## 本里程碑交付物
1. `validator-core-output`
2. `validator-domain-wxt-manifest`
3. Validation Plan 生成逻辑
4. Activation Result 组装逻辑
5. Validator Result 结构
6. 固定主任务的 regression fixtures
7. 闭环验收记录

## 实施步骤

### Step 1：先做 core validator
它只解决一个问题：
- 结果是不是结构化对象
- 该有的 descriptor / summary / status 是否存在

### Step 2：再做 domain validator
它只解决一个问题：
- manifest 类输出是否有依据、有边界、有落点
而不是写一个复杂审计框架。

### Step 3：把 Validation Plan 接入运行链
plan 不应事后拼接，而应由：
- pack 声明
- blueprint 推荐
- 当前 route / bundle
共同生成。

### Step 4：组装 Activation Result
必须根据 validator 状态决定：
- completed
- partial
- handoff
- failed

### Step 5：建立测试矩阵
至少覆盖：
- 主成功路径
- 上下文不足
- 错路由防御
- artifact 缺失
- validator fail
- 正常 handoff

### Step 6：跑主任务闭环
固定主任务：
`review WXT manifest permissions for browser store submission`

必须完整跑通：
activation -> route -> bundle -> artifact -> validation -> activation result

## 不要做的事
- 不要把 validator 数量做大
- 不要先做 permissions / runtime / release 一堆 domain validators
- 不要把 domain validator 变成自由文本批注
- 不要在 validator 阶段又改 pack 边界
- 不要用“人工判断差不多”代替结构化状态

## 完成定义
只有同时满足以下条件，本里程碑才算完成：
1. core + domain 两个 validator 稳定运行
2. Validation Plan 可以生成
3. Activation Result 状态由验证结果约束
4. 主任务闭环跑通
5. 至少有一组 golden cases 和一组负例
6. M4 可以基于这些结果决定是否冻结 Phase 1

## 向 M4 的 handoff
M4 可直接依赖：
- 已稳定 validators
- 已稳定 validation plan
- 已稳定 activation result
- 已跑通主任务闭环
- 已建立 regression matrix

M4 不应再改：
- 主任务
- 主 pack
- 主 validators
