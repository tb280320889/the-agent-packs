# Requirements: the-agent-packs

**Defined:** 2026-03-16  
**Core Value:** 在不泄露全仓语义的前提下，让消费侧 agent 始终拿到目标域相关、最小但不遗漏、可验证的上下文交付结果

## v1 Requirements

Requirements for initial release. Each maps to roadmap phases.

### Parsing & Index

- [x] **PARS-01**: Maintainer can parse `package.yaml` with standard YAML parser and fail fast on invalid fields
- [x] **PARS-02**: Maintainer can parse Blueprint frontmatter with robust YAML semantics (not string-split)
- [x] **INDX-01**: Maintainer can rebuild index transactionally so failed compile does not leave partial DB
- [x] **INDX-02**: Maintainer can detect index build/report write failure with explicit error outcome

### Routing Governance

- [x] **ROUT-01**: Maintainer can enforce candidate-space-first routing (scope/mode filter before scoring)
- [x] **ROUT-02**: Maintainer can guarantee capability is attach-only after primary domain selection
- [x] **ROUT-03**: Maintainer can explain why main domain/main package/capabilities were selected
- [x] **ROUT-04**: Maintainer can return explicit error/partial when no canonical registry mapping exists (no silent fallback)

### Context Delivery Contract

- [x] **CONT-01**: Consumer agent can receive context bundle containing required domain knowledge only
- [x] **CONT-02**: Consumer agent can inspect include/exclude rationale for delivered context
- [x] **CONT-03**: Maintainer can verify “minimal yet complete” delivery with repeatable checks

### Validation & Runtime Governance

- [x] **VALD-01**: Maintainer can execute core + domain validators from registry-defined plans
- [x] **VALD-02**: Maintainer can trace validation results to activation artifacts and handoff outputs
- [x] **GOVR-01**: Maintainer can update runtime ledgers (assumption/decision/change/validation) for every key change

### Domain Expansion Readiness

- [x] **DOMN-01**: Maintainer can onboard a second primary domain using existing routing and registry governance rules
- [x] **DOMN-02**: Maintainer can prove second-domain onboarding does not break WXT sample contract

## v2 Requirements

Deferred to future release. Tracked but not in current roadmap.

### Platformization

- **PLAT-01**: Maintainer can register validators dynamically without core code edits
- **PLAT-02**: Maintainer can apply index schema versioning and migration workflows safely
- **PLAT-03**: Maintainer can support configurable domain inference without hardcoded trigger logic
- **PLAT-04**: Maintainer can provide stronger observability metrics for route quality and context precision

## Out of Scope

Explicitly excluded. Documented to prevent scope creep.

| Feature | Reason |
|---------|--------|
| Full runtime rewrite in non-Go stack | High migration risk, no direct milestone value |
| Consumer-side default access to entire AIDP corpus | Violates progressive disclosure contract |
| Capability-first global routing | Breaks domain boundary and explainability guarantees |
| Multi-domain mass rollout in one iteration | Too risky before foundation hardening |

## Traceability

Which phases cover which requirements. Updated during roadmap creation.

| Requirement | Phase | Status |
|-------------|-------|--------|
| PARS-01 | Phase 1 | Complete |
| PARS-02 | Phase 1 | Complete |
| INDX-01 | Phase 1 | Complete |
| INDX-02 | Phase 1 | Complete |
| ROUT-01 | Phase 2 | Complete |
| ROUT-02 | Phase 2 | Complete |
| ROUT-03 | Phase 2 | Complete |
| ROUT-04 | Phase 2 | Complete |
| CONT-01 | Phase 3 | Complete |
| CONT-02 | Phase 3 | Complete |
| CONT-03 | Phase 3 | Complete |
| VALD-01 | Phase 4 | Complete |
| VALD-02 | Phase 4 | Complete |
| GOVR-01 | Phase 4 | Complete |
| DOMN-01 | Phase 5 | Complete |
| DOMN-02 | Phase 5 | Complete |

**Coverage:**
- v1 requirements: 16 total
- Mapped to phases: 16
- Unmapped: 0 ✓

---
*Requirements defined: 2026-03-16*  
*Last updated: 2026-03-16 after initial definition*
