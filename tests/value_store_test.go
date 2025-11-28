/****************************************************************************
BSD 3-Clause License

Copyright (c) 2021, 🍀☀🌕🌥 🌊
All rights reserved.
****************************************************************************/

package tests

import (
	"testing"

	"github.com/kcenon/go_container_system/container/core"
	"github.com/kcenon/go_container_system/container/values"
)

func TestValueStoreBasics(t *testing.T) {
	t.Run("CreateEmptyStore", func(t *testing.T) {
		store := core.NewValueStore()
		if store.Size() != 0 {
			t.Errorf("Expected size 0, got %d", store.Size())
		}
		if !store.Empty() {
			t.Error("Expected empty store")
		}
	})

	t.Run("AddAndGetValue", func(t *testing.T) {
		store := core.NewValueStore()
		value := values.NewInt32Value("count", 42)
		store.Add("count", value)

		if store.Size() != 1 {
			t.Errorf("Expected size 1, got %d", store.Size())
		}

		retrieved := store.Get("count")
		if retrieved == nil {
			t.Fatal("Expected value, got nil")
		}

		intVal, err := retrieved.ToInt32()
		if err != nil {
			t.Fatalf("Failed to convert to int: %v", err)
		}
		if intVal != 42 {
			t.Errorf("Expected 42, got %d", intVal)
		}
	})

	t.Run("Contains", func(t *testing.T) {
		store := core.NewValueStore()
		store.Add("name", values.NewStringValue("name", "test"))

		if !store.Contains("name") {
			t.Error("Expected to contain 'name'")
		}
		if store.Contains("missing") {
			t.Error("Should not contain 'missing'")
		}
	})

	t.Run("Remove", func(t *testing.T) {
		store := core.NewValueStore()
		store.Add("key1", values.NewInt32Value("key1", 1))
		store.Add("key2", values.NewInt32Value("key2", 2))

		if store.Size() != 2 {
			t.Errorf("Expected size 2, got %d", store.Size())
		}

		result := store.Remove("key1")
		if !result {
			t.Error("Expected Remove to return true")
		}
		if store.Size() != 1 {
			t.Errorf("Expected size 1, got %d", store.Size())
		}
		if store.Contains("key1") {
			t.Error("Should not contain 'key1' after removal")
		}

		result = store.Remove("nonexistent")
		if result {
			t.Error("Expected Remove to return false for non-existent key")
		}
	})

	t.Run("Clear", func(t *testing.T) {
		store := core.NewValueStore()
		store.Add("a", values.NewInt32Value("a", 1))
		store.Add("b", values.NewInt32Value("b", 2))
		store.Add("c", values.NewInt32Value("c", 3))

		if store.Size() != 3 {
			t.Errorf("Expected size 3, got %d", store.Size())
		}

		store.Clear()

		if store.Size() != 0 {
			t.Errorf("Expected size 0 after clear, got %d", store.Size())
		}
		if !store.Empty() {
			t.Error("Expected empty store after clear")
		}
	})

	t.Run("KeysAndValues", func(t *testing.T) {
		store := core.NewValueStore()
		store.Add("name", values.NewStringValue("name", "John"))
		store.Add("age", values.NewInt32Value("age", 30))

		keys := store.Keys()
		if len(keys) != 2 {
			t.Errorf("Expected 2 keys, got %d", len(keys))
		}

		vals := store.Values()
		if len(vals) != 2 {
			t.Errorf("Expected 2 values, got %d", len(vals))
		}
	})
}

func TestValueStoreTypes(t *testing.T) {
	t.Run("BoolValue", func(t *testing.T) {
		store := core.NewValueStore()
		store.Add("flag", values.NewBoolValue("flag", true))

		value := store.Get("flag")
		if value == nil {
			t.Fatal("Expected value, got nil")
		}

		boolVal, err := value.ToBool()
		if err != nil {
			t.Fatalf("Failed to convert to bool: %v", err)
		}
		if !boolVal {
			t.Error("Expected true, got false")
		}
	})

	t.Run("NumericValues", func(t *testing.T) {
		store := core.NewValueStore()

		store.Add("short", values.NewInt16Value("short", -100))
		store.Add("ushort", values.NewUInt16Value("ushort", 100))
		store.Add("int", values.NewInt32Value("int", -1000))
		store.Add("uint", values.NewUInt32Value("uint", 1000))
		store.Add("float", values.NewFloat32Value("float", 3.14))
		store.Add("double", values.NewFloat64Value("double", 3.14159265359))

		// Verify short
		shortVal := store.Get("short")
		if v, _ := shortVal.ToInt16(); v != -100 {
			t.Errorf("Expected -100, got %d", v)
		}

		// Verify uint
		uintVal := store.Get("uint")
		if v, _ := uintVal.ToUInt32(); v != 1000 {
			t.Errorf("Expected 1000, got %d", v)
		}

		// Verify double
		doubleVal := store.Get("double")
		if v, _ := doubleVal.ToFloat64(); v < 3.14 || v > 3.15 {
			t.Errorf("Expected ~3.14, got %f", v)
		}
	})

	t.Run("StringValue", func(t *testing.T) {
		store := core.NewValueStore()
		store.Add("message", values.NewStringValue("message", "Hello, World!"))

		value := store.Get("message")
		if value == nil {
			t.Fatal("Expected value, got nil")
		}

		strVal, err := value.ToString()
		if err != nil {
			t.Fatalf("Failed to convert to string: %v", err)
		}
		if strVal != "Hello, World!" {
			t.Errorf("Expected 'Hello, World!', got '%s'", strVal)
		}
	})

	t.Run("BytesValue", func(t *testing.T) {
		store := core.NewValueStore()
		data := []byte{0x01, 0x02, 0x03, 0x04}
		store.Add("data", values.NewBytesValue("data", data))

		value := store.Get("data")
		if value == nil {
			t.Fatal("Expected value, got nil")
		}

		// ToBytes returns full serialized format with headers
		// Check raw data via Data() method instead
		rawData := value.Data()
		if len(rawData) != 4 {
			t.Errorf("Expected 4 bytes in raw data, got %d", len(rawData))
		}
	})
}

func TestValueStoreSerialization(t *testing.T) {
	t.Run("JSONSerialization", func(t *testing.T) {
		store := core.NewValueStore()
		store.Add("name", values.NewStringValue("name", "Alice"))
		store.Add("age", values.NewInt32Value("age", 25))
		store.Add("active", values.NewBoolValue("active", true))

		jsonStr, err := store.Serialize()
		if err != nil {
			t.Fatalf("Failed to serialize: %v", err)
		}
		if jsonStr == "" {
			t.Error("Expected non-empty JSON string")
		}
	})

	t.Run("BinarySerialization", func(t *testing.T) {
		store := core.NewValueStore()
		store.Add("count", values.NewInt32Value("count", 42))
		store.Add("message", values.NewStringValue("message", "test"))

		binary, err := store.SerializeBinary()
		if err != nil {
			t.Fatalf("Failed to serialize binary: %v", err)
		}
		if len(binary) == 0 {
			t.Error("Expected non-empty binary data")
		}

		// Verify version byte
		if binary[0] != core.BinaryVersion {
			t.Errorf("Expected version %d, got %d", core.BinaryVersion, binary[0])
		}
	})
}

func TestValueStoreThreadSafety(t *testing.T) {
	t.Run("EnableDisableThreadSafety", func(t *testing.T) {
		store := core.NewValueStore()

		if store.IsThreadSafe() {
			t.Error("Expected thread safety to be disabled by default")
		}

		store.EnableThreadSafety()
		if !store.IsThreadSafe() {
			t.Error("Expected thread safety to be enabled")
		}

		store.DisableThreadSafety()
		if store.IsThreadSafe() {
			t.Error("Expected thread safety to be disabled")
		}
	})

	t.Run("ThreadSafeOperations", func(t *testing.T) {
		store := core.NewValueStore()
		store.EnableThreadSafety()

		store.Add("key", values.NewInt32Value("key", 1))
		value := store.Get("key")
		if value == nil {
			t.Fatal("Expected value, got nil")
		}

		intVal, _ := value.ToInt32()
		if intVal != 1 {
			t.Errorf("Expected 1, got %d", intVal)
		}

		store.Remove("key")
		if store.Contains("key") {
			t.Error("Should not contain 'key' after removal")
		}
	})
}

func TestValueStoreStatistics(t *testing.T) {
	t.Run("ReadWriteCounts", func(t *testing.T) {
		store := core.NewValueStore()

		if store.GetReadCount() != 0 {
			t.Errorf("Expected read count 0, got %d", store.GetReadCount())
		}
		if store.GetWriteCount() != 0 {
			t.Errorf("Expected write count 0, got %d", store.GetWriteCount())
		}

		store.Add("a", values.NewInt32Value("a", 1))
		if store.GetWriteCount() != 1 {
			t.Errorf("Expected write count 1, got %d", store.GetWriteCount())
		}

		store.Add("b", values.NewInt32Value("b", 2))
		if store.GetWriteCount() != 2 {
			t.Errorf("Expected write count 2, got %d", store.GetWriteCount())
		}

		store.Get("a")
		if store.GetReadCount() != 1 {
			t.Errorf("Expected read count 1, got %d", store.GetReadCount())
		}

		store.Get("b")
		if store.GetReadCount() != 2 {
			t.Errorf("Expected read count 2, got %d", store.GetReadCount())
		}
	})

	t.Run("SerializationCount", func(t *testing.T) {
		store := core.NewValueStore()
		store.Add("x", values.NewInt32Value("x", 10))

		if store.GetSerializationCount() != 0 {
			t.Errorf("Expected serialization count 0, got %d", store.GetSerializationCount())
		}

		store.Serialize()
		if store.GetSerializationCount() != 1 {
			t.Errorf("Expected serialization count 1, got %d", store.GetSerializationCount())
		}

		store.SerializeBinary()
		if store.GetSerializationCount() != 2 {
			t.Errorf("Expected serialization count 2, got %d", store.GetSerializationCount())
		}
	})

	t.Run("ResetStatistics", func(t *testing.T) {
		store := core.NewValueStore()
		store.Add("test", values.NewInt32Value("test", 1))
		store.Get("test")
		store.Serialize()

		store.ResetStatistics()

		if store.GetWriteCount() != 0 {
			t.Errorf("Expected write count 0 after reset, got %d", store.GetWriteCount())
		}
		if store.GetReadCount() != 0 {
			t.Errorf("Expected read count 0 after reset, got %d", store.GetReadCount())
		}
		if store.GetSerializationCount() != 0 {
			t.Errorf("Expected serialization count 0 after reset, got %d", store.GetSerializationCount())
		}
	})
}

func TestValueStoreRange(t *testing.T) {
	store := core.NewValueStore()
	store.Add("a", values.NewInt32Value("a", 1))
	store.Add("b", values.NewInt32Value("b", 2))
	store.Add("c", values.NewInt32Value("c", 3))

	count := 0
	store.Range(func(key string, value core.Value) bool {
		count++
		return true
	})

	if count != 3 {
		t.Errorf("Expected Range to visit 3 items, visited %d", count)
	}

	// Test early termination
	count = 0
	store.Range(func(key string, value core.Value) bool {
		count++
		return count < 2 // Stop after 2
	})

	if count != 2 {
		t.Errorf("Expected Range to stop after 2 items, visited %d", count)
	}
}

// Benchmarks
func BenchmarkValueStoreAdd(b *testing.B) {
	store := core.NewValueStore()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store.Add("key", values.NewInt32Value("key", int32(i)))
	}
}

func BenchmarkValueStoreGet(b *testing.B) {
	store := core.NewValueStore()
	store.Add("key", values.NewInt32Value("key", 42))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store.Get("key")
	}
}

func BenchmarkValueStoreSerialize(b *testing.B) {
	store := core.NewValueStore()
	store.Add("name", values.NewStringValue("name", "test"))
	store.Add("count", values.NewInt32Value("count", 42))
	store.Add("flag", values.NewBoolValue("flag", true))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store.Serialize()
	}
}

func BenchmarkValueStoreSerializeBinary(b *testing.B) {
	store := core.NewValueStore()
	store.Add("name", values.NewStringValue("name", "test"))
	store.Add("count", values.NewInt32Value("count", 42))
	store.Add("flag", values.NewBoolValue("flag", true))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store.SerializeBinary()
	}
}
