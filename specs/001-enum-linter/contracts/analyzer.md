# Analyzer Contract

**Feature**: 001-enum-linter  
**Date**: 2025-11-23  
**Purpose**: Define the analyzer interface, behavior contract, and configuration for quasi-enum type safety

## Overview

This document specifies the contract for the quasi-enum type safety analyzer, including its interface, behavior, inputs, outputs, configuration flags, and integration points.

## Analyzer Specification

### Basic Information

**Name**: `enumsafety`  
**Doc**: `check that quasi-enum types are only assigned their defined constants and satisfy definition constraints`  
**URL**: `https://github.com/Djarvur/go-enumsafety`  
**Requires**: `[]*analysis.Analyzer{inspect.Analyzer}`

### Analyzer Interface

```go
var Analyzer = &analysis.Analyzer{
    Name:     "enumsafety",
    Doc:      "check that quasi-enum types are only assigned their defined constants and satisfy definition constraints",
    URL:      "https://github.com/Djarvur/go-enumsafety",
    Requires: []*analysis.Analyzer{inspect.Analyzer},
    Run:      run,
    Flags:    *flagSet(),
}

func run(pass *analysis.Pass) (interface{}, error)
```

### Configuration Flags

The analyzer MUST support 10 command-line flags for configuring detection and constraints:

#### Detection Technique Flags

**Flag**: `-disable-constants-detection`  
**Type**: `bool`  
**Default**: `false`  
**Description**: Disable DT-001 (constants-based detection)  
**Effect**: Types will not be detected as quasi-enums based on having 2+ constants

**Flag**: `-disable-suffix-detection`  
**Type**: `bool`  
**Default**: `false`  
**Description**: Disable DT-002 (name suffix detection)  
**Effect**: Types with "enum" suffix will not be detected as quasi-enums

**Flag**: `-disable-inline-comment-detection`  
**Type**: `bool`  
**Default**: `false`  
**Description**: Disable DT-003 (inline comment detection)  
**Effect**: Types with inline "enum" comment will not be detected as quasi-enums

**Flag**: `-disable-preceding-comment-detection`  
**Type**: `bool`  
**Default**: `false`  
**Description**: Disable DT-004 (preceding comment detection)  
**Effect**: Types with preceding "enum" comment will not be detected as quasi-enums

**Flag**: `-disable-named-comment-detection`  
**Type**: `bool`  
**Default**: `false`  
**Description**: Disable DT-005 (named comment detection)  
**Effect**: Types with "TypeName enum" comment will not be detected as quasi-enums

#### Definition Constraint Flags

**Flag**: `-disable-min-constants-check`  
**Type**: `bool`  
**Default**: `false`  
**Description**: Disable DC-001 (minimum 2 constants check)  
**Effect**: Quasi-enums with fewer than 2 constants will not be reported as violations

**Flag**: `-disable-same-block-check`  
**Type**: `bool`  
**Default**: `false`  
**Description**: Disable DC-002 (same const block check)  
**Effect**: Quasi-enums with constants in different blocks will not be reported as violations

**Flag**: `-disable-same-file-check`  
**Type**: `bool`  
**Default**: `false`  
**Description**: Disable DC-003 (same file check)  
**Effect**: Quasi-enums with type and constants in different files will not be reported as violations

**Flag**: `-disable-exclusive-block-check`  
**Type**: `bool`  
**Default**: `false`  
**Description**: Disable DC-004 (exclusive const block check)  
**Effect**: Quasi-enums with mixed-type const blocks will not be reported as violations

**Flag**: `-disable-proximity-check`  
**Type**: `bool`  
**Default**: `false`  
**Description**: Disable DC-005 (proximity check)  
**Effect**: Quasi-enums with code between type and const block will not be reported as violations

### Input Contract

**Receives from analysis.Pass**:
- `Files`: `[]*ast.File` - AST of all files in the package
- `Fset`: `*token.FileSet` - File set for position information
- `Pkg`: `*types.Package` - Type-checked package
- `TypesInfo`: `*types.Info` - Type information for expressions
- `ResultOf[inspect.Analyzer]`: `*inspector.Inspector` - AST inspector

**Preconditions**:
- Package MUST be successfully type-checked (no type errors)
- All imports MUST be resolved
- Type information MUST be complete

### Output Contract

**Returns**:
- `interface{}`: Always `nil` (analyzer produces diagnostics, not results)
- `error`: Non-nil only for internal analyzer errors (not for detected violations)

**Produces via analysis.Pass.Report()**:
- `analysis.Diagnostic` for each violation found (usage + constraint)
- Diagnostics include:
  - `Pos`: Exact position of violation
  - `Message`: Human-readable error message
  - `Category`: Violation type (for tooling)
  - `SuggestedFixes`: Optional fixes (future enhancement)

**Postconditions**:
- All quasi-enum usage violations in package are reported (based on enabled detection)
- All quasi-enum constraint violations are reported (based on enabled constraints)
- No false positives (100% precision per SC-003, SC-005, SC-007)
- No false negatives for enabled techniques (100% recall per SC-001, SC-004, SC-006)

## Behavior Contract

### Quasi-Enum Detection

**MUST detect when enabled** (each technique can be disabled):

**DT-001: Constants-Based** (default enabled):
- Named types derived from basic types
- With 2 or more constants in same package

**DT-002: Name Suffix** (default enabled):
- Named types derived from basic types
- With name ending in "enum" (case-insensitive)
- Examples: `StatusEnum`, `PriorityENUM`, `Colorenum`

**DT-003: Inline Comment** (default enabled):
- Named types derived from basic types
- With inline comment starting with "enum" (case-insensitive)
- Example: `type Status uint8 // enum for user status`

**DT-004: Preceding Comment** (default enabled):
- Named types derived from basic types
- With doc comment starting with "enum" (case-insensitive)

**DT-005: Named Comment** (default enabled):
- Named types derived from basic types
- With doc comment matching "TypeName enum" pattern (case-insensitive)

**MUST NOT detect**:
- Types with zero constants (unless other technique matches)
- Struct types, interface types, function types
- Type aliases (using `type X = Y` syntax)
- Types when all detection techniques are disabled

### Constraint Validation

**MUST enforce when enabled** (each constraint can be disabled):

**DC-001: Minimum Constants** (default enabled):
- At least 2 constants required for detected quasi-enum
- Violation if fewer than 2 constants

**DC-002: Same Const Block** (default enabled):
- All constants must be in same const block
- Violation if constants scattered across multiple blocks

**DC-003: Same File** (default enabled):
- Const block must be in same file as type definition
- Violation if type and constants in different files

**DC-004: Exclusive Const Block** (default enabled):
- No other type's constants in the quasi-enum's const block
- Violation if mixed-type const block

**DC-005: Proximity** (default enabled):
- Only empty lines and comments between type and const block
- Violation if other code (functions, variables, types) in between

### Usage Violation Detection

**MUST detect** (always enabled, not configurable):

1. **Literal Assignment**:
   ```go
   var s Status = 5  // VIOLATION
   ```

2. **Literal Conversion**:
   ```go
   s := Status(5)  // VIOLATION
   ```

3. **Literal Function Argument**:
   ```go
   func SetStatus(s Status) {}
   SetStatus(3)  // VIOLATION
   ```

4. **Untyped Constant**:
   ```go
   const myVal = 5
   var s Status = myVal  // VIOLATION
   ```

5. **Variable Conversion**:
   ```go
   var x uint8 = 5
   s := Status(x)  // VIOLATION
   ```

6. **Composite Literal**:
   ```go
   statuses := []Status{1, 2, 3}  // VIOLATION
   ```

**MUST NOT report violations**:

1. **Valid Constant Assignment**:
   ```go
   var s Status = StatusActive  // OK
   ```

2. **Zero Value**:
   ```go
   var s Status  // OK - defaults to 0
   ```

3. **Constant from Same Quasi-Enum**:
   ```go
   s := StatusActive  // OK
   ```

4. **Cross-Package Quasi-Enum Constants**:
   ```go
   import "other/pkg"
   var s pkg.Status = pkg.StatusActive  // OK
   ```

### Error Message Contract

#### Usage Violation Format

```
<violation-description>
  quasi-enum type: <TypeName>
  valid constants: <Constant1>, <Constant2>, ...
  suggestion: <actionable-fix>
```

**Example**:
```
literal assignment to quasi-enum type not allowed
  quasi-enum type: Status
  valid constants: StatusActive, StatusInactive, StatusPending
  suggestion: use Status constants instead of literal value
```

#### Constraint Violation Format

```
quasi-enum definition constraint violated: <TypeName>
  constraint: <ConstraintID> (<constraint-name>)
  problem: <specific-issue>
  fix: <actionable-fix>
```

**Example**:
```
quasi-enum definition constraint violated: Status
  constraint: DC-002 (same const block)
  problem: constants found in 2 different const blocks
  fix: move all Status constants into a single const block
```

**Requirements**:
- MUST include quasi-enum type name
- MUST list valid constants (up to 10, then "... and N more") for usage violations
- MUST identify specific constraint for constraint violations
- MUST provide actionable suggestion
- MUST be clear and concise

### Performance Contract

**Requirements** (SC-012):
- MUST complete analysis in <100ms for files under 1000 lines
- MUST complete analysis in <1s for typical Go packages
- SHOULD use memory proportional to analyzed code size
- MUST NOT leak memory between package analyses

**Implementation Constraints**:
- Parallel detection technique execution
- Lazy constraint validation (only for detected quasi-enums)
- Efficient type lookups using maps
- Early exit if all detection disabled

## Integration Contracts

### CLI Integration

**Command Line**:
```bash
go-enumsafety [flags] <packages>
```

**Standard Flags** (provided by singlechecker):
- `-json`: Output diagnostics in JSON format
- `-c=N`: Display offending line with N lines of context
- `-V`: Print analyzer version and exit

**Custom Flags** (10 quasi-enum flags):
- `-disable-constants-detection`
- `-disable-suffix-detection`
- `-disable-inline-comment-detection`
- `-disable-preceding-comment-detection`
- `-disable-named-comment-detection`
- `-disable-min-constants-check`
- `-disable-same-block-check`
- `-disable-same-file-check`
- `-disable-exclusive-block-check`
- `-disable-proximity-check`

**Exit Codes**:
- `0`: No violations found
- `1`: Violations found
- `2`: Analysis error (internal error, not violations)

**Output Format** (human-readable):
```
path/to/file.go:10:5: literal assignment to quasi-enum type not allowed
  quasi-enum type: Status
  valid constants: StatusActive, StatusInactive

path/to/file.go:5:1: quasi-enum definition constraint violated: Status
  constraint: DC-002 (same const block)
  problem: constants found in 2 different const blocks
```

**Output Format** (JSON with `-json` flag):
```json
{
  "posn": "path/to/file.go:10:5",
  "message": "literal assignment to quasi-enum type not allowed\n  quasi-enum type: Status\n  valid constants: StatusActive, StatusInactive",
  "category": "enumsafety",
  "suggested_fixes": []
}
```

**Flag Combination Examples**:
```bash
# Disable suffix detection
go-enumsafety -disable-suffix-detection ./...

# Disable multiple constraints
go-enumsafety -disable-proximity-check -disable-exclusive-block-check ./...

# Disable all comment-based detection
go-enumsafety -disable-inline-comment-detection \
         -disable-preceding-comment-detection \
         -disable-named-comment-detection ./...
```

### go vet Integration

**Usage**:
```bash
go vet -vettool=$(which go-enumsafety) ./...
```

**With Flags**:
```bash
go vet -vettool=$(which go-enumsafety) \
       -enumsafety.disable-suffix-detection \
       -enumsafety.disable-proximity-check ./...
```

**Requirements** (FR-026, SC-014):
- MUST work as `-vettool` for `go vet`
- MUST follow `go vet` conventions
- MUST integrate seamlessly with other vet checks
- Flags MUST be prefixed with analyzer name (`-enumsafety.`)

### IDE Integration

**Compatibility** (via analysis framework):
- VS Code with gopls
- GoLand/IntelliJ IDEA
- Vim/Neovim with LSP
- Emacs with LSP

**Behavior**:
- Real-time diagnostics as code is written
- Inline error messages
- Quick fixes (future enhancement)

## Testing Contract

### Unit Test Requirements

**Coverage** (Constitution IV):
- >80% code coverage for `internal/analyzer/` package
- 100% coverage of detection techniques
- 100% coverage of constraint validation
- 100% coverage of usage violation patterns

**Test Structure**:
```go
func TestAnalyzer(t *testing.T) {
    testdata := analysistest.TestData()
    analysistest.Run(t, testdata, Analyzer, "a")
}

func TestDetectionTechniques(t *testing.T) {
    // Test each technique individually
}

func TestConstraints(t *testing.T) {
    // Test each constraint individually
}

func TestFlags(t *testing.T) {
    // Test flag combinations
}
```

**Test Fixtures** (in `internal/testdata/src/a/`):
- `enum.go`: Valid quasi-enum patterns (all 5 detection techniques)
- `violations.go`: Usage violations with `// want` comments
- `constraints.go`: Constraint violations with `// want` comments
- `iota.go`: Iota-based enums
- `expressions.go`: Expression-based constants
- `flags/disabled_dt.go`: Tests with detection disabled
- `flags/disabled_dc.go`: Tests with constraints disabled
- `crosspackage/types.go`: Cross-package enum usage

### Integration Test Requirements

**CLI Tests**:
- Test standalone execution
- Test with various flags
- Test flag combinations
- Test exit codes

**go vet Tests**:
- Test `-vettool` integration
- Test combined with other vet checks
- Test flag passing via `-enumsafety.` prefix

## Versioning Contract

**Semantic Versioning** (Constitution, Development Workflow):
- **MAJOR**: Breaking changes to detection logic, constraint logic, or CLI interface
- **MINOR**: New detection techniques, new constraints, new flags (backward compatible)
- **PATCH**: Bug fixes, documentation, performance improvements

**Backward Compatibility**:
- Detection techniques MUST remain stable within MAJOR version
- Constraints MUST remain stable within MAJOR version
- New techniques/constraints MAY be added in MINOR versions (disabled by default if breaking)
- Error message format MAY change in MINOR versions (not breaking)

## Error Handling Contract

**Internal Errors**:
- Return `error` from `run()` function
- Log diagnostic information
- Do NOT crash on malformed code

**Type Errors**:
- Skip analysis if package has type errors
- Report via `pass.Reportf()` if needed
- Do NOT report quasi-enum violations in packages with type errors

**Edge Cases**:
- Handle nil pointers gracefully
- Handle missing type information
- Handle circular dependencies
- Handle comments with special characters

## Compliance Summary

This contract ensures:
- ✅ All functional requirements (FR-001 through FR-030) are met
- ✅ All success criteria (SC-001 through SC-016) are verifiable
- ✅ Constitution principles are upheld
- ✅ Standard Go tooling integration is guaranteed
- ✅ Performance targets are specified and measurable
- ✅ Configuration is flexible and well-defined
