package validator

import "the-agent-packs/internal/model"

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

	if input.RequestedHandoff && len(input.Handoff) == 0 {
		findings = append(findings, model.Finding{
			Severity:    "error",
			Code:        "handoff-missing",
			Message:     "Handoff is requested but handoff payload is missing.",
			ArtifactRef: "manifest-review.md",
		})
		repair = append(repair, "Provide handoff payload with carry_context for downstream packs.")
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
