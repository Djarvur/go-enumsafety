// Package analyzer implements the quasi-enum type safety analyzer.
package analyzer

import (
	"flag"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Configuration flags for detection techniques
var (
	disableConstantsDetection        bool
	disableSuffixDetection           bool
	disableInlineCommentDetection    bool
	disablePrecedingCommentDetection bool
	disableNamedCommentDetection     bool
)

// Configuration flags for definition constraints
var (
	disableMinConstantsCheck   bool
	disableSameBlockCheck      bool
	disableSameFileCheck       bool
	disableExclusiveBlockCheck bool
	disableProximityCheck      bool
)

// Configuration flags for quality-of-life checks (US4-US6)
var (
	disableUint8Suggestion      bool
	disableStringMethodCheck    bool
	disableUnmarshalMethodCheck bool
)

// Configuration for detection keyword customization (FR-070, FR-131)
var enumKeyword string

// Analyzer is the quasi-enum type safety analyzer.
var Analyzer = &analysis.Analyzer{
	Name:     "enumsafety",
	Doc:      "check that quasi-enum types are only assigned their defined constants and satisfy definition constraints",
	URL:      "https://github.com/Djarvur/go-enumsafety",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
	Flags:    makeFlags(),
}

// makeFlags creates and returns a flag.FlagSet with all analyzer flags.
func makeFlags() flag.FlagSet {
	var fs flag.FlagSet

	// Detection technique flags
	fs.BoolVar(&disableConstantsDetection, "disable-constants-detection", false,
		"disable DT-001: constants-based detection")
	fs.BoolVar(&disableSuffixDetection, "disable-suffix-detection", false,
		"disable DT-002: name suffix detection")
	fs.BoolVar(&disableInlineCommentDetection, "disable-inline-comment-detection", false,
		"disable DT-003: inline comment detection")
	fs.BoolVar(&disablePrecedingCommentDetection, "disable-preceding-comment-detection", false,
		"disable DT-004: preceding comment detection")
	fs.BoolVar(&disableNamedCommentDetection, "disable-named-comment-detection", false,
		"disable DT-005: named comment detection")

	// Definition constraint flags
	fs.BoolVar(&disableMinConstantsCheck, "disable-min-constants-check", false,
		"disable DC-001: minimum 2 constants check")
	fs.BoolVar(&disableSameBlockCheck, "disable-same-block-check", false,
		"disable DC-002: same const block check")
	fs.BoolVar(&disableSameFileCheck, "disable-same-file-check", false,
		"disable DC-003: same file check")
	fs.BoolVar(&disableExclusiveBlockCheck, "disable-exclusive-block-check", false,
		"disable DC-004: exclusive const block check")
	fs.BoolVar(&disableProximityCheck, "disable-proximity-check", false,
		"disable DC-005: proximity check")

	// Quality-of-life check flags (US4-US6)
	fs.BoolVar(&disableUint8Suggestion, "disable-uint8-suggestion", false,
		"disable US4: uint8 optimization suggestion")
	fs.BoolVar(&disableStringMethodCheck, "disable-string-method-check", false,
		"disable US5: String() method check")
	fs.BoolVar(&disableUnmarshalMethodCheck, "disable-unmarshal-method-check", false,
		"disable US6: UnmarshalText() method check")

	// Keyword customization flag (FR-070, FR-131)
	fs.StringVar(&enumKeyword, "enum-keyword", "enum",
		"customize the detection keyword (default: 'enum')")

	return fs
}

// run is the main analyzer entry point.
func run(pass *analysis.Pass) (interface{}, error) {
	// Check if all detection techniques are disabled
	if disableConstantsDetection && disableSuffixDetection &&
		disableInlineCommentDetection && disablePrecedingCommentDetection &&
		disableNamedCommentDetection {
		// Report error and exit with code 2 (configuration error)
		pass.Reportf(0, "all detection techniques disabled - no quasi-enums will be detected")
		return nil, flag.ErrHelp // Signals configuration error
	}

	// Create configuration
	detectionConfig := NewDetectionConfig()
	constraintConfig := NewConstraintConfig()

	// Step 1: Detect quasi-enum types
	detectedTypes := detectQuasiEnums(pass, detectionConfig)
	if len(detectedTypes) == 0 {
		// No quasi-enums detected, nothing to do
		return nil, nil
	}

	// Step 2: Build QuasiEnumRegistry
	registry := NewQuasiEnumRegistry(detectionConfig, constraintConfig)

	// For each detected type, collect constants and build QuasiEnumType
	for namedType, techniques := range detectedTypes {
		qe := buildQuasiEnumType(pass, namedType, techniques)
		if qe != nil {
			registry.RegisterQuasiEnum(qe)
		}
	}

	// Step 3: Validate definition constraints
	for _, qe := range registry.QuasiEnums {
		violations := qe.ValidateConstraints(
			constraintConfig,
			pass.Fset,
			qe.TypeDecl,
			qe.ConstBlock,
			qe.File,
			pass.TypesInfo,
		)

		// Report constraint violations as warnings
		for _, violation := range violations {
			reportConstraintViolation(pass, qe, violation)
		}
	}

	// Step 4: Check for usage violations (US1, US2, US3)
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.AssignStmt:
				checkAssignment(pass, registry, node)
			case *ast.GenDecl:
				checkVarDecl(pass, registry, node)
			case *ast.CallExpr:
				checkCallExpr(pass, registry, node)
			case *ast.CompositeLit:
				checkCompositeLit(pass, registry, node)
			}
			return true
		})
	}

	// Step 5: Check for quality-of-life improvements (US4, US5, US6)
	checkUint8Optimization(pass, registry)
	checkStringMethod(pass, registry)
	checkUnmarshalTextMethod(pass, registry)

	return nil, nil
}
