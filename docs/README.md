# Go Container System Documentation

> **Language:** **English** | [한국어](../README_KO.md)

**Version:** 1.0.0
**Last Updated:** 2025-11-26
**Status:** Comprehensive

Welcome to the Go Container System documentation! A high-performance, type-safe container framework for Go with comprehensive data management capabilities for messaging systems and general-purpose applications.

---

## Quick Navigation

| I want to... | Document |
|--------------|----------|
| Understand the system design | [Architecture](ARCHITECTURE.md) |
| Learn about all features | [Features](FEATURES.md) |
| See the complete API | [API Reference](API_REFERENCE.md) |
| Understand project structure | [Project Structure](PROJECT_STRUCTURE.md) |
| Find answers to common questions | [FAQ](guides/FAQ.md) |
| Fix issues | [Troubleshooting](guides/TROUBLESHOOTING.md) |
| Review performance | [Benchmarks](performance/BENCHMARKS.md) |
| Migrate from C++ | [Migration Guide](../MIGRATION_GUIDE.md) |

---

## Documentation Structure

### Core Documentation

| Document | Description | Lines |
|----------|-------------|-------|
| [ARCHITECTURE.md](ARCHITECTURE.md) | System architecture, design patterns, Go idioms | 1290+ |
| [FEATURES.md](FEATURES.md) | Complete feature documentation with examples | 650+ |
| [API_REFERENCE.md](API_REFERENCE.md) | Complete API documentation with Go examples | 900+ |
| [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) | Project organization, file structure, build system | 700+ |
| [CHANGELOG.md](CHANGELOG.md) | Version history and release notes | 200+ |

### User Guides

| Document | Description | Lines |
|----------|-------------|-------|
| [FAQ.md](guides/FAQ.md) | Frequently asked questions | 300+ |
| [TROUBLESHOOTING.md](guides/TROUBLESHOOTING.md) | Common issues and solutions | 250+ |
| [ARRAY_VALUE_GUIDE.md](ARRAY_VALUE_GUIDE.md) | ArrayValue implementation guide | 200+ |

### Contributing

| Document | Description | Lines |
|----------|-------------|-------|
| [CONTRIBUTING.md](../CONTRIBUTING.md) | Contribution guidelines | 200+ |
| [TESTING.md](contributing/TESTING.md) | Testing strategy and guide | 400+ |

### Performance

| Document | Description | Lines |
|----------|-------------|-------|
| [BENCHMARKS.md](performance/BENCHMARKS.md) | Performance benchmarks and analysis | 500+ |

---

## Project Information

### Current Status

- **Version**: 1.0.0
- **Go Version**: 1.25+
- **License**: BSD 3-Clause
- **Status**: Production Ready

### Key Features

- Type-safe container system with 15 value types
- Multiple serialization formats (String, Binary, JSON, XML, MessagePack)
- Wire protocol compatible with C++, Python, .NET implementations
- Thread-safe operations (opt-in, < 2% overhead)
- Interface-based design for flexibility
- Comprehensive test suite (95%+ coverage)
- High performance (2M MessagePack ops/s)
- Cross-platform support (Linux, macOS, Windows, BSD)

### Language-Specific Advantages

Go implementation provides:
- Simple interface-based design (no complex inheritance)
- Automatic garbage collection (no manual memory management)
- Built-in concurrency with goroutines
- Fast compilation and development cycle
- Explicit error handling (no hidden exceptions)
- Cross-compilation support
- Smaller binary sizes

---

## Quick Start

### Installation

```bash
go get github.com/kcenon/go_container_system
```

### Basic Usage

```go
import (
    "github.com/kcenon/go_container_system/container/core"
    "github.com/kcenon/go_container_system/container/values"
)

// Create container
container := core.NewValueContainer()
container.SetMessageType("user_data")

// Add values
container.AddValue(values.NewStringValue("username", "alice"))
container.AddValue(values.NewInt32Value("age", 30))
container.AddValue(values.NewFloat64Value("balance", 1500.75))

// Serialize
jsonData, _ := container.ToJSON()
binaryData, _ := container.SerializeArray()

// Deserialize
restored := core.NewValueContainer()
restored.DeserializeArray(binaryData)

// Access values
username := restored.GetValue("username", 0)
name, _ := username.ToString()
```

See [README.md](../README.md) for more quick start examples.

---

## Documentation by Use Case

### For New Users

1. **Start Here**: [README.md](../README.md) - Project overview and installation
2. **Learn Features**: [FEATURES.md](FEATURES.md) - Understand what the system can do
3. **Try Examples**: [examples/](../examples/) - Working code examples
4. **Get Help**: [FAQ.md](guides/FAQ.md) - Common questions answered

### For Developers

1. **System Design**: [ARCHITECTURE.md](ARCHITECTURE.md) - Understand the architecture
2. **API Details**: [API_REFERENCE.md](API_REFERENCE.md) - Complete API documentation
3. **Project Layout**: [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) - File organization
4. **Performance**: [BENCHMARKS.md](performance/BENCHMARKS.md) - Performance characteristics

### For Contributors

1. **Guidelines**: [CONTRIBUTING.md](../CONTRIBUTING.md) - How to contribute
2. **Testing**: [TESTING.md](contributing/TESTING.md) - Testing requirements
3. **Architecture**: [ARCHITECTURE.md](ARCHITECTURE.md) - Design principles
4. **Structure**: [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) - Code organization

### For Cross-Language Integration

1. **Migration**: [MIGRATION_GUIDE.md](../MIGRATION_GUIDE.md) - Wire Protocol migration
2. **Wire Protocol**: [API_REFERENCE.md](API_REFERENCE.md#wire-protocol) - Binary format details
3. **Compatibility**: [ARCHITECTURE.md](ARCHITECTURE.md#comparison-with-c-version) - Cross-language comparison
4. **Testing**: [tests/cross_language_test.go](../tests/cross_language_test.go) - Compatibility tests

---

## Feature Documentation

### Core Capabilities

| Feature | Documentation | Code Example |
|---------|--------------|--------------|
| **Value Types** | [FEATURES.md#value-types](FEATURES.md#value-types) | [API_REFERENCE.md#value-types](API_REFERENCE.md#value-types) |
| **Type Safety** | [FEATURES.md#type-safety](FEATURES.md#type-safety) | [API_REFERENCE.md#value-interface](API_REFERENCE.md#value-interface) |
| **Serialization** | [FEATURES.md#serialization-formats](FEATURES.md#serialization-formats) | [API_REFERENCE.md#serialization](API_REFERENCE.md#serialization) |
| **Thread Safety** | [FEATURES.md#thread-safety](FEATURES.md#thread-safety) | [API_REFERENCE.md#thread-safety](API_REFERENCE.md#thread-safety) |
| **Wire Protocol** | [MIGRATION_GUIDE.md](../MIGRATION_GUIDE.md) | [API_REFERENCE.md#wire-protocol](API_REFERENCE.md#wire-protocol) |

### Advanced Features

| Feature | Documentation | Code Example |
|---------|--------------|--------------|
| **Nested Containers** | [FEATURES.md#nested-containers](FEATURES.md#nested-containers) | [examples/advanced_usage.go](../examples/advanced_usage.go) |
| **Array Values** | [ARRAY_VALUE_GUIDE.md](ARRAY_VALUE_GUIDE.md) | [tests/](../tests/) |
| **MessagePack** | [FEATURES.md#serialization-formats](FEATURES.md#serialization-formats) | [API_REFERENCE.md](API_REFERENCE.md) |
| **File I/O** | [FEATURES.md#file-io-operations](FEATURES.md#file-io-operations) | [API_REFERENCE.md](API_REFERENCE.md) |

---

## Architecture Overview

### Design Philosophy

The Go Container System is built on three principles:

1. **Simplicity Through Interfaces**
   - Interface-based polymorphism for type safety
   - Struct embedding for code reuse
   - Package-based organization
   - Explicit error handling

2. **Performance with Safety**
   - Efficient slice operations
   - Automatic garbage collection
   - Opt-in thread safety (< 2% overhead)
   - Zero-copy where possible

3. **Idiomatic Go Design**
   - Multiple return values for errors
   - Interface segregation
   - Zero values are useful
   - Standard library integration

See [ARCHITECTURE.md](ARCHITECTURE.md) for complete details.

### Package Structure

```
container/
├── core/           # Core types and interfaces
│   ├── value_types.go
│   ├── value.go
│   └── container.go
├── values/         # Concrete value implementations
│   ├── bool_value.go
│   ├── numeric_value.go
│   ├── string_value.go
│   ├── bytes_value.go
│   ├── container_value.go
│   └── array_value.go
└── wireprotocol/   # Binary wire protocol
    └── wire_protocol.go
```

See [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) for complete details.

---

## Performance Characteristics

### Serialization Performance (Apple M4 Max)

| Format | Throughput | Size | Best For |
|--------|-----------|------|----------|
| **MessagePack** | 2.0M ops/s | 70% | High-performance binary |
| **Binary** | 500K ops/s | 95% | Network transmission |
| **String** | 330K ops/s | 100% | Debugging, logging |
| **JSON** | 83K ops/s | 180% | Web APIs, human-readable |
| **XML** | 55K ops/s | 220% | Enterprise integration |

### Operation Performance

| Operation | Latency | Memory |
|-----------|---------|--------|
| Container creation | 75 ns | 168 bytes |
| Value addition | 20 ns | 0 bytes (amortized) |
| Value retrieval | 15 ns | 0 bytes |
| Thread-safe ops | +2% | 24 bytes (RWMutex) |

See [BENCHMARKS.md](performance/BENCHMARKS.md) for complete benchmark data.

---

## Cross-Language Compatibility

The Go implementation maintains 100% compatibility with other implementations:

| Feature | C++ | Go | Python | .NET |
|---------|-----|-----|--------|------|
| 15 Value Types | ✅ | ✅ | ✅ | ✅ |
| Wire Protocol | ✅ | ✅ | ✅ | ✅ |
| Binary Serialization | ✅ | ✅ | ✅ | ✅ |
| JSON/XML Support | ✅ | ✅ | ✅ | ✅ |
| Nested Containers | ✅ | ✅ | ✅ | ✅ |
| Thread Safety | ✅ | ✅ | ✅ | ✅ |
| MessagePack | ✅ | ✅ | ✅ | ✅ |

**Related Implementations**:
- [container_system](https://github.com/kcenon/container_system) - C++ implementation
- [python_container_system](https://github.com/kcenon/python_container_system) - Python implementation
- [dotnet_container_system](https://github.com/kcenon/dotnet_container_system) - .NET implementation

See [ARCHITECTURE.md#comparison-with-c-version](ARCHITECTURE.md#comparison-with-c-version) for detailed comparison.

---

## Testing

### Test Coverage

- **Unit Tests**: 95%+ code coverage
- **Integration Tests**: Cross-language compatibility
- **Benchmark Tests**: Performance regression detection
- **Property Tests**: Edge case validation

### Running Tests

```bash
# Run all tests
go test ./tests -v

# Run specific tests
go test ./tests -v -run TestValueContainer

# Run benchmarks
go test ./tests -bench=. -benchmem

# Generate coverage report
go test ./tests -coverprofile=coverage.out
go tool cover -html=coverage.out
```

See [TESTING.md](contributing/TESTING.md) for complete testing guide.

---

## Common Use Cases

### Messaging Systems

```go
message := core.NewValueContainerFull(
    "client_app", "instance_1",
    "server_api", "v2",
    "user_login",
)
message.AddValue(values.NewStringValue("username", "alice"))
data, _ := message.SerializeArray()
// Send data over network
```

### Data Serialization

```go
config := core.NewValueContainer()
config.AddValue(values.NewStringValue("host", "localhost"))
config.AddValue(values.NewInt32Value("port", 8080))
config.SaveToFileJSON("config.json")
```

### Cross-Language Data Exchange

```go
// Receive binary data from C++ implementation
data := receiveFromCpp()
container := core.NewValueContainer()
container.DeserializeArray(data)

// Process and respond
response := container.Copy(false)
response.SwapHeader()
response.AddValue(values.NewStringValue("status", "ok"))
sendToCpp(response.SerializeArray())
```

See [FEATURES.md#real-world-use-cases](FEATURES.md#real-world-use-cases) for more examples.

---

## Getting Help

### Documentation Issues

If documentation is unclear or incomplete:
1. Check the [FAQ](guides/FAQ.md) for common questions
2. Review [TROUBLESHOOTING.md](guides/TROUBLESHOOTING.md) for known issues
3. Search [GitHub Issues](https://github.com/kcenon/go_container_system/issues)
4. Open a new issue with "Documentation" label

### Technical Support

- **GitHub Issues**: [Report bugs or request features](https://github.com/kcenon/go_container_system/issues)
- **GitHub Discussions**: [Ask questions and share ideas](https://github.com/kcenon/go_container_system/discussions)
- **Email**: kcenon@naver.com

### Contributing Documentation

Documentation improvements are always welcome! See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

---

## Version History

See [CHANGELOG.md](CHANGELOG.md) for complete version history and release notes.

**Recent Releases**:
- **v1.0.0** (2025-11-26): Production release with complete feature set
- **v0.1.0** (2025-10-27): Initial release with Long/ULong types

---

## License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](../LICENSE) file for details.

---

## Acknowledgments

- Based on the C++ [container_system](https://github.com/kcenon/container_system)
- Designed for compatibility with the messaging system ecosystem
- Wire protocol compatible with C++, Python, and .NET implementations

---

**Last Updated**: 2025-11-26
**Version**: 1.0.0
**Maintainer**: kcenon (kcenon@naver.com)
