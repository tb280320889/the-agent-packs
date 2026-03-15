package model

type Node struct {
	ID                  string
	Level               string
	Domain              string
	Subdomain           *string
	Capability          *string
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
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Summary string   `json:"summary"`
	Score   float64  `json:"score"`
	Reason  []string `json:"reason"`
}

type RouteResult struct {
	Candidates  []RouteCandidate `json:"candidates"`
	MustInclude []string         `json:"must_include"`
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
	RecommendedValidators []string      `json:"recommended_validators"`
	RecommendedArtifacts  []string      `json:"recommended_artifacts"`
}

type ActivationResult struct {
	RequestID         string `json:"request_id"`
	Status            string `json:"status"`
	MainPack          any    `json:"main_pack"`
	Artifacts         []any  `json:"artifacts"`
	ValidationResults []any  `json:"validation_results"`
	Handoff           any    `json:"handoff"`
	Summary           string `json:"summary"`
}
