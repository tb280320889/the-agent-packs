package tests

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"the-agent-packs/internal/activation"
	"the-agent-packs/internal/model"
)

func TestM4ValidationTraceContractStructs(t *testing.T) {
	envelopeType := reflect.TypeOf(model.ValidationEnvelope{})
	requiredEnvelopeFields := []string{"RunID", "EvidenceRefs", "MachineView", "HumanView"}
	for _, field := range requiredEnvelopeFields {
		if _, ok := envelopeType.FieldByName(field); !ok {
			t.Fatalf("ValidationEnvelope missing field %s", field)
		}
	}

	resultType := reflect.TypeOf(model.ActivationResult{})
	requiredResultFields := []string{"ValidationRunHistory", "CurrentValidationRunID"}
	for _, field := range requiredResultFields {
		if _, ok := resultType.FieldByName(field); !ok {
			t.Fatalf("ActivationResult missing field %s", field)
		}
	}

	result := model.ActivationResult{CurrentValidationRunID: "test-run"}
	result.ValidationRunHistory = []model.ValidationEnvelope{{
		RunID:        "test-run",
		EvidenceRefs: []model.ValidationEvidenceRef{{RefID: "runtime-ledger:test-run"}},
	}}

	if result.CurrentValidationRunID == "" {
		t.Fatalf("CurrentValidationRunID should not be empty")
	}
	if len(result.ValidationRunHistory) == 0 || len(result.ValidationRunHistory[0].EvidenceRefs) == 0 {
		t.Fatalf("ValidationRunHistory should carry EvidenceRefs")
	}
}

func TestM4ValidationTraceLinksArtifactsAndHandoff(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	requestPath := writeRequestFile(t, map[string]any{
		"request_id":                "req-m4-trace-1",
		"task":                      "handoff after manifest review for security and release",
		"target_domain":             "wxt",
		"phase_id":                  "04",
		"plan_id":                   "01",
		"validation_trigger_kind":   "milestone_auto",
		"validation_trigger_reason": "plan_milestone_validation",
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

	if result.CurrentValidationRunID == "" {
		t.Fatalf("CurrentValidationRunID should not be empty")
	}
	if len(result.ValidationRunHistory) == 0 {
		t.Fatalf("ValidationRunHistory should not be empty")
	}

	envelope := result.ValidationRunHistory[0]
	if len(envelope.EvidenceRefs) == 0 {
		t.Fatalf("EvidenceRefs should not be empty")
	}

	hasArtifact := false
	hasHandoff := false
	hasRuntimeLedger := false
	for _, ref := range envelope.EvidenceRefs {
		if strings.HasPrefix(ref.RefID, "artifact:") {
			hasArtifact = true
		}
		if strings.HasPrefix(ref.RefID, "handoff:") {
			hasHandoff = true
		}
		if strings.HasPrefix(ref.RefID, "runtime-ledger:") {
			hasRuntimeLedger = true
		}
	}

	if !hasArtifact {
		t.Fatalf("expected artifact evidence ref, got %+v", envelope.EvidenceRefs)
	}
	if !hasHandoff {
		t.Fatalf("expected handoff evidence ref, got %+v", envelope.EvidenceRefs)
	}
	if !hasRuntimeLedger {
		t.Fatalf("expected runtime-ledger evidence ref, got %+v", envelope.EvidenceRefs)
	}

	if envelope.MachineView.Status != result.Status {
		t.Fatalf("machine status should match envelope status: machine=%s activation=%s", envelope.MachineView.Status, result.Status)
	}
	if strings.TrimSpace(envelope.HumanView.Summary) == "" {
		t.Fatalf("human summary should not be empty")
	}
}
