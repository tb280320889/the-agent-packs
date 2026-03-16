# 维护 local Project AIDP

## 改动分类
修改 local AIDP 前，先判断属于哪一类：
- 项目语义变化
- 协议/agent 行为变化
- 运行态工件变化
- GSD 协同映射变化
- 版本治理变化
- 仅说明性修正

## 分层规则
- `core/` 只放项目长期稳定语义
- `protocol/` 只放 agent 行为与文档更新规则
- `runtime/` 只放当前状态、决策、变更、验证和迭代记录
- `adapters/gsd/` 只放与 GSD 的协同和映射，不替代 GSD 自身工作流文档

## 修改前必须先做的事
1. 读取当前 local AIDP 基线
2. 判断是否属于 Type A / Type B / Type C 变更
3. 明确受影响文档
4. 判断是否需要升级 local AIDP 版本

## 修改后必须做的事
1. 更新受影响文档
2. 更新 `runtime/03-变更摘要.md`
3. 必要时更新 `runtime/02-决策日志.md`
4. 如属于有记录意义的版本变化，更新 `CHANGELOG.md` 与 `VERSION.md`

## 质量标准
每次新增文档至少要回答：
- 它解决什么问题
- 何时使用
- 不该拿来做什么
- 它与 GSD / local AIDP / 代码事实的关系是什么
