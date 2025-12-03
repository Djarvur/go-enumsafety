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
