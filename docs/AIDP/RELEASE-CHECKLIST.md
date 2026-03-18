# Local Project AIDP Release Checklist

## 版本发布前检查
- [ ] 已确认本次变化属于 Major / Minor / Patch 中哪一类
- [ ] 已检查是否与上游 system AIDP 思想保持一致
- [ ] 已检查是否破坏 local AIDP 与 GSD 的职责边界
- [ ] 已更新 `CHANGELOG.md`
- [ ] 已更新 `VERSION.md`
- [ ] 已更新必要的 runtime 工件
- [ ] 已确认没有新增多处真相源

## 若涉及旧文档退役
- [ ] 已对照 `runtime/07-旧文档退役检查表.md`
- [ ] 已确认删除旧目录不会影响 agent 注入与 GSD 协同
