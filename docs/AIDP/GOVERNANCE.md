# Local Project AIDP Governance

## 文档目标
规定 local project AIDP 的稳定层级、评审重点和版本治理原则。

## 稳定等级

### Class A - 高稳定
这些文件直接定义项目长期语义和 agent 注入入口，应谨慎变更：
- `README.md`
- `AGENTS.md`
- `AIDP-MANIFEST.yaml`
- `core/`
- `protocol/`
- `adapters/gsd/00-协同原则.md`

### Class B - 中稳定
这些文件是运行态与本地治理文件，会持续演化，但结构应保持可识别：
- `runtime/`
- `VERSION.md`
- `CHANGELOG.md`
- `COMPATIBILITY.md`
- `RELEASE-CHECKLIST.md`

### Class C - 快速演化
这些文件可更频繁迭代以改善使用体验：
- 补充说明
- 示例性质文档
- 退役检查与临时治理说明

## 评审重点
- 是否破坏了 local AIDP 作为业务语义真实源的角色
- 是否错误把 GSD 任务流逻辑写进 `core/`
- 是否把临时操作决策固化为稳定协议
- 是否引入不必要的重复描述与多处真相源

## 上游与下游关系
- 上游 system AIDP 负责模板协议与通用治理
- 下游 local AIDP 负责项目实例真相
- 不要把上游模板演化与当前项目状态混为一谈

## 退役政策
- 旧文档不应直接删除，必须先通过 `runtime/07-旧文档退役检查表.md`
- 任何稳定入口退役前，必须提供替代路径并更新 `CHANGELOG.md`
