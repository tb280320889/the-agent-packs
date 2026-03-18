package validator

import "the-agent-packs/internal/model"

type BoundedContextSnapshot struct {
	SelectedFiles   []string
	ConfigFragments []string
	HostHints       []string
	BrowserHints    []string
	ContextHints    []string
}

type ExecutionInput struct {
	Task             string
	MainPack         string
	PhaseID          string
	PlanID           string
	TriggerKind      string
	TriggerReason    string
	ContractBundle   *model.ContextBundle
	Artifacts        []model.Artifact
	RequiredPacks    []string
	BoundedContext   BoundedContextSnapshot
	RequestedHandoff bool
	Handoff          map[string]any
}

type ValidatorFunc func(plan model.ValidationPlan, input ExecutionInput) model.ValidatorResult
