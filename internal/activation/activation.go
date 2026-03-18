package activation

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"the-agent-packs/internal/model"
	"the-agent-packs/internal/query"
	"the-agent-packs/internal/validator"
)

type requestShape struct {
	RequestID       string          `json:"request_id"`
	Task            string          `json:"task"`
	PhaseID         string          `json:"phase_id"`
	PlanID          string          `json:"plan_id"`
	TriggerKind     string          `json:"validation_trigger_kind"`
	TriggerReason   string          `json:"validation_trigger_reason"`
	ManualRerun     bool            `json:"validation_manual_rerun"`
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

type LedgerWriteMode string

const (
	LedgerWriteModeImmediate     LedgerWriteMode = "immediate"
	LedgerWriteModeBatchFinalize LedgerWriteMode = "batch_finalize"

	runtimeLedgerDefaultDeferredWindow = "24h"
)

func runtimeLedgerDeferredWindow() time.Duration {
	raw := strings.TrimSpace(os.Getenv("RUNTIME_LEDGER_DEFER_WINDOW"))
	if raw == "" {
		raw = runtimeLedgerDefaultDeferredWindow
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 24 * time.Hour
	}
	return d
}

type RuntimeLedgerBuildInput struct {
	TraceID             string
	RunID               string
	TriggerKind         string
	MachineStatus       string
	Timestamp           time.Time
	SourceRefs          []string
	Finalized           bool
	DeferredReason      string
	DeferredDeadline    string
	HasBlockingDecision bool
}

func resolveLedgerWriteMode(triggerKind, machineStatus string, hasBlockingDecision bool) LedgerWriteMode {
	if triggerKind == model.ValidationTriggerRuleChangeAuto {
		return LedgerWriteModeImmediate
	}
	if machineStatus == model.ValidationStatusFailed {
		return LedgerWriteModeImmediate
	}
	if hasBlockingDecision {
		return LedgerWriteModeImmediate
	}
	return LedgerWriteModeBatchFinalize
}

func BuildRuntimeLedgerEntries(existing []model.RuntimeLedgerEntry, input RuntimeLedgerBuildInput) ([]model.RuntimeLedgerEntry, model.RuntimeLedgerEntry, LedgerWriteMode) {
	now := input.Timestamp.UTC()
	if now.IsZero() {
		now = time.Now().UTC()
	}
	traceID := strings.TrimSpace(input.TraceID)
	if traceID == "" {
		traceID = "runtime-ledger:unknown"
	}
	mode := resolveLedgerWriteMode(input.TriggerKind, input.MachineStatus, input.HasBlockingDecision)
	recordTypes := determineRuntimeLedgerRecordTypes(input, mode)

	updated := append([]model.RuntimeLedgerEntry{}, existing...)
	latestByType := map[string]model.RuntimeLedgerEntry{}
	for _, recordType := range recordTypes {
		entry := model.RuntimeLedgerEntry{
			TraceID:    traceID,
			RunID:      input.RunID,
			RecordType: recordType,
			Timestamp:  now.Format(time.RFC3339),
			SourceRefs: append([]string{}, input.SourceRefs...),
			Status:     input.MachineStatus,
		}
		if !model.IsRuntimeLedgerRecordType(entry.RecordType) {
			continue
		}
		if mode == LedgerWriteModeBatchFinalize && !input.Finalized {
			entry.Status = "deferred"
			reason := strings.TrimSpace(input.DeferredReason)
			if reason == "" {
				reason = "awaiting plan finalization"
			}
			entry.DeferredReason = reason

			deadline := strings.TrimSpace(input.DeferredDeadline)
			if deadline == "" {
				deadline = now.Add(runtimeLedgerDeferredWindow()).UTC().Format(time.RFC3339)
			}
			for _, prior := range updated {
				if prior.TraceID == traceID && prior.RecordType == entry.RecordType && prior.IsCurrent && strings.TrimSpace(prior.DeferredDeadline) != "" {
					deadline = prior.DeferredDeadline
					break
				}
			}
			entry.DeferredDeadline = deadline
			if parsed, err := time.Parse(time.RFC3339, deadline); err == nil && now.After(parsed) {
				entry.RiskEscalated = true
			}
		}

		updated = appendRuntimeLedgerVersion(updated, entry)
		latestByType[recordType] = updated[len(updated)-1]
	}

	if validationEntry, ok := latestByType[model.RuntimeLedgerRecordTypeValidation]; ok {
		if hasEscalationInEntries(latestByType) {
			validationEntry.RiskEscalated = true
		}
		return updated, validationEntry, mode
	}

	if len(updated) == 0 {
		return existing, model.RuntimeLedgerEntry{}, mode
	}
	return updated, updated[len(updated)-1], mode
}

func determineRuntimeLedgerRecordTypes(input RuntimeLedgerBuildInput, mode LedgerWriteMode) []string {
	typeSet := map[string]bool{
		model.RuntimeLedgerRecordTypeValidation: true,
	}

	if input.TriggerKind == model.ValidationTriggerRuleChangeAuto || input.TriggerKind == model.ValidationTriggerValidatorManifestChangeAuto {
		typeSet[model.RuntimeLedgerRecordTypeChange] = true
	}

	if input.MachineStatus == model.ValidationStatusFailed || input.TriggerKind == model.ValidationTriggerManualRerun {
		typeSet[model.RuntimeLedgerRecordTypeDecision] = true
	}

	hasDeferredInfo := strings.TrimSpace(input.DeferredReason) != "" || strings.TrimSpace(input.DeferredDeadline) != ""
	if mode == LedgerWriteModeBatchFinalize && hasDeferredInfo {
		typeSet[model.RuntimeLedgerRecordTypeAssumption] = true
	}

	ordered := make([]string, 0, len(model.RuntimeLedgerRecordTypes))
	for _, recordType := range model.RuntimeLedgerRecordTypes {
		if typeSet[recordType] {
			ordered = append(ordered, recordType)
		}
	}
	return ordered
}

func hasEscalationInEntries(entries map[string]model.RuntimeLedgerEntry) bool {
	for _, entry := range entries {
		if entry.RiskEscalated {
			return true
		}
	}
	return false
}

func appendRuntimeLedgerVersion(existing []model.RuntimeLedgerEntry, incoming model.RuntimeLedgerEntry) []model.RuntimeLedgerEntry {
	updated := make([]model.RuntimeLedgerEntry, len(existing))
	copy(updated, existing)
	maxVersion := 0
	for i := range updated {
		entry := updated[i]
		if entry.TraceID == incoming.TraceID && entry.RecordType == incoming.RecordType {
			if entry.Version > maxVersion {
				maxVersion = entry.Version
			}
			updated[i].IsCurrent = false
		}
	}
	incoming.Version = maxVersion + 1
	incoming.IsCurrent = true
	updated = append(updated, incoming)
	return updated
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

func buildValidationPlan(req *requestShape, mainPack string, artifacts []model.Artifact, requiredPacks []string, bundleValidators []string) model.ValidationPlan {
	profile := query.RecommendedValidators(mainPack, requiredPacks)
	validatorSet := map[string]bool{}
	validatorNames := make([]string, 0, len(profile.Validators)+len(bundleValidators)+1)
	appendValidator := func(name string) {
		name = strings.TrimSpace(name)
		if name == "" || validatorSet[name] {
			return
		}
		validatorSet[name] = true
		validatorNames = append(validatorNames, name)
	}
	appendValidator("validator-core-output")
	for _, name := range profile.Validators {
		appendValidator(name)
	}
	for _, name := range bundleValidators {
		appendValidator(name)
	}
	rest := make([]string, 0, len(validatorNames))
	for _, name := range validatorNames {
		if name == "validator-core-output" {
			continue
		}
		rest = append(rest, name)
	}
	sort.Strings(rest)
	validatorNames = append([]string{"validator-core-output"}, rest...)

	validators := make([]model.ValidatorPlan, 0, len(validatorNames))
	for _, name := range validatorNames {
		scope := "domain"
		reason := "Registry declared validator for core output integrity."
		if name == "validator-core-output" {
			scope = "artifact"
		}
		if strings.HasPrefix(name, "validator-domain-") || strings.HasPrefix(name, "validator-contract-") {
			scope = "domain"
			reason = fmt.Sprintf("Registry declared validator for %s", strings.TrimPrefix(name, "validator-domain-"))
		}
		if src := profile.Sources[name]; len(src) > 0 && name != "validator-core-output" {
			reason = fmt.Sprintf("Registry declared validator for %s", strings.Join(src, "+"))
		}
		if name == "validator-domain-wxt-manifest" {
			reason = "Registry declared validator for wxt-manifest"
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
		PlanReason: fmt.Sprintf("registry-defined plan for core+domain validators (signature=%s)", profile.Signature),
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

func deriveStatuses(requestInvalid bool, routeMissing bool, plan model.ValidationPlan, results []model.ValidatorResult) (string, string) {
	if requestInvalid {
		return model.ActivationStatusFailed, model.ValidationStatusFailed
	}
	if routeMissing {
		return model.ActivationStatusFailed, model.ValidationStatusFailed
	}
	if hasBlockingFailure(plan.SeverityPolicy, results) {
		return model.ActivationStatusFailed, model.ValidationStatusFailed
	}
	if hasWarnFindings(results) && plan.SeverityPolicy["warn"] == "allow_partial" {
		return model.ActivationStatusPartial, model.ValidationStatusWarned
	}
	return model.ActivationStatusCompleted, model.ValidationStatusPassed
}

func resolveTrigger(req *requestShape) (string, string) {
	if req.ManualRerun {
		return model.ValidationTriggerManualRerun, "manual validation rerun requested"
	}
	kind := strings.TrimSpace(req.TriggerKind)
	if kind == "" {
		kind = model.ValidationTriggerMilestoneAuto
	}
	reason := strings.TrimSpace(req.TriggerReason)
	if reason == "" {
		switch kind {
		case model.ValidationTriggerRuleChangeAuto:
			reason = "auto-triggered by rule change"
		case model.ValidationTriggerValidatorManifestChangeAuto:
			reason = "auto-triggered by validator manifest change"
		case model.ValidationTriggerManualRerun:
			reason = "manual validation rerun requested"
		default:
			kind = model.ValidationTriggerMilestoneAuto
			reason = "auto-triggered by milestone validation"
		}
	}
	allowed := map[string]bool{
		model.ValidationTriggerMilestoneAuto:               true,
		model.ValidationTriggerRuleChangeAuto:              true,
		model.ValidationTriggerValidatorManifestChangeAuto: true,
		model.ValidationTriggerManualRerun:                 true,
	}
	if !allowed[kind] {
		kind = model.ValidationTriggerMilestoneAuto
		reason = "auto-triggered by milestone validation"
	}
	return kind, reason
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

func fallbackOrDefault(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func collectCodesAndSuggestions(results []model.ValidatorResult) ([]string, []string, []string) {
	errorCodes := []string{}
	ruleCodes := []string{}
	repairSuggestions := []string{}

	appendUnique := func(target *[]string, values ...string) {
		for _, value := range values {
			trimmed := strings.TrimSpace(value)
			if trimmed == "" {
				continue
			}
			exists := false
			for _, existing := range *target {
				if existing == trimmed {
					exists = true
					break
				}
			}
			if !exists {
				*target = append(*target, trimmed)
			}
		}
	}

	for _, result := range results {
		appendUnique(&repairSuggestions, result.RepairSuggestions...)
		for _, finding := range result.Findings {
			if finding.Severity == "error" {
				appendUnique(&errorCodes, finding.Code)
			}
			appendUnique(&ruleCodes, finding.RuleRef, finding.SourceRule)
		}
	}

	return errorCodes, ruleCodes, repairSuggestions
}

func makeInputDigest(req *requestShape, artifacts []model.Artifact) string {
	snapshot := map[string]any{
		"request_id":       req.RequestID,
		"task":             req.Task,
		"phase_id":         req.PhaseID,
		"plan_id":          req.PlanID,
		"trigger_kind":     req.TriggerKind,
		"trigger_reason":   req.TriggerReason,
		"selected_files":   req.SelectedFiles,
		"config_fragments": req.ConfigFragments,
		"context_hints":    req.ContextHints,
		"artifacts":        artifacts,
	}
	payload, _ := json.Marshal(snapshot)
	digest := sha256.Sum256(payload)
	return fmt.Sprintf("sha256:%x", digest)
}

func buildEvidenceRefs(requestID, runID string, artifacts []model.Artifact, handoff map[string]any) []model.ValidationEvidenceRef {
	evidenceRefs := make([]model.ValidationEvidenceRef, 0, len(artifacts)+2)
	for _, artifact := range artifacts {
		if strings.TrimSpace(artifact.Name) == "" {
			continue
		}
		evidenceRefs = append(evidenceRefs, model.ValidationEvidenceRef{
			RefID:      "artifact:" + artifact.Name,
			RefType:    "artifact",
			RefPath:    artifact.Name,
			StrongLink: false,
		})
	}

	if len(handoff) > 0 {
		evidenceRefs = append(evidenceRefs, model.ValidationEvidenceRef{
			RefID:      "handoff:" + requestID,
			RefType:    "handoff",
			RefPath:    "handoff",
			StrongLink: true,
		})
	}

	evidenceRefs = append(evidenceRefs, model.ValidationEvidenceRef{
		RefID:      "runtime-ledger:" + runID,
		RefType:    "runtime-ledger",
		RefPath:    "runtime/ledger/" + runID,
		StrongLink: true,
	})

	return evidenceRefs
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
	if len(artifacts) == 0 {
		artifacts = append(artifacts, model.Artifact{Name: "activation-output.json", Kind: "artifact"})
	}

	handoff := buildHandoff(mainPack, req.Task, bundle.RequiredPacks, artifacts)
	contextInsufficient := len(req.BoundedContext.SelectedFiles) == 0 && len(req.BoundedContext.ConfigFragments) == 0

	plan := buildValidationPlan(req, mainPack, artifacts, append([]string{}, bundle.RequiredPacks...), append([]string{}, bundle.RecommendedValidators...))
	triggerKind, triggerReason := resolveTrigger(req)
	vInput := validator.ExecutionInput{
		Task:           req.Task,
		MainPack:       mainPack,
		PhaseID:        fallbackOrDefault(req.PhaseID, "04"),
		PlanID:         fallbackOrDefault(req.PlanID, "unknown"),
		TriggerKind:    triggerKind,
		TriggerReason:  triggerReason,
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
		RequestedHandoff: strings.Contains(strings.ToLower(req.Task), "handoff"),
		Handoff:          handoff,
	}
	validatorResults := validator.Run(plan, vInput)
	status, machineStatus := deriveStatuses(false, false, plan, validatorResults)
	runID := fmt.Sprintf("%s:validation:%d", req.RequestID, time.Now().Unix())
	errorCodes, ruleCodes, repairSuggestions := collectCodesAndSuggestions(validatorResults)
	evidenceRefs := buildEvidenceRefs(req.RequestID, runID, vInput.Artifacts, handoff)
	humanSummary := fmt.Sprintf("本次校验结论为 %s，已生成 %d 条证据引用并可追溯到运行账本。", machineStatus, len(evidenceRefs))
	nextActions := []string{"如存在 error 级别问题，请先修复后再触发下一次 validation。"}
	if machineStatus == model.ValidationStatusPassed {
		nextActions = []string{"当前校验通过，可继续后续计划执行。"}
	}
	if machineStatus == model.ValidationStatusWarned {
		nextActions = []string{fmt.Sprintf("记录 warned 处理意见并链接 run_id=%s", runID)}
	}

	runtimeLedger, ledgerEntry, _ := BuildRuntimeLedgerEntries(nil, RuntimeLedgerBuildInput{
		TraceID:             fmt.Sprintf("runtime-ledger:%s:%s:%s", req.RequestID, vInput.PhaseID, vInput.PlanID),
		RunID:               runID,
		TriggerKind:         vInput.TriggerKind,
		MachineStatus:       machineStatus,
		Timestamp:           time.Now().UTC(),
		SourceRefs:          []string{routeResult.DocsRef, "docs/AIDP/runtime/06-验证记录.md", "docs/AIDP/runtime/03-变更摘要.md"},
		Finalized:           false,
		DeferredReason:      "awaiting plan finalization runtime writeback",
		HasBlockingDecision: status == model.ActivationStatusFailed,
	})
	if ledgerEntry.RiskEscalated {
		nextActions = append(nextActions, "runtime-ledger-overdue")
		humanSummary = humanSummary + " runtime-ledger-overdue"
	}

	currentValidation := model.ValidationEnvelope{
		RunID:              runID,
		PhaseID:            vInput.PhaseID,
		PlanID:             vInput.PlanID,
		TriggerKind:        vInput.TriggerKind,
		TriggerReason:      vInput.TriggerReason,
		IsCurrentEffective: true,
		InputDigest:        makeInputDigest(req, vInput.Artifacts),
		EvidenceRefs:       evidenceRefs,
		MachineView: model.ValidationMachineView{
			Status:            machineStatus,
			ErrorCodes:        errorCodes,
			RuleCodes:         ruleCodes,
			Trigger:           vInput.TriggerKind,
			RepairSuggestions: repairSuggestions,
		},
		HumanView: model.ValidationHumanView{
			Summary:     humanSummary,
			NextActions: nextActions,
		},
		ValidationPlan:   plan,
		ValidatorResults: validatorResults,
	}

	summary := "route to L1.wxt.manifest with required cross-cutting lines"
	if contextInsufficient {
		summary = "bounded context missing required evidence"
	}
	if len(handoff) > 0 && status == model.ActivationStatusCompleted {
		summary = "handoff requested by task boundary"
	}
	if len(bundle.RecommendedValidators) == 0 && len(bundle.RecommendedArtifacts) == 0 {
		summary = "route succeeded with minimal bundle"
	}

	return model.ActivationResult{
		RequestID:              req.RequestID,
		Status:                 status,
		MainPack:               ptr(mainPack),
		RouteStatus:            routeResult.Status,
		RouteErrorCode:         routeResult.ErrorCode,
		RouteNextAction:        routeResult.NextAction,
		RouteDecision:          routeResult.DecisionBasis,
		RouteTraceID:           routeResult.DecisionTraceID,
		RouteDocsRef:           routeResult.DocsRef,
		Artifacts:              artifacts,
		ValidationResults:      []model.ValidationEnvelope{currentValidation},
		ValidationRunHistory:   []model.ValidationEnvelope{currentValidation},
		CurrentValidationRunID: runID,
		RuntimeLedger:          runtimeLedger,
		Handoff:                handoff,
		Summary:                summary,
	}, nil
}

func inferArtifactKind(name string) string {
	if strings.HasSuffix(strings.ToLower(name), ".md") {
		return "review-report"
	}
	return "artifact"
}
