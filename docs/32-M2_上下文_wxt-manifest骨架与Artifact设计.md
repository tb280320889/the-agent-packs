# 32 M2 上下文：wxt-manifest 骨架与 Artifact 设计

## 作用域
本文只回答：
1. wxt-manifest 的目录骨架
2. package.yaml 的建议内容
3. 必需 artifact 的结构
4. handoff 何时触发

## 一、目录骨架
建议固定如下：

```text
packages/domain-packs/wxt/wxt-manifest/
├─ README.md
├─ package.yaml
├─ skill/
│  ├─ instructions.md
│  ├─ routing.md
│  ├─ fallback.md
│  └─ examples.md
├─ validators/
│  ├─ rules/
│  ├─ checks/
│  └─ repair/
├─ contracts/
│  ├─ input.schema.json
│  ├─ output.schema.json
│  ├─ handoff.schema.json
│  └─ artifact.schema.json
├─ templates/
│  ├─ reports/
│  │  └─ manifest-review.md
│  └─ checklists/
│     └─ store-release-checklist.md
└─ tests/
   ├─ fixtures/
   ├─ routing/
   └─ validation/
```

## 二、每个文件的职责

### README.md
给人看：
- 这个 pack 做什么
- 不做什么
- 主 artifact 是什么
- 何时 handoff

### package.yaml
给系统看：
- goal
- inputs
- resources
- validators
- artifacts
- exit_criteria
- handoff

### skill/instructions.md
给调用 agent：
- 进入该 pack 后优先回答哪些问题
- 当前 pack 的停止条件是什么
- 什么时候必须交接

### skill/routing.md
说明：
- 什么任务应该 route 到这里
- 什么任务不应该 route 到这里
- 常见横线挂接对象是谁

### skill/fallback.md
说明：
- 上下文不足时怎么办
- 任务更像别的 pack 时怎么办
- handoff 触发条件是什么

### skill/examples.md
只给 activation 场景例子，不给业务应用样例。

## 三、wxt-manifest 的稳定 goal
建议固定为一句话：

> Review and enhance manifest-related decisions for WXT browser extension tasks under bounded context.

不要扩成一整段。

## 四、wxt-manifest 要优先回答的问题
1. 这是不是 manifest-centered task
2. 当前提供的上下文能支持哪些 manifest 结论
3. 权限风险是否显式存在
4. 兼容性假设是否显式存在
5. release-facing 风险是否显式存在
6. 当前应该输出什么 artifact
7. 是否需要 handoff

## 五、必要 Artifact 设计

### 1. manifest-review.md
必须包含这些块：
- Scope
- Provided Context
- Findings
- Risks
- Permission Notes
- Release Notes
- Recommended Next Actions
- Optional Handoff

不应包含：
- 无依据的宿主结构推测
- 通用浏览器扩展教程
- 与 task 无关的宽泛建议

### 2. permission-audit.md
第一阶段处理方式：
- 可以由 wxt-manifest 推荐
- 也可以由 security-permissions 真正产出
- 第一阶段不要求一定生成完整文档，但必须有明确推荐逻辑

### 3. store-release-checklist.md
更合理的归属：
- release-engineering 产出
- wxt-manifest 识别需求并 handoff
第一阶段允许 wxt-manifest 输出轻量 checklist 建议，但不应把 release 主问题全部吞掉。

## 六、review / audit / checklist 统一结构

### review
- Scope
- Provided Context
- Findings
- Risks
- Notes
- Recommended Next Actions
- Suggested Handoff

### audit
- Scope
- Evaluated Surface
- Findings
- Risk Notes
- Minimization / Correction Notes
- Open Questions

### checklist
- Package Scope
- Required Checks
- Optional Checks
- Known Risks
- Ready / Not Ready Notes

## 七、示例 package.yaml
```yaml
name: wxt-manifest
kind: enhancement-pack
domain: wxt
subdomain: manifest
layer: domain-pack
version: 0.1.0
goal: Review and enhance manifest-related decisions for WXT browser extension tasks under bounded context.
inputs:
  - activation_request
  - selected_files
  - config_fragments
resources:
  - L1.wxt.manifest
  - L1.security.permissions
  - L1.release.browser-store
validators:
  - validator-core-output
  - validator-domain-wxt-manifest
artifacts:
  - manifest-review.md
exit_criteria:
  - manifest 范围已清楚
  - 主要风险已列出
  - 至少一个 artifact 可生成
handoff:
  - security-permissions
  - release-engineering
recommended_cross_cutting:
  - security-permissions
  - release-engineering
```

## 八、handoff 触发条件
出现以下任一情况时，应 handoff 而不是继续扩张：
- 权限最小化已成为主问题
- store submission / packaging 已成为主问题
- 当前 pack 缺必要上下文，而别的横线更适合继续
- 当前 pack 本职职责已经完成

## 非目标
本文不定义：
- validator 结果结构
- activation result 结构
- regression 测试矩阵
这些由 M3 文档定义。
