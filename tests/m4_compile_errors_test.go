package tests

import (
	"path/filepath"
	"testing"

	"the-agent-packs/internal/compiler"
)

func TestM4CompileErrors(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "blueprint")
	badPath := filepath.Join(root, "L1", "demo", "bad.md")
	content := "---\nid: L1.demo.bad\nlevel: L1\ndomain: demo\nsubdomain: bad\ncapability: null\nnode_kind: workflow-entry\nvisibility_scope: domain-scoped\nactivation_mode: direct\ntitle: Bad\nsummary: missing aliases\naliases: []\ntriggers: []\nanti_triggers: []\nrequired_with: []\nmay_include: []\nchildren: []\nentry_conditions: []\nstop_conditions: []\nunknown_field: nope\n---\nbody\n"
	if err := writeBlueprintFixture(badPath, content); err != nil {
		t.Fatalf("write blueprint fixture failed: %v", err)
	}

	result, err := compiler.Compile(root, filepath.Join(base, "index", "blueprint.db"), filepath.Join(base, "index"))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}
	if len(result.Errors) == 0 {
		t.Fatalf("expected compile errors")
	}
	foundPhase := false
	foundPath := false
	foundCode := false
	for _, entry := range result.Errors {
		if entry.Phase != "" {
			foundPhase = true
		}
		if entry.Path == badPath {
			foundPath = true
		}
		if entry.Code != "" {
			foundCode = true
		}
	}
	if !foundPhase || !foundPath || !foundCode {
		t.Fatalf("expected structured error with phase/path/code, got %+v", result.Errors)
	}
}
