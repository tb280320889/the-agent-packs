---
phase: 05
slug: domain-expansion-pilot
status: draft
nyquist_compliant: true
wave_0_complete: true
created: 2026-03-18
---

# Phase 05 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test |
| **Config file** | none — Go testing uses go.mod module defaults |
| **Quick run command** | `go test ./... -run "TestM5DomainExpansion.*" -count=1` |
| **Full suite command** | `go test ./... -count=1` |
| **Estimated runtime** | ~90 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./... -run "TestM5DomainExpansion.*" -count=1`
- **After every plan wave:** Run `go test ./... -count=1`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 90 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 05-01-01 | 01 | 1 | DOMN-01 | integration | `go test ./... -run "TestM2RegistryLoadsAndValidates" -count=1` | ✅ | ⬜ pending |
| 05-01-02 | 01 | 1 | DOMN-01 | integration | `go test ./... -run "TestM5DomainExpansionFeatureSwitchRollback" -count=1` | ✅ | ⬜ pending |
| 05-01-03 | 01 | 1 | DOMN-01 | integration | `go test ./... -run "TestM5DomainExpansionOnboardMonorepoRouteAndActivation" -count=1` | ✅ | ⬜ pending |
| 05-02-01 | 02 | 2 | DOMN-02 | regression | `go test ./... -run "TestM5DomainExpansionWXTNonRegression" -count=1` | ✅ | ⬜ pending |
| 05-02-02 | 02 | 2 | DOMN-01, DOMN-02 | contract-doc | `go test ./... -run "TestM5DomainExpansion.*" -count=1` | ✅ | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

Existing infrastructure covers all phase requirements.

---

## Manual-Only Verifications

All phase behaviors have automated verification.

---

## Validation Sign-Off

- [x] All tasks have `<automated>` verify or Wave 0 dependencies
- [x] Sampling continuity: no 3 consecutive tasks without automated verify
- [x] Wave 0 covers all MISSING references
- [x] No watch-mode flags
- [x] Feedback latency < 90s
- [x] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
