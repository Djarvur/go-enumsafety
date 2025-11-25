# Quickstart Guide: go-enumsafety Quasi-Enum Type Safety Linter

**Feature**: 001-enum-linter  
**Date**: 2025-11-23  
**Purpose**: Quick guide to using the go-enumsafety linter with quasi-enum detection and configuration

## What is go-enumsafety?

go-enumsafety is a Go linter that enforces type safety for quasi-enum patterns. It detects quasi-enums using 5 configurable techniques, enforces 5 definition constraints, and prevents invalid value assignments.

**Problem it solves**:
```go
type Status uint8
const (
    StatusActive Status = 1
    StatusInactive Status = 2
)

// Without go-enumsafety - compiles but defeats enum purpose
var s Status = 5  // ❌ Should use StatusActive or StatusInactive

// With go-enumsafety - caught at lint time
var s Status = 5  // ❌ ERROR: literal assignment to quasi-enum type not allowed
```

## Installation

### From Source (Development)

```bash
git clone https://github.com/Djarvur/go-enumsafety.git
cd go-enumsafety
go install ./cmd/enumsafety
```

### Using go install (Once Published)

```bash
go install github.com/Djarvur/go-enumsafety/cmd/enumsafety@latest
```

### Verify Installation

```bash
go-enumsafety -V
# Output: go-enumsafety version X.Y.Z
```

## Basic Usage

### Analyze a Single Package

```bash
go-enumsafety ./path/to/package
```

### Analyze Multiple Packages

```bash
go-enumsafety ./...
```

### Analyze Specific Files

```bash
go-enumsafety ./path/to/file.go
```

## Quasi-Enum Detection

go-enumsafety uses 5 detection techniques (all enabled by default):

### DT-001: Constants-Based Detection

Detects types with 2+ constants:

```go
type Status uint8  // Detected as quasi-enum
const (
    StatusActive Status = 1
    StatusInactive Status = 2
)
```

### DT-002: Name Suffix Detection

Detects types with "enum" suffix:

```go
type StatusEnum uint8  // Detected by suffix
const StatusActive StatusEnum = 1
```

### DT-003: Inline Comment Detection

Detects types with inline "enum" comment:

```go
type Status uint8 // enum for user status  // Detected by inline comment
```

### DT-004: Preceding Comment Detection

Detects types with preceding "enum" comment:

```go
// enum of valid statuses
type Status uint8  // Detected by preceding comment
```

### DT-005: Named Comment Detection

Detects types with "TypeName enum" comment:

```go
// Status enum for user states
type Status uint8  // Detected by named comment
```

## Definition Constraints

Once detected, quasi-enums must satisfy 5 constraints (all enabled by default):

### DC-001: Minimum 2 Constants

```go
// ❌ VIOLATION: Only 1 constant
type Status uint8
const StatusActive Status = 1

// ✅ OK: 2+ constants
type Status uint8
const (
    StatusActive Status = 1
    StatusInactive Status = 2
)
```

### DC-002: Same Const Block

```go
// ❌ VIOLATION: Different const blocks
const StatusActive Status = 1
// ... other code ...
const StatusInactive Status = 2

// ✅ OK: Same const block
const (
    StatusActive Status = 1
    StatusInactive Status = 2
)
```

### DC-003: Same File

```go
// ❌ VIOLATION: Type in file1.go, constants in file2.go

// ✅ OK: Type and constants in same file
```

### DC-004: Exclusive Const Block

```go
// ❌ VIOLATION: Mixed types
const (
    StatusActive Status = 1
    PriorityHigh Priority = 1  // Different type
)

// ✅ OK: Only Status constants
const (
    StatusActive Status = 1
    StatusInactive Status = 2
)
```

### DC-005: Proximity

```go
// ❌ VIOLATION: Code between type and const block
type Status uint8

var x = 5  // Non-comment code

const (
    StatusActive Status = 1
)

// ✅ OK: Only empty lines and comments
type Status uint8

// Valid constants below

const (
    StatusActive Status = 1
)
```

## Configuration Flags

### Disabling Detection Techniques

```bash
# Disable suffix detection
go-enumsafety -disable-suffix-detection ./...

# Disable all comment-based detection
go-enumsafety -disable-inline-comment-detection \
         -disable-preceding-comment-detection \
         -disable-named-comment-detection ./...

# Only use constants-based detection
go-enumsafety -disable-suffix-detection \
         -disable-inline-comment-detection \
         -disable-preceding-comment-detection \
         -disable-named-comment-detection ./...
```

### Disabling Definition Constraints

```bash
# Disable proximity check (allow code between type and const block)
go-enumsafety -disable-proximity-check ./...

# Disable exclusive block check (allow mixed-type const blocks)
go-enumsafety -disable-exclusive-block-check ./...

# Disable multiple constraints
go-enumsafety -disable-proximity-check \
         -disable-exclusive-block-check \
         -disable-same-file-check ./...
```

### Common Configuration Scenarios

**Strict Mode** (default - all checks enabled):
```bash
go-enumsafety ./...
```

**Relaxed Mode** (disable organizational constraints):
```bash
go-enumsafety -disable-proximity-check \
         -disable-exclusive-block-check \
         -disable-same-file-check ./...
```

**Legacy Codebase** (only detect by suffix, minimal constraints):
```bash
go-enumsafety -disable-constants-detection \
         -disable-inline-comment-detection \
         -disable-preceding-comment-detection \
         -disable-named-comment-detection \
         -disable-proximity-check \
         -disable-exclusive-block-check ./...
```

**Documentation-Driven** (only comment-based detection):
```bash
go-enumsafety -disable-constants-detection \
         -disable-suffix-detection ./...
```

## Integration with go vet

### Run as Part of go vet

```bash
go vet -vettool=$(which go-enumsafety) ./...
```

### With Configuration Flags

```bash
go vet -vettool=$(which go-enumsafety) \
       -enumsafety.disable-suffix-detection \
       -enumsafety.disable-proximity-check ./...
```

**Note**: Flags must be prefixed with `-enumsafety.` when using with `go vet`.

### Add to Makefile

```makefile
.PHONY: lint
lint:
\tgo vet -vettool=$(which go-enumsafety) \
\t       -enumsafety.disable-proximity-check ./...
```

## Output Formats

### Human-Readable (Default)

```bash
go-enumsafety ./...
```

**Usage Violation Output**:
```
models/status.go:15:5: literal assignment to quasi-enum type not allowed
  quasi-enum type: Status
  valid constants: StatusActive, StatusInactive, StatusPending
  suggestion: use Status constants instead of literal value
```

**Constraint Violation Output**:
```
models/status.go:5:1: quasi-enum definition constraint violated: Status
  constraint: DC-002 (same const block)
  problem: constants found in 2 different const blocks
  fix: move all Status constants into a single const block
```

### JSON Format (for Tooling)

```bash
go-enumsafety -json ./...
```

**Output**:
```json
[
  {
    "posn": "models/status.go:15:5",
    "message": "literal assignment to quasi-enum type not allowed\n  quasi-enum type: Status\n  valid constants: StatusActive, StatusInactive, StatusPending",
    "category": "enumsafety",
    "suggested_fixes": []
  }
]
```

## What go-enumsafety Catches

### ❌ Usage Violations

**Literal Assignment**:
```go
var s Status = 5  // ERROR
```

**Literal Type Conversion**:
```go
s := Status(5)  // ERROR
```

**Literal Function Arguments**:
```go
func SetStatus(s Status) {}
SetStatus(3)  // ERROR
```

**Untyped Constants**:
```go
const myValue = 5
var s Status = myValue  // ERROR
```

**Variable Conversions**:
```go
var x uint8 = 5
s := Status(x)  // ERROR
```

**Composite Literals**:
```go
statuses := []Status{1, 2, 3}  // ERROR
```

### ❌ Constraint Violations

**Too Few Constants** (DC-001):
```go
type Status uint8  // Detected as quasi-enum
const StatusActive Status = 1  // ERROR: only 1 constant
```

**Different Const Blocks** (DC-002):
```go
const StatusActive Status = 1
const StatusInactive Status = 2  // ERROR: different block
```

**Different Files** (DC-003):
```go
// file1.go: type Status uint8
// file2.go: const StatusActive Status = 1  // ERROR: different file
```

**Mixed Const Block** (DC-004):
```go
const (
    StatusActive Status = 1
    PriorityHigh Priority = 1  // ERROR: different type
)
```

**Code Between Type and Const** (DC-005):
```go
type Status uint8
var x = 5  // ERROR: non-comment code
const (StatusActive Status = 1)
```

## What go-enumsafety Allows

### ✅ Valid Patterns

**Enum Constant Assignment**:
```go
var s Status = StatusActive  // OK
s = StatusInactive           // OK
```

**Zero Value Initialization**:
```go
var s Status  // OK - defaults to 0
```

**Cross-Package Enum Constants**:
```go
import "example.com/models"
var s models.Status = models.StatusActive  // OK
```

**Enum Constants in Expressions**:
```go
if s == StatusActive {  // OK
    // ...
}

switch s {
case StatusActive:   // OK
case StatusInactive: // OK
}
```

## Best Practices

### 1. Choose Detection Technique

**Explicit Naming** (recommended for new code):
```go
// Status enum for user states
type StatusEnum uint8
const (
    StatusActive StatusEnum = 1
    StatusInactive StatusEnum = 2
)
```

**Constants-Based** (works with existing code):
```go
type Status uint8
const (
    StatusActive Status = 1
    StatusInactive Status = 2
)
```

### 2. Organize Enums in Dedicated Files

```
models/
├── status.go      # Status enum and constants
├── priority.go    # Priority enum and constants
└── permission.go  # Permission enum and constants
```

### 3. Document Enum Purpose

```go
// Status represents the current state of a user account.
// Status enum
type Status uint8

const (
    // StatusActive indicates an active account.
    StatusActive Status = 1
    // StatusInactive indicates a deactivated account.
    StatusInactive Status = 2
)
```

### 4. Use Zero Value Meaningfully

```go
// Good - zero value has meaning
type Status uint8
const (
    StatusUnknown Status = 0  // Default/uninitialized
    StatusActive Status = 1
)

var s Status  // OK - defaults to StatusUnknown
```

### 5. Configure for Your Codebase

Start strict, relax as needed:
```bash
# Start with all checks
go-enumsafety ./...

# If too strict, disable specific checks
go-enumsafety -disable-proximity-check ./...
```

## Troubleshooting

### "command not found: go-enumsafety"

**Solution**: Ensure `go-enumsafety` is installed and in PATH:
```bash
go install github.com/Djarvur/go-enumsafety/cmd/enumsafety@latest
which go-enumsafety
```

### "quasi-enum definition constraint violated"

**Solution**: Check which constraint is violated and fix:
- **DC-001**: Add more constants (minimum 2)
- **DC-002**: Move constants to same const block
- **DC-003**: Move constants to same file as type
- **DC-004**: Remove other type's constants from block
- **DC-005**: Remove code between type and const block

**Or disable the constraint**:
```bash
go-enumsafety -disable-proximity-check ./...
```

### Too Many False Positives

**Solution**: Disable detection techniques you don't use:
```bash
# Only use constants-based detection
go-enumsafety -disable-suffix-detection \
         -disable-inline-comment-detection \
         -disable-preceding-comment-detection \
         -disable-named-comment-detection ./...
```

### Performance Issues

**Solution**: go-enumsafety should complete in <100ms for typical files. If slow:
- Check file size (very large files may take longer)
- Ensure you're using the latest version
- File an issue with performance profile

## IDE Integration

### VS Code with gopls

Works automatically with gopls. Violations appear as inline diagnostics.

**Setup**: No additional configuration needed if `go-enumsafety` is in PATH.

### GoLand/IntelliJ IDEA

Configure as external tool or use `go vet` integration.

**Settings** → **Go** → **Vet** → Add `-vettool=$(which go-enumsafety)`

### Vim/Neovim with LSP

Works automatically with gopls LSP integration.

## Next Steps

- Read the [full specification](spec.md) for detailed requirements
- Review the [implementation plan](plan.md) for architecture details
- Check the [data model](data-model.md) for internal structures
- See the [analyzer contract](contracts/analyzer.md) for integration details

## Support

- **Issues**: https://github.com/Djarvur/go-enumsafety/issues
- **Discussions**: https://github.com/Djarvur/go-enumsafety/discussions
- **Contributing**: See CONTRIBUTING.md (once available)

## License

MIT License - see LICENSE file for details
