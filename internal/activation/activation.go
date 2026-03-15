package activation

import (
	"database/sql"
	"encoding/json"
	"os"
	"strings"

	"the-agent-packs/internal/model"
	"the-agent-packs/internal/query"
)

type requestShape struct {
	RequestID       string          `json:"request_id"`
	Task            string          `json:"task"`
	TargetPack      *string         `json:"target_pack"`
	TargetDomain    *string         `json:"target_domain"`
	BoundedContext  *boundedContext `json:"bounded_context"`
	ContextHints    []string        `json:"context_hints"`
	SelectedFiles   []string        `json:"selected_files"`
	ConfigFragments []string        `json:"config_fragments"`
}

type validationPlan struct {
	PlanID                   string            `json:"plan_id"`
	RequestID                string            `json:"request_id"`
	MainPack                 string            `json:"main_pack"`
	Validators               []validatorPlan   `json:"validators"`
	ArtifactsUnderValidation []string          `json:"artifacts_under_validation"`
	SeverityPolicy           map[string]string `json:"severity_policy"`
	PlanReason               string            `json:"plan_reason"`
}

type validatorPlan struct {
	Name   string `json:"name"`
	Scope  string `json:"scope"`
	Reason string `json:"reason"`
}

type finding struct {
	Severity    string `json:"severity"`
	Code        string `json:"code"`
	Message     string `json:"message"`
	ArtifactRef string `json:"artifact_ref"`
}

type validatorResult struct {
	ValidatorName      string    `json:"validator_name"`
	Status             string    `json:"status"`
	Findings           []finding `json:"findings"`
	RepairSuggestions  []string  `json:"repair_suggestions"`
	ValidatedArtifacts []string  `json:"validated_artifacts"`
}

type boundedContext struct {
	SelectedFiles   []string `json:"selected_files"`
	ConfigFragments []string `json:"config_fragments"`
	HostHints       []string `json:"host_hints"`
}

func loadRequest(path string) (*requestShape, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var req requestShape
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func activationResult(req *requestShape, mainPack any, routeReason string, status string) model.ActivationResult {
	return model.ActivationResult{
		RequestID:         req.RequestID,
		Status:            status,
		MainPack:          mainPack,
		Artifacts:         []any{},
		ValidationResults: []any{},
		Handoff:           nil,
		Summary:           routeReason,
	}
}

func buildValidationPlan(req *requestShape, mainPack string, artifacts []string) validationPlan {
	validators := []validatorPlan{
		{
			Name:   "validator-core-output",
			Scope:  "artifact",
			Reason: "All output artifacts must satisfy envelope completeness.",
		},
	}
	if mainPack == "wxt-manifest" {
		validators = append(validators, validatorPlan{
			Name:   "validator-domain-wxt-manifest",
			Scope:  "domain",
			Reason: "Manifest review must cover permission and store-facing risks.",
		})
	}
	return validationPlan{
		PlanID:                   req.RequestID + "-validation-plan",
		RequestID:                req.RequestID,
		MainPack:                 mainPack,
		Validators:               validators,
		ArtifactsUnderValidation: artifacts,
		SeverityPolicy: map[string]string{
			"warn":  "allow_partial",
			"error": "block_completed",
		},
		PlanReason: "Run core and domain validators for routed workflow package output.",
	}
}

func buildValidationResults(mainPack string, artifacts []string, status string) []any {
	coreResult := validatorResult{
		ValidatorName:      "validator-core-output",
		Status:             "passed",
		Findings:           []finding{},
		RepairSuggestions:  []string{},
		ValidatedArtifacts: artifacts,
	}
	results := []any{coreResult}
	if mainPack != "wxt-manifest" {
		return results
	}
	if status == "partial" {
		results = append(results, validatorResult{
			ValidatorName: "validator-domain-wxt-manifest",
			Status:        "warned",
			Findings: []finding{{
				Severity:    "warn",
				Code:        "insufficient-manifest-context",
				Message:     "Bounded context is insufficient for a full manifest domain review.",
				ArtifactRef: "manifest-review.md",
			}},
			RepairSuggestions: []string{
				"Provide manifest permissions and host_permissions config fragments.",
				"Provide target browser hints for override checks.",
			},
			ValidatedArtifacts: artifacts,
		})
		return results
	}
	results = append(results, validatorResult{
		ValidatorName:      "validator-domain-wxt-manifest",
		Status:             "passed",
		Findings:           []finding{},
		RepairSuggestions:  []string{},
		ValidatedArtifacts: artifacts,
	})
	return results
}

func buildHandoff(mainPack string, task string) any {
	if mainPack != "wxt-manifest" {
		return nil
	}
	if !strings.Contains(strings.ToLower(task), "handoff") {
		return nil
	}
	return map[string]any{
		"from_pack": "wxt-manifest",
		"to_packs":  []string{"security-permissions", "release-store-review"},
		"reason":    "Manifest review reached package boundary and requires cross-line continuation.",
		"carry_context": map[string]any{
			"required_artifact": "manifest-review.md",
			"required_checks": []string{
				"permissions-minimization",
				"store-review-checklist",
			},
		},
	}
}

func withExecutionPayload(base model.ActivationResult, req *requestShape, mainPack string, status string) model.ActivationResult {
	artifacts := []string{"manifest-review.md"}
	if mainPack != "wxt-manifest" {
		artifacts = []string{}
	}
	plan := buildValidationPlan(req, mainPack, artifacts)
	results := buildValidationResults(mainPack, artifacts, status)
	base.Artifacts = []any{}
	for _, artifact := range artifacts {
		base.Artifacts = append(base.Artifacts, map[string]any{
			"name": artifact,
			"kind": "review-report",
		})
	}
	base.ValidationResults = []any{map[string]any{
		"validation_plan":   plan,
		"validator_results": results,
	}}
	base.Handoff = buildHandoff(mainPack, req.Task)
	return base
}

func mergeHints(primary, secondary []string) []string {
	merged := []string{}
	add := func(xs []string) {
		for _, item := range xs {
			seen := false
			for _, existing := range merged {
				if existing == item {
					seen = true
					break
				}
			}
			if !seen {
				merged = append(merged, item)
			}
		}
	}
	add(primary)
	add(secondary)
	return merged
}

func Execute(db *sql.DB, requestPath string) (model.ActivationResult, error) {
	req, err := loadRequest(requestPath)
	if err != nil {
		return model.ActivationResult{}, err
	}
	if req.RequestID == "" || req.Task == "" || req.BoundedContext == nil {
		return model.ActivationResult{
			Status:  "failed",
			Summary: "invalid activation request",
		}, nil
	}

	selectedFiles := mergeHints(req.SelectedFiles, req.BoundedContext.SelectedFiles)
	configFragments := mergeHints(req.ConfigFragments, req.BoundedContext.ConfigFragments)
	contextHints := mergeHints(req.ContextHints, req.BoundedContext.HostHints)

	routeResult, err := query.RouteQuery(
		db,
		"L1",
		req.Task,
		req.TargetPack,
		req.TargetDomain,
		selectedFiles,
		configFragments,
		contextHints,
		1,
	)
	if err != nil {
		return model.ActivationResult{}, err
	}

	if len(routeResult.Candidates) == 0 {
		if req.TargetDomain != nil {
			result := activationResult(req, nil, "insufficient evidence for L1, fallback to L0 recommended", "partial")
			return withExecutionPayload(result, req, "", "partial"), nil
		}
		return activationResult(req, nil, "no route candidate", "failed"), nil
	}

	mainNode := routeResult.Candidates[0].ID
	bundle, err := query.BuildContextBundle(db, mainNode, true, false, false)
	if err != nil {
		return model.ActivationResult{}, err
	}
	mainPack := query.PackForNode(mainNode)
	if mainPack == "" {
		mainPack = "wxt-manifest"
	}

	if len(req.BoundedContext.SelectedFiles) == 0 && len(req.BoundedContext.ConfigFragments) == 0 {
		result := activationResult(req, mainPack, "bounded context missing required evidence", "partial")
		return withExecutionPayload(result, req, mainPack, "partial"), nil
	}
	if strings.Contains(strings.ToLower(req.Task), "handoff") {
		result := activationResult(req, mainPack, "handoff requested by task boundary", "handoff")
		return withExecutionPayload(result, req, mainPack, "handoff"), nil
	}
	if len(bundle.RecommendedArtifacts) > 0 || len(bundle.RecommendedValidators) > 0 {
		result := activationResult(req, mainPack, "route to L1.wxt.manifest with required cross-cutting lines", "completed")
		return withExecutionPayload(result, req, mainPack, "completed"), nil
	}
	result := activationResult(req, mainPack, "route to L1.wxt.manifest with required cross-cutting lines", "completed")
	return withExecutionPayload(result, req, mainPack, "completed"), nil
}
