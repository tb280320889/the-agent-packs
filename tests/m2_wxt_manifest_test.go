package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"the-agent-packs/internal/activation"
	"the-agent-packs/internal/query"
)

func TestM2BundleHasRecommendedValidatorsAndArtifacts(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	bundle, err := query.BuildContextBundle(db, "L1.wxt.manifest", true, false, false)
	if err != nil {
		t.Fatalf("build bundle failed: %v", err)
	}
	if len(bundle.RecommendedValidators) != 2 {
		t.Fatalf("expected 2 recommended validators, got %v", bundle.RecommendedValidators)
	}
	if len(bundle.RecommendedArtifacts) != 1 || bundle.RecommendedArtifacts[0] != "manifest-review.md" {
		t.Fatalf("unexpected recommended artifacts: %v", bundle.RecommendedArtifacts)
	}
}

func TestM2ActivationCompletedCarriesValidationPayload(t *testing.T) {
	root := projectRoot(t)
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	result, err := activation.Execute(db, filepath.Join(root, "fixtures", "activation-request.sample.json"))
	if err != nil {
		t.Fatalf("activation execute failed: %v", err)
	}
	if result.Status != "completed" {
		t.Fatalf("expected completed, got %s", result.Status)
	}
	if len(result.Artifacts) == 0 {
		t.Fatalf("expected artifacts in completed result")
	}
	if len(result.ValidationResults) == 0 {
		t.Fatalf("expected validation_results in completed result")
	}
}

func TestM2ActivationHandoffCarriesCarryContext(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	request := map[string]any{
		"request_id":    "req-handoff-m2",
		"task":          "handoff to next pack after manifest review",
		"target_domain": "wxt",
		"bounded_context": map[string]any{
			"selected_files":   []string{"wxt.config.ts"},
			"config_fragments": []string{"manifest.permissions"},
		},
	}
	tempFile, err := os.CreateTemp("", "activation-m2-*.json")
	if err != nil {
		t.Fatalf("create temp file failed: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	enc := json.NewEncoder(tempFile)
	if err := enc.Encode(request); err != nil {
		t.Fatalf("write request failed: %v", err)
	}

	result, err := activation.Execute(db, tempFile.Name())
	if err != nil {
		t.Fatalf("activation execute failed: %v", err)
	}
	if result.Status != "handoff" {
		t.Fatalf("expected handoff, got %s", result.Status)
	}
	if result.Handoff == nil {
		t.Fatalf("expected handoff payload")
	}
	h := result.Handoff
	if _, ok := h["carry_context"]; !ok {
		t.Fatalf("handoff payload missing carry_context")
	}
}
