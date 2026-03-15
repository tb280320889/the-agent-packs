# Blueprint 目录说明

本目录存放 Blueprint 节点（L0/L1/L2/L3）。

固定结构：

```text
blueprint/
├─ L0/
├─ L1/
├─ L2/
└─ L3/
```

编写约束：
- 路径与 id 必须严格一致
- level 与目录一致
- summary 必须可直接进入最小上下文包
- L3 默认只在 deferred 中出现

frontmatter 规范与示例见：`blueprint/frontmatter-examples.md`
