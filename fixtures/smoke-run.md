# Smoke Case（M1）

本文件描述最薄链路的本地验证步骤（非自动化）。

## 1. 生成索引
```bash
python tools/compiler/compiler.py --root blueprint --db blueprint/index/blueprint.db --report-dir blueprint/index
```

## 2. 路由（L1）
```bash
python tools/query-mcp/query_mcp.py route_query --db blueprint/index/blueprint.db --level L1 --task "review WXT manifest permissions for browser store submission" --target-domain wxt
```

## 3. 构建最小 bundle
```bash
python tools/query-mcp/query_mcp.py build_context_bundle --db blueprint/index/blueprint.db --node-id L1.wxt.manifest --include-required
```

## 4. 触发 activation entry
```bash
python tools/activation-entry/activation_entry.py --db blueprint/index/blueprint.db --request fixtures/activation-request.sample.json
```
