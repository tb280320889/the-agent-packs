package tests

import (
	"os"
	"path/filepath"
	"testing"

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
	if err := os.WriteFile(oldDB, []byte("old-index"), 0o644); err != nil {
		t.Fatalf("seed old index failed: %v", err)
	}
	_, err := os.Stat(oldDB)
	if err != nil {
		t.Fatalf("stat old index failed: %v", err)
	}

	if _, err := compiler.Compile(workDir, oldDB, oldDB); err == nil {
		t.Fatalf("expected compile to fail when report write fails")
	}

	_, err = os.Stat(oldDB)
	if err != nil {
		t.Fatalf("stat old index after failure failed: %v", err)
	}
	data, readErr := os.ReadFile(oldDB)
	if readErr != nil {
		t.Fatalf("expected old index to remain intact")
	}
	if len(data) == 0 {
		t.Fatalf("expected old index to remain intact")
	}
	if _, err := os.Stat(oldDB + ".bak"); err == nil {
		t.Fatalf("expected no backup file after report failure")
	}
	if _, err := os.Stat(oldDB + ".tmp"); err == nil {
		t.Fatalf("expected no temp file after report failure")
	}
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
