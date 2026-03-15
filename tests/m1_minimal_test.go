package tests

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"the-agent-packs/internal/activation"
	"the-agent-packs/internal/compiler"
	"the-agent-packs/internal/query"
)

func projectRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	return filepath.Clean(filepath.Join(wd, ".."))
}

func compileMainIndex(t *testing.T) string {
	t.Helper()
	root := projectRoot(t)
	dbPath := filepath.Join(root, "blueprint", "index", "blueprint.db")
	reportDir := filepath.Join(root, "blueprint", "index")
	errList, err := compiler.Compile(filepath.Join(root, "blueprint"), dbPath, reportDir)
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}
	if len(errList) > 0 {
		t.Fatalf("compile has errors: %+v", errList)
	}
	return dbPath
}

func openDB(t *testing.T, dbPath string) *sql.DB {
	t.Helper()
	db, err := query.OpenDB(dbPath)
	if err != nil {
		t.Fatalf("open db failed: %v", err)
	}
	return db
}

func TestCompilerReportsNoErrors(t *testing.T) {
	root := projectRoot(t)
	_ = compileMainIndex(t)

	validationRaw, err := os.ReadFile(filepath.Join(root, "blueprint", "index", "validation-report.json"))
	if err != nil {
		t.Fatalf("read validation report failed: %v", err)
	}
	missingRaw, err := os.ReadFile(filepath.Join(root, "blueprint", "index", "missing-reference-report.json"))
	if err != nil {
		t.Fatalf("read missing report failed: %v", err)
	}

	var validation map[string]any
	if err := json.Unmarshal(validationRaw, &validation); err != nil {
		t.Fatalf("decode validation report failed: %v", err)
	}
	var missing map[string]any
	if err := json.Unmarshal(missingRaw, &missing); err != nil {
		t.Fatalf("decode missing report failed: %v", err)
	}

	errorsList, _ := validation["errors"].([]any)
	if len(errorsList) != 0 {
		t.Fatalf("expected no validation errors, got %v", errorsList)
	}
	missingList, _ := missing["missing"].([]any)
	if len(missingList) != 0 {
		t.Fatalf("expected no missing edges, got %v", missingList)
	}
}

func TestRouteTargetPackHasHighestPriority(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	tp := "security-permissions"
	result, err := query.RouteQuery(db, "L1", "unrelated task", &tp, nil, nil, nil, nil, 3)
	if err != nil {
		t.Fatalf("route query failed: %v", err)
	}
	if len(result.Candidates) == 0 {
		t.Fatalf("expected candidates")
	}
	if result.Candidates[0].ID != "L1.security.permissions" {
		t.Fatalf("unexpected candidate: %s", result.Candidates[0].ID)
	}
	found := false
	for _, r := range result.Candidates[0].Reason {
		if r == "target_pack match" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("target_pack match reason missing")
	}
}

func TestRouteRespectsTargetDomain(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	td := "wxt"
	result, err := query.RouteQuery(db, "L1", "permissions review for browser extension", nil, &td, nil, nil, nil, 3)
	if err != nil {
		t.Fatalf("route query failed: %v", err)
	}
	if len(result.Candidates) == 0 {
		t.Fatalf("expected candidates")
	}
	for _, c := range result.Candidates {
		segments := strings.Split(c.ID, ".")
		if len(segments) < 2 || segments[1] != "wxt" {
			t.Fatalf("candidate not in wxt domain: %s", c.ID)
		}
	}
}

func TestRouteTargetPackOverridesTargetDomainConflict(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	tp := "security-permissions"
	td := "wxt"
	result, err := query.RouteQuery(db, "L1", "manifest permissions", &tp, &td, nil, nil, nil, 3)
	if err != nil {
		t.Fatalf("route query failed: %v", err)
	}
	if len(result.Candidates) == 0 || result.Candidates[0].ID != "L1.security.permissions" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestRouteAntiTriggerExclusion(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	td := "wxt"
	result, err := query.RouteQuery(db, "L1", "tauri permission setup", nil, &td, nil, nil, nil, 3)
	if err != nil {
		t.Fatalf("route query failed: %v", err)
	}
	if len(result.Candidates) != 0 {
		t.Fatalf("expected no candidates, got %+v", result.Candidates)
	}
}

func TestFrontmatterSummaryWithColon(t *testing.T) {
	tempDir := t.TempDir()
	l0Dir := filepath.Join(tempDir, "L0", "demo")
	if err := os.MkdirAll(l0Dir, 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	filePath := filepath.Join(l0Dir, "overview.md")
	content := "---\nid: L0.demo\nlevel: L0\ndomain: demo\nsubdomain: null\ncapability: null\ntitle: Demo\nsummary: This summary has colon: valid content\naliases: []\ntriggers:\n  - demo\nanti_triggers: []\nrequired_with: []\nmay_include: []\nchildren: []\nentry_conditions: []\nstop_conditions: []\n---\n\ndemo body\n"
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	dbPath := filepath.Join(tempDir, "index", "blueprint.db")
	reportDir := filepath.Join(tempDir, "index")
	errList, err := compiler.Compile(tempDir, dbPath, reportDir)
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}
	if len(errList) != 0 {
		t.Fatalf("expected no errors, got %+v", errList)
	}
}

func TestFrontmatterMissingRequiredKeysDetected(t *testing.T) {
	tempDir := t.TempDir()
	l0Dir := filepath.Join(tempDir, "L0", "demo")
	if err := os.MkdirAll(l0Dir, 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	filePath := filepath.Join(l0Dir, "overview.md")
	content := "---\nid: L0.demo\nlevel: L0\ndomain: demo\nsummary: missing required keys\n---\n\ndemo body\n"
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	dbPath := filepath.Join(tempDir, "index", "blueprint.db")
	reportDir := filepath.Join(tempDir, "index")
	errList, err := compiler.Compile(tempDir, dbPath, reportDir)
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}
	if len(errList) == 0 {
		t.Fatalf("expected missing keys errors")
	}
}

func TestActivationCompletedPath(t *testing.T) {
	root := projectRoot(t)
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	result, err := activation.Execute(db, filepath.Join(root, "fixtures", "activation-request.sample.json"))
	if err != nil {
		t.Fatalf("activation execute failed: %v", err)
	}
	if result.Status != "completed" {
		t.Fatalf("expected completed, got %s", result.Status)
	}
}

func TestActivationHandoffPath(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	request := map[string]any{
		"request_id":    "req-handoff",
		"task":          "handoff to next pack after manifest review",
		"target_domain": "wxt",
		"bounded_context": map[string]any{
			"selected_files":   []string{"wxt.config.ts"},
			"config_fragments": []string{"manifest.permissions"},
		},
	}
	tempFile, err := os.CreateTemp("", "activation-*.json")
	if err != nil {
		t.Fatalf("create temp file failed: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	enc := json.NewEncoder(tempFile)
	if err := enc.Encode(request); err != nil {
		t.Fatalf("write request failed: %v", err)
	}

	result, err := activation.Execute(db, tempFile.Name())
	if err != nil {
		t.Fatalf("activation execute failed: %v", err)
	}
	if result.Status != "handoff" {
		t.Fatalf("expected handoff, got %s", result.Status)
	}
}

func TestActivationPartialWhenDomainKnownButNoCandidate(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	request := map[string]any{
		"request_id":    "req-partial",
		"task":          "totally unrelated text",
		"target_domain": "wxt",
		"bounded_context": map[string]any{
			"selected_files":   []string{},
			"config_fragments": []string{},
		},
	}
	tempFile, err := os.CreateTemp("", "activation-*.json")
	if err != nil {
		t.Fatalf("create temp file failed: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	enc := json.NewEncoder(tempFile)
	if err := enc.Encode(request); err != nil {
		t.Fatalf("write request failed: %v", err)
	}

	result, err := activation.Execute(db, tempFile.Name())
	if err != nil {
		t.Fatalf("activation execute failed: %v", err)
	}
	if result.Status != "partial" {
		t.Fatalf("expected partial, got %s", result.Status)
	}
}

func TestActivationFailedWhenRequestInvalid(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	request := map[string]any{
		"request_id": "req-failed",
		"task":       "manifest",
	}
	tempFile, err := os.CreateTemp("", "activation-*.json")
	if err != nil {
		t.Fatalf("create temp file failed: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	enc := json.NewEncoder(tempFile)
	if err := enc.Encode(request); err != nil {
		t.Fatalf("write request failed: %v", err)
	}

	result, err := activation.Execute(db, tempFile.Name())
	if err != nil {
		t.Fatalf("activation execute failed: %v", err)
	}
	if result.Status != "failed" {
		t.Fatalf("expected failed, got %s", result.Status)
	}
}
