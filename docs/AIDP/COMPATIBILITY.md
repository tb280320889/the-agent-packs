# Local Project AIDP Compatibility Policy

## 目标
定义本项目 local AIDP 的兼容范围、breaking change 判定，以及在上游 system AIDP 升级或本地语义升级时应如何处理。

## 兼容性判断面
1. 本地目录结构与稳定入口
2. 文件的语义职责
3. local AIDP 与 GSD 的协同边界
4. local AIDP 与上游 system AIDP 的兼容关系

## 稳定入口
以下路径在 local AIDP 主版本内应视为稳定：
- `README.md`
- `AGENTS.md`
- `AIDP-MANIFEST.yaml`
- `core/`
- `protocol/`
- `runtime/`
- `adapters/gsd/`

## breaking change 判定
出现以下任一情况时，视为 local AIDP 的 breaking change：
- 重命名或删除稳定入口文件
- 将 `core/`、`protocol/`、`runtime/`、`adapters/gsd/` 中某文件改为完全不同的语义职责
- 让 GSD agent 无法再通过 `AGENTS.md + local AIDP` 获得稳定业务语义注入
- 更改 local AIDP 与 GSD 的职责边界，导致 GSD 工作流需要被动改变

## 非 breaking change
- 增加新文档或新章节
- 增加新的运行态工件
- 增强局部规则说明但不改变原有语义
- 增加增强开发、迭代或治理相关能力

## 上游升级处理
当 system AIDP 升级时：
1. 先阅读上游 `VERSION.md`、`CHANGELOG.md`、`COMPATIBILITY.md`
2. 判断变化属于结构变化、治理变化、协议变化还是提示词变化
3. 比较其思想是否仍与本项目“local AIDP 注入 GSD 工作流”的模型一致
4. 只在兼容且有价值时同步到 local AIDP
5. 将本地差异记录到 `CHANGELOG.md` 与 `runtime/03-变更摘要.md`

## 本地覆盖规则
- local AIDP 可以覆盖模板中的措辞和项目实例化内容
- local AIDP 不应静默背离上游 system AIDP 的核心职责分层
- 若确需偏离，应在 `runtime/02-决策日志.md` 中记录原因
