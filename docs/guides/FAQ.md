# Go Container System - Frequently Asked Questions

> **Version:** 1.0
> **Last Updated:** 2025-11-26
> **Audience:** Users, Developers

This FAQ addresses common questions about the Go Container System, covering containers, values, serialization, and integration.

---

## Table of Contents

1. [General Questions](#general-questions)
2. [Container Basics](#container-basics)
3. [Value Types](#value-types)
4. [Serialization](#serialization)
5. [Performance](#performance)
6. [Integration](#integration)
7. [Advanced Topics](#advanced-topics)

---

## General Questions

### 1. What is the Go Container System?

A type-safe, high-performance container and serialization system for Go:
- **Type-safe containers** using Go interfaces
- **Multiple serialization formats** (Binary, JSON, XML, Wire Protocol)
- **Cross-language compatibility** with C++, .NET, and Python implementations
- **15 value types** with full type checking

```go
import (
    "github.com/kcenon/go_container_system/container/core"
    "github.com/kcenon/go_container_system/container/values"
)

container := core.NewValueContainer()
container.AddValue(values.NewStringValue("name", "John"))
container.AddValue(values.NewInt32Value("age", 30))
container.AddValue(values.NewFloat64Value("score", 95.5))

// Serialize to JSON
json, _ := container.ToJSON()
```

---

### 2. What Go version is required?

**Required:** Go 1.21+

**Recommended:** Go 1.22+ for best performance

---

### 3. How does it differ from encoding/json or other libraries?

| Feature | Go Container System | encoding/json | gob |
|---------|---------------------|---------------|-----|
| Type Safety | Interface-based | Reflection | Reflection |
| Cross-language | Yes (C++, .NET, Python) | JSON only | Go only |
| Binary Format | Yes (Wire Protocol) | No | Yes |
| Nested Containers | Yes | Yes | Yes |
| Type Metadata | Preserved | Lost | Preserved |
| Memory Efficiency | High | Medium | Medium |

---

## Container Basics

### 4. How do I create a container?

```go
import (
    "github.com/kcenon/go_container_system/container/core"
    "github.com/kcenon/go_container_system/container/values"
)

// Empty container
container1 := core.NewValueContainer()

// Container with message type
container2 := core.NewValueContainerWithType("user_profile")

// Full container with header
container3 := core.NewValueContainerFull(
    "client_app", "instance_1",  // Source ID, Sub ID
    "server_api", "v2",           // Target ID, Sub ID
    "user_registration",          // Message type
)
```

---

### 5. How do I add values to a container?

```go
c := core.NewValueContainer()

// Primitive types
c.AddValue(values.NewStringValue("name", "Alice"))
c.AddValue(values.NewInt32Value("age", 25))
c.AddValue(values.NewFloat64Value("height", 165.5))
c.AddValue(values.NewBoolValue("active", true))

// Binary data
c.AddValue(values.NewBytesValue("data", []byte{0x01, 0x02, 0x03}))

// Nested container
userData := values.NewContainerValue("user",
    values.NewStringValue("name", "Bob"),
    values.NewInt32Value("age", 30),
)
c.AddValue(userData)
```

---

### 6. How do I retrieve values?

```go
// Get value by name (first occurrence)
name := c.GetValue("name", 0)
if str, err := name.ToString(); err == nil {
    fmt.Printf("Name: %s\n", str)
}

// Get all values with same name
allTags := c.GetValues("tag")
for _, tag := range allTags {
    if str, err := tag.ToString(); err == nil {
        fmt.Println(str)
    }
}

// Check if value exists
if c.HasValue("email") {
    email := c.GetValue("email", 0)
    // ... use email
}

// Iterate all values
for _, val := range c.Values() {
    fmt.Printf("%s: %v\n", val.Name(), val.Data())
}
```

---

## Value Types

### 7. What value types are supported?

**Numeric Types:**
- `Int16Value`, `UInt16Value` (16-bit integers)
- `Int32Value`, `UInt32Value` (32-bit integers)
- `Int64Value`, `UInt64Value` (64-bit integers)
- `Float32Value`, `Float64Value` (floating point)

**Other Types:**
- `BoolValue` (boolean)
- `StringValue` (UTF-8 string)
- `BytesValue` (binary data)
- `ContainerValue` (nested container)
- `NullValue` (null/empty value)

**Example:**
```go
// All supported types
c.AddValue(values.NewBoolValue("enabled", true))
c.AddValue(values.NewInt16Value("short", 123))
c.AddValue(values.NewInt32Value("int", 12345))
c.AddValue(values.NewInt64Value("long", 123456789))
c.AddValue(values.NewFloat32Value("float", 3.14))
c.AddValue(values.NewFloat64Value("double", 3.14159265359))
c.AddValue(values.NewStringValue("text", "Hello"))
c.AddValue(values.NewBytesValue("binary", []byte{0xAB, 0xCD}))
```

---

### 8. How do I work with arrays?

For array values, use the ArrayValue type:

```go
import "github.com/kcenon/go_container_system/container/values"

// Create array of integers
intArray := values.NewArrayValue("scores",
    values.NewInt32Value("", 85),
    values.NewInt32Value("", 90),
    values.NewInt32Value("", 95),
)
c.AddValue(intArray)

// Create array of strings
stringArray := values.NewArrayValue("tags",
    values.NewStringValue("", "go"),
    values.NewStringValue("", "container"),
    values.NewStringValue("", "serialization"),
)
c.AddValue(stringArray)

// Iterate array
if arr := c.GetValue("scores", 0); arr != nil {
    for i := 0; i < arr.ChildCount(); i++ {
        child := arr.GetChild("", i)
        if val, err := child.ToInt32(); err == nil {
            fmt.Printf("Score %d: %d\n", i, val)
        }
    }
}
```

See [ARRAY_VALUE_GUIDE](../ARRAY_VALUE_GUIDE.md) for detailed usage.

---

### 9. How do I work with nested containers?

```go
// Create nested structure
root := core.NewValueContainerWithType("root")

// Create user with address
address := values.NewContainerValue("address",
    values.NewStringValue("city", "Seoul"),
    values.NewStringValue("country", "Korea"),
    values.NewStringValue("zip", "12345"),
)

user := values.NewContainerValue("user",
    values.NewStringValue("name", "Alice"),
    values.NewInt32Value("age", 30),
    address,
)

root.AddValue(user)

// Access nested values
userVal := root.GetValue("user", 0)
if userVal != nil {
    nameVal := userVal.GetChild("name", 0)
    if name, err := nameVal.ToString(); err == nil {
        fmt.Printf("User name: %s\n", name)
    }

    addrVal := userVal.GetChild("address", 0)
    if addrVal != nil {
        cityVal := addrVal.GetChild("city", 0)
        if city, err := cityVal.ToString(); err == nil {
            fmt.Printf("City: %s\n", city)
        }
    }
}
```

---

## Serialization

### 10. What serialization formats are supported?

**String Format** (human-readable):
```go
serialized, _ := c.Serialize()
restored := core.NewValueContainer()
restored.Deserialize(serialized)
```

**Binary Format** (Wire Protocol, fastest):
```go
bytes, _ := c.SerializeArray()
restored := core.NewValueContainer()
restored.DeserializeArray(bytes)
```

**JSON Format** (interoperable):
```go
json, _ := c.ToJSON()
restored := core.NewValueContainer()
restored.FromJSON(json)
```

**XML Format** (legacy support):
```go
xml, _ := c.ToXML()
restored := core.NewValueContainer()
restored.FromXML(xml)
```

---

### 11. Which format should I use?

| Format | Speed | Size | Use Case |
|--------|-------|------|----------|
| **Binary (Wire Protocol)** | Fastest | Smallest | Network, cross-language |
| **String** | Fast | Medium | Debugging, logging |
| **JSON** | Medium | Medium | REST APIs, configuration |
| **XML** | Slowest | Largest | Legacy systems, SOAP |

---

### 12. How do I serialize to/from JSON?

```go
// Serialize
c := core.NewValueContainerWithType("user_data")
c.AddValue(values.NewStringValue("name", "Alice"))
c.AddValue(values.NewInt32Value("age", 30))

jsonStr, err := c.ToJSON()
if err != nil {
    log.Fatalf("JSON serialization failed: %v", err)
}
// Result: {"message_type":"user_data","values":[{"name":"name","type":"string","value":"Alice"},{"name":"age","type":"int","value":30}]}

// Deserialize
restored := core.NewValueContainer()
if err := restored.FromJSON(jsonStr); err != nil {
    log.Fatalf("JSON deserialization failed: %v", err)
}

name := restored.GetValue("name", 0)
```

---

## Performance

### 13. What is the performance?

**Benchmarks** (Apple M1, Go 1.22):

| Operation | Throughput | Notes |
|-----------|------------|-------|
| Container creation | ~500K/s | Empty container |
| Value addition | ~2M/s | Single value |
| Binary serialize | ~300K/s | Medium container |
| JSON serialize | ~100K/s | Medium container |
| Value access | ~5M/s | By name lookup |

---

### 14. How do I optimize performance?

```go
// 1. Use binary format for speed
bytes, _ := c.SerializeArray()  // Fastest

// 2. Reuse containers when possible
c.ClearValues()  // Clear values but keep container

// 3. Use appropriate value types
// Use Int32Value instead of Int64Value when values fit
c.AddValue(values.NewInt32Value("count", 100))  // More efficient

// 4. Batch operations
// Add multiple values before serialization
for i := 0; i < 100; i++ {
    c.AddValue(values.NewInt32Value(fmt.Sprintf("item_%d", i), i))
}
serialized, _ := c.SerializeArray()  // Single serialization
```

---

## Integration

### 15. How do I integrate with C++ container_system?

Use the Wire Protocol format for cross-language compatibility:

```go
// Go: Serialize
c := core.NewValueContainerWithType("request")
c.AddValue(values.NewStringValue("action", "query"))
bytes, _ := c.SerializeArray()

// Send bytes to C++ server...

// C++ side:
// auto container = value_container::deserialize_from_bytes(bytes);
// auto action = container->get_value<std::string>("action");
```

See [MIGRATION_GUIDE](../MIGRATION_GUIDE.md) for detailed integration.

---

### 16. Can I use MessagePack?

MessagePack support is planned but not yet implemented. Use the Wire Protocol binary format for similar performance:

```go
// Use Wire Protocol (binary) instead
bytes, _ := c.SerializeArray()
```

---

### 17. How do I integrate with REST APIs?

```go
import (
    "encoding/json"
    "net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
    // Parse incoming container
    c := core.NewValueContainer()
    body, _ := io.ReadAll(r.Body)
    c.FromJSON(string(body))

    // Process request
    action := c.GetValue("action", 0)
    // ...

    // Send response
    response := core.NewValueContainerWithType("response")
    response.AddValue(values.NewStringValue("status", "success"))

    jsonResponse, _ := response.ToJSON()
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(jsonResponse))
}
```

---

## Advanced Topics

### 18. Is it thread-safe?

**Read-only operations**: Thread-safe
**Modifications**: Not thread-safe (use external synchronization)

```go
import "sync"

var mutex sync.RWMutex

// Safe read
mutex.RLock()
value := c.GetValue("key", 0)
mutex.RUnlock()

// Safe write
mutex.Lock()
c.AddValue(values.NewStringValue("key", "value"))
mutex.Unlock()

// Or use sync.Map for concurrent access to multiple containers
var containers sync.Map
containers.Store("user_123", userContainer)
```

---

### 19. How do I swap headers for response messages?

```go
// Incoming request
request := core.NewValueContainerFull(
    "client", "client_sub",
    "server", "server_sub",
    "query",
)

// Create response by swapping headers
response := request.Copy(false)  // Copy without values
response.SwapHeader()            // Swap source/target

// Now: source="server/server_sub", target="client/client_sub"
response.SetMessageType("query_response")
response.AddValue(values.NewStringValue("result", "success"))
```

---

### 20. How do I debug containers?

```go
// Print container structure
fmt.Printf("Container: %+v\n", c)

// Print serialized form (human-readable)
serialized, _ := c.Serialize()
fmt.Println(serialized)

// Print JSON (formatted)
jsonStr, _ := c.ToJSON()
var prettyJSON bytes.Buffer
json.Indent(&prettyJSON, []byte(jsonStr), "", "  ")
fmt.Println(prettyJSON.String())

// Print all values
fmt.Printf("Value count: %d\n", len(c.Values()))
for _, v := range c.Values() {
    fmt.Printf("  %s (%s): %v\n", v.Name(), v.Type(), v.Data())
}
```

---

### 21. What are the memory considerations?

**Memory Usage (approximate):**
- **Empty container**: ~200 bytes
- **Per value**: 50-100 bytes (depends on type)
- **String value**: 50 bytes + string length
- **Container value**: 200 bytes + nested values

**Optimization:**
```go
// Use string interning for repeated keys
// (Go's string handling is already efficient)

// Clear containers for reuse
c.ClearValues()

// Use smaller numeric types when possible
c.AddValue(values.NewInt16Value("count", 100))  // 2 bytes vs 8 bytes
```

---

### 22. Where can I find more examples?

**Documentation:**
- [Quick Start](../README.md#quick-start) - 5-minute guide
- [Architecture](../ARCHITECTURE.md) - System design
- [ARRAY_VALUE_GUIDE](../ARRAY_VALUE_GUIDE.md) - Array patterns

**Examples:**
- `examples/basic_usage.go` - Basic operations
- `examples/advanced_usage.go` - Advanced patterns

**Support:**
- [GitHub Issues](https://github.com/kcenon/go_container_system/issues)
- [GitHub Discussions](https://github.com/kcenon/go_container_system/discussions)

---

**Last Updated:** 2025-11-26
**Next Review:** 2026-02-26
