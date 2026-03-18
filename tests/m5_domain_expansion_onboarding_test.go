package tests

import (
	"encoding/json"
	"os"
	"testing"

	"the-agent-packs/internal/activation"
	"the-agent-packs/internal/query"
)

func TestM5DomainExpansionOnboardMonorepoRouteAndActivation(t *testing.T) {
	t.Setenv("DOMAIN_MONOREPO_ENABLED", "true")

	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	td := "monorepo"
	route, err := query.RouteQuery(
		db,
		"L1",
		"need oss governance policy and monorepo contribution standards",
		nil,
		&td,
		[]string{"CODEOWNERS", "CONTRIBUTING.md"},
		[]string{"monorepo governance"},
		[]string{"oss compliance"},
		3,
	)
	if err != nil {
		t.Fatalf("route query failed: %v", err)
	}
	if len(route.Candidates) == 0 {
		t.Fatalf("expected route candidates")
	}
	if route.Candidates[0].Pack != "monorepo-oss-governance" {
		t.Fatalf("expected monorepo-oss-governance as main pack, got %+v", route.Candidates)
	}
	if route.DecisionBasis == "" {
		t.Fatalf("expected non-empty decision basis")
	}

	request := map[string]any{
		"request_id":    "req-m5-onboard",
		"task":          "need oss governance policy and monorepo contribution standards",
		"target_domain": "monorepo",
		"bounded_context": map[string]any{
			"selected_files":   []string{"CODEOWNERS", "CONTRIBUTING.md"},
			"config_fragments": []string{"monorepo governance"},
			"host_hints":       []string{"github"},
		},
	}
	tempFile, err := os.CreateTemp("", "activation-m5-monorepo-*.json")
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
	if result.MainPack == nil || *result.MainPack != "monorepo-oss-governance" {
		t.Fatalf("expected monorepo-oss-governance main pack, got %+v", result.MainPack)
	}
	if result.RouteStatus != "completed" {
		t.Fatalf("expected route status completed, got %s", result.RouteStatus)
	}
}

func TestM5DomainExpansionFeatureSwitchRollback(t *testing.T) {
	t.Setenv("DOMAIN_MONOREPO_ENABLED", "false")

	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	td := "monorepo"
	route, err := query.RouteQuery(
		db,
		"L1",
		"need oss governance policy and monorepo contribution standards",
		nil,
		&td,
		[]string{"CODEOWNERS", "CONTRIBUTING.md"},
		[]string{"monorepo governance"},
		[]string{"oss compliance"},
		3,
	)
	if err != nil {
		t.Fatalf("route query failed: %v", err)
	}
	if route.ErrorCode != "ROUTE_NO_PRIMARY_CANDIDATE" {
		t.Fatalf("expected ROUTE_NO_PRIMARY_CANDIDATE, got %s", route.ErrorCode)
	}
	if len(route.Candidates) != 0 {
		t.Fatalf("expected no candidates when feature switch is disabled, got %+v", route.Candidates)
	}

	for _, decision := range route.CapabilityDecisions {
		if decision.Attached {
			t.Fatalf("expected capability not attached without primary candidate, got %+v", route.CapabilityDecisions)
		}
	}
}
