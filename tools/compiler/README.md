# Blueprint Compiler（M1）

本目录定义 M1 最薄 compiler 的职责边界与输入输出。

## 运行方式（Go 模式B）
- 本能力由统一二进制 `agent-pack-mcp` 提供。
- 编译命令：`go run ./cmd/agent-pack-mcp compile --root blueprint --db blueprint/index/blueprint.db --report-dir blueprint/index`
- 推荐发布后使用预编译二进制执行同样子命令，不依赖 Python/Node 运行时。

## 目标
- 解析 Blueprint Markdown 的 frontmatter
- 校验路径与 id 一致性
- 生成最小 nodes/node_meta/edges
- 写入 SQLite 最小索引
- 输出校验报告与失效引用报告

## 输入
- `blueprint/` 目录下的 Markdown 文件

## 输出
- SQLite 索引文件（位置待定）
- 校验报告（JSON）
- 失效引用报告（JSON）

## 最薄实现边界
- 不做复杂 FTS 排序
- 不做增量重建优化
- 不做图可视化
