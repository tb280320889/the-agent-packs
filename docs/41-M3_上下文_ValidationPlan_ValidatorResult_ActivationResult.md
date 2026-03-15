# 41 M3 上下文：Validation Plan、Validator Result、Activation Result

## 文档类型
单里程碑所需上下文文档

## 作用域
定义 M3 需要固定的三个核心对象，以及它们之间的组合关系。

## 本文必须回答的问题
1. Validation Plan 要长什么样。
2. Validator Result 要长什么样。
3. Activation Result 如何统一承载它们。
4. validator 的结果如何影响最终状态。
5. 什么情况下最终结果是 completed / partial / handoff / failed。

---

## 一、Validation Plan

### 角色
Validation Plan 是一次运行的验证计划声明。  
它不是执行结果，也不是 artifact。

### 最小字段
- `plan_id`
- `request_id`
- `main_pack`
- `validators`
- `artifacts_under_validation`
- `severity_policy`
- `plan_reason`

### `validators` 建议结构
```json
[
  {
    "name": "validator-core-output",
    "scope": "artifact",
    "reason": "All output artifacts must satisfy envelope completeness."
  },
  {
    "name": "validator-domain-wxt-manifest",
    "scope": "domain",
    "reason": "Manifest review must cover permission and store-facing risks."
  }
]
```

---

## 二、Validator Result

### 角色
Validator Result 表达单个 validator 的运行结果。

### 最小字段
- `validator_name`
- `status`
- `findings`
- `repair_suggestions`
- `validated_artifacts`

### `status` 建议枚举
- `passed`
- `warned`
- `failed`
- `skipped`

### `findings` 建议结构
```json
[
  {
    "severity": "warn",
    "code": "missing-browser-override-note",
    "message": "Browser-specific override note is absent.",
    "artifact_ref": "manifest-review.md"
  }
]
```

---

## 三、Activation Result

### 角色
Activation Result 是一次增强运行的统一结果容器。  
它不等于某个 artifact，而是整个 activation 的总输出。

### 最小字段
- `request_id`
- `status`
- `main_pack`
- `artifacts`
- `validation_results`
- `handoff`
- `summary`

### `status` 的计算原则

#### `completed`
- 核心 artifacts 已生成
- 关键 validator 未失败
- 不需要 handoff

#### `partial`
- 已有部分结果
- 但上下文不足、依赖缺失、或 validator 有未阻塞但重要警告
- 当前不能宣称完整结束

#### `handoff`
- 当前 pack 正常到达边界
- 必须将工作移交给下一 pack
- handoff bundle 已形成

#### `failed`
- activation request 不成立
- route 不成立
- validator 的失败阻止当前输出被接受
- 无法进入合理 handoff

---

## 四、三者的组合顺序

固定顺序应为：

1. route
2. bundle
3. pack 生成 artifacts
4. 生成 Validation Plan
5. 跑 validators
6. 收集 Validator Results
7. 组装 Activation Result

### 不要颠倒
- 不要先跑 validator，再临时补 plan
- 不要先写 Activation Result，再回填 validator 结果

---

## 五、M3 至少要实现的两个 validator

### 1. validator-core-output
职责：
- 检查 artifact 是否符合统一输出壳层
- 检查 summary / sections / status 是否完整
- 检查 handoff / partial 语义是否合理

### 2. validator-domain-wxt-manifest
职责：
- 检查 manifest-review 是否覆盖关键领域项
- 检查 permissions / host permissions / browser overrides / store-facing 风险是否出现
- 检查建议与结论是否可执行

---

## 六、M3 在这些对象上不能做什么
- 不能改 M0 冻结的 envelope 语义
- 不能让 domain validator 改写 Activation Result 的顶层意义
- 不能把 findings 写成无结构长文本
