package main

import (
	"fmt"
	"log"

	"github.com/kcenon/go_container_system/container/core"
	"github.com/kcenon/go_container_system/container/values"
)

func main() {
	fmt.Println("=== Go Container System - Basic Usage Example ===")

	// Example 1: Creating and using simple values
	fmt.Println("1. Creating Simple Values:")
	boolVal := values.NewBoolValue("enabled", true)
	intVal := values.NewInt32Value("count", 42)
	stringVal := values.NewStringValue("message", "Hello, Go Container System!")

	if val, err := boolVal.ToBool(); err == nil {
		fmt.Printf("   Bool: %v\n", val)
	}
	if val, err := intVal.ToInt32(); err == nil {
		fmt.Printf("   Int: %d\n", val)
	}
	if val, err := stringVal.ToString(); err == nil {
		fmt.Printf("   String: %s\n", val)
	}

	// Example 2: Creating a container value with children
	fmt.Println("\n2. Creating Container Value:")
	containerVal := values.NewContainerValue("user_data",
		values.NewStringValue("name", "Alice"),
		values.NewInt32Value("age", 30),
		values.NewStringValue("email", "alice@example.com"),
	)
	fmt.Printf("   Container '%s' has %d children\n", containerVal.Name(), containerVal.ChildCount())

	// Example 3: Creating a message container
	fmt.Println("\n3. Creating Message Container:")
	container := core.NewValueContainerFull(
		"client_app", "instance_1",
		"server_api", "v2",
		"user_registration",
	)

	container.AddValue(values.NewStringValue("username", "bob"))
	container.AddValue(values.NewStringValue("email", "bob@example.com"))
	container.AddValue(values.NewInt32Value("age", 25))

	fmt.Printf("   Source: %s/%s\n", container.SourceID(), container.SourceSubID())
	fmt.Printf("   Target: %s/%s\n", container.TargetID(), container.TargetSubID())
	fmt.Printf("   Message Type: %s\n", container.MessageType())
	fmt.Printf("   Value Count: %d\n", len(container.Values()))

	// Example 4: Serialization to string
	fmt.Println("\n4. Serialization:")
	serialized, err := container.Serialize()
	if err != nil {
		log.Fatalf("Serialization failed: %v", err)
	}
	fmt.Printf("   Serialized (first 100 chars): %s...\n", serialized[:min(100, len(serialized))])

	// Example 5: JSON conversion
	fmt.Println("\n5. JSON Conversion:")
	jsonStr, err := container.ToJSON()
	if err != nil {
		log.Fatalf("JSON conversion failed: %v", err)
	}
	fmt.Printf("   JSON (first 200 chars): %s...\n", jsonStr[:min(200, len(jsonStr))])

	// Example 6: XML conversion
	fmt.Println("\n6. XML Conversion:")
	xmlStr, err := container.ToXML()
	if err != nil {
		log.Fatalf("XML conversion failed: %v", err)
	}
	fmt.Printf("   XML (first 200 chars): %s...\n", xmlStr[:min(200, len(xmlStr))])

	// Example 7: Header swap
	fmt.Println("\n7. Header Swap:")
	fmt.Printf("   Before swap - Source: %s, Target: %s\n", container.SourceID(), container.TargetID())
	container.SwapHeader()
	fmt.Printf("   After swap  - Source: %s, Target: %s\n", container.SourceID(), container.TargetID())

	// Example 8: Getting specific values
	fmt.Println("\n8. Retrieving Values:")
	username := container.GetValue("username", 0)
	if username.Type() != core.NullValue {
		if val, err := username.ToString(); err == nil {
			fmt.Printf("   Username: %s\n", val)
		}
	}

	age := container.GetValue("age", 0)
	if age.Type() != core.NullValue {
		if val, err := age.ToInt32(); err == nil {
			fmt.Printf("   Age: %d\n", val)
		}
	}

	// Example 9: Container copy
	fmt.Println("\n9. Container Copy:")
	copiedContainer := container.Copy(true)
	fmt.Printf("   Original has %d values\n", len(container.Values()))
	fmt.Printf("   Copy has %d values\n", len(copiedContainer.Values()))

	headerOnlyCopy := container.Copy(false)
	fmt.Printf("   Header-only copy has %d values\n", len(headerOnlyCopy.Values()))

	fmt.Println("\n=== Example completed successfully! ===")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
