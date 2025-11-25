# Feature Specification: Enum Type Safety Linter

**Feature Branch**: `001-enum-linter`  
**Created**: 2025-11-22  
**Status**: Draft  
**Input**: User description: "I'm building the linter for golang source code which must compensate the lack of standard enum type in golang. To emulate enums we are creating types aliasing the standard ones, mostly uint8, with a bunch of constants of this types providing all the possible values this type can accept. The problem is we should not be able to use the values other the defined constants for this type, like literals, untyped constants, type assigned values and all the other stuff like this. So we need a linter which will prohibit to init the variables and func call params with something but the predefined constants of the proper type."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Detect Literal Assignment to Quasi-Enum Types (Priority: P1)

As a Go developer using enum-like types, I want the linter to catch when I accidentally assign a literal value (like `5`) to a quasi-enum variable instead of using the defined constants, so that I maintain type safety and prevent runtime bugs.

**Why this priority**: This is the most common mistake when working with Go enum patterns. Developers often write `status := 5` instead of `status := StatusActive`, which bypasses the entire purpose of the enum pattern. This is the core value proposition of the linter.

**Independent Test**: Can be fully tested by creating a Go file with an enum type and constants, then attempting to assign a literal value. The linter should report an error. This delivers immediate value as an MVP.

**Acceptance Scenarios**:

1. **Given** a Go file with a quasi-enum type `type Status uint8` and constants `const (StatusActive Status = 1; StatusInactive Status = 2)`, **When** a developer writes `var s Status = 5`, **Then** the linter reports an error indicating literal assignment to quasi-enum type
2. **Given** the same enum definition, **When** a developer writes `var s Status = StatusActive`, **Then** the linter accepts this as valid
3. **Given** the same enum definition, **When** a developer writes `s := Status(5)`, **Then** the linter reports an error for explicit type conversion from literal
4. **Given** a function `func SetStatus(s Status)`, **When** called with `SetStatus(3)`, **Then** the linter reports an error for literal argument

---

### User Story 2 - Detect Untyped Constant Assignment (Priority: P2)

As a Go developer, I want the linter to prevent assignment of untyped constants (that aren't the defined enum constants) to quasi-enum variables, so that only explicitly defined enum values can be used.

**Why this priority**: After catching literal assignments, the next common mistake is using untyped constants defined elsewhere in the codebase. This is more subtle than direct literals but equally problematic.

**Independent Test**: Can be tested by creating an enum type, defining an untyped constant outside the enum definition, and attempting to assign it. The linter should flag this while allowing the defined enum constants.

**Acceptance Scenarios**:

1. **Given** an enum type `type Priority uint8` with constants `const (PriorityLow Priority = 1; PriorityHigh Priority = 2)` and a separate untyped constant `const myValue = 3`, **When** a developer writes `var p Priority = myValue`, **Then** the linter reports an error
2. **Given** the same setup, **When** a developer writes `var p Priority = PriorityLow`, **Then** the linter accepts this as valid
3. **Given** a typed constant of the wrong type `const otherValue Status = 5`, **When** assigned to a `Priority` variable, **Then** the linter reports a type mismatch error

---

### User Story 3 - Detect Variable Assignment to Quasi-Enum Types (Priority: P3)

As a Go developer, I want the linter to catch when I assign a variable of the underlying type (like `uint8`) to a quasi-enum variable, so that I can't bypass enum safety through variable indirection.

**Why this priority**: This is a less common but still important case. Developers might try to work around type safety by using variables of the underlying type. This should be caught to maintain complete enum integrity.

**Independent Test**: Can be tested by creating an enum type, a variable of the underlying type, and attempting assignment. The linter should flag this pattern.

**Acceptance Scenarios**:

1. **Given** an enum type `type Color uint8` with constants, and a variable `var x uint8 = 5`, **When** a developer writes `var c Color = Color(x)`, **Then** the linter reports an error for conversion from variable
2. **Given** the same setup, **When** a developer writes `var c Color = ColorRed`, **Then** the linter accepts this as valid
3. **Given** two different enum types with the same underlying type, **When** converting between them, **Then** the linter reports an error

---

### User Story 4 - Suggest uint8 Optimization (Priority: P4)

As a Go developer, I want the linter to suggest using `uint8` as the base type for quasi-enums with fewer than 256 constants when using a larger type, so that I can optimize memory usage.

**Why this priority**: Memory optimization is important but not critical for correctness. This is a helpful suggestion that can improve code quality.

**Independent Test**: Can be tested by creating an enum with `int` or `uint` base type and fewer than 256 constants. The linter should suggest using `uint8`.

**Acceptance Scenarios**:

1. **Given** an enum type `type Status int` with 3 constants, **When** the linter analyzes the code, **Then** it suggests using `uint8` instead of `int` with autofix capability
2. **Given** an enum type `type Color uint8` with 3 constants, **When** the linter analyzes the code, **Then** it does NOT suggest any change
3. **Given** an enum type `type Flags int` with 300 constants, **When** the linter analyzes the code, **Then** it does NOT suggest using `uint8`
4. **Given** the flag `-disable-uint8-suggestion` is set, **When** the linter analyzes the code, **Then** it does NOT suggest uint8 optimization

---

### User Story 5 - Warn Missing String() Method (Priority: P4)

As a Go developer, I want the linter to warn me when a quasi-enum type doesn't have a `String()` method, so that I can provide human-readable string representations.

**Why this priority**: String() methods are very useful for debugging and logging but not required for type safety. This is a quality-of-life improvement.

**Independent Test**: Can be tested by creating an enum without a String() method. The linter should warn and suggest using stringer tools.

**Acceptance Scenarios**:

1. **Given** an enum type `type Status int` without a `String()` method, **When** the linter analyzes the code, **Then** it warns and suggests using `golang.org/x/tools/cmd/stringer` or `github.com/Djarvur/go-silly-enum`
2. **Given** an enum type `type Color uint8` with a `String() string` method, **When** the linter analyzes the code, **Then** it does NOT warn
3. **Given** the flag `-disable-string-method-check` is set, **When** the linter analyzes the code, **Then** it does NOT warn about missing String() method

---

### User Story 6 - Warn Missing UnmarshalText() Method (Priority: P4)

As a Go developer, I want the linter to warn me when a quasi-enum type doesn't have an `UnmarshalText([]byte) error` method, so that I can properly handle JSON/text unmarshaling.

**Why this priority**: UnmarshalText is important for JSON/config parsing but not required for basic type safety. This helps ensure enums work correctly with serialization.

**Independent Test**: Can be tested by creating an enum without an UnmarshalText method. The linter should warn and suggest using go-silly-enum.

**Acceptance Scenarios**:

1. **Given** an enum type `type Priority uint8` without an `UnmarshalText([]byte) error` method, **When** the linter analyzes the code, **Then** it warns and suggests using `github.com/Djarvur/go-silly-enum`
2. **Given** an enum type `type Status int` with an `UnmarshalText([]byte) error` method, **When** the linter analyzes the code, **Then** it does NOT warn
3. **Given** the flag `-disable-unmarshal-method-check` is set, **When** the linter analyzes the code, **Then** it does NOT warn about missing UnmarshalText() method

---

### Edge Cases

**Detection Edge Cases**:
- What happens if a type has both "enum" and "not enum" comments? → "not enum" takes precedence (opt-out)
- What if a type is detected by multiple techniques? → All techniques that detected it are recorded
- What happens with multi-line comments containing "enum"? → First line only is checked for "enum" keyword
- How does case-insensitivity work for different locales (e.g., Turkish İ vs I)? → ASCII case-insensitive comparison (English locale)
- What if all detection techniques are disabled via flags? → Linter reports error and exits with code 2 (configuration error)
- What if a type is detected only by DT-001 (constants-based)? → Linter suggests adding "// enum" comment for explicit documentation
- In what order are detection techniques checked? → DT-002 through DT-005 (explicit) first, then DT-001 (implicit) for performance
- How does "not enum" opt-out work? → Inline comment `// not enum` prevents detection by any technique

**Constraint Edge Cases**:
- What happens when an enum constant is defined using an expression (e.g., `const StatusActive Status = 1 << 0`)?
- How does the linter handle iota-based enum definitions (e.g., `const (StatusA Status = iota; StatusB; StatusC)`)?
- What happens with zero values (e.g., `var s Status` which defaults to `0`)?
- How are enum types used in composite literals handled (e.g., `[]Status{1, 2, 3}` vs `[]Status{StatusA, StatusB}`)?
- What happens when enum values are used in switch statements with literal cases?
- How does the linter handle enum types passed through interfaces or type assertions?
- What happens with enum constants defined in different packages?
- What constitutes "proximity" - how many empty lines are allowed between type and const block?
- What if there are multiple const blocks in the same file - which one should contain the enum constants?
- How are const blocks with mixed types but same underlying type handled?
- What if a type is detected by one technique but violates constraints - is it still treated as quasi-enum? → YES, linter reports BOTH constraint violation AND enforces usage checks
- How are disabled detection techniques and constraints combined? → Detection techniques are OR-ed (any one can detect); constraints are AND-ed (all enabled must pass)
- How are nested violations handled (e.g., `SetStatus(Status(5))`)? → Report only the first (innermost) violation to avoid redundant messages
- What happens when a type violates multiple constraints simultaneously? → All violations are reported independently (FR-051)

## Requirements *(mandatory)*

### Functional Requirements

**Quasi-Enum Detection**:
- **FR-001**: Linter MUST implement DT-001 (constants-based detection): detect types with 2+ constants
- **FR-002**: Linter MUST implement DT-002 (name suffix detection): detect types with "enum" suffix (case-insensitive)
- **FR-003**: Linter MUST implement DT-003 (inline comment detection): detect types with inline "enum" comment
- **FR-004**: Linter MUST implement DT-004 (preceding comment detection): detect types with preceding "enum" comment
- **FR-005**: Linter MUST implement DT-005 (named comment detection): detect types with "TypeName enum" comment
- **FR-006**: All detection techniques MUST be enabled by default
- **FR-007**: Each detection technique MUST be individually disableable via command-line flag
- **FR-046**: Linter MUST NOT detect types with inline "not enum" comment (opt-out mechanism)
- **FR-047**: Linter MUST suggest adding inline "// enum" comment when type is detected only by DT-001 (constants-based)
- **FR-048**: Linter MUST check explicit detection techniques (DT-002 through DT-005) before DT-001 for performance optimization

**Definition Constraints**:
- **FR-008**: Linter MUST enforce DC-001 (minimum 2 constants) for detected quasi-enums
- **FR-009**: Linter MUST enforce DC-002 (same const block) for detected quasi-enums
- **FR-010**: Linter MUST enforce DC-003 (same file) for detected quasi-enums
- **FR-011**: Linter MUST enforce DC-004 (exclusive const block) for detected quasi-enums
- **FR-012**: Linter MUST enforce DC-005 (proximity) for detected quasi-enums
- **FR-013**: All definition constraints MUST be enabled by default
- **FR-014**: Each definition constraint MUST be individually disableable via command-line flag
- **FR-051**: When a type violates multiple constraints simultaneously, linter MUST report all violations independently (not just the first one)

**Usage Violations**:
- **FR-015**: Linter MUST detect direct literal assignment to quasi-enum type variables (e.g., `var s Status = 5`)
- **FR-016**: Linter MUST detect literal arguments passed to functions expecting quasi-enum types (e.g., `SetStatus(3)`)
- **FR-017**: Linter MUST detect explicit type conversions from literals to quasi-enum types (e.g., `Status(5)`)
- **FR-018**: Linter MUST allow assignment of defined enum constants to quasi-enum variables (e.g., `var s Status = StatusActive`)
- **FR-019**: Linter MUST detect assignment of untyped constants (not part of enum definition) to quasi-enum variables
- **FR-020**: Linter MUST detect type conversions from variables of underlying type to quasi-enum type (e.g., `Color(x)` where `x` is `uint8`)
- **FR-021**: Linter MUST detect literals in composite literals for quasi-enum slices/arrays (e.g., `[]Status{1, 2}`)
- **FR-049**: For nested violations (e.g., literal in composite literal in function call), linter MUST report only the first (innermost) discovered violation to avoid redundant error messages
- **FR-092**: Linter MUST apply violation checks to all nested contexts including composite literals within function calls and chained method calls
- **FR-093**: Linter MUST apply violation checks to method call chains involving quasi-enum types

**General Requirements**:
- **FR-022**: Linter MUST support iota-based enum definitions
- **FR-023**: Linter MUST handle enum constants defined using expressions (e.g., bit flags with `1 << 0`)
- **FR-024**: Linter MUST provide clear, actionable error messages indicating which enum constant should be used OR which constraint is violated
- **FR-025**: Linter MUST support quasi-enum types across multiple packages (detection only, constants must be in same package as type)
- **FR-026**: Linter MUST integrate with `golang.org/x/tools/go/analysis` framework for standard tooling compatibility
- **FR-027**: Linter MUST suggest using `uint8` as base type for quasi-enums with non-uint8 base type and fewer than 256 constants
- **FR-028**: uint8 optimization suggestion MUST provide autofix capability to change the base type
- **FR-029**: uint8 optimization check MUST be disableable via `-disable-uint8-suggestion` flag
- **FR-114**: Linter MUST suggest uint8 even when enum has exactly 256 constants (256 fits in uint8: 0-255)
- **FR-115**: Linter MUST suggest uint8 when current type is int8 (if all values are non-negative)
- **FR-116**: Linter MUST NOT suggest uint8 when enum uses negative constant values
- **FR-030**: Linter MUST warn when a quasi-enum type lacks a `String() string` method
- **FR-031**: String() method warning MUST suggest using `golang.org/x/tools/cmd/stringer` or `github.com/Djarvur/go-silly-enum`
- **FR-032**: String() method check MUST be disableable via `-disable-string-method-check` flag
- **FR-117**: Linter MUST warn when String() method has non-standard signature (not `String() string`)
- **FR-033**: Linter MUST warn when a quasi-enum type lacks an `UnmarshalText([]byte) error` method
- **FR-034**: UnmarshalText() method warning MUST suggest using `github.com/Djarvur/go-silly-enum`
- **FR-035**: UnmarshalText() method check MUST be disableable via `-disable-unmarshal-method-check` flag
- **FR-118**: Linter MUST warn when UnmarshalText() method has non-standard signature (not `UnmarshalText([]byte) error`)
- **FR-119**: Linter MUST ignore String() method defined on pointer receiver (only value receiver counts)

**Testing Requirements**:
- **FR-036**: Linter MUST have unit tests for all 5 detection techniques (DT-001 through DT-005)
- **FR-037**: Linter MUST have unit tests for all 5 definition constraints (DC-001 through DC-005)
- **FR-038**: Linter MUST have unit tests for all usage violation checks (US1, US2, US3)
- **FR-039**: Linter MUST have unit tests for uint8 optimization suggestion (US4)
- **FR-040**: Linter MUST have unit tests for String() method check (US5)
- **FR-041**: Linter MUST have unit tests for UnmarshalText() method check (US6)

**Error Handling Requirements**:
- **FR-052**: Linter MUST ignore lines with syntax errors (rely on Go compiler for syntax validation)
- **FR-053**: Linter MUST ignore type checking failures (rely on Go compiler for type validation)
- **FR-054**: Linter MUST fail completely and exit with non-zero code on file I/O errors
- **FR-055**: Linter MUST NOT implement graceful degradation - exit immediately on fatal errors
- **FR-056**: Linter MAY produce partial analysis results when some files in a package fail (non-fatal errors)

**Boundary Condition Requirements**:
- **FR-057**: Linter MUST NOT impose limits on number of constants beyond Go language limitations
- **FR-058**: Linter MUST handle empty files and packages without error (no quasi-enums detected is valid)
- **FR-059**: Linter MUST NOT impose limits on type name or comment length (rely on Go compiler and other tools)

**Concurrency & Encoding Requirements**:
- **FR-060**: Linter MUST ensure thread-safety of quasi-enum registry for concurrent package analysis
- **FR-061**: Linter MAY rely on `golang.org/x/tools/go/analysis` framework for concurrent analysis behavior (no special requirements)
- **Memory**: Memory usage proportional to analyzed code (no special constraints)
- **Compatibility**: All platforms supported by Go 1.22+

## Glossary

**Quasi-Enum Type**: A named type derived from a Go base type (integer, float, string, complex) that matches at least one detection technique (DT-001 through DT-005). Also referred to as "enum type" in this specification.

**Enum Constant**: A constant of a quasi-enum type that satisfies all enabled definition constraints (DC-001 through DC-005). These are the valid values for the quasi-enum type.

**Base Type**: One of Go's basic types: integer (int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64), float (float32, float64), string, or complex (complex64, complex128).

**Detection Technique (DT)**: A method used to identify whether a type should be treated as a quasi-enum. Five techniques are defined (DT-001 through DT-005).

**Definition Constraint (DC)**: A validation rule that enum constants must satisfy. Five constraints are defined (DC-001 through DC-005).

**Violation**: An instance where:
- A non-enum-constant value is assigned to or used as a quasi-enum type (usage violation), OR
- A quasi-enum type fails to satisfy one or more definition constraints (constraint violation)

**Usage Violation**: Occurs when code uses a literal, untyped constant, or type conversion instead of a defined enum constant for a quasi-enum type.

**Constraint Violation**: Occurs when a detected quasi-enum type fails to satisfy one or more enabled definition constraints.

**Const Block**: A group of constant declarations within a single `const ()` statement in Go.

**Proximity**: The requirement that a type declaration and its const block be close together in the source file, with only empty lines and comments allowed between them.

**Type Conversion**: An explicit conversion from one type to another using the syntax `Type(value)`, e.g., `Status(5)`.

**Untyped Constant**: A Go constant declared without an explicit type, e.g., `const MyValue = 5` (as opposed to `const MyValue Status = 5`).

**Literal**: A fixed value written directly in code, e.g., `5`, `"active"`, `1.5`.

**Composite Literal**: A Go construct for creating values of structs, arrays, slices, and maps, e.g., `[]Status{1, 2, 3}`.

**Autofix**: An automatic code correction that the linter can apply to fix a detected issue (e.g., changing `int` to `uint8`).

**Quality-of-Life Check**: A non-critical suggestion to improve code quality, such as optimization suggestions or missing helper method warnings (US4, US5, US6).

**Thread-Safety**: The property of code that ensures correct behavior when accessed by multiple goroutines concurrently.

**Registry**: An internal data structure that stores information about detected quasi-enum types during analysis.

**Analysis Framework**: The `golang.org/x/tools/go/analysis` package that provides infrastructure for writing Go static analysis tools.

**gopls**: The official Go language server that provides IDE features like code completion, navigation, and diagnostics.

**go vet**: A standard Go tool for running static analysis checks on Go source code.

## Edge Cases *(important)*
- **FR-062**: Linter MAY rely on Go standard library for Unicode/encoding handling (no special requirements beyond ASCII case-insensitive for "enum" keyword)

**Security & Reliability Requirements**:
- **FR-063**: Linter MAY rely on Go compiler and runtime for code security (no special security requirements for analyzing untrusted code)
- **FR-064**: Linter MAY rely on `golang.org/x/tools/go/analysis` framework for reliability and determinism (no special requirements)

**Dependency & Integration Requirements**:
- **FR-065**: Linter MUST ignore missing or incompatible dependencies (continue analysis with available information)
- **FR-066**: Linter MAY rely on Go standard package structure assumptions (no special requirements)
- **FR-067**: Linter MUST accept both `//` and `/* */` comment syntax for detection (examples use `//` syntax for suggestions)
- **FR-068**: Linter integration with IDEs via gopls is not specially supported (standard `go vet` integration only)
- **FR-069**: Linter MAY be used in CI/CD systems via standard Go tooling (no special requirements)
- **FR-070**: Linter MUST support customizing the detection keyword via `-enum-keyword` flag (default: "enum", case-insensitive matching)

**Edge Case Clarifications - Comment Parsing**:
- **FR-071**: Linter has no special requirements for Unicode/special characters in comments beyond ASCII case-insensitive "enum" keyword matching
- **FR-072**: Linter MUST ignore leading whitespace before "enum" keyword in comments when checking first word
- **FR-073**: Opt-out mechanism (`// not enum`) MUST work only with inline comments (not preceding comments)
- **FR-074**: Opt-out mechanism MUST work only with single-line comments (not multi-line comment blocks)
- **FR-075**: Opt-out keyword "not enum" MUST be matched case-insensitively (like "enum" keyword)

**Edge Case Clarifications - Type System**:
- **FR-076**: Linter MUST apply all detection and constraint checks to type aliases (e.g., `type Status = uint8`) same as named types
- **FR-077**: Linter MUST restrict quasi-enum detection to basic underlying types only (int, uint, float, string, complex) because only basic types can be const
- **FR-078**: Linter MUST apply quasi-enum detection to unexported (lowercase) types same as exported types

**Edge Case Clarifications - Constraint Validation**:
- **FR-079**: Linter MUST support iota-based constant counting for DC-001 (minimum constants check)
- **FR-080**: Linter MUST support expression-based constants (e.g., `1 << 0`) for all constraint checks
- **FR-081**: Linter enforces DC-003 (same file) which inherently requires constants in same package (cross-package constants violate this check)
- **FR-082**: Const blocks cannot be nested in Go syntax (no special handling required)
- **FR-083**: Linter MUST accept any const declaration style (grouped `const ()` or individual `const`) for DC-002 validation
- **FR-084**: Linter MUST allow comments and empty lines within const blocks when validating DC-002 and DC-004
- **FR-085**: Linter MUST require all constants in quasi-enum const block to be of the quasi-enum type (for DC-004 exclusive block check)

**Edge Case Clarifications - Proximity**:
- **FR-086**: Linter has no special line count requirements for DC-005 proximity check (structural rule: only empty lines and comments allowed)
- **FR-087**: Linter MUST treat both blank lines and whitespace-only lines as "empty lines" for DC-005 proximity check
- **FR-088**: Linter MUST enforce that only one const block is allowed to define quasi-enum values (DC-002 same block requirement)
- **FR-089**: Linter MUST require type declaration and const block to be in same scope for DC-005 proximity check

**Edge Case Clarifications - Package Boundaries**:
- **FR-090**: When type is in one package and constants imported from another, this violates DC-003 (same file) check
- **FR-091**: Linter has no special requirements for handling build tags (standard Go build tag behavior applies)

**General Requirements**:
- **FR-042**: ~~Linter MUST allow zero value initialization (e.g., `var s Status`) as this is standard Go behavior~~ **SUPERSEDED by FR-094**
- **FR-094**: Linter MUST NOT allow zero value initialization without explicit constant (e.g., `var s Status` is violation unless `const StatusZero Status = 0` exists)
- **FR-095**: When zero (0) is a valid enum constant, it MUST be explicitly listed in the const block (e.g., `const StatusUnknown Status = 0`)

**Edge Case Clarifications - Package Boundaries & Type Conversions**:
- **FR-096**: Quasi-enum definitions (type + constants) MUST NOT cross package boundaries (all must be in same package)
- **FR-097**: Linter MUST allow duplicate constant values for quasi-enum types
- **FR-098**: Linter MUST NOT allow conversions between different quasi-enum types even with same underlying type (e.g., `Status(colorValue)` is violation)
- **FR-099**: Linter MUST NOT allow conversions through interface{} or any to quasi-enum types
- **FR-100**: Linter MUST allow type assertions FROM quasi-enum types but MUST NOT allow type assertions TO quasi-enum types
- **FR-101**: Linter MUST NOT allow unsafe pointer conversions to quasi-enum types

**Edge Case Clarifications - Control Flow & Composite Literals**:
- **FR-102**: Linter MUST apply same violation detection rules to map literals with enum keys or values (e.g., `map[Status]string{1: "active"}` is violation)
- **FR-103**: Linter MUST apply same violation detection rules to struct literals with enum fields (e.g., `Config{Status: 1}` is violation)
- **FR-104**: Linter MUST apply same violation detection rules to switch statement cases with quasi-enum types (e.g., `case 1:` is violation if switching on Status)
- **FR-105**: Linter MUST apply same violation detection rules to if/for conditions with quasi-enum comparisons (e.g., `if status == 1` is violation)
- **FR-106**: Linter MUST apply same violation detection rules to range loops over quasi-enum slices/arrays

**Edge Case Clarifications - Error Handling**:
- **FR-107**: Linter MUST ignore files where AST parsing fails (same as syntax errors)
- **FR-108**: Linter MUST ignore types with incomplete or missing type information
- **FR-109**: Linter has no special requirements for handling file permission errors (covered by FR-054 I/O errors)
- **FR-110**: Linter has no special requirements for handling package import resolution failures (standard Go tooling behavior)
- **FR-111**: Linter has no special requirements for handling out-of-memory conditions (rely on Go runtime)
- **FR-112**: Linter has no special requirements for validating custom enum keyword values beyond "must be valid Go identifier" (FR-070)
- **FR-113**: Linter has no special requirements for handling conflicting flag values (flags are independent per FR-050)

**Edge Case Clarifications - Performance & Boundary Conditions**:
- **FR-120**: Linter has no special requirements for very large const blocks (>1000 constants) beyond Go language limits
- **FR-121**: Linter has no special requirements for analyzing very large files (>10k lines) beyond performance targets
- **FR-122**: When type is detected by non-DT-001 techniques but has fewer than 2 constants, linter MUST report DC-001 violation (minimum 2 constants)
- **FR-123**: Linter has no special requirements for reporting order of multiple violations (implementation-defined)
- **FR-124**: Linter has no special requirements for duplicate violation detection (implementation may report same issue multiple times)

**Edge Case Clarifications - State Management & Concurrency**:
- **FR-125**: Linter code except registry is read-only, so race conditions are not expected outside registry access
- **FR-126**: When quasi-enum registry state becomes inconsistent, linter MUST exit with error
- **FR-127**: Quasi-enum registry is shared across all analyzed packages (not cleaned up between packages)
- **FR-128**: Linter has no special requirements for circular package dependencies (disallowed by Go language)

**Edge Case Clarifications - Integration & CLI**:
- **FR-129**: Linter has no special requirements for unknown/unrecognized flags (handled by analysis framework)
- **FR-130**: Linter has no special requirements for flag value validation beyond type checking (e.g., negative numbers accepted if type allows)
- **FR-131**: Custom enum keyword (via `-enum-keyword` flag) MUST NOT be empty string
- **FR-132**: Linter has no special requirements when used with other vet analyzers that modify AST
- **FR-133**: Linter has no special requirements for handling vet-specific flag prefixing (standard analysis framework behavior)
- **FR-134**: Linter has no special requirements for exit code conflicts with other analyzers

**Edge Case Clarifications - Output & Reporting**:
- **FR-135**: Linter has no special requirements for message formatting when enum has >10 constants
- **FR-136**: Linter has no special requirements for JSON output with special characters in messages
- **FR-137**: Linter has no special requirements for handling very long file paths in diagnostics

**General Requirements**:
- **FR-043**: Linter MUST provide CLI interface for running analysis on Go source files
- **FR-044**: Linter MUST use `golang.org/x/tools/go/analysis` standard flag mechanism for all command-line options
- **FR-045**: Linter MUST support combining multiple disable flags in a single invocation

### Quasi-Enum Detection Techniques

The linter identifies "quasi-enum" types using multiple detection techniques. A type is considered a quasi-enum if it matches ANY of the following criteria:

**DT-001: Constants-Based Detection** (default enabled)
- Type derived from Go base type (integer, float, string, complex)
- Has 2 or more constants of this type defined in the same package

**DT-002: Name Suffix Detection** (default enabled)
- Type derived from Go base type
- Type name has configured keyword suffix (default "enum", case-insensitive, configurable via `-enum-keyword`)
- Examples with default keyword "enum": `StatusEnum`, `PriorityENUM`, `Colorenum`
- Examples with custom keyword "enumeration": `StatusEnumeration`, `ColorENUMERATION`

**DT-003: Inline Comment Detection** (default enabled)
- Type derived from Go base type
- Has comment on the same line with first word matching configured keyword (default "enum", case-insensitive, configurable via `-enum-keyword`)
- Example with default: `type Status uint8 // enum for user status`
- Example with custom keyword: `type Status uint8 // enumeration for user status`

**DT-004: Preceding Comment Detection** (default enabled)
- Type derived from Go base type
- Has comment line(s) immediately before type definition with first word matching configured keyword (default "enum", case-insensitive, configurable via `-enum-keyword`)
- Example with default:
  ```go
  // enum of valid statuses
  type Status uint8
  ```

**DT-005: Named Comment Detection** (default enabled)
- Type derived from Go base type
- Has comment line(s) immediately before type definition with first two words: type name followed by configured keyword (default "enum", case-insensitive, configurable via `-enum-keyword`)
- Example with default:
  ```go
  // Status enum for user states
  type Status uint8
  ```

### Detection Technique Matching Rules

**Note**: The keyword used in DT-002 through DT-005 is configurable via the `-enum-keyword` flag (default: "enum"). All examples below use the default "enum" keyword.

**DT-002 (Name Suffix) Matching**:
- Match: Type name ends with configured keyword (default "enum", case-insensitive)
- Algorithm: `strings.HasSuffix(strings.ToLower(typeName), strings.ToLower(keyword))`
- Examples (with default "enum"):
  - ✅ `StatusEnum` - matches
  - ✅ `PriorityENUM` - matches
  - ✅ `Colorenum` - matches
  - ❌ `EnumStatus` - does NOT match (prefix, not suffix)
  - ❌ `MyEnumeration` - does NOT match ("enum" not at end)

**DT-003 (Inline Comment) Matching**:
- Match: Comment on same line as type declaration, first word is "enum" (case-insensitive)
- Algorithm: Extract comment text, split by whitespace, check if first word equals "enum" (case-insensitive)
- Examples:
  - ✅ `type Status uint8 // enum for user status` - matches
  - ✅ `type Color int // Enum of colors` - matches
  - ✅ `type Priority uint8 // ENUM` - matches
  - ❌ `type Status uint8 // This is an enum` - does NOT match (first word is "This")
  - ❌ `type Status uint8 // my_enum_type` - does NOT match (first word is "my_enum_type")

**DT-004 (Preceding Comment) Matching**:
- Match: Comment immediately before type declaration, first word is "enum" (case-insensitive)
- Algorithm: Check comment line(s) directly above type, extract first word, compare case-insensitively
- Multi-line: Only first line of multi-line comment is checked
- Examples:
  - ✅ `// enum of valid statuses` - matches
  - ✅ `// Enum for user priorities` - matches
  - ✅ `/* enum */` - matches
  - ❌ `// This is an enum` - does NOT match (first word is "This")
  - ❌ `// valid enum statuses` - does NOT match (first word is "valid")

**DT-005 (Named Comment) Matching**:
- Match: Comment immediately before type declaration, first word is type name, second word is "enum" (both case-insensitive)
- Algorithm: Extract first two words, compare first with type name and second with "enum" (case-insensitive)
- Examples:
  - ✅ `// Status enum for users` (type Status) - matches
  - ✅ `// status ENUM` (type Status) - matches (case-insensitive)
  - ✅ `/* Priority enum */` (type Priority) - matches
  - ❌ `// enum Status` (type Status) - does NOT match (wrong order)
  - ❌ `// Status type` (type Status) - does NOT match (second word not "enum")
  - ❌ `// MyStatus enum` (type Status) - does NOT match (first word doesn't match type name)

### Quasi-Enum Definition Constraints

Once a type is detected as a quasi-enum, the following constraints are enforced:

**DC-001: Minimum Constants** (default enabled)
- At least 2 constants MUST be defined for the quasi-enum type
- Violation: Type detected as quasi-enum but has fewer than 2 constants

**DC-002: Same Const Block** (default enabled)
- All constants for the quasi-enum type MUST be in the same const block
- Violation: Constants for same quasi-enum scattered across multiple const blocks
- Example violation:
  ```go
  const StatusActive Status = 1
  // ... other code ...
  const StatusInactive Status = 2  // ERROR: different const block
  ```

**DC-003: Same File** (default enabled)
- The const block MUST be in the same file as the type definition
- Violation: Type defined in one file, constants in another

**DC-004: Exclusive Const Block** (default enabled)
- No other type's constants MUST be in the quasi-enum's const block
- Violation: Mixed constants of different types in same block
- Example violation:
  ```go
  const (
      StatusActive Status = 1
      PriorityHigh Priority = 1  // ERROR: different type in same block
  )
  ```

**DC-005: Proximity** (default enabled)
- Only empty lines and comments MUST be between type definition and const block (any number allowed)
- Violation: Other code (functions, variables, other types) between type and constants
- Example violation:
  ```go
  type Status uint8
  
  var x = 5  // ERROR: non-comment code between type and const block
  
  const (
      StatusActive Status = 1
  )
  ```

### Configuration Options

The linter MUST support command-line flags to disable detection techniques, constraints, and quality-of-life checks individually:

**Detection Technique Flags**:
- `-disable-constants-detection`: Disable DT-001 (constants-based detection)
- `-disable-suffix-detection`: Disable DT-002 (name suffix detection)
- `-disable-inline-comment-detection`: Disable DT-003 (inline comment detection)
- `-disable-preceding-comment-detection`: Disable DT-004 (preceding comment detection)
- `-disable-named-comment-detection`: Disable DT-005 (named comment detection)

**Definition Constraint Flags**:
- `-disable-min-constants-check`: Disable DC-001 (minimum 2 constants)
- `-disable-same-block-check`: Disable DC-002 (same const block)
- `-disable-same-file-check`: Disable DC-003 (same file)
- `-disable-exclusive-block-check`: Disable DC-004 (exclusive const block)
- `-disable-proximity-check`: Disable DC-005 (proximity check)

**Quality-of-Life Check Flags**:
- `-disable-uint8-suggestion`: Disable US4 (uint8 optimization suggestion)
- `-disable-string-method-check`: Disable US5 (String() method check)
- `-disable-unmarshal-method-check`: Disable US6 (UnmarshalText() method check)

**Detection Keyword Configuration**:
- `-enum-keyword=<word>`: Customize the keyword used for detection (default: "enum")
  - This keyword is used in DT-002 (name suffix), DT-003 (inline comment), DT-004 (preceding comment), and DT-005 (named comment)
  - Matching is case-insensitive (e.g., "enum", "Enum", "ENUM" all match)
  - Example: `-enum-keyword=enumeration` would detect types like `StatusEnumeration` or `// enumeration of statuses`
  - Must be a valid Go identifier (alphanumeric and underscore only)

**Flag Implementation**:
- MUST use `golang.org/x/tools/go/analysis` standard flag mechanism
- Flags MUST be registered via `Analyzer.Flags` field
- Flag values MUST be accessible during analysis
- Multiple flags MAY be combined
- Example usage: `go-enumsafety -disable-suffix-detection -disable-proximity-check -disable-uint8-suggestion ./...`

**Flag Independence** (FR-050):
- All flags are completely independent of each other
- No flag has precedence over any other flag
- Flags can be combined in any order with the same result
- Each flag controls only its specific feature (detection technique, constraint, or quality-of-life check)

### Key Entities

- **Quasi-Enum Type**: A named type derived from a Go base type (integer, float, string, complex) that matches at least one detection technique (DT-001 through DT-005)
- **Enum Constant**: A constant of a quasi-enum type, subject to definition constraints (DC-001 through DC-005)
- **Violation**: An instance where a non-enum-constant value is assigned to or used as a quasi-enum type, OR where definition constraints are violated

## Success Criteria *(mandatory)*

### Measurable Outcomes

**Detection Accuracy**:
- **SC-001**: Linter correctly identifies quasi-enums using all 5 detection techniques (DT-001 to DT-005) with 100% accuracy in test cases
- **SC-002**: Linter correctly applies case-insensitive matching for "enum" keyword in all comment-based detection techniques
- **SC-003**: Linter produces zero false positives for types that don't match any detection technique

**Constraint Enforcement**:
- **SC-004**: Linter correctly enforces all 5 definition constraints (DC-001 to DC-005) with 100% accuracy in test cases
- **SC-005**: Linter produces zero false positives for quasi-enums that satisfy all enabled constraints

**Usage Violation Detection**:
- **SC-006**: Linter correctly identifies 100% of literal assignments to quasi-enum types in test cases (zero false negatives)
- **SC-007**: Linter produces zero false positives for valid enum constant usage in test suite

**Configuration**:
- **SC-008**: All 14 command-line flags (5 detection + 5 constraint + 3 quality-of-life + 1 keyword config) work correctly when used individually
- **SC-009**: Multiple flags can be combined in a single invocation without conflicts
- **SC-010**: Disabling a detection technique prevents that technique from identifying quasi-enums
- **SC-011**: Disabling a constraint check prevents that constraint from being enforced

**Performance & Integration**:
- **SC-012**: Linter analysis completes in under 100ms for files under 1000 lines (per constitution performance target)
- **SC-013**: Error messages include the quasi-enum type name and either suggest valid constants OR explain which constraint is violated
- **SC-014**: Linter successfully integrates with `go vet` and can be run as part of standard Go toolchain
- **SC-015**: Linter handles iota-based enum definitions correctly in 100% of test cases
- **SC-016**: Linter can analyze quasi-enum types defined in external packages (detection only)
