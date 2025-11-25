// Package analyzer implements the quasi-enum type safety analyzer.
package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// detectByConstants implements DT-001: constants-based detection.
// Detects types with 2+ constants in the same package.
func detectByConstants(pass *analysis.Pass) map[*types.Named]bool {
	candidates := make(map[*types.Named]bool)

	// Count constants per type
	constantCounts := make(map[*types.Named]int)
	for _, name := range pass.Pkg.Scope().Names() {
		obj := pass.Pkg.Scope().Lookup(name)
		if c, ok := obj.(*types.Const); ok {
			if named, ok := c.Type().(*types.Named); ok {
				if isBasicType(named.Underlying()) {
					constantCounts[named]++
				}
			}
		}
	}

	// Mark types with 2+ constants
	for typ, count := range constantCounts {
		if count >= 2 {
			candidates[typ] = true
		}
	}

	return candidates
}

// detectBySuffix implements DT-002: name suffix detection.
// Detects types with name ending in "enum" (case-insensitive).
func detectBySuffix(pass *analysis.Pass) map[*types.Named]bool {
	candidates := make(map[*types.Named]bool)

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				// Check if name ends with enum keyword (case-insensitive)
				if strings.HasSuffix(strings.ToLower(typeSpec.Name.Name), strings.ToLower(enumKeyword)) {
					obj := pass.TypesInfo.Defs[typeSpec.Name]
					if named, ok := obj.Type().(*types.Named); ok {
						if isBasicType(named.Underlying()) {
							candidates[named] = true
						}
					}
				}
			}
		}
	}

	return candidates
}

// detectByInlineComment implements DT-003: inline comment detection.
// Detects types with inline comment starting with "enum".
func detectByInlineComment(pass *analysis.Pass) map[*types.Named]bool {
	candidates := make(map[*types.Named]bool)

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				// Check line comment
				if typeSpec.Comment != nil {
					for _, comment := range typeSpec.Comment.List {
						text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
						if startsWithEnumKeyword(text) {
							obj := pass.TypesInfo.Defs[typeSpec.Name]
							if named, ok := obj.Type().(*types.Named); ok {
								if isBasicType(named.Underlying()) {
									candidates[named] = true
								}
							}
						}
					}
				}
			}
		}
	}

	return candidates
}

// detectByPrecedingComment implements DT-004: preceding comment detection.
// Detects types with doc comment starting with "enum".
func detectByPrecedingComment(pass *analysis.Pass) map[*types.Named]bool {
	candidates := make(map[*types.Named]bool)

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			// Check doc comment
			if genDecl.Doc != nil {
				for _, comment := range genDecl.Doc.List {
					text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
					if startsWithEnumKeyword(text) {
						for _, spec := range genDecl.Specs {
							typeSpec, ok := spec.(*ast.TypeSpec)
							if !ok {
								continue
							}

							obj := pass.TypesInfo.Defs[typeSpec.Name]
							if named, ok := obj.Type().(*types.Named); ok {
								if isBasicType(named.Underlying()) {
									candidates[named] = true
								}
							}
						}
					}
				}
			}
		}
	}

	return candidates
}

// detectByNamedComment implements DT-005: named comment detection.
// Detects types with comment matching "TypeName enum" pattern.
func detectByNamedComment(pass *analysis.Pass) map[*types.Named]bool {
	candidates := make(map[*types.Named]bool)

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			if genDecl.Doc != nil {
				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					typeName := typeSpec.Name.Name
					for _, comment := range genDecl.Doc.List {
						text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
						if startsWithTypeNameEnumKeyword(text, typeName) {
							obj := pass.TypesInfo.Defs[typeSpec.Name]
							if named, ok := obj.Type().(*types.Named); ok {
								if isBasicType(named.Underlying()) {
									candidates[named] = true
								}
							}
						}
					}
				}
			}
		}
	}

	return candidates
}

// startsWithEnumKeyword checks if text starts with the configured enum keyword (case-insensitive).
func startsWithEnumKeyword(text string) bool {
	lower := strings.ToLower(strings.TrimSpace(text))
	keywordLower := strings.ToLower(enumKeyword)
	return strings.HasPrefix(lower, keywordLower+" ") || lower == keywordLower
}

// startsWithTypeNameEnumKeyword checks if text starts with "TypeName <keyword>" pattern.
func startsWithTypeNameEnumKeyword(text string, typeName string) bool {
	lower := strings.ToLower(strings.TrimSpace(text))
	pattern := strings.ToLower(typeName) + " " + strings.ToLower(enumKeyword)
	return strings.HasPrefix(lower, pattern)
}

// isBasicType checks if a type is a Go basic type.
func isBasicType(t types.Type) bool {
	_, ok := t.(*types.Basic)
	return ok
}

// detectQuasiEnums orchestrates all detection techniques.
func detectQuasiEnums(pass *analysis.Pass, config *DetectionConfig) map[*types.Named][]DetectionTechnique {
	allCandidates := make(map[*types.Named][]DetectionTechnique)

	// Apply enabled detection techniques
	if config.ConstantsDetectionEnabled {
		for named := range detectByConstants(pass) {
			allCandidates[named] = append(allCandidates[named], DT001ConstantsBased)
		}
	}

	if config.SuffixDetectionEnabled {
		for named := range detectBySuffix(pass) {
			allCandidates[named] = append(allCandidates[named], DT002NameSuffix)
		}
	}

	if config.InlineCommentDetectionEnabled {
		for named := range detectByInlineComment(pass) {
			allCandidates[named] = append(allCandidates[named], DT003InlineComment)
		}
	}

	if config.PrecedingCommentDetectionEnabled {
		for named := range detectByPrecedingComment(pass) {
			allCandidates[named] = append(allCandidates[named], DT004PrecedingComment)
		}
	}

	if config.NamedCommentDetectionEnabled {
		for named := range detectByNamedComment(pass) {
			allCandidates[named] = append(allCandidates[named], DT005NamedComment)
		}
	}

	return allCandidates
}
