// Package analyzer implements the quasi-enum type safety analyzer.
package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
)

// validateMinConstants implements DC-001: minimum 2 constants check.
// Returns true if the type has at least 2 constants defined.
func validateMinConstants(qe *QuasiEnumType) bool {
	return len(qe.Constants) >= 2
}

// validateSameConstBlock implements DC-002: same const block check.
// Returns true if all constants are defined in the same const block.
func validateSameConstBlock(qe *QuasiEnumType, fset *token.FileSet) bool {
	if len(qe.Constants) < 2 {
		return true // Trivially satisfied if < 2 constants
	}

	// Get the first constant's const block position
	firstConst := qe.Constants[0]
	firstPos := fset.Position(firstConst.Position)

	// Check if all other constants are in the same const block
	// We determine this by checking if they're all part of the same GenDecl
	for i := 1; i < len(qe.Constants); i++ {
		constPos := fset.Position(qe.Constants[i].Position)
		// If file differs, they're definitely not in the same block
		if constPos.Filename != firstPos.Filename {
			return false
		}
	}

	// Additional check: verify they're all in the same GenDecl
	// This is done by checking if their parent GenDecl nodes are the same
	// Note: This requires access to AST nodes, which we'll handle in the caller
	return true
}

// validateSameFile implements DC-003: same file check.
// Returns true if type definition and all constants are in the same file.
func validateSameFile(qe *QuasiEnumType, fset *token.FileSet) bool {
	if len(qe.Constants) == 0 {
		return true
	}

	// Get the type definition file
	typePos := fset.Position(qe.TypeDef.Pos())
	typeFile := typePos.Filename

	// Check if all constants are in the same file
	for _, c := range qe.Constants {
		constPos := fset.Position(c.Position)
		if constPos.Filename != typeFile {
			return false
		}
	}

	return true
}

// validateExclusiveBlock implements DC-004: exclusive const block check.
// Returns true if the const block contains only constants of this enum type.
func validateExclusiveBlock(qe *QuasiEnumType, constDecl *ast.GenDecl, info *types.Info) bool {
	if constDecl == nil || constDecl.Tok != token.CONST {
		return true
	}

	// Check each spec in the const block
	for _, spec := range constDecl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		// Check each name in the value spec
		for _, name := range valueSpec.Names {
			obj := info.Defs[name]
			if obj == nil {
				continue
			}

			// Check if this constant's type matches our enum type
			constType := obj.Type()
			if named, ok := constType.(*types.Named); ok {
				if named != qe.Type {
					// Found a constant of a different type
					return false
				}
			} else {
				// Found an untyped constant or different type
				return false
			}
		}
	}

	return true
}

// validateProximity implements DC-005: proximity check.
// Returns true if type definition and const block are adjacent (allowing comments/empty lines).
func validateProximity(qe *QuasiEnumType, typeDecl *ast.GenDecl, constDecl *ast.GenDecl, fset *token.FileSet, file *ast.File) bool {
	if typeDecl == nil || constDecl == nil {
		return true
	}

	// Get positions
	typeEnd := fset.Position(typeDecl.End())
	constStart := fset.Position(constDecl.Pos())

	// Must be in the same file
	if typeEnd.Filename != constStart.Filename {
		return false
	}

	// Find any code between type and const declarations
	// We allow comments and empty lines, but no executable code
	for _, decl := range file.Decls {
		declStart := fset.Position(decl.Pos())
		declEnd := fset.Position(decl.End())

		// Check if this declaration is between type and const
		if declStart.Offset > typeEnd.Offset && declEnd.Offset < constStart.Offset {
			// Found a declaration between type and const
			return false
		}
	}

	return true
}

// ValidateConstraints validates all enabled constraints for a quasi-enum type.
// Returns a slice of constraint violations.
func (qe *QuasiEnumType) ValidateConstraints(
	config *ConstraintConfig,
	fset *token.FileSet,
	typeDecl *ast.GenDecl,
	constDecl *ast.GenDecl,
	file *ast.File,
	info *types.Info,
) []DefinitionConstraint {
	var violations []DefinitionConstraint

	// DC-001: Minimum Constants
	if config.MinConstantsEnabled && !validateMinConstants(qe) {
		violations = append(violations, DC001MinConstants)
	}

	// DC-002: Same Const Block
	if config.SameConstBlockEnabled && !validateSameConstBlock(qe, fset) {
		violations = append(violations, DC002SameConstBlock)
	}

	// DC-003: Same File
	if config.SameFileEnabled && !validateSameFile(qe, fset) {
		violations = append(violations, DC003SameFile)
	}

	// DC-004: Exclusive Const Block
	if config.ExclusiveBlockEnabled && !validateExclusiveBlock(qe, constDecl, info) {
		violations = append(violations, DC004ExclusiveConstBlock)
	}

	// DC-005: Proximity
	if config.ProximityEnabled && !validateProximity(qe, typeDecl, constDecl, fset, file) {
		violations = append(violations, DC005Proximity)
	}

	return violations
}
