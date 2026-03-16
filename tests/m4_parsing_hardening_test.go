package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"the-agent-packs/internal/compiler"
	"the-agent-packs/internal/registry"
)

func TestM4PackageYamlStrict(t *testing.T) {
	root := projectRoot(t)
	regPath := filepath.Join(root, "workflow-packages", "registry.json")
	reg, err := registry.Load(regPath)
	if err != nil {
		t.Fatalf("load registry failed: %v", err)
	}

	unknownField := "unknown_field_should_fail"

	badRoot := t.TempDir()
	badReg := reg
	badReg.ReservedNames = []string{}
	badReg.Packages = []registry.PackageEntry{
		{
			Name:                   "demo-test",
			Kind:                   "workflow-package",
			Domain:                 "demo",
			Subdomain:              "test",
			Category:               "workflow",
			CanonicalBlueprintNode: "L1.demo.test",
			VisibilityScope:        "domain-scoped",
			ActivationMode:         "direct",
			Aliases:                []string{},
			RecommendedValidators:  []string{},
			RecommendedArtifacts:   []string{},
			RequiredPacks:          []string{},
		},
	}

	badRegPath := filepath.Join(badRoot, "workflow-packages", "registry.json")
	badPkgDir := filepath.Join(badRoot, "workflow-packages", "demo-test")
	if err := writeRegistryFixture(badRegPath, badReg); err != nil {
		t.Fatalf("write registry fixture failed: %v", err)
	}
	if err := writePackageFixture(badPkgDir, strings.Join([]string{
		"name: demo-test",
		"kind: workflow-package",
		"domain: demo",
		"subdomain: test",
		"" + unknownField + ": should-fail",
		"depends_on:",
		"  - security-permissions",
		"validators:",
		"  - validator-core-output",
		"artifacts:",
		"  - demo.md",
	}, "\n")); err != nil {
		t.Fatalf("write package fixture failed: %v", err)
	}
	if err := registry.Validate(badReg, filepath.Join(badRoot, "workflow-packages")); err == nil {
		t.Fatalf("expected unknown field to fail strict YAML decode")
	}
}

func TestM4PackageYamlListParsing(t *testing.T) {
	root := t.TempDir()
	reg := registry.Registry{
		ReservedNames: []string{},
		Packages: []registry.PackageEntry{
			{
				Name:                   "demo-list",
				Kind:                   "workflow-package",
				Domain:                 "demo",
				Subdomain:              "list",
				Category:               "workflow",
				CanonicalBlueprintNode: "L1.demo.list",
				VisibilityScope:        "domain-scoped",
				ActivationMode:         "direct",
				Aliases:                []string{},
				RecommendedValidators:  []string{"validator-core-output"},
				RecommendedArtifacts:   []string{"demo.md"},
				RequiredPacks:          []string{},
			},
		},
	}
	regPath := filepath.Join(root, "workflow-packages", "registry.json")
	if err := writeRegistryFixture(regPath, reg); err != nil {
		t.Fatalf("write registry fixture failed: %v", err)
	}
	if err := writePackageFixture(filepath.Join(root, "workflow-packages", "demo-list"), strings.Join([]string{
		"name: demo-list",
		"kind: workflow-package",
		"domain: demo",
		"subdomain: list",
		"validators:",
		"  - validator-core-output",
		"artifacts:",
		"  - demo.md",
	}, "\n")); err != nil {
		t.Fatalf("write package fixture failed: %v", err)
	}
	if err := registry.Validate(reg, filepath.Join(root, "workflow-packages")); err != nil {
		t.Fatalf("expected list fields to parse: %v", err)
	}
}

func TestM4FrontmatterStrict(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "blueprint")
	filePath := filepath.Join(root, "L1", "demo", "sample.md")
	if err := writeBlueprintFixture(filePath, strings.Join([]string{
		"---",
		"id: L1.demo.sample",
		"level: L1",
		"domain: demo",
		"subdomain: sample",
		"capability: null",
		"node_kind: workflow-entry",
		"visibility_scope: domain-scoped",
		"activation_mode: direct",
		"title: \"Demo Sample\"",
		"summary: |",
		"  multi-line summary",
		"  continues here",
		"aliases:",
		"  - \"demo alias\"",
		"triggers:",
		"  - manifest",
		"anti_triggers: []",
		"required_with: []",
		"may_include: []",
		"children: []",
		"entry_conditions:",
		"  - entry_ok",
		"stop_conditions:",
		"  - stop_ok",
		"---",
		"body",
	}, "\n")); err != nil {
		t.Fatalf("write blueprint fixture failed: %v", err)
	}

	unknownPath := filepath.Join(root, "L1", "demo", "unknown.md")
	if err := writeBlueprintFixture(unknownPath, strings.Join([]string{
		"---",
		"id: L1.demo.unknown",
		"level: L1",
		"domain: demo",
		"subdomain: unknown",
		"capability: null",
		"node_kind: workflow-entry",
		"visibility_scope: domain-scoped",
		"activation_mode: direct",
		"title: Demo Unknown",
		"summary: test",
		"aliases: []",
		"triggers: []",
		"anti_triggers: []",
		"required_with: []",
		"may_include: []",
		"children: []",
		"entry_conditions: []",
		"stop_conditions: []",
		"unknown_field: nope",
		"---",
		"body",
	}, "\n")); err != nil {
		t.Fatalf("write blueprint fixture failed: %v", err)
	}

	errList, err := compiler.Compile(root, filepath.Join(base, "index", "blueprint.db"), filepath.Join(base, "index"))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}
	if len(errList) == 0 {
		t.Fatalf("expected frontmatter unknown field to be reported")
	}
	found := false
	for _, entry := range errList {
		if entry["path"] == unknownPath {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected unknown frontmatter error for %s", unknownPath)
	}
}

func writeRegistryFixture(path string, reg registry.Registry) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(reg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0o644)
}

func writePackageFixture(dir, content string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "package.yaml"), []byte(content+"\n"), 0o644)
}

func writeBlueprintFixture(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content+"\n"), 0o644)
}
