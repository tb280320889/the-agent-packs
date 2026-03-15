# 32 M2 上下文：wxt-manifest Pack 规格、artifact、handoff

## 文档类型
单里程碑所需上下文文档

## 作用域
定义 `wxt-manifest` 这个首个完整 pack 的具体规格。  
这是 M2 最直接的实现依据。

## 本文必须回答的问题
1. `wxt-manifest` 的范围是什么。
2. 进入条件与退出条件是什么。
3. 主要 artifact 是什么。
4. 该 pack 依赖哪些横线。
5. 何时 handoff 给相邻 pack。

---

## 一、pack 目标

`wxt-manifest` 负责审查和生成与以下对象相关的增强输出：
- manifest 生成规则
- browser-specific overrides
- permissions
- host permissions
- CSP 与 store-facing 风险提示

它的首要目标不是直接产出最终业务配置，而是产出 **结构化评审与修复建议**。

---

## 二、进入条件

进入本 pack 至少满足以下之一：
- task 明确提到 manifest
- task 明确提到 permissions / host permissions
- task 明确提到 browser store submission
- route 命中 `L1.wxt.manifest`

### 不进入的典型情况
- task 主要是 content script 注入策略
- task 主要是 background runtime 设计
- task 主要是 browser API 适配
- task 主要是整体扩展 UX

---

## 三、依赖与横线

### 必带横线
- `security-permissions`
- `release-store-review`

### 可能带上
- `wxt-permissions`
- `wxt-store-release`

### 原则
`wxt-manifest` 可以发起与汇总，但不能吞并横线职责。

---

## 四、主 artifact

### 主 artifact
`manifest-review.md`

### 该 artifact 至少应包含
1. 当前任务摘要
2. main pack / required packs
3. manifest 相关风险摘要
4. permissions / host permissions 审查
5. browser-specific overrides 提示
6. store-facing 风险与建议
7. validator 预期项
8. 下一步建议 / 可能 handoff

### 可附带次级 artifact
- `permission-checklist.md`
- `browser-overrides-notes.md`

首期建议以一个主 artifact 为主，避免面过宽。

---

## 五、退出条件
以下条件都成立时，方可认为本 pack 工作完成：

1. manifest 核心风险已被结构化说明
2. 权限最小化建议已给出
3. browser-specific 差异已提醒
4. store-facing 风险已标出
5. validation 可执行
6. 该 pack 自己不再需要继续深挖

---

## 六、handoff 条件

### handoff 到 security-permissions
当以下情况成立时：
- 风险核心是权限最小化
- 需要更细的 host permission 审查
- 涉及 unsafe surface review

### handoff 到 release-store-review
当以下情况成立时：
- 风险核心是商店提交 / 审核约束
- 需要更细的浏览器商店检查表
- manifest 审查已完成自身职责，但需要发布前校验

---

## 七、为 M3 预留的验证面

### validator-core-output 关注
- artifact 结构是否完整
- summary 是否清楚
- handoff / partial / failed 是否合理表达

### validator-domain-wxt-manifest 关注
- 是否覆盖 permissions / host permissions
- 是否覆盖 browser-specific overrides
- 是否覆盖 store-facing 风险
- 是否给出下一步建议

---

## 八、M2 实现时必须守住的范围
- 做到“完整 pack”
- 不做到“完整 WXT 体系”
- 把边界写清楚，比多写功能更重要
