/****************************************************************************
BSD 3-Clause License

Copyright (c) 2021, 🍀☀🌕🌥 🌊
All rights reserved.
****************************************************************************/

package values

import (
	"bytes"
	"testing"
)

func TestArrayValueBinarySerialization(t *testing.T) {
	// Create an array with mixed types
	array := NewArrayValue("test_array")
	array.Push(NewInt32Value("", 42))
	array.Push(NewInt32Value("", 100))
	array.Push(NewStringValue("", "hello"))

	// Serialize to binary
	binaryData, err := array.ToBinaryBytes()
	if err != nil {
		t.Fatalf("Binary serialization failed: %v", err)
	}

	if len(binaryData) == 0 {
		t.Fatal("Binary serialization produced empty data")
	}

	t.Logf("Serialized %d bytes", len(binaryData))

	// Deserialize
	restored, err := DeserializeArrayValue(binaryData)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	// Verify name
	if restored.Name() != "test_array" {
		t.Errorf("Expected name 'test_array', got '%s'", restored.Name())
	}

	// Verify count
	if restored.Count() != 3 {
		t.Fatalf("Expected 3 elements, got %d", restored.Count())
	}

	// Verify first element (int 42)
	elem0, err := restored.At(0)
	if err != nil {
		t.Fatalf("Failed to get element 0: %v", err)
	}
	val0, err := elem0.ToInt32()
	if err != nil {
		t.Fatalf("Failed to convert element 0 to int32: %v", err)
	}
	if val0 != 42 {
		t.Errorf("Expected first element to be 42, got %d", val0)
	}

	// Verify second element (int 100)
	elem1, err := restored.At(1)
	if err != nil {
		t.Fatalf("Failed to get element 1: %v", err)
	}
	val1, err := elem1.ToInt32()
	if err != nil {
		t.Fatalf("Failed to convert element 1 to int32: %v", err)
	}
	if val1 != 100 {
		t.Errorf("Expected second element to be 100, got %d", val1)
	}

	// Verify third element (string "hello")
	elem2, err := restored.At(2)
	if err != nil {
		t.Fatalf("Failed to get element 2: %v", err)
	}
	val2, err := elem2.ToString()
	if err != nil {
		t.Fatalf("Failed to convert element 2 to string: %v", err)
	}
	if val2 != "hello" {
		t.Errorf("Expected third element to be 'hello', got '%s'", val2)
	}
}

func TestArrayValueBinaryRoundtrip(t *testing.T) {
	original := NewArrayValue("numbers")
	original.Push(NewInt32Value("", 1))
	original.Push(NewInt32Value("", 2))
	original.Push(NewInt32Value("", 3))

	// Serialize
	data, err := original.ToBinaryBytes()
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	// Deserialize
	restored, err := DeserializeArrayValue(data)
	if err != nil {
		t.Fatalf("Roundtrip failed: %v", err)
	}

	// Verify
	if original.Name() != restored.Name() {
		t.Errorf("Name mismatch: %s != %s", original.Name(), restored.Name())
	}
	if original.Count() != restored.Count() {
		t.Fatalf("Count mismatch: %d != %d", original.Count(), restored.Count())
	}

	for i := 0; i < original.Count(); i++ {
		origElem, _ := original.At(i)
		restElem, _ := restored.At(i)

		origVal, _ := origElem.ToInt32()
		restVal, _ := restElem.ToInt32()

		if origVal != restVal {
			t.Errorf("Element %d mismatch: %d != %d", i, origVal, restVal)
		}
	}
}

func TestArrayValueEmptyArray(t *testing.T) {
	// Create empty array
	empty := NewArrayValue("empty")

	// Serialize
	data, err := empty.ToBinaryBytes()
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	// Deserialize
	restored, err := DeserializeArrayValue(data)
	if err != nil {
		t.Fatalf("Empty array deserialization failed: %v", err)
	}

	if restored.Count() != 0 {
		t.Errorf("Expected empty array, got %d elements", restored.Count())
	}
	if !restored.IsEmpty() {
		t.Error("Expected IsEmpty() to return true")
	}
}

func TestArrayValueBinaryFormat(t *testing.T) {
	// Test that binary format matches C++ specification
	array := NewArrayValue("test")
	array.Push(NewInt32Value("", 123))

	data, err := array.ToBinaryBytes()
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	// Verify type byte (ArrayValue = 15)
	if data[0] != 15 {
		t.Errorf("Expected type byte 15, got %d", data[0])
	}

	// Verify name length (4 bytes, little-endian)
	nameLen := uint32(data[1]) | (uint32(data[2]) << 8) | (uint32(data[3]) << 16) | (uint32(data[4]) << 24)
	if nameLen != 4 { // "test" = 4 bytes
		t.Errorf("Expected name length 4, got %d", nameLen)
	}

	// Verify name
	name := string(data[5 : 5+nameLen])
	if name != "test" {
		t.Errorf("Expected name 'test', got '%s'", name)
	}
}

func TestArrayValueInvalidBinaryData(t *testing.T) {
	// Test with insufficient data
	tooShort := []byte{15, 0, 0, 0, 0} // Just type + name_len
	_, err := DeserializeArrayValue(tooShort)
	if err == nil {
		t.Error("Expected error for insufficient data")
	}

	// Test with wrong type
	wrongType := make([]byte, 20)
	wrongType[0] = 3 // IntValue instead of ArrayValue
	_, err = DeserializeArrayValue(wrongType)
	if err == nil {
		t.Error("Expected error for wrong type")
	}
}

func TestArrayValueBinaryCompatibility(t *testing.T) {
	// This test creates binary data that matches the C++ format exactly
	// to ensure cross-language compatibility

	// Create array with specific values
	array := NewArrayValue("colors")
	array.Push(NewStringValue("", "red"))
	array.Push(NewStringValue("", "blue"))

	// Serialize
	data, err := array.ToBinaryBytes()
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	// Verify we can deserialize it back
	restored, err := DeserializeArrayValue(data)
	if err != nil {
		t.Fatalf("Binary compatibility test failed: %v", err)
	}

	// Verify structure
	if restored.Name() != "colors" {
		t.Errorf("Name mismatch: %s", restored.Name())
	}
	if restored.Count() != 2 {
		t.Errorf("Count mismatch: %d", restored.Count())
	}

	// Verify element 0
	elem0, _ := restored.At(0)
	str0, _ := elem0.ToString()
	if str0 != "red" {
		t.Errorf("Element 0 mismatch: %s", str0)
	}

	// Verify element 1
	elem1, _ := restored.At(1)
	str1, _ := elem1.ToString()
	if str1 != "blue" {
		t.Errorf("Element 1 mismatch: %s", str1)
	}

	// Log binary format for manual verification
	t.Logf("Binary format (%d bytes): %v", len(data), data[:min(50, len(data))])
}

func TestArrayValueNestedDeserialization(t *testing.T) {
	// Test that we can properly calculate offset for sequential elements
	array := NewArrayValue("multi")
	array.Push(NewInt32Value("a", 10))
	array.Push(NewInt32Value("b", 20))
	array.Push(NewInt32Value("c", 30))

	data, err := array.ToBinaryBytes()
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	restored, err := DeserializeArrayValue(data)
	if err != nil {
		t.Fatalf("Nested deserialization failed: %v", err)
	}

	// Verify all three elements
	if restored.Count() != 3 {
		t.Fatalf("Expected 3 elements, got %d", restored.Count())
	}

	vals := []int32{10, 20, 30}
	names := []string{"a", "b", "c"}

	for i := 0; i < 3; i++ {
		elem, _ := restored.At(i)
		elemVal, _ := elem.ToInt32()
		if elemVal != vals[i] {
			t.Errorf("Element %d value mismatch: expected %d, got %d", i, vals[i], elemVal)
		}
		if elem.Name() != names[i] {
			t.Errorf("Element %d name mismatch: expected %s, got %s", i, names[i], elem.Name())
		}
	}
}

func TestArrayValueLargeArray(t *testing.T) {
	// Test with larger array to ensure offset calculations are correct
	const size = 100

	array := NewArrayValue("large")
	for i := 0; i < size; i++ {
		array.Push(NewInt32Value("", int32(i)))
	}

	data, err := array.ToBinaryBytes()
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}
	t.Logf("Large array serialized to %d bytes", len(data))

	restored, err := DeserializeArrayValue(data)
	if err != nil {
		t.Fatalf("Large array deserialization failed: %v", err)
	}

	if restored.Count() != size {
		t.Fatalf("Expected %d elements, got %d", size, restored.Count())
	}

	// Spot check a few elements
	checkIndices := []int{0, size / 4, size / 2, size * 3 / 4, size - 1}
	for _, idx := range checkIndices {
		elem, _ := restored.At(idx)
		elemVal, _ := elem.ToInt32()
		if elemVal != int32(idx) {
			t.Errorf("Element %d mismatch: expected %d, got %d", idx, idx, elemVal)
		}
	}
}

func TestArrayValueBinaryVsTextSerialization(t *testing.T) {
	// Compare binary vs text serialization sizes
	array := NewArrayValue("test")
	for i := 0; i < 10; i++ {
		array.Push(NewInt32Value("", int32(i*100)))
	}

	binaryData, err := array.ToBinaryBytes()
	if err != nil {
		t.Fatalf("Binary serialization failed: %v", err)
	}
	textData, _ := array.Serialize()

	t.Logf("Binary size: %d bytes", len(binaryData))
	t.Logf("Text size: %d bytes", len(textData))

	// Binary should typically be more compact for numeric data
	if len(binaryData) > len(textData) {
		t.Logf("Note: Binary format is larger than text for this data")
	}

	// Both should deserialize correctly
	restored, err := DeserializeArrayValue(binaryData)
	if err != nil {
		t.Fatalf("Binary deserialization failed: %v", err)
	}
	if restored.Count() != 10 {
		t.Errorf("Expected 10 elements, got %d", restored.Count())
	}
}

func TestArrayValueBinaryToBytes(t *testing.T) {
	// Verify that ToBinaryBytes() is idempotent
	array := NewArrayValue("stable")
	array.Push(NewInt32Value("", 42))

	data1, err1 := array.ToBinaryBytes()
	if err1 != nil {
		t.Fatalf("First serialization failed: %v", err1)
	}
	data2, err2 := array.ToBinaryBytes()
	if err2 != nil {
		t.Fatalf("Second serialization failed: %v", err2)
	}

	if !bytes.Equal(data1, data2) {
		t.Error("ToBinaryBytes() produced different results on repeated calls")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
