# Research: Quasi-Enum Type Safety Linter

**Feature**: 001-enum-linter  
**Date**: 2025-11-23  
**Purpose**: Research technical approaches for implementing quasi-enum type safety analysis with multiple detection techniques and definition constraints

## Overview

This document consolidates research findings for implementing a Go linter that enforces type safety for quasi-enum patterns using the `golang.org/x/tools/go/analysis` framework. The linter supports 5 detection techniques, 5 definition constraints, and 10 configuration flags.

## Key Research Areas

### 1. golang.org/x/tools/go/analysis Framework

**Decision**: Use `golang.org/x/tools/go/analysis` as the foundation (unchanged from original research)

**Rationale**: Standard framework, rich API, integration, testing support, community

**Key APIs**:
- `analysis.Analyzer`: Main analyzer struct with Run function and **Flags field** for configuration
- `analysis.Pass`: Provides access to AST, type info, and reporting
- `inspect.Analyzer`: Efficient AST traversal using visitor pattern
- `analysistest.Run()`: Testing framework for analyzers

**New Requirement - Flag Registration**:
```go
var Analyzer = &analysis.Analyzer{
    Name: "enumsafety",
    Flags: flagSet(),  // Register 10 configuration flags
    // ...
}

func flagSet() flag.FlagSet {
    fs := flag.NewFlagSet("enumsafety", flag.ExitOnError)
    fs.BoolVar(&disableConstantsDetection, "disable-constants-detection", false, "...")
    // ... 9 more flags
    return *fs
}
```

**References**:
- https://pkg.go.dev/golang.org/x/tools/go/analysis
- https://pkg.go.dev/golang.org/x/tools/go/analysis/analysistest

### 2. Quasi-Enum Detection Strategies

**Decision**: Implement 5 detection techniques, all enabled by default, individually disableable

#### DT-001: Constants-Based Detection

**Algorithm**:
1. Find all type definitions: `type X BasicType`
2. Count constants of each type in same package
3. If count >= 2, mark as quasi-enum candidate

**Implementation**:
```go
func detectByConstants(pass *analysis.Pass) map[*types.Named]bool {
    candidates := make(map[*types.Named]bool)
    
    // Count constants per type
    constantCounts := make(map[*types.Named]int)
    for _, obj := range pass.Pkg.Scope().Names() {
        if c, ok := pass.Pkg.Scope().Lookup(obj).(*types.Const); ok {
            if named, ok := c.Type().(*types.Named); ok {
                if isBasicType(named.Underlying()) {
                    constantCounts[named]++
                }
            }
        }
    }
    
    // Mark types with 2+ constants
    for typ, count := range constantCounts {
        if count >= 2 {
            candidates[typ] = true
        }
    }
    
    return candidates
}
```

#### DT-002: Name Suffix Detection

**Algorithm**:
1. Find all type definitions
2. Check if type name ends with "enum" (case-insensitive)
3. Mark as quasi-enum candidate

**Implementation**:
```go
func detectBySuffix(pass *analysis.Pass) map[*types.Named]bool {
    candidates := make(map[*types.Named]bool)
    
    for _, file := range pass.Files {
        for _, decl := range file.Decls {
            if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
                for _, spec := range genDecl.Specs {
                    if typeSpec, ok := spec.(*ast.TypeSpec); ok {
                        if strings.HasSuffix(strings.ToLower(typeSpec.Name.Name), "enum") {
                            if named := pass.TypesInfo.Defs[typeSpec.Name].Type().(*types.Named); named != nil {
                                if isBasicType(named.Underlying()) {
                                    candidates[named] = true
                                }
                            }
                        }
                    }
                }
            }
        }
    }
    
    return candidates
}
```

#### DT-003: Inline Comment Detection

**Algorithm**:
1. Find all type definitions
2. Check for comment on same line as type definition
3. If comment starts with "enum" (case-insensitive), mark as quasi-enum

**Implementation**:
```go
func detectByInlineComment(pass *analysis.Pass) map[*types.Named]bool {
    candidates := make(map[*types.Named]bool)
    
    for _, file := range pass.Files {
        for _, decl := range file.Decls {
            if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
                for _, spec := range genDecl.Specs {
                    if typeSpec, ok := spec.(*ast.TypeSpec); ok {
                        // Check line comment
                        if typeSpec.Comment != nil {
                            for _, comment := range typeSpec.Comment.List {
                                text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
                                if startsWithEnum(text) {
                                    if named := pass.TypesInfo.Defs[typeSpec.Name].Type().(*types.Named); named != nil {
                                        if isBasicType(named.Underlying()) {
                                            candidates[named] = true
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }
    
    return candidates
}
```

#### DT-004: Preceding Comment Detection

**Algorithm**:
1. Find all type definitions
2. Check doc comment (immediately preceding lines)
3. If any line starts with "enum" (case-insensitive), mark as quasi-enum

**Implementation**:
```go
func detectByPrecedingComment(pass *analysis.Pass) map[*types.Named]bool {
    candidates := make(map[*types.Named]bool)
    
    for _, file := range pass.Files {
        for _, decl := range file.Decls {
            if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
                // Check doc comment
                if genDecl.Doc != nil {
                    for _, comment := range genDecl.Doc.List {
                        text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
                        if startsWithEnum(text) {
                            for _, spec := range genDecl.Specs {
                                if typeSpec, ok := spec.(*ast.TypeSpec); ok {
                                    if named := pass.TypesInfo.Defs[typeSpec.Name].Type().(*types.Named); named != nil {
                                        if isBasicType(named.Underlying()) {
                                            candidates[named] = true
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }
    
    return candidates
}
```

#### DT-005: Named Comment Detection

**Algorithm**:
1. Find all type definitions
2. Check doc comment for pattern "TypeName enum" (case-insensitive)
3. Mark as quasi-enum

**Implementation**:
```go
func detectByNamedComment(pass *analysis.Pass) map[*types.Named]bool {
    candidates := make(map[*types.Named]bool)
    
    for _, file := range pass.Files {
        for _, decl := range file.Decls {
            if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
                if genDecl.Doc != nil {
                    for _, spec := range genDecl.Specs {
                        if typeSpec, ok := spec.(*ast.TypeSpec); ok {
                            typeName := typeSpec.Name.Name
                            for _, comment := range genDecl.Doc.List {
                                text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
                                if startsWithTypeNameEnum(text, typeName) {
                                    if named := pass.TypesInfo.Defs[typeSpec.Name].Type().(*types.Named); named != nil {
                                        if isBasicType(named.Underlying()) {
                                            candidates[named] = true
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }
    
    return candidates
}
```

**Helper Functions**:
```go
func startsWithEnum(text string) bool {
    lower := strings.ToLower(text)
    return strings.HasPrefix(lower, "enum ")  || lower == "enum"
}

func startsWithTypeNameEnum(text string, typeName string) bool {
    lower := strings.ToLower(text)
    pattern := strings.ToLower(typeName) + " enum"
    return strings.HasPrefix(lower, pattern)
}

func isBasicType(t types.Type) bool {
    _, ok := t.(*types.Basic)
    return ok
}
```

### 3. Definition Constraint Validation

**Decision**: Implement 5 constraints, all enabled by default, individually disableable

#### DC-001: Minimum Constants Check

**Algorithm**:
```go
func validateMinConstants(quasiEnum *QuasiEnumType) *ConstraintViolation {
    if len(quasiEnum.Constants) < 2 {
        return &ConstraintViolation{
            Constraint: DC001MinConstants,
            Message: fmt.Sprintf("quasi-enum %s has only %d constant(s), minimum 2 required",
                quasiEnum.TypeName.Obj().Name(), len(quasiEnum.Constants)),
        }
    }
    return nil
}
```

#### DC-002: Same Const Block Check

**Algorithm**:
```go
func validateSameConstBlock(quasiEnum *QuasiEnumType, pass *analysis.Pass) *ConstraintViolation {
    if len(quasiEnum.Constants) == 0 {
        return nil
    }
    
    // Find const blocks for each constant
    blocks := make(map[*ast.GenDecl]int)
    for _, c := range quasiEnum.Constants {
        block := findConstBlock(c.Position, pass)
        if block != nil {
            blocks[block]++
        }
    }
    
    if len(blocks) > 1 {
        return &ConstraintViolation{
            Constraint: DC002SameBlock,
            Message: fmt.Sprintf("quasi-enum %s has constants in %d different const blocks, must be in same block",
                quasiEnum.TypeName.Obj().Name(), len(blocks)),
        }
    }
    return nil
}
```

#### DC-003: Same File Check

**Algorithm**:
```go
func validateSameFile(quasiEnum *QuasiEnumType, pass *analysis.Pass) *ConstraintViolation {
    typeFile := pass.Fset.File(quasiEnum.Position)
    
    for _, c := range quasiEnum.Constants {
        constFile := pass.Fset.File(c.Position)
        if typeFile != constFile {
            return &ConstraintViolation{
                Constraint: DC003SameFile,
                Message: fmt.Sprintf("quasi-enum %s type and constants must be in same file",
                    quasiEnum.TypeName.Obj().Name()),
            }
        }
    }
    return nil
}
```

#### DC-004: Exclusive Const Block Check

**Algorithm**:
```go
func validateExclusiveBlock(quasiEnum *QuasiEnumType, pass *analysis.Pass) *ConstraintViolation {
    // Find the const block containing enum constants
    block := findConstBlock(quasiEnum.Constants[0].Position, pass)
    if block == nil {
        return nil
    }
    
    // Check all constants in block
    for _, spec := range block.Specs {
        if valueSpec, ok := spec.(*ast.ValueSpec); ok {
            for _, name := range valueSpec.Names {
                obj := pass.TypesInfo.Defs[name]
                if c, ok := obj.(*types.Const); ok {
                    if named, ok := c.Type().(*types.Named); ok {
                        if named != quasiEnum.TypeName {
                            return &ConstraintViolation{
                                Constraint: DC004ExclusiveBlock,
                                Message: fmt.Sprintf("quasi-enum %s const block contains constants of other types",
                                    quasiEnum.TypeName.Obj().Name()),
                            }
                        }
                    }
                }
            }
        }
    }
    return nil
}
```

#### DC-005: Proximity Check

**Algorithm**:
```go
func validateProximity(quasiEnum *QuasiEnumType, pass *analysis.Pass) *ConstraintViolation {
    typePos := quasiEnum.Position
    constBlock := findConstBlock(quasiEnum.Constants[0].Position, pass)
    if constBlock == nil {
        return nil
    }
    
    // Find all nodes between type and const block
    file := findFile(typePos, pass)
    between := findNodesBetween(file, typePos, constBlock.Pos(), pass)
    
    // Check if only empty lines and comments
    for _, node := range between {
        if !isCommentOrEmpty(node) {
            return &ConstraintViolation{
                Constraint: DC005Proximity,
                Message: fmt.Sprintf("quasi-enum %s has non-comment code between type and const block",
                    quasiEnum.TypeName.Obj().Name()),
            }
        }
    }
    return nil
}
```

### 4. Configuration Flag System

**Decision**: Use `analysis.Analyzer.Flags` with standard `flag` package

**Implementation**:
```go
var (
    // Detection technique flags
    disableConstantsDetection        bool
    disableSuffixDetection          bool
    disableInlineCommentDetection   bool
    disablePrecedingCommentDetection bool
    disableNamedCommentDetection    bool
    
    // Definition constraint flags
    disableMinConstantsCheck    bool
    disableSameBlockCheck       bool
    disableSameFileCheck        bool
    disableExclusiveBlockCheck  bool
    disableProximityCheck       bool
)

func init() {
    Analyzer.Flags.BoolVar(&disableConstantsDetection, "disable-constants-detection", false,
        "disable DT-001: constants-based detection")
    Analyzer.Flags.BoolVar(&disableSuffixDetection, "disable-suffix-detection", false,
        "disable DT-002: name suffix detection")
    // ... 8 more flags
}
```

**Usage in Detection**:
```go
func detectQuasiEnums(pass *analysis.Pass) map[*types.Named]*QuasiEnumType {
    candidates := make(map[*types.Named]bool)
    
    // Apply enabled detection techniques
    if !disableConstantsDetection {
        merge(candidates, detectByConstants(pass))
    }
    if !disableSuffixDetection {
        merge(candidates, detectBySuffix(pass))
    }
    // ... other techniques
    
    return buildQuasiEnums(candidates, pass)
}
```

### 5. Error Message Design

**Decision**: Different message formats for usage violations vs constraint violations

**Usage Violation Format** (unchanged):
```
invalid assignment to quasi-enum type 'Status': literal value not allowed
  valid constants: StatusActive, StatusInactive, StatusPending
  suggestion: use Status constants instead of literal values
```

**Constraint Violation Format** (new):
```
quasi-enum definition constraint violated: Status
  constraint: DC-002 (same const block)
  problem: constants found in 2 different const blocks
  fix: move all Status constants into a single const block
```

### 6. Testing Strategy

**Decision**: Extend analysistest with constraint violation testing

**Test Structure** (updated):
```
internal/testdata/src/a/
├── enum.go              # Valid quasi-enum patterns (all detection techniques)
├── violations.go        # Usage violations with want comments
├── constraints.go       # Constraint violations with want comments
├── iota.go              # Iota-based enums
├── expressions.go       # Expression-based constants
├── flags/               # Flag combination tests
│   ├── disabled_dt.go   # Tests with detection disabled
│   └── disabled_dc.go   # Tests with constraints disabled
└── crosspackage/        # Cross-package enum usage
    └── types.go
```

**New Test Categories**:
- Detection technique tests (each technique individually)
- Constraint violation tests (each constraint individually)
- Flag combination tests
- Edge case tests (proximity, exclusive block, etc.)

### 7. Performance Optimization

**Decision**: Parallel detection, lazy constraint validation

**Optimization Strategy**:
1. **Parallel Detection**: Run all 5 detection techniques concurrently
2. **Early Exit**: If all detection disabled, skip analysis
3. **Lazy Constraints**: Only validate constraints for detected quasi-enums
4. **Caching**: Cache AST node lookups (const blocks, files)

**Implementation**:
```go
func detectQuasiEnumsParallel(pass *analysis.Pass) map[*types.Named]*QuasiEnumType {
    var wg sync.WaitGroup
    results := make([]map[*types.Named]bool, 5)
    
    techniques := []func(*analysis.Pass) map[*types.Named]bool{
        detectByConstants,
        detectBySuffix,
        detectByInlineComment,
        detectByPrecedingComment,
        detectByNamedComment,
    }
    
    for i, technique := range techniques {
        if !isDetectionDisabled(i) {
            wg.Add(1)
            go func(idx int, fn func(*analysis.Pass) map[*types.Named]bool) {
                defer wg.Done()
                results[idx] = fn(pass)
            }(i, technique)
        }
    }
    
    wg.Wait()
    
    // Merge results
    candidates := make(map[*types.Named]bool)
    for _, result := range results {
        merge(candidates, result)
    }
    
    return buildQuasiEnums(candidates, pass)
}
```

## Technology Stack Summary (Updated)

| Component | Technology | Rationale |
|-----------|------------|-----------|
| Language | Go 1.22+ | Constitution requirement, latest features |
| Analysis Framework | golang.org/x/tools/go/analysis | Standard, well-supported, constitution requirement |
| Configuration | flag package via Analyzer.Flags | Standard Go flag mechanism |
| Comment Parsing | ast.CommentGroup | Built-in AST comment support |
| CLI | singlechecker.Main() | Standard Go tool UX, minimal code |
| Testing | analysistest + testing package | Standard analyzer testing, constitution requirement |
| AST Traversal | inspect.Analyzer | Efficient visitor pattern |
| Type Checking | go/types via analysis.Pass | Standard type information |

## Implementation Phases (Updated)

Based on this research, implementation will proceed as:

1. **Phase 1**: Basic analyzer skeleton + flag registration
2. **Phase 2**: Implement 5 detection techniques (DT-001 to DT-005)
3. **Phase 3**: Implement 5 definition constraints (DC-001 to DC-005)
4. **Phase 4**: Literal assignment detection (P1 user story)
5. **Phase 5**: Untyped constant detection (P2 user story)
6. **Phase 6**: Variable conversion detection (P3 user story)
7. **Phase 7**: Edge cases (iota, expressions, composite literals)
8. **Phase 8**: Error message enhancement and CLI polish

## Open Questions

**Resolved** - All technical decisions made through research and spec requirements.

## References

- [Writing Go analyzers](https://disaev.me/p/writing-useful-go-analysis-linter/)
- [Go analysis framework guide](https://pkg.go.dev/golang.org/x/tools/go/analysis)
- [Example analyzers](https://github.com/golang/tools/tree/master/go/analysis/passes)
- [analysistest documentation](https://pkg.go.dev/golang.org/x/tools/go/analysis/analysistest)
- [Flag package documentation](https://pkg.go.dev/flag)
- [AST comment handling](https://pkg.go.dev/go/ast#CommentGroup)
