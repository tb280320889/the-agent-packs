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
}
