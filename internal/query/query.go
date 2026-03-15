package query

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "modernc.org/sqlite"

	"the-agent-packs/internal/model"
	"the-agent-packs/internal/registry"
)

type packProfile struct {
	PackName              string
	RecommendedValidators []string
	RecommendedArtifacts  []string
}

type nodeRecord struct {
	ID              string
	Level           string
	Domain          string
	Subdomain       string
	NodeKind        string
	VisibilityScope string
	ActivationMode  string
	Title           string
	Summary         string
}

func PackForNode(nodeID string) string {
	profile, ok := profileForNode(nodeID)
	if !ok {
		return ""
	}
	return profile.PackName
}

func profileForNode(nodeID string) (packProfile, bool) {
	reg, err := registry.Default()
	if err != nil {
		return packProfile{}, false
	}
	entry, ok := registry.FindByNode(reg, nodeID)
	if !ok {
		return packProfile{}, false
	}
	return packProfile{
		PackName:              entry.Name,
		RecommendedValidators: append([]string{}, entry.RecommendedValidators...),
		RecommendedArtifacts:  append([]string{}, entry.RecommendedArtifacts...),
	}, true
}

func OpenDB(dbPath string) (*sql.DB, error) {
	if _, err := os.Stat(dbPath); err != nil {
		return nil, err
	}
	return sql.Open("sqlite", dbPath)
}

func fetchNodes(db *sql.DB, level *string) ([]nodeRecord, error) {
	rows := &sql.Rows{}
	var err error
	if level != nil {
		rows, err = db.Query("SELECT id, level, domain, subdomain, node_kind, visibility_scope, activation_mode, title, summary FROM nodes WHERE level = ?", *level)
	} else {
		rows, err = db.Query("SELECT id, level, domain, subdomain, node_kind, visibility_scope, activation_mode, title, summary FROM nodes")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]nodeRecord, 0)
	for rows.Next() {
		var id, lv, domain, nodeKind, visibilityScope, activationMode string
		var subdomain sql.NullString
		var title, summary string
		if err := rows.Scan(&id, &lv, &domain, &subdomain, &nodeKind, &visibilityScope, &activationMode, &title, &summary); err != nil {
			return nil, err
		}
		sub := ""
		if subdomain.Valid {
			sub = subdomain.String
		}
		result = append(result, nodeRecord{
			ID:              id,
			Level:           lv,
			Domain:          domain,
			Subdomain:       sub,
			NodeKind:        nodeKind,
			VisibilityScope: visibilityScope,
			ActivationMode:  activationMode,
			Title:           title,
			Summary:         summary,
		})
	}
	return result, nil
}

func fetchMeta(db *sql.DB, nodeID string) (map[string][]string, error) {
	row := db.QueryRow("SELECT aliases, triggers, anti_triggers FROM node_meta WHERE node_id = ?", nodeID)
	var aliasesRaw, triggersRaw, antiRaw sql.NullString
	if err := row.Scan(&aliasesRaw, &triggersRaw, &antiRaw); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return map[string][]string{"aliases": {}, "triggers": {}, "anti_triggers": {}}, nil
		}
		return nil, err
	}
	decode := func(v sql.NullString) []string {
		if !v.Valid || strings.TrimSpace(v.String) == "" {
			return []string{}
		}
		arr := []string{}
		_ = json.Unmarshal([]byte(v.String), &arr)
		return arr
	}
	return map[string][]string{
		"aliases":       decode(aliasesRaw),
		"triggers":      decode(triggersRaw),
		"anti_triggers": decode(antiRaw),
	}, nil
}

func fetchEdges(db *sql.DB, sourceID, edgeType string) ([]string, error) {
	rows, err := db.Query("SELECT target_id FROM edges WHERE source_id = ? AND edge_type = ?", sourceID, edgeType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]string, 0)
	for rows.Next() {
		var target string
		if err := rows.Scan(&target); err != nil {
			return nil, err
		}
		list = append(list, target)
	}
	return list, nil
}

func fetchNodeRecord(db *sql.DB, nodeID string) (nodeRecord, error) {
	row := db.QueryRow("SELECT id, level, domain, subdomain, node_kind, visibility_scope, activation_mode, title, summary FROM nodes WHERE id = ?", nodeID)
	var rec nodeRecord
	var subdomain sql.NullString
	if err := row.Scan(&rec.ID, &rec.Level, &rec.Domain, &subdomain, &rec.NodeKind, &rec.VisibilityScope, &rec.ActivationMode, &rec.Title, &rec.Summary); err != nil {
		return nodeRecord{}, err
	}
	if subdomain.Valid {
		rec.Subdomain = subdomain.String
	}
	return rec, nil
}

func tokenize(s string) string {
	return strings.ToLower(s)
}

func listTokenScore(taskText string, values []string, weight float64, reasonPrefix string) (float64, []string) {
	score := 0.0
	reason := make([]string, 0)
	for _, value := range values {
		token := strings.ToLower(value)
		if token != "" && strings.Contains(taskText, token) {
			score += weight
			reason = append(reason, reasonPrefix+":"+value)
		}
	}
	return score, reason
}

func evidenceScore(taskText string, selectedFiles, configFragments, contextHints []string) (float64, []string) {
	score := 0.0
	reason := make([]string, 0)
	for _, filePath := range selectedFiles {
		fileToken := strings.ToLower(filepath.Base(filePath))
		if fileToken != "" && strings.Contains(taskText, fileToken) {
			score += 0.4
			reason = append(reason, "selected_file:"+fileToken)
		}
	}
	for _, fragment := range configFragments {
		token := strings.ToLower(fragment)
		if token != "" && strings.Contains(taskText, token) {
			score += 0.3
			reason = append(reason, "config_fragment:"+fragment)
		}
	}
	for _, hint := range contextHints {
		token := strings.ToLower(hint)
		if token != "" && strings.Contains(taskText, token) {
			score += 0.2
			reason = append(reason, "context_hint:"+hint)
		}
	}
	return score, reason
}

func scoreCandidate(task string, candidate nodeRecord, meta map[string][]string, targetDomain *string, selectedFiles, configFragments, contextHints []string) (float64, []string) {
	score := 0.0
	reason := make([]string, 0)
	taskText := tokenize(task)
	if targetDomain != nil && candidate.Domain == *targetDomain {
		score += 3.0
		reason = append(reason, "target_domain match")
	}
	triggerScore, triggerReason := listTokenScore(taskText, meta["triggers"], 1.0, "trigger")
	aliasScore, aliasReason := listTokenScore(taskText, meta["aliases"], 0.5, "alias")
	evidence, evidenceReason := evidenceScore(taskText, selectedFiles, configFragments, contextHints)
	score += triggerScore + aliasScore + evidence
	reason = append(reason, triggerReason...)
	reason = append(reason, aliasReason...)
	reason = append(reason, evidenceReason...)
	return score, reason
}

func domainNodeAllowedInGlobal(candidate nodeRecord) bool {
	if candidate.ActivationMode == "attach-only" {
		return false
	}
	return candidate.NodeKind == "domain-root" || candidate.NodeKind == "domain-orchestrator"
}

func workflowNodeAllowedInDomain(candidate nodeRecord, activeDomain string) bool {
	if candidate.Domain != activeDomain {
		return false
	}
	if candidate.ActivationMode == "attach-only" {
		return false
	}
	if candidate.VisibilityScope == "domain-scoped" || candidate.VisibilityScope == "global" {
		return candidate.NodeKind == "workflow-entry" || candidate.NodeKind == "domain-orchestrator"
	}
	return false
}

func capabilityAttachAllowed(candidate nodeRecord, activeDomain string) bool {
	if activeDomain == "" {
		return false
	}
	if candidate.ActivationMode != "attach-only" {
		return false
	}
	return candidate.VisibilityScope == "capability-scoped" || candidate.VisibilityScope == "domain-scoped"
}

func inferMainDomain(task string, selectedFiles, configFragments, contextHints []string) string {
	taskText := tokenize(task)
	joined := strings.Join(append(append([]string{}, selectedFiles...), append(configFragments, contextHints...)...), " ")
	evidence := tokenize(joined)
	if strings.Contains(taskText, "wxt") || strings.Contains(taskText, "browser extension") || strings.Contains(evidence, "wxt") || strings.Contains(evidence, "manifest") {
		return "wxt"
	}
	return ""
}

func buildCandidate(candidate nodeRecord, score float64, reason []string) model.RouteCandidate {
	return model.RouteCandidate{
		ID:              candidate.ID,
		Title:           candidate.Title,
		Summary:         candidate.Summary,
		Score:           score,
		Reason:          reason,
		NodeKind:        candidate.NodeKind,
		VisibilityScope: candidate.VisibilityScope,
		ActivationMode:  candidate.ActivationMode,
	}
}

func RouteQuery(db *sql.DB, level string, task string, targetPack *string, targetDomain *string, selectedFiles, configFragments, contextHints []string, maxResults int) (model.RouteResult, error) {
	if targetPack != nil {
		reg, regErr := registry.Default()
		if regErr == nil {
			if entry, ok := registry.FindByName(reg, *targetPack); ok {
				mappedNode := entry.CanonicalBlueprintNode
				directRec, recErr := fetchNodeRecord(db, mappedNode)
				direct, err := ReadNode(db, mappedNode, "summary")
				if err == nil && recErr == nil && direct != nil && direct.Level == level {
					must, _ := fetchEdges(db, direct.ID, "required_with")
					reason := []string{"target_pack match"}
					if targetDomain != nil {
						reason = append(reason, fmt.Sprintf("target_domain=%s bypassed by explicit target_pack", *targetDomain))
					}
					return model.RouteResult{
						Candidates: []model.RouteCandidate{{
							ID:              direct.ID,
							Title:           direct.Title,
							Summary:         direct.Summary,
							Score:           99.0,
							Reason:          reason,
							NodeKind:        directRec.NodeKind,
							VisibilityScope: directRec.VisibilityScope,
							ActivationMode:  directRec.ActivationMode,
						}},
						MustInclude: must,
					}, nil
				}
			}
		}
	}

	levelArg := level
	nodes, err := fetchNodes(db, &levelArg)
	if err != nil {
		return model.RouteResult{}, err
	}

	activeDomain := ""
	if targetDomain != nil {
		activeDomain = *targetDomain
	} else {
		activeDomain = inferMainDomain(task, selectedFiles, configFragments, contextHints)
	}

	globalCandidates := make([]model.RouteCandidate, 0)
	workflowCandidates := make([]model.RouteCandidate, 0)
	attachIDs := make([]string, 0)

	for _, row := range nodes {
		meta, err := fetchMeta(db, row.ID)
		if err != nil {
			return model.RouteResult{}, err
		}
		taskText := tokenize(task)
		excluded := false
		for _, anti := range meta["anti_triggers"] {
			if strings.Contains(taskText, strings.ToLower(anti)) {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}
		score, reason := scoreCandidate(task, row, meta, targetDomain, selectedFiles, configFragments, contextHints)
		if score <= 0 {
			continue
		}

		if level == "L0" && domainNodeAllowedInGlobal(row) {
			reason = append(reason, "global_candidate_space")
			globalCandidates = append(globalCandidates, buildCandidate(row, score, reason))
			continue
		}
		if level == "L1" && workflowNodeAllowedInDomain(row, activeDomain) {
			reason = append(reason, "domain_candidate_space")
			workflowCandidates = append(workflowCandidates, buildCandidate(row, score, reason))
			continue
		}
		if level == "L1" && capabilityAttachAllowed(row, activeDomain) {
			attachIDs = append(attachIDs, row.ID)
		}
	}

	sort.Slice(globalCandidates, func(i, j int) bool {
		if globalCandidates[i].Score == globalCandidates[j].Score {
			return globalCandidates[i].ID < globalCandidates[j].ID
		}
		return globalCandidates[i].Score > globalCandidates[j].Score
	})
	sort.Slice(workflowCandidates, func(i, j int) bool {
		if workflowCandidates[i].Score == workflowCandidates[j].Score {
			return workflowCandidates[i].ID < workflowCandidates[j].ID
		}
		return workflowCandidates[i].Score > workflowCandidates[j].Score
	})

	if maxResults <= 0 {
		maxResults = 3
	}

	selected := workflowCandidates
	if level == "L0" {
		selected = globalCandidates
	}
	if len(selected) > maxResults {
		selected = selected[:maxResults]
	}

	must := []string{}
	if len(selected) > 0 {
		must, _ = fetchEdges(db, selected[0].ID, "required_with")
		seen := map[string]bool{}
		for _, id := range must {
			seen[id] = true
		}
		for _, id := range attachIDs {
			if !seen[id] {
				must = append(must, id)
			}
		}
	}

	return model.RouteResult{Candidates: selected, MustInclude: must}, nil
}

func ReadNode(db *sql.DB, nodeID string, section string) (*model.NodeSummary, error) {
	row := db.QueryRow("SELECT id, title, summary, body_md, level FROM nodes WHERE id = ?", nodeID)
	var id, title, summary, body, level string
	if err := row.Scan(&id, &title, &summary, &body, &level); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if section == "summary" {
		return &model.NodeSummary{ID: id, Title: title, Summary: summary, Level: level}, nil
	}
	return &model.NodeSummary{ID: id, Title: title, Body: body, Level: level}, nil
}

func BuildContextBundle(db *sql.DB, mainNode string, includeRequired, includeMayInclude, includeChildren bool) (model.ContextBundle, error) {
	bundle := model.ContextBundle{
		Main:                  nil,
		Required:              []model.NodeSummary{},
		ExecutionChildren:     []model.NodeSummary{},
		Deferred:              []model.NodeSummary{},
		RecommendedValidators: []string{},
		RecommendedArtifacts:  []string{},
	}

	main, err := ReadNode(db, mainNode, "summary")
	if err != nil {
		return bundle, err
	}
	if main == nil {
		return bundle, nil
	}
	bundle.Main = main
	if profile, ok := profileForNode(mainNode); ok {
		bundle.RecommendedValidators = append(bundle.RecommendedValidators, profile.RecommendedValidators...)
		bundle.RecommendedArtifacts = append(bundle.RecommendedArtifacts, profile.RecommendedArtifacts...)
	}

	appendNode := func(targetID string, collection *[]model.NodeSummary) error {
		node, err := ReadNode(db, targetID, "summary")
		if err != nil {
			return err
		}
		if node == nil {
			return nil
		}
		if node.Level == "L3" {
			bundle.Deferred = append(bundle.Deferred, *node)
			return nil
		}
		*collection = append(*collection, *node)
		return nil
	}

	if includeRequired {
		targets, err := fetchEdges(db, mainNode, "required_with")
		if err != nil {
			return bundle, err
		}
		for _, t := range targets {
			if err := appendNode(t, &bundle.Required); err != nil {
				return bundle, err
			}
		}
	}

	if includeChildren {
		targets, err := fetchEdges(db, mainNode, "child")
		if err != nil {
			return bundle, err
		}
		for _, t := range targets {
			if err := appendNode(t, &bundle.ExecutionChildren); err != nil {
				return bundle, err
			}
		}
	}

	if includeMayInclude {
		targets, err := fetchEdges(db, mainNode, "may_include")
		if err != nil {
			return bundle, err
		}
		for _, t := range targets {
			if err := appendNode(t, &bundle.ExecutionChildren); err != nil {
				return bundle, err
			}
		}
	}

	return bundle, nil
}

func ExpandNode(db *sql.DB, nodeID, edgeType string) ([]string, error) {
	return fetchEdges(db, nodeID, edgeType)
}
