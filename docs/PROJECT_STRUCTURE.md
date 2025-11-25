# Go Container System Project Structure

**Last Updated**: 2025-11-26

## Overview

This document provides comprehensive information about the Go Container System project structure, including directory organization, file descriptions, module dependencies, and build artifacts.

## Directory Tree

```
go_container_system/
├── 📁 container/                # Main package
│   ├── 📁 core/                 # Core types and interfaces
│   │   ├── value_types.go       # Value type enumeration (15 types)
│   │   ├── value.go             # Value interface and BaseValue implementation
│   │   └── container.go         # ValueContainer implementation
│   ├── 📁 values/               # Concrete value implementations
│   │   ├── bool_value.go        # Boolean value
│   │   ├── numeric_value.go     # All numeric types (10 types)
│   │   ├── string_value.go      # String value
│   │   ├── bytes_value.go       # Binary data value
│   │   ├── container_value.go   # Nested container value
│   │   ├── array_value.go       # Array value (wire protocol v2)
│   │   └── *_test.go            # Unit tests for values
│   └── 📁 wireprotocol/         # Binary wire protocol
│       └── wire_protocol.go     # Binary serialization/deserialization
├── 📁 examples/                 # Example applications
│   ├── basic_usage.go           # Basic container usage
│   └── advanced_usage.go        # Advanced features (thread safety, etc.)
├── 📁 tests/                    # Test suite
│   ├── container_test.go        # Container unit tests
│   ├── container_bench_test.go  # Performance benchmarks
│   ├── interop_test.go          # Cross-language compatibility tests
│   ├── binary_interop_test.go   # Binary format tests
│   └── cross_language_test.go   # Wire protocol compatibility tests
├── 📁 docs/                     # Documentation
│   ├── README.md                # Documentation index
│   ├── FEATURES.md              # Feature documentation
│   ├── API_REFERENCE.md         # API documentation
│   ├── PROJECT_STRUCTURE.md     # This file
│   ├── ARCHITECTURE.md          # Architecture documentation
│   ├── CHANGELOG.md             # Version history
│   ├── ARRAY_VALUE_GUIDE.md     # ArrayValue implementation guide
│   ├── 📁 guides/               # User guides
│   │   ├── FAQ.md               # Frequently asked questions
│   │   └── TROUBLESHOOTING.md   # Troubleshooting guide
│   ├── 📁 contributing/         # Contribution guides
│   │   └── TESTING.md           # Testing guide
│   └── 📁 performance/          # Performance documentation
│       └── BENCHMARKS.md        # Performance benchmarks
├── 📄 go.mod                    # Go module definition
├── 📄 go.sum                    # Dependency checksums
├── 📄 .gitignore                # Git ignore rules
├── 📄 LICENSE                   # BSD 3-Clause license
├── 📄 README.md                 # Main README
├── 📄 README_KO.md              # Korean README
├── 📄 CONTRIBUTING.md           # Contribution guidelines
└── 📄 MIGRATION_GUIDE.md        # Wire Protocol migration guide
```

## Core Module Files

### Container Core (`container/core/`)

#### `value_types.go`

**Purpose**: Defines the 15 value types supported by the system

**Key Features**:
- Type enumeration using iota (0-14)
- Type checking methods (IsNumeric)
- Human-readable type names
- Cross-language compatibility

**Enumerations**:
```go
const (
    NullValue      ValueType = 0
    BoolValue      ValueType = 1
    ShortValue     ValueType = 2   // int16
    UShortValue    ValueType = 3   // uint16
    IntValue       ValueType = 4   // int32
    UIntValue      ValueType = 5   // uint32
    LongValue      ValueType = 6   // int32 (compatibility)
    ULongValue     ValueType = 7   // uint32 (compatibility)
    LLongValue     ValueType = 8   // int64
    ULLongValue    ValueType = 9   // uint64
    FloatValue     ValueType = 10  // float32
    DoubleValue    ValueType = 11  // float64
    BytesValue     ValueType = 12  // []byte
    StringValue    ValueType = 13  // string
    ContainerValue ValueType = 14  // nested container
)
```

**Public Methods**:
- `String() string` - Returns type name
- `IsNumeric() bool` - Checks if type is numeric

#### `value.go`

**Purpose**: Defines Value interface and BaseValue implementation

**Key Features**:
- Interface-based polymorphism for all value types
- BaseValue provides default implementations
- Type conversions with error returns
- Serialization support
- Container operations

**Public Interface**:
```go
type Value interface {
    // Basic accessors
    Name() string
    Type() ValueType
    Data() []byte
    Size() int

    // Type checking
    IsNull() bool
    IsBytes() bool
    IsBoolean() bool
    IsNumeric() bool
    IsString() bool
    IsContainer() bool

    // Type conversions
    ToBool() (bool, error)
    ToInt16() (int16, error)
    ToUInt16() (uint16, error)
    ToInt32() (int32, error)
    ToUInt32() (uint32, error)
    ToInt64() (int64, error)
    ToUInt64() (uint64, error)
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

**BaseValue Structure**:
```go
type BaseValue struct {
    name   string
    vtype  ValueType
    data   []byte
    parent Value
    units  []Value
}
```

#### `container.go`

**Purpose**: Main container class for managing values with header

**Key Features**:
- Header management (source, target, message_type, version)
- Slice-based value storage
- Multiple serialization formats
- Thread-safe operations (opt-in)
- Copy operations
- Header swapping for request/response

**Public Interface**:
```go
type ValueContainer struct {
    // Unexported fields
}

// Constructor functions
func NewValueContainer() *ValueContainer
func NewValueContainerWithType(messageType string, units ...Value) *ValueContainer
func NewValueContainerWithTarget(targetID, targetSubID, messageType string, units ...Value) *ValueContainer
func NewValueContainerFull(sourceID, sourceSubID, targetID, targetSubID, messageType string, units ...Value) *ValueContainer

// Header methods
func (c *ValueContainer) SetSource(sourceID, sourceSubID string)
func (c *ValueContainer) SetTarget(targetID, targetSubID string)
func (c *ValueContainer) SetMessageType(messageType string)
func (c *ValueContainer) SwapHeader()

// Value management
func (c *ValueContainer) AddValue(value Value)
func (c *ValueContainer) RemoveValue(name string)
func (c *ValueContainer) GetValue(name string, index int) Value
func (c *ValueContainer) GetValues(name string) []Value
func (c *ValueContainer) Values() []Value
func (c *ValueContainer) ClearValues()

// Container operations
func (c *ValueContainer) Copy(containingValues bool) *ValueContainer

// Thread safety
func (c *ValueContainer) EnableThreadSafe()
func (c *ValueContainer) DisableThreadSafe()
func (c *ValueContainer) IsThreadSafe() bool

// Serialization
func (c *ValueContainer) Serialize() (string, error)
func (c *ValueContainer) Deserialize(data string) error
func (c *ValueContainer) SerializeArray() ([]byte, error)
func (c *ValueContainer) DeserializeArray(data []byte) error
func (c *ValueContainer) ToJSON() (string, error)
func (c *ValueContainer) ToXML() (string, error)
func (c *ValueContainer) ToMessagePack() ([]byte, error)
func (c *ValueContainer) FromMessagePack(data []byte) error

// File I/O
func (c *ValueContainer) SaveToFile(filePath string) error
func (c *ValueContainer) LoadFromFile(filePath string) error
func (c *ValueContainer) SaveToFileJSON(filePath string) error
func (c *ValueContainer) SaveToFileXML(filePath string) error
func (c *ValueContainer) SaveToFileMessagePack(filePath string) error
func (c *ValueContainer) LoadFromFileMessagePack(filePath string) error
```

## Value Implementation Files

### Primitive Values (`container/values/`)

#### `bool_value.go`

**Purpose**: Boolean value implementation

**Key Features**:
- Single-byte storage
- Type-safe boolean operations
- Efficient serialization

**Public Interface**:
```go
type BoolValue struct {
    *core.BaseValue
}

func NewBoolValue(name string, value bool) *BoolValue
func (v *BoolValue) Value() bool
```

### Numeric Values (`container/values/numeric_value.go`)

**Purpose**: All numeric type implementations

**Key Features**:
- 10 numeric types (int16 through float64)
- Little-endian binary encoding
- Type-safe conversions
- Consistent API across all types

**Implemented Types**:
- `Int16Value` (int16, 2 bytes)
- `UInt16Value` (uint16, 2 bytes)
- `Int32Value` (int32, 4 bytes)
- `UInt32Value` (uint32, 4 bytes)
- `LongValue` (int32, 4 bytes, compatibility)
- `ULongValue` (uint32, 4 bytes, compatibility)
- `Int64Value` (int64, 8 bytes)
- `UInt64Value` (uint64, 8 bytes)
- `Float32Value` (float32, 4 bytes)
- `Float64Value` (float64, 8 bytes)

**Example**:
```go
type Int32Value struct {
    *core.BaseValue
    value int32
}

func NewInt32Value(name string, value int32) *Int32Value
func (v *Int32Value) Value() int32
func (v *Int32Value) ToInt32() (int32, error)
func (v *Int32Value) ToInt64() (int64, error)
```

#### `string_value.go`

**Purpose**: UTF-8 string value implementation

**Key Features**:
- Immutable string storage
- Proper JSON/XML escaping
- UTF-8 validation
- Efficient string operations

**Public Interface**:
```go
type StringValue struct {
    *core.BaseValue
    value string
}

func NewStringValue(name string, value string) *StringValue
func (v *StringValue) Value() string
func (v *StringValue) ToString() (string, error)
```

#### `bytes_value.go`

**Purpose**: Binary data value implementation

**Key Features**:
- Copy-on-access for safety
- Base64 encoding for text representation
- Efficient binary operations
- Size tracking

**Public Interface**:
```go
type BytesValue struct {
    *core.BaseValue
    data []byte
}

func NewBytesValue(name string, data []byte) *BytesValue
func (v *BytesValue) Value() []byte
func (v *BytesValue) ToBytes() ([]byte, error)
```

#### `container_value.go`

**Purpose**: Nested container support

**Key Features**:
- Heterogeneous child types
- Recursive serialization
- Query by name with index support
- Child management operations

**Public Interface**:
```go
type ContainerValue struct {
    *core.BaseValue
    children []core.Value
}

func NewContainerValue(name string, children ...core.Value) *ContainerValue
func (v *ContainerValue) Children() []core.Value
func (v *ContainerValue) GetChild(name string, index int) core.Value
func (v *ContainerValue) AddChild(child core.Value) error
func (v *ContainerValue) RemoveChild(name string) error
```

#### `array_value.go`

**Purpose**: Array value for efficient bulk data (wire protocol v2)

**Key Features**:
- Efficient storage for homogeneous arrays
- Support for numeric and string arrays
- Wire protocol v2 compatibility
- Zero-copy where possible

**Public Interface**:
```go
type ArrayValue struct {
    *core.BaseValue
    elementType core.ValueType
    elements    []interface{}
}

func NewInt32ArrayValue(name string, values []int32) *ArrayValue
func NewStringArrayValue(name string, values []string) *ArrayValue
// ... other array type constructors
```

## Wire Protocol Module

### `container/wireprotocol/wire_protocol.go`

**Purpose**: Binary wire protocol serialization

**Key Features**:
- Little-endian byte order
- Type-length-value encoding
- Cross-language compatibility (C++, Python, .NET)
- Version support
- Efficient for network transmission

**Key Functions**:
```go
func SerializeContainer(container *core.ValueContainer) ([]byte, error)
func DeserializeContainer(data []byte) (*core.ValueContainer, error)
func SerializeValue(value core.Value) ([]byte, error)
func DeserializeValue(data []byte) (core.Value, error)
```

**Wire Protocol Format**:
```
Container:
  [4 bytes: version]
  [4 bytes: source_id length] [source_id]
  [4 bytes: source_sub_id length] [source_sub_id]
  [4 bytes: target_id length] [target_id]
  [4 bytes: target_sub_id length] [target_sub_id]
  [4 bytes: message_type length] [message_type]
  [4 bytes: value count]
  [values...]

Value:
  [1 byte: type]
  [4 bytes: name length] [name]
  [4 bytes: data length] [data]
```

## Examples

### `examples/basic_usage.go`

**Purpose**: Demonstrates basic container usage

**Features Covered**:
- Container creation
- Adding different value types
- Serialization to multiple formats
- Value retrieval
- Type conversions

**Usage**:
```bash
go run examples/basic_usage.go
```

### `examples/advanced_usage.go`

**Purpose**: Demonstrates advanced features

**Features Covered**:
- Thread-safe operations
- Nested containers
- Binary wire protocol
- MessagePack serialization
- File I/O operations

**Usage**:
```bash
go run examples/advanced_usage.go
```

## Test Organization

### Unit Tests (`tests/`)

#### `container_test.go`

**Coverage**: Core container functionality

**Key Test Cases**:
- Container creation and initialization
- Header management
- Value addition and retrieval
- Value removal
- Container copy operations
- Header swapping
- Thread safety

**Running Tests**:
```bash
go test ./tests -v -run TestValueContainer
```

#### `interop_test.go`

**Coverage**: Cross-language compatibility

**Key Test Cases**:
- Binary serialization compatibility
- Value type mapping
- Wire protocol compatibility
- Round-trip serialization

**Running Tests**:
```bash
go test ./tests -v -run TestInterop
```

#### `binary_interop_test.go`

**Coverage**: Binary format interoperability

**Key Test Cases**:
- Binary encoding/decoding
- Little-endian byte order
- Type-length-value encoding
- Wire protocol v2

#### `cross_language_test.go`

**Coverage**: Cross-language data exchange

**Key Test Cases**:
- Data exchange with C++ implementation
- Data exchange with Python implementation
- Data exchange with .NET implementation

### Benchmarks (`tests/container_bench_test.go`)

**Coverage**: Performance benchmarks

**Key Benchmarks**:
- Container creation
- Value addition
- Value retrieval
- String serialization
- Binary serialization
- JSON serialization
- MessagePack serialization
- Thread-safe operations

**Running Benchmarks**:
```bash
go test ./tests -bench=. -benchmem
```

**Sample Output**:
```
BenchmarkContainerCreation-10           15000000        75 ns/op        168 B/op        3 allocs/op
BenchmarkValueAddition-10               50000000        20 ns/op          0 B/op        0 allocs/op
BenchmarkStringSerialization-10           300000      3000 ns/op       1024 B/op       10 allocs/op
BenchmarkBinarySerialization-10           500000      2000 ns/op        512 B/op        5 allocs/op
BenchmarkMessagePackSerialization-10     2000000       500 ns/op        256 B/op        3 allocs/op
```

### Value-Specific Tests

#### `container/values/long_value_test.go`

**Coverage**: LongValue and ULongValue types

**Key Test Cases**:
- Range validation (32-bit limits)
- Error handling for overflow
- Serialization/deserialization
- Type conversions

#### `container/values/array_value_binary_test.go`

**Coverage**: ArrayValue binary serialization

**Key Test Cases**:
- Array serialization
- Array deserialization
- Type-specific arrays
- Wire protocol compatibility

## Build Configuration

### Go Module (`go.mod`)

```go
module github.com/kcenon/go_container_system

go 1.25

require (
    github.com/vmihailenco/msgpack/v5 v5.4.1
)
```

**Dependencies**:
- `msgpack/v5` - MessagePack serialization (optional, for high-performance binary format)

### Build Commands

#### Standard Build

```bash
go build ./...
```

#### Run Tests

```bash
go test ./tests -v
```

#### Run Benchmarks

```bash
go test ./tests -bench=. -benchmem
```

#### Run Examples

```bash
go run examples/basic_usage.go
go run examples/advanced_usage.go
```

#### Generate Coverage Report

```bash
go test ./tests -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Module Dependencies

### Internal Dependencies

```
core (value_types, value, container)
    ↓
values (bool, numeric, string, bytes, container, array)
    ↓
wireprotocol (binary serialization)
    ↓
Application Layer
```

**Dependency Rules**:
- `core` has no internal dependencies
- `values` depends on `core`
- `wireprotocol` depends on `core` and `values`
- `examples` and `tests` depend on all packages

### External Dependencies

| Dependency | Purpose | Version | Required |
|-----------|---------|---------|----------|
| **msgpack/v5** | MessagePack binary serialization | v5.4.1 | Optional |

**Standard Library Dependencies**:
- `encoding/json` - JSON serialization
- `encoding/xml` - XML serialization
- `encoding/binary` - Binary encoding
- `sync` - Thread safety (RWMutex)
- `os` - File I/O
- `fmt`, `errors`, `strings` - Utilities

## File Naming Conventions

### Source Files

- **Value types**: `<type>_value.go` (e.g., `bool_value.go`, `numeric_value.go`)
- **Core files**: `<component>.go` (e.g., `value.go`, `container.go`)
- **Tests**: `<component>_test.go` (e.g., `container_test.go`)
- **Benchmarks**: `<component>_bench_test.go` (e.g., `container_bench_test.go`)

### Documentation Files

- **Markdown**: `<TOPIC>.md` (e.g., `FEATURES.md`)
- **Korean Translation**: `<TOPIC>_KO.md` (e.g., `README_KO.md`)

## Code Organization Best Practices

### Package Structure

```go
// Package declaration
package core

// Imports (grouped)
import (
    // Standard library
    "encoding/json"
    "fmt"
    "sync"

    // Third-party
    "github.com/vmihailenco/msgpack/v5"

    // Internal
    "github.com/kcenon/go_container_system/container/values"
)

// Types and constants
type ValueType int

const (
    NullValue ValueType = 0
    // ...
)

// Exported types
type Value interface {
    // ...
}

// Exported functions
func NewValueContainer() *ValueContainer {
    // ...
}
```

### File Size Guidelines

- **Source Files**: < 500 lines
- **Test Files**: < 1000 lines per module
- **Example Files**: < 300 lines
- **Documentation**: No hard limit (clarity over brevity)

## Cross-Platform Support

### Operating Systems

- Linux (x86_64, ARM64)
- macOS (Intel, Apple Silicon)
- Windows (x86_64)
- FreeBSD, NetBSD, OpenBSD

### Architectures

- amd64 (x86_64)
- arm64 (ARMv8)
- arm (ARMv7)
- 386 (x86)

### Cross-Compilation

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build

# Build for Windows
GOOS=windows GOARCH=amd64 go build

# Build for macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build
```

## See Also

- [FEATURES.md](FEATURES.md) - Complete feature documentation
- [API_REFERENCE.md](API_REFERENCE.md) - Complete API documentation
- [ARCHITECTURE.md](ARCHITECTURE.md) - Architecture guide
- [BENCHMARKS.md](performance/BENCHMARKS.md) - Performance benchmarks
- [TESTING.md](contributing/TESTING.md) - Testing guide

---

**Last Updated**: 2025-11-26
**Version**: 1.0.0
