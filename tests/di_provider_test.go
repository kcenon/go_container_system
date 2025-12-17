package tests

import (
	"testing"

	"github.com/kcenon/go_container_system/container/core"
	"github.com/kcenon/go_container_system/container/di"
	"github.com/kcenon/go_container_system/container/messaging"
)

func TestNewContainerFactory(t *testing.T) {
	factory := di.NewContainerFactory()

	if factory == nil {
		t.Fatal("NewContainerFactory returned nil")
	}
}

func TestContainerFactoryNewContainer(t *testing.T) {
	factory := di.NewContainerFactory()
	container := factory.NewContainer()

	if container == nil {
		t.Fatal("NewContainer returned nil")
	}

	if container.SourceID() != "" {
		t.Errorf("Expected empty source ID, got '%s'", container.SourceID())
	}
	if container.TargetID() != "" {
		t.Errorf("Expected empty target ID, got '%s'", container.TargetID())
	}
	if container.MessageType() != "" {
		t.Errorf("Expected empty message type, got '%s'", container.MessageType())
	}
	if container.Version() != "1.0.0.0" {
		t.Errorf("Expected version '1.0.0.0', got '%s'", container.Version())
	}
}

func TestContainerFactoryNewContainerWithType(t *testing.T) {
	factory := di.NewContainerFactory()
	container := factory.NewContainerWithType("test_message")

	if container == nil {
		t.Fatal("NewContainerWithType returned nil")
	}

	if container.MessageType() != "test_message" {
		t.Errorf("Expected message type 'test_message', got '%s'", container.MessageType())
	}
}

func TestContainerFactoryNewContainerWithTarget(t *testing.T) {
	factory := di.NewContainerFactory()
	container := factory.NewContainerWithTarget("server", "main", "request")

	if container == nil {
		t.Fatal("NewContainerWithTarget returned nil")
	}

	if container.TargetID() != "server" {
		t.Errorf("Expected target ID 'server', got '%s'", container.TargetID())
	}
	if container.TargetSubID() != "main" {
		t.Errorf("Expected target sub ID 'main', got '%s'", container.TargetSubID())
	}
	if container.MessageType() != "request" {
		t.Errorf("Expected message type 'request', got '%s'", container.MessageType())
	}
}

func TestContainerFactoryNewContainerFull(t *testing.T) {
	factory := di.NewContainerFactory()
	container := factory.NewContainerFull("client", "1", "server", "main", "request")

	if container == nil {
		t.Fatal("NewContainerFull returned nil")
	}

	if container.SourceID() != "client" {
		t.Errorf("Expected source ID 'client', got '%s'", container.SourceID())
	}
	if container.SourceSubID() != "1" {
		t.Errorf("Expected source sub ID '1', got '%s'", container.SourceSubID())
	}
	if container.TargetID() != "server" {
		t.Errorf("Expected target ID 'server', got '%s'", container.TargetID())
	}
	if container.TargetSubID() != "main" {
		t.Errorf("Expected target sub ID 'main', got '%s'", container.TargetSubID())
	}
	if container.MessageType() != "request" {
		t.Errorf("Expected message type 'request', got '%s'", container.MessageType())
	}
}

func TestContainerFactoryNewBuilder(t *testing.T) {
	factory := di.NewContainerFactory()
	builder := factory.NewBuilder()

	if builder == nil {
		t.Fatal("NewBuilder returned nil")
	}

	container, err := builder.
		WithSource("client", "1").
		WithTarget("server", "main").
		WithType("request").
		Build()

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if container.SourceID() != "client" {
		t.Errorf("Expected source ID 'client', got '%s'", container.SourceID())
	}
	if container.TargetID() != "server" {
		t.Errorf("Expected target ID 'server', got '%s'", container.TargetID())
	}
}

func TestContainerFactoryInterface(t *testing.T) {
	var factory di.ContainerFactory = di.NewContainerFactory()

	if factory == nil {
		t.Fatal("Factory should implement ContainerFactory interface")
	}

	container := factory.NewContainer()
	if container == nil {
		t.Fatal("NewContainer returned nil")
	}
}

func TestContainerFactoryMocking(t *testing.T) {
	var factory di.ContainerFactory = &mockContainerFactory{}

	container := factory.NewContainer()
	if container == nil {
		t.Fatal("Mock factory should return a container")
	}

	if container.MessageType() != "mocked" {
		t.Errorf("Expected mocked message type, got '%s'", container.MessageType())
	}
}

type mockContainerFactory struct{}

func (m *mockContainerFactory) NewContainer() *core.ValueContainer {
	return core.NewValueContainerWithType("mocked")
}

func (m *mockContainerFactory) NewContainerWithType(messageType string) *core.ValueContainer {
	return core.NewValueContainerWithType("mocked_" + messageType)
}

func (m *mockContainerFactory) NewContainerWithTarget(targetID, targetSubID, messageType string) *core.ValueContainer {
	return core.NewValueContainerWithTarget(targetID, targetSubID, "mocked_"+messageType)
}

func (m *mockContainerFactory) NewContainerFull(sourceID, sourceSubID, targetID, targetSubID, messageType string) *core.ValueContainer {
	return core.NewValueContainerFull(sourceID, sourceSubID, targetID, targetSubID, "mocked_"+messageType)
}

func (m *mockContainerFactory) NewBuilder() *messaging.ContainerBuilder {
	return messaging.NewContainerBuilder()
}
