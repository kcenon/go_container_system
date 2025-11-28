package main

import (
	"fmt"
	"log"

	"github.com/kcenon/go_container_system/container/values"
)

func main() {
	fmt.Println("=== Go ArrayValue Quick Tests ===")

	// Test 1: Basic creation
	fmt.Println("Test 1: Basic creation")
	arr := values.NewArrayValue("numbers")
	if arr.Count() != 0 || !arr.IsEmpty() {
		log.Fatal("Empty array test failed")
	}
	fmt.Println("✓ Empty array created successfully")

	// Test 2: Append elements
	fmt.Println("Test 2: Append elements")
	arr.Append(values.NewInt32Value("", 10))
	arr.Append(values.NewInt32Value("", 20))
	arr.Append(values.NewInt32Value("", 30))

	if arr.Count() != 3 {
		log.Fatalf("Expected 3 elements, got %d", arr.Count())
	}

	elem0, err := arr.At(0)
	if err != nil {
		log.Fatal("Failed to get element 0")
	}
	elem1, err := arr.At(1)
	if err != nil {
		log.Fatal("Failed to get element 1")
	}
	elem2, err := arr.At(2)
	if err != nil {
		log.Fatal("Failed to get element 2")
	}

	// Just check they exist and are non-nil
	if elem0 == nil || elem1 == nil || elem2 == nil {
		log.Fatal("Element values incorrect")
	}
	fmt.Println("✓ Appending elements works")

	// Test 3: Heterogeneous array
	fmt.Println("Test 3: Heterogeneous array")
	mixed := values.NewArrayValue("mixed")
	mixed.Append(values.NewInt32Value("", 42))
	mixed.Append(values.NewStringValue("", "hello"))
	mixed.Append(values.NewFloat64Value("", 3.14))
	mixed.Append(values.NewBoolValue("", true))

	if mixed.Count() != 4 {
		log.Fatalf("Expected 4 elements, got %d", mixed.Count())
	}
	fmt.Println("✓ Mixed-type array works")

	// Test 4: Serialization
	fmt.Println("Test 4: Serialization")
	testArr := values.NewArrayValue("test")
	testArr.Append(values.NewInt32Value("val", 100))
	testArr.Append(values.NewInt32Value("val", 200))

	serialized, err := testArr.Serialize()
	if err != nil {
		log.Fatalf("Serialization failed: %v", err)
	}

	if len(serialized) == 0 {
		log.Fatal("Serialization returned empty string")
	}
	preview := serialized
	if len(serialized) > 50 {
		preview = serialized[:50]
	}
	fmt.Printf("  Serialized: %s\n", preview)
	fmt.Println("✓ Serialization works")

	// Test 5: Index out of range
	fmt.Println("Test 5: Index out of range")
	_, err = arr.At(100)
	if err == nil {
		log.Fatal("Expected error for out of range index")
	}
	fmt.Println("✓ Index bounds checking works")

	// Test 6: Clear
	fmt.Println("Test 6: Clear")
	arr.Clear()
	if arr.Count() != 0 || !arr.IsEmpty() {
		log.Fatal("Clear failed")
	}
	fmt.Println("✓ Clear works")

	// Test 7: ToJSON
	fmt.Println("Test 7: ToJSON")
	jsonArr := values.NewArrayValue("data")
	jsonArr.Append(values.NewInt32Value("num", 42))
	jsonArr.Append(values.NewStringValue("text", "test"))

	jsonStr, err := jsonArr.ToJSON()
	if err != nil {
		log.Fatalf("ToJSON failed: %v", err)
	}

	if len(jsonStr) == 0 {
		log.Fatal("ToJSON returned empty string")
	}
	jsonPreview := jsonStr
	if len(jsonStr) > 80 {
		jsonPreview = jsonStr[:80]
	}
	fmt.Printf("  JSON: %s\n", jsonPreview)
	fmt.Println("✓ ToJSON works")

	fmt.Println("=== All Tests Passed ===")
}
