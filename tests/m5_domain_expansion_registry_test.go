package tests

import (
	"path/filepath"
	"testing"

	"the-agent-packs/internal/registry"
)

func TestM5DomainExpansionRegistryOnboarding(t *testing.T) {
	root := projectRoot(t)
	reg, err := registry.Load(filepath.Join(root, "workflow-packages", "registry.json"))
	if err != nil {
		t.Fatalf("load registry failed: %v", err)
	}

	entry, ok := registry.FindByName(reg, "monorepo-oss-governance")
	if !ok {
		t.Fatalf("expected monorepo-oss-governance in registry")
	}
	if entry.CanonicalBlueprintNode != "L1.monorepo.oss-governance" {
		t.Fatalf("unexpected canonical blueprint node: %s", entry.CanonicalBlueprintNode)
	}
	if entry.Category != "workflow" {
		t.Fatalf("expected workflow category, got %s", entry.Category)
	}
	if entry.ActivationMode != "direct" {
		t.Fatalf("expected direct activation mode, got %s", entry.ActivationMode)
	}

	securityEntry, ok := registry.FindByName(reg, "security-permissions")
	if !ok {
		t.Fatalf("expected security-permissions in registry")
	}
	if securityEntry.ActivationMode != "attach-only" {
		t.Fatalf("security-permissions must remain attach-only, got %s", securityEntry.ActivationMode)
	}

	releaseEntry, ok := registry.FindByName(reg, "release-store-review")
	if !ok {
		t.Fatalf("expected release-store-review in registry")
	}
	if releaseEntry.ActivationMode != "attach-only" {
		t.Fatalf("release-store-review must remain attach-only, got %s", releaseEntry.ActivationMode)
	}
}
