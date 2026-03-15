package validator

import "the-agent-packs/internal/model"

func Run(plan model.ValidationPlan, input ExecutionInput) []model.ValidatorResult {
	all := make([]model.ValidatorResult, 0, len(plan.Validators))
	r := registry()
	for _, validatorPlan := range plan.Validators {
		fn, ok := r[validatorPlan.Name]
		if !ok {
			all = append(all, model.ValidatorResult{
				ValidatorName:      validatorPlan.Name,
				Status:             "skipped",
				Findings:           []model.Finding{},
				RepairSuggestions:  []string{"Validator is not registered."},
				ValidatedArtifacts: plan.ArtifactsUnderValidation,
			})
			continue
		}
		result := fn(plan, input)
		if result.ValidatorName == "" {
			result.ValidatorName = validatorPlan.Name
		}
		if result.Status == "" {
			result.Status = "passed"
		}
		all = append(all, result)
	}
	return all
}
