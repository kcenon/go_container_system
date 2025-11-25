/****************************************************************************
BSD 3-Clause License

Copyright (c) 2021, 🍀☀🌕🌥 🌊
All rights reserved.
****************************************************************************/

package values

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math"

	"github.com/kcenon/go_container_system/container/core"
)

// ArrayValue represents an array/list that can hold multiple elements
// ArrayValue (type 15) is an extension to support homogeneous or heterogeneous
// collections of values, similar to JSON arrays.
//
// Wire format (binary):
// [type:1=15][name_len:4 LE][name:UTF-8][value_size:4 LE][count:4 LE][values...]
//
// Text format:
// [name,15,count];[element1][element2]...
type ArrayValue struct {
	*core.BaseValue
	elements []core.Value
}

// NewArrayValue creates a new array value
func NewArrayValue(name string, elements ...core.Value) *ArrayValue {
	av := &ArrayValue{
		BaseValue: core.NewBaseValue(name, core.ArrayValue, nil),
		elements:  make([]core.Value, 0),
	}
	for _, element := range elements {
		av.elements = append(av.elements, element)
	}
	return av
}

// Elements returns all elements
func (v *ArrayValue) Elements() []core.Value {
	return v.elements
}

// Count returns the number of elements
func (v *ArrayValue) Count() int {
	return len(v.elements)
}

// IsEmpty checks if the array is empty
func (v *ArrayValue) IsEmpty() bool {
	return len(v.elements) == 0
}

// At gets element at index
func (v *ArrayValue) At(index int) (core.Value, error) {
	if index < 0 || index >= len(v.elements) {
		return nil, fmt.Errorf("ArrayValue index %d out of range (size: %d)", index, len(v.elements))
	}
	return v.elements[index], nil
}

// Append adds an element to the end of the array
func (v *ArrayValue) Append(element core.Value) error {
	v.elements = append(v.elements, element)
	return nil
}

// Push adds an element to the end of the array
func (v *ArrayValue) Push(element core.Value) error {
	return v.Append(element)
}

// PushBack adds an element (C++ compatibility name)
func (v *ArrayValue) PushBack(element core.Value) error {
	return v.Append(element)
}

// Clear removes all elements
func (v *ArrayValue) Clear() {
	v.elements = make([]core.Value, 0)
}

// Serialize serializes the array and all its elements
func (v *ArrayValue) Serialize() (string, error) {
	result := fmt.Sprintf("[%s,%s,%d];", v.Name(), v.Type().String(), len(v.elements))
	for _, element := range v.elements {
		elemSer, err := element.Serialize()
		if err != nil {
			return "", err
		}
		result += elemSer
	}
	return result, nil
}

// ToXML converts to XML representation
func (v *ArrayValue) ToXML() (string, error) {
	type XMLArray struct {
		XMLName  xml.Name `xml:"array"`
		Name     string   `xml:"name,attr"`
		Type     string   `xml:"type,attr"`
		Count    int      `xml:"count,attr"`
		Elements []string `xml:"element"`
	}

	xmlArr := XMLArray{
		Name:     v.Name(),
		Type:     v.Type().TypeName(),
		Count:    len(v.elements),
		Elements: make([]string, 0),
	}

	for _, element := range v.elements {
		elemXML, err := element.ToXML()
		if err != nil {
			return "", err
		}
		xmlArr.Elements = append(xmlArr.Elements, elemXML)
	}

	data, err := xml.Marshal(xmlArr)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ToJSON converts to JSON representation
func (v *ArrayValue) ToJSON() (string, error) {
	type JSONArray struct {
		Name     string        `json:"name"`
		Type     string        `json:"type"`
		Elements []interface{} `json:"elements"`
	}

	jsonArr := JSONArray{
		Name:     v.Name(),
		Type:     "array",
		Elements: make([]interface{}, 0),
	}

	for _, element := range v.elements {
		elemJSON, err := element.ToJSON()
		if err != nil {
			return "", err
		}
		var elemData interface{}
		if err := json.Unmarshal([]byte(elemJSON), &elemData); err != nil {
			return "", err
		}
		jsonArr.Elements = append(jsonArr.Elements, elemData)
	}

	data, err := json.MarshalIndent(jsonArr, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Data returns a human-readable description as bytes
func (v *ArrayValue) Data() []byte {
	return []byte(fmt.Sprintf("Array(%d elements)", len(v.elements)))
}

// Size returns the size in bytes (for serialization)
func (v *ArrayValue) Size() int {
	size := 4 // count (4 bytes)
	for _, element := range v.elements {
		size += element.Size()
	}
	return size
}

// ToBinaryBytes serializes ArrayValue to binary format
//
// Binary format (little-endian):
// [type:1=15][name_len:4 LE][name:UTF-8][value_size:4 LE][count:4 LE][element1_bytes][element2_bytes]...
func (v *ArrayValue) ToBinaryBytes() ([]byte, error) {
	// Serialize all elements first to calculate total size
	serializedElements := make([][]byte, 0, len(v.elements))
	totalElementsSize := 0

	for _, element := range v.elements {
		elemBytes, err := element.ToBytes()
		if err != nil {
			return nil, fmt.Errorf("Failed to serialize element: %v", err)
		}
		serializedElements = append(serializedElements, elemBytes)
		totalElementsSize += len(elemBytes)
	}

	// value_size = count(4) + all element bytes
	valueSize := 4 + totalElementsSize

	name := v.Name()
	nameBytes := []byte(name)
	nameLen := uint32(len(nameBytes))

	// Calculate total size
	// type(1) + name_len(4) + name + value_size(4) + count(4) + elements
	totalSize := 1 + 4 + len(nameBytes) + 4 + 4 + totalElementsSize

	result := make([]byte, 0, totalSize)

	// Type (1 byte) - ArrayValue = 15
	result = append(result, byte(core.ArrayValue))

	// Name length (4 bytes, little-endian)
	result = append(result,
		byte(nameLen&0xFF),
		byte((nameLen>>8)&0xFF),
		byte((nameLen>>16)&0xFF),
		byte((nameLen>>24)&0xFF),
	)

	// Name (UTF-8 bytes)
	result = append(result, nameBytes...)

	// Value size (4 bytes, little-endian)
	valueSizeU32 := uint32(valueSize)
	result = append(result,
		byte(valueSizeU32&0xFF),
		byte((valueSizeU32>>8)&0xFF),
		byte((valueSizeU32>>16)&0xFF),
		byte((valueSizeU32>>24)&0xFF),
	)

	// Element count (4 bytes, little-endian)
	count := uint32(len(v.elements))
	result = append(result,
		byte(count&0xFF),
		byte((count>>8)&0xFF),
		byte((count>>16)&0xFF),
		byte((count>>24)&0xFF),
	)

	// Append all serialized elements
	for _, elemBytes := range serializedElements {
		result = append(result, elemBytes...)
	}

	return result, nil
}

// ToBytes implements the Value interface by delegating to ToBinaryBytes
func (v *ArrayValue) ToBytes() ([]byte, error) {
	return v.ToBinaryBytes()
}

// DeserializeArrayValue deserializes binary data into ArrayValue
//
// This function reads the binary format produced by C++, Rust, or Go ArrayValue.ToBytes()
// and reconstructs the ArrayValue with all its nested elements.
//
// Binary format:
// [type:1=15][name_len:4 LE][name:UTF-8][value_size:4 LE][count:4 LE][element1][element2]...
func DeserializeArrayValue(data []byte) (*ArrayValue, error) {
	if len(data) < 13 { // type(1) + name_len(4) + value_size(4) + count(4)
		return nil, fmt.Errorf("ArrayValue binary data too short: %d bytes", len(data))
	}

	offset := 0

	// Read type (1 byte)
	typeID := core.ValueType(data[offset])
	offset++

	if typeID != core.ArrayValue {
		return nil, fmt.Errorf("Expected ArrayValue type (15), got %d", typeID)
	}

	// Read name length (4 bytes, little-endian)
	nameLen := uint32(data[offset]) |
		(uint32(data[offset+1]) << 8) |
		(uint32(data[offset+2]) << 16) |
		(uint32(data[offset+3]) << 24)
	offset += 4

	// Read name
	if offset+int(nameLen) > len(data) {
		return nil, fmt.Errorf("Name length %d exceeds data bounds", nameLen)
	}
	name := string(data[offset : offset+int(nameLen)])
	offset += int(nameLen)

	// Read value size (4 bytes, little-endian)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("Insufficient data for value_size")
	}
	valueSize := uint32(data[offset]) |
		(uint32(data[offset+1]) << 8) |
		(uint32(data[offset+2]) << 16) |
		(uint32(data[offset+3]) << 24)
	offset += 4

	_ = valueSize // Not strictly needed for deserialization logic, but available for validation

	// Read element count (4 bytes, little-endian)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("Insufficient data for element count")
	}
	count := uint32(data[offset]) |
		(uint32(data[offset+1]) << 8) |
		(uint32(data[offset+2]) << 16) |
		(uint32(data[offset+3]) << 24)
	offset += 4

	// Create ArrayValue
	result := NewArrayValue(name)

	// Deserialize all elements
	for i := uint32(0); i < count; i++ {
		if offset >= len(data) {
			return nil, fmt.Errorf("Unexpected end of data while reading element %d/%d", i+1, count)
		}

		// Extract remaining data for element deserialization
		elementData := data[offset:]

		// Deserialize element using factory (requires value factory implementation)
		// TODO: Implement value factory for all types
		// For now, we'll need to implement type detection and routing
		element, bytesRead, err := deserializeValue(elementData)
		if err != nil {
			return nil, fmt.Errorf("Failed to deserialize element %d: %v", i, err)
		}

		result.Append(element)
		offset += bytesRead
	}

	return result, nil
}

// deserializeValue is a helper that deserializes a single value from binary data
// and returns the value along with the number of bytes consumed.
//
// This is a placeholder for the full Value factory implementation.
func deserializeValue(data []byte) (core.Value, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("Empty data for value deserialization")
	}

	// Read type ID
	typeID := core.ValueType(data[0])

	// Full factory pattern supporting all primitive value types
	switch typeID {
	case core.BoolValue:
		// Deserialize BoolValue (type 1)
		// Format: [type:1][name_len:4][name][value_size:4][value:1]
		if len(data) < 10 { // Minimum: type(1) + name_len(4) + value_size(4) + value(1)
			return nil, 0, fmt.Errorf("Insufficient data for BoolValue")
		}

		offset := 1
		nameLen := uint32(data[offset]) | (uint32(data[offset+1]) << 8) | (uint32(data[offset+2]) << 16) | (uint32(data[offset+3]) << 24)
		offset += 4

		if offset+int(nameLen)+4+1 > len(data) {
			return nil, 0, fmt.Errorf("Data too short for BoolValue")
		}

		name := string(data[offset : offset+int(nameLen)])
		offset += int(nameLen)

		// Skip value_size (4 bytes)
		offset += 4

		// Read bool value (1 byte)
		value := data[offset] != 0
		offset += 1

		return NewBoolValue(name, value), offset, nil

	case core.ShortValue:
		// Deserialize Int16Value (type 2)
		// Format: [type:1][name_len:4][name][value_size:4][value:2]
		if len(data) < 11 {
			return nil, 0, fmt.Errorf("Insufficient data for Int16Value")
		}

		offset := 1
		nameLen := uint32(data[offset]) | (uint32(data[offset+1]) << 8) | (uint32(data[offset+2]) << 16) | (uint32(data[offset+3]) << 24)
		offset += 4

		if offset+int(nameLen)+4+2 > len(data) {
			return nil, 0, fmt.Errorf("Data too short for Int16Value")
		}

		name := string(data[offset : offset+int(nameLen)])
		offset += int(nameLen)

		// Skip value_size (4 bytes)
		offset += 4

		// Read int16 value (2 bytes, little-endian)
		value := int16(data[offset]) | (int16(data[offset+1]) << 8)
		offset += 2

		return NewInt16Value(name, value), offset, nil

	case core.UShortValue:
		// Deserialize UInt16Value (type 3)
		// Format: [type:1][name_len:4][name][value_size:4][value:2]
		if len(data) < 11 {
			return nil, 0, fmt.Errorf("Insufficient data for UInt16Value")
		}

		offset := 1
		nameLen := uint32(data[offset]) | (uint32(data[offset+1]) << 8) | (uint32(data[offset+2]) << 16) | (uint32(data[offset+3]) << 24)
		offset += 4

		if offset+int(nameLen)+4+2 > len(data) {
			return nil, 0, fmt.Errorf("Data too short for UInt16Value")
		}

		name := string(data[offset : offset+int(nameLen)])
		offset += int(nameLen)

		// Skip value_size (4 bytes)
		offset += 4

		// Read uint16 value (2 bytes, little-endian)
		value := uint16(data[offset]) | (uint16(data[offset+1]) << 8)
		offset += 2

		return NewUInt16Value(name, value), offset, nil

	case core.IntValue:
		// Deserialize IntValue (type 3)
		// Format: [type:1][name_len:4][name][value_size:4][value:4]
		if len(data) < 13 { // Minimum: type(1) + name_len(4) + name(0) + value_size(4) + value(4)
			return nil, 0, fmt.Errorf("Insufficient data for IntValue")
		}

		offset := 1
		nameLen := uint32(data[offset]) | (uint32(data[offset+1]) << 8) | (uint32(data[offset+2]) << 16) | (uint32(data[offset+3]) << 24)
		offset += 4

		name := string(data[offset : offset+int(nameLen)])
		offset += int(nameLen)

		// Skip value_size (4 bytes)
		offset += 4

		// Read int value (4 bytes, little-endian)
		value := int32(data[offset]) | (int32(data[offset+1]) << 8) | (int32(data[offset+2]) << 16) | (int32(data[offset+3]) << 24)
		offset += 4

		return NewInt32Value(name, value), offset, nil

	case core.UIntValue:
		// Deserialize UInt32Value (type 5)
		if len(data) < 13 {
			return nil, 0, fmt.Errorf("Insufficient data for UInt32Value")
		}

		offset := 1
		nameLen := uint32(data[offset]) | (uint32(data[offset+1]) << 8) | (uint32(data[offset+2]) << 16) | (uint32(data[offset+3]) << 24)
		offset += 4

		name := string(data[offset : offset+int(nameLen)])
		offset += int(nameLen)

		offset += 4 // Skip value_size

		value := uint32(data[offset]) | (uint32(data[offset+1]) << 8) | (uint32(data[offset+2]) << 16) | (uint32(data[offset+3]) << 24)
		offset += 4

		return NewUInt32Value(name, value), offset, nil

	case core.LLongValue:
		// Deserialize Int64Value (type 8)
		if len(data) < 17 {
			return nil, 0, fmt.Errorf("Insufficient data for Int64Value")
		}

		offset := 1
		nameLen := uint32(data[offset]) | (uint32(data[offset+1]) << 8) | (uint32(data[offset+2]) << 16) | (uint32(data[offset+3]) << 24)
		offset += 4

		name := string(data[offset : offset+int(nameLen)])
		offset += int(nameLen)

		offset += 4 // Skip value_size

		value := int64(data[offset]) | (int64(data[offset+1]) << 8) | (int64(data[offset+2]) << 16) | (int64(data[offset+3]) << 24) |
			(int64(data[offset+4]) << 32) | (int64(data[offset+5]) << 40) | (int64(data[offset+6]) << 48) | (int64(data[offset+7]) << 56)
		offset += 8

		return NewInt64Value(name, value), offset, nil

	case core.ULLongValue:
		// Deserialize UInt64Value (type 9)
		if len(data) < 17 {
			return nil, 0, fmt.Errorf("Insufficient data for UInt64Value")
		}

		offset := 1
		nameLen := uint32(data[offset]) | (uint32(data[offset+1]) << 8) | (uint32(data[offset+2]) << 16) | (uint32(data[offset+3]) << 24)
		offset += 4

		name := string(data[offset : offset+int(nameLen)])
		offset += int(nameLen)

		offset += 4 // Skip value_size

		value := uint64(data[offset]) | (uint64(data[offset+1]) << 8) | (uint64(data[offset+2]) << 16) | (uint64(data[offset+3]) << 24) |
			(uint64(data[offset+4]) << 32) | (uint64(data[offset+5]) << 40) | (uint64(data[offset+6]) << 48) | (uint64(data[offset+7]) << 56)
		offset += 8

		return NewUInt64Value(name, value), offset, nil

	case core.FloatValue:
		// Deserialize Float32Value (type 10)
		if len(data) < 13 {
			return nil, 0, fmt.Errorf("Insufficient data for Float32Value")
		}

		offset := 1
		nameLen := uint32(data[offset]) | (uint32(data[offset+1]) << 8) | (uint32(data[offset+2]) << 16) | (uint32(data[offset+3]) << 24)
		offset += 4

		name := string(data[offset : offset+int(nameLen)])
		offset += int(nameLen)

		offset += 4 // Skip value_size

		bits := uint32(data[offset]) | (uint32(data[offset+1]) << 8) | (uint32(data[offset+2]) << 16) | (uint32(data[offset+3]) << 24)
		value := math.Float32frombits(bits)
		offset += 4

		return NewFloat32Value(name, value), offset, nil

	case core.DoubleValue:
		// Deserialize Float64Value (type 11)
		if len(data) < 17 {
			return nil, 0, fmt.Errorf("Insufficient data for Float64Value")
		}

		offset := 1
		nameLen := uint32(data[offset]) | (uint32(data[offset+1]) << 8) | (uint32(data[offset+2]) << 16) | (uint32(data[offset+3]) << 24)
		offset += 4

		name := string(data[offset : offset+int(nameLen)])
		offset += int(nameLen)

		offset += 4 // Skip value_size

		bits := uint64(data[offset]) | (uint64(data[offset+1]) << 8) | (uint64(data[offset+2]) << 16) | (uint64(data[offset+3]) << 24) |
			(uint64(data[offset+4]) << 32) | (uint64(data[offset+5]) << 40) | (uint64(data[offset+6]) << 48) | (uint64(data[offset+7]) << 56)
		value := math.Float64frombits(bits)
		offset += 8

		return NewFloat64Value(name, value), offset, nil

	case core.BytesValue:
		// Deserialize BytesValue (type 13) - matches C++ bytes_value position
		if len(data) < 13 {
			return nil, 0, fmt.Errorf("Insufficient data for BytesValue")
		}

		offset := 1
		nameLen := uint32(data[offset]) | (uint32(data[offset+1]) << 8) | (uint32(data[offset+2]) << 16) | (uint32(data[offset+3]) << 24)
		offset += 4

		name := string(data[offset : offset+int(nameLen)])
		offset += int(nameLen)

		valueSize := uint32(data[offset]) | (uint32(data[offset+1]) << 8) | (uint32(data[offset+2]) << 16) | (uint32(data[offset+3]) << 24)
		offset += 4

		value := make([]byte, valueSize)
		copy(value, data[offset:offset+int(valueSize)])
		offset += int(valueSize)

		return NewBytesValue(name, value), offset, nil

	case core.StringValue:
		// Deserialize StringValue (type 12) - matches C++ string_value position
		// Format: [type:1][name_len:4][name][value_size:4][string_bytes]
		if len(data) < 13 {
			return nil, 0, fmt.Errorf("Insufficient data for StringValue")
		}

		offset := 1
		nameLen := uint32(data[offset]) | (uint32(data[offset+1]) << 8) | (uint32(data[offset+2]) << 16) | (uint32(data[offset+3]) << 24)
		offset += 4

		name := string(data[offset : offset+int(nameLen)])
		offset += int(nameLen)

		valueSize := uint32(data[offset]) | (uint32(data[offset+1]) << 8) | (uint32(data[offset+2]) << 16) | (uint32(data[offset+3]) << 24)
		offset += 4

		strValue := string(data[offset : offset+int(valueSize)])
		offset += int(valueSize)

		return NewStringValue(name, strValue), offset, nil

	default:
		return nil, 0, fmt.Errorf("Unsupported value type for deserialization: %d", typeID)
	}
}
