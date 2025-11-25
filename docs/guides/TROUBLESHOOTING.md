# Go Container System - Troubleshooting Guide

> **Version:** 1.0
> **Last Updated:** 2025-11-26

This guide consolidates the most common issues reported while using the Go Container System and explains how to resolve them.

---

## Table of Contents

1. [Serialization Issues](#1-serialization-issues)
2. [Type Conversion Errors](#2-type-conversion-errors)
3. [Cross-Language Compatibility](#3-cross-language-compatibility)
4. [Performance Issues](#4-performance-issues)
5. [Build and Installation](#5-build-and-installation)
6. [Common Error Messages](#6-common-error-messages)

---

## 1. Serialization Issues

### Serialization Returns Empty or Invalid Data

**Symptoms:**
- `Serialize()` returns an empty string
- `SerializeArray()` returns empty byte slice
- Deserialization fails with invalid data

**Checklist:**
1. Ensure the container has values added before serialization
2. Check that value names are not empty strings
3. Verify binary data doesn't contain invalid UTF-8 (for string format)

```go
// Good: Non-empty name
c.AddValue(values.NewStringValue("name", "Alice"))

// Bad: Empty name (may cause issues)
c.AddValue(values.NewStringValue("", "Alice"))

// Check before serialization
if len(c.Values()) == 0 {
    log.Println("Warning: Serializing empty container")
}
```

### JSON Serialization Fails

**Symptoms:**
- `ToJSON()` returns error
- Special characters cause parsing issues

**Solution:**
```go
// Ensure strings don't contain invalid JSON characters
name := strings.ReplaceAll(rawName, "\x00", "")  // Remove null bytes
c.AddValue(values.NewStringValue("name", name))

// For binary data, use BytesValue (automatically base64 encoded in JSON)
c.AddValue(values.NewBytesValue("data", binaryData))
```

### Deserialization Fails with Type Mismatch

**Symptoms:**
- `Deserialize()` returns error
- Values have wrong types after deserialization

**Solution:**
```go
// Check for errors during deserialization
err := c.Deserialize(data)
if err != nil {
    log.Printf("Deserialization error: %v", err)
    // Check if data format matches (string vs binary vs JSON)
}

// Use the correct deserialization method
if strings.HasPrefix(data, "{") {
    c.FromJSON(data)  // JSON format
} else if isBinaryData(data) {
    c.DeserializeArray([]byte(data))  // Binary format
} else {
    c.Deserialize(data)  // String format
}
```

---

## 2. Type Conversion Errors

### Numeric Conversion Fails

**Symptoms:**
- `ToInt32()` returns error for numeric value
- Value appears correct but conversion fails

**Checklist:**
1. Check if the value type matches the conversion method
2. Verify numeric ranges (Int32 vs Int64)
3. Check for nil values

```go
val := c.GetValue("count", 0)
if val == nil {
    log.Println("Value not found")
    return
}

// Check value type before conversion
switch val.Type() {
case core.TypeInt32:
    if num, err := val.ToInt32(); err == nil {
        fmt.Printf("Int32: %d\n", num)
    }
case core.TypeInt64:
    if num, err := val.ToInt64(); err == nil {
        fmt.Printf("Int64: %d\n", num)
    }
case core.TypeString:
    // String to number conversion
    if str, err := val.ToString(); err == nil {
        num, _ := strconv.Atoi(str)
        fmt.Printf("Parsed: %d\n", num)
    }
}
```

### String Conversion Fails for Binary Data

**Symptoms:**
- `ToString()` returns error for BytesValue
- Binary data corrupted when converted to string

**Solution:**
```go
val := c.GetValue("data", 0)
if val.Type() == core.TypeBytes {
    // Use ToBytes() for binary data
    if bytes, err := val.ToBytes(); err == nil {
        // If you need string representation, use base64
        str := base64.StdEncoding.EncodeToString(bytes)
        fmt.Printf("Base64: %s\n", str)
    }
} else {
    // Safe to use ToString()
    str, _ := val.ToString()
    fmt.Println(str)
}
```

---

## 3. Cross-Language Compatibility

### Wire Protocol Data Not Compatible with C++

**Symptoms:**
- C++ server cannot deserialize Go-serialized data
- Type mismatches between languages

**Checklist:**
1. Use `SerializeArray()` for Wire Protocol binary format
2. Ensure value type mappings match (see table below)
3. Check byte order (little-endian)

**Type Mapping:**

| Go Type | C++ Type | Notes |
|---------|----------|-------|
| Int16Value | short_value | 16-bit signed |
| Int32Value | int_value | 32-bit signed |
| Int64Value | llong_value | 64-bit signed |
| Float32Value | float_value | 32-bit float |
| Float64Value | double_value | 64-bit float |
| StringValue | string_value | UTF-8 encoded |
| BytesValue | bytes_value | Raw binary |

```go
// Correct: Use Wire Protocol for C++ interop
bytes, err := c.SerializeArray()
if err != nil {
    log.Fatalf("Wire Protocol serialization failed: %v", err)
}
// Send bytes to C++ server
```

### JSON Format Differs Between Languages

**Symptoms:**
- JSON from Go doesn't parse correctly in other languages
- Field names or structure don't match

**Solution:**
Check JSON format version. Go Container System uses JSON v2.0 format by default:

```json
{
  "message_type": "example",
  "source_id": "go_client",
  "target_id": "server",
  "values": [
    {"name": "key", "type": "string", "value": "data"}
  ]
}
```

See [MIGRATION_GUIDE](../MIGRATION_GUIDE.md) for format details.

---

## 4. Performance Issues

### Slow Serialization Performance

**Symptoms:**
- Serialization takes longer than expected
- High CPU usage during serialization

**Checklist:**
1. Use binary format instead of JSON/XML for speed
2. Avoid serializing unnecessarily large containers
3. Consider batching small containers

```go
// Slow: JSON serialization
json, _ := c.ToJSON()  // ~100K ops/sec

// Fast: Binary serialization
bytes, _ := c.SerializeArray()  // ~300K ops/sec

// Optimization: Batch small containers
batch := core.NewValueContainerWithType("batch")
for _, item := range items {
    itemContainer := values.NewContainerValue(
        fmt.Sprintf("item_%d", item.ID),
        values.NewStringValue("name", item.Name),
        values.NewInt32Value("value", item.Value),
    )
    batch.AddValue(itemContainer)
}
bytes, _ := batch.SerializeArray()  // Single serialization
```

### High Memory Usage

**Symptoms:**
- Memory grows when processing many containers
- Garbage collection pressure

**Solution:**
```go
// Reuse containers when possible
var containerPool = sync.Pool{
    New: func() interface{} {
        return core.NewValueContainer()
    },
}

// Get from pool
c := containerPool.Get().(*core.ValueContainer)

// Process...

// Return to pool (clear first)
c.ClearValues()
containerPool.Put(c)
```

---

## 5. Build and Installation

### Module Not Found

**Symptoms:**
- `go get` fails with "module not found"
- Import errors in IDE

**Solution:**
```bash
# Ensure Go modules are enabled
export GO111MODULE=on

# Clear module cache
go clean -modcache

# Re-fetch dependencies
go mod download
go mod tidy
```

### Version Conflicts

**Symptoms:**
- Multiple versions of dependencies
- Incompatible API changes

**Solution:**
```bash
# Check current version
go list -m github.com/kcenon/go_container_system

# Update to specific version
go get github.com/kcenon/go_container_system@v1.0.0

# Update to latest
go get -u github.com/kcenon/go_container_system
```

---

## 6. Common Error Messages

### "value not found"

**Cause:** Trying to get a value that doesn't exist

**Solution:**
```go
val := c.GetValue("key", 0)
if val == nil {
    // Handle missing value
    log.Println("Key not found, using default")
    val = values.NewStringValue("key", "default")
}
```

### "type conversion not supported"

**Cause:** Converting between incompatible types

**Solution:**
```go
val := c.GetValue("data", 0)
switch val.Type() {
case core.TypeString:
    str, _ := val.ToString()
    // Handle string
case core.TypeInt32:
    num, _ := val.ToInt32()
    // Handle integer
case core.TypeBytes:
    bytes, _ := val.ToBytes()
    // Handle binary
default:
    log.Printf("Unexpected type: %v", val.Type())
}
```

### "invalid wire protocol format"

**Cause:** Corrupted or incomplete binary data

**Solution:**
```go
// Verify data length before deserialization
if len(data) < 4 {
    log.Println("Data too short for Wire Protocol")
    return
}

// Check for corruption
err := c.DeserializeArray(data)
if err != nil {
    log.Printf("Wire Protocol error: %v", err)
    // Try alternative formats
    if jsonErr := c.FromJSON(string(data)); jsonErr == nil {
        log.Println("Data was actually JSON format")
    }
}
```

### "index out of range"

**Cause:** Accessing value at invalid index

**Solution:**
```go
// Check value count before access
values := c.GetValues("item")
if len(values) > index {
    val := values[index]
    // Safe to use
} else {
    log.Printf("Only %d values, requested index %d", len(values), index)
}

// Or use GetValue with index check
val := c.GetValue("item", index)
if val != nil {
    // Safe to use
}
```

---

## Getting Help

If your issue isn't covered here:

1. **Search existing issues**: [GitHub Issues](https://github.com/kcenon/go_container_system/issues)
2. **Ask in discussions**: [GitHub Discussions](https://github.com/kcenon/go_container_system/discussions)
3. **Report new issue** with:
   - Go version (`go version`)
   - Container System version
   - Minimal reproducible example
   - Error messages and stack traces

---

**Last Updated:** 2025-11-26
