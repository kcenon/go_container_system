# Go Container System API Reference

> **Version**: 1.0.0
> **Last Updated**: 2025-11-26
> **Go Version**: 1.25+

## Table of Contents

1. [Package Structure](#package-structure)
2. [Value Types](#value-types)
3. [Value Interface](#value-interface)
4. [ValueContainer](#valuecontainer)
5. [Concrete Value Types](#concrete-value-types)
6. [Serialization](#serialization)
7. [Wire Protocol](#wire-protocol)

---

## Package Structure

### Import Paths

```go
import (
    "github.com/kcenon/go_container_system/container/core"
    "github.com/kcenon/go_container_system/container/values"
    "github.com/kcenon/go_container_system/container/wireprotocol"
)
```

**Package Organization**:
- `container/core` - Core types, interfaces, and container implementation
- `container/values` - Concrete value type implementations
- `container/wireprotocol` - Binary wire protocol serialization

---

## Value Types

### ValueType Enumeration

**Package**: `github.com/kcenon/go_container_system/container/core`

**Definition**:
```go
type ValueType int

const (
    NullValue      ValueType = 0
    BoolValue      ValueType = 1
    ShortValue     ValueType = 2  // int16
    UShortValue    ValueType = 3  // uint16
    IntValue       ValueType = 4  // int32
    UIntValue      ValueType = 5  // uint32
    LongValue      ValueType = 6  // int32 (compatibility)
    ULongValue     ValueType = 7  // uint32 (compatibility)
    LLongValue     ValueType = 8  // int64
    ULLongValue    ValueType = 9  // uint64
    FloatValue     ValueType = 10 // float32
    DoubleValue    ValueType = 11 // float64
    BytesValue     ValueType = 12 // []byte
    StringValue    ValueType = 13 // string
    ContainerValue ValueType = 14 // nested container
)
```

**Methods**:

#### `String() string`

Returns human-readable type name.

```go
vtype := core.IntValue
fmt.Println(vtype.String()) // Output: "int"
```

#### `IsNumeric() bool`

Checks if type is numeric (int16 through float64).

```go
if core.IntValue.IsNumeric() {
    // Process numeric value
}
```

---

## Value Interface

### Interface Definition

**Package**: `github.com/kcenon/go_container_system/container/core`

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

### Basic Accessors

#### `Name() string`

Returns the name/key of the value.

```go
value := values.NewStringValue("username", "alice")
fmt.Println(value.Name()) // Output: "username"
```

#### `Type() ValueType`

Returns the type enumeration of the value.

```go
value := values.NewInt32Value("count", 42)
fmt.Println(value.Type()) // Output: 4 (IntValue)
```

#### `Data() []byte`

Returns the raw byte representation of the value.

```go
value := values.NewInt32Value("count", 42)
data := value.Data() // []byte{42, 0, 0, 0} (little-endian)
```

#### `Size() int`

Returns the size of the data in bytes.

```go
value := values.NewStringValue("name", "Alice")
fmt.Println(value.Size()) // Output: 5
```

### Type Checking

#### `IsNull() bool`

Checks if value is null type.

```go
value := core.NewBaseValue("empty", core.NullValue, nil)
if value.IsNull() {
    fmt.Println("Value is null")
}
```

#### `IsBoolean() bool`

Checks if value is boolean type.

```go
value := values.NewBoolValue("active", true)
if value.IsBoolean() {
    fmt.Println("Value is boolean")
}
```

#### `IsNumeric() bool`

Checks if value is any numeric type (int16 through float64).

```go
value := values.NewInt32Value("count", 42)
if value.IsNumeric() {
    // Can convert to numeric types
}
```

#### `IsString() bool`

Checks if value is string type.

```go
value := values.NewStringValue("name", "Alice")
if value.IsString() {
    str, _ := value.ToString()
}
```

#### `IsBytes() bool`

Checks if value is bytes type.

```go
value := values.NewBytesValue("data", []byte{0x01, 0x02})
if value.IsBytes() {
    data, _ := value.ToBytes()
}
```

#### `IsContainer() bool`

Checks if value is container type.

```go
nested := core.NewValueContainer()
value := values.NewContainerValue("nested", nested)
if value.IsContainer() {
    children := value.Children()
}
```

### Type Conversions

All conversion methods return `(value, error)` where error is non-nil if conversion fails.

#### `ToBool() (bool, error)`

Converts value to boolean.

```go
value := values.NewBoolValue("active", true)
result, err := value.ToBool()
if err != nil {
    log.Fatal(err)
}
fmt.Println(result) // Output: true
```

**Error Cases**:
- Returns error if value is not boolean type
- Returns error if value is null

#### `ToInt32() (int32, error)`

Converts value to int32.

```go
value := values.NewInt32Value("count", 42)
result, err := value.ToInt32()
if err != nil {
    log.Fatal(err)
}
fmt.Println(result) // Output: 42
```

**Supported Conversions**:
- int16 → int32 (widening)
- uint16 → int32 (widening)
- int32 → int32 (direct)
- Returns error for other types

#### `ToInt64() (int64, error)`

Converts value to int64.

```go
value := values.NewInt64Value("timestamp", 1234567890)
result, err := value.ToInt64()
if err != nil {
    log.Fatal(err)
}
```

**Supported Conversions**:
- All integer types → int64 (widening)
- Returns error for non-integer types

#### `ToFloat64() (float64, error)`

Converts value to float64.

```go
value := values.NewFloat64Value("price", 99.99)
result, err := value.ToFloat64()
if err != nil {
    log.Fatal(err)
}
```

**Supported Conversions**:
- All numeric types → float64
- Returns error for non-numeric types

#### `ToString() (string, error)`

Converts value to string.

```go
value := values.NewStringValue("name", "Alice")
result, err := value.ToString()
if err != nil {
    log.Fatal(err)
}
fmt.Println(result) // Output: "Alice"
```

#### `ToBytes() ([]byte, error)`

Converts value to byte slice.

```go
value := values.NewBytesValue("data", []byte{0x01, 0x02, 0x03})
result, err := value.ToBytes()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%v\n", result) // Output: [1 2 3]
```

### Serialization Methods

#### `Serialize() (string, error)`

Serializes value to native string format.

```go
value := values.NewStringValue("name", "Alice")
serialized, err := value.Serialize()
if err != nil {
    log.Fatal(err)
}
fmt.Println(serialized) // Output: "name|13|5|Alice"
```

**Format**: `name|type|size|data`

#### `ToJSON() (string, error)`

Serializes value to JSON format.

```go
value := values.NewInt32Value("count", 42)
json, err := value.ToJSON()
if err != nil {
    log.Fatal(err)
}
// Output: {"name":"count","type":"int","data":"..."}
```

#### `ToXML() (string, error)`

Serializes value to XML format.

```go
value := values.NewStringValue("name", "Alice")
xml, err := value.ToXML()
if err != nil {
    log.Fatal(err)
}
// Output: <value name="name" type="string">Alice</value>
```

### Container Operations

These methods are primarily for ContainerValue types. Default implementations return empty/error for non-container values.

#### `Children() []Value`

Returns all child values.

```go
container := values.NewContainerValue("parent",
    values.NewStringValue("name", "Alice"),
    values.NewInt32Value("age", 30),
)
children := container.Children()
fmt.Println(len(children)) // Output: 2
```

#### `ChildCount() int`

Returns the number of child values.

```go
fmt.Println(container.ChildCount()) // Output: 2
```

#### `GetChild(name string, index int) Value`

Returns child value by name and index (for duplicate names).

```go
child := container.GetChild("name", 0)
if child != nil {
    name, _ := child.ToString()
    fmt.Println(name) // Output: "Alice"
}
```

#### `AddChild(child Value) error`

Adds a child value to container.

```go
err := container.AddChild(values.NewStringValue("email", "alice@example.com"))
if err != nil {
    log.Fatal(err)
}
```

#### `RemoveChild(name string) error`

Removes all child values with given name.

```go
err := container.RemoveChild("name")
if err != nil {
    log.Fatal(err)
}
```

---

## ValueContainer

### Container Structure

**Package**: `github.com/kcenon/go_container_system/container/core`

```go
type ValueContainer struct {
    // Unexported fields
}
```

### Constructor Functions

#### `NewValueContainer() *ValueContainer`

Creates an empty container with default version.

```go
container := core.NewValueContainer()
```

#### `NewValueContainerWithType(messageType string, units ...Value) *ValueContainer`

Creates container with message type and optional initial values.

```go
container := core.NewValueContainerWithType("user_data",
    values.NewStringValue("name", "Alice"),
    values.NewInt32Value("age", 30),
)
```

#### `NewValueContainerWithTarget(targetID, targetSubID, messageType string, units ...Value) *ValueContainer`

Creates container with target information.

```go
container := core.NewValueContainerWithTarget(
    "server_api", "v2",
    "user_registration",
    values.NewStringValue("username", "alice"),
)
```

#### `NewValueContainerFull(sourceID, sourceSubID, targetID, targetSubID, messageType string, units ...Value) *ValueContainer`

Creates container with full header information.

```go
container := core.NewValueContainerFull(
    "client_app", "instance_1",  // Source
    "server_api", "v2",           // Target
    "user_registration",          // Message type
    values.NewStringValue("username", "alice"),
    values.NewInt32Value("age", 30),
)
```

### Header Methods

#### `SetSource(sourceID, sourceSubID string)`

Sets the source ID and sub-ID.

```go
container.SetSource("auth_service", "instance_01")
```

#### `SourceID() string` / `SourceSubID() string`

Returns source ID or sub-ID.

```go
fmt.Println(container.SourceID())    // "auth_service"
fmt.Println(container.SourceSubID()) // "instance_01"
```

#### `SetTarget(targetID, targetSubID string)`

Sets the target ID and sub-ID.

```go
container.SetTarget("user_service", "main")
```

#### `TargetID() string` / `TargetSubID() string`

Returns target ID or sub-ID.

```go
fmt.Println(container.TargetID())    // "user_service"
fmt.Println(container.TargetSubID()) // "main"
```

#### `SetMessageType(messageType string)`

Sets the message type.

```go
container.SetMessageType("user_login")
```

#### `MessageType() string`

Returns the message type.

```go
fmt.Println(container.MessageType()) // "user_login"
```

#### `SetVersion(version string)`

Sets the protocol version.

```go
container.SetVersion("2.0.0.0")
```

#### `Version() string`

Returns the protocol version.

```go
fmt.Println(container.Version()) // "1.0.0.0" (default)
```

#### `SwapHeader()`

Swaps source and target (useful for response messages).

```go
// Request: client -> server
request.SetSource("client", "1")
request.SetTarget("server", "main")

// Create response by swapping
response := request.Copy(false)
response.SwapHeader()
// Response: server -> client
```

### Value Management

#### `AddValue(value Value)`

Adds a value to the container.

```go
container.AddValue(values.NewStringValue("username", "alice"))
container.AddValue(values.NewInt32Value("age", 30))
```

#### `RemoveValue(name string)`

Removes all values with the given name.

```go
container.RemoveValue("age")
```

#### `GetValue(name string, index int) Value`

Retrieves value by name and index (for duplicate names).

```go
value := container.GetValue("username", 0)
if value != nil {
    name, _ := value.ToString()
    fmt.Println(name)
}
```

**Returns**: Value if found, nil otherwise

#### `GetValues(name string) []Value`

Retrieves all values with the given name.

```go
values := container.GetValues("tag")
for _, val := range values {
    tag, _ := val.ToString()
    fmt.Println(tag)
}
```

#### `Values() []Value`

Returns all values in the container.

```go
allValues := container.Values()
fmt.Println(len(allValues))
```

#### `ClearValues()`

Removes all values from the container.

```go
container.ClearValues()
```

### Container Operations

#### `Copy(containingValues bool) *ValueContainer`

Creates a copy of the container.

```go
// Copy with values
fullCopy := container.Copy(true)

// Copy header only (no values)
headerOnly := container.Copy(false)
```

**Parameters**:
- `containingValues`: If true, copies all values; if false, copies only header

### Thread Safety

#### `EnableThreadSafe()`

Enables thread-safe mode with RWMutex protection.

```go
container.EnableThreadSafe()
// Now safe for concurrent goroutine access
```

#### `DisableThreadSafe()`

Disables thread-safe mode.

```go
container.DisableThreadSafe()
```

#### `IsThreadSafe() bool`

Checks if thread-safe mode is enabled.

```go
if container.IsThreadSafe() {
    fmt.Println("Thread-safe mode enabled")
}
```

### Serialization Methods

#### `Serialize() (string, error)`

Serializes container to native string format.

```go
serialized, err := container.Serialize()
if err != nil {
    log.Fatal(err)
}
fmt.Println(serialized)
```

**Format**: Pipe-delimited header followed by values

#### `Deserialize(data string) error`

Deserializes container from string format.

```go
err := container.Deserialize(serialized)
if err != nil {
    log.Fatal(err)
}
```

#### `SerializeArray() ([]byte, error)`

Serializes container to binary byte array using wire protocol.

```go
binaryData, err := container.SerializeArray()
if err != nil {
    log.Fatal(err)
}
```

#### `DeserializeArray(data []byte) error`

Deserializes container from binary byte array.

```go
err := container.DeserializeArray(binaryData)
if err != nil {
    log.Fatal(err)
}
```

#### `ToJSON() (string, error)`

Serializes container to JSON format.

```go
json, err := container.ToJSON()
if err != nil {
    log.Fatal(err)
}
fmt.Println(json)
```

**Example Output**:
```json
{
  "source_id": "client",
  "source_sub_id": "1",
  "target_id": "server",
  "target_sub_id": "main",
  "message_type": "user_login",
  "version": "1.0.0.0",
  "values": [...]
}
```

#### `ToXML() (string, error)`

Serializes container to XML format.

```go
xml, err := container.ToXML()
if err != nil {
    log.Fatal(err)
}
fmt.Println(xml)
```

#### `ToMessagePack() ([]byte, error)`

Serializes container to MessagePack format.

```go
msgpack, err := container.ToMessagePack()
if err != nil {
    log.Fatal(err)
}
```

#### `FromMessagePack(data []byte) error`

Deserializes container from MessagePack format.

```go
err := container.FromMessagePack(msgpack)
if err != nil {
    log.Fatal(err)
}
```

### File I/O Methods

#### `SaveToFile(filePath string) error`

Saves container to file in string format.

```go
err := container.SaveToFile("data.txt")
if err != nil {
    log.Fatal(err)
}
```

#### `LoadFromFile(filePath string) error`

Loads container from file in string format.

```go
err := container.LoadFromFile("data.txt")
if err != nil {
    log.Fatal(err)
}
```

#### `SaveToFileJSON(filePath string) error`

Saves container to file in JSON format.

```go
err := container.SaveToFileJSON("data.json")
```

#### `SaveToFileXML(filePath string) error`

Saves container to file in XML format.

```go
err := container.SaveToFileXML("data.xml")
```

#### `SaveToFileMessagePack(filePath string) error`

Saves container to file in MessagePack format (most efficient).

```go
err := container.SaveToFileMessagePack("data.msgpack")
```

#### `LoadFromFileMessagePack(filePath string) error`

Loads container from MessagePack file.

```go
err := container.LoadFromFileMessagePack("data.msgpack")
```

---

## Concrete Value Types

### Boolean Value

**Package**: `github.com/kcenon/go_container_system/container/values`

#### `NewBoolValue(name string, value bool) *BoolValue`

Creates a boolean value.

```go
value := values.NewBoolValue("is_active", true)
```

**Methods**:
- `Value() bool` - Returns the boolean value

### Integer Values

#### `NewInt16Value(name string, value int16) *Int16Value`

Creates a 16-bit signed integer value.

```go
value := values.NewInt16Value("temperature", 25)
```

#### `NewUInt16Value(name string, value uint16) *UInt16Value`

Creates a 16-bit unsigned integer value.

```go
value := values.NewUInt16Value("port", 8080)
```

#### `NewInt32Value(name string, value int32) *Int32Value`

Creates a 32-bit signed integer value.

```go
value := values.NewInt32Value("count", 42)
```

#### `NewUInt32Value(name string, value uint32) *UInt32Value`

Creates a 32-bit unsigned integer value.

```go
value := values.NewUInt32Value("user_id", 12345)
```

#### `NewLongValue(name string, value int64) (*LongValue, error)`

Creates a 32-bit signed long value (compatibility type).

```go
value, err := values.NewLongValue("timestamp", 1234567890)
if err != nil {
    // Value out of 32-bit range
}
```

**Note**: Returns error if value exceeds int32 range

#### `NewULongValue(name string, value uint64) (*ULongValue, error)`

Creates a 32-bit unsigned long value (compatibility type).

```go
value, err := values.NewULongValue("counter", 4000000000)
if err != nil {
    // Value out of 32-bit range
}
```

#### `NewInt64Value(name string, value int64) *Int64Value`

Creates a 64-bit signed integer value.

```go
value := values.NewInt64Value("large_number", 9_000_000_000_000)
```

#### `NewUInt64Value(name string, value uint64) *UInt64Value`

Creates a 64-bit unsigned integer value.

```go
value := values.NewUInt64Value("huge_number", 18_000_000_000_000_000_000)
```

**Common Methods** (all integer types):
- `Value() T` - Returns the numeric value of type T

### Floating-Point Values

#### `NewFloat32Value(name string, value float32) *Float32Value`

Creates a 32-bit floating-point value.

```go
value := values.NewFloat32Value("temperature", 98.6)
```

#### `NewFloat64Value(name string, value float64) *Float64Value`

Creates a 64-bit floating-point value.

```go
value := values.NewFloat64Value("price", 99.99)
```

**Common Methods**:
- `Value() float32 | float64` - Returns the floating-point value

### String Value

#### `NewStringValue(name string, value string) *StringValue`

Creates a UTF-8 string value.

```go
value := values.NewStringValue("username", "alice")
```

**Methods**:
- `Value() string` - Returns the string value

### Bytes Value

#### `NewBytesValue(name string, data []byte) *BytesValue`

Creates a binary data value.

```go
data := []byte{0x01, 0x02, 0x03, 0x04}
value := values.NewBytesValue("binary_data", data)
```

**Methods**:
- `Value() []byte` - Returns a copy of the byte data

### Container Value

#### `NewContainerValue(name string, children ...Value) *ContainerValue`

Creates a nested container value.

```go
value := values.NewContainerValue("user",
    values.NewStringValue("name", "Alice"),
    values.NewInt32Value("age", 30),
)
```

**Methods**:
- All container operations from Value interface

---

## Serialization

### Format Comparison

| Format | Speed | Size | Human Readable | Use Case |
|--------|-------|------|----------------|----------|
| String | Fast | Medium | Yes | Debugging, logging |
| Binary | Fastest | Smallest | No | Network, performance |
| MessagePack | Very Fast | Small | No | High-performance RPC |
| JSON | Medium | Large | Yes | Web APIs, config |
| XML | Slow | Largest | Yes | Enterprise integration |

### Complete Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/kcenon/go_container_system/container/core"
    "github.com/kcenon/go_container_system/container/values"
)

func main() {
    // Create container with full header
    container := core.NewValueContainerFull(
        "client_app", "instance_1",
        "server_api", "v2",
        "user_registration",
    )

    // Add various value types
    container.AddValue(values.NewStringValue("username", "alice"))
    container.AddValue(values.NewInt32Value("age", 30))
    container.AddValue(values.NewFloat64Value("balance", 1500.75))
    container.AddValue(values.NewBoolValue("active", true))

    // Add nested container
    address := core.NewValueContainer()
    address.AddValue(values.NewStringValue("street", "123 Main St"))
    address.AddValue(values.NewStringValue("city", "Seattle"))
    container.AddValue(values.NewContainerValue("address", address))

    // Serialize to different formats
    stringData, _ := container.Serialize()
    fmt.Println("String format:", stringData)

    binaryData, _ := container.SerializeArray()
    fmt.Printf("Binary format: %d bytes\n", len(binaryData))

    jsonData, _ := container.ToJSON()
    fmt.Println("JSON format:", jsonData)

    // Deserialize
    restored := core.NewValueContainer()
    restored.DeserializeArray(binaryData)

    // Access values
    username := restored.GetValue("username", 0)
    if username != nil {
        name, _ := username.ToString()
        fmt.Println("Username:", name)
    }
}
```

---

## Wire Protocol

### Binary Serialization

**Package**: `github.com/kcenon/go_container_system/container/wireprotocol`

The wire protocol provides efficient binary serialization compatible with C++, Python, and .NET implementations.

#### Key Features

- Little-endian byte order
- Type-length-value encoding
- Cross-language compatibility
- Version support
- Efficient for network transmission

#### Usage

```go
import "github.com/kcenon/go_container_system/container/wireprotocol"

// Serialize using wire protocol
data, err := container.SerializeArray()
if err != nil {
    log.Fatal(err)
}

// Send over network
conn.Write(data)

// Deserialize
received := core.NewValueContainer()
err = received.DeserializeArray(data)
```

---

## Error Handling Patterns

### Pattern 1: Immediate Check

```go
value := container.GetValue("user_id", 0)
id, err := value.ToInt32()
if err != nil {
    return fmt.Errorf("failed to get user_id: %w", err)
}
```

### Pattern 2: Type Check Before Conversion

```go
value := container.GetValue("count", 0)
if value != nil && value.IsNumeric() {
    count, _ := value.ToInt32()
    fmt.Println("Count:", count)
}
```

### Pattern 3: Default Value

```go
value := container.GetValue("optional_field", 0)
var result int32 = 0 // default
if value != nil {
    if val, err := value.ToInt32(); err == nil {
        result = val
    }
}
```

### Pattern 4: Error Accumulation

```go
var errors []error
for _, value := range container.Values() {
    if _, err := value.Serialize(); err != nil {
        errors = append(errors, err)
    }
}
if len(errors) > 0 {
    return fmt.Errorf("serialization errors: %v", errors)
}
```

---

## Performance Tips

1. **Reuse Containers**: Use `Copy(false)` to reuse container headers
2. **Enable Thread Safety Only When Needed**: < 2% overhead when enabled
3. **Use Binary Format**: 6x faster than JSON for serialization
4. **Preallocate When Possible**: Reduces allocations
5. **Use Appropriate Types**: Match data range to type size
6. **Batch Operations**: Add multiple values before serialization

---

## See Also

- [FEATURES.md](FEATURES.md) - Feature documentation
- [ARCHITECTURE.md](ARCHITECTURE.md) - System architecture
- [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) - Project organization
- [FAQ.md](guides/FAQ.md) - Frequently asked questions
- [BENCHMARKS.md](performance/BENCHMARKS.md) - Performance data

---

**Version**: 1.0.0
**Last Updated**: 2025-11-26
**Go Version**: 1.25+
