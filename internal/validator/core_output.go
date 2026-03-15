package validator

import (
	"fmt"

	"the-agent-packs/internal/model"
)

func validateCoreOutput(plan model.ValidationPlan, input ExecutionInput) model.ValidatorResult {
	findings := []model.Finding{}
	repair := []string{}

	if len(input.Artifacts) == 0 {
		findings = append(findings, model.Finding{
			Severity:    "error",
			Code:        "missing-artifact",
			Message:     "No output artifacts were generated.",
			ArtifactRef: "",
		})
		repair = append(repair, "Generate the primary artifact before running validators.")
	}

	for _, artifact := range input.Artifacts {
		if artifact.Name == "" {
			findings = append(findings, model.Finding{
				Severity:    "error",
				Code:        "artifact-name-empty",
				Message:     "Artifact name is required.",
				ArtifactRef: "",
			})
		}
		if artifact.Kind == "" {
			findings = append(findings, model.Finding{
				Severity:    "error",
				Code:        "artifact-kind-empty",
				Message:     "Artifact kind is required.",
				ArtifactRef: artifact.Name,
			})
		}
	}

	if input.MainPack != "" && len(input.RequiredPacks) == 0 {
		findings = append(findings, model.Finding{
			Severity:    "warn",
			Code:        "required-packs-missing",
			Message:     "Main pack is missing required pack declarations in validator input.",
			ArtifactRef: input.MainPack,
		})
		repair = append(repair, "Populate required packs from registry-backed context bundle before validation.")
	}

	if input.RequestedHandoff && len(input.Handoff) == 0 {
		findings = append(findings, model.Finding{
			Severity:    "error",
			Code:        "handoff-missing",
			Message:     "Handoff is requested but handoff payload is missing.",
			ArtifactRef: "manifest-review.md",
		})
		repair = append(repair, "Provide handoff payload with carry_context for downstream packs.")
	}

	if len(input.RequiredPacks) > 0 && len(input.Handoff) > 0 {
		toPacks, ok := input.Handoff["to_packs"].([]string)
		if !ok {
			if raw, rawOK := input.Handoff["to_packs"].([]any); rawOK {
				toPacks = make([]string, 0, len(raw))
				for _, item := range raw {
					if s, sOK := item.(string); sOK {
						toPacks = append(toPacks, s)
					}
				}
			}
		}
		if !samePackSet(input.RequiredPacks, toPacks) {
			findings = append(findings, model.Finding{
				Severity:    "error",
				Code:        "handoff-required-packs-mismatch",
				Message:     fmt.Sprintf("Handoff to_packs must match required packs. required=%v handoff=%v", input.RequiredPacks, toPacks),
				ArtifactRef: "manifest-review.md",
			})
			repair = append(repair, "Align handoff to_packs with registry required_packs declarations.")
		}
	}

	status := "passed"
	for _, f := range findings {
		if f.Severity == "error" {
			status = "failed"
			break
		}
	}

	validated := plan.ArtifactsUnderValidation
	if len(validated) == 0 {
		validated = []string{}
	}

	return model.ValidatorResult{
		ValidatorName:      "validator-core-output",
		Status:             status,
		Findings:           findings,
		RepairSuggestions:  repair,
		ValidatedArtifacts: validated,
	}
}

func samePackSet(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	seen := map[string]int{}
	for _, item := range left {
		seen[item]++
	}
	for _, item := range right {
		seen[item]--
	}
	for _, count := range seen {
		if count != 0 {
			return false
		}
	}
	return true
}
