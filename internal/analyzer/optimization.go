package analyzer

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// checkUint8Optimization suggests using uint8 for enums with <256 constants using larger types (US4).
func checkUint8Optimization(pass *analysis.Pass, registry *QuasiEnumRegistry) {
	if disableUint8Suggestion {
		return
	}

	for _, qe := range registry.QuasiEnums {
		// Check if already using uint8
		if qe.UnderlyingType == types.Uint8 {
			continue
		}

		// Count constants
		constCount := len(qe.Constants)
		if constCount >= 256 {
			continue
		}

		// Check if using a larger integer type
		if !isLargerIntegerType(qe.UnderlyingType) {
			continue
		}

		// Suggest uint8 with autofix
		suggestUint8(pass, qe)
	}
}

// isLargerIntegerType checks if the type is larger than uint8.
func isLargerIntegerType(kind types.BasicKind) bool {
	switch kind {
	case types.Int, types.Int16, types.Int32, types.Int64,
		types.Uint, types.Uint16, types.Uint32, types.Uint64:
		return true
	default:
		return false
	}
}

// suggestUint8 creates a suggestion to use uint8 with autofix capability.
func suggestUint8(pass *analysis.Pass, qe *QuasiEnumType) {
	typeName := qe.Type.Obj().Name()

	// Get the string representation of the underlying type
	var currentType string
	switch qe.UnderlyingType {
	case types.Int:
		currentType = "int"
	case types.Int8:
		currentType = "int8"
	case types.Int16:
		currentType = "int16"
	case types.Int32:
		currentType = "int32"
	case types.Int64:
		currentType = "int64"
	case types.Uint:
		currentType = "uint"
	case types.Uint16:
		currentType = "uint16"
	case types.Uint32:
		currentType = "uint32"
	case types.Uint64:
		currentType = "uint64"
	default:
		currentType = "unknown"
	}

	msg := fmt.Sprintf(
		"quasi-enum type %s uses %s but has only %d constants; consider using uint8 for memory optimization",
		typeName,
		currentType,
		len(qe.Constants),
	)

	// Create diagnostic with suggested fix
	diagnostic := analysis.Diagnostic{
		Pos:      qe.Position,
		Message:  msg,
		Category: "optimization",
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: fmt.Sprintf("Change %s base type to uint8", typeName),
				TextEdits: []analysis.TextEdit{
					{
						Pos:     qe.TypeDecl.Pos(),
						End:     qe.TypeDecl.End(),
						NewText: []byte(generateUint8TypeDecl(qe)),
					},
				},
			},
		},
	}

	pass.Report(diagnostic)
}

// generateUint8TypeDecl generates the type declaration with uint8 base type.
func generateUint8TypeDecl(qe *QuasiEnumType) string {
	typeName := qe.Type.Obj().Name()

	// Find the type spec within the GenDecl
	for _, spec := range qe.TypeDecl.Specs {
		if typeSpec, ok := spec.(*ast.TypeSpec); ok {
			if typeSpec.Name.Name == typeName {
				// Generate: type TypeName uint8
				return fmt.Sprintf("type %s uint8", typeName)
			}
		}
	}

	// Fallback
	return fmt.Sprintf("type %s uint8", typeName)
}
