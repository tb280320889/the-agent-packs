# Activation Entry（M1）

本目录定义 M1 最薄 activation entry 的职责边界。

## 最薄职责
1. 接收 Activation Request
2. 解析 task/target_pack/target_domain/bounded_context
3. 调用 route_query
4. 构建最小 context bundle
5. 返回 activation result 或 route result

## 返回状态
- `completed`
- `partial`
- `handoff`
- `failed`

## 最薄实现边界
- 不做编排引擎
- 不做长链执行
- 不做 UI
