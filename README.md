# Go Container System

[![License](https://img.shields.io/badge/license-BSD%203--Clause-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.25+-blue.svg)](https://golang.org/dl/)

> **Language:** **English** | [한국어](README_KO.md)

## Overview

The Go Container System is a high-performance, type-safe container framework for Go that provides comprehensive data management capabilities for messaging systems and general-purpose applications. This is a Go implementation of the [C++ container_system](https://github.com/kcenon/container_system), designed to provide identical functionality while leveraging Go's strengths.

## Features

### Core Capabilities

- **Type-Safe Value System**: 15 built-in value types with compile-time type checking
  - null, bool, short, ushort, int, uint, long, ulong, llong, ullong
  - float32, float64, bytes, string, container
- **Message Container**: Full-featured message container with header support
  - Source/Target IDs with sub-IDs
  - Message type and version tracking
  - Multiple serialization formats
- **Serialization**: Multiple format support
  - String-based serialization
  - Byte array serialization
  - JSON conversion
  - XML conversion
- **Value Operations**: Rich set of value operations
  - Type conversions
  - Child value management
  - Value queries by name
- **Container Operations**: Comprehensive container management
  - Header manipulation
  - Value add/remove/query
  - Container copy (with/without values)
  - Header swap for response messages
- **Fluent Builder API**: ContainerBuilder pattern for readable container construction
  - Chainable methods for source, target, type, and values
  - Optional thread-safe mode
- **Dependency Injection Support**: Standard interfaces and providers for DI frameworks
  - ContainerFactory interface for easy mocking and testing
  - Google Wire provider set for automatic wiring
  - Compatible with Uber Dig and other DI containers

## Installation

```bash
go get github.com/kcenon/go_container_system
```

## Quick Start

### Creating Simple Values

```go
import (
    "github.com/kcenon/go_container_system/container/core"
    "github.com/kcenon/go_container_system/container/values"
)

// Create different value types
boolVal := values.NewBoolValue("enabled", true)
intVal := values.NewInt32Value("count", 42)
stringVal := values.NewStringValue("message", "Hello!")

// Type conversions
if val, err := intVal.ToInt32(); err == nil {
    fmt.Printf("Value: %d\n", val)
}
```

### Creating a Message Container

```go
// Create container with full header
container := core.NewValueContainerFull(
    "client_app", "instance_1",  // Source
    "server_api", "v2",           // Target
    "user_registration",          // Message type
)

// Add values
container.AddValue(values.NewStringValue("username", "alice"))
container.AddValue(values.NewInt32Value("age", 30))
container.AddValue(values.NewStringValue("email", "alice@example.com"))

// Serialize
serialized, _ := container.Serialize()
jsonStr, _ := container.ToJSON()
xmlStr, _ := container.ToXML()
```

### Working with Container Values

```go
// Create a nested structure
userData := values.NewContainerValue("user",
    values.NewStringValue("name", "Bob"),
    values.NewInt32Value("age", 25),
)

// Add to parent container
container.AddValue(userData)

// Retrieve values
name := container.GetValue("name", 0)
if str, err := name.ToString(); err == nil {
    fmt.Printf("Name: %s\n", str)
}
```

## Architecture

### Package Structure

```
go_container_system/
├── container/
│   ├── core/           # Core types and interfaces
│   │   ├── value_types.go   # Value type enumeration
│   │   ├── value.go         # Value interface and base implementation
│   │   └── container.go     # ValueContainer implementation
│   ├── di/             # Dependency injection support
│   │   └── provider.go      # ContainerFactory interface and provider
│   ├── messaging/      # Fluent builder API
│   │   └── builder.go       # ContainerBuilder implementation
│   └── values/         # Concrete value implementations
│       ├── bool_value.go
│       ├── numeric_value.go
│       ├── string_value.go
│       ├── bytes_value.go
│       └── container_value.go
├── examples/           # Usage examples
├── tests/             # Test suites
└── README.md
```

### Value Type Hierarchy

```
Value (interface)
├── BaseValue (base implementation)
│   ├── BoolValue
│   ├── Int16Value, UInt16Value
│   ├── Int32Value, UInt32Value
│   ├── Int64Value, UInt64Value
│   ├── Float32Value, Float64Value
│   ├── StringValue
│   ├── BytesValue
│   └── ContainerValue
└── ValueContainer (message container)
```

## Value Types

### Numeric Types

| Type | Go Type | Size | Description |
|------|---------|------|-------------|
| ShortValue | int16 | 2 bytes | 16-bit signed integer |
| UShortValue | uint16 | 2 bytes | 16-bit unsigned integer |
| IntValue | int32 | 4 bytes | 32-bit signed integer |
| UIntValue | uint32 | 4 bytes | 32-bit unsigned integer |
| LongValue | int32 | 4 bytes | 32-bit signed integer (compatibility) |
| ULongValue | uint32 | 4 bytes | 32-bit unsigned integer (compatibility) |
| LLongValue | int64 | 8 bytes | 64-bit signed integer |
| ULLongValue | uint64 | 8 bytes | 64-bit unsigned integer |
| FloatValue | float32 | 4 bytes | 32-bit floating point |
| DoubleValue | float64 | 8 bytes | 64-bit floating point |

### Other Types

- **BoolValue**: Boolean (true/false)
- **StringValue**: UTF-8 string
- **BytesValue**: Binary data
- **ContainerValue**: Nested container with child values
- **NullValue**: Empty/null value

## API Reference

### Value Interface

```go
type Value interface {
    // Basic accessors
    Name() string
    Type() ValueType
    Data() []byte
    Size() int

    // Type checking
    IsNull() bool
    IsBoolean() bool
    IsNumeric() bool
    IsString() bool
    IsBytes() bool
    IsContainer() bool

    // Type conversions
    ToBool() (bool, error)
    ToInt32() (int32, error)
    ToInt64() (int64, error)
    ToFloat32() (float32, error)
    ToFloat64() (float64, error)
    ToString() (string, error)
    ToBytes() ([]byte, error)

    // Serialization
    Serialize() (string, error)
    ToXML() (string, error)
    ToJSON() (string, error)

    // Container operations
    Children() []Value
    ChildCount() int
    GetChild(name string, index int) Value
    AddChild(child Value) error
    RemoveChild(name string) error
}
```

### ValueContainer

```go
// Creation
container := core.NewValueContainer()
container := core.NewValueContainerWithType(messageType, values...)
container := core.NewValueContainerFull(sourceID, sourceSubID, targetID, targetSubID, messageType, values...)

// Header operations
container.SetSource(sourceID, sourceSubID)
container.SetTarget(targetID, targetSubID)
container.SetMessageType(messageType)
container.SwapHeader()

// Value operations
container.AddValue(value)
container.RemoveValue(name)
value := container.GetValue(name, index)
values := container.GetValues(name)
container.ClearValues()

// Container operations
copy := container.Copy(containingValues)

// Serialization
serialized, _ := container.Serialize()
bytes, _ := container.SerializeArray()
json, _ := container.ToJSON()
xml, _ := container.ToXML()

// Deserialization
container.Deserialize(data)
container.DeserializeArray(bytes)
```

### ContainerBuilder (Fluent API)

```go
import "github.com/kcenon/go_container_system/container/messaging"

// Create container using fluent builder pattern
container, err := messaging.NewContainerBuilder().
    WithSource("client", "1").
    WithTarget("server", "main").
    WithType("request").
    WithValues(
        values.NewStringValue("action", "login"),
        values.NewStringValue("user", "alice"),
    ).
    WithThreadSafe(true).
    Build()
```

### Dependency Injection

```go
import "github.com/kcenon/go_container_system/container/di"

// Using ContainerFactory directly
factory := di.NewContainerFactory()
container := factory.NewContainer()
container = factory.NewContainerWithType("request")
builder := factory.NewBuilder()

// Using with Google Wire
// wire.go
//go:build wireinject
// +build wireinject

package main

import (
    "github.com/google/wire"
    "github.com/kcenon/go_container_system/container/di"
)

var ProviderSet = wire.NewSet(
    di.NewContainerFactory,
    wire.Bind(new(di.ContainerFactory), new(*di.DefaultContainerFactory)),
)

func InitializeApp() (*App, error) {
    wire.Build(ProviderSet, NewApp)
    return nil, nil
}

// Using with Uber Dig
container := dig.New()
container.Provide(di.NewContainerFactory)
```

## Examples

See the [examples](examples/) directory for complete working examples:

- `basic_usage.go`: Comprehensive example showing all major features

Run examples:
```bash
go run examples/basic_usage.go
```

## Testing

Run all tests:
```bash
go test ./tests -v
```

Run specific tests:
```bash
go test ./tests -v -run TestBoolValue
go test ./tests -v -run TestValueContainer
```

## Compatibility with C++ Version

This Go implementation provides the same functionality as the C++ container_system:

### Identical Features
- ✅ 15 value types with same semantics
- ✅ Value container with header support
- ✅ String and byte array serialization
- ✅ XML and JSON conversion
- ✅ Container copy operations
- ✅ Header swap functionality
- ✅ Value query by name and index

### Go-Specific Improvements
- 🔹 Interface-based design for better type safety
- 🔹 Error handling using Go idioms (error returns)
- 🔹 Garbage collection (no manual memory management)
- 🔹 Simplified API using Go conventions

### Not Yet Implemented
- ⏳ MessagePack serialization (planned)
- ⏳ File load/save operations (planned)
- ⏳ Thread-safe operations with mutexes (planned)
- ⏳ Memory pool optimization (not needed in Go)

## Project Ecosystem

This container system is designed to work with other ecosystem components:

- **[container_system](https://github.com/kcenon/container_system)**: Original C++ implementation
- **[messaging_system](https://github.com/kcenon/messaging_system)**: Message passing framework
- **[network_system](https://github.com/kcenon/network_system)**: Network communication layer

## Use Cases

- **Message Passing**: Structured message containers for IPC
- **Network Protocols**: Binary serialization for network communication
- **Configuration**: Flexible configuration data structures
- **Data Exchange**: Cross-language data serialization
- **API Communication**: JSON/XML serialization for REST APIs

## Performance Considerations

- **Type Safety**: Compile-time type checking prevents runtime errors
- **Memory Efficiency**: Go's garbage collector manages memory automatically
- **Zero-Copy**: Byte slices use copy-on-write when possible
- **Serialization**: Efficient binary serialization format

## License

This project is licensed under the BSD 3-Clause License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Author

**kcenon**
- Email: kcenon@naver.com
- GitHub: [@kcenon](https://github.com/kcenon)

## Acknowledgments

- Based on the C++ [container_system](https://github.com/kcenon/container_system)
- Designed for compatibility with the messaging system ecosystem
