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

	if envelope.MachineView.Status != "passed" {
		t.Fatalf("expected machine status passed, got %s", envelope.MachineView.Status)
	}
	if result.Status != "completed" {
		t.Fatalf("expected activation status completed, got %s", result.Status)
	}
	if strings.TrimSpace(envelope.HumanView.Summary) == "" {
		t.Fatalf("human summary should not be empty")
	}
}

func TestM4ValidationPlanGeneratedFromRegistry(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	requestPath := writeRequestFile(t, map[string]any{
		"request_id":                "req-m4-plan-registry-1",
		"task":                      "review WXT manifest permissions for browser store submission",
		"target_domain":             "wxt",
		"phase_id":                  "04",
		"plan_id":                   "02",
		"validation_trigger_kind":   "validator_manifest_change_auto",
		"validation_trigger_reason": "validator manifest changed in registry",
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
	if len(result.ValidationResults) != 1 {
		t.Fatalf("expected 1 validation result, got %d", len(result.ValidationResults))
	}

	envelope := result.ValidationResults[0]
	if envelope.TriggerKind != "validator_manifest_change_auto" {
		t.Fatalf("expected trigger kind validator_manifest_change_auto, got %s", envelope.TriggerKind)
	}
	if envelope.MachineView.Status != "passed" {
		t.Fatalf("expected machine status passed, got %s", envelope.MachineView.Status)
	}
	if len(envelope.HumanView.NextActions) == 0 {
		t.Fatalf("expected human next actions to be present")
	}

	if len(envelope.ValidationPlan.Validators) < 2 {
		t.Fatalf("expected core+domain validators from registry, got %+v", envelope.ValidationPlan.Validators)
	}
	if envelope.ValidationPlan.Validators[0].Name != "validator-core-output" {
		t.Fatalf("expected validator-core-output first, got %+v", envelope.ValidationPlan.Validators)
	}
	hasDomainPlan := false
	for _, item := range envelope.ValidationPlan.Validators {
		if strings.HasPrefix(item.Name, "validator-domain-") {
			hasDomainPlan = true
			break
		}
	}
	if !hasDomainPlan {
		t.Fatalf("expected at least one domain validator in plan, got %+v", envelope.ValidationPlan.Validators)
	}
	if !strings.Contains(envelope.ValidationPlan.PlanReason, "registry-defined plan") {
		t.Fatalf("expected plan reason mention registry-defined plan, got %s", envelope.ValidationPlan.PlanReason)
	}

	hasCoreResult := false
	hasDomainResult := false
	for _, item := range envelope.ValidatorResults {
		if item.ValidatorName == "validator-core-output" {
			hasCoreResult = true
		}
		if strings.HasPrefix(item.ValidatorName, "validator-domain-") {
			hasDomainResult = true
		}
	}
	if !hasCoreResult || !hasDomainResult {
		t.Fatalf("expected validator results include core+domain entries, got %+v", envelope.ValidatorResults)
	}

	if result.CurrentValidationRunID == "" || envelope.RunID == "" {
		t.Fatalf("run id should not be empty: envelope=%s current=%s", envelope.RunID, result.CurrentValidationRunID)
	}
	if envelope.RunID != result.CurrentValidationRunID {
		t.Fatalf("expected envelope run id == current run id, envelope=%s current=%s", envelope.RunID, result.CurrentValidationRunID)
	}
}

func TestM4ValidationTriggerKindsAndPlanScopedBlocking(t *testing.T) {
	t.Run("milestone_auto passed->completed", func(t *testing.T) {
		dbPath := compileMainIndex(t)
		db := openDB(t, dbPath)
		defer db.Close()

		requestPath := writeRequestFile(t, map[string]any{
			"request_id":                "req-m4-trigger-milestone-1",
			"task":                      "review WXT manifest permissions for browser store submission",
			"target_domain":             "wxt",
			"phase_id":                  "04",
			"plan_id":                   "02",
			"validation_trigger_kind":   "milestone_auto",
			"validation_trigger_reason": "plan milestone reached",
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
		envelope := result.ValidationResults[0]
		if envelope.TriggerKind != "milestone_auto" {
			t.Fatalf("expected trigger kind milestone_auto, got %s", envelope.TriggerKind)
		}
		if envelope.MachineView.Status != "passed" {
			t.Fatalf("expected machine status passed, got %s", envelope.MachineView.Status)
		}
		if result.Status != "completed" {
			t.Fatalf("expected activation completed, got %s", result.Status)
		}
		if len(envelope.HumanView.NextActions) == 0 {
			t.Fatalf("expected human next actions")
		}
	})

	t.Run("manual_rerun warned->partial", func(t *testing.T) {
		dbPath := compileMainIndex(t)
		db := openDB(t, dbPath)
		defer db.Close()

		requestPath := writeRequestFile(t, map[string]any{
			"request_id":                "req-m4-trigger-manual-1",
			"task":                      "review WXT manifest permissions",
			"target_domain":             "wxt",
			"phase_id":                  "04",
			"plan_id":                   "02",
			"validation_manual_rerun":   true,
			"validation_trigger_reason": "operator requested rerun",
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
		envelope := result.ValidationResults[0]
		if envelope.TriggerKind != "manual_rerun" {
			t.Fatalf("expected trigger kind manual_rerun, got %s", envelope.TriggerKind)
		}
		if envelope.MachineView.Status != "warned" {
			t.Fatalf("expected machine status warned, got %s", envelope.MachineView.Status)
		}
		if result.Status != "partial" {
			t.Fatalf("expected activation partial, got %s", result.Status)
		}
		if len(envelope.HumanView.NextActions) == 0 {
			t.Fatalf("expected warned path next actions")
		}
		if !strings.Contains(strings.ToLower(strings.Join(envelope.HumanView.NextActions, " ")), "run_id") {
			t.Fatalf("expected warned path next action contains run_id trace, got %+v", envelope.HumanView.NextActions)
		}
	})

	t.Run("rule_change_auto failed->failed and no cross-phase field", func(t *testing.T) {
		dbPath := compileMainIndex(t)
		db := openDB(t, dbPath)
		defer db.Close()

		requestPath := writeRequestFile(t, map[string]any{
			"request_id":                "req-m4-trigger-failed-1",
			"task":                      "handoff security permissions follow-up",
			"target_pack":               "security-permissions",
			"target_domain":             "security",
			"phase_id":                  "04",
			"plan_id":                   "02",
			"validation_trigger_kind":   "rule_change_auto",
			"validation_trigger_reason": "source_rule changed in validator input",
			"bounded_context": map[string]any{
				"selected_files":   []string{"manifest.json"},
				"config_fragments": []string{"manifest.permissions"},
				"host_hints":       []string{"browser-extension"},
				"browser_hints":    []string{"chrome"},
			},
			"context_hints": []string{"policy update"},
		})
		defer os.Remove(requestPath)

		result, err := activation.Execute(db, requestPath)
		if err != nil {
			t.Fatalf("activation execute failed: %v", err)
		}
		envelope := result.ValidationResults[0]
		if envelope.TriggerKind != "rule_change_auto" {
			t.Fatalf("expected trigger kind rule_change_auto, got %s", envelope.TriggerKind)
		}
		if envelope.MachineView.Status != "failed" {
			t.Fatalf("expected machine status failed, got %s", envelope.MachineView.Status)
		}
		if result.Status != "failed" {
			t.Fatalf("expected activation failed, got %s", result.Status)
		}
		if len(envelope.HumanView.NextActions) == 0 {
			t.Fatalf("expected human next actions")
		}
		if strings.Contains(strings.ToLower(strings.Join(envelope.HumanView.NextActions, "|")), "cross_phase") {
			t.Fatalf("unexpected cross phase action in next actions: %+v", envelope.HumanView.NextActions)
		}
	})
}
