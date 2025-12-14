package analyzer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Djarvur/go-enumsafety/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

// Existing tests for US1, US2, US3
func TestUS1_LiteralAssignment(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "..", "testdata")
	analysistest.Run(t, testdata, analyzer.Analyzer, "a")
}

func TestUS2_UntypedConstant(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "..", "testdata")
	analysistest.Run(t, testdata, analyzer.Analyzer, "a")
}

func TestUS3_VariableConversion(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "..", "testdata")
	analysistest.Run(t, testdata, analyzer.Analyzer, "a")
}

// New tests for US4, US5, US6
func TestUS4_Uint8Optimization(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "..", "testdata")
	analysistest.Run(t, testdata, analyzer.Analyzer, "optimization")
}

func TestUS5_StringMethod(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "..", "testdata")
	analysistest.Run(t, testdata, analyzer.Analyzer, "helpers")
}

func TestUS6_UnmarshalTextMethod(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "..", "testdata")
	analysistest.Run(t, testdata, analyzer.Analyzer, "helpers")
}

// New tests for comprehensive coverage

func TestConstraintViolations(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "..", "testdata")
	analysistest.Run(t, testdata, analyzer.Analyzer, "constraints_full")
}

func TestVarDeclarations(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "..", "testdata")
	analysistest.Run(t, testdata, analyzer.Analyzer, "vardecl")
}

func TestCompositeLiterals(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "..", "testdata")
	analysistest.Run(t, testdata, analyzer.Analyzer, "composite")
}

func TestDetectionEdgeCases(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "..", "testdata")
	analysistest.Run(t, testdata, analyzer.Analyzer, "detection_edge")
}

func TestCallExpressions(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join(wd, "..", "..", "testdata")
	analysistest.Run(t, testdata, analyzer.Analyzer, "call_expr")
}
