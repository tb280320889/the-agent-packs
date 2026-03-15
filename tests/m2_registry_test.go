package tests

import (
	"path/filepath"
	"testing"

	"the-agent-packs/internal/query"
	"the-agent-packs/internal/registry"
)

func TestM2RegistryLoadsAndValidates(t *testing.T) {
	root := projectRoot(t)
	reg, err := registry.Load(filepath.Join(root, "workflow-packages", "registry.json"))
	if err != nil {
		t.Fatalf("load registry failed: %v", err)
	}
	if len(reg.Packages) != 3 {
		t.Fatalf("expected 3 registered packages, got %d", len(reg.Packages))
	}
	entry, ok := registry.FindByName(reg, "wxt-manifest")
	if !ok {
		t.Fatalf("expected wxt-manifest in registry")
	}
	if entry.Category != "workflow" || entry.CanonicalBlueprintNode != "L1.wxt.manifest" {
		t.Fatalf("unexpected entry: %+v", entry)
	}
	if len(entry.RequiredPacks) != 2 || entry.RequiredPacks[0] != "security-permissions" || entry.RequiredPacks[1] != "release-store-review" {
		t.Fatalf("expected registry required packs from manifest, got %+v", entry.RequiredPacks)
	}
}

func TestM2RegistryReservedBareNameRejected(t *testing.T) {
	root := projectRoot(t)
	reg, err := registry.Load(filepath.Join(root, "workflow-packages", "registry.json"))
	if err != nil {
		t.Fatalf("load registry failed: %v", err)
	}
	bad := registry.PackageEntry{
		Name:                   "security",
		Kind:                   "workflow-package",
		Domain:                 "security",
		Subdomain:              "permissions",
		Category:               "capability",
		CanonicalBlueprintNode: "L1.security.permissions",
		VisibilityScope:        "capability-scoped",
		ActivationMode:         "attach-only",
		Aliases:                []string{"security review"},
	}
	reg.Packages = append(reg.Packages, bad)
	err = registry.Validate(reg, filepath.Join(root, "workflow-packages"))
	if err == nil {
		t.Fatalf("expected reserved bare name validation failure")
	}
}

func TestM2RegistryCapabilityAttachOnlyAndContextBundleUsesRegistry(t *testing.T) {
	reg, err := registry.Default()
	if err != nil {
		t.Fatalf("default registry failed: %v", err)
	}
	entry, ok := registry.FindByNode(reg, "L1.security.permissions")
	if !ok {
		t.Fatalf("expected registry entry for L1.security.permissions")
	}
	if entry.ActivationMode != "attach-only" || entry.VisibilityScope != "capability-scoped" {
		t.Fatalf("unexpected capability entry: %+v", entry)
	}

	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()
	bundle, err := query.BuildContextBundle(db, "L1.wxt.manifest", true, false, false)
	if err != nil {
		t.Fatalf("build bundle failed: %v", err)
	}
	if len(bundle.RecommendedValidators) != 2 {
		t.Fatalf("expected registry-backed validators, got %+v", bundle.RecommendedValidators)
	}
	if len(bundle.RecommendedArtifacts) != 1 || bundle.RecommendedArtifacts[0] != "manifest-review.md" {
		t.Fatalf("expected registry-backed artifacts, got %+v", bundle.RecommendedArtifacts)
	}
}

func TestM2RouteTargetPackUsesRegistryCanonicalNode(t *testing.T) {
	dbPath := compileMainIndex(t)
	db := openDB(t, dbPath)
	defer db.Close()

	tp := "release-store-review"
	result, err := query.RouteQuery(db, "L1", "submission review", &tp, nil, nil, nil, nil, 3)
	if err != nil {
		t.Fatalf("route query failed: %v", err)
	}
	if len(result.Candidates) == 0 || result.Candidates[0].ID != "L1.release.store-review" {
		t.Fatalf("expected registry-mapped target pack route, got %+v", result.Candidates)
	}
}
