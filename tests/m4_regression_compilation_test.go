package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"the-agent-packs/internal/compiler"
	"the-agent-packs/internal/registry"
)

func TestM4RegressionParsing(t *testing.T) {
	t.Run("frontmatter fixture parses", func(t *testing.T) {
		root := projectRoot(t)
		fixturePath := filepath.Join(root, "fixtures", "blueprint", "frontmatter-multi-line.md")
		raw, err := os.ReadFile(fixturePath)
		if err != nil {
			t.Fatalf("read frontmatter fixture failed: %v", err)
		}

		base := t.TempDir()
		blueprintRoot := filepath.Join(base, "blueprint")
		targetPath := filepath.Join(blueprintRoot, "L1", "demo", "multiline.md")
		if err := writeBlueprintFixture(targetPath, string(raw)); err != nil {
			t.Fatalf("write blueprint fixture failed: %v", err)
		}

		result, err := compiler.Compile(blueprintRoot, filepath.Join(base, "index", "blueprint.db"), filepath.Join(base, "index"))
		if err != nil {
			t.Fatalf("compile failed: %v", err)
		}
		if len(result.Errors) != 0 {
			t.Fatalf("expected no frontmatter errors, got %+v", result.Errors)
		}
	})

	t.Run("package unknown field fails", func(t *testing.T) {
		root := projectRoot(t)
		fixturePath := filepath.Join(root, "fixtures", "registry", "package-with-unknown-field.yaml")
		raw, err := os.ReadFile(fixturePath)
		if err != nil {
			t.Fatalf("read package fixture failed: %v", err)
		}

		reg := registry.Registry{
			ReservedNames: []string{},
			Packages: []registry.PackageEntry{
				{
					Name:                   "demo-unknown",
					Kind:                   "workflow-package",
					Domain:                 "demo",
					Subdomain:              "unknown",
					Category:               "workflow",
					CanonicalBlueprintNode: "L1.demo.unknown",
					VisibilityScope:        "domain-scoped",
					ActivationMode:         "direct",
					Aliases:                []string{},
					RecommendedValidators:  []string{"validator-core-output"},
					RecommendedArtifacts:   []string{"demo.md"},
					RequiredPacks:          []string{},
				},
			},
		}

		base := t.TempDir()
		regPath := filepath.Join(base, "workflow-packages", "registry.json")
		if err := writeRegistryFixture(regPath, reg); err != nil {
			t.Fatalf("write registry fixture failed: %v", err)
		}

		pkgDir := filepath.Join(base, "workflow-packages", "demo-unknown")
		if err := writePackageFixture(pkgDir, string(raw)); err != nil {
			t.Fatalf("write package fixture failed: %v", err)
		}

		err = registry.Validate(reg, filepath.Join(base, "workflow-packages"))
		if err == nil {
			t.Fatalf("expected unknown field to fail strict YAML decode")
		}
		if !strings.Contains(err.Error(), "unknown_field") {
			t.Fatalf("expected error to include unknown field name, got %v", err)
		}
	})
}
