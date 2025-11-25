# Go Container System - Testing Guide

> **Last Updated:** 2025-11-26
> **Maintainer:** Container System Development Team

## Overview

This document provides a comprehensive guide to the Go Container System's testing infrastructure, including unit tests, integration tests, performance benchmarks, and testing best practices.

---

## Test Architecture

The Go Container System employs a multi-layered testing strategy:

```
┌─────────────────────────────────────────────┐
│           Testing Infrastructure            │
├─────────────────────────────────────────────┤
│  Unit Tests (tests/)                        │
│  ├─ Value type tests                        │
│  ├─ Container operations                    │
│  ├─ Serialization formats                   │
│  └─ Wire Protocol compatibility             │
├─────────────────────────────────────────────┤
│  Integration Tests (tests/)                 │
│  ├─ Cross-language compatibility            │
│  ├─ Binary interoperability                 │
│  └─ End-to-end scenarios                    │
├─────────────────────────────────────────────┤
│  Benchmark Tests (tests/)                   │
│  ├─ Container operations                    │
│  ├─ Serialization performance               │
│  └─ Memory allocation                       │
└─────────────────────────────────────────────┘
```

---

## Test Organization

### Test Files

| Test File | Focus Area | Description |
|-----------|------------|-------------|
| `container_test.go` | Core container | Container creation, values, headers |
| `container_bench_test.go` | Performance | Benchmarks for all operations |
| `interop_test.go` | Interoperability | Cross-language compatibility |
| `binary_interop_test.go` | Wire Protocol | Binary format validation |
| `cross_language_test.go` | Full interop | Complete cross-language tests |

### Test Categories

#### 1. Value Type Tests
Tests for all 15 supported value types:
- Null, Boolean, Numeric (int16-int64, uint16-uint64)
- Float32, Float64
- String and Bytes values
- Container values (nested structures)
- Type conversions between types

#### 2. Container Operations Tests
Core container functionality:
- Container creation (empty, with type, full header)
- Value addition and retrieval
- Header management (source/target/message type)
- Serialization/deserialization
- Copy operations (shallow/deep)
- Header swapping

#### 3. Serialization Tests
Format-specific tests:
- String serialization (legacy format)
- Binary serialization (Wire Protocol)
- JSON serialization/deserialization
- XML serialization/deserialization
- Round-trip validation

#### 4. Cross-Language Tests
Interoperability validation:
- C++ Wire Protocol compatibility
- Binary format byte-level validation
- Type mapping verification
- Edge cases (empty containers, large data)

---

## Running Tests

### Run All Tests

```bash
# Run all tests
go test ./tests -v

# Run with coverage
go test ./tests -v -cover

# Run with race detection
go test ./tests -v -race

# Run with timeout
go test ./tests -v -timeout 30s
```

### Run Specific Tests

```bash
# Run specific test file
go test ./tests -v -run TestBoolValue

# Run tests matching pattern
go test ./tests -v -run "Test.*Value"

# Run container tests only
go test ./tests -v -run TestValueContainer

# Run serialization tests
go test ./tests -v -run "Test.*Serialize"
```

### Run Benchmarks

```bash
# Run all benchmarks
go test ./tests -bench=. -benchmem

# Run specific benchmark
go test ./tests -bench=BenchmarkSerialize -benchmem

# Run with custom iterations
go test ./tests -bench=. -benchtime=5s

# Compare benchmarks (using benchstat)
go test ./tests -bench=. -count=10 > old.txt
# Make changes...
go test ./tests -bench=. -count=10 > new.txt
benchstat old.txt new.txt
```

---

## Test Coverage

### Current Coverage Targets

| Component | Target | Status |
|-----------|--------|--------|
| Core Container | 85%+ | ✅ |
| Value Types | 90%+ | ✅ |
| Serialization | 85%+ | ✅ |
| Wire Protocol | 80%+ | ✅ |
| **Overall** | **85%+** | ✅ |

### Generating Coverage Report

```bash
# Generate coverage profile
go test ./tests -coverprofile=coverage.out

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# View coverage by package
go test ./... -cover
```

---

## Writing Tests

### Test Naming Conventions

```go
// Format: Test<Component>_<Operation>_<Scenario>
func TestBoolValue_Creation(t *testing.T) { ... }
func TestContainer_AddValue_MultipleValues(t *testing.T) { ... }
func TestSerialization_JSON_NestedContainer(t *testing.T) { ... }
```

### Test Structure (AAA Pattern)

```go
func TestContainer_AddValue(t *testing.T) {
    // Arrange: Set up test data
    c := core.NewValueContainer()
    val := values.NewStringValue("name", "Alice")

    // Act: Perform the operation
    c.AddValue(val)

    // Assert: Verify results
    retrieved := c.GetValue("name", 0)
    if retrieved == nil {
        t.Fatal("Expected value to be found")
    }

    str, err := retrieved.ToString()
    if err != nil {
        t.Fatalf("ToString failed: %v", err)
    }
    if str != "Alice" {
        t.Errorf("Expected 'Alice', got '%s'", str)
    }
}
```

### Table-Driven Tests

```go
func TestNumericConversions(t *testing.T) {
    tests := []struct {
        name     string
        value    core.Value
        expected interface{}
    }{
        {"Int32 to Int32", values.NewInt32Value("n", 42), int32(42)},
        {"Int64 to Int64", values.NewInt64Value("n", 12345), int64(12345)},
        {"Float32 to Float32", values.NewFloat32Value("n", 3.14), float32(3.14)},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test conversion based on type
            switch expected := tt.expected.(type) {
            case int32:
                result, err := tt.value.ToInt32()
                if err != nil {
                    t.Fatalf("Conversion error: %v", err)
                }
                if result != expected {
                    t.Errorf("Expected %d, got %d", expected, result)
                }
            // ... other cases
            }
        })
    }
}
```

### Subtests for Related Tests

```go
func TestSerialization(t *testing.T) {
    c := createTestContainer()

    t.Run("String", func(t *testing.T) {
        serialized, err := c.Serialize()
        if err != nil {
            t.Fatalf("Serialize failed: %v", err)
        }
        // Verify...
    })

    t.Run("Binary", func(t *testing.T) {
        bytes, err := c.SerializeArray()
        if err != nil {
            t.Fatalf("SerializeArray failed: %v", err)
        }
        // Verify...
    })

    t.Run("JSON", func(t *testing.T) {
        json, err := c.ToJSON()
        if err != nil {
            t.Fatalf("ToJSON failed: %v", err)
        }
        // Verify...
    })
}
```

---

## Benchmark Guidelines

### Writing Benchmarks

```go
func BenchmarkContainerCreation(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = core.NewValueContainer()
    }
}

func BenchmarkSerialize(b *testing.B) {
    c := createTestContainer()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = c.Serialize()
    }
}

func BenchmarkSerializeArray(b *testing.B) {
    c := createTestContainer()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = c.SerializeArray()
    }
}

// Benchmark with memory allocation reporting
func BenchmarkAddValues(b *testing.B) {
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        c := core.NewValueContainer()
        for j := 0; j < 10; j++ {
            c.AddValue(values.NewInt32Value(fmt.Sprintf("val_%d", j), int32(j)))
        }
    }
}
```

### Expected Performance Baselines

| Operation | Target | Notes |
|-----------|--------|-------|
| Container creation | >500K/s | Empty container |
| Value addition | >2M/s | Single value |
| Binary serialize | >300K/s | Medium container |
| Binary deserialize | >250K/s | Medium container |
| JSON serialize | >100K/s | Medium container |
| JSON deserialize | >80K/s | Medium container |

---

## Test Fixtures and Helpers

### Common Test Helpers

```go
// tests/helpers_test.go

func createTestContainer() *core.ValueContainer {
    c := core.NewValueContainerFull(
        "test_source", "sub1",
        "test_target", "sub2",
        "test_message",
    )
    c.AddValue(values.NewStringValue("name", "Test"))
    c.AddValue(values.NewInt32Value("count", 42))
    c.AddValue(values.NewBoolValue("active", true))
    return c
}

func createNestedContainer() *core.ValueContainer {
    c := core.NewValueContainerWithType("nested_test")

    inner := values.NewContainerValue("inner",
        values.NewStringValue("key", "value"),
        values.NewInt32Value("num", 123),
    )
    c.AddValue(inner)

    return c
}

func assertContainersEqual(t *testing.T, expected, actual *core.ValueContainer) {
    t.Helper()

    // Compare headers
    if expected.SourceID() != actual.SourceID() {
        t.Errorf("SourceID mismatch: expected %s, got %s",
            expected.SourceID(), actual.SourceID())
    }
    // ... compare other fields
}
```

---

## Continuous Integration

### CI/CD Pipeline

Tests are automatically executed in the CI/CD pipeline:

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go: ['1.21', '1.22']

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}

      - name: Run Tests
        run: go test ./tests -v -race -coverprofile=coverage.out

      - name: Upload Coverage
        uses: codecov/codecov-action@v4
        with:
          file: coverage.out
```

---

## Debugging Failed Tests

### Common Issues

1. **Race Conditions**
   ```bash
   # Run with race detector
   go test ./tests -race -v
   ```

2. **Flaky Tests**
   ```bash
   # Run test multiple times
   go test ./tests -v -run TestFlaky -count=100
   ```

3. **Platform-Specific Failures**
   - Check byte order (endianness)
   - Verify path separators
   - Check timeout values

### Debugging Commands

```bash
# Verbose output
go test ./tests -v -run TestFailing

# Print test output
go test ./tests -v -run TestFailing 2>&1 | tee test.log

# Run with dlv debugger
dlv test ./tests -- -test.run TestFailing
```

---

## Test Quality Metrics

### Goals

- **Test Execution Time**: <30 seconds for all unit tests
- **Test Reliability**: >99% pass rate
- **Code Coverage**: Maintain >85% line coverage
- **Benchmark Stability**: <5% variance between runs

### Monitoring

```bash
# Check test execution time
go test ./tests -v 2>&1 | grep -E "^(ok|FAIL)"

# Check coverage
go test ./tests -cover | grep coverage

# Run benchmarks with statistics
go test ./tests -bench=. -benchmem -count=5 | tee bench.txt
```

---

## References

### Related Documentation

- [ARCHITECTURE.md](../../ARCHITECTURE.md) - System architecture
- [FAQ.md](../guides/FAQ.md) - Frequently asked questions
- [BENCHMARKS.md](../performance/BENCHMARKS.md) - Performance benchmarks

### External Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
- [Table-Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)

---

**Document Version:** 1.0
**Last Updated:** 2025-11-26
**Contact:** kcenon@naver.com
