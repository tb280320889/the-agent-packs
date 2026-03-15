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
			return activationResult(req, nil, "insufficient evidence for L1, fallback to L0 recommended", "partial"), nil
		}
		return activationResult(req, nil, "no route candidate", "failed"), nil
	}

	mainNode := routeResult.Candidates[0].ID
	_, err = query.BuildContextBundle(db, mainNode, true, false, false)
	if err != nil {
		return model.ActivationResult{}, err
	}

	if len(req.BoundedContext.SelectedFiles) == 0 && len(req.BoundedContext.ConfigFragments) == 0 {
		return activationResult(req, "wxt-manifest", "bounded context missing required evidence", "partial"), nil
	}
	if strings.Contains(strings.ToLower(req.Task), "handoff") {
		return activationResult(req, "wxt-manifest", "handoff requested by task boundary", "handoff"), nil
	}
	return activationResult(req, "wxt-manifest", "route to L1.wxt.manifest with required cross-cutting lines", "completed"), nil
}
