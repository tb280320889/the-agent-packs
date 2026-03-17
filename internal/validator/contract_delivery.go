package validator

import (
	"fmt"
	"strings"

	"the-agent-packs/internal/model"
)

func validateContractDelivery(plan model.ValidationPlan, input ExecutionInput) model.ValidatorResult {
	validated := plan.ArtifactsUnderValidation
	if len(validated) == 0 {
		validated = []string{}
	}

	bundle := input.ContractBundle
	if bundle == nil {
		return model.ValidatorResult{
			ValidatorName:      "validator-contract-delivery",
			Status:             "failed",
			Findings:           []model.Finding{{Severity: "error", Code: "CONTRACT_BUNDLE_MISSING", Message: "Contract bundle is missing from validator input.", RuleRef: "CONT-03", SourceRule: "CONT-03", ArtifactRef: "context-bundle"}},
			RepairSuggestions:  []string{"在验证输入中传递 context bundle contract 数据后重试。"},
			ValidatedArtifacts: validated,
		}
	}

	findings := []model.Finding{}
	repair := []string{}

	mainDomain := domainFromNodeID("")
	if bundle.Main != nil {
		mainDomain = domainFromNodeID(bundle.Main.ID)
	}

	includedByNode := map[string]model.ContractDecision{}
	for _, d := range bundle.IncludedDecisions {
		includedByNode[d.NodeID] = d
		if !isTraceableRule(d.SourceRule) {
			findings = append(findings, model.Finding{
				Severity:    "error",
				Code:        "CONTRACT_SOURCE_RULE_UNTRACEABLE",
				Message:     fmt.Sprintf("Included decision source_rule is not traceable: %q", d.SourceRule),
				RuleRef:     "CONT-03",
				SourceRule:  d.SourceRule,
				ArtifactRef: d.NodeID,
			})
		}
		if strings.TrimSpace(d.HumanNote) == "" || len([]rune(strings.TrimSpace(d.HumanNote))) < 8 {
			findings = append(findings, model.Finding{
				Severity:    "warn",
				Code:        "CONTRACT_HUMAN_NOTE_WEAK",
				Message:     "Human note is too short for manual review context.",
				RuleRef:     "CONT-02",
				SourceRule:  d.SourceRule,
				ArtifactRef: d.NodeID,
			})
		}
		nodeDomain := domainFromNodeID(d.NodeID)
		if mainDomain != "" && nodeDomain != "" && nodeDomain != mainDomain && d.Scope != "attach_only_capability" {
			findings = append(findings, model.Finding{
				Severity:    "error",
				Code:        "CONTRACT_CROSS_DOMAIN_INCLUDED",
				Message:     fmt.Sprintf("Included node %q is outside target domain %q without legal attach-only scope.", d.NodeID, mainDomain),
				RuleRef:     "CONT-01",
				SourceRule:  d.SourceRule,
				ArtifactRef: d.NodeID,
			})
		}
	}

	for _, d := range bundle.ExcludedDecisions {
		if !isTraceableRule(d.SourceRule) {
			findings = append(findings, model.Finding{
				Severity:    "error",
				Code:        "CONTRACT_SOURCE_RULE_UNTRACEABLE",
				Message:     fmt.Sprintf("Excluded decision source_rule is not traceable: %q", d.SourceRule),
				RuleRef:     "CONT-03",
				SourceRule:  d.SourceRule,
				ArtifactRef: d.NodeID,
			})
		}
		if strings.TrimSpace(d.HumanNote) == "" || len([]rune(strings.TrimSpace(d.HumanNote))) < 8 {
			findings = append(findings, model.Finding{
				Severity:    "warn",
				Code:        "CONTRACT_HUMAN_NOTE_WEAK",
				Message:     "Human note is too short for manual review context.",
				RuleRef:     "CONT-02",
				SourceRule:  d.SourceRule,
				ArtifactRef: d.NodeID,
			})
		}
	}

	if bundle.Main != nil {
		mainDecision, ok := includedByNode[bundle.Main.ID]
		if !ok || mainDecision.Action != "include" {
			findings = append(findings, model.Finding{
				Severity:    "error",
				Code:        "CONTRACT_REQUIRED_MISSING",
				Message:     "Main node is missing include decision.",
				RuleRef:     "CONT-03",
				SourceRule:  "CONT-03",
				ArtifactRef: bundle.Main.ID,
			})
		}
	}

	for _, required := range bundle.Required {
		d, ok := includedByNode[required.ID]
		if !ok || d.Action != "include" {
			findings = append(findings, model.Finding{
				Severity:    "error",
				Code:        "CONTRACT_REQUIRED_MISSING",
				Message:     fmt.Sprintf("Required node %q is missing include decision.", required.ID),
				RuleRef:     "CONT-03",
				SourceRule:  "CONT-03",
				ArtifactRef: required.ID,
			})
		}
	}

	for _, f := range findings {
		switch f.Code {
		case "CONTRACT_CROSS_DOMAIN_INCLUDED":
			repair = appendUnique(repair, "移除跨域 include 或将其改为合法 attach_only_capability，并补充依据。")
		case "CONTRACT_REQUIRED_MISSING":
			repair = appendUnique(repair, "为 main/required 节点补齐 include decision，确保最小且完整可复验。")
		case "CONTRACT_SOURCE_RULE_UNTRACEABLE":
			repair = appendUnique(repair, "将 source_rule 修正为可追溯规则 ID（如 BR-xx/CONT-xx）。")
		case "CONTRACT_HUMAN_NOTE_WEAK":
			repair = appendUnique(repair, "补充更清晰的 human_note，说明包含或排除的最小修复理由。")
		}
	}

	status := "passed"
	hasWarn := false
	for _, f := range findings {
		if f.Severity == "error" {
			status = "failed"
			break
		}
		if f.Severity == "warn" {
			hasWarn = true
		}
	}
	if status != "failed" && hasWarn {
		status = "warned"
	}

	return model.ValidatorResult{
		ValidatorName:      "validator-contract-delivery",
		Status:             status,
		Findings:           findings,
		RepairSuggestions:  repair,
		ValidatedArtifacts: validated,
	}
}

func domainFromNodeID(nodeID string) string {
	parts := strings.Split(nodeID, ".")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

func isTraceableRule(sourceRule string) bool {
	s := strings.TrimSpace(strings.ToUpper(sourceRule))
	if s == "" {
		return false
	}
	return strings.HasPrefix(s, "BR-") ||
		strings.HasPrefix(s, "CONT-") ||
		strings.HasPrefix(s, "ROUT-") ||
		strings.HasPrefix(s, "VALD-") ||
		strings.HasPrefix(s, "GOVR-") ||
		strings.HasPrefix(s, "PARS-") ||
		strings.HasPrefix(s, "INDX-") ||
		strings.HasPrefix(s, "DOMN-") ||
		strings.HasPrefix(s, "REQ-")
}

func appendUnique(list []string, item string) []string {
	for _, existing := range list {
		if existing == item {
			return list
		}
	}
	return append(list, item)
}
