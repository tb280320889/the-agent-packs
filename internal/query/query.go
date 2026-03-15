package query

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "modernc.org/sqlite"

	"the-agent-packs/internal/model"
)

var packNodeMap = map[string]string{
	"wxt-manifest":         "L1.wxt.manifest",
	"security-permissions": "L1.security.permissions",
	"release-store-review": "L1.release.store-review",
}

func OpenDB(dbPath string) (*sql.DB, error) {
	if _, err := os.Stat(dbPath); err != nil {
		return nil, err
	}
	return sql.Open("sqlite", dbPath)
}

func fetchNodes(db *sql.DB, level *string) ([][6]string, error) {
	rows := &sql.Rows{}
	var err error
	if level != nil {
		rows, err = db.Query("SELECT id, level, domain, subdomain, title, summary FROM nodes WHERE level = ?", *level)
	} else {
		rows, err = db.Query("SELECT id, level, domain, subdomain, title, summary FROM nodes")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([][6]string, 0)
	for rows.Next() {
		var id, lv, domain string
		var subdomain sql.NullString
		var title, summary string
		if err := rows.Scan(&id, &lv, &domain, &subdomain, &title, &summary); err != nil {
			return nil, err
		}
		sub := ""
		if subdomain.Valid {
			sub = subdomain.String
		}
		result = append(result, [6]string{id, lv, domain, sub, title, summary})
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

func scoreCandidate(task string, candidate [6]string, meta map[string][]string, targetDomain *string, selectedFiles, configFragments, contextHints []string) (float64, []string) {
	score := 0.0
	reason := make([]string, 0)
	taskText := tokenize(task)
	if targetDomain != nil && candidate[2] == *targetDomain {
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

func RouteQuery(db *sql.DB, level string, task string, targetPack *string, targetDomain *string, selectedFiles, configFragments, contextHints []string, maxResults int) (model.RouteResult, error) {
	if targetPack != nil {
		if mappedNode, ok := packNodeMap[*targetPack]; ok {
			direct, err := ReadNode(db, mappedNode, "summary")
			if err == nil && direct != nil && direct.Level == level {
				must, _ := fetchEdges(db, direct.ID, "required_with")
				return model.RouteResult{
					Candidates: []model.RouteCandidate{{
						ID:      direct.ID,
						Title:   direct.Title,
						Summary: direct.Summary,
						Score:   99.0,
						Reason:  []string{"target_pack match"},
					}},
					MustInclude: must,
				}, nil
			}
		}
	}

	levelArg := level
	nodes, err := fetchNodes(db, &levelArg)
	if err != nil {
		return model.RouteResult{}, err
	}

	candidates := make([]model.RouteCandidate, 0)
	for _, row := range nodes {
		if targetDomain != nil && row[2] != *targetDomain {
			continue
		}
		meta, err := fetchMeta(db, row[0])
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
		if score > 0 {
			candidates = append(candidates, model.RouteCandidate{
				ID:      row[0],
				Title:   row[4],
				Summary: row[5],
				Score:   score,
				Reason:  reason,
			})
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Score == candidates[j].Score {
			return candidates[i].ID < candidates[j].ID
		}
		return candidates[i].Score > candidates[j].Score
	})

	if maxResults <= 0 {
		maxResults = 3
	}
	if len(candidates) > maxResults {
		candidates = candidates[:maxResults]
	}

	must := []string{}
	if len(candidates) > 0 {
		must, _ = fetchEdges(db, candidates[0].ID, "required_with")
	}

	return model.RouteResult{Candidates: candidates, MustInclude: must}, nil
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
