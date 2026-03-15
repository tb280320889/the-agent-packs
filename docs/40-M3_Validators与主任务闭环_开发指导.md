# 40 M3 Validators 与主任务闭环 开发指导

## 文档类型
开发里程碑开发指导文档

## 作用域
M3 负责把 **validation plan -> validators -> activation result** 与固定主任务闭环跑通。  
这一步是整个首期系统从“能组织文档与 pack”升级到“有可验证输出”的关键。

## 本里程碑的目标
1. 固定 Validation Plan
2. 实现至少两个 validator
3. 把 validator 结果写回 Activation Result
4. 跑通固定主任务闭环
5. 建立 regression 基线

## 必读文档
- 01
- 02
- 40
- 41
- 42

## 参考输入
- M2 的 `wxt-manifest` 完整 package
- M1 的 entry / route / bundle 骨架
- M0 的冻结协议

## 必须产出
1. `validator-core-output`
2. `validator-domain-wxt-manifest`
3. Validation Plan 的稳定写法
4. Activation Result 的稳定组合逻辑
5. 固定主任务的闭环用例
6. regression cases

## 明确不产出
- 第二个完整领域 pack
- 控制台 / 控制平面
- 复杂多 pack 协同编排

---

## 执行顺序

### 第一步：先固定 Validation Plan 与 Activation Result 的接线关系
不要先写一堆 validator 规则，再思考结果怎么回收。

### 第二步：实现 core validator
先保证任何 artifact 至少有统一输出质量门槛。

### 第三步：实现 domain validator
再保证 `wxt-manifest` 有领域特定质量门槛。

### 第四步：跑固定主任务
选一个固定任务，把：
- request
- route
- bundle
- artifact
- validators
- activation result
跑通。

### 第五步：建立 regression
把成功、负例、partial、handoff 都变成固定样例。

---

## 完成标准
1. validation plan 已成为正式一等对象
2. 至少 2 个 validator 可运行
3. activation result 能统一回收验证结果
4. 固定主任务闭环已跑通
5. regression 基线已建立

## handoff 给 M4 的内容
- 稳定 validator 列表
- activation result 写法
- regression 样例
- 当前冻结点
- 仍可允许扩展的面

## M3 冻结对象清单（交付后不得随意改形）
1. `Validation Plan` 为运行前声明对象，固定最小字段：
   - `plan_id`
   - `request_id`
   - `main_pack`
   - `validators`
   - `artifacts_under_validation`
   - `severity_policy`
   - `plan_reason`
2. `Validator Result` 为单 validator 结构化输出，固定最小字段：
   - `validator_name`
   - `status`（`passed` / `warned` / `failed` / `skipped`）
   - `findings`
   - `repair_suggestions`
   - `validated_artifacts`
3. `Activation Result` 顶层结果容器固定字段：
   - `request_id`
   - `status`（`completed` / `partial` / `handoff` / `failed`）
   - `main_pack`
   - `artifacts`
   - `validation_results`
   - `handoff`
   - `summary`
4. M3 固定 validator 列表：
   - `validator-core-output`
   - `validator-domain-wxt-manifest`
5. 状态裁决优先级冻结为：
   - `failed > handoff > partial > completed`
6. 固定主任务冻结为：
   - `review WXT manifest permissions for browser store submission`

## 本里程碑绝对不要做的事
- 重改 M0/M1/M2 已冻结 envelope
- 先造很多 validator，再补 validation plan
- 用人工口头判断替代结构化 validator 输出
