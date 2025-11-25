# enumsafety

A Go static analysis linter that enforces type safety for quasi-enum patterns, compensating for Go's lack of native enum support.

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue)](https://go.dev/)
[![CI](https://github.com/Djarvur/enumsafety/workflows/CI/badge.svg)](https://github.com/Djarvur/enumsafety/actions)
[![Coverage Status](https://coveralls.io/repos/github/Djarvur/enumsafety/badge.svg?branch=main)](https://coveralls.io/github/Djarvur/enumsafety?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/Djarvur/enumsafety)](https://goreportcard.com/report/github.com/Djarvur/enumsafety)
[![GoDoc](https://pkg.go.dev/badge/github.com/Djarvur/go-enumsafety)](https://pkg.go.dev/github.com/Djarvur/go-enumsafety)

## Overview

Go doesn't have native enum types. Developers typically emulate enums using type aliases with constants:

```go
type Status int

const (
    StatusActive   Status = iota
    StatusInactive
    StatusPending
)
```

However, Go's type system allows unsafe assignments that bypass enum safety:

```go
var s Status = 5              // ❌ Should only use defined constants
s = Status(someVariable)      // ❌ Bypasses enum safety
const myVal = 3
var s2 Status = myVal         // ❌ Untyped constant, not part of enum
```

**enumsafety** detects and prevents these violations, ensuring quasi-enum types are only assigned their defined constant values.

## Features

✅ **Literal Detection** - Catches direct literal assignments (`var s Status = 5`)  
✅ **Untyped Constant Detection** - Prevents use of constants not part of enum definition  
✅ **Variable Conversion Detection** - Blocks conversions from variables to enum types  
✅ **Cross-Enum Detection** - Prevents conversions between different enum types  
✅ **Flexible Detection** - 5 configurable techniques to identify quasi-enums  
✅ **Definition Validation** - 5 constraints to ensure proper enum structure  
✅ **Quality-of-Life Checks** - Suggests uint8 optimization, String(), and UnmarshalText() methods  
✅ **Configurable** - 14 command-line flags for fine-grained control  
✅ **go vet Integration** - Works seamlessly with standard Go tooling

## Installation

```bash
go install github.com/Djarvur/enumsafety/cmd/enumsafety@latest
```

## Quick Start

### 1. Define Your Enum

```go
// Status enum
type Status int

const (
    StatusActive   Status = iota
    StatusInactive
    StatusPending
)
```

### 2. Run the Linter

```bash
# Analyze current package
enumsafety ./...

# Or use with go vet
go vet -vettool=$(which enumsafety) ./...
```

### 3. Fix Violations

```go
// ❌ Before
var s Status = 5

// ✅ After
var s Status = StatusActive
```

## Detection Techniques

The linter identifies quasi-enums using multiple techniques (all enabled by default):

### 1. Constants-Based (DT-001)
Detects types with 2+ constants of the same type:
```go
type Priority uint8
const (
    PriorityLow  Priority = 1
    PriorityHigh Priority = 2
)
```

### 2. Name Suffix (DT-002)
Detects types with "enum" suffix (case-insensitive):
```go
type StatusEnum int
```

### 3. Inline Comment (DT-003)
Detects types with inline `// enum` comment:
```go
type Color uint8 // enum
```

### 4. Preceding Comment (DT-004)
Detects types with preceding `// enum` comment:
```go
// enum
type Level int
```

### 5. Named Comment (DT-005)
Detects types with `// TypeName enum` pattern:
```go
// Status enum
type Status int
```

### Opt-Out Mechanism
Prevent detection with `// not enum` comment:
```go
type NotAnEnum int // not enum
```

## Definition Constraints

The linter validates that detected quasi-enums follow best practices:

### DC-001: Minimum 2 Constants
Enums must have at least 2 constants.

### DC-002: Same Const Block
All enum constants must be in the same `const` block.

### DC-003: Same File
Type and constants must be in the same file.

### DC-004: Exclusive Block
The const block should only contain constants of the enum type.

### DC-005: Proximity
Type declaration and const block should be close together (within 10 lines).

## Violation Detection

### Literal Assignment (US1)
```go
var s Status = 5  // ❌ Error: literal value assigned to quasi-enum type Status
```

### Untyped Constant (US2)
```go
const myValue = 3
var s Status = myValue  // ❌ Error: untyped constant assigned to quasi-enum type Status
```

### Variable Conversion (US3)
```go
var x uint8 = 5
s := Status(x)  // ❌ Error: variable converted to quasi-enum type Status
```

### Cross-Enum Conversion
```go
var c Color = ColorRed
l := Level(c)  // ❌ Error: variable converted to quasi-enum type Level
```

## Quality-of-Life Features

### uint8 Optimization (US4)
Suggests using `uint8` for memory efficiency:
```go
type Status int  // ⚠️ Suggestion: use uint8 (only 3 constants)
```

### String() Method (US5)
Warns about missing `String()` method:
```go
type Status int  // ⚠️ Warning: lacks String() method
```

### UnmarshalText() Method (US6)
Warns about missing `UnmarshalText()` for JSON/config parsing:
```go
type Status int  // ⚠️ Warning: lacks UnmarshalText([]byte) error method
```

## Configuration Flags

### Detection Technique Flags

Disable specific detection techniques:

```bash
-disable-constants-detection        # Disable DT-001
-disable-suffix-detection           # Disable DT-002
-disable-inline-comment-detection   # Disable DT-003
-disable-preceding-comment-detection # Disable DT-004
-disable-named-comment-detection    # Disable DT-005
```

### Definition Constraint Flags

Disable specific constraints:

```bash
-disable-min-constants-check     # Disable DC-001 (minimum 2 constants)
-disable-same-block-check        # Disable DC-002 (same const block)
-disable-same-file-check         # Disable DC-003 (same file)
-disable-exclusive-block-check   # Disable DC-004 (exclusive block)
-disable-proximity-check         # Disable DC-005 (proximity)
```

### Quality-of-Life Flags

Disable quality-of-life checks:

```bash
-disable-uint8-suggestion        # Disable uint8 optimization suggestions
-disable-string-method-check     # Disable String() method warnings
-disable-unmarshal-method-check  # Disable UnmarshalText() warnings
```

### Keyword Customization

Customize the detection keyword (default: "enum"):

```bash
-enum-keyword=enumeration  # Use "enumeration" instead of "enum"
```

## Usage Examples

### Basic Analysis

```bash
# Analyze current package
enumsafety ./...

# Analyze specific package
enumsafety ./internal/models

# Analyze with verbose output
enumsafety -v ./...
```

### With go vet

```bash
# Use as vet tool
go vet -vettool=$(which enumsafety) ./...

# Combine with other vet checks
go vet -vettool=$(which enumsafety) -all ./...
```

### Custom Configuration

```bash
# Disable uint8 suggestions
enumsafety -disable-uint8-suggestion ./...

# Only check for violations, skip quality checks
enumsafety \
  -disable-uint8-suggestion \
  -disable-string-method-check \
  -disable-unmarshal-method-check \
  ./...

# Use custom keyword
enumsafety -enum-keyword=enumeration ./...

# Minimal checks (only constants-based detection)
enumsafety \
  -disable-suffix-detection \
  -disable-inline-comment-detection \
  -disable-preceding-comment-detection \
  -disable-named-comment-detection \
  ./...
```

## Example Output

Given this code:

```go
package main

// Status enum
type Status int

const (
    StatusActive   Status = iota
    StatusInactive
    StatusPending
)

func main() {
    var s Status = 5  // Violation
}
```

The linter reports:

```
main.go:12:6: literal value assigned to quasi-enum type Status
main.go:4:6: quasi-enum type Status uses int but has only 3 constants; consider using uint8 for memory optimization
main.go:4:6: quasi-enum type Status lacks a String() method; consider using golang.org/x/tools/cmd/stringer or github.com/Djarvur/go-silly-enum to generate it
main.go:4:6: quasi-enum type Status lacks an UnmarshalText([]byte) error method; consider using github.com/Djarvur/go-silly-enum to generate it
```

## Best Practices

### 1. Use Explicit Detection Comments

Make enum intent clear:

```go
// Status enum for user states
type Status uint8

const (
    StatusActive   Status = 1
    StatusInactive Status = 2
)
```

### 2. Implement Helper Methods

Add `String()` and `UnmarshalText()` for better DX:

```go
func (s Status) String() string {
    switch s {
    case StatusActive:
        return "Active"
    case StatusInactive:
        return "Inactive"
    default:
        return "Unknown"
    }
}

func (s *Status) UnmarshalText(text []byte) error {
    switch string(text) {
    case "Active":
        *s = StatusActive
    case "Inactive":
        *s = StatusInactive
    default:
        return fmt.Errorf("unknown status: %s", text)
    }
    return nil
}
```

### 3. Use uint8 for Small Enums

Optimize memory usage:

```go
type Status uint8  // ✅ Better than int for 2-256 values
```

### 4. Keep Enums Simple

Follow the constraints:
- Minimum 2 constants
- All constants in same block
- Type and constants in same file
- Exclusive const block
- Close proximity

## Integration with CI/CD

### GitHub Actions

```yaml
name: Lint

on: [push, pull_request]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      - name: Install enumsafety
        run: go install github.com/Djarvur/enumsafety/cmd/enumsafety@latest
      - name: Run linter
        run: enumsafety ./...
```

### Pre-commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

enumsafety ./...
if [ $? -ne 0 ]; then
    echo "enumsafety found violations. Please fix before committing."
    exit 1
fi
```

## FAQ

### Q: Why use this instead of just being careful?

**A**: Type safety catches bugs at compile time. enumsafety enforces enum safety that Go's type system doesn't provide natively.

### Q: Can I use this with existing codebases?

**A**: Yes! Start by running with all quality checks disabled, then gradually enable them:

```bash
enumsafety \
  -disable-uint8-suggestion \
  -disable-string-method-check \
  -disable-unmarshal-method-check \
  ./...
```

### Q: What if I have a type that looks like an enum but isn't?

**A**: Use the opt-out comment:

```go
type NotAnEnum int // not enum
```

### Q: Does this work with iota?

**A**: Yes! The linter fully supports iota-based enums.

### Q: Can I customize the detection keyword?

**A**: Yes! Use `-enum-keyword=yourword` to use a different keyword than "enum".

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Related Tools

- [golang.org/x/tools/cmd/stringer](https://pkg.go.dev/golang.org/x/tools/cmd/stringer) - Generate String() methods
- [github.com/Djarvur/go-silly-enum](https://github.com/Djarvur/go-silly-enum) - Generate enum helper methods

## Acknowledgments

This project was developed using:
- **AI Assistant**: Google Gemini 2.0 Flash (Experimental)
- **Development Framework**: SpecKit - AI-assisted software specification and implementation toolkit

Built with [golang.org/x/tools/go/analysis](https://pkg.go.dev/golang.org/x/tools/go/analysis) framework.
