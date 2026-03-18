package compiler

import "the-agent-packs/internal/model"

type CompilePhase string

const (
	PhaseParse  CompilePhase = "parse"
	PhaseIndex  CompilePhase = "index"
	PhaseReport CompilePhase = "report"
)

type CompilerError = model.CompilerError

type CompileResult struct {
	Errors []model.CompilerError `json:"errors"`
}
