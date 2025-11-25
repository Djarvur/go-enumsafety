// Package analyzer implements the quasi-enum type safety analyzer.
package analyzer

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
)

// DetectionTechnique represents the technique used to identify a quasi-enum type.
type DetectionTechnique int

const (
	DT001ConstantsBased DetectionTechnique = iota
	DT002NameSuffix
	DT003InlineComment
	DT004PrecedingComment
	DT005NamedComment
)

func (dt DetectionTechnique) String() string {
	switch dt {
	case DT001ConstantsBased:
		return "DT-001 (constants-based)"
	case DT002NameSuffix:
		return "DT-002 (name suffix)"
	case DT003InlineComment:
		return "DT-003 (inline comment)"
	case DT004PrecedingComment:
		return "DT-004 (preceding comment)"
	case DT005NamedComment:
		return "DT-005 (named comment)"
	default:
		return "unknown"
	}
}

// DefinitionConstraint represents a constraint that quasi-enum definitions must satisfy.
type DefinitionConstraint int

const (
	DC001MinConstants DefinitionConstraint = iota
	DC002SameConstBlock
	DC003SameFile
	DC004ExclusiveConstBlock
	DC005Proximity
)

func (dc DefinitionConstraint) String() string {
	switch dc {
	case DC001MinConstants:
		return "DC-001 (minimum 2 constants)"
	case DC002SameConstBlock:
		return "DC-002 (same const block)"
	case DC003SameFile:
		return "DC-003 (same file)"
	case DC004ExclusiveConstBlock:
		return "DC-004 (exclusive const block)"
	case DC005Proximity:
		return "DC-005 (proximity)"
	default:
		return "unknown"
	}
}

// QuasiEnumType represents a detected quasi-enum type with its constants and metadata.
type QuasiEnumType struct {
	Type           *types.Named    // The named type
	TypeDef        *types.TypeName // Type definition object
	UnderlyingType types.BasicKind
	PackagePath    string
	Constants      []EnumConstant
	Position       token.Pos
	DetectedBy     []DetectionTechnique
	TypeDecl       *ast.GenDecl // Type declaration node (for constraint validation)
	ConstBlock     *ast.GenDecl // Const block node (for constraint validation)
	File           *ast.File    // File containing the type (for constraint validation)

	// Helper method tracking (US5, US6)
	HasStringMethod        bool // Whether type has String() string method
	HasUnmarshalTextMethod bool // Whether type has UnmarshalText([]byte) error method

	// FR-047: Suggest adding enum comment when detected only by constants-based
	SuggestEnumComment bool
}

// EnumConstant represents a valid constant value for a quasi-enum type.
type EnumConstant struct {
	Name          string
	Value         constant.Value
	QuasiEnumType *types.Named
	Position      token.Pos
	IsIota        bool
	Expression    string
	ConstBlock    *ast.GenDecl
}

// DetectionConfig holds configuration for detection techniques.
type DetectionConfig struct {
	ConstantsDetectionEnabled        bool
	SuffixDetectionEnabled           bool
	InlineCommentDetectionEnabled    bool
	PrecedingCommentDetectionEnabled bool
	NamedCommentDetectionEnabled     bool
}

// NewDetectionConfig creates a new DetectionConfig with defaults.
func NewDetectionConfig() *DetectionConfig {
	return &DetectionConfig{
		ConstantsDetectionEnabled:        !disableConstantsDetection,
		SuffixDetectionEnabled:           !disableSuffixDetection,
		InlineCommentDetectionEnabled:    !disableInlineCommentDetection,
		PrecedingCommentDetectionEnabled: !disablePrecedingCommentDetection,
		NamedCommentDetectionEnabled:     !disableNamedCommentDetection,
	}
}

// ConstraintConfig holds configuration for definition constraints.
type ConstraintConfig struct {
	MinConstantsEnabled   bool
	SameConstBlockEnabled bool
	SameFileEnabled       bool
	ExclusiveBlockEnabled bool
	ProximityEnabled      bool
}

// NewConstraintConfig creates a new ConstraintConfig with defaults.
func NewConstraintConfig() *ConstraintConfig {
	return &ConstraintConfig{
		MinConstantsEnabled:   !disableMinConstantsCheck,
		SameConstBlockEnabled: !disableSameBlockCheck,
		SameFileEnabled:       !disableSameFileCheck,
		ExclusiveBlockEnabled: !disableExclusiveBlockCheck,
		ProximityEnabled:      !disableProximityCheck,
	}
}

// QuasiEnumRegistry is the central registry of all quasi-enum types.
type QuasiEnumRegistry struct {
	QuasiEnums       map[*types.Named]*QuasiEnumType
	ConstantLookup   map[*types.Named]map[string]*EnumConstant
	Packages         map[string][]*QuasiEnumType
	DetectionConfig  *DetectionConfig
	ConstraintConfig *ConstraintConfig
}

// NewQuasiEnumRegistry creates a new registry.
func NewQuasiEnumRegistry(detectionConfig *DetectionConfig, constraintConfig *ConstraintConfig) *QuasiEnumRegistry {
	return &QuasiEnumRegistry{
		QuasiEnums:       make(map[*types.Named]*QuasiEnumType),
		ConstantLookup:   make(map[*types.Named]map[string]*EnumConstant),
		Packages:         make(map[string][]*QuasiEnumType),
		DetectionConfig:  detectionConfig,
		ConstraintConfig: constraintConfig,
	}
}

// RegisterQuasiEnum adds a quasi-enum to the registry.
func (r *QuasiEnumRegistry) RegisterQuasiEnum(qe *QuasiEnumType) {
	r.QuasiEnums[qe.Type] = qe

	// Build constant lookup
	if _, exists := r.ConstantLookup[qe.Type]; !exists {
		r.ConstantLookup[qe.Type] = make(map[string]*EnumConstant)
	}
	for i := range qe.Constants {
		r.ConstantLookup[qe.Type][qe.Constants[i].Name] = &qe.Constants[i]
	}

	// Add to package list
	r.Packages[qe.PackagePath] = append(r.Packages[qe.PackagePath], qe)
}

// IsQuasiEnumType checks if a type is a quasi-enum.
func (r *QuasiEnumRegistry) IsQuasiEnumType(t types.Type) bool {
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	_, exists := r.QuasiEnums[named]
	return exists
}

// GetEnumConstants returns the valid constants for a quasi-enum type.
func (r *QuasiEnumRegistry) GetEnumConstants(t *types.Named) []EnumConstant {
	if qe, exists := r.QuasiEnums[t]; exists {
		return qe.Constants
	}
	return nil
}
