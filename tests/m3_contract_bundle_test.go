package tests

import (
	"strings"
	"testing"

	"the-agent-packs/internal/model"
	"the-agent-packs/internal/query"
)

func TestContractBundleWXTPositive(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	td := "wxt"
	routeResult, err := query.RouteQuery(
		db,
		"L1",
		"review WXT manifest permissions for browser store submission",
		nil,
		&td,
		[]string{"wxt.config.ts"},
		[]string{"manifest.permissions", "manifest.host_permissions"},
		[]string{"browser-extension", "store-submission"},
		3,
	)
	if err != nil {
		t.Fatalf("route query failed: %v", err)
	}
	if len(routeResult.Candidates) == 0 {
		t.Fatalf("expected route candidates")
	}
	if routeResult.Candidates[0].ID != "L1.wxt.manifest" {
		t.Fatalf("expected L1.wxt.manifest as primary, got %+v", routeResult.Candidates)
	}
	if len(routeResult.MustInclude) != 2 || routeResult.MustInclude[0] != "L1.release.store-review" || routeResult.MustInclude[1] != "L1.security.permissions" {
		t.Fatalf("expected stable attach-only must_include, got %+v", routeResult.MustInclude)
	}

	bundle, err := query.BuildContextBundle(db, "L1.wxt.manifest", true, false, true)
	if err != nil {
		t.Fatalf("build bundle failed: %v", err)
	}
	if bundle.Main == nil {
		t.Fatalf("expected main node in context bundle")
	}
	if !strings.HasPrefix(bundle.Main.ID, "L1.wxt") {
		t.Fatalf("main node should stay in target domain, got %s", bundle.Main.ID)
	}
	allowedAttach := map[string]bool{
		"L1.security.permissions": true,
		"L1.release.store-review": true,
	}
	for _, n := range bundle.Required {
		if !strings.HasPrefix(n.ID, "L1.wxt") && !strings.HasPrefix(n.ID, "L2.wxt") && !allowedAttach[n.ID] {
			t.Fatalf("required node leaked outside target domain: %s", n.ID)
		}
	}
	if len(bundle.IncludedDecisions) == 0 {
		t.Fatalf("expected included decisions")
	}
	if len(bundle.ExcludedDecisions) == 0 {
		t.Fatalf("expected excluded decisions")
	}

	seenMainDecision := false
	seenCompletenessRelaxation := false
	for _, d := range bundle.IncludedDecisions {
		assertContractDecisionFields(t, d)
		if d.NodeID == bundle.Main.ID && d.ReasonCode == "INCLUDE_PRIMARY_CONTEXT" {
			seenMainDecision = true
		}
		if d.DecisionBasis == "completeness_over_minimality" {
			seenCompletenessRelaxation = true
		}
	}
	if !seenMainDecision {
		t.Fatalf("expected primary include decision for main node")
	}
	if !seenCompletenessRelaxation {
		t.Fatalf("expected explicit minimality relaxation record when completeness requires expansion")
	}

	for _, d := range bundle.ExcludedDecisions {
		assertContractDecisionFields(t, d)
		if strings.HasPrefix(d.NodeID, "L1.tauri") || strings.HasPrefix(d.NodeID, "L2.tauri") {
			if d.ReasonCode == "" || d.SourceRule == "" {
				t.Fatalf("cross-domain exclusion must carry rationale, got %+v", d)
			}
		}
	}
}

func assertContractDecisionFields(t *testing.T, d model.ContractDecision) {
	t.Helper()
	if d.Action != "include" && d.Action != "exclude" {
		t.Fatalf("invalid action in decision %+v", d)
	}
	if d.NodeID == "" {
		t.Fatalf("decision missing node_id: %+v", d)
	}
	if d.ReasonCode == "" || d.SourceRule == "" || d.Scope == "" || d.DecisionBasis == "" {
		t.Fatalf("machine-readable fields incomplete: %+v", d)
	}
	if strings.TrimSpace(d.HumanNote) == "" {
		t.Fatalf("human-readable note missing: %+v", d)
	}
	if !strings.HasPrefix(d.SourceRule, "BR-") && !strings.HasPrefix(d.SourceRule, "CONT-") {
		t.Fatalf("source_rule must map to stable contract rule id, got %s", d.SourceRule)
	}
	if !strings.Contains(d.Scope, "domain") && d.Scope != "attach_only_capability" {
		t.Fatalf("unexpected scope value: %s", d.Scope)
	}
	if d.ReasonCode == "INCLUDE_COMPLETENESS_RELAXATION" && d.DecisionBasis != "completeness_over_minimality" {
		t.Fatalf("completeness relaxation must record decision basis, got %+v", d)
	}
	if d.ReasonCode == "EXCLUDE_MINIMALITY_BOUNDARY" && d.DecisionBasis != "minimality_guard" {
		t.Fatalf("minimality exclusion must keep stable decision basis, got %+v", d)
	}
}
