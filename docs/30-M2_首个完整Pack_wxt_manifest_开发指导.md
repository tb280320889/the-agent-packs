# 30 M2 首个完整 Pack（wxt-manifest）开发指导

## 文档类型
开发里程碑开发指导文档

## 作用域
M2 负责在 M0/M1 已冻结的壳层之上，实现首个完整 workflow package：`wxt-manifest`。  
它的目标不是扩张 WXT 全线，而是证明 **单包模板 + bundle + artifact + validator 接入** 可以成立。

## 本里程碑的目标
1. 按统一 package 模板实现首个完整 pack
2. 明确 pack 的输入、输出、artifact、handoff、exit criteria
3. 把 pack 接到已有 route / bundle 骨架上
4. 为 M3 准备 validator 与主任务闭环的稳定目标物

## 必读文档
- 01
- 02
- 03
- 30
- 31
- 32

## 参考输入
- M1 交付的 Blueprint / compiler / query / entry 骨架

## 默认不读
- `40+`

## 必须产出
1. `wxt-manifest` 包目录
2. `README.md`
3. `package.yaml`
4. `skill/instructions.md`
5. `skill/routing.md`
6. `skill/fallback.md`
7. `skill/examples.md`
8. `contracts/`
9. `templates/`
10. `fixtures/`
11. `tests/`
12. 首个 artifact 模板

## 明确不产出
- 整个 WXT 领域线的所有 pack
- 大量 browser API / content-script / UX 细节
- 复杂 store release 平台适配实现

---

## 执行顺序

### 第一步：先把 pack 边界写死
先写：
- 这个 pack 做什么
- 不做什么
- 何时进入
- 何时交给相邻 pack

### 第二步：再按统一模板建包
不能先写零散文档。  
必须先按 package 模板立起来。

### 第三步：固定 artifact
先固定该 pack 的主要 artifact：
- `manifest-review.md`

可附带次级 artifact，但不宜过多。

### 第四步：接入 route / bundle
让主任务可以 route 到它，并拿到合适 bundle。

### 第五步：准备给 M3 的验证面
让 validator-core-output 与 validator-domain-wxt-manifest 有稳定输入面。

---

## 完成标准
1. `wxt-manifest` 已是一个完整 package，而不是一堆散文件
2. 进入条件、边界、handoff 条件明确
3. 主 artifact 已稳定
4. M3 可以直接对它跑验证与闭环

## handoff 给 M3 的内容
- pack 目录与 contracts
- artifact 模板
- examples / fixtures
- 预期 validator 输入面
- 已知限制与边界

## 本里程碑绝对不要做的事
- 借着做 WXT pack 扩成整条 WXT 产品线
- 动系统层母版
- 回头修改 M0/M1 已冻结 envelope
- 把 artifact 做成一个大而泛的报告集合
