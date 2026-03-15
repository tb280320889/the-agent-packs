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
	Artifacts        []model.Artifact
	BoundedContext   BoundedContextSnapshot
	RequestedHandoff bool
	Handoff          map[string]any
}

type ValidatorFunc func(plan model.ValidationPlan, input ExecutionInput) model.ValidatorResult
