# Implementation Plan: Enum Type Safety Linter

**Branch**: `001-enum-linter` | **Date**: 2025-11-25 | **Spec**: [spec.md](./spec.md)

## Summary

Build a Go static analysis linter that enforces type safety for quasi-enum patterns by detecting enum-like types through multiple techniques (name suffix, comments, constants) and preventing usage of literals, untyped constants, or type conversions instead of defined enum constants. The linter integrates with `golang.org/x/tools/go/analysis` framework for standard Go tooling compatibility and provides quality-of-life features including uint8 optimization suggestions and helper method warnings.

## Technical Context

**Language/Version**: Go 1.22+  
**Primary Dependencies**: `golang.org/x/tools/go/analysis`, `golang.org/x/tools/go/analysis/passes/inspect`, `golang.org/x/tools/go/analysis/analysistest`  
**Storage**: N/A (stateless analysis tool)  
**Testing**: `analysistest` framework from `golang.org/x/tools`  
**Target Platform**: All platforms supported by Go 1.22+ (Linux, macOS, Windows, BSD)  
**Project Type**: Standalone CLI tool (single binary)  
**Performance Goals**: <100ms for files under 1000 lines, <1s for typical packages  
**Constraints**: Zero false positives/negatives for core detection, memory proportional to analyzed code  
**Scale/Scope**: Analyze Go codebases of any size, 5 detection techniques, 5 definition constraints, 14 configuration flags

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Principle I: Maintainability-First
✅ **PASS** - Design uses standard `golang.org/x/tools/go/analysis` framework with clear separation of concerns:
- Detection logic in `internal/analyzer/detection.go`
- Constraint validation in `internal/analyzer/constraints.go`
- Usage violation checks in `internal/analyzer/usage.go`
- Quality-of-life checks in separate modules (`optimization.go`, `helpers.go`)
- Comprehensive documentation required for all exported functions

### Principle II: Accuracy
✅ **PASS** - Specification defines:
- 100% accuracy targets for all detection techniques (SC-001, SC-003, SC-005, SC-007)
- Zero false positives/negatives requirements
- Precise matching algorithms for all detection techniques
- Comprehensive edge case handling (15+ scenarios documented)
- All 70 functional requirements are testable and measurable

### Principle III: Developer Experience
✅ **PASS** - Design prioritizes DX:
- Clear, actionable error messages with suggestions (FR-024, SC-013)
- Autofix capability for uint8 optimization (US4)
- 14 configuration flags for customization
- Performance target: <100ms for <1000 lines (SC-012)
- Integration with standard `go vet` tooling (SC-014)

### Principle IV: Pragmatic Testing
✅ **PASS** - Testing strategy defined:
- Unit tests for all 5 detection techniques (FR-036)
- Unit tests for all 5 definition constraints (FR-037)
- Unit tests for all 6 user stories (FR-038 to FR-041)
- Uses `analysistest` framework (constitution-compliant)
- Test fixtures for all scenarios
- Target: >80% coverage for core packages

### Principle V: Simplicity
✅ **PASS** - Design follows YAGNI:
- Minimal dependencies (only `golang.org/x/tools` and stdlib)
- No premature optimization (relies on analysis framework for concurrency)
- Clear feature scope (enum detection and validation only)
- Pragmatic reliance on Go ecosystem (compiler, stdlib, framework)
- No library API guarantees (internal packages)

**Constitution Compliance**: ✅ **ALL GATES PASSED**

## Project Structure

### Documentation (this feature)

```text
specs/001-enum-linter/
├── spec.md              # Complete specification (490 lines, 70 FRs)
├── plan.md              # This file
├── research.md          # Phase 0 output (technology decisions)
├── data-model.md        # Phase 1 output (core data structures)
├── quickstart.md        # Phase 1 output (getting started guide)
├── contracts/           # Phase 1 output (analyzer contract)
│   └── analyzer.md      # Analyzer interface specification
├── tasks.md             # Phase 2 output (implementation tasks)
└── checklists/          # Quality validation
    ├── quality.md       # Requirements quality checklist (110/110 PASS)
    ├── requirements.md  # Requirements completeness checklist
    └── enhancements.md  # Enhancement quality checklist
```

### Source Code (repository root)

```text
go-enumsafety/
├── cmd/
│   └── go-enumsafety/
│       └── main.go           # CLI entry point
├── internal/
│   └── analyzer/
│       ├── analyzer.go       # Main analyzer with flag registration
│       ├── detection.go      # 5 detection techniques (DT-001 to DT-005)
│       ├── constraints.go    # 5 definition constraints (DC-001 to DC-005)
│       ├── usage.go          # Usage violation checks (US1, US2, US3)
│       ├── optimization.go   # uint8 optimization (US4)
│       ├── helpers.go        # Helper method checks (US5, US6)
│       ├── enum.go           # QuasiEnumType data structure
│       └── violation.go      # Violation reporting
├── tests/
│   └── unit/
│       └── analyzer_test.go  # Unit tests using analysistest
└── internal/testdata/
    └── src/
        ├── a/                # Test fixtures for user stories
        ├── detection/        # Test fixtures for detection techniques
        ├── constraints/      # Test fixtures for constraints
        ├── optimization/     # Test fixtures for US4
        └── helpers/          # Test fixtures for US5, US6
```

## Phase 0: Research & Technology Decisions

**Status**: ✅ COMPLETE

### Research Topics

1. **Go Analysis Framework Best Practices**
   - **Decision**: Use `golang.org/x/tools/go/analysis` standard framework
   - **Rationale**: 
     - Official Go static analysis infrastructure
     - Automatic integration with `go vet`
     - Robust AST traversal and type information
     - Standard testing framework (`analysistest`)
   - **Alternatives Considered**: Custom AST walker (rejected: reinventing wheel, no tooling integration)

2. **Detection Strategy**
   - **Decision**: Multiple detection techniques (5 total) with OR logic
   - **Rationale**:
     - Flexibility for different coding styles
     - Explicit techniques (suffix, comments) checked first for performance
     - Implicit technique (constants-based) as fallback
     - Opt-out mechanism via `// not enum` comment
   - **Alternatives Considered**: Single detection method (rejected: too rigid)

3. **Constraint Validation Approach**
   - **Decision**: 5 independent constraints with AND logic (all enabled must pass)
   - **Rationale**:
     - Ensures enum definitions follow best practices
     - Each constraint independently disableable
     - Clear violation messages for each constraint
   - **Alternatives Considered**: Hardcoded validation (rejected: not flexible)

4. **Testing Framework**
   - **Decision**: `analysistest` from `golang.org/x/tools`
   - **Rationale**:
     - Standard testing approach for Go analyzers
     - Automatic test fixture discovery
     - Built-in diagnostic verification
     - Constitution-compliant (Principle IV)
   - **Alternatives Considered**: Custom test harness (rejected: unnecessary complexity)

5. **Configuration Mechanism**
   - **Decision**: Command-line flags via `analysis.Analyzer.Flags`
   - **Rationale**:
     - Standard Go analysis flag mechanism
     - Automatic help text generation
     - Integration with `go vet` and other tools
     - 14 total flags for fine-grained control
   - **Alternatives Considered**: Config file (rejected: adds complexity, not standard for analyzers)

6. **Performance Optimization**
   - **Decision**: Check explicit detection techniques before implicit (FR-048)
   - **Rationale**:
     - Explicit techniques (suffix, comments) are fast string operations
     - Implicit technique (constants-based) requires full package analysis
     - Reduces unnecessary work for most cases
   - **Alternatives Considered**: Random order (rejected: performance impact)

**Output**: research.md with all technology decisions documented

## Phase 1: Design & Contracts

**Status**: IN PROGRESS

### Core Data Model

**Primary Entities**:

1. **QuasiEnumType**
   - Type name and underlying Go type
   - Detection techniques that matched (DT-001 to DT-005)
   - List of enum constants
   - Flags: HasStringMethod, HasUnmarshalTextMethod, SuggestEnumComment
   - Position information for reporting

2. **EnumConstant**
   - Constant name and value
   - Type reference
   - Position in source code

3. **Violation**
   - Violation type (usage vs constraint)
   - Position and message
   - Suggested fix (if applicable)

**Relationships**:
- QuasiEnumType has many EnumConstants (1:N)
- Violations reference QuasiEnumType (N:1)

### API Contracts

**Analyzer Interface** (`contracts/analyzer.md`):

```go
// Analyzer is the main analysis entry point
var Analyzer = &analysis.Analyzer{
    Name: "enumlinter",
    Doc:  "Enforces type safety for quasi-enum patterns in Go",
    Run:  run,
    Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// Configuration flags (14 total)
// Detection technique flags (5)
var disableConstantsDetection bool
var disableSuffixDetection bool
var disableInlineCommentDetection bool
var disablePrecedingCommentDetection bool
var disableNamedCommentDetection bool

// Constraint flags (5)
var disableMinConstantsCheck bool
var disableSameBlockCheck bool
var disableSameFileCheck bool
var disableExclusiveBlockCheck bool
var disableProximityCheck bool

// Quality-of-life flags (3)
var disableUint8Suggestion bool
var disableStringMethodCheck bool
var disableUnmarshalMethodCheck bool

// Keyword configuration (1)
var enumKeyword string // default: "enum"
```

**Detection Interface**:
```go
// detectQuasiEnums identifies all quasi-enum types in the package
func detectQuasiEnums(pass *analysis.Pass) []*QuasiEnumType

// Detection technique functions
func detectByConstants(pass *analysis.Pass) []*QuasiEnumType  // DT-001
func detectBySuffix(pass *analysis.Pass) []*QuasiEnumType     // DT-002
func detectByInlineComment(pass *analysis.Pass) []*QuasiEnumType    // DT-003
func detectByPrecedingComment(pass *analysis.Pass) []*QuasiEnumType // DT-004
func detectByNamedComment(pass *analysis.Pass) []*QuasiEnumType     // DT-005
func hasOptOutComment(typeSpec *ast.TypeSpec) bool  // FR-046
```

**Constraint Interface**:
```go
// validateConstraints checks all enabled constraints for a quasi-enum
func validateConstraints(qe *QuasiEnumType, pass *analysis.Pass) []Violation

// Constraint validation functions
func checkMinConstants(qe *QuasiEnumType) *Violation      // DC-001
func checkSameBlock(qe *QuasiEnumType) *Violation         // DC-002
func checkSameFile(qe *QuasiEnumType) *Violation          // DC-003
func checkExclusiveBlock(qe *QuasiEnumType) *Violation    // DC-004
func checkProximity(qe *QuasiEnumType) *Violation         // DC-005
```

**Usage Validation Interface**:
```go
// checkUsageViolations scans for improper enum usage
func checkUsageViolations(qe *QuasiEnumType, pass *analysis.Pass) []Violation

// Violation detection functions
func checkLiteralAssignment(pass *analysis.Pass) []Violation       // FR-015
func checkLiteralArgument(pass *analysis.Pass) []Violation         // FR-016
func checkTypeConversion(pass *analysis.Pass) []Violation          // FR-017
func checkUntypedConstant(pass *analysis.Pass) []Violation         // FR-019
func checkVariableConversion(pass *analysis.Pass) []Violation      // FR-020
func checkCompositeLiteral(pass *analysis.Pass) []Violation        // FR-021
```

### Quickstart Guide

**Installation**:
```bash
go install github.com/Djarvur/go-enumsafety/cmd/enumsafety@latest
```

**Basic Usage**:
```bash
# Analyze current package
go-enumsafety ./...

# With go vet
go vet -vettool=$(which go-enumsafety) ./...

# Disable specific checks
go-enumsafety -disable-suffix-detection -disable-uint8-suggestion ./...

# Custom detection keyword
go-enumsafety -enum-keyword=enumeration ./...
```

**Example Enum Pattern**:
```go
// Status enum for user states
type Status uint8

const (
    StatusActive   Status = 1
    StatusInactive Status = 2
    StatusPending  Status = 3
)

// ✅ Valid usage
var s Status = StatusActive

// ❌ Invalid - linter will report error
var s Status = 5  // literal assignment
```

## Phase 2: Implementation Breakdown

**Note**: Detailed task breakdown will be generated by `/speckit.tasks` command.

### High-Level Implementation Phases

**Phase 2.1: Core Infrastructure** (Estimated: 2-3 days)
- Set up analyzer structure with flag registration
- Implement QuasiEnumType data structure
- Create violation reporting mechanism
- Set up test infrastructure with `analysistest`

**Phase 2.2: Detection Techniques** (Estimated: 3-4 days)
- Implement DT-002 (name suffix detection)
- Implement DT-003 (inline comment detection)
- Implement DT-004 (preceding comment detection)
- Implement DT-005 (named comment detection)
- Implement DT-001 (constants-based detection)
- Implement FR-046 (opt-out mechanism)
- Implement FR-047 (suggest enum comment)
- Implement FR-048 (performance optimization)
- Unit tests for all detection techniques

**Phase 2.3: Definition Constraints** (Estimated: 2-3 days)
- Implement DC-001 (minimum constants)
- Implement DC-002 (same const block)
- Implement DC-003 (same file)
- Implement DC-004 (exclusive const block)
- Implement DC-005 (proximity)
- Implement FR-051 (multiple violation reporting)
- Unit tests for all constraints

**Phase 2.4: Usage Violation Detection** (Estimated: 3-4 days)
- Implement literal assignment detection (US1)
- Implement untyped constant detection (US2)
- Implement variable conversion detection (US3)
- Implement FR-049 (nested violation handling)
- Unit tests for all usage violations

**Phase 2.5: Quality-of-Life Features** (Estimated: 2 days)
- Implement uint8 optimization suggestion (US4) with autofix
- Implement String() method check (US5)
- Implement UnmarshalText() method check (US6)
- Unit tests for all QoL features

**Phase 2.6: Configuration & Integration** (Estimated: 1-2 days)
- Implement all 14 command-line flags
- Implement FR-070 (enum-keyword customization)
- Implement FR-050 (flag independence)
- Integration tests for flag combinations

**Phase 2.7: Documentation & Polish** (Estimated: 1-2 days)
- Complete README with examples
- Create CONTRIBUTING guide
- Add usage examples
- Performance optimization if needed

**Total Estimated Effort**: 14-20 days

## Verification Plan

### Automated Tests

1. **Unit Tests** (using `analysistest`):
   - All 5 detection techniques (FR-036)
   - All 5 definition constraints (FR-037)
   - All 6 user stories (FR-038 to FR-041)
   - Edge cases (15+ scenarios)
   - Target: >80% coverage

2. **Integration Tests**:
   - CLI flag combinations (SC-008, SC-009)
   - go vet integration (SC-014)
   - Performance benchmarks (SC-012)

3. **Regression Tests**:
   - Test fixtures for all reported bugs
   - Continuous validation of accuracy targets

### Manual Verification

1. **Real-World Testing**:
   - Run on existing Go codebases with enum patterns
   - Verify zero false positives/negatives
   - Validate error message clarity

2. **Performance Testing**:
   - Benchmark on files of various sizes
   - Verify <100ms target for <1000 lines
   - Memory profiling

3. **Usability Testing**:
   - Developer feedback on error messages
   - CLI ergonomics validation
   - Documentation completeness review

### Success Criteria Validation

All 16 success criteria (SC-001 to SC-016) must pass:
- ✅ 100% detection accuracy for all techniques
- ✅ Zero false positives/negatives
- ✅ All 14 flags work correctly
- ✅ Performance targets met
- ✅ Error messages are actionable
- ✅ go vet integration works

## Risk Assessment

### Technical Risks

1. **AST Complexity** (Medium Risk)
   - **Mitigation**: Use `golang.org/x/tools/go/analysis` framework which handles AST complexity
   - **Fallback**: Extensive test fixtures to validate edge cases

2. **Performance on Large Codebases** (Low Risk)
   - **Mitigation**: FR-048 (check explicit techniques first), rely on analysis framework
   - **Fallback**: Performance profiling and optimization if needed

3. **False Positives/Negatives** (Medium Risk)
   - **Mitigation**: Comprehensive test suite, precise matching algorithms
   - **Fallback**: Regression tests for all reported issues

### Project Risks

1. **Scope Creep** (Low Risk)
   - **Mitigation**: Constitution Principle V (Simplicity), clear feature scope in spec
   - **Fallback**: Defer non-essential features to future versions

2. **Maintenance Burden** (Low Risk)
   - **Mitigation**: Constitution Principle I (Maintainability), comprehensive documentation
   - **Fallback**: Community contribution guidelines

## Dependencies

### External Dependencies

- `golang.org/x/tools/go/analysis` - Core analysis framework
- `golang.org/x/tools/go/analysis/passes/inspect` - AST inspection
- `golang.org/x/tools/go/analysis/singlechecker` - CLI runner
- `golang.org/x/tools/go/analysis/analysistest` - Testing framework

All dependencies are from `golang.org/x` (semi-official Go packages), meeting constitution dependency policy.

### Internal Dependencies

None - standalone tool with no library dependencies.

## Next Steps

1. ✅ **Phase 0 Complete**: Research and technology decisions documented
2. 🔄 **Phase 1 In Progress**: Design artifacts being generated
   - ✅ data-model.md
   - ✅ contracts/analyzer.md
   - ✅ quickstart.md
3. ⏭️ **Phase 2 Next**: Run `/speckit.tasks` to generate detailed implementation tasks

## Appendix

### Glossary Reference

See [spec.md §Glossary](./spec.md#glossary) for comprehensive definitions of:
- Quasi-Enum Type, Enum Constant, Base Type
- Detection Technique (DT), Definition Constraint (DC)
- Violation types, Thread-Safety, Registry
- And 15+ more technical terms

### Quality Checklist Status

- ✅ **Requirements Quality**: 110/110 items passing (100%)
- ✅ **Constitution Compliance**: All 5 principles satisfied
- ✅ **Specification Completeness**: 70 FRs, 16 SCs, comprehensive edge cases

---

**Plan Version**: 1.0.0 | **Generated**: 2025-11-25 | **Status**: Ready for Implementation
