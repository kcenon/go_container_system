# ArrayValue Implementation Guide (Go)

## Overview

`ArrayValue` provides Go's implementation of type-15 heterogeneous array collections with full cross-language compatibility. It leverages Go's interface system, slices, and idiomatic patterns.

## Architecture

### Type System

```
interface Value
├── BoolValue (type BoolValue)
├── Numeric types (Int16Value, Int32Value, etc.)
├── BytesValue (type BytesValue)
├── StringValue (type StringValue)
└── ArrayValue (type ArrayValue) ← Slice-based heterogeneous collection
```

### Package Structure

```
go_container_system/
├── container/
│   ├── core/
│   │   ├── value.go          # Value interface
│   │   ├── value_types.go    # ValueType constants with ArrayValue
│   │   ├── container.go      # ValueContainer
│   │   └── base_value.go     # BaseValue embedded struct
│   ├── values/
│   │   ├── array_value.go    # ArrayValue implementation
│   │   ├── int_value.go
│   │   ├── string_value.go
│   │   └── ...
│   └── wireprotocol/
│       └── wire_protocol.go  # Serialization with Array support
```

### Struct Diagram

```
┌────────────────────────────────────┐
│      interface Value               │
├────────────────────────────────────┤
│ Name() string                      │
│ Type() ValueType                   │
│ ToBytes() []byte                   │
│ ToJSON() (string, error)          │
│ ToXML() (string, error)           │
└────────────────────────────────────┘
                ▲
                │ embeds
                │
┌────────────────────────────────────┐
│      BaseValue                     │
├────────────────────────────────────┤
│ name string                        │
│ valueType ValueType                │
│ data []byte                        │
└────────────────────────────────────┘
                ▲
                │ embeds
                │
┌────────────────────────────────────┐
│      ArrayValue                    │
├────────────────────────────────────┤
│ *BaseValue                         │
│ elements []Value                   │
├────────────────────────────────────┤
│ NewArrayValue(name, ...elems)     │
│ Elements() []Value                 │
│ Count() int                        │
│ IsEmpty() bool                     │
│ Push(element Value)               │
│ PushBack(element Value)           │
│ At(index int) (Value, error)      │
│ Clear()                            │
│ ToBytes() []byte                   │
│ ToJSON() (string, error)          │
│ ToXML() (string, error)           │
└────────────────────────────────────┘
```

## Usage Examples

### Basic Creation

```go
import (
    "fmt"
    "github.com/kcenon/go_container_system/container/core"
    "github.com/kcenon/go_container_system/container/values"
)

// Create empty array
numbers := values.NewArrayValue("numbers")

// Add elements
numbers.Push(values.NewInt32Value("", 10))
numbers.Push(values.NewInt32Value("", 20))
numbers.Push(values.NewInt32Value("", 30))

fmt.Printf("Array has %d elements\n", numbers.Count())
// Output: Array has 3 elements
```

### Constructor with Variadic Parameters

```go
// Create array with initial values
colors := values.NewArrayValue("colors",
    values.NewStringValue("", "red"),
    values.NewStringValue("", "green"),
    values.NewStringValue("", "blue"),
)

fmt.Printf("Created array with %d colors\n", colors.Count())
```

### Heterogeneous Collections

```go
// Mix different value types
mixed := values.NewArrayValue("mixed")

mixed.Push(values.NewInt32Value("", 42))
mixed.Push(values.NewStringValue("", "hello"))
mixed.Push(values.NewFloat64Value("", 3.14))
mixed.Push(values.NewBoolValue("", true))

// Type assertion for element access
if elem, err := mixed.At(0); err == nil {
    if intVal, ok := elem.(*values.Int32Value); ok {
        fmt.Printf("First element: %d\n", intVal.Value())
    }
}
```

### Integration with ValueContainer

```go
import (
    "github.com/kcenon/go_container_system/container/wireprotocol"
)

// Create container
container := core.NewValueContainer()
container.SetMessageType("user_data")

// Create array
scores := values.NewArrayValue("test_scores")
scores.Push(values.NewInt32Value("", 95))
scores.Push(values.NewInt32Value("", 87))
scores.Push(values.NewInt32Value("", 92))

// Add to container
container.AddValue(scores)

// Serialize to wire protocol
wireData, err := wireprotocol.SerializeCppWire(container)
if err != nil {
    panic(err)
}

fmt.Println("Wire format:", wireData)
// Output: @header={{...}};@data={{[test_scores,array_value,3];}};
```

## Iteration and Access

### Safe Element Access

```go
array := values.NewArrayValue("data",
    values.NewInt32Value("", 10),
    values.NewInt32Value("", 20),
    values.NewInt32Value("", 30),
)

// Error-based access (idiomatic Go)
if element, err := array.At(0); err == nil {
    fmt.Println("First element:", element.Name())
} else {
    fmt.Println("Error accessing element:", err)
}

// Out of bounds returns error
_, err := array.At(100)
if err != nil {
    fmt.Println("Index out of bounds")
}
```

### Iterating Elements

```go
array := values.NewArrayValue("items")
// ... populate array ...

// Iterate through slice
for i, element := range array.Elements() {
    fmt.Printf("Element %d: %s\n", i, element.Name())
}

// Get count
count := array.Count()
fmt.Printf("Total elements: %d\n", count)

// Check if empty
if array.IsEmpty() {
    fmt.Println("Array is empty")
}
```

### Mutable Operations

```go
array := values.NewArrayValue("mutable")

// Add elements
array.Push(values.NewInt32Value("", 1))
array.PushBack(values.NewInt32Value("", 2)) // Alias for Push

// Clear all elements
array.Clear()
fmt.Println("Is empty:", array.IsEmpty()) // Output: Is empty: true
```

## Serialization

### Binary Format

```go
array := values.NewArrayValue("data",
    values.NewInt32Value("", 42),
    values.NewStringValue("", "test"),
)

// Serialize to bytes
bytes := array.ToBytes()
fmt.Printf("Serialized %d bytes\n", len(bytes))
```

### JSON Format

```go
colors := values.NewArrayValue("colors",
    values.NewStringValue("", "red"),
    values.NewStringValue("", "blue"),
)

jsonStr, err := colors.ToJSON()
if err != nil {
    panic(err)
}
fmt.Println("JSON:", jsonStr)
// Output: {"name":"colors","type":"array","elements":[...]}
```

### XML Format

```go
scores := values.NewArrayValue("scores",
    values.NewInt32Value("", 95),
    values.NewInt32Value("", 87),
)

xmlStr, err := scores.ToXML()
if err != nil {
    panic(err)
}
fmt.Println("XML:", xmlStr)
// Output: <array name="scores" count="2">...</array>
```

### Wire Protocol Format

```go
import "github.com/kcenon/go_container_system/container/wireprotocol"

container := core.NewValueContainer()
array := values.NewArrayValue("items",
    values.NewInt32Value("", 1),
    values.NewInt32Value("", 2),
)
container.AddValue(array)

// Serialize to C++ compatible format
wireFormat, err := wireprotocol.SerializeCppWire(container)
if err != nil {
    panic(err)
}

fmt.Println(wireFormat)
// Output: @header={{...}};@data={{[items,array_value,2];}};
```

## Type Assertions and Switches

### Type Assertion Pattern

```go
element, _ := array.At(0)

// Type assertion with ok idiom
if intVal, ok := element.(*values.Int32Value); ok {
    fmt.Printf("Integer value: %d\n", intVal.Value())
}
```

### Type Switch Pattern

```go
for i := 0; i < array.Count(); i++ {
    elem, _ := array.At(i)

    switch v := elem.(type) {
    case *values.Int32Value:
        fmt.Printf("Int32: %d\n", v.Value())
    case *values.StringValue:
        fmt.Printf("String: %s\n", v.Value())
    case *values.Float64Value:
        fmt.Printf("Float64: %f\n", v.Value())
    case *values.ArrayValue:
        fmt.Printf("Nested array with %d elements\n", v.Count())
    default:
        fmt.Printf("Unknown type: %d\n", v.Type())
    }
}
```

### ValueType Check

```go
element, _ := array.At(0)

if element.Type() == core.IntValue {
    if intVal, ok := element.(*values.Int32Value); ok {
        fmt.Println("It's an integer:", intVal.Value())
    }
}
```

## Error Handling

```go
func processArray(array *values.ArrayValue) error {
    if array.IsEmpty() {
        return fmt.Errorf("array cannot be empty")
    }

    for i := 0; i < array.Count(); i++ {
        elem, err := array.At(i)
        if err != nil {
            return fmt.Errorf("failed to access element %d: %w", i, err)
        }

        // Process element
        jsonStr, err := elem.ToJSON()
        if err != nil {
            return fmt.Errorf("failed to serialize element %d: %w", i, err)
        }
        fmt.Println("Element JSON:", jsonStr)
    }

    return nil
}
```

## Cross-Language Interoperability

### Receiving from C++/Rust

```go
import "github.com/kcenon/go_container_system/container/wireprotocol"

// Receive wire format from C++/Rust
cppWireData := `@header={{[5,test];}};@data={{[nums,array_value,3];}};`

// Deserialize
container, err := wireprotocol.DeserializeCppWire(cppWireData)
if err != nil {
    panic(err)
}

// Extract ArrayValue (creates empty placeholder currently)
for _, value := range container.Values() {
    if arrayVal, ok := value.(*values.ArrayValue); ok {
        fmt.Printf("Received array '%s' with count %d\n",
            arrayVal.Name(), arrayVal.Count())
    }
}
```

### Sending to Python/.NET

```go
container := core.NewValueContainer()
array := values.NewArrayValue("data",
    values.NewInt32Value("", 100),
    values.NewStringValue("", "test"),
)
container.AddValue(array)

// Serialize for cross-language transmission
wireData, err := wireprotocol.SerializeCppWire(container)
if err != nil {
    panic(err)
}

// Send wireData over network/IPC
// Python: container = ValueContainer.from_string(wire_data)
// C#: var container = ValueContainer.FromString(wireData);
```

## Best Practices

### 1. Use Variadic Constructor

```go
// Good: Concise initialization
array := values.NewArrayValue("colors",
    values.NewStringValue("", "red"),
    values.NewStringValue("", "green"),
    values.NewStringValue("", "blue"),
)

// Avoid: Verbose step-by-step
array := values.NewArrayValue("colors")
array.Push(values.NewStringValue("", "red"))
array.Push(values.NewStringValue("", "green"))
array.Push(values.NewStringValue("", "blue"))
```

### 2. Check Errors

```go
// Good: Always check errors from At()
if elem, err := array.At(index); err == nil {
    // Use elem
} else {
    return fmt.Errorf("index error: %w", err)
}

// Avoid: Ignoring errors
elem, _ := array.At(index) // Dangerous if index invalid
```

### 3. Use Range for Iteration

```go
// Good: Idiomatic Go
for i, elem := range array.Elements() {
    fmt.Printf("%d: %s\n", i, elem.Name())
}

// Avoid: Manual indexing (more error-prone)
for i := 0; i < array.Count(); i++ {
    elem, _ := array.At(i)
    fmt.Printf("%d: %s\n", i, elem.Name())
}
```

### 4. Type Descriptive Names

```go
// Good
userIDs := values.NewArrayValue("user_ids")
testScores := values.NewArrayValue("test_scores")

// Avoid
arr := values.NewArrayValue("a")
data := values.NewArrayValue("d")
```

## Testing

```go
package values_test

import (
    "testing"
    "github.com/kcenon/go_container_system/container/values"
)

func TestArrayCreation(t *testing.T) {
    array := values.NewArrayValue("test")
    if array.Count() != 0 {
        t.Errorf("Expected empty array, got %d elements", array.Count())
    }
    if !array.IsEmpty() {
        t.Error("Expected array to be empty")
    }
}

func TestHeterogeneousArray(t *testing.T) {
    array := values.NewArrayValue("mixed",
        values.NewInt32Value("", 42),
        values.NewStringValue("", "hello"),
    )

    if array.Count() != 2 {
        t.Errorf("Expected 2 elements, got %d", array.Count())
    }

    // Test first element type
    elem, err := array.At(0)
    if err != nil {
        t.Fatalf("Failed to get element: %v", err)
    }

    if _, ok := elem.(*values.Int32Value); !ok {
        t.Error("Expected first element to be Int32Value")
    }
}

func TestOutOfBounds(t *testing.T) {
    array := values.NewArrayValue("test")

    _, err := array.At(0)
    if err == nil {
        t.Error("Expected error for out of bounds access")
    }
}

func TestClear(t *testing.T) {
    array := values.NewArrayValue("test",
        values.NewInt32Value("", 1),
    )

    if array.IsEmpty() {
        t.Error("Array should not be empty")
    }

    array.Clear()

    if !array.IsEmpty() {
        t.Error("Array should be empty after Clear()")
    }
}
```

## Performance Considerations

- **Slice growth**: Go slices grow dynamically; no manual capacity management needed
- **Interface overhead**: Method calls through `Value` interface have small runtime cost
- **Pointer receivers**: All methods use pointer receivers for efficiency
- **Memory**: Each element is a pointer to Value interface

## Common Patterns

### Building from Existing Slice

```go
numbers := []int32{1, 2, 3, 4, 5}

array := values.NewArrayValue("numbers")
for _, num := range numbers {
    array.Push(values.NewInt32Value("", num))
}
```

### Filtering Elements

```go
func filterIntegers(array *values.ArrayValue) *values.ArrayValue {
    filtered := values.NewArrayValue("filtered")

    for _, elem := range array.Elements() {
        if elem.Type() == core.IntValue {
            filtered.Push(elem)
        }
    }

    return filtered
}
```

### Mapping Elements

```go
func doubleIntegers(array *values.ArrayValue) *values.ArrayValue {
    result := values.NewArrayValue("doubled")

    for _, elem := range array.Elements() {
        if intVal, ok := elem.(*values.Int32Value); ok {
            doubled := intVal.Value() * 2
            result.Push(values.NewInt32Value("", doubled))
        }
    }

    return result
}
```

## Migration from Raw Slices

### Before

```go
var values []core.Value
values = append(values, NewInt32Value("", 10))
values = append(values, NewInt32Value("", 20))
```

### After

```go
array := values.NewArrayValue("values")
array.Push(values.NewInt32Value("", 10))
array.Push(values.NewInt32Value("", 20))
```

**Benefits:**
- Type safety with `ValueType::ArrayValue`
- Serialization methods (ToBytes, ToJSON, ToXML)
- Cross-language wire protocol support
- Consistent API with other container systems

## Integration with Go Ecosystem

### JSON Marshal/Unmarshal

```go
import "encoding/json"

// Custom JSON serialization
type JSONArray struct {
    Name     string   `json:"name"`
    Elements []string `json:"elements"`
}

func arrayToJSON(array *values.ArrayValue) ([]byte, error) {
    ja := JSONArray{
        Name:     array.Name(),
        Elements: make([]string, 0, array.Count()),
    }

    for _, elem := range array.Elements() {
        ja.Elements = append(ja.Elements, elem.Name())
    }

    return json.Marshal(ja)
}
```

### Context Usage

```go
import "context"

func processArrayWithContext(ctx context.Context, array *values.ArrayValue) error {
    for i := 0; i < array.Count(); i++ {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            elem, err := array.At(i)
            if err != nil {
                return err
            }
            // Process element
        }
    }
    return nil
}
```

## See Also

- [Go Value Interface Documentation](../container/core/value.go)
- [ValueType Constants](../container/core/value_types.go)
- [Wire Protocol Implementation](../container/wireprotocol/wire_protocol.go)
- [Container Architecture](ARCHITECTURE.md)
