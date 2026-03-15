package compiler

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"the-agent-packs/internal/model"
)

var requiredKeys = []string{
	"id",
	"level",
	"domain",
	"subdomain",
	"capability",
	"title",
	"summary",
	"aliases",
	"triggers",
	"anti_triggers",
	"required_with",
	"may_include",
	"children",
	"entry_conditions",
	"stop_conditions",
}

func parseFrontmatter(text string) (map[string]any, string, error) {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 || strings.TrimSpace(strings.TrimRight(lines[0], "\r")) != "---" {
		return nil, "", errors.New("frontmatter missing or not starting with ---")
	}

	end := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(strings.TrimRight(lines[i], "\r")) == "---" {
			end = i
			break
		}
	}
	if end == -1 {
		return nil, "", errors.New("frontmatter not closed with ---")
	}

	fm := map[string]any{}
	currentKey := ""
	for _, raw := range lines[1:end] {
		line := strings.TrimRight(raw, "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(strings.TrimLeft(line, " \t"), "-") && currentKey != "" {
			arr, ok := fm[currentKey].([]string)
			if !ok {
				return nil, "", fmt.Errorf("invalid list field: %s", currentKey)
			}
			item := strings.TrimSpace(strings.TrimPrefix(strings.TrimLeft(line, " \t"), "-"))
			fm[currentKey] = append(arr, item)
			continue
		}

		k, v, found := strings.Cut(line, ":")
		if !found {
			return nil, "", fmt.Errorf("invalid frontmatter line: %s", line)
		}
		key := strings.TrimSpace(k)
		val := strings.TrimSpace(v)
		currentKey = key
		switch val {
		case "", "[]":
			fm[key] = []string{}
		case "null":
			fm[key] = nil
		default:
			fm[key] = val
		}
	}

	body := strings.TrimSpace(strings.Join(lines[end+1:], "\n"))
	return fm, body, nil
}

func deriveID(level, domain string, subdomain *string, stem string) string {
	if level == "L0" {
		return level + "." + domain
	}
	if level == "L1" {
		return level + "." + domain + "." + stem
	}
	if subdomain == nil {
		return level + "." + domain + "." + stem
	}
	return level + "." + domain + "." + *subdomain + "." + stem
}

func computeParentID(level, domain string, subdomain *string) *string {
	if level == "L1" {
		v := "L0." + domain
		return &v
	}
	if level == "L2" && subdomain != nil {
		v := "L1." + domain + "." + *subdomain
		return &v
	}
	return nil
}

func checksumText(text string) string {
	h := sha256.Sum256([]byte(text))
	return hex.EncodeToString(h[:])
}

func ensureSchema(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS nodes (
			id TEXT PRIMARY KEY,
			level TEXT,
			domain TEXT,
			subdomain TEXT,
			capability TEXT,
			title TEXT,
			summary TEXT,
			path TEXT,
			parent_id TEXT,
			body_md TEXT,
			entry_conditions_json TEXT,
			stop_conditions_json TEXT,
			checksum TEXT,
			updated_at TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS node_meta (
			node_id TEXT PRIMARY KEY,
			aliases TEXT,
			triggers TEXT,
			anti_triggers TEXT,
			tags TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS edges (
			source_id TEXT,
			target_id TEXT,
			edge_type TEXT
		)`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	return nil
}

func loadBlueprintFiles(rootDir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

func asString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func asStringPtr(v any) *string {
	if v == nil {
		return nil
	}
	s := asString(v)
	if s == "" {
		return nil
	}
	return &s
}

func asStringSlice(v any) []string {
	if v == nil {
		return []string{}
	}
	if xs, ok := v.([]string); ok {
		return xs
	}
	if s, ok := v.(string); ok {
		if s == "" {
			return []string{}
		}
		return []string{s}
	}
	return []string{}
}

func validateAndCollect(rootDir string) ([]model.Node, []model.NodeMeta, []model.Edge, []map[string]string, error) {
	files, err := loadBlueprintFiles(rootDir)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	nodes := make([]model.Node, 0)
	metas := make([]model.NodeMeta, 0)
	edges := make([]model.Edge, 0)
	errs := make([]map[string]string, 0)

	for _, p := range files {
		base := filepath.Base(p)
		if base == "README.md" || base == "frontmatter-examples.md" || base == "schema.md" {
			continue
		}

		rel, err := filepath.Rel(rootDir, p)
		if err != nil {
			errs = append(errs, map[string]string{"path": p, "error": err.Error()})
			continue
		}
		parts := strings.Split(rel, string(filepath.Separator))
		if len(parts) < 3 {
			errs = append(errs, map[string]string{"path": p, "error": "invalid path structure"})
			continue
		}

		levelDir := parts[0]
		domainDir := parts[1]
		stem := strings.TrimSuffix(parts[len(parts)-1], filepath.Ext(parts[len(parts)-1]))

		var subdomainDir *string
		if levelDir == "L2" || levelDir == "L3" {
			v := parts[2]
			subdomainDir = &v
		}
		if levelDir == "L1" {
			v := stem
			subdomainDir = &v
		}

		raw, err := os.ReadFile(p)
		if err != nil {
			errs = append(errs, map[string]string{"path": p, "error": err.Error()})
			continue
		}
		content := string(raw)
		fm, body, err := parseFrontmatter(content)
		if err != nil {
			errs = append(errs, map[string]string{"path": p, "error": err.Error()})
			continue
		}

		missing := make([]string, 0)
		for _, k := range requiredKeys {
			if _, ok := fm[k]; !ok {
				missing = append(missing, k)
			}
		}
		if len(missing) > 0 {
			errs = append(errs, map[string]string{
				"path":  p,
				"error": fmt.Sprintf("missing keys: %v", missing),
			})
			continue
		}

		derivedID := deriveID(levelDir, domainDir, subdomainDir, stem)
		if asString(fm["id"]) != derivedID {
			errs = append(errs, map[string]string{"path": p, "error": fmt.Sprintf("id mismatch: %s != %s", asString(fm["id"]), derivedID)})
			continue
		}
		if asString(fm["level"]) != levelDir {
			errs = append(errs, map[string]string{"path": p, "error": "level mismatch"})
			continue
		}
		if asString(fm["domain"]) != domainDir {
			errs = append(errs, map[string]string{"path": p, "error": "domain mismatch"})
			continue
		}
		if levelDir != "L0" {
			if subdomainDir == nil || asString(fm["subdomain"]) != *subdomainDir {
				errs = append(errs, map[string]string{"path": p, "error": "subdomain mismatch"})
				continue
			}
		}

		entryConditionsJSON, _ := json.Marshal(asStringSlice(fm["entry_conditions"]))
		stopConditionsJSON, _ := json.Marshal(asStringSlice(fm["stop_conditions"]))
		aliasesJSON, _ := json.Marshal(asStringSlice(fm["aliases"]))
		triggersJSON, _ := json.Marshal(asStringSlice(fm["triggers"]))
		antiTriggersJSON, _ := json.Marshal(asStringSlice(fm["anti_triggers"]))
		tagsJSON, _ := json.Marshal([]string{})

		nodes = append(nodes, model.Node{
			ID:                  asString(fm["id"]),
			Level:               asString(fm["level"]),
			Domain:              asString(fm["domain"]),
			Subdomain:           asStringPtr(fm["subdomain"]),
			Capability:          asStringPtr(fm["capability"]),
			Title:               asString(fm["title"]),
			Summary:             asString(fm["summary"]),
			Path:                filepath.ToSlash(filepath.Join("blueprint", rel)),
			ParentID:            computeParentID(levelDir, domainDir, subdomainDir),
			BodyMD:              body,
			EntryConditionsJSON: string(entryConditionsJSON),
			StopConditionsJSON:  string(stopConditionsJSON),
			Checksum:            checksumText(content),
			UpdatedAt:           time.Now().UTC().Format(time.RFC3339Nano),
		})

		metas = append(metas, model.NodeMeta{
			NodeID:       asString(fm["id"]),
			Aliases:      string(aliasesJSON),
			Triggers:     string(triggersJSON),
			AntiTriggers: string(antiTriggersJSON),
			Tags:         string(tagsJSON),
		})

		for _, t := range asStringSlice(fm["children"]) {
			edges = append(edges, model.Edge{SourceID: asString(fm["id"]), TargetID: t, EdgeType: "child"})
		}
		for _, t := range asStringSlice(fm["required_with"]) {
			edges = append(edges, model.Edge{SourceID: asString(fm["id"]), TargetID: t, EdgeType: "required_with"})
		}
		for _, t := range asStringSlice(fm["may_include"]) {
			edges = append(edges, model.Edge{SourceID: asString(fm["id"]), TargetID: t, EdgeType: "may_include"})
		}
	}

	return nodes, metas, edges, errs, nil
}

func writeIndex(dbPath string, nodes []model.Node, metas []model.NodeMeta, edges []model.Edge) error {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return err
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := ensureSchema(db); err != nil {
		return err
	}

	if _, err := db.Exec("DELETE FROM nodes"); err != nil {
		return err
	}
	if _, err := db.Exec("DELETE FROM node_meta"); err != nil {
		return err
	}
	if _, err := db.Exec("DELETE FROM edges"); err != nil {
		return err
	}

	for _, n := range nodes {
		_, err := db.Exec(
			`INSERT INTO nodes (id, level, domain, subdomain, capability, title, summary, path, parent_id, body_md,
			entry_conditions_json, stop_conditions_json, checksum, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			n.ID,
			n.Level,
			n.Domain,
			n.Subdomain,
			n.Capability,
			n.Title,
			n.Summary,
			n.Path,
			n.ParentID,
			n.BodyMD,
			n.EntryConditionsJSON,
			n.StopConditionsJSON,
			n.Checksum,
			n.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	for _, m := range metas {
		_, err := db.Exec(
			`INSERT INTO node_meta (node_id, aliases, triggers, anti_triggers, tags) VALUES (?, ?, ?, ?, ?)`,
			m.NodeID,
			m.Aliases,
			m.Triggers,
			m.AntiTriggers,
			m.Tags,
		)
		if err != nil {
			return err
		}
	}

	for _, e := range edges {
		_, err := db.Exec(
			`INSERT INTO edges (source_id, target_id, edge_type) VALUES (?, ?, ?)`,
			e.SourceID,
			e.TargetID,
			e.EdgeType,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeReports(reportDir string, errs []map[string]string, edges []model.Edge, nodes []model.Node) error {
	if err := os.MkdirAll(reportDir, 0o755); err != nil {
		return err
	}

	nodeIDs := map[string]bool{}
	for _, n := range nodes {
		nodeIDs[n.ID] = true
	}
	missing := make([]model.Edge, 0)
	for _, e := range edges {
		if !nodeIDs[e.TargetID] {
			missing = append(missing, e)
		}
	}

	validationRaw, _ := json.MarshalIndent(map[string]any{"errors": errs}, "", "  ")
	missingRaw, _ := json.MarshalIndent(map[string]any{"missing": missing}, "", "  ")

	if err := os.WriteFile(filepath.Join(reportDir, "validation-report.json"), validationRaw, 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(reportDir, "missing-reference-report.json"), missingRaw, 0o644); err != nil {
		return err
	}
	return nil
}

func Compile(rootDir, dbPath, reportDir string) ([]map[string]string, error) {
	nodes, metas, edges, errs, err := validateAndCollect(rootDir)
	if err != nil {
		return nil, err
	}
	if err := writeIndex(dbPath, nodes, metas, edges); err != nil {
		return nil, err
	}
	if err := writeReports(reportDir, errs, edges, nodes); err != nil {
		return nil, err
	}
	return errs, nil
}
