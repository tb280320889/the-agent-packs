package tests

import (
	"encoding/json"
	"os"
	"testing"

	"the-agent-packs/internal/activation"
	"the-agent-packs/internal/query"
)

func TestM5DomainExpansionWXTNonRegression(t *testing.T) {
	t.Setenv("DOMAIN_MONOREPO_ENABLED", "true")

	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	td := "wxt"
	route, err := query.RouteQuery(
		db,
		"L1",
		"review browser extension manifest and permissions for store release",
		nil,
		&td,
		[]string{"wxt.config.ts", "manifest.json"},
		[]string{"manifest.permissions"},
		[]string{"browser-extension"},
		3,
	)
	if err != nil {
		t.Fatalf("route query failed: %v", err)
	}
	if len(route.Candidates) == 0 {
		t.Fatalf("expected route candidates")
	}
	if route.Candidates[0].Pack != "wxt-manifest" {
		t.Fatalf("expected wxt-manifest to remain primary, got %+v", route.Candidates)
	}
}

func TestM5DomainExpansionWXTBundleContractNonRegression(t *testing.T) {
	t.Setenv("DOMAIN_MONOREPO_ENABLED", "true")

	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	bundle, err := query.BuildContextBundle(db, "L1.wxt.manifest", true, false, true)
	if err != nil {
		t.Fatalf("build context bundle failed: %v", err)
	}

	foundSecurity := false
	foundRelease := false
	for _, pack := range bundle.RequiredPacks {
		if pack == "security-permissions" {
			foundSecurity = true
		}
		if pack == "release-store-review" {
			foundRelease = true
		}
	}
	if !foundSecurity || !foundRelease {
		t.Fatalf("required packs regression: %+v", bundle.RequiredPacks)
	}
	if len(bundle.IncludedDecisions) == 0 || len(bundle.ExcludedDecisions) == 0 {
		t.Fatalf("expected non-empty include/exclude decisions, got include=%d exclude=%d", len(bundle.IncludedDecisions), len(bundle.ExcludedDecisions))
	}
}

func TestM5DomainExpansionWXTActivationTraceNonRegression(t *testing.T) {
	t.Setenv("DOMAIN_MONOREPO_ENABLED", "true")

	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	request := map[string]any{
		"request_id":    "req-m5-wxt-activation",
		"task":          "review browser extension manifest and permissions for store release",
		"target_domain": "wxt",
		"bounded_context": map[string]any{
			"selected_files":   []string{"wxt.config.ts", "manifest.json"},
			"config_fragments": []string{"manifest.permissions"},
			"host_hints":       []string{"browser-extension"},
		},
	}
	tempFile, err := os.CreateTemp("", "activation-m5-wxt-*.json")
	if err != nil {
		t.Fatalf("create temp file failed: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	if err := json.NewEncoder(tempFile).Encode(request); err != nil {
		t.Fatalf("write request failed: %v", err)
	}

	result, err := activation.Execute(db, tempFile.Name())
	if err != nil {
		t.Fatalf("activation execute failed: %v", err)
	}
	if result.CurrentValidationRunID == "" {
		t.Fatalf("expected non-empty CurrentValidationRunID")
	}
	if len(result.RuntimeLedger) == 0 {
		t.Fatalf("expected runtime ledger entries")
	}
	validationCurrentFound := false
	for _, entry := range result.RuntimeLedger {
		if entry.RecordType == "validation" && entry.IsCurrent {
			validationCurrentFound = true
			break
		}
	}
	if !validationCurrentFound {
		t.Fatalf("expected current validation record in runtime ledger: %+v", result.RuntimeLedger)
	}
}

func TestM5DomainExpansionConflictDoesNotStealWXT(t *testing.T) {
	t.Setenv("DOMAIN_MONOREPO_ENABLED", "true")

	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	td := "wxt"
	route, err := query.RouteQuery(
		db,
		"L1",
		"monorepo docs mention extension manifest permissions and browser extension release",
		nil,
		&td,
		[]string{"wxt.config.ts", "manifest.json", "CONTRIBUTING.md"},
		[]string{"manifest.permissions", "monorepo governance"},
		[]string{"browser-extension", "oss governance"},
		3,
	)
	if err != nil {
		t.Fatalf("route query failed: %v", err)
	}
	if len(route.Candidates) == 0 {
		t.Fatalf("expected route candidates")
	}
	if route.Candidates[0].Pack != "wxt-manifest" {
		t.Fatalf("expected conflict not steal WXT route, got %+v", route.Candidates)
	}
}
