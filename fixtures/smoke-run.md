# Smoke Case（M1）

本文件描述最薄链路的本地验证步骤（非自动化）。

> 当前主路径为 Go 模式B。开发态使用 `go run`，发布态可直接替换为预编译二进制 `agent-pack-mcp`。

## 1. 生成索引
```bash
go run ./cmd/agent-pack-mcp compile --root blueprint --db blueprint/index/blueprint.db --report-dir blueprint/index
```

## 2. 路由（L1）
```bash
go run ./cmd/agent-pack-mcp route_query --db blueprint/index/blueprint.db --level L1 --task "review WXT manifest permissions for browser store submission" --target-domain wxt
```

## 3. 构建最小 bundle
```bash
go run ./cmd/agent-pack-mcp build_context_bundle --db blueprint/index/blueprint.db --node-id L1.wxt.manifest --include-required
```

## 4. 触发 activation entry
```bash
go run ./cmd/agent-pack-mcp activate --db blueprint/index/blueprint.db --request fixtures/activation-request.sample.json
```

## 5. 启动 MCP server（stdio）
```bash
go run ./cmd/agent-pack-mcp mcp --db blueprint/index/blueprint.db
```
