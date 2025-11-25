// Package analyzer implements the quasi-enum type safety analyzer.
package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
)

// ViolationType represents the category of violation.
type ViolationType int

const (
	// Usage violations
	VTLiteralAssignment ViolationType = iota
	VTLiteralConversion
	VTLiteralArgument
	VTLiteralCompositeField
	VTUntypedConstant
	VTVariableConversion

	// Constraint violation
	VTConstraint
)

func (vt ViolationType) String() string {
	switch vt {
	case VTLiteralAssignment:
		return "literal assignment"
	case VTLiteralConversion:
		return "literal conversion"
	case VTLiteralArgument:
		return "literal argument"
	case VTLiteralCompositeField:
		return "literal composite field"
	case VTUntypedConstant:
		return "untyped constant"
	case VTVariableConversion:
		return "variable conversion"
	case VTConstraint:
		return "constraint violation"
	default:
		return "unknown"
	}
}

// Violation represents a detected violation.
type Violation struct {
	Type          ViolationType
	Position      token.Pos
	QuasiEnumType *types.Named
	InvalidValue  ast.Expr
	Constraint    *DefinitionConstraint
	Context       ViolationContext
	SuggestedFix  string
}

// ViolationContext provides additional context for a violation.
type ViolationContext struct {
	VariableName      string
	FunctionName      string
	ParameterName     string
	Statement         string
	LineNumber        int
	ConstraintDetails string
}
