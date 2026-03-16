# Local Project AIDP Version

Current Local Version: 1.1.0
Base System AIDP Version: 2.0.1
Release Date: 2026-03-16
Status: Active

## Version Intent
本文件用于定义 **项目实例级 local AIDP** 的版本，而不是上游 system AIDP 模板仓库的版本。

它回答两类问题：
- 当前项目实例语义基线到了哪个版本
- 当项目进入新里程碑、结构升级或语义重构时，应该如何标记本地文档基线变化

## 双层版本模型

### 1. System AIDP 版本
- 由上游 `docs/AIDP-zh-v2.0.0/AIDP-zh/` 定义
- 负责模板协议、目录职责、治理规则与通用工作方式

### 2. Local Project AIDP 版本
- 由 `docs/AIDP/` 定义
- 负责当前项目的业务语义、边界、假设、决策、运行态上下文与本地演化历史

## Local Versioning Policy
- Major：项目语义结构重组、关键职责重分配、文档职责发生 breaking 变化、删除稳定入口文件
- Minor：新增协议模块、增加运行态工件、引入新的里程碑治理能力、重要规则增强但不破坏既有入口
- Patch：措辞修正、结构补充、非 breaking 的对齐与说明增强

## 当前版本说明
- `1.0.0`：建立 local project AIDP 首版，切换默认入口到 `docs/AIDP/`
- `1.1.0`：对齐 system AIDP v2.0.1 的版本治理与增强迭代机制，明确 local AIDP 与 GSD 协同边界

## 使用规则
- 升级 local AIDP 时，先判断是 Type A / Type B / Type C 变更，再决定是否需要升级 local version
- 如果只是代码实现变化但不影响项目语义与文档职责，不必升级 local AIDP 版本
- 如果上游 system AIDP 升级，不应直接覆盖 local AIDP；必须先评估兼容性与本地 delta
