# 30 M2 wxt-manifest 完整 Pack 开发指导

## 作用域
本里程碑只负责：
- 把抽象骨架落成首个完整 pack
- 固化 wxt-manifest 的目录、职责、模板、局部 contracts
- 让 WXT 首条领域线第一次具备真实可运行实体

本里程碑不负责：
- WXT 全线做完
- runtime / background 深建模
- 多个 pack 同时做完整实现
- 大规模 domain validators

## 本里程碑唯一目标
把 `wxt-manifest` 做成第一阶段唯一完整样板 pack。

## 必读文档
- `01-通用开发方法论与里程碑解耦规则.md`
- `31-M2_上下文_WXT领域线与Pack边界.md`
- `32-M2_上下文_wxt-manifest骨架与Artifact设计.md`

## 输入前提
- M1 已完成首批节点、route、bundle、最薄入口
- 主任务固定为 manifest / permissions / store submission 闭环
- `wxt-manifest` 为固定主 pack

## 本里程碑交付物
1. `wxt-manifest` 完整目录骨架
2. `README.md`
3. `package.yaml`
4. `skill/instructions.md`
5. `skill/routing.md`
6. `skill/fallback.md`
7. `skill/examples.md`
8. `templates/reports/manifest-review.md`
9. 必要的 pack 局部 contracts
10. 薄版 `wxt-core`、`wxt-content-script`、`security-permissions`、`release-engineering` 连接位

## 实施步骤

### Step 1：只把 wxt-manifest 做完整
其他 pack 只做“能挂接、能 handoff、能声明”的薄骨架。

### Step 2：固定 pack goal
wxt-manifest 的一句话目标必须保持稳定，不能写成多段说明。

### Step 3：固定 pack 输入
只接受 manifest-centered 的 bounded context，不把 runtime / repo 全局信息吞进来。

### Step 4：固定 artifact
优先固定：
- manifest-review.md
- （可选）permission-audit.md 的推荐逻辑
- （可选）store-release-checklist.md 的 handoff 逻辑

### Step 5：固定停止条件
一旦以下条件满足必须停止：
- manifest 范围已清楚
- 主要风险已列出
- 至少一个 artifact 可生成
- handoff 条件已明确（如需要）

### Step 6：固定 handoff
当 release 或 permissions 成为主问题时，wxt-manifest 必须交接，不继续吞并。

## 不要做的事
- 不要让 wxt-manifest 吞掉 runtime
- 不要让它变成 store submission all-in-one
- 不要让它主动扫描宿主
- 不要在本阶段把 wxt-content-script 做得比它更深
- 不要引入真实插件代码

## 完成定义
只有同时满足以下条件，本里程碑才算完成：
1. route 能稳定进入 wxt-manifest
2. pack goal 清晰且边界稳定
3. 至少一个主 artifact 模板稳定
4. pack 能在 bounded context 下工作
5. handoff 规则清楚
6. M3 可在此基础上直接加 validator 并跑主任务闭环

## 向 M3 的 handoff
M3 可直接依赖：
- wxt-manifest 的固定目录
- package.yaml 的固定字段
- artifact 模板
- handoff 条件
- 薄版横线 pack 连接位

M3 不应再改：
- wxt-manifest goal
- 首批 artifact 名称
- pack 基本边界
