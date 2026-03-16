# Requirements: the-agent-packs

**Defined:** 2026-03-16  
**Core Value:** 在不泄露全仓语义的前提下，让消费侧 agent 始终拿到目标域相关、最小但不遗漏、可验证的上下文交付结果

## v1 Requirements

Requirements for initial release. Each maps to roadmap phases.

### Parsing & Index

- [ ] **PARS-01**: Maintainer can parse `package.yaml` with standard YAML parser and fail fast on invalid fields
- [ ] **PARS-02**: Maintainer can parse Blueprint frontmatter with robust YAML semantics (not string-split)
- [ ] **INDX-01**: Maintainer can rebuild index transactionally so failed compile does not leave partial DB
- [ ] **INDX-02**: Maintainer can detect index build/report write failure with explicit error outcome

### Routing Governance

- [ ] **ROUT-01**: Maintainer can enforce candidate-space-first routing (scope/mode filter before scoring)
- [ ] **ROUT-02**: Maintainer can guarantee capability is attach-only after primary domain selection
- [ ] **ROUT-03**: Maintainer can explain why main domain/main package/capabilities were selected
- [ ] **ROUT-04**: Maintainer can return explicit error/partial when no canonical registry mapping exists (no silent fallback)

### Context Delivery Contract

- [ ] **CONT-01**: Consumer agent can receive context bundle containing required domain knowledge only
- [ ] **CONT-02**: Consumer agent can inspect include/exclude rationale for delivered context
- [ ] **CONT-03**: Maintainer can verify “minimal yet complete” delivery with repeatable checks

### Validation & Runtime Governance

- [ ] **VALD-01**: Maintainer can execute core + domain validators from registry-defined plans
- [ ] **VALD-02**: Maintainer can trace validation results to activation artifacts and handoff outputs
- [ ] **GOVR-01**: Maintainer can update runtime ledgers (assumption/decision/change/validation) for every key change

### Domain Expansion Readiness

- [ ] **DOMN-01**: Maintainer can onboard a second primary domain using existing routing and registry governance rules
- [ ] **DOMN-02**: Maintainer can prove second-domain onboarding does not break WXT sample contract

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
| PARS-01 | Phase 1 | Pending |
| PARS-02 | Phase 1 | Pending |
| INDX-01 | Phase 1 | Pending |
| INDX-02 | Phase 1 | Pending |
| ROUT-01 | Phase 2 | Pending |
| ROUT-02 | Phase 2 | Pending |
| ROUT-03 | Phase 2 | Pending |
| ROUT-04 | Phase 2 | Pending |
| CONT-01 | Phase 3 | Pending |
| CONT-02 | Phase 3 | Pending |
| CONT-03 | Phase 3 | Pending |
| VALD-01 | Phase 4 | Pending |
| VALD-02 | Phase 4 | Pending |
| GOVR-01 | Phase 4 | Pending |
| DOMN-01 | Phase 5 | Pending |
| DOMN-02 | Phase 5 | Pending |

**Coverage:**
- v1 requirements: 16 total
- Mapped to phases: 16
- Unmapped: 0 ✓

---
*Requirements defined: 2026-03-16*  
*Last updated: 2026-03-16 after initial definition*
