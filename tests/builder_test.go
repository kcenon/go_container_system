package tests

import (
	"testing"

	"github.com/kcenon/go_container_system/container/messaging"
	"github.com/kcenon/go_container_system/container/values"
)

func TestContainerBuilderBasic(t *testing.T) {
	container, err := messaging.NewContainerBuilder().
		WithSource("source1", "sub1").
		WithTarget("target1", "sub2").
		WithType("test_message").
		Build()

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if container.SourceID() != "source1" {
		t.Errorf("Expected source 'source1', got '%s'", container.SourceID())
	}
	if container.SourceSubID() != "sub1" {
		t.Errorf("Expected source sub ID 'sub1', got '%s'", container.SourceSubID())
	}
	if container.TargetID() != "target1" {
		t.Errorf("Expected target 'target1', got '%s'", container.TargetID())
	}
	if container.TargetSubID() != "sub2" {
		t.Errorf("Expected target sub ID 'sub2', got '%s'", container.TargetSubID())
	}
	if container.MessageType() != "test_message" {
		t.Errorf("Expected message type 'test_message', got '%s'", container.MessageType())
	}
}

func TestContainerBuilderWithValues(t *testing.T) {
	container, err := messaging.NewContainerBuilder().
		WithType("value_test").
		WithValues(
			values.NewStringValue("name", "John"),
			values.NewInt32Value("age", 30),
		).
		Build()

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	vals := container.Values()
	if len(vals) != 2 {
		t.Errorf("Expected 2 values, got %d", len(vals))
	}

	nameVal := container.GetValue("name", 0)
	if nameVal.Name() != "name" {
		t.Errorf("Expected value name 'name', got '%s'", nameVal.Name())
	}

	name, err := nameVal.ToString()
	if err != nil {
		t.Errorf("ToString failed: %v", err)
	}
	if name != "John" {
		t.Errorf("Expected 'John', got '%s'", name)
	}
}

func TestContainerBuilderThreadSafe(t *testing.T) {
	container, err := messaging.NewContainerBuilder().
		WithType("thread_safe_test").
		WithThreadSafe(true).
		Build()

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if !container.IsThreadSafe() {
		t.Error("Expected thread-safe mode to be enabled")
	}
}

func TestContainerBuilderThreadSafeDisabled(t *testing.T) {
	container, err := messaging.NewContainerBuilder().
		WithType("thread_safe_disabled_test").
		WithThreadSafe(false).
		Build()

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if container.IsThreadSafe() {
		t.Error("Expected thread-safe mode to be disabled")
	}
}

func TestContainerBuilderPartialConfiguration(t *testing.T) {
	t.Run("SourceOnly", func(t *testing.T) {
		container, err := messaging.NewContainerBuilder().
			WithSource("client", "1").
			Build()

		if err != nil {
			t.Fatalf("Build failed: %v", err)
		}

		if container.SourceID() != "client" {
			t.Errorf("Expected source 'client', got '%s'", container.SourceID())
		}
		if container.TargetID() != "" {
			t.Errorf("Expected empty target, got '%s'", container.TargetID())
		}
	})

	t.Run("TargetOnly", func(t *testing.T) {
		container, err := messaging.NewContainerBuilder().
			WithTarget("server", "main").
			Build()

		if err != nil {
			t.Fatalf("Build failed: %v", err)
		}

		if container.TargetID() != "server" {
			t.Errorf("Expected target 'server', got '%s'", container.TargetID())
		}
		if container.SourceID() != "" {
			t.Errorf("Expected empty source, got '%s'", container.SourceID())
		}
	})

	t.Run("TypeOnly", func(t *testing.T) {
		container, err := messaging.NewContainerBuilder().
			WithType("simple_message").
			Build()

		if err != nil {
			t.Fatalf("Build failed: %v", err)
		}

		if container.MessageType() != "simple_message" {
			t.Errorf("Expected message type 'simple_message', got '%s'", container.MessageType())
		}
	})
}

func TestContainerBuilderEmptyBuild(t *testing.T) {
	container, err := messaging.NewContainerBuilder().Build()

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if container.SourceID() != "" {
		t.Errorf("Expected empty source, got '%s'", container.SourceID())
	}
	if container.TargetID() != "" {
		t.Errorf("Expected empty target, got '%s'", container.TargetID())
	}
	if container.MessageType() != "" {
		t.Errorf("Expected empty message type, got '%s'", container.MessageType())
	}
}

func TestContainerBuilderChaining(t *testing.T) {
	// Test that all methods return the builder for chaining
	builder := messaging.NewContainerBuilder()

	result := builder.
		WithSource("src", "s1").
		WithTarget("tgt", "t1").
		WithType("msg").
		WithValues(values.NewBoolValue("flag", true)).
		WithThreadSafe(true)

	if result != builder {
		t.Error("Expected method chaining to return the same builder instance")
	}
}

func TestContainerBuilderMultipleValues(t *testing.T) {
	container, err := messaging.NewContainerBuilder().
		WithType("multi_value_test").
		WithValues(values.NewStringValue("first", "1")).
		WithValues(values.NewStringValue("second", "2")).
		WithValues(
			values.NewStringValue("third", "3"),
			values.NewStringValue("fourth", "4"),
		).
		Build()

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	vals := container.Values()
	if len(vals) != 4 {
		t.Errorf("Expected 4 values, got %d", len(vals))
	}
}
