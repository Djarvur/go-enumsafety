# go-enumsafety Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-12-04

## Active Technologies
- N/A (stateless analysis tool) (001-enum-linter)
- Go 1.22+ + `golang.org/x/tools/go/analysis`, `golang.org/x/tools/go/analysis/passes/inspect`, `golang.org/x/tools/go/analysis/analysistest` (001-enum-linter)
- Go 1.22+ + `golang.org/x/tools/go/analysis` (standard Go analysis framework) (001-enum-linter)

## Project Structure

```text
internal/
  analyzer/         # Core analyzer implementation
  testdata/         # Test data for analysistest
    src/
      a/            # US1, US2, US3 violation tests
      constraints/  # DC-001 through DC-005 tests
      constraints_full/  # Comprehensive constraint violation tests
      detection/    # DT-001 through DT-005 tests
      detection_edge/    # Detection edge case tests
      helpers/      # US5, US6 helper method tests
      optimization/ # US4 uint8 optimization tests
      vardecl/      # Variable declaration edge cases
      composite/    # Composite literal edge cases
      call_expr/    # Function call expression edge cases
tests/
  unit/            # Unit tests
cmd/
  enumsafety/      # CLI entry point
```

## Commands

```bash
# Run tests with coverage
go test -coverprofile=coverage.out ./analyzer
go tool cover -func=coverage.out

# Run all tests
go test ./...

# Build the linter
go build ./cmd/enumsafety

# Run the linter
go vet -vettool=$(which enumsafety) ./...
```

## Code Style

Go 1.22+: Follow standard conventions

## Test Coverage

Current coverage: **81.0%** (target: ≥75%)

### Test Packages
- `a`: US1, US2, US3 violation tests
- `constraints_full`: Comprehensive constraint validation tests
- `vardecl`: Variable declaration edge cases
- `composite`: Composite literal edge cases
- `detection_edge`: Detection technique variations
- `call_expr`: Function call expression edge cases
- `optimization`: US4 uint8 optimization tests
- `helpers`: US5, US6 helper method tests
- `detection`: DT-001 through DT-005 detection tests
- `constraints`: DC-001 through DC-005 constraint tests

## Recent Changes
- 2025-12-04: Increased test coverage from 74.3% to 81.0% by adding 5 new comprehensive test packages
- 2025-12-04: Refactored message formatting to use `fmt.Sprintf()` instead of manual string manipulation
- 2025-11-22: Added Go 1.22+ + `golang.org/x/tools/go/analysis`, `golang.org/x/tools/go/analysis/passes/inspect`, `golang.org/x/tools/go/analysis/analysistest`
- 2025-11-22: Added Go 1.22+ + `golang.org/x/tools/go/analysis` (standard Go analysis framework)

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
