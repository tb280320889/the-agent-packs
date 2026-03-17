package activation

import (
	"database/sql"
	"encoding/json"
	"os"
	"strings"

	"the-agent-packs/internal/model"
	"the-agent-packs/internal/query"
	"the-agent-packs/internal/validator"
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

type boundedContext struct {
	SelectedFiles   []string `json:"selected_files"`
	ConfigFragments []string `json:"config_fragments"`
	HostHints       []string `json:"host_hints"`
	BrowserHints    []string `json:"browser_hints"`
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

func ptr(s string) *string {
	if s == "" {
		return nil
	}
	v := s
	return &v
}

func buildValidationPlan(req *requestShape, mainPack string, artifacts []model.Artifact, bundleValidators []string) model.ValidationPlan {
	if len(bundleValidators) == 0 {
		bundleValidators = []string{"validator-core-output"}
	}
	validators := make([]model.ValidatorPlan, 0, len(bundleValidators))
	for _, name := range bundleValidators {
		scope := "artifact"
		reason := "All output artifacts must satisfy envelope completeness."
		if name != "validator-core-output" {
			scope = "domain"
			reason = "Domain validator must verify workflow-specific runtime constraints."
		}
		if name == "validator-domain-wxt-manifest" {
			reason = "Manifest review must cover permission and store-facing risks."
		}
		validators = append(validators, model.ValidatorPlan{Name: name, Scope: scope, Reason: reason})
	}
	artifactNames := make([]string, 0, len(artifacts))
	for _, artifact := range artifacts {
		artifactNames = append(artifactNames, artifact.Name)
	}
	return model.ValidationPlan{
		PlanID:                   req.RequestID + "-validation-plan",
		RequestID:                req.RequestID,
		MainPack:                 mainPack,
		Validators:               validators,
		ArtifactsUnderValidation: artifactNames,
		SeverityPolicy: map[string]string{
			"warn":  "allow_partial",
			"error": "block_completed",
		},
		PlanReason: "Run core and domain validators for routed workflow package output.",
	}
}

func buildHandoff(mainPack string, task string, requiredPacks []string, artifacts []model.Artifact) map[string]any {
	if len(requiredPacks) == 0 {
		return nil
	}
	if !strings.Contains(strings.ToLower(task), "handoff") {
		return nil
	}
	requiredArtifact := ""
	if len(artifacts) > 0 {
		requiredArtifact = artifacts[0].Name
	}
	requiredChecks := make([]string, 0, len(requiredPacks))
	for _, pack := range requiredPacks {
		requiredChecks = append(requiredChecks, pack+"-ready")
	}
	return map[string]any{
		"from_pack": mainPack,
		"to_packs":  append([]string{}, requiredPacks...),
		"reason":    "Workflow review reached package boundary and requires registered continuation.",
		"carry_context": map[string]any{
			"required_artifact": requiredArtifact,
			"required_checks":   requiredChecks,
		},
	}
}

func hasBlockingFailure(policy map[string]string, results []model.ValidatorResult) bool {
	for _, result := range results {
		if result.Status == "failed" {
			return true
		}
		for _, f := range result.Findings {
			if f.Severity == "error" && policy["error"] == "block_completed" {
				return true
			}
		}
	}
	return false
}

func hasWarnFindings(results []model.ValidatorResult) bool {
	for _, result := range results {
		if result.Status == "warned" {
			return true
		}
		for _, f := range result.Findings {
			if f.Severity == "warn" {
				return true
			}
		}
	}
	return false
}

func deriveStatus(requestInvalid bool, routeMissing bool, contextInsufficient bool, handoff map[string]any, plan model.ValidationPlan, results []model.ValidatorResult) string {
	if requestInvalid {
		return "failed"
	}
	if routeMissing {
		return "failed"
	}
	if hasBlockingFailure(plan.SeverityPolicy, results) {
		return "failed"
	}
	if len(handoff) > 0 {
		return "handoff"
	}
	if contextInsufficient {
		return "partial"
	}
	if hasWarnFindings(results) && plan.SeverityPolicy["warn"] == "allow_partial" {
		return "partial"
	}
	return "completed"
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

	invalid := req.RequestID == "" || req.Task == "" || req.BoundedContext == nil
	if invalid {
		return model.ActivationResult{
			RequestID:         req.RequestID,
			Status:            "failed",
			MainPack:          nil,
			RouteStatus:       "failed",
			RouteErrorCode:    "ACTIVATION_REQUEST_INVALID",
			RouteNextAction:   "补全 request_id/task/bounded_context 后重试",
			RouteDecision:     "activation_request_validation",
			RouteTraceID:      "activation:request:invalid",
			RouteDocsRef:      "",
			Artifacts:         []model.Artifact{},
			ValidationResults: []model.ValidationEnvelope{},
			Handoff:           nil,
			Summary:           "invalid activation request",
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
			return model.ActivationResult{
				RequestID:         req.RequestID,
				Status:            "partial",
				MainPack:          nil,
				RouteStatus:       routeResult.Status,
				RouteErrorCode:    routeResult.ErrorCode,
				RouteNextAction:   routeResult.NextAction,
				RouteDecision:     routeResult.DecisionBasis,
				RouteTraceID:      routeResult.DecisionTraceID,
				RouteDocsRef:      routeResult.DocsRef,
				Artifacts:         []model.Artifact{},
				ValidationResults: []model.ValidationEnvelope{},
				Handoff:           nil,
				Summary:           routeResult.Message,
			}, nil
		}
		return model.ActivationResult{
			RequestID:         req.RequestID,
			Status:            "failed",
			MainPack:          nil,
			RouteStatus:       routeResult.Status,
			RouteErrorCode:    routeResult.ErrorCode,
			RouteNextAction:   routeResult.NextAction,
			RouteDecision:     routeResult.DecisionBasis,
			RouteTraceID:      routeResult.DecisionTraceID,
			RouteDocsRef:      routeResult.DocsRef,
			Artifacts:         []model.Artifact{},
			ValidationResults: []model.ValidationEnvelope{},
			Handoff:           nil,
			Summary:           routeResult.Message,
		}, nil
	}

	mainNode := routeResult.Candidates[0].ID
	bundle, err := query.BuildContextBundle(db, mainNode, true, false, true)
	if err != nil {
		return model.ActivationResult{}, err
	}

	mainPack := query.PackForNode(mainNode)
	if mainPack == "" {
		mainPack = "wxt-manifest"
	}

	artifacts := make([]model.Artifact, 0, len(bundle.RecommendedArtifacts))
	for _, artifactName := range bundle.RecommendedArtifacts {
		artifacts = append(artifacts, model.Artifact{Name: artifactName, Kind: inferArtifactKind(artifactName)})
	}
	if len(artifacts) == 0 && mainPack == "wxt-manifest" {
		artifacts = append(artifacts, model.Artifact{Name: "manifest-review.md", Kind: "review-report"})
	}

	handoff := buildHandoff(mainPack, req.Task, bundle.RequiredPacks, artifacts)
	contextInsufficient := len(req.BoundedContext.SelectedFiles) == 0 && len(req.BoundedContext.ConfigFragments) == 0

	plan := buildValidationPlan(req, mainPack, artifacts, append([]string{}, bundle.RecommendedValidators...))
	vInput := validator.ExecutionInput{
		Task:           req.Task,
		MainPack:       mainPack,
		ContractBundle: &bundle,
		Artifacts:      artifacts,
		RequiredPacks:  append([]string{}, bundle.RequiredPacks...),
		BoundedContext: validator.BoundedContextSnapshot{
			SelectedFiles:   req.BoundedContext.SelectedFiles,
			ConfigFragments: req.BoundedContext.ConfigFragments,
			HostHints:       req.BoundedContext.HostHints,
			BrowserHints:    req.BoundedContext.BrowserHints,
			ContextHints:    req.ContextHints,
		},
		RequestedHandoff: len(handoff) > 0,
		Handoff:          handoff,
	}
	validatorResults := validator.Run(plan, vInput)
	status := deriveStatus(false, false, contextInsufficient, handoff, plan, validatorResults)

	summary := "route to L1.wxt.manifest with required cross-cutting lines"
	if contextInsufficient {
		summary = "bounded context missing required evidence"
	}
	if len(handoff) > 0 {
		summary = "handoff requested by task boundary"
	}
	if len(bundle.RecommendedValidators) == 0 && len(bundle.RecommendedArtifacts) == 0 {
		summary = "route succeeded with minimal bundle"
	}

	return model.ActivationResult{
		RequestID:       req.RequestID,
		Status:          status,
		MainPack:        ptr(mainPack),
		RouteStatus:     routeResult.Status,
		RouteErrorCode:  routeResult.ErrorCode,
		RouteNextAction: routeResult.NextAction,
		RouteDecision:   routeResult.DecisionBasis,
		RouteTraceID:    routeResult.DecisionTraceID,
		RouteDocsRef:    routeResult.DocsRef,
		Artifacts:       artifacts,
		ValidationResults: []model.ValidationEnvelope{{
			ValidationPlan:   plan,
			ValidatorResults: validatorResults,
		}},
		Handoff: handoff,
		Summary: summary,
	}, nil
}

func inferArtifactKind(name string) string {
	if strings.HasSuffix(strings.ToLower(name), ".md") {
		return "review-report"
	}
	return "artifact"
}
