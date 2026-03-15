# 31 M2 上下文：WXT 领域线与 Pack 边界

## 作用域
本文只回答：
1. WXT 第一条领域线的固定定位
2. 首批 packs 各自负责什么
3. 各 pack 不负责什么
4. 第一阶段固定主任务为何选 manifest

## 一、WXT 的固定定位
WXT 在本项目中是：
- 第一条 domain pack line
- 第一条参考实现骨架
- 用来验证可复用方法的试点领域

WXT 不是：
- 浏览器插件模板
- 业务插件工程
- starter repo
- 宿主必须采用的项目结构

## 二、第一阶段固定 pack 集合

### 领域 packs
- `wxt-core`
- `wxt-manifest`
- `wxt-content-script`

### 横线 packs
- `security-permissions`
- `release-engineering`

第一阶段只允许 `wxt-manifest` 做成完整样板。  
其他 pack 只做到“能 route / 能 bundle / 能 handoff / 有声明”。

## 三、各 pack 边界

### wxt-core
负责：
- 判断任务是否属于 WXT
- 初步判断更接近 manifest / content-script / runtime / release
- 补挂必要横线

不负责：
- 深入做 manifest 分析
- 深入做 content script 分析
- 生成最终 artifact

### wxt-manifest
负责：
- manifest-centered task
- manifest 结构 / 权限 / 兼容性 / 发布相关风险的增强分析
- 生成结构化工件
- 必要时 handoff 给 release 或 security

不负责：
- page injection 主问题
- runtime messaging 主问题
- 宿主仓库整体架构分析
- 默认全仓扫描

### wxt-content-script
负责：
- content script 与页面注入边界相关任务
- site compatibility 风险
- page context 假设
- 必要时 handoff 到 runtime / permissions

不负责：
- manifest 主审视
- release 主审视
- 全局扩展架构总设计

### security-permissions
负责：
- 权限范围与最小化的横向能力
- 接收 required_with / handoff
- 输出 permission audit 或权限风险说明

不负责：
- 接管 WXT 主领域分析

### release-engineering
负责：
- browser store readiness
- packaging / submission checklist
- release-facing 风险说明

不负责：
- 接管 manifest 主分析

## 四、第一阶段固定主任务
固定主任务：
`review WXT manifest permissions for browser store submission`

## 五、为什么它是最佳主任务
因为它天然覆盖：
- WXT 领域
- manifest 子问题
- permissions 横线
- release 横线
- review artifact
- validator 可写性
- handoff 可验证性

## 六、WXT 领域线的防跑偏规则
- 不允许引入真实 popup / background / content 业务代码
- 不允许让 wxt-core 变成万能包
- 不允许让 wxt-manifest 吞掉 release 与 permissions
- 不允许把 L2 节点写成教程
- 不允许把第一阶段做成“WXT 全线平台”

## 七、本文向 32 文档的接口
本文只定义“谁负责什么”。  
目录、artifact 模板、pack 内文件骨架由 `32` 文档定义。
