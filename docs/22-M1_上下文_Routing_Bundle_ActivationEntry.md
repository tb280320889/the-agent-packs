# 22 M1 上下文：Routing、Bundle、Activation Entry

## 文档类型
单里程碑所需上下文文档

## 作用域
定义 M1 所需的最薄 route、bundle、activation entry 行为。  
不涉及完整 pack 逻辑，只涉及骨架层如何正确工作。

## 本文必须回答的问题
1. activation entry 最薄版该做什么。
2. route 最薄版如何使用 M0 裁决规则。
3. bundle 如何保持最小化。
4. 何时返回 partial / handoff / failed。
5. M1 的 entry 为什么不能变厚。

---

## 一、Activation Entry 的最薄职责

首期 activation entry 只需要做到：
1. 接 activation request
2. 解析输入
3. 识别 `target_pack` 或 `target_domain`
4. 调 route
5. 构建最小 context bundle
6. 返回最基本 activation result 或 route result

### 不做
- 复杂 CLI 编排
- 服务编排
- 长链执行
- 复杂 UI

---

## 二、route 的最薄职责

### 输入来源
- task
- target_pack
- target_domain
- bounded_context
- context_hints
- selected_files
- config_fragments
- blueprint triggers / anti_triggers

### 输出
- `main_pack`
- `main_blueprint_node`
- `required_packs`
- `route_reason`
- `recommended_validators`
- `recommended_artifacts`

### 解释性要求
每次 route 至少说明：
- 为什么选中这个主节点
- 为什么挂上这些横线
- 为什么没有选其他候选

---

## 三、bundle 的最薄职责

### 最小结构
```json
{
  "main": {},
  "required": [],
  "execution_children": [],
  "deferred": [],
  "recommended_validators": [],
  "recommended_artifacts": []
}
```

### 最小化原则
1. 默认只给摘要
2. 不默认带 may_include
3. 只带当前步必要的 L2
4. L3 默认进 deferred

---

## 四、何时退回 L0，何时继续下钻

### 退回 L0
适用于：
- 证据不足以确定 L1
- target_domain 不清楚
- selected_files 只有模糊信号
- anti_trigger 排除了当前深路由

### 继续到 L1
适用于：
- 主域已成立
- 任务词与 triggers 明确匹配
- 已有足够 bounded context

### 继续到 L2
适用于：
- 已进入执行态
- 当前 artifact / validator 需要 L2 执行策略
- 不需要开 L3

---

## 五、partial / handoff / failed 的 entry 级行为

### partial
entry 发现方向大致对，但上下文不足时返回：
- 当前能确认的主域 / 主节点
- 仍缺哪些上下文
- 当前可给出的局部结论

### handoff
M1 不主动做业务 handoff，但必须保留 envelope 能力。  
也就是说，M1 的 entry 要支持返回 handoff shape，即使首期多半不会深用。

### failed
发生以下情况时可直接 failed：
- activation request 不成立
- 明确命中 anti_trigger 且无其他合理候选
- 任务不属于当前系统处理面

---

## 六、示例流程

### 输入
- task: `review WXT manifest permissions for browser store submission`
- target_domain: `wxt`
- bounded_context: 相关配置片段 + 选中文件

### 最薄行为
1. route 到 `L0.wxt`
2. route 到 `L1.wxt.manifest`
3. 自动挂 `L1.security.permissions`、`L1.release.store-review`
4. bundle 只带这些 L1 摘要
5. 执行态再补 `L2.wxt.manifest.permissions-review`

---

## 七、M1 不得把 entry 做成什么
- 不是控制平面
- 不是 workflow 引擎
- 不是万能 orchestrator
- 不是完整 tool runner

M1 只要做到“把骨架入口打通”，就算成功。
