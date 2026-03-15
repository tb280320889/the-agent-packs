# 41 M3 上下文：Validation Plan、Activation Result、Validator Result

## 作用域
本文只回答：
1. Validation Plan 如何组成
2. Activation Result 如何组装
3. Validator Result 如何表示
4. 哪些验证结果会改变 activation 状态

## 一、Validator Result

### 最小字段
```json
{
  "validator": "validator-domain-wxt-manifest",
  "status": "pass",
  "issues": [],
  "warnings": [],
  "summary": "Manifest 输出已覆盖权限、兼容性与 release-facing 要点。"
}
```

### 状态枚举
- `pass`
- `warn`
- `fail`

### 语义
- `pass`：满足该 validator 的最低要求
- `warn`：基本成立，但存在信息缺口、边界说明不足或建议修正
- `fail`：关键结构或关键领域结论缺失，当前结果不能声明完整完成

## 二、Validation Plan 生成规则

### 来源
Validation Plan = `package.yaml` + `blueprint` + `route/bundle context`

### 必含对象
- `required_artifacts`
- `required_validators`
- `gating_rules`
- `boundary_assertions`

### 建议结构
```json
{
  "plan_id": "vp-001",
  "request_id": "req-001",
  "main_pack": "wxt-manifest",
  "required_artifacts": ["manifest-review.md"],
  "required_validators": [
    "validator-core-output",
    "validator-domain-wxt-manifest"
  ],
  "gating_rules": [
    "缺 manifest-review.md 则不能 completed",
    "任一 required validator fail 则不能 completed"
  ],
  "boundary_assertions": [
    "不得默认推断未提供的宿主结构",
    "release 主问题必须 handoff 给 release-engineering"
  ]
}
```

## 三、Activation Result 组装规则

### 核心原则
Activation Result 不是“打包所有输出”，而是：
- 对本次 activation 的最终可消费说明
- 对 artifacts / validations / handoff 的统一封装
- 对当前状态的最终约束表达

### 最小结构
```json
{
  "request_id": "req-001",
  "status": "partial",
  "main_pack": "wxt-manifest",
  "route": {
    "main_blueprint_node": "L1.wxt.manifest",
    "required_packs": ["security-permissions", "release-engineering"],
    "route_reason": "manifest-centered task"
  },
  "artifacts": [
    {
      "type": "analysis_report",
      "name": "manifest-review.md",
      "format": "markdown",
      "producer_pack": "wxt-manifest",
      "status": "generated",
      "summary": "Manifest 风险与建议已输出"
    }
  ],
  "validation_results": [],
  "handoff": null,
  "summary": "已完成 manifest 主分析，但发布上下文不足，无法声明 completed。",
  "missing_context": ["browser targets", "submission constraints"],
  "recommended_next_request": {
    "target_pack": "release-engineering",
    "desired_outputs": ["checklist"]
  }
}
```

## 四、Activation 状态判定表

### completed
只有同时满足以下条件才允许：
- required artifacts 已出现
- required validators 全部 `pass` 或只有可接受 `warn`
- 没有关键 missing_context 阻断结论
- 当前主问题仍属于当前 pack
- 不存在必须 handoff 但未 handoff 的情况

### partial
适用于：
- 已产出部分 artifact
- validator 有 `warn` 或某些关键上下文缺失
- 当前 pack 仍是主问题，但不能声明完整完成

### handoff
适用于：
- 当前 pack 已到边界
- 下一个 pack 更适合继续
- handoff bundle 已具备 reason / confirmed_assumptions / open_questions

### failed
适用于：
- request 不合法
- 当前任务与当前 pack 不匹配
- required validator fail 且无法形成有意义的 partial
- 缺少最低必要输入，连 L0 级建议都无法给出

## 五、validator 与 activation 状态的关系

### 强制规则
- required validator 只要 `fail`，结果就不能是 `completed`
- required artifact 缺失，结果就不能是 `completed`
- 必须 handoff 却未生成 handoff，结果就不能是 `completed`
- 所有 validator 通过但边界说明缺失，结果仍可降为 `partial`

## 六、首批两个 validator 的最小职责

### validator-core-output
检查：
- activation result 顶层结构
- status 是否合法
- artifacts descriptor 是否齐全
- summary 是否存在

### validator-domain-wxt-manifest
检查：
- 是否明确表述当前任务确实是 manifest-centered
- 是否提到 permission implications 或诚实说明没有权限信息
- 是否提到 compatibility assumptions 或诚实说明无法判断
- 是否提到 release-facing implications 或说明不适用
- 是否给出至少一个 next action 或 handoff

## 七、本文的作用
完成本文后，M3 agent 可以把“看起来像结果”变成“真正带有验证约束的结果”。

## 非目标
本文不定义：
- route 与 bundle 的实现细节
- wxt-manifest 的目录与模板细节
这些已在前面里程碑固定。
