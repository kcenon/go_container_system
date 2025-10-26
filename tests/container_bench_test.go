package tests

import (
	"testing"

	"github.com/kcenon/go_container_system/container/core"
	"github.com/kcenon/go_container_system/container/values"
)

// Benchmark value creation
func BenchmarkBoolValueCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = values.NewBoolValue("test", true)
	}
}

func BenchmarkInt32ValueCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = values.NewInt32Value("test", 42)
	}
}

func BenchmarkStringValueCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = values.NewStringValue("test", "Hello, World!")
	}
}

// Benchmark container operations
func BenchmarkContainerAddValue(b *testing.B) {
	container := core.NewValueContainerWithType("bench_test")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		container.AddValue(values.NewInt32Value("data", int32(i)))
	}
}

func BenchmarkContainerAddValueThreadSafe(b *testing.B) {
	container := core.NewValueContainerWithType("bench_test")
	container.EnableThreadSafe()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		container.AddValue(values.NewInt32Value("data", int32(i)))
	}
}

func BenchmarkContainerGetValue(b *testing.B) {
	container := core.NewValueContainerWithType("bench_test")
	for i := 0; i < 100; i++ {
		container.AddValue(values.NewInt32Value("data", int32(i)))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = container.GetValue("data", i%100)
	}
}

func BenchmarkContainerGetValueThreadSafe(b *testing.B) {
	container := core.NewValueContainerWithType("bench_test")
	container.EnableThreadSafe()
	for i := 0; i < 100; i++ {
		container.AddValue(values.NewInt32Value("data", int32(i)))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = container.GetValue("data", i%100)
	}
}

// Benchmark serialization
func BenchmarkSerializeString(b *testing.B) {
	container := createBenchContainer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = container.Serialize()
	}
}

func BenchmarkSerializeJSON(b *testing.B) {
	container := createBenchContainer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = container.ToJSON()
	}
}

func BenchmarkSerializeXML(b *testing.B) {
	container := createBenchContainer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = container.ToXML()
	}
}

func BenchmarkSerializeMessagePack(b *testing.B) {
	container := createBenchContainer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = container.ToMessagePack()
	}
}

// Benchmark deserialization
func BenchmarkDeserializeMessagePack(b *testing.B) {
	container := createBenchContainer()
	data, _ := container.ToMessagePack()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newContainer := core.NewValueContainer()
		_ = newContainer.FromMessagePack(data)
	}
}

// Benchmark file I/O
func BenchmarkSaveToFileMessagePack(b *testing.B) {
	container := createBenchContainer()
	filePath := "/tmp/bench_container.msgpack"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = container.SaveToFileMessagePack(filePath)
	}
}

func BenchmarkLoadFromFileMessagePack(b *testing.B) {
	container := createBenchContainer()
	filePath := "/tmp/bench_container.msgpack"
	_ = container.SaveToFileMessagePack(filePath)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newContainer := core.NewValueContainer()
		_ = newContainer.LoadFromFileMessagePack(filePath)
	}
}

// Benchmark container copy
func BenchmarkContainerCopyWithValues(b *testing.B) {
	container := createBenchContainer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = container.Copy(true)
	}
}

func BenchmarkContainerCopyHeaderOnly(b *testing.B) {
	container := createBenchContainer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = container.Copy(false)
	}
}

// Benchmark parallel operations
func BenchmarkParallelAddValue(b *testing.B) {
	container := core.NewValueContainerWithType("parallel_bench")
	container.EnableThreadSafe()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			container.AddValue(values.NewInt32Value("data", int32(i)))
			i++
		}
	})
}

func BenchmarkParallelGetValue(b *testing.B) {
	container := core.NewValueContainerWithType("parallel_bench")
	container.EnableThreadSafe()
	for i := 0; i < 1000; i++ {
		container.AddValue(values.NewInt32Value("data", int32(i)))
	}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_ = container.GetValue("data", i%1000)
			i++
		}
	})
}

// Helper function
func createBenchContainer() *core.ValueContainer {
	container := core.NewValueContainerFull(
		"bench_source", "sub1",
		"bench_target", "sub2",
		"benchmark_message",
	)
	container.AddValue(values.NewStringValue("name", "Benchmark Test"))
	container.AddValue(values.NewInt32Value("count", 12345))
	container.AddValue(values.NewFloat64Value("pi", 3.14159265359))
	container.AddValue(values.NewBoolValue("enabled", true))
	return container
}
