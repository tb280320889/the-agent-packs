# Testing Patterns

**Analysis Date:** 2026-03-16

## Test Framework

**Runner:**
- Go `testing` 标准库
- 配置：未检测到独立配置文件（无 `go test` 配置文件）

**Assertion Library:**
- 标准库 `testing` + 手写断言（`t.Fatalf`）

**Run Commands:**
```bash
go test ./...              # Run all tests
go test ./... -run TestM2  # Filter by name
go test ./... -count=1     # Disable test cache
```

## Test File Organization

**Location:**
- 独立 `tests/` 目录集中管理：`tests/m1_minimal_test.go`、`tests/m3_validation_closure_test.go`

**Naming:**
- `*_test.go` 文件命名；测试函数以 `Test` 开头：`tests/m2_registry_test.go`。

**Structure:**
```
tests/
├── m1_minimal_test.go
├── m2_registry_test.go
├── m2_wxt_manifest_test.go
└── m3_validation_closure_test.go
```

## Test Structure

**Suite Organization:**
```go
func TestXxx(t *testing.T) {
    dbPath := compileMainIndex(t)
    db := openDB(t, dbPath)
    defer db.Close()
    // arrange
    // act
    // assert with t.Fatalf
}
```
示例来自 `tests/m1_minimal_test.go`、`tests/m3_validation_closure_test.go`。

**Patterns:**
- 共享测试工具函数：`tests/m1_minimal_test.go` 中的 `projectRoot`、`compileMainIndex`、`openDB`。
- 断言失败使用 `t.Fatalf` 并包含上下文：`tests/m2_wxt_manifest_test.go`。

## Mocking

**Framework:**
- 未使用 mocking 框架。

**Patterns:**
```go
// 通过临时文件/临时目录模拟输入
tempFile, _ := os.CreateTemp("", "activation-*.json")
enc := json.NewEncoder(tempFile)
_ = enc.Encode(request)
```
示例来自 `tests/m1_minimal_test.go`、`tests/m2_wxt_manifest_test.go`。

**What to Mock:**
- 通过 `os.CreateTemp`、`t.TempDir` 模拟文件系统输入：`tests/m1_minimal_test.go`。

**What NOT to Mock:**
- 核心路由/编译流程以集成方式跑真实 SQLite 与文件扫描：`tests/m1_minimal_test.go` 的 `compiler.Compile` + `query.OpenDB`。

## Fixtures and Factories

**Test Data:**
```go
request := map[string]any{
  "request_id": "req-handoff",
  "task": "handoff to next pack after manifest review",
  "bounded_context": map[string]any{...},
}
```
示例来自 `tests/m1_minimal_test.go`、`tests/m3_validation_closure_test.go`。

**Location:**
- 固定样例 JSON：`fixtures/activation-request.sample.json`（由测试读取）。

## Coverage

**Requirements:** None enforced

**View Coverage:**
```bash
go test ./... -cover
```

## Test Types

**Unit Tests:**
- 以小模块断言为主（registry/load、route 评分与筛选规则）：`tests/m2_registry_test.go`、`tests/m1_minimal_test.go`。

**Integration Tests:**
- 编译蓝图 + 打开 SQLite + 执行 activation 流程：`tests/m3_validation_closure_test.go`。

**E2E Tests:**
- 未检测到（无 Playwright/Jest/Vitest）。

## Common Patterns

**Async Testing:**
- 不涉及并发/异步测试（纯同步流程）。

**Error Testing:**
```go
if err == nil {
    t.Fatalf("expected ...")
}
```
示例来自 `tests/m2_registry_test.go`、`tests/m1_minimal_test.go`。

---

*Testing analysis: 2026-03-16*
