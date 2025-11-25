# Contributing to enumsafety

Thank you for your interest in contributing to enumsafety! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Code Style](#code-style)
- [Release Process](#release-process)

## Getting Started

### Prerequisites

- Go 1.22 or later
- Git
- Familiarity with Go's `go/analysis` framework

### Development Setup

1. **Fork and Clone**

```bash
git clone https://github.com/YOUR_USERNAME/enumsafety.git
cd enumsafety
```

2. **Install Dependencies**

```bash
go mod download
```

3. **Build the Project**

```bash
go build ./cmd/enumsafety
```

4. **Run Tests**

```bash
go test -v ./tests/unit/...
```

## Project Structure

```
enumsafety/
├── cmd/
│   └── enumsafety/
│       └── main.go              # CLI entry point
├── internal/
│   └── analyzer/
│       ├── analyzer.go          # Main analyzer with flag registration
│       ├── detection.go         # Detection techniques (DT-001 to DT-005)
│       ├── constraints.go       # Definition constraints (DC-001 to DC-005)
│       ├── usage.go             # Violation detection (US1-US3)
│       ├── optimization.go      # uint8 optimization (US4)
│       ├── helpers.go           # Helper method checks (US5-US6)
│       ├── enum.go              # QuasiEnumType data structure
│       └── violation.go         # Violation reporting
├── tests/
│   └── unit/
│       └── analyzer_test.go     # Unit tests
├── internal/testdata/
│   └── src/                     # Test fixtures
│       ├── a/                   # US1-US3 test cases
│       ├── optimization/        # US4 test cases
│       ├── helpers/             # US5-US6 test cases
│       ├── detection/           # Detection technique tests
│       └── constraints/         # Constraint validation tests
└── specs/
    └── 001-enum-linter/         # Feature specification
```

## Making Changes

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Your Changes

Follow these guidelines:

- **Keep changes focused**: One feature or bug fix per PR
- **Write tests**: All new features must have tests
- **Update documentation**: Update README.md if adding features
- **Follow code style**: Run `gofmt` and `golint`

### 3. Add Tests

All changes must include tests. Test files should be placed in:

- `tests/unit/analyzer_test.go` - For integration tests
- `internal/testdata/src/` - For test fixtures

Example test:

```go
func TestNewFeature(t *testing.T) {
    wd, err := os.Getwd()
    if err != nil {
        t.Fatal(err)
    }
    testdata := filepath.Join(wd, "..", "..", "internal", "testdata")
    analysistest.Run(t, testdata, analyzer.Analyzer, "yourpackage")
}
```

### 4. Update Test Fixtures

Test fixtures use the `analysistest` framework. Add `// want` comments for expected diagnostics:

```go
type Status int // want "quasi-enum type Status uses int but has only 3 constants"

var s Status = 5 // want "literal value assigned to quasi-enum type Status"
```

## Testing

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run specific test
go test -v ./tests/unit -run TestUS1

# Run with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...
```

### Test Requirements

- All tests must pass
- New features must have >80% test coverage
- Test fixtures must be comprehensive
- No false positives/negatives allowed

### Writing Good Tests

1. **Test one thing at a time**
2. **Use descriptive names**: `TestLiteralAssignmentDetection`
3. **Include edge cases**: Zero values, nil, empty strings, etc.
4. **Use table-driven tests** when appropriate
5. **Add comments** explaining complex test scenarios

## Submitting Changes

### 1. Commit Your Changes

Use clear, descriptive commit messages:

```bash
git commit -m "feat: add support for custom enum keywords

- Implement -enum-keyword flag
- Update detection logic to use configurable keyword
- Add tests for custom keyword detection
- Update documentation

Fixes #123"
```

Commit message format:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Test additions/changes
- `refactor:` - Code refactoring
- `perf:` - Performance improvements

### 2. Push to Your Fork

```bash
git push origin feature/your-feature-name
```

### 3. Create Pull Request

- Provide clear description of changes
- Reference related issues
- Include test results
- Add screenshots if UI-related

### PR Checklist

- [ ] Tests pass locally
- [ ] Code follows project style
- [ ] Documentation updated
- [ ] Commit messages are clear
- [ ] No merge conflicts
- [ ] Added tests for new features
- [ ] Updated CHANGELOG.md (if applicable)

## Code Style

### Go Style Guidelines

Follow standard Go conventions:

- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use meaningful variable names
- Add comments for exported functions
- Keep functions small and focused

### Analyzer-Specific Guidelines

1. **Flag Naming**: Use `-disable-*` pattern for boolean flags
2. **Error Messages**: Be clear and actionable
3. **Performance**: Minimize AST traversals
4. **Type Safety**: Use type assertions carefully

### Example Code Style

```go
// detectQuasiEnums identifies all quasi-enum types in the package.
// It applies all enabled detection techniques and returns a map of
// detected types to their detection techniques.
func detectQuasiEnums(pass *analysis.Pass, config *DetectionConfig) map[*types.Named][]DetectionTechnique {
    detected := make(map[*types.Named][]DetectionTechnique)
    
    // Apply each enabled detection technique
    if !config.DisableConstantsDetection {
        for namedType := range detectByConstants(pass) {
            detected[namedType] = append(detected[namedType], DT001Constants)
        }
    }
    
    return detected
}
```

## Release Process

### Version Numbering

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Checklist

1. Update version in relevant files
2. Update CHANGELOG.md
3. Run full test suite
4. Create git tag: `git tag v1.0.0`
5. Push tag: `git push origin v1.0.0`
6. Create GitHub release with notes

## Getting Help

- **Issues**: Open an issue for bugs or feature requests
- **Discussions**: Use GitHub Discussions for questions
- **Documentation**: Check specs/001-enum-linter/ for detailed specs

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on the code, not the person
- Help others learn and grow

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Thank You!

Your contributions make enumsafety better for everyone. Thank you for taking the time to contribute!
