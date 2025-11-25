# Contributing to Go Container System

Thank you for your interest in contributing to the Go Container System! This document provides guidelines and instructions for contributing.

---

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [Development Setup](#development-setup)
4. [How to Contribute](#how-to-contribute)
5. [Coding Standards](#coding-standards)
6. [Testing Requirements](#testing-requirements)
7. [Pull Request Process](#pull-request-process)
8. [Release Process](#release-process)

---

## Code of Conduct

This project adheres to a code of conduct that all contributors are expected to follow:

- **Be respectful**: Treat everyone with respect and consideration
- **Be constructive**: Provide helpful feedback and suggestions
- **Be inclusive**: Welcome newcomers and help them get started
- **Be patient**: Remember that everyone has different skill levels

---

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- A GitHub account

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:

```bash
git clone https://github.com/YOUR_USERNAME/go_container_system.git
cd go_container_system
```

3. Add the upstream remote:

```bash
git remote add upstream https://github.com/kcenon/go_container_system.git
```

---

## Development Setup

### Install Dependencies

```bash
# Download dependencies
go mod download

# Verify installation
go mod verify
```

### Build and Test

```bash
# Build
go build ./...

# Run tests
go test ./tests -v

# Run with race detection
go test ./tests -race

# Run benchmarks
go test ./tests -bench=. -benchmem
```

### IDE Setup

**VS Code:**
1. Install the Go extension
2. Open the project folder
3. VS Code will automatically detect Go and configure

**GoLand:**
1. Open as a Go module project
2. Wait for indexing to complete
3. Configure Go SDK if needed

---

## How to Contribute

### Reporting Bugs

Before reporting a bug:
1. Search [existing issues](https://github.com/kcenon/go_container_system/issues)
2. Check the [troubleshooting guide](docs/guides/TROUBLESHOOTING.md)

When reporting:
- Use the bug report template
- Include Go version (`go version`)
- Provide a minimal reproducible example
- Include error messages and stack traces

### Suggesting Features

1. Check [existing discussions](https://github.com/kcenon/go_container_system/discussions)
2. Open a new discussion in the "Ideas" category
3. Describe the feature and its use case
4. Wait for community feedback before implementing

### Contributing Code

1. **Find an issue**: Look for issues labeled `good first issue` or `help wanted`
2. **Discuss**: Comment on the issue to express interest
3. **Implement**: Follow the coding standards
4. **Test**: Add tests for your changes
5. **Submit**: Create a pull request

---

## Coding Standards

### Go Style

Follow the official Go style guidelines:
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Formatting

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run
```

### Naming Conventions

```go
// Package names: lowercase, single word
package values

// Exported types: PascalCase
type ValueContainer struct { ... }

// Exported functions: PascalCase
func NewValueContainer() *ValueContainer { ... }

// Private members: camelCase
type container struct {
    sourceID    string
    targetID    string
    values      []Value
}

// Constants: PascalCase for exported, camelCase for private
const MaxValueCount = 1000
const defaultBufferSize = 4096
```

### Documentation

```go
// Package documentation
// Package values provides concrete implementations of the Value interface.
package values

// Type documentation
// StringValue represents a UTF-8 encoded string value.
// It implements the Value interface and supports all standard
// string operations including serialization to multiple formats.
type StringValue struct {
    name  string
    value string
}

// Function documentation
// NewStringValue creates a new string value with the given name and value.
// The name is used as the key when storing in a container.
// The value must be a valid UTF-8 string.
//
// Example:
//
//	val := NewStringValue("greeting", "Hello, World!")
//	container.AddValue(val)
func NewStringValue(name, value string) *StringValue {
    return &StringValue{name: name, value: value}
}
```

### Error Handling

```go
// Return errors, don't panic
func (v *StringValue) ToInt32() (int32, error) {
    return 0, fmt.Errorf("cannot convert string to int32")
}

// Wrap errors with context
func (c *ValueContainer) Deserialize(data string) error {
    if len(data) == 0 {
        return fmt.Errorf("deserialize: empty data")
    }
    // ...
    if err := parseHeader(data); err != nil {
        return fmt.Errorf("deserialize: %w", err)
    }
    return nil
}

// Check errors explicitly
bytes, err := c.SerializeArray()
if err != nil {
    return fmt.Errorf("failed to serialize: %w", err)
}
```

---

## Testing Requirements

### Test Coverage

All contributions must include tests:
- Minimum 80% coverage for new code
- All public APIs must have tests
- Edge cases and error conditions must be tested

### Test Structure

```go
func TestNewFeature(t *testing.T) {
    // Arrange
    input := createTestInput()

    // Act
    result, err := NewFeature(input)

    // Assert
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result != expected {
        t.Errorf("expected %v, got %v", expected, result)
    }
}
```

### Running Tests

```bash
# All tests
go test ./tests -v

# With coverage
go test ./tests -cover -coverprofile=coverage.out
go tool cover -html=coverage.out

# With race detection
go test ./tests -race

# Specific test
go test ./tests -run TestSpecific -v
```

See [TESTING.md](docs/contributing/TESTING.md) for detailed testing guidelines.

---

## Pull Request Process

### Before Submitting

1. **Update from upstream**:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run all checks**:
   ```bash
   go fmt ./...
   go vet ./...
   go test ./tests -v -race
   golangci-lint run
   ```

3. **Update documentation** if needed

### Creating the PR

1. Push your branch:
   ```bash
   git push origin feature/your-feature
   ```

2. Create a pull request on GitHub

3. Fill out the PR template:
   - Description of changes
   - Related issue(s)
   - Testing performed
   - Breaking changes (if any)

### PR Title Format

```
type(scope): description

Examples:
feat(values): add ArrayValue type
fix(serialization): handle empty containers
docs(readme): update installation instructions
test(container): add header swap tests
refactor(core): simplify value interface
```

### Review Process

1. **Automated checks**: CI must pass
2. **Code review**: At least one maintainer approval
3. **Testing**: Manual testing for complex changes
4. **Documentation**: Verify docs are updated

### After Merge

- Delete your branch
- Update local repository:
  ```bash
  git checkout main
  git pull upstream main
  git branch -d feature/your-feature
  ```

---

## Release Process

### Version Numbering

We follow [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking API changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Checklist

1. Update CHANGELOG.md
2. Update version in relevant files
3. Create release tag
4. Write release notes
5. Publish release on GitHub

---

## Getting Help

- **Documentation**: [docs/README.md](docs/README.md)
- **FAQ**: [docs/guides/FAQ.md](docs/guides/FAQ.md)
- **Discussions**: [GitHub Discussions](https://github.com/kcenon/go_container_system/discussions)
- **Issues**: [GitHub Issues](https://github.com/kcenon/go_container_system/issues)

---

## Recognition

Contributors are recognized in:
- CHANGELOG.md for each release
- GitHub contributors page
- README.md acknowledgments section

Thank you for contributing to the Go Container System!

---

**Last Updated:** 2025-11-26
