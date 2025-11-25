# Data Model: Quasi-Enum Type Safety Linter

**Feature**: 001-enum-linter  
**Date**: 2025-11-23  
**Purpose**: Define the core data structures and entities for quasi-enum type safety analysis with detection techniques and definition constraints

## Overview

This document defines the data model for representing quasi-enum types, their valid constants, detection techniques, definition constraints, and violations. The model supports 5 detection techniques (DT-001 to DT-005), 5 definition constraints (DC-001 to DC-005), and 10 configuration flags.

## Core Entities

### 1. QuasiEnumType

Represents a Go type identified as a quasi-enum through one or more detection techniques.

**Attributes**:
- `TypeName`: `*types.Named` - The named type (e.g., `Status`, `Priority`)
- `UnderlyingType`: `types.BasicKind` - The basic type being aliased (e.g., `uint8`, `int`, `string`)
- `PackagePath`: `string` - Full package path where the quasi-enum is defined
- `Constants`: `[]EnumConstant` - List of valid enum constants
- `Position`: `token.Pos` - Source position of type definition
- `DetectedBy`: `[]DetectionTechnique` - Which techniques detected this type
- `ConstBlock`: `*ast.GenDecl` - The const block containing enum constants (for constraint validation)

**Validation Rules**:
- TypeName MUST be a named type
- UnderlyingType MUST be a basic type (int, uint, string, etc.)
- DetectedBy MUST contain at least one detection technique
- Constants list SHOULD have at least 2 elements (enforced by DC-001 if enabled)

**Example**:
```go
// Source code:
// Status enum for user states
type StatusEnum uint8
const (
    StatusActive StatusEnum = 1
    StatusInactive StatusEnum = 2
)

// QuasiEnumType representation:
QuasiEnumType{
    TypeName: types.Named("StatusEnum"),
    UnderlyingType: types.Uint8,
    PackagePath: "example.com/myapp/models",
    Constants: []EnumConstant{...},
    Position: token.Pos(42),
    DetectedBy: []DetectionTechnique{DT002NameSuffix, DT005NamedComment},
    ConstBlock: &ast.GenDecl{...},
}
```

### 2. DetectionTechnique (Enum)

Represents the technique used to identify a quasi-enum type.

**Values**:
- `DT001ConstantsBased`: Type has 2+ constants in same package
- `DT002NameSuffix`: Type name ends with "enum" (case-insensitive)
- `DT003InlineComment`: Type has inline comment starting with "enum"
- `DT004PrecedingComment`: Type has preceding comment starting with "enum"
- `DT005NamedComment`: Type has comment with "TypeName enum" pattern

**Usage**: Track which techniques detected each quasi-enum for debugging and reporting

**Example**:
```go
type DetectionTechnique int

const (
    DT001ConstantsBased DetectionTechnique = iota
    DT002NameSuffix
    DT003InlineComment
    DT004PrecedingComment
    DT005NamedComment
)

func (dt DetectionTechnique) String() string {
    switch dt {
    case DT001ConstantsBased:
        return "DT-001 (constants-based)"
    case DT002NameSuffix:
        return "DT-002 (name suffix)"
    // ... etc
    }
}
```

### 3. DefinitionConstraint (Enum)

Represents a constraint that quasi-enum definitions must satisfy.

**Values**:
- `DC001MinConstants`: At least 2 constants required
- `DC002SameBlock`: All constants in same const block
- `DC003SameFile`: Const block in same file as type
- `DC004ExclusiveBlock`: No other type's constants in block
- `DC005Proximity`: Only empty lines/comments between type and const block

**Usage**: Track which constraints are violated for error reporting

**Example**:
```go
type DefinitionConstraint int

const (
    DC001MinConstants DefinitionConstraint = iota
    DC002SameBlock
    DC003SameFile
    DC004ExclusiveBlock
    DC005Proximity
)

func (dc DefinitionConstraint) String() string {
    switch dc {
    case DC001MinConstants:
        return "DC-001 (minimum 2 constants)"
    case DC002SameBlock:
        return "DC-002 (same const block)"
    // ... etc
    }
}
```

### 4. EnumConstant

Represents a valid constant value for a quasi-enum type.

**Attributes**:
- `Name`: `string` - Constant identifier (e.g., `StatusActive`)
- `Value`: `constant.Value` - The constant's value (from `go/constant` package)
- `QuasiEnumType`: `*types.Named` - Reference to the parent quasi-enum type
- `Position`: `token.Pos` - Source position of constant definition
- `IsIota`: `bool` - Whether constant uses iota
- `Expression`: `string` - Original expression (e.g., `1 << 0` for bit flags)
- `ConstBlock`: `*ast.GenDecl` - The const block containing this constant

**Validation Rules**:
- Name MUST be a valid Go identifier
- Value MUST be assignable to QuasiEnumType's underlying type
- QuasiEnumType MUST NOT be nil

**Example**:
```go
// Source code:
const StatusActive Status = 1 << 0

// EnumConstant representation:
EnumConstant{
    Name: "StatusActive",
    Value: constant.MakeInt64(1),
    QuasiEnumType: types.Named("Status"),
    Position: token.Pos(100),
    IsIota: false,
    Expression: "1 << 0",
    ConstBlock: &ast.GenDecl{...},
}
```

### 5. Violation

Represents a detected violation of quasi-enum type safety rules OR definition constraints.

**Attributes**:
- `Type`: `ViolationType` - Category of violation
- `Position`: `token.Pos` - Source position of the violation
- `QuasiEnumType`: `*types.Named` - The quasi-enum type involved
- `InvalidValue`: `ast.Expr` - The AST node of the invalid value (for usage violations)
- `Constraint`: `*DefinitionConstraint` - The violated constraint (for constraint violations)
- `Context`: `ViolationContext` - Additional context
- `SuggestedFix`: `string` - Optional suggestion for fixing

**Validation Rules**:
- Type MUST be a valid ViolationType
- Position MUST be valid
- QuasiEnumType MUST NOT be nil
- For usage violations: InvalidValue MUST NOT be nil
- For constraint violations: Constraint MUST NOT be nil

**Example (Usage Violation)**:
```go
// Source code:
var s Status = 5  // VIOLATION

// Violation representation:
Violation{
    Type: ViolationLiteralAssignment,
    Position: token.Pos(200),
    QuasiEnumType: types.Named("Status"),
    InvalidValue: ast.BasicLit{Value: "5"},
    Constraint: nil,
    Context: ViolationContext{VariableName: "s"},
    SuggestedFix: "use Status constants: StatusActive, StatusInactive",
}
```

**Example (Constraint Violation)**:
```go
// Source code:
// Status type with constants in different blocks

// Violation representation:
Violation{
    Type: ViolationConstraint,
    Position: token.Pos(150),
    QuasiEnumType: types.Named("Status"),
    InvalidValue: nil,
    Constraint: &DC002SameBlock,
    Context: ViolationContext{Statement: "const block analysis"},
    SuggestedFix: "move all Status constants into a single const block",
}
```

### 6. ViolationType (Enum)

Categories of violations (usage + constraints).

**Values**:

*Usage Violations*:
- `ViolationLiteralAssignment`: Direct literal assigned to quasi-enum variable
- `ViolationLiteralConversion`: Literal converted to quasi-enum type
- `ViolationLiteralArgument`: Literal passed as quasi-enum function parameter
- `ViolationUntypedConstant`: Untyped constant (not in enum definition) assigned
- `ViolationVariableConversion`: Variable of underlying type converted to quasi-enum
- `ViolationCompositeLiteral`: Literal in composite literal for quasi-enum slice/array

*Constraint Violations*:
- `ViolationConstraint`: Definition constraint violated (see Constraint field for which one)

**Example**:
```go
type ViolationType int

const (
    // Usage violations
    ViolationLiteralAssignment ViolationType = iota
    ViolationLiteralConversion
    ViolationLiteralArgument
    ViolationUntypedConstant
    ViolationVariableConversion
    ViolationCompositeLiteral
    
    // Constraint violation (generic - check Constraint field)
    ViolationConstraint
)
```

### 7. ViolationContext

Additional context information for a violation.

**Attributes**:
- `VariableName`: `string` - Name of variable being assigned (if applicable)
- `FunctionName`: `string` - Name of function being called (if applicable)
- `ParameterName`: `string` - Name of parameter (if applicable)
- `Statement`: `string` - Full statement containing the violation
- `LineNumber`: `int` - Line number in source file
- `ConstraintDetails`: `string` - Additional details for constraint violations

**Purpose**: Provides context for generating helpful error messages.

### 8. QuasiEnumRegistry

Central registry of all quasi-enum types discovered in analyzed packages.

**Attributes**:
- `QuasiEnums`: `map[*types.Named]*QuasiEnumType` - Map from type to quasi-enum definition
- `ConstantLookup`: `map[*types.Named]map[string]*EnumConstant` - Fast constant lookup
- `Packages`: `map[string][]*QuasiEnumType` - Quasi-enums grouped by package
- `DetectionConfig`: `DetectionConfig` - Which detection techniques are enabled
- `ConstraintConfig`: `ConstraintConfig` - Which constraints are enabled

**Operations**:
- `RegisterQuasiEnum(qe *QuasiEnumType)`: Add quasi-enum to registry
- `IsQuasiEnumType(t types.Type) bool`: Check if type is a quasi-enum
- `GetEnumConstants(t *types.Named) []EnumConstant`: Get valid constants for quasi-enum
- `IsValidConstant(t *types.Named, name string) bool`: Check if constant is valid for quasi-enum
- `ValidateConstraints(qe *QuasiEnumType) []Violation`: Validate all enabled constraints

**Lifecycle**: Built once per package analysis, used for all violation detection.

### 9. DetectionConfig

Configuration for detection techniques (from command-line flags).

**Attributes**:
- `ConstantsDetectionEnabled`: `bool` - DT-001 enabled
- `SuffixDetectionEnabled`: `bool` - DT-002 enabled
- `InlineCommentDetectionEnabled`: `bool` - DT-003 enabled
- `PrecedingCommentDetectionEnabled`: `bool` - DT-004 enabled
- `NamedCommentDetectionEnabled`: `bool` - DT-005 enabled

**Default**: All `true`

**Example**:
```go
type DetectionConfig struct {
    ConstantsDetectionEnabled        bool
    SuffixDetectionEnabled          bool
    InlineCommentDetectionEnabled   bool
    PrecedingCommentDetectionEnabled bool
    NamedCommentDetectionEnabled    bool
}

func NewDetectionConfig() *DetectionConfig {
    return &DetectionConfig{
        ConstantsDetectionEnabled:        true,
        SuffixDetectionEnabled:          true,
        InlineCommentDetectionEnabled:   true,
        PrecedingCommentDetectionEnabled: true,
        NamedCommentDetectionEnabled:    true,
    }
}
```

### 10. ConstraintConfig

Configuration for definition constraints (from command-line flags).

**Attributes**:
- `MinConstantsCheckEnabled`: `bool` - DC-001 enabled
- `SameBlockCheckEnabled`: `bool` - DC-002 enabled
- `SameFileCheckEnabled`: `bool` - DC-003 enabled
- `ExclusiveBlockCheckEnabled`: `bool` - DC-004 enabled
- `ProximityCheckEnabled`: `bool` - DC-005 enabled

**Default**: All `true`

**Example**:
```go
type ConstraintConfig struct {
    MinConstantsCheckEnabled    bool
    SameBlockCheckEnabled       bool
    SameFileCheckEnabled        bool
    ExclusiveBlockCheckEnabled  bool
    ProximityCheckEnabled       bool
}

func NewConstraintConfig() *ConstraintConfig {
    return &ConstraintConfig{
        MinConstantsCheckEnabled:    true,
        SameBlockCheckEnabled:       true,
        SameFileCheckEnabled:        true,
        ExclusiveBlockCheckEnabled:  true,
        ProximityCheckEnabled:       true,
    }
}
```

## Relationships

```
QuasiEnumRegistry
    │
    ├─── contains ───> QuasiEnumType (1:N)
    │                      │
    │                      ├─── has ───> EnumConstant (1:N)
    │                      └─── detected by ───> DetectionTechnique (1:N)
    │
    ├─── uses ───> DetectionConfig (1:1)
    ├─── uses ───> ConstraintConfig (1:1)
    │
    └─── referenced by ───> Violation (1:N)
                                │
                                ├─── references ───> QuasiEnumType (N:1)
                                ├─── has type ───> ViolationType (N:1)
                                └─── may reference ───> DefinitionConstraint (N:0..1)
```

## State Transitions

### QuasiEnumType Lifecycle

1. **Discovery**: Type definition found during AST traversal
2. **Detection**: Apply enabled detection techniques (DT-001 to DT-005)
3. **Validation**: Check if type matches any enabled technique
4. **Constant Collection**: Find all constants of this type in same package
5. **Constraint Validation**: Apply enabled constraints (DC-001 to DC-005)
6. **Registration**: Add to QuasiEnumRegistry if valid
7. **Usage**: Referenced during violation detection

### Violation Lifecycle

1. **Detection**: Invalid usage OR constraint violation found
2. **Classification**: Determine ViolationType
3. **Context Collection**: Gather context information
4. **Reporting**: Convert to diagnostic message via `analysis.Pass.Reportf()`

## Data Flow

```
Source Code
    │
    ├─── AST Traversal ───> Discover Type Definitions
    │                           │
    │                           ├───> Apply Detection Techniques (DT-001 to DT-005)
    │                           │         │
    │                           │         └───> Create QuasiEnumType
    │                           │                   │
    │                           │                   ├───> Validate Constraints (DC-001 to DC-005)
    │                           │                   │
    │                           │                   └───> Register in QuasiEnumRegistry
    │
    └─── AST Traversal ───> Discover Assignments/Conversions
                                │
                                ├─── Check against QuasiEnumRegistry
                                │
                                └─── If invalid ───> Create Violation
                                                         │
                                                         └───> Report Diagnostic
```

## Implementation Notes

### Type Checking

Use `go/types` package for accurate type information:
- `types.Named`: For quasi-enum type names
- `types.Basic`: For underlying types
- `constant.Value`: For constant values
- `types.Identical()`: For type comparison

### AST Nodes

Key AST nodes for detection and validation:
- `ast.TypeSpec`: Type definitions
- `ast.GenDecl`: Const blocks and type declarations
- `ast.ValueSpec`: Constant and variable declarations
- `ast.CommentGroup`: Inline and preceding comments
- `ast.AssignStmt`: Assignments
- `ast.CallExpr`: Function calls and type conversions
- `ast.CompositeLit`: Composite literals

### Performance Considerations

- **Parallel detection**: Run detection techniques concurrently
- **Lazy constraint validation**: Only validate constraints for detected quasi-enums
- **Caching**: Cache const block lookups, file lookups
- **Early exit**: If all detection disabled, skip analysis entirely

## Example Usage

```go
// Building the registry with configuration
config := &DetectionConfig{
    ConstantsDetectionEnabled: true,
    SuffixDetectionEnabled: true,
    // ... other techniques
}
constraintConfig := &ConstraintConfig{
    MinConstantsCheckEnabled: true,
    SameBlockCheckEnabled: true,
    // ... other constraints
}

registry := NewQuasiEnumRegistry(config, constraintConfig)

// During AST traversal - type definition found
quasiEnum := &QuasiEnumType{
    TypeName: namedType,
    UnderlyingType: types.Uint8,
    PackagePath: pkg.Path(),
    Constants: collectConstants(namedType, pkg),
    Position: typeSpec.Pos(),
    DetectedBy: []DetectionTechnique{DT001ConstantsBased, DT002NameSuffix},
}

// Validate constraints
violations := registry.ValidateConstraints(quasiEnum)
for _, v := range violations {
    reportConstraintViolation(pass, v)
}

// Register if valid
if len(violations) == 0 || !config.StrictMode {
    registry.RegisterQuasiEnum(quasiEnum)
}

// During AST traversal - assignment found
if registry.IsQuasiEnumType(lhsType) {
    if isLiteral(rhsExpr) {
        violation := &Violation{
            Type: ViolationLiteralAssignment,
            Position: assignStmt.Pos(),
            QuasiEnumType: lhsType.(*types.Named),
            InvalidValue: rhsExpr,
            Context: buildContext(assignStmt),
        }
        reportViolation(pass, violation)
    }
}
```

## Validation Summary

All entities have clear validation rules ensuring:
- **Type safety**: Proper Go types used throughout
- **Completeness**: No nil references where not allowed
- **Consistency**: Relationships between entities maintained
- **Accuracy**: Precise source positions for error reporting
- **Configurability**: Detection and constraints can be individually disabled
