# Fixtures（M1）

本目录用于放置 M1 阶段的最小 fixtures 与 smoke cases。

当前默认运行面为 Go 模式B：
- 开发态：`go run ./cmd/agent-pack-mcp ...`
- 发布态：预编译二进制 `agent-pack-mcp ...`

## 目标
- 覆盖 route / bundle / entry 的最小闭环
- 不引入完整 pack 实现
