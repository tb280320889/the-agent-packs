package model

type Node struct {
	ID                  string
	Level               string
	Domain              string
	Subdomain           *string
	Capability          *string
	NodeKind            string
	VisibilityScope     string
	ActivationMode      string
	Title               string
	Summary             string
	Path                string
	ParentID            *string
	BodyMD              string
	EntryConditionsJSON string
	StopConditionsJSON  string
	Checksum            string
	UpdatedAt           string
}

type NodeMeta struct {
	NodeID       string
	Aliases      string
	Triggers     string
	AntiTriggers string
	Tags         string
}

type Edge struct {
	SourceID string `json:"source_id"`
	TargetID string `json:"target_id"`
	EdgeType string `json:"edge_type"`
}

type RouteCandidate struct {
	ID              string   `json:"id"`
	Pack            string   `json:"pack,omitempty"`
	Title           string   `json:"title"`
	Summary         string   `json:"summary"`
	Score           float64  `json:"score"`
	Reason          []string `json:"reason"`
	ReasonCode      string   `json:"reason_code,omitempty"`
	RuleRef         string   `json:"rule_ref,omitempty"`
	NextAction      string   `json:"next_action,omitempty"`
	DocsRef         string   `json:"docs_ref"`
	NodeKind        string   `json:"node_kind,omitempty"`
	VisibilityScope string   `json:"visibility_scope,omitempty"`
	ActivationMode  string   `json:"activation_mode,omitempty"`
}

type RouteCapabilityDecision struct {
	Pack       string `json:"pack"`
	NodeID     string `json:"node_id"`
	Attached   bool   `json:"attached"`
	ReasonCode string `json:"reason_code"`
	RuleRef    string `json:"rule_ref"`
	Message    string `json:"message"`
	NextAction string `json:"next_action,omitempty"`
	DocsRef    string `json:"docs_ref"`
}

type RouteResult struct {
	Status              string                    `json:"status"`
	ErrorCode           string                    `json:"error_code,omitempty"`
	Message             string                    `json:"message,omitempty"`
	NextAction          string                    `json:"next_action,omitempty"`
	DecisionBasis       string                    `json:"decision_basis,omitempty"`
	DecisionTraceID     string                    `json:"decision_trace_id,omitempty"`
	DocsRef             string                    `json:"docs_ref"`
	Candidates          []RouteCandidate          `json:"candidates"`
	MustInclude         []string                  `json:"must_include"`
	CapabilityDecisions []RouteCapabilityDecision `json:"capability_decisions,omitempty"`
	Details             map[string]any            `json:"details,omitempty"`
}

type NodeSummary struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Summary string `json:"summary,omitempty"`
	Body    string `json:"body,omitempty"`
	Level   string `json:"level"`
}

type ContextBundle struct {
	Main                  *NodeSummary  `json:"main"`
	Required              []NodeSummary `json:"required"`
	ExecutionChildren     []NodeSummary `json:"execution_children"`
	Deferred              []NodeSummary `json:"deferred"`
	RequiredPacks         []string      `json:"required_packs"`
	RecommendedValidators []string      `json:"recommended_validators"`
	RecommendedArtifacts  []string      `json:"recommended_artifacts"`
}

type ActivationResult struct {
	RequestID         string               `json:"request_id"`
	Status            string               `json:"status"`
	MainPack          *string              `json:"main_pack"`
	RouteStatus       string               `json:"route_status,omitempty"`
	RouteErrorCode    string               `json:"route_error_code,omitempty"`
	RouteNextAction   string               `json:"route_next_action,omitempty"`
	RouteDecision     string               `json:"route_decision_basis,omitempty"`
	RouteTraceID      string               `json:"route_decision_trace_id,omitempty"`
	RouteDocsRef      string               `json:"route_docs_ref"`
	Artifacts         []Artifact           `json:"artifacts"`
	ValidationResults []ValidationEnvelope `json:"validation_results"`
	Handoff           map[string]any       `json:"handoff"`
	Summary           string               `json:"summary"`
}

type Artifact struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type ValidationPlan struct {
	PlanID                   string            `json:"plan_id"`
	RequestID                string            `json:"request_id"`
	MainPack                 string            `json:"main_pack"`
	Validators               []ValidatorPlan   `json:"validators"`
	ArtifactsUnderValidation []string          `json:"artifacts_under_validation"`
	SeverityPolicy           map[string]string `json:"severity_policy"`
	PlanReason               string            `json:"plan_reason"`
}

type ValidatorPlan struct {
	Name   string `json:"name"`
	Scope  string `json:"scope"`
	Reason string `json:"reason"`
}

type Finding struct {
	Severity    string `json:"severity"`
	Code        string `json:"code"`
	Message     string `json:"message"`
	ArtifactRef string `json:"artifact_ref"`
}

type ValidatorResult struct {
	ValidatorName      string    `json:"validator_name"`
	Status             string    `json:"status"`
	Findings           []Finding `json:"findings"`
	RepairSuggestions  []string  `json:"repair_suggestions"`
	ValidatedArtifacts []string  `json:"validated_artifacts"`
}

type ValidationEnvelope struct {
	ValidationPlan   ValidationPlan    `json:"validation_plan"`
	ValidatorResults []ValidatorResult `json:"validator_results"`
}

type CompilerError struct {
	Phase   string `json:"phase"`
	Path    string `json:"path"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Line    int    `json:"line,omitempty"`
	Column  int    `json:"column,omitempty"`
}

type CompileResult struct {
	Errors []CompilerError `json:"errors"`
}
