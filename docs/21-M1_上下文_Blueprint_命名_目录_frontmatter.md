# 21 M1 上下文：Blueprint 命名、目录、frontmatter

## 文档类型
单里程碑所需上下文文档

## 作用域
定义 M1 需要固定的：
- Blueprint 目录结构
- 节点命名规则
- frontmatter 最小字段
- L0/L1/L2/L3 正文模板
- 首批必要节点清单

## 本文必须回答的问题
1. Blueprint 文件怎么放。
2. id、level、domain、subdomain 如何对应路径。
3. frontmatter 最小需要哪些字段。
4. 正文每层写什么，不写什么。
5. 首批节点应该写到什么程度。

---

## 一、目录结构

首期固定：
```text
blueprint/
├─ L0/
├─ L1/
├─ L2/
└─ L3/
```

### 各层职责
- `L0/`：领域入口
- `L1/`：子领域编排
- `L2/`：执行策略
- `L3/`：边缘与升级

---

## 二、id 与路径推导规则

### 规则
`id` 必须可从路径稳定推导。  
例如：

- `blueprint/L0/wxt/overview.md` -> `L0.wxt`
- `blueprint/L1/wxt/manifest.md` -> `L1.wxt.manifest`
- `blueprint/L2/wxt/manifest/permissions-review.md` -> `L2.wxt.manifest.permissions-review`

### 硬约束
1. 文件路径与 id 必须一致
2. `level` 与所在目录必须一致
3. `domain` 必须从一级目录可推导
4. `subdomain` / `capability` 不能靠自由文本猜

---

## 三、frontmatter 最小字段

每个节点都必须有 YAML frontmatter。  
首期最小字段如下：

```yaml
id: L1.wxt.manifest
level: L1
domain: wxt
subdomain: manifest
capability: null
title: WXT Manifest
summary: Review manifest generation rules, browser overrides, permissions and store-facing risks.
aliases:
  - web extension manifest
  - wxt config
triggers:
  - manifest
  - permissions
  - host permissions
anti_triggers:
  - tauri
  - telegram miniapp
required_with:
  - L1.security.permissions
  - L1.release.store-review
may_include:
  - L2.wxt.manifest.permissions-review
children:
  - L2.wxt.manifest.permissions-review
entry_conditions:
  - browser_extension_host_confirmed
stop_conditions:
  - execution_plan_possible
```

### 说明
- `aliases / triggers / anti_triggers / required_with / may_include / children` 允许为空数组，但字段语义必须稳定。
- `summary` 必须是可直接进入 bundle 的摘要，不能写成长段正文。
- `entry_conditions / stop_conditions` 首期可以很薄，但字段必须存在。

---

## 四、frontmatter 非法示例

### 非法示例 A：id 与路径不一致
```yaml
id: L1.browser.manifest
level: L1
domain: wxt
subdomain: manifest
```

### 非法示例 B：字段语义漂移
```yaml
required_with:
  - maybe security
```

### 非法示例 C：把正文塞进 summary
`summary` 不能变成长段实现说明。

---

## 五、正文模板

### L0 正文
只写：
- 领域定义
- 何时进入
- 何时不要进入
- 必带横线
- 推荐下一跳

### L1 正文
只写：
- 子领域目标
- 输入输出
- 依赖
- 常见分支
- 升级条件
- 停止条件

### L2 正文
只写：
- 执行步骤
- 配置点
- 工具调用点
- validator 前置

### L3 正文
只写：
- edge cases
- 平台差异
- 深度排障
- 回退路径

---

## 六、首批必要节点清单

首期只落首个闭环所需最小节点集。建议至少有：

### L0
- `L0.wxt`
- `L0.security`
- `L0.release`

### L1
- `L1.wxt.manifest`
- `L1.security.permissions`
- `L1.release.store-review`

### L2
- `L2.wxt.manifest.permissions-review`
- `L2.wxt.manifest.browser-overrides`
- `L2.security.permissions.minimization`
- `L2.release.store-review.browser-extension-checklist`

### L3
- `L3.wxt.manifest.edge-cases`
- `L3.release.store-review.browser-specific-edge-cases`

> 这套节点不是长期领域树的全部，只是首个闭环所需最小闭包。

---

## 七、M1 在这份文档上必须守住的底线
1. 首批节点要少
2. frontmatter 语义要硬
3. summary 要能直接进 bundle
4. 路径与 id 绝对不能漂
5. 不因为未来可能有用就先把整棵树写出来
