package tests

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
	"the-agent-packs/internal/compiler"
)

func TestM4IndexRebuildTransactional(t *testing.T) {
	root := projectRoot(t)
	base := t.TempDir()
	workDir := filepath.Join(base, "blueprint")
	if err := copyDir(filepath.Join(root, "blueprint"), workDir); err != nil {
		t.Fatalf("copy blueprint failed: %v", err)
	}

	indexDir := filepath.Join(workDir, "index")
	if err := os.MkdirAll(indexDir, 0o755); err != nil {
		t.Fatalf("mkdir index failed: %v", err)
	}
	oldDB := filepath.Join(indexDir, "blueprint.db")
	seedReportDir := filepath.Join(base, "seed-report")
	if _, err := compiler.Compile(workDir, oldDB, seedReportDir); err != nil {
		t.Fatalf("seed old index failed: %v", err)
	}
	seedTitle, err := queryNodeTitle(oldDB, "L0.wxt")
	if err != nil {
		t.Fatalf("query seed index failed: %v", err)
	}
	if seedTitle == "" {
		t.Fatalf("expected seeded index to be queryable")
	}

	if err := mutateNodeTitle(filepath.Join(workDir, "L0", "wxt", "overview.md"), "WXT Mutated Title"); err != nil {
		t.Fatalf("mutate blueprint failed: %v", err)
	}

	if _, err := os.Stat(oldDB); err != nil {
		t.Fatalf("stat old index failed: %v", err)
	}

	if _, err := compiler.Compile(workDir, oldDB, oldDB); err == nil {
		t.Fatalf("expected compile to fail when report write fails")
	}

	_, err = os.Stat(oldDB)
	if err != nil {
		t.Fatalf("stat old index after failure failed: %v", err)
	}
	afterFailureTitle, queryErr := queryNodeTitle(oldDB, "L0.wxt")
	if queryErr != nil {
		t.Fatalf("expected old index to remain queryable: %v", queryErr)
	}
	if afterFailureTitle != seedTitle {
		t.Fatalf("expected old index title to remain %q, got %q", seedTitle, afterFailureTitle)
	}
	if _, err := os.Stat(oldDB + ".bak"); err == nil {
		t.Fatalf("expected no backup file after report failure")
	}
	if _, err := os.Stat(oldDB + ".tmp"); err == nil {
		t.Fatalf("expected no temp file after report failure")
	}

	successReportDir := filepath.Join(base, "success-report")
	if _, err := compiler.Compile(workDir, oldDB, successReportDir); err != nil {
		t.Fatalf("expected compile success after fixing report dir: %v", err)
	}
	afterSuccessTitle, queryErr := queryNodeTitle(oldDB, "L0.wxt")
	if queryErr != nil {
		t.Fatalf("query replaced index failed: %v", queryErr)
	}
	if afterSuccessTitle != "WXT Mutated Title" {
		t.Fatalf("expected replaced index title %q, got %q", "WXT Mutated Title", afterSuccessTitle)
	}
}

func queryNodeTitle(dbPath, nodeID string) (string, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var title string
	err = db.QueryRow(`SELECT title FROM nodes WHERE id = ?`, nodeID).Scan(&title)
	if err != nil {
		return "", err
	}
	return title, nil
}

func mutateNodeTitle(path, title string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(raw), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "title:") {
			lines[i] = "title: " + title
			break
		}
	}
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode())
	})
}
