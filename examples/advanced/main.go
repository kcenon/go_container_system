package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/kcenon/go_container_system/container/core"
	"github.com/kcenon/go_container_system/container/values"
)

func main() {
	fmt.Println("=== Advanced Go Container System Examples ===")

	// Example 1: MessagePack binary serialization
	messagePackExample()

	// Example 2: File I/O with different formats
	fileIOExample()

	// Example 3: Thread-safe concurrent operations
	threadSafeExample()

	// Example 4: Complex nested containers
	nestedContainersExample()

	fmt.Println("\n=== All advanced examples completed successfully! ===")
}

func messagePackExample() {
	fmt.Println("1. MessagePack Serialization:")

	container := core.NewValueContainerFull(
		"app_client", "v1.0",
		"api_server", "v2.0",
		"data_sync",
	)
	container.AddValue(values.NewStringValue("user_id", "user123"))
	container.AddValue(values.NewInt32Value("timestamp", 1234567890))
	container.AddValue(values.NewBoolValue("compressed", true))

	// Serialize to MessagePack
	msgpackData, err := container.ToMessagePack()
	if err != nil {
		log.Fatalf("MessagePack serialization failed: %v", err)
	}
	fmt.Printf("   MessagePack binary size: %d bytes\n", len(msgpackData))

	// Deserialize
	newContainer := core.NewValueContainer()
	if err := newContainer.FromMessagePack(msgpackData); err != nil {
		log.Fatalf("MessagePack deserialization failed: %v", err)
	}
	fmt.Printf("   Deserialized message type: %s\n", newContainer.MessageType())
	fmt.Println()
}

func fileIOExample() {
	fmt.Println("2. File I/O with Multiple Formats:")

	container := core.NewValueContainerWithType("file_storage")
	container.AddValue(values.NewStringValue("filename", "data.bin"))
	container.AddValue(values.NewInt64Value("size", 1024000))
	container.AddValue(values.NewFloat64Value("version", 2.5))

	// Save in different formats
	formats := []struct {
		name     string
		path     string
		saveFunc func(string) error
	}{
		{"MessagePack", "/tmp/container.msgpack", container.SaveToFileMessagePack},
		{"JSON", "/tmp/container.json", container.SaveToFileJSON},
		{"XML", "/tmp/container.xml", container.SaveToFileXML},
		{"String", "/tmp/container.txt", container.SaveToFile},
	}

	for _, format := range formats {
		if err := format.saveFunc(format.path); err != nil {
			log.Printf("   Failed to save %s: %v", format.name, err)
			continue
		}
		fmt.Printf("   ✓ Saved %s format to %s\n", format.name, format.path)
	}
	fmt.Println()
}

func threadSafeExample() {
	fmt.Println("3. Thread-Safe Concurrent Operations:")

	container := core.NewValueContainerWithType("concurrent_data")
	container.EnableThreadSafe()

	var wg sync.WaitGroup
	numWorkers := 10
	opsPerWorker := 100

	// Concurrent writers
	fmt.Printf("   Starting %d workers, each adding %d values...\n", numWorkers, opsPerWorker)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < opsPerWorker; j++ {
				container.AddValue(values.NewInt32Value(
					fmt.Sprintf("worker_%d_value_%d", workerID, j),
					int32(workerID*opsPerWorker+j),
				))
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("   ✓ Successfully added %d values concurrently\n", len(container.Values()))

	// Concurrent readers
	numReaders := 20
	fmt.Printf("   Starting %d concurrent readers...\n", numReaders)
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				_ = container.GetValue(fmt.Sprintf("worker_0_value_%d", j), 0)
			}
		}()
	}

	wg.Wait()
	fmt.Println("   ✓ Concurrent reads completed successfully")
	fmt.Println()
}

func nestedContainersExample() {
	fmt.Println("4. Complex Nested Containers:")

	// Create user profile
	userProfile := values.NewContainerValue("profile",
		values.NewStringValue("name", "Alice Johnson"),
		values.NewInt32Value("age", 28),
		values.NewStringValue("email", "alice@example.com"),
	)

	// Create user preferences
	preferences := values.NewContainerValue("preferences",
		values.NewBoolValue("notifications", true),
		values.NewStringValue("theme", "dark"),
		values.NewStringValue("language", "en"),
	)

	// Create user statistics
	statistics := values.NewContainerValue("statistics",
		values.NewInt64Value("login_count", 150),
		values.NewInt64Value("messages_sent", 1250),
		values.NewFloat64Value("average_session_time", 23.5),
	)

	// Create main user container
	userContainer := core.NewValueContainerFull(
		"mobile_app", "v3.2",
		"user_service", "v1.0",
		"user_data",
	)
	userContainer.AddValue(userProfile)
	userContainer.AddValue(preferences)
	userContainer.AddValue(statistics)

	fmt.Println("   Created nested container structure:")
	fmt.Printf("   - Main container: %s\n", userContainer.MessageType())
	fmt.Printf("   - Child containers: %d\n", len(userContainer.Values()))
	for _, child := range userContainer.Values() {
		if child.IsContainer() {
			fmt.Printf("     - %s (%d children)\n", child.Name(), child.ChildCount())
		}
	}

	// Serialize and measure sizes
	jsonSize, _ := userContainer.ToJSON()
	msgpackSize, _ := userContainer.ToMessagePack()
	fmt.Printf("   Serialization comparison:\n")
	fmt.Printf("   - JSON: %d bytes\n", len(jsonSize))
	fmt.Printf("   - MessagePack: %d bytes (%.1f%% smaller)\n",
		len(msgpackSize),
		100.0*(1.0-float64(len(msgpackSize))/float64(len(jsonSize))))
	fmt.Println()
}
