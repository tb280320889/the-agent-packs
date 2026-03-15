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
