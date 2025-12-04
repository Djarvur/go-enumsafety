// Package main provides the CLI entry point for go-enumsafety.
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/Djarvur/go-enumsafety/internal/analyzer"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
