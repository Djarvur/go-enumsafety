// Package analyzer implements the quasi-enum type safety analyzer.
package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// buildQuasiEnumType constructs a QuasiEnumType from a detected type.
func buildQuasiEnumType(pass *analysis.Pass, namedType *types.Named, techniques []DetectionTechnique) *QuasiEnumType {
	// Find the type definition
	typeName := namedType.Obj()
	if typeName == nil {
		return nil
	}

	// Get underlying basic type
	basicType, ok := namedType.Underlying().(*types.Basic)
	if !ok {
		return nil
	}

	// Collect constants of this type
	var constants []EnumConstant
	var constBlock *ast.GenDecl
	var typeDecl *ast.GenDecl
	var file *ast.File

	// Find type declaration and constants
	for _, f := range pass.Files {
		// Find type declaration
		for _, decl := range f.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			// Check for type declaration
			if genDecl.Tok == token.TYPE {
				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					obj := pass.TypesInfo.Defs[typeSpec.Name]
					if obj != nil && obj.Type() == namedType {
						typeDecl = genDecl
						file = f
					}
				}
			}

			// Check for const declarations
			if genDecl.Tok == token.CONST {
				for _, spec := range genDecl.Specs {
					valueSpec, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}

					for i, name := range valueSpec.Names {
						obj := pass.TypesInfo.Defs[name]
						if obj == nil {
							continue
						}

						constObj, ok := obj.(*types.Const)
						if !ok {
							continue
						}

						// Check if this constant is of our enum type
						if constObj.Type() == namedType {
							expr := ""
							if i < len(valueSpec.Values) {
								expr = types.ExprString(valueSpec.Values[i])
							}

							constants = append(constants, EnumConstant{
								Name:          name.Name,
								Value:         constObj.Val(),
								QuasiEnumType: namedType,
								Position:      name.Pos(),
								IsIota:        expr == "iota" || expr == "",
								Expression:    expr,
								ConstBlock:    genDecl,
							})

							// Track the const block (use the first one found)
							if constBlock == nil {
								constBlock = genDecl
							}
						}
					}
				}
			}
		}
	}

	qe := &QuasiEnumType{
		Type:           namedType,
		TypeDef:        typeName,
		UnderlyingType: basicType.Kind(),
		PackagePath:    typeName.Pkg().Path(),
		Constants:      constants,
		Position:       typeName.Pos(),
		DetectedBy:     techniques,
		TypeDecl:       typeDecl,
		ConstBlock:     constBlock,
		File:           file,
	}

	// Detect helper methods (US5, US6)
	detectHelperMethods(namedType, qe)

	return qe
}

// checkVarDecl checks variable declarations for literal values and untyped constants.
func checkVarDecl(pass *analysis.Pass, registry *QuasiEnumRegistry, decl *ast.GenDecl) {
	// Only check var declarations
	if decl.Tok != token.VAR {
		return
	}

	for _, spec := range decl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		// Check each name/value pair
		for i, name := range valueSpec.Names {
			// Get the type of the variable
			obj := pass.TypesInfo.Defs[name]
			if obj == nil {
				continue
			}

			varType := obj.Type()
			if !registry.IsQuasiEnumType(varType) {
				continue
			}

			// Check if there's an initial value
			if i >= len(valueSpec.Values) {
				continue
			}

			value := valueSpec.Values[i]

			// Check for type conversion first: Status(5) or Status(x)
			if callExpr, ok := value.(*ast.CallExpr); ok {
				if isTypeConversion(pass, callExpr, varType) && len(callExpr.Args) > 0 {
					arg := callExpr.Args[0]

					// Check for literal (US1)
					if isLiteralValue(pass, arg) {
						reportUsageViolation(pass, registry, callExpr, varType, VTLiteralConversion)
						continue
					}

					// Check for variable conversion (US3)
					if ident, ok := arg.(*ast.Ident); ok {
						if isVariableConversion(pass, registry, ident, varType) {
							reportUsageViolation(pass, registry, callExpr, varType, VTVariableConversion)
							continue
						}
					}
				}
			}

			// Check for untyped constant (US2)
			if ident, ok := value.(*ast.Ident); ok {
				if isUntypedConstant(pass, registry, ident, varType) {
					reportUsageViolation(pass, registry, value, varType, VTUntypedConstant)
					continue
				}
			}

			// Then check if value is a literal
			if isLiteralValue(pass, value) {
				reportUsageViolation(pass, registry, value, varType, VTLiteralAssignment)
			}
		}
	}
}

// checkAssignment checks variable assignments for literal values, untyped constants, and variable conversions.
func checkAssignment(pass *analysis.Pass, registry *QuasiEnumRegistry, stmt *ast.AssignStmt) {
	for i, rhs := range stmt.Rhs {
		if i >= len(stmt.Lhs) {
			break
		}

		// Get the type of the LHS
		lhsType := pass.TypesInfo.TypeOf(stmt.Lhs[i])
		if lhsType == nil {
			continue
		}

		// Check if LHS is a quasi-enum type
		if !registry.IsQuasiEnumType(lhsType) {
			continue
		}

		// Check for type conversion first: Status(5) or Status(x)
		if callExpr, ok := rhs.(*ast.CallExpr); ok {
			if isTypeConversion(pass, callExpr, lhsType) && len(callExpr.Args) > 0 {
				arg := callExpr.Args[0]

				// Check for literal (US1)
				if isLiteralValue(pass, arg) {
					reportUsageViolation(pass, registry, callExpr, lhsType, VTLiteralConversion)
					continue
				}

				// Check for variable conversion (US3)
				if ident, ok := arg.(*ast.Ident); ok {
					if isVariableConversion(pass, registry, ident, lhsType) {
						reportUsageViolation(pass, registry, callExpr, lhsType, VTVariableConversion)
						continue
					}
				}
			}
		}

		// Check for untyped constant (US2)
		if ident, ok := rhs.(*ast.Ident); ok {
			if isUntypedConstant(pass, registry, ident, lhsType) {
				reportUsageViolation(pass, registry, rhs, lhsType, VTUntypedConstant)
				continue
			}
		}

		// Then check if RHS is a literal
		if isLiteralValue(pass, rhs) {
			reportUsageViolation(pass, registry, rhs, lhsType, VTLiteralAssignment)
		}
	}
}

// checkCallExpr checks function call arguments for literal values and untyped constants.
func checkCallExpr(pass *analysis.Pass, registry *QuasiEnumRegistry, call *ast.CallExpr) {
	// Get the function signature
	fnType := pass.TypesInfo.TypeOf(call.Fun)
	if fnType == nil {
		return
	}

	sig, ok := fnType.(*types.Signature)
	if !ok {
		return
	}

	// Check each argument
	params := sig.Params()
	for i, arg := range call.Args {
		if i >= params.Len() {
			break
		}

		paramType := params.At(i).Type()
		if !registry.IsQuasiEnumType(paramType) {
			continue
		}

		// Check for untyped constant (US2)
		if ident, ok := arg.(*ast.Ident); ok {
			if isUntypedConstant(pass, registry, ident, paramType) {
				reportUsageViolation(pass, registry, arg, paramType, VTUntypedConstant)
				continue
			}
		}

		// Check if argument is a literal
		if isLiteralValue(pass, arg) {
			reportUsageViolation(pass, registry, arg, paramType, VTLiteralArgument)
		}
	}
}

// checkCompositeLit checks composite literals for literal field values.
func checkCompositeLit(pass *analysis.Pass, registry *QuasiEnumRegistry, lit *ast.CompositeLit) {
	// Get the composite type
	litType := pass.TypesInfo.TypeOf(lit)
	if litType == nil {
		return
	}

	// Handle struct types
	structType, ok := litType.Underlying().(*types.Struct)
	if !ok {
		return
	}

	// Check each element
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		// Get field name
		fieldName, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}

		// Find the field in the struct
		for i := 0; i < structType.NumFields(); i++ {
			field := structType.Field(i)
			if field.Name() == fieldName.Name {
				// Check if field type is a quasi-enum
				if registry.IsQuasiEnumType(field.Type()) {
					// Check if value is a literal
					if isLiteralValue(pass, kv.Value) {
						reportUsageViolation(pass, registry, kv.Value, field.Type(), VTLiteralCompositeField)
					}
				}
				break
			}
		}
	}
}

// isUntypedConstant checks if an identifier is an untyped constant that's not a valid enum value.
// This implements US2: detecting untyped constants that aren't part of the enum definition.
func isUntypedConstant(pass *analysis.Pass, registry *QuasiEnumRegistry, ident *ast.Ident, enumType types.Type) bool {
	// Get the object this identifier refers to
	obj := pass.TypesInfo.Uses[ident]
	if obj == nil {
		return false
	}

	// Check if it's a constant
	_, ok := obj.(*types.Const)
	if !ok {
		return false
	}

	// Get the named enum type
	namedType, ok := enumType.(*types.Named)
	if !ok {
		return false
	}

	// Check if this constant is one of the defined enum constants
	qe := registry.QuasiEnums[namedType]
	if qe == nil {
		return false
	}

	// Check if the constant is in the enum's constant list
	for _, enumConst := range qe.Constants {
		if enumConst.Name == ident.Name {
			// This is a valid enum constant
			return false
		}
	}

	// This is a constant, but not one of the enum's defined constants
	// It's an untyped constant violation
	return true
}

// isVariableConversion checks if an identifier is a variable being converted to an enum type.
// This implements US3: detecting variable conversions like Color(x) where x is a variable.
func isVariableConversion(pass *analysis.Pass, registry *QuasiEnumRegistry, ident *ast.Ident, enumType types.Type) bool {
	// Get the object this identifier refers to
	obj := pass.TypesInfo.Uses[ident]
	if obj == nil {
		return false
	}

	// Check if it's a variable (not a constant)
	varObj, ok := obj.(*types.Var)
	if !ok {
		// Not a variable, could be a constant or type
		return false
	}

	// Get the variable's type
	varType := varObj.Type()

	// Get the named enum type
	namedEnumType, ok := enumType.(*types.Named)
	if !ok {
		return false
	}

	// Check if this is a valid enum constant (should not be flagged)
	qe := registry.QuasiEnums[namedEnumType]
	if qe == nil {
		return false
	}

	// Check if the variable is one of the enum constants (shouldn't happen for variables, but be safe)
	for _, enumConst := range qe.Constants {
		if enumConst.Name == ident.Name {
			return false
		}
	}

	// Check if the variable type matches the enum's underlying type
	// This catches: var x uint8 = 5; Color(x) where Color is uint8
	if varType == namedEnumType.Underlying() {
		return true
	}

	// Check if the variable is of a different enum type with the same underlying type
	// This catches: var s Status = StatusActive; Priority(s)
	if namedVarType, ok := varType.(*types.Named); ok {
		if registry.IsQuasiEnumType(namedVarType) {
			// It's a different enum type - this is a cross-enum conversion
			return true
		}
	}

	return false
}

// isLiteralValue checks if an expression is a literal value.
func isLiteralValue(pass *analysis.Pass, expr ast.Expr) bool {
	// Check for basic literals
	if _, ok := expr.(*ast.BasicLit); ok {
		return true
	}

	// Check for constant values that are not identifiers
	tv, ok := pass.TypesInfo.Types[expr]
	if !ok {
		return false
	}

	// If it's a constant value but not an identifier, it's a literal
	if tv.Value != nil {
		if _, isIdent := expr.(*ast.Ident); !isIdent {
			return true
		}
	}

	return false
}

// isTypeConversion checks if a call expression is a type conversion.
func isTypeConversion(pass *analysis.Pass, call *ast.CallExpr, targetType types.Type) bool {
	// Type conversions have exactly one argument
	if len(call.Args) != 1 {
		return false
	}

	// Check if the function is actually a type
	if ident, ok := call.Fun.(*ast.Ident); ok {
		obj := pass.TypesInfo.Uses[ident]
		if typeName, ok := obj.(*types.TypeName); ok {
			return typeName.Type() == targetType
		}
	}

	return false
}

// reportUsageViolation reports a usage violation.
func reportUsageViolation(pass *analysis.Pass, registry *QuasiEnumRegistry, node ast.Node, enumType types.Type, violationType ViolationType) {
	namedType := enumType.(*types.Named)
	qe := registry.QuasiEnums[namedType]
	if qe == nil {
		return
	}

	// Build the list of valid constants
	validConstants := make([]string, len(qe.Constants))
	for i, c := range qe.Constants {
		validConstants[i] = c.Name
	}

	msg := formatUsageViolation(violationType, namedType.Obj().Name(), validConstants)
	pass.Reportf(node.Pos(), "%s", msg)
}

// reportConstraintViolation reports a definition constraint violation.
func reportConstraintViolation(pass *analysis.Pass, qe *QuasiEnumType, violation DefinitionConstraint) {
	msg := formatConstraintViolation(qe.Type.Obj().Name(), violation)
	pass.Reportf(qe.Position, "%s", msg)
}

// formatUsageViolation formats a usage violation message.
func formatUsageViolation(vt ViolationType, typeName string, validConstants []string) string {
	switch vt {
	case VTLiteralAssignment:
		return formatMessage("literal value assigned to quasi-enum type %s; use one of: %v", typeName, validConstants)
	case VTLiteralConversion:
		return formatMessage("literal value converted to quasi-enum type %s; use one of: %v", typeName, validConstants)
	case VTLiteralArgument:
		return formatMessage("literal value passed as quasi-enum type %s; use one of: %v", typeName, validConstants)
	case VTLiteralCompositeField:
		return formatMessage("literal value in composite literal for quasi-enum type %s; use one of: %v", typeName, validConstants)
	case VTUntypedConstant:
		return formatMessage("untyped constant assigned to quasi-enum type %s; use one of: %v", typeName, validConstants)
	case VTVariableConversion:
		return formatMessage("variable converted to quasi-enum type %s; use one of: %v", typeName, validConstants)
	default:
		return formatMessage("invalid usage of quasi-enum type %s", typeName, validConstants)
	}
}

// formatConstraintViolation formats a constraint violation message.
func formatConstraintViolation(typeName string, dc DefinitionConstraint) string {
	switch dc {
	case DC001MinConstants:
		return formatMessage("quasi-enum type %s violates %s: must have at least 2 constants", typeName, dc.String())
	case DC002SameConstBlock:
		return formatMessage("quasi-enum type %s violates %s: all constants must be in the same const block", typeName, dc.String())
	case DC003SameFile:
		return formatMessage("quasi-enum type %s violates %s: type and constants must be in the same file", typeName, dc.String())
	case DC004ExclusiveConstBlock:
		return formatMessage("quasi-enum type %s violates %s: const block must contain only constants of this type", typeName, dc.String())
	case DC005Proximity:
		return formatMessage("quasi-enum type %s violates %s: type definition and const block must be adjacent", typeName, dc.String())
	default:
		return formatMessage("quasi-enum type %s violates constraint %s", typeName, dc.String())
	}
}

// formatMessage is a helper to format messages consistently.
func formatMessage(format string, args ...interface{}) string {
	// For now, just use a simple format
	// In the future, this could be enhanced with more sophisticated formatting
	if len(args) == 0 {
		return format
	}

	// Handle the case where we have a type name and valid constants
	if len(args) == 2 {
		if validConstants, ok := args[1].([]string); ok {
			return formatWithConstants(format, args[0].(string), validConstants)
		}
		// Handle two string arguments (e.g., typeName and constraint name)
		if str1, ok1 := args[0].(string); ok1 {
			if str2, ok2 := args[1].(string); ok2 {
				return formatWithTwoStrings(format, str1, str2)
			}
		}
	}

	// For constraint violations, just use the type name
	if len(args) == 1 {
		typeName := args[0].(string)
		return formatSimple(format, typeName)
	}

	return format
}

// formatWithConstants formats a message with valid constants.
func formatWithConstants(format string, typeName string, validConstants []string) string {
	return fmt.Sprintf(format, typeName, strings.Join(validConstants, ", "))
}

// formatSimple formats a simple message with just a type name.
func formatSimple(format string, typeName string) string {
	return fmt.Sprintf(format, typeName)
}

// formatWithTwoStrings formats a message with two string arguments.
func formatWithTwoStrings(format string, str1 string, str2 string) string {
	return fmt.Sprintf(format, str1, str2)
}
