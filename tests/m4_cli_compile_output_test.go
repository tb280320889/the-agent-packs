package tests

import (
	"encoding/json"
	"os/exec"
	"strings"
	"testing"
)

func TestM4CLICompileOutput(t *testing.T) {
	cmd := exec.Command("go", "run", "../cmd/agent-pack-mcp", "compile", "--root", "../blueprint", "--db", "../blueprint/index/blueprint.db", "--report-dir", "../blueprint/index")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("compile command failed: %v, output: %s", err, string(output))
	}
	text := strings.TrimSpace(string(output))
	if text == "" {
		t.Fatalf("expected compile output")
	}
	jsonStart := strings.Index(text, "{")
	if jsonStart == -1 {
		t.Fatalf("expected JSON output, got %s", text)
	}
	jsonPayload := text[jsonStart:]
	var decoded struct {
		Errors []map[string]any `json:"errors"`
	}
	if err := json.Unmarshal([]byte(jsonPayload), &decoded); err != nil {
		t.Fatalf("expected JSON compile result, got %s", text)
	}
}
