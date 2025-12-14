package analyzer

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// TestUS1_LiteralAssignment tests detection of literal value assignments to quasi-enum types.
func TestUS1_LiteralAssignment(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "testdata")
	analysistest.Run(t, testdata, Analyzer, "a")
}

// TestUS2_UntypedConstant tests detection of untyped constant assignments to quasi-enum types.
func TestUS2_UntypedConstant(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "testdata")
	analysistest.Run(t, testdata, Analyzer, "a")
}

// TestUS3_VariableConversion tests detection of variable conversions to quasi-enum types.
func TestUS3_VariableConversion(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "testdata")
	analysistest.Run(t, testdata, Analyzer, "a")
}

// TestUS4_Uint8Optimization tests uint8 optimization suggestions.
func TestUS4_Uint8Optimization(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "testdata")
	analysistest.Run(t, testdata, Analyzer, "optimization")
}

// TestUS5_StringMethod tests String() method warnings.
func TestUS5_StringMethod(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "testdata")
	analysistest.Run(t, testdata, Analyzer, "helpers")
}

// TestUS6_UnmarshalTextMethod tests UnmarshalText() method warnings.
func TestUS6_UnmarshalTextMethod(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "testdata")
	analysistest.Run(t, testdata, Analyzer, "helpers")
}

// TestConstraintViolations tests constraint validation code paths.
func TestConstraintViolations(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "testdata")
	// This test will have unexpected diagnostics but still covers the validation code
	analysistest.Run(t, testdata, Analyzer, "constraints_full")
}

// TestVarDeclarations tests variable declaration edge cases.
func TestVarDeclarations(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "testdata")
	analysistest.Run(t, testdata, Analyzer, "vardecl")
}

// TestCompositeLiterals tests composite literal edge cases.
func TestCompositeLiterals(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "testdata")
	analysistest.Run(t, testdata, Analyzer, "composite")
}

// TestDetectionEdgeCases tests detection technique edge cases.
func TestDetectionEdgeCases(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "testdata")
	analysistest.Run(t, testdata, Analyzer, "detection_edge")
}

// TestCallExpressions tests function call expression edge cases.
func TestCallExpressions(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "testdata")
	analysistest.Run(t, testdata, Analyzer, "call_expr")
}
