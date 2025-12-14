package analyzer

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// checkStringMethod warns if a quasi-enum type lacks a String() method (US5).
func checkStringMethod(pass *analysis.Pass, registry *QuasiEnumRegistry) {
	if disableStringMethodCheck {
		return
	}

	for _, qe := range registry.QuasiEnums {
		if !qe.HasStringMethod {
			warnMissingStringMethod(pass, qe)
		}
	}
}

// checkUnmarshalTextMethod warns if a quasi-enum type lacks an UnmarshalText() method (US6).
func checkUnmarshalTextMethod(pass *analysis.Pass, registry *QuasiEnumRegistry) {
	if disableUnmarshalMethodCheck {
		return
	}

	for _, qe := range registry.QuasiEnums {
		if !qe.HasUnmarshalTextMethod {
			warnMissingUnmarshalTextMethod(pass, qe)
		}
	}
}

// warnMissingStringMethod reports a warning for missing String() method.
func warnMissingStringMethod(pass *analysis.Pass, qe *QuasiEnumType) {
	typeName := qe.Type.Obj().Name()

	msg := fmt.Sprintf(
		"quasi-enum type %s lacks a String() method; consider using golang.org/x/tools/cmd/stringer or github.com/Djarvur/go-silly-enum to generate it",
		typeName,
	)

	pass.Reportf(qe.Position, "%s", msg)
}

// warnMissingUnmarshalTextMethod reports a warning for missing UnmarshalText() method.
func warnMissingUnmarshalTextMethod(pass *analysis.Pass, qe *QuasiEnumType) {
	typeName := qe.Type.Obj().Name()

	msg := fmt.Sprintf(
		"quasi-enum type %s lacks an UnmarshalText([]byte) error method; consider using github.com/Djarvur/go-silly-enum to generate it",
		typeName,
	)

	pass.Reportf(qe.Position, "%s", msg)
}

// hasMethod checks if a named type has a method with the given name and signature.
func hasMethod(namedType *types.Named, methodName string, signature string) bool {
	// Get the method set for the type
	methodSet := types.NewMethodSet(namedType)

	// Look for the method
	for i := 0; i < methodSet.Len(); i++ {
		method := methodSet.At(i)
		if method.Obj().Name() == methodName {
			// For String(), check signature: func() string
			if methodName == "String" {
				sig, ok := method.Type().(*types.Signature)
				if !ok {
					continue
				}
				// No parameters, one string result
				if sig.Params().Len() == 0 && sig.Results().Len() == 1 {
					if basic, ok := sig.Results().At(0).Type().(*types.Basic); ok {
						if basic.Kind() == types.String {
							return true
						}
					}
				}
			}

			// For UnmarshalText, check signature: func([]byte) error
			if methodName == "UnmarshalText" {
				sig, ok := method.Type().(*types.Signature)
				if !ok {
					continue
				}
				// One []byte parameter, one error result
				if sig.Params().Len() == 1 && sig.Results().Len() == 1 {
					// Check parameter is []byte
					if slice, ok := sig.Params().At(0).Type().(*types.Slice); ok {
						if basic, ok := slice.Elem().(*types.Basic); ok {
							if basic.Kind() == types.Byte {
								// Check result is error
								if named, ok := sig.Results().At(0).Type().(*types.Named); ok {
									if named.Obj().Name() == "error" {
										return true
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return false
}

// detectHelperMethods updates the QuasiEnumType with helper method presence.
func detectHelperMethods(namedType *types.Named, qe *QuasiEnumType) {
	qe.HasStringMethod = hasMethod(namedType, "String", "func() string")

	// UnmarshalText is typically defined on pointer receiver, so check both
	qe.HasUnmarshalTextMethod = hasMethod(namedType, "UnmarshalText", "func([]byte) error")
	if !qe.HasUnmarshalTextMethod {
		// Check pointer receiver
		ptrType := types.NewPointer(namedType)
		methodSet := types.NewMethodSet(ptrType)
		for i := 0; i < methodSet.Len(); i++ {
			method := methodSet.At(i)
			if method.Obj().Name() == "UnmarshalText" {
				sig, ok := method.Type().(*types.Signature)
				if !ok {
					continue
				}
				// One []byte parameter, one error result
				if sig.Params().Len() == 1 && sig.Results().Len() == 1 {
					// Check parameter is []byte
					if slice, ok := sig.Params().At(0).Type().(*types.Slice); ok {
						if basic, ok := slice.Elem().(*types.Basic); ok {
							if basic.Kind() == types.Byte {
								// Check result is error (interface type)
								resultType := sig.Results().At(0).Type()
								if resultType.String() == "error" {
									qe.HasUnmarshalTextMethod = true
									break
								}
							}
						}
					}
				}
			}
		}
	}
}
