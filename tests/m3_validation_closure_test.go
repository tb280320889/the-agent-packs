package tests

import (
	"encoding/json"
	"os"
	"testing"

	"the-agent-packs/internal/activation"
	"the-agent-packs/internal/query"
)

func writeRequestFile(t *testing.T, payload map[string]any) string {
	t.Helper()
	tmp, err := os.CreateTemp("", "activation-m3-*.json")
	if err != nil {
		t.Fatalf("create temp file failed: %v", err)
	}
	defer tmp.Close()
	enc := json.NewEncoder(tmp)
	if err := enc.Encode(payload); err != nil {
		t.Fatalf("write request failed: %v", err)
	}
	return tmp.Name()
}

func TestM3GoldenCompleted(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	requestPath := writeRequestFile(t, map[string]any{
		"request_id":    "req-m3-golden-1",
		"task":          "review WXT manifest permissions for browser store submission",
		"target_domain": "wxt",
		"bounded_context": map[string]any{
			"selected_files":   []string{"wxt.config.ts"},
			"config_fragments": []string{"manifest.permissions", "manifest.host_permissions"},
			"host_hints":       []string{"browser-extension"},
			"browser_hints":    []string{"chrome", "firefox"},
		},
		"context_hints": []string{"store submission"},
	})
	defer os.Remove(requestPath)

	result, err := activation.Execute(db, requestPath)
	if err != nil {
		t.Fatalf("activation execute failed: %v", err)
	}
	if result.Status != "completed" {
		t.Fatalf("expected completed, got %s", result.Status)
	}
	if len(result.ValidationResults) != 1 {
		t.Fatalf("expected 1 validation envelope, got %d", len(result.ValidationResults))
	}
	if len(result.ValidationResults[0].ValidatorResults) != 2 {
		t.Fatalf("expected 2 validator results, got %d", len(result.ValidationResults[0].ValidatorResults))
	}
}

func TestM3GoldenWarnLeadsToPartialByPolicy(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	requestPath := writeRequestFile(t, map[string]any{
		"request_id":    "req-m3-golden-2",
		"task":          "review WXT manifest permissions",
		"target_domain": "wxt",
		"bounded_context": map[string]any{
			"selected_files":   []string{"wxt.config.ts"},
			"config_fragments": []string{"manifest.permissions"},
			"host_hints":       []string{"browser-extension"},
			"browser_hints":    []string{},
		},
		"context_hints": []string{},
	})
	defer os.Remove(requestPath)

	result, err := activation.Execute(db, requestPath)
	if err != nil {
		t.Fatalf("activation execute failed: %v", err)
	}
	if result.Status != "partial" {
		t.Fatalf("expected partial with warn policy, got %s", result.Status)
	}
}

func TestM3NegativeWrongTargetPackNotOverriddenByHint(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	tp := "security-permissions"
	td := "wxt"
	result, err := query.RouteQuery(db, "L1", "manifest permissions review", &tp, &td, nil, nil, []string{"wxt"}, 3)
	if err != nil {
		t.Fatalf("route query failed: %v", err)
	}
	if len(result.Candidates) == 0 || result.Candidates[0].ID != "L1.security.permissions" {
		t.Fatalf("expected target_pack to keep priority, got %+v", result.Candidates)
	}
	if result.Candidates[0].ActivationMode != "attach-only" {
		t.Fatalf("expected attach-only metadata to be preserved, got %+v", result.Candidates[0])
	}
	if result.DecisionBasis != "target_pack>canonical_exact" {
		t.Fatalf("expected target_pack decision basis, got %s", result.DecisionBasis)
	}
	if result.Candidates[0].ReasonCode != "TARGET_PACK_EXACT" {
		t.Fatalf("expected target_pack reason code, got %s", result.Candidates[0].ReasonCode)
	}
	if result.Candidates[0].RuleRef != "BR-04" {
		t.Fatalf("expected BR-04 rule ref, got %s", result.Candidates[0].RuleRef)
	}
}

func TestRouteTargetPackCanonicalUnavailableAtLevelReturnsNoPrimaryCandidate(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	tp := "security-permissions"
	result, err := query.RouteQuery(db, "L0", "manifest permissions review", &tp, nil, nil, nil, nil, 3)
	if err != nil {
		t.Fatalf("route query failed: %v", err)
	}
	if len(result.Candidates) != 0 {
		t.Fatalf("expected no candidates when canonical mapping is unavailable for level, got %+v", result.Candidates)
	}
	if result.Status != "failed" {
		t.Fatalf("expected failed status when canonical missing, got %s", result.Status)
	}
	if result.ErrorCode != "ROUTE_CANONICAL_MISSING" {
		t.Fatalf("expected ROUTE_CANONICAL_MISSING, got %s", result.ErrorCode)
	}
	if result.NextAction != "检查 registry canonical 映射" {
		t.Fatalf("unexpected next action: %s", result.NextAction)
	}
	if result.DocsRef != "" {
		t.Fatalf("expected empty docs_ref placeholder, got %q", result.DocsRef)
	}
	if result.DecisionBasis != "target_pack>canonical_missing_hard_fail" {
		t.Fatalf("unexpected decision basis for unresolved target pack: %s", result.DecisionBasis)
	}
}

func TestM3NegativeBundleNotDefaultDeepRead(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	bundle, err := query.BuildContextBundle(db, "L1.wxt.manifest", true, false, false)
	if err != nil {
		t.Fatalf("build bundle failed: %v", err)
	}
	if len(bundle.ExecutionChildren) != 0 {
		t.Fatalf("expected no execution children by default, got %d", len(bundle.ExecutionChildren))
	}
	if len(bundle.Deferred) != 0 {
		t.Fatalf("expected no deferred nodes by default, got %d", len(bundle.Deferred))
	}
}

func TestM3PartialWhenContextInsufficient(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	requestPath := writeRequestFile(t, map[string]any{
		"request_id":    "req-m3-partial-1",
		"task":          "review WXT manifest permissions for browser store submission",
		"target_domain": "wxt",
		"bounded_context": map[string]any{
			"selected_files":   []string{},
			"config_fragments": []string{},
			"host_hints":       []string{},
			"browser_hints":    []string{},
		},
		"context_hints": []string{},
	})
	defer os.Remove(requestPath)

	result, err := activation.Execute(db, requestPath)
	if err != nil {
		t.Fatalf("activation execute failed: %v", err)
	}
	if result.Status != "partial" {
		t.Fatalf("expected partial, got %s", result.Status)
	}
	if result.RouteStatus != "completed" {
		t.Fatalf("expected route status completed, got %s", result.RouteStatus)
	}
	if result.RouteErrorCode != "" {
		t.Fatalf("expected empty route error code, got %s", result.RouteErrorCode)
	}
	if result.RouteDocsRef != "" {
		t.Fatalf("expected empty route docs_ref placeholder, got %q", result.RouteDocsRef)
	}
}

func TestM3HandoffContainsCarryContext(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	requestPath := writeRequestFile(t, map[string]any{
		"request_id":    "req-m3-handoff-1",
		"task":          "handoff after manifest review for security and release",
		"target_domain": "wxt",
		"bounded_context": map[string]any{
			"selected_files":   []string{"wxt.config.ts"},
			"config_fragments": []string{"manifest.permissions", "manifest.host_permissions"},
			"host_hints":       []string{"browser-extension"},
			"browser_hints":    []string{"chrome"},
		},
		"context_hints": []string{"store submission"},
	})
	defer os.Remove(requestPath)

	result, err := activation.Execute(db, requestPath)
	if err != nil {
		t.Fatalf("activation execute failed: %v", err)
	}
	if result.Status != "handoff" {
		t.Fatalf("expected handoff, got %s", result.Status)
	}
	if result.Handoff == nil {
		t.Fatalf("expected handoff payload")
	}
	if _, ok := result.Handoff["carry_context"]; !ok {
		t.Fatalf("expected carry_context in handoff")
	}
	toPacksRaw, ok := result.Handoff["to_packs"].([]string)
	if !ok {
		t.Fatalf("expected typed to_packs slice in handoff, got %#v", result.Handoff["to_packs"])
	}
	if len(toPacksRaw) != 2 || toPacksRaw[0] != "security-permissions" || toPacksRaw[1] != "release-store-review" {
		t.Fatalf("unexpected to_packs in handoff: %+v", toPacksRaw)
	}
	plan := result.ValidationResults[0].ValidationPlan
	if len(plan.Validators) != 2 {
		t.Fatalf("expected registry-aligned validators, got %+v", plan.Validators)
	}
}

func TestM6CapabilityPackActivationProducesRegisteredHandoff(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	requestPath := writeRequestFile(t, map[string]any{
		"request_id":    "req-m6-capability-handoff-1",
		"task":          "handoff after security permissions review",
		"target_pack":   "wxt-manifest",
		"target_domain": "wxt",
		"bounded_context": map[string]any{
			"selected_files":   []string{"wxt.config.ts"},
			"config_fragments": []string{"manifest.permissions", "manifest.host_permissions"},
			"host_hints":       []string{"browser-extension"},
			"browser_hints":    []string{"chrome"},
		},
		"context_hints": []string{"store submission"},
	})
	defer os.Remove(requestPath)

	result, err := activation.Execute(db, requestPath)
	if err != nil {
		t.Fatalf("activation execute failed: %v", err)
	}
	if result.Status != "handoff" {
		t.Fatalf("expected handoff, got %s", result.Status)
	}
	if result.Handoff["from_pack"] != "wxt-manifest" {
		t.Fatalf("expected wxt-manifest handoff source, got %+v", result.Handoff)
	}
	toPacks, ok := result.Handoff["to_packs"].([]string)
	if !ok {
		t.Fatalf("expected typed to_packs slice, got %#v", result.Handoff["to_packs"])
	}
	if len(toPacks) != 2 || toPacks[0] != "security-permissions" || toPacks[1] != "release-store-review" {
		t.Fatalf("unexpected to_packs: %+v", toPacks)
	}
	carry, ok := result.Handoff["carry_context"].(map[string]any)
	if !ok {
		t.Fatalf("expected carry_context map, got %#v", result.Handoff["carry_context"])
	}
	if carry["required_artifact"] != "manifest-review.md" {
		t.Fatalf("unexpected required_artifact: %+v", carry)
	}
	checks, ok := carry["required_checks"].([]string)
	if !ok {
		t.Fatalf("expected typed required_checks slice, got %#v", carry["required_checks"])
	}
	if len(checks) != 2 || checks[0] != "security-permissions-ready" || checks[1] != "release-store-review-ready" {
		t.Fatalf("unexpected required_checks: %+v", checks)
	}
}
