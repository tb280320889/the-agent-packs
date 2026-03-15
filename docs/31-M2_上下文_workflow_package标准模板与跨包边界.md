# 31 M2 上下文：workflow package 标准模板与跨包边界

## 文档类型
单里程碑所需上下文文档

## 作用域
定义单个 workflow package 的标准内部结构、`package.yaml` 基线字段，以及跨包职责边界。  
它既服务 `wxt-manifest`，也服务后续第二、第三个 pack 的复制。

## 本文必须回答的问题
1. 单个 workflow package 的最小内部结构是什么。
2. `package.yaml` 应至少包含什么。
3. skill / MCP / validator / templates / tests 各自放哪里。
4. 跨包边界如何判断。
5. 为什么一个 package 不能退化成一段 prompt。

---

## 一、标准内部结构

每个 workflow package 首期统一使用以下结构：

```text
<package>/
├─ README.md
├─ package.yaml
├─ skill/
│  ├─ instructions.md
│  ├─ routing.md
│  ├─ fallback.md
│  └─ examples.md
├─ mcp/
│  ├─ tools/
│  ├─ resources/
│  └─ prompts/
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
│  ├─ output/
│  ├─ config/
│  ├─ docs/
│  └─ checklists/
├─ fixtures/
│  ├─ pass/
│  └─ fail/
├─ tests/
│  ├─ routing/
│  ├─ validation/
│  └─ regression/
└─ CHANGELOG.md
```

---

## 二、`package.yaml` 建议最小字段

```yaml
name: wxt-manifest
kind: workflow-package
domain: wxt
subdomain: manifest
layer: product-line
version: 0.1.0
goal: >
  Review WXT manifest generation rules, browser-specific overrides,
  permissions, host permissions, and store-facing risks.
inputs:
  - user_goal
  - project_state
  - host_environment
  - available_tools
  - extension_constraints
depends_on:
  - security-permissions
  - release-store-review
mcp:
  tools:
    - inspect_manifest_inputs
    - compare_browser_overrides
  resources:
    - wxt_manifest_guidelines
    - browser_store_constraints
  prompts:
    - manifest_review_prompt
    - repair_prompt
validators:
  - validator-domain-wxt-manifest
  - validator-core-output
handoff:
  incoming:
    - global-handoff
  outgoing:
    - security-permissions-input
    - release-store-review-input
artifacts:
  - manifest-review.md
exit_criteria:
  - permissions_reviewed
  - store_risks_flagged
  - validation_passed
```

---

## 三、skill / MCP / validator 的职责边界

### skill 负责
- 任务识别
- 路径决策
- 何时调用哪些 MCP
- 何时触发验证
- 输出结构约束
- 失败回退策略

### MCP 负责
- 暴露文档 / schema / 配置 / 项目上下文
- 暴露构建、分析、发布、查询等工具
- 暴露可复用 prompts

### validator 负责
- 判定产物是否合格
- 生成修复反馈
- 决定能否进入下一 workflow

### 结论
一个 package 不是一段 prompt，而是：
**编排 + 能力面 + 验证 + 交接 + 模板** 的最小组合体。

---

## 四、跨包边界的判断规则

只有在以下任一条件成立时，才值得拆出相邻 package：
- 规则集明显不同
- 工具集明显不同
- 验证标准明显不同
- 交付物明显不同
- 安全边界明显不同
- 宿主环境明显不同

否则不要继续下钻。

---

## 五、跨包交接的最低要求
handoff 至少应说明：
- 从哪个 pack 来
- 交给哪个 pack
- 为什么交
- 带上哪些上下文
- 下一包应从哪里继续

没有这个协议，package 树会变成孤岛集合。

---

## 六、对 `wxt-manifest` 的直接约束
- 它只处理 manifest / permissions / host permissions / browser overrides / store-facing 风险
- 它不负责整个 content-script 架构
- 它不负责整个 background runtime
- 它可以挂 security / release 横线，但不能吞掉它们的职责
