# 22 M1 上下文：Routing 分层化与候选集裁剪

## 文档类型
单里程碑所需上下文文档

## 作用域
定义从当前平铺式 route 升级到两段式 route 的规则。

## 一、目标流程

### 第一轮：全局 route
只允许以下节点进入候选：
- `domain-root`
- `domain-orchestrator`

输出：
- 主域
- 主域 orchestrator
- 原因

### 第二轮：域内 route
只允许以下节点进入候选：
- 当前主域下的 `workflow-entry`
- 当前主域显式允许的 `domain-scoped` 节点

输出：
- main workflow package
- required cross-cutting attaches
- may_include

### 第三轮：横线挂接
不再和主 workflow package 竞争，而是：
- 根据 `required_with`
- 根据主域 orchestrator 配置
- 根据 bounded context

附挂 `attach-only` 横线能力。

## 二、裁剪原则
1. 先裁候选集，再评分
2. 横线能力不参与第一轮主竞争
3. aliases / triggers 不能突破可见性规则
4. target_pack 仍保持最高优先级，但必须落在合法候选空间内

## 三、必须可解释
route 至少回答：
- 为什么选中这个主域
- 为什么在该主域内选中这个 workflow package
- 为什么挂上这些横线而不是让它们直接成为主包

## 四、从当前实现到目标模型的差异摘要
结合当前仓库实现，M1 至少应明确以下差异，供后续实现 bead 消费：

| 当前事实 | M1 目标 |
|---|---|
| `internal/query/query.go` 以平铺 `packNodeMap` 与 `nodeProfileMap` 做命中 | 先按候选空间裁剪，再在合法空间内评分 |
| `target_domain` 只是一层简单过滤 | 第一轮先确认主域，第二轮才进入域内 workflow 选择 |
| `required_with` 主要在主节点命中后直接带出 | 横线挂接应由 `attach-only` 规则与 orchestrator 决策共同决定 |
| 横线节点与主域节点都可能出现在同层 L1 命中空间 | 横线节点默认不参加第一轮主竞争 |

## 五、M1 输出给实现层的最小规则摘要
1. 第一轮候选只允许 `domain-root` / `domain-orchestrator`。
2. 第二轮候选只允许当前主域下的 `workflow-entry` 与显式允许的 `domain-scoped` 节点。
3. `attach-only` capability 只能在主域与主 workflow 已确定后附挂。
4. `target_pack` 即使保持最高优先级，也只能命中合法候选空间。
5. route 返回结果必须同时解释“主域选择”“主包选择”“横线挂接”三件事。
