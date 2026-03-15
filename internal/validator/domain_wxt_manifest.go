package validator

import (
	"strings"

	"the-agent-packs/internal/model"
)

func containsAny(values []string, needles ...string) bool {
	for _, v := range values {
		low := strings.ToLower(v)
		for _, needle := range needles {
			if strings.Contains(low, strings.ToLower(needle)) {
				return true
			}
		}
	}
	return false
}

func validateDomainWXTManifest(plan model.ValidationPlan, input ExecutionInput) model.ValidatorResult {
	if input.MainPack != "wxt-manifest" {
		return model.ValidatorResult{
			ValidatorName:      "validator-domain-wxt-manifest",
			Status:             "skipped",
			Findings:           []model.Finding{},
			RepairSuggestions:  []string{},
			ValidatedArtifacts: plan.ArtifactsUnderValidation,
		}
	}

	findings := []model.Finding{}
	repair := []string{}

	fragments := input.BoundedContext.ConfigFragments
	hints := input.BoundedContext.ContextHints
	hostHints := input.BoundedContext.HostHints
	browserHints := input.BoundedContext.BrowserHints
	taskLower := strings.ToLower(input.Task)

	if !containsAny(fragments, "manifest.permissions", "permissions") {
		findings = append(findings, model.Finding{
			Severity:    "warn",
			Code:        "missing-permissions-coverage",
			Message:     "Manifest permissions coverage is missing.",
			ArtifactRef: "manifest-review.md",
		})
		repair = append(repair, "Provide manifest.permissions fragment for permission review.")
	}

	if !containsAny(fragments, "manifest.host_permissions", "host_permissions") {
		findings = append(findings, model.Finding{
			Severity:    "warn",
			Code:        "missing-host-permissions-coverage",
			Message:     "Manifest host_permissions coverage is missing.",
			ArtifactRef: "manifest-review.md",
		})
		repair = append(repair, "Provide manifest.host_permissions fragment for host permission review.")
	}

	if len(browserHints) == 0 && !containsAny(hostHints, "chrome", "firefox", "edge", "safari") {
		findings = append(findings, model.Finding{
			Severity:    "warn",
			Code:        "missing-browser-override-note",
			Message:     "Browser-specific override hints are absent.",
			ArtifactRef: "manifest-review.md",
		})
		repair = append(repair, "Provide browser hints to evaluate browser-specific overrides.")
	}

	if !strings.Contains(taskLower, "store") && !containsAny(hints, "store", "submission", "review") {
		findings = append(findings, model.Finding{
			Severity:    "warn",
			Code:        "missing-store-facing-risk",
			Message:     "Store-facing review context is insufficient.",
			ArtifactRef: "manifest-review.md",
		})
		repair = append(repair, "Include store submission or policy review hints.")
	}

	status := "passed"
	if len(findings) > 0 {
		status = "warned"
	}

	return model.ValidatorResult{
		ValidatorName:      "validator-domain-wxt-manifest",
		Status:             status,
		Findings:           findings,
		RepairSuggestions:  repair,
		ValidatedArtifacts: plan.ArtifactsUnderValidation,
	}
}
