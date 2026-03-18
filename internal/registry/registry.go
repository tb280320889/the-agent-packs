package registry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

type PackageEntry struct {
	Name                   string   `json:"name"`
	Kind                   string   `json:"kind"`
	Domain                 string   `json:"domain"`
	Subdomain              string   `json:"subdomain"`
	Category               string   `json:"category"`
	CanonicalBlueprintNode string   `json:"canonical_blueprint_node"`
	VisibilityScope        string   `json:"visibility_scope"`
	ActivationMode         string   `json:"activation_mode"`
	Aliases                []string `json:"aliases"`
	Reserved               bool     `json:"reserved"`
	RecommendedValidators  []string `json:"recommended_validators"`
	RecommendedArtifacts   []string `json:"recommended_artifacts"`
	RequiredPacks          []string `json:"required_packs"`
}

type Registry struct {
	ReservedNames []string       `json:"reserved_names"`
	Packages      []PackageEntry `json:"packages"`
}

type packageManifest struct {
	Name         string           `yaml:"name"`
	Kind         string           `yaml:"kind"`
	Domain       string           `yaml:"domain"`
	Subdomain    string           `yaml:"subdomain"`
	Layer        string           `yaml:"layer"`
	Version      string           `yaml:"version"`
	Goal         string           `yaml:"goal"`
	Inputs       []string         `yaml:"inputs"`
	DependsOn    []string         `yaml:"depends_on"`
	MCP          *manifestMCP     `yaml:"mcp"`
	Validators   []string         `yaml:"validators"`
	Handoff      *manifestHandoff `yaml:"handoff"`
	Artifacts    []string         `yaml:"artifacts"`
	ExitCriteria []string         `yaml:"exit_criteria"`
}

type manifestMCP struct {
	Tools     []string `yaml:"tools"`
	Resources []string `yaml:"resources"`
	Prompts   []string `yaml:"prompts"`
}

type manifestHandoff struct {
	Incoming []string `yaml:"incoming"`
	Outgoing []string `yaml:"outgoing"`
}

var (
	defaultOnce sync.Once
	defaultReg  Registry
	defaultErr  error
)

func Default() (Registry, error) {
	defaultOnce.Do(func() {
		root, err := findProjectRoot()
		if err != nil {
			defaultErr = err
			return
		}
		defaultReg, defaultErr = Load(filepath.Join(root, "workflow-packages", "registry.json"))
	})
	return defaultReg, defaultErr
}

func Load(path string) (Registry, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Registry{}, err
	}
	var reg Registry
	if err := json.Unmarshal(raw, &reg); err != nil {
		return Registry{}, err
	}
	baseDir := filepath.Dir(path)
	if err := Validate(reg, baseDir); err != nil {
		return Registry{}, err
	}
	return reg, nil
}

func Validate(reg Registry, packagesRoot string) error {
	reserved := map[string]bool{}
	for _, name := range reg.ReservedNames {
		normalized := strings.TrimSpace(name)
		if normalized == "" {
			return errors.New("reserved name cannot be empty")
		}
		reserved[normalized] = true
	}

	allNames := map[string]bool{}
	for _, entry := range reg.Packages {
		allNames[entry.Name] = true
	}

	seenNames := map[string]bool{}
	aliases := map[string]string{}
	for _, entry := range reg.Packages {
		if err := validateEntry(entry, packagesRoot, reserved, seenNames, allNames, aliases); err != nil {
			return err
		}
		seenNames[entry.Name] = true
		for _, alias := range entry.Aliases {
			aliases[alias] = entry.Name
		}
	}
	return nil
}

func FindByName(reg Registry, name string) (PackageEntry, bool) {
	for _, entry := range reg.Packages {
		if entry.Name == name {
			return entry, true
		}
	}
	return PackageEntry{}, false
}

func FindByNode(reg Registry, nodeID string) (PackageEntry, bool) {
	for _, entry := range reg.Packages {
		if entry.CanonicalBlueprintNode == nodeID {
			return entry, true
		}
	}
	return PackageEntry{}, false
}

func validateEntry(entry PackageEntry, packagesRoot string, reserved, seenNames, allNames map[string]bool, aliases map[string]string) error {
	if strings.TrimSpace(entry.Name) == "" {
		return errors.New("package name cannot be empty")
	}
	if reserved[entry.Name] {
		return fmt.Errorf("package name %q is reserved and cannot be registered directly", entry.Name)
	}
	if seenNames[entry.Name] {
		return fmt.Errorf("duplicate package name %q", entry.Name)
	}
	if strings.TrimSpace(entry.Kind) == "" || strings.TrimSpace(entry.Category) == "" {
		return fmt.Errorf("package %q must define kind and category", entry.Name)
	}
	if strings.TrimSpace(entry.CanonicalBlueprintNode) == "" {
		return fmt.Errorf("package %q must define canonical_blueprint_node", entry.Name)
	}
	if strings.TrimSpace(entry.VisibilityScope) == "" || strings.TrimSpace(entry.ActivationMode) == "" {
		return fmt.Errorf("package %q must define visibility_scope and activation_mode", entry.Name)
	}
	if err := validateNameShape(entry); err != nil {
		return err
	}
	for _, alias := range entry.Aliases {
		if strings.TrimSpace(alias) == "" {
			return fmt.Errorf("package %q has empty alias", entry.Name)
		}
		if alias == entry.Name {
			return fmt.Errorf("package %q alias must not equal canonical name", entry.Name)
		}
		if allNames[alias] {
			return fmt.Errorf("package %q alias %q conflicts with canonical name", entry.Name, alias)
		}
		if owner, ok := aliases[alias]; ok {
			return fmt.Errorf("package %q alias %q already used by %q", entry.Name, alias, owner)
		}
	}
	manifestPath := filepath.Join(packagesRoot, entry.Name, "package.yaml")
	manifest, err := readPackageManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("package %q manifest invalid: %w", entry.Name, err)
	}
	if manifest.Name != entry.Name {
		return fmt.Errorf("package %q manifest name mismatch: %q", entry.Name, manifest.Name)
	}
	if manifest.Kind != entry.Kind {
		return fmt.Errorf("package %q manifest kind mismatch: %q", entry.Name, manifest.Kind)
	}
	if manifest.Domain != entry.Domain {
		return fmt.Errorf("package %q manifest domain mismatch: %q", entry.Name, manifest.Domain)
	}
	if manifest.Subdomain != entry.Subdomain {
		return fmt.Errorf("package %q manifest subdomain mismatch: %q", entry.Name, manifest.Subdomain)
	}
	if !sameStringSet(manifest.Validators, entry.RecommendedValidators) {
		return fmt.Errorf("package %q manifest validators mismatch: manifest=%v registry=%v", entry.Name, manifest.Validators, entry.RecommendedValidators)
	}
	if !sameStringSet(manifest.Artifacts, entry.RecommendedArtifacts) {
		return fmt.Errorf("package %q manifest artifacts mismatch: manifest=%v registry=%v", entry.Name, manifest.Artifacts, entry.RecommendedArtifacts)
	}
	if !sameStringSet(manifest.DependsOn, entry.RequiredPacks) {
		return fmt.Errorf("package %q manifest depends_on mismatch: manifest=%v registry=%v", entry.Name, manifest.DependsOn, entry.RequiredPacks)
	}
	for _, dep := range entry.RequiredPacks {
		if dep == entry.Name {
			return fmt.Errorf("package %q must not require itself", entry.Name)
		}
		if !allNames[dep] {
			return fmt.Errorf("package %q requires unknown pack %q", entry.Name, dep)
		}
	}
	return nil
}

func validateNameShape(entry PackageEntry) error {
	switch entry.Category {
	case "orchestrator":
		if entry.Domain == "" {
			return fmt.Errorf("orchestrator package %q must define domain", entry.Name)
		}
		expected := entry.Domain + "-orchestrator"
		if entry.Name != expected {
			return fmt.Errorf("orchestrator package %q must use name %q", entry.Name, expected)
		}
	case "workflow":
		if entry.Domain == "" || entry.Subdomain == "" {
			return fmt.Errorf("workflow package %q must define domain and subdomain", entry.Name)
		}
		expected := entry.Domain + "-" + entry.Subdomain
		if entry.Name != expected {
			return fmt.Errorf("workflow package %q must use name %q", entry.Name, expected)
		}
		if entry.ActivationMode == "attach-only" {
			return fmt.Errorf("workflow package %q must not be attach-only", entry.Name)
		}
	case "capability":
		if entry.Domain == "" || entry.Subdomain == "" {
			return fmt.Errorf("capability package %q must define capability line and subdomain", entry.Name)
		}
		expected := entry.Domain + "-" + entry.Subdomain
		if entry.Name != expected {
			return fmt.Errorf("capability package %q must use name %q", entry.Name, expected)
		}
		if entry.ActivationMode != "attach-only" {
			return fmt.Errorf("capability package %q must be attach-only", entry.Name)
		}
		if entry.VisibilityScope != "capability-scoped" {
			return fmt.Errorf("capability package %q must be capability-scoped", entry.Name)
		}
	default:
		return fmt.Errorf("package %q has unsupported category %q", entry.Name, entry.Category)
	}
	return nil
}

func readPackageManifest(path string) (packageManifest, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return packageManifest{}, err
	}
	var manifest packageManifest
	decoder := yaml.NewDecoder(bytes.NewReader(raw))
	decoder.KnownFields(true)
	if err := decoder.Decode(&manifest); err != nil {
		return packageManifest{}, err
	}
	if manifest.Name == "" || manifest.Kind == "" || manifest.Domain == "" || manifest.Subdomain == "" {
		return packageManifest{}, errors.New("package manifest must include name, kind, domain, subdomain")
	}
	return manifest, nil
}

func sameStringSet(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	seen := map[string]int{}
	for _, item := range left {
		seen[item]++
	}
	for _, item := range right {
		seen[item]--
	}
	for _, count := range seen {
		if count != 0 {
			return false
		}
	}
	return true
}

func findProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	current := wd
	for {
		if _, err := os.Stat(filepath.Join(current, "go.mod")); err == nil {
			if _, err := os.Stat(filepath.Join(current, "workflow-packages", "registry.json")); err == nil {
				return current, nil
			}
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return "", errors.New("project root not found")
}
