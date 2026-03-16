# Archive: 2026-03 phase 01.1

本目录用于 Phase 01.1 的“先归档后删除”治理。

## Rules

- 任何疑似冗余对象先登记到 `ARCHIVE-LIST.md`。
- 状态统一从 `archived-pending-delete` 开始。
- 删除动作必须满足 `delete_gate=human-approval`。
- 本目录仅存放归档记录与必要快照，不替代长期语义文档。

## Recovery

- 若误判可通过 `original_path` 与 `archive_path` 回迁。
- 回迁后需同步更新 `ARCHIVE-LIST.md` 状态与说明。
