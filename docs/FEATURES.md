# Go Container System Features

**Last Updated**: 2025-11-26

## Overview

This document provides comprehensive details about all features and capabilities of the Go Container System, including value types, serialization formats, advanced features, and integration capabilities.

## Core Capabilities

### Type Safety

- **Strongly-typed value system** with interface-based polymorphism
- **Compile-time type checking** through Go's type system
- **Interface-based storage** with minimal overhead
- **Type validation** at construction and deserialization
- **Explicit error handling** through multiple return values

### Thread Safety

- **Opt-in thread safety** with RWMutex (< 2% overhead)
- **Concurrent read access** - multiple goroutines can safely read simultaneously
- **Thread-safe write operations** via EnableThreadSafe()
- **Lock-free by default** - no synchronization overhead when not needed
- **Race detector validated** - zero data races detected

### Memory Efficiency

- **Interface-based storage** with automatic memory management
- **Garbage collection** - automatic resource cleanup
- **Move semantics** through Go's assignment
- **Zero-copy operations** where possible
- **Slice-based storage** for efficient value management

### Serialization Formats

- **String format** - native pipe-delimited format
- **Binary format** - efficient wire protocol serialization
- **JSON format** - human-readable with standard library support
- **XML format** - standard XML with proper escaping
- **MessagePack format** - high-performance binary format (6x faster than JSON)
- **Automatic format detection** for deserialization

## Value Types

The Go Container System supports 15 distinct value types covering all common data scenarios:

### Primitive Types

| Type | Code | Go Type | Size | Description |
|------|------|---------|------|-------------|
| `NullValue` | 0 | - | 0 bytes | Null/empty value, represents absence of data |
| `BoolValue` | 1 | `bool` | 1 byte | Boolean true/false |

### Integer Types

| Type | Code | Go Type | Size | Range |
|------|------|---------|------|-------|
| `ShortValue` | 2 | `int16` | 2 bytes | -32,768 to 32,767 |
| `UShortValue` | 3 | `uint16` | 2 bytes | 0 to 65,535 |
| `IntValue` | 4 | `int32` | 4 bytes | -2^31 to 2^31-1 |
| `UIntValue` | 5 | `uint32` | 4 bytes | 0 to 2^32-1 |
| `LongValue` | 6 | `int32` | 4 bytes | -2^31 to 2^31-1 (compatibility) |
| `ULongValue` | 7 | `uint32` | 4 bytes | 0 to 2^32-1 (compatibility) |
| `LLongValue` | 8 | `int64` | 8 bytes | -2^63 to 2^63-1 |
| `ULLongValue` | 9 | `uint64` | 8 bytes | 0 to 2^64-1 |

### Floating-Point Types

| Type | Code | Go Type | Size | Precision |
|------|------|---------|------|-----------|
| `FloatValue` | 10 | `float32` | 4 bytes | ~7 decimal digits |
| `DoubleValue` | 11 | `float64` | 8 bytes | ~15 decimal digits |

### Complex Types

| Type | Code | Go Type | Size | Description |
|------|------|---------|------|-------------|
| `BytesValue` | 12 | `[]byte` | Variable | Raw byte array, binary data |
| `StringValue` | 13 | `string` | Variable | UTF-8 encoded string |
| `ContainerValue` | 14 | Nested container | Variable | Nested container for hierarchical data |

### Value Type Usage Examples

```go
import (
    "github.com/kcenon/go_container_system/container/core"
    "github.com/kcenon/go_container_system/container/values"
)

// Primitive types
nullVal := core.NewBaseValue("empty", core.NullValue, nil)
boolVal := values.NewBoolValue("is_active", true)

// Integer types
int16Val := values.NewInt16Value("tiny", 127)
uint64Val := values.NewUInt64Value("big_number", 1000000000000)

// Floating-point types
float32Val := values.NewFloat32Value("pi_approx", 3.14159)
float64Val := values.NewFloat64Value("price", 99.99)

// Complex types
bytesVal := values.NewBytesValue("data", []byte{0x01, 0x02, 0x03})
stringVal := values.NewStringValue("name", "John Doe")

// Nested container
nested := core.NewValueContainer()
nested.AddValue(values.NewStringValue("city", "Seattle"))
containerVal := values.NewContainerValue("address", nested)
```

## Enhanced Features

### Message Container

- **Header management** - source, target, message_type, version
- **Message routing** - automatic header swapping for responses
- **Multiple values per key** - support for arrays of values
- **Value queries** - efficient retrieval by name and index
- **Container operations** - copy with/without values

```go
import "github.com/kcenon/go_container_system/container/core"

// Create container with full header
message := core.NewValueContainerFull(
    "order_service", "instance_01",      // Source
    "fulfillment_service", "warehouse",  // Target
    "order_create",                       // Message type
)

message.AddValue(values.NewStringValue("order_id", "ORD-2025-001234"))
message.AddValue(values.NewFloat64Value("total_amount", 299.99))

// Serialize for network transmission
data, err := message.Serialize()
```

### Fluent Builder API

The ContainerBuilder provides a fluent, chainable API for constructing ValueContainer instances with improved readability:

```go
import "github.com/kcenon/go_container_system/container/messaging"

// Create container using fluent builder pattern
container, err := messaging.NewContainerBuilder().
    WithSource("client_app", "instance_1").
    WithTarget("server_api", "v2").
    WithType("user_registration").
    WithValues(
        values.NewStringValue("username", "alice"),
        values.NewInt32Value("age", 30),
    ).
    WithThreadSafe(true).
    Build()

if err != nil {
    log.Fatal(err)
}
```

**Builder Methods**:
- `WithSource(id, subID string)` - Set source identification
- `WithTarget(id, subID string)` - Set target identification
- `WithType(messageType string)` - Set message type
- `WithValues(values ...Value)` - Add multiple values (can be called multiple times)
- `WithThreadSafe(enabled bool)` - Enable/disable thread-safe mode
- `Build()` - Construct the final ValueContainer

**Benefits**:
- Readable, self-documenting code
- Optional parameters without constructor overloading
- Compile-time type safety
- Easy to extend with new configuration options

### Wire Protocol Integration

- **Binary serialization** with efficient encoding
- **Cross-language compatibility** with C++, Python, .NET implementations
- **Version support** for protocol evolution
- **Little-endian encoding** for consistency across platforms
- **Array value support** for efficient bulk data transfer

### File I/O Operations

```go
// String format
err := container.SaveToFile("data.txt")
err = container.LoadFromFile("data.txt")

// JSON format
err := container.SaveToFileJSON("data.json")

// XML format
err := container.SaveToFileXML("data.xml")

// MessagePack format (most efficient)
err := container.SaveToFileMessagePack("data.msgpack")
err = container.LoadFromFileMessagePack("data.msgpack")
```

## Advanced Capabilities

### Nested Containers

```go
// Create hierarchical data structures
root := core.NewValueContainer()
root.SetMessageType("order")

customer := core.NewValueContainer()
customer.SetMessageType("customer_info")
customer.AddValue(values.NewStringValue("name", "Alice"))
customer.AddValue(values.NewStringValue("email", "alice@example.com"))

address := core.NewValueContainer()
address.SetMessageType("address")
address.AddValue(values.NewStringValue("street", "123 Main St"))
address.AddValue(values.NewStringValue("city", "Seattle"))

customer.AddValue(values.NewContainerValue("address", address))
root.AddValue(values.NewContainerValue("customer", customer))
```

### Thread-Safe Operations

```go
import (
    "sync"
    "github.com/kcenon/go_container_system/container/core"
)

container := core.NewValueContainer()
container.EnableThreadSafe()

// Multiple goroutines can safely access
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()

        // Thread-safe write
        container.AddValue(values.NewInt32Value(
            fmt.Sprintf("thread_%d", id), int32(id)))

        // Thread-safe read
        value := container.GetValue(fmt.Sprintf("thread_%d", id), 0)
    }(i)
}

wg.Wait()
```

### Dependency Injection Support

The container system provides first-class support for dependency injection frameworks through the `ContainerFactory` interface:

```go
import "github.com/kcenon/go_container_system/container/di"

// Create factory instance
factory := di.NewContainerFactory()

// Use factory methods
container := factory.NewContainer()
containerWithType := factory.NewContainerWithType("request")
containerFull := factory.NewContainerFull("client", "1", "server", "main", "message")

// Get builder from factory
builder := factory.NewBuilder()
```

**ContainerFactory Interface**:
```go
type ContainerFactory interface {
    NewContainer() *core.ValueContainer
    NewContainerWithType(messageType string) *core.ValueContainer
    NewContainerWithTarget(targetID, targetSubID, messageType string) *core.ValueContainer
    NewContainerFull(sourceID, sourceSubID, targetID, targetSubID, messageType string) *core.ValueContainer
    NewBuilder() *messaging.ContainerBuilder
}
```

**Integration with Google Wire**:
```go
// wire.go
//go:build wireinject

package main

import (
    "github.com/google/wire"
    "github.com/kcenon/go_container_system/container/di"
)

var ProviderSet = wire.NewSet(
    di.NewContainerFactory,
    wire.Bind(new(di.ContainerFactory), new(*di.DefaultContainerFactory)),
)
```

**Integration with Uber Dig**:
```go
import "go.uber.org/dig"

container := dig.New()
container.Provide(di.NewContainerFactory)
```

**Benefits**:
- Easy mocking for unit tests
- Decoupled container creation
- Standard interface for DI frameworks
- Consistent factory pattern across the application

### Memory Optimization

```go
// Create container with known size
container := core.NewValueContainer()

// Use efficient value types for known data
container.AddValue(values.NewInt32Value("count", 100))     // 4 bytes
container.AddValue(values.NewFloat64Value("value", 99.99)) // 8 bytes

// Efficient bulk operations
for i := 0; i < 100; i++ {
    container.AddValue(values.NewInt32Value(fmt.Sprintf("val_%d", i), int32(i)))
}
```

## Real-World Use Cases

### Financial Trading System

```go
marketData := core.NewValueContainerFull(
    "trading_engine", "session_001",
    "risk_monitor", "main",
    "market_tick",
)

marketData.AddValue(values.NewStringValue("symbol", "AAPL"))
marketData.AddValue(values.NewFloat64Value("price", 175.50))
marketData.AddValue(values.NewInt64Value("volume", 1000000))
marketData.AddValue(values.NewFloat64Value("bid", 175.48))
marketData.AddValue(values.NewFloat64Value("ask", 175.52))

// High-frequency serialization
binaryData, _ := marketData.SerializeArray()
```

### IoT Sensor Data

```go
sensorReading := core.NewValueContainer()
sensorReading.SetSource("sensor_array", "building_A_floor_3")
sensorReading.SetMessageType("environmental_reading")

// Bulk sensor data
temperatures := make([]float64, 1000)
for i := range temperatures {
    temperatures[i] = 20.0 + float64(i)*0.01
}

// Store as bytes for efficiency
tempBytes := make([]byte, len(temperatures)*8)
for i, t := range temperatures {
    binary.LittleEndian.PutUint64(tempBytes[i*8:], math.Float64bits(t))
}

sensorReading.AddValue(values.NewBytesValue("temp_data", tempBytes))
```

### Web API Response

```go
apiResponse := core.NewValueContainerWithType("api_response")
apiResponse.AddValue(values.NewInt32Value("status", 200))
apiResponse.AddValue(values.NewBoolValue("success", true))

// Nested data object
userData := core.NewValueContainer()
userData.AddValue(values.NewInt32Value("user_id", 12345))
userData.AddValue(values.NewStringValue("username", "john_doe"))
userData.AddValue(values.NewStringValue("email", "john@example.com"))

apiResponse.AddValue(values.NewContainerValue("data", userData))

// JSON serialization for HTTP response
jsonResponse, _ := apiResponse.ToJSON()
```

### Database Storage

```go
record := core.NewValueContainer()
record.SetMessageType("user_record")
record.AddValue(values.NewInt64Value("id", 12345))
record.AddValue(values.NewStringValue("username", "alice"))
record.AddValue(values.NewFloat64Value("balance", 1500.75))
record.AddValue(values.NewBoolValue("active", true))

// Compact binary format for BLOB storage
blobData, _ := record.SerializeArray()
// Store blobData in database BLOB field
```

## Integration Capabilities

### With Messaging Systems

```go
// Create message for inter-service communication
message := core.NewValueContainerFull(
    "auth_service", "instance_01",
    "user_service", "main",
    "user_login",
)

message.AddValue(values.NewStringValue("username", "alice"))
message.AddValue(values.NewStringValue("token", "jwt_token_here"))

// Serialize and send
data, _ := message.Serialize()
// Send data over network/message queue
```

### With Network Systems

```go
// Serialize and send over network
data, _ := container.Serialize()
conn.Write(data)

// Receive and deserialize
buffer := make([]byte, 4096)
n, _ := conn.Read(buffer)
restored := core.NewValueContainer()
restored.Deserialize(string(buffer[:n]))
```

### With gRPC Systems

```go
// Convert container to Protocol Buffer message
protoMsg := &pb.ContainerMessage{
    SourceId: container.SourceID(),
    TargetId: container.TargetID(),
    MessageType: container.MessageType(),
    Data: binaryData,
}

// Send via gRPC
client.SendMessage(ctx, protoMsg)
```

## Performance Characteristics

### Serialization Performance (Apple M4 Max)

| Format | Throughput (ops/s) | Size Overhead | Best Use Case |
|--------|-------------------|---------------|---------------|
| **String** | ~330K | 100% (baseline) | Simple debugging, logging |
| **Binary** | ~500K | ~95% | High-performance, network transmission |
| **MessagePack** | 2.0M | ~70% | High-performance binary serialization |
| **JSON** | ~83K | ~180% | Human-readable, debugging, APIs |
| **XML** | ~55K | ~220% | Schema validation, enterprise integration |

### Memory Characteristics

| Component | Memory Usage | Notes |
|-----------|--------------|-------|
| Empty Container | ~168 bytes | Header + minimal allocations |
| String Value | ~48 bytes + length | Includes key + value |
| Numeric Value | ~32 bytes | Fixed-size allocation |
| Nested Container | Recursive | Sum of all child containers |

### Thread Safety Overhead

| Operation | Without Thread Safety | With Thread Safety | Overhead |
|-----------|---------------------|-------------------|----------|
| Container creation | 75 ns | 77 ns | <2% |
| Value addition | 20 ns | 21 ns | ~5% |
| Value retrieval | 15 ns | 16 ns | ~7% |

## Error Handling

### Go Idiomatic Error Handling

The Go Container System uses Go's multiple return values for explicit error handling:

```go
// Type conversions return (value, error)
value := container.GetValue("user_id", 0)
id, err := value.ToInt32()
if err != nil {
    log.Printf("Conversion error: %v", err)
    return err
}

// Serialization returns (data, error)
jsonData, err := container.ToJSON()
if err != nil {
    return fmt.Errorf("JSON serialization failed: %w", err)
}

// File operations return error
err = container.SaveToFile("data.txt")
if err != nil {
    return fmt.Errorf("file save failed: %w", err)
}
```

### Error Categories

| Category | Example | Handling |
|----------|---------|----------|
| Type Conversion | Converting string to int | Check error, use default value |
| Null Value | Operating on null values | Check IsNull() before conversion |
| Serialization | JSON/XML encoding errors | Propagate error with context |
| File I/O | File read/write failures | Wrap error with file path |
| Range Validation | Long/ULong overflow | Check range before creation |

## Technology Stack

### Modern Go Foundation

- **Go 1.25+** - generics, interfaces, goroutines
- **Interface-based polymorphism** - clean abstractions
- **Garbage collection** - automatic memory management
- **Goroutines** - lightweight concurrency
- **Standard library** - encoding/json, encoding/xml, sync
- **Third-party libraries** - msgpack for efficient binary serialization

### Design Patterns

- **Interface Pattern** - Value interface for type abstraction
- **Factory Pattern** - Constructor functions for value creation
- **Strategy Pattern** - Multiple serialization formats
- **Composite Pattern** - Nested container support
- **Opt-in Pattern** - Thread safety only when needed

## Language-Specific Features

### Go Advantages

1. **Simplicity**: Clean syntax, easy to read and maintain
2. **Fast Compilation**: Instant feedback during development
3. **Garbage Collection**: No manual memory management
4. **Built-in Concurrency**: Goroutines and channels
5. **Cross-compilation**: Build for multiple platforms easily
6. **Standard Library**: Rich standard library with JSON/XML support

### Go Idioms Used

- **Multiple return values** for error handling
- **Struct embedding** for composition
- **Interface segregation** for flexibility
- **Package-level functions** for constructors
- **Zero values are useful** (empty containers are valid)
- **Defer statements** for resource cleanup

## Cross-Language Compatibility

The Go Container System maintains 100% compatibility with other implementations:

| Feature | C++ | Go | Python | .NET |
|---------|-----|-----|--------|------|
| 15 Value Types | ✅ | ✅ | ✅ | ✅ |
| Wire Protocol | ✅ | ✅ | ✅ | ✅ |
| Binary Serialization | ✅ | ✅ | ✅ | ✅ |
| JSON/XML Support | ✅ | ✅ | ✅ | ✅ |
| Nested Containers | ✅ | ✅ | ✅ | ✅ |
| Thread Safety | ✅ | ✅ | ✅ | ✅ |
| MessagePack | ✅ | ✅ | ✅ | ✅ |

## See Also

- [API_REFERENCE.md](API_REFERENCE.md) - Complete API documentation
- [ARCHITECTURE.md](ARCHITECTURE.md) - System design and architecture
- [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) - File organization
- [BENCHMARKS.md](performance/BENCHMARKS.md) - Performance benchmarks
- [FAQ.md](guides/FAQ.md) - Frequently asked questions

---

**Last Updated**: 2025-11-26
**Version**: 1.0.0
