package tests

import (
	"testing"

	"github.com/kcenon/go_container_system/container/core"
	"github.com/kcenon/go_container_system/container/values"
)

func TestBoolValue(t *testing.T) {
	bv := values.NewBoolValue("test_bool", true)

	if bv.Name() != "test_bool" {
		t.Errorf("Expected name 'test_bool', got '%s'", bv.Name())
	}

	if bv.Type() != core.BoolValue {
		t.Errorf("Expected type BoolValue, got %v", bv.Type())
	}

	val, err := bv.ToBool()
	if err != nil {
		t.Errorf("ToBool failed: %v", err)
	}
	if !val {
		t.Error("Expected true, got false")
	}

	str, err := bv.ToString()
	if err != nil {
		t.Errorf("ToString failed: %v", err)
	}
	if str != "true" {
		t.Errorf("Expected 'true', got '%s'", str)
	}
}

func TestNumericValues(t *testing.T) {
	// Int32
	iv := values.NewInt32Value("test_int", 42)
	val, err := iv.ToInt32()
	if err != nil {
		t.Errorf("ToInt32 failed: %v", err)
	}
	if val != 42 {
		t.Errorf("Expected 42, got %d", val)
	}

	// UInt32
	uv := values.NewUInt32Value("test_uint", 100)
	uval, err := uv.ToUInt32()
	if err != nil {
		t.Errorf("ToUInt32 failed: %v", err)
	}
	if uval != 100 {
		t.Errorf("Expected 100, got %d", uval)
	}

	// Float32
	fv := values.NewFloat32Value("test_float", 3.14)
	fval, err := fv.ToFloat32()
	if err != nil {
		t.Errorf("ToFloat32 failed: %v", err)
	}
	if fval < 3.13 || fval > 3.15 {
		t.Errorf("Expected ~3.14, got %f", fval)
	}

	// Float64
	dv := values.NewFloat64Value("test_double", 2.718281828)
	dval, err := dv.ToFloat64()
	if err != nil {
		t.Errorf("ToFloat64 failed: %v", err)
	}
	if dval < 2.71 || dval > 2.72 {
		t.Errorf("Expected ~2.718, got %f", dval)
	}
}

func TestStringValue(t *testing.T) {
	sv := values.NewStringValue("test_string", "Hello, World!")

	if sv.Name() != "test_string" {
		t.Errorf("Expected name 'test_string', got '%s'", sv.Name())
	}

	val, err := sv.ToString()
	if err != nil {
		t.Errorf("ToString failed: %v", err)
	}
	if val != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got '%s'", val)
	}
}

func TestBytesValue(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04}
	bv := values.NewBytesValue("test_bytes", data)

	if bv.Name() != "test_bytes" {
		t.Errorf("Expected name 'test_bytes', got '%s'", bv.Name())
	}

	val, err := bv.ToBytes()
	if err != nil {
		t.Errorf("ToBytes failed: %v", err)
	}
	if len(val) != 4 {
		t.Errorf("Expected 4 bytes, got %d", len(val))
	}
	for i, b := range data {
		if val[i] != b {
			t.Errorf("Byte mismatch at index %d: expected 0x%02x, got 0x%02x", i, b, val[i])
		}
	}
}

func TestContainerValue(t *testing.T) {
	// Create children
	child1 := values.NewStringValue("child1", "value1")
	child2 := values.NewInt32Value("child2", 123)

	// Create container
	cv := values.NewContainerValue("test_container", child1, child2)

	if cv.Name() != "test_container" {
		t.Errorf("Expected name 'test_container', got '%s'", cv.Name())
	}

	if cv.ChildCount() != 2 {
		t.Errorf("Expected 2 children, got %d", cv.ChildCount())
	}

	// Get child
	retrievedChild := cv.GetChild("child1", 0)
	if retrievedChild.Name() != "child1" {
		t.Errorf("Expected child name 'child1', got '%s'", retrievedChild.Name())
	}
}

func TestValueContainer(t *testing.T) {
	// Create container
	container := core.NewValueContainerFull(
		"source1", "sub1",
		"target1", "sub2",
		"test_message",
	)

	// Add values
	container.AddValue(values.NewStringValue("name", "John"))
	container.AddValue(values.NewInt32Value("age", 30))

	// Check header
	if container.SourceID() != "source1" {
		t.Errorf("Expected source 'source1', got '%s'", container.SourceID())
	}
	if container.MessageType() != "test_message" {
		t.Errorf("Expected message type 'test_message', got '%s'", container.MessageType())
	}

	// Get values
	name := container.GetValue("name", 0)
	if name.Name() != "name" {
		t.Errorf("Expected value name 'name', got '%s'", name.Name())
	}

	// Serialize
	serialized, err := container.Serialize()
	if err != nil {
		t.Errorf("Serialize failed: %v", err)
	}
	if len(serialized) == 0 {
		t.Error("Serialized string is empty")
	}

	// Swap header
	container.SwapHeader()
	if container.SourceID() != "target1" {
		t.Errorf("After swap, expected source 'target1', got '%s'", container.SourceID())
	}
	if container.TargetID() != "source1" {
		t.Errorf("After swap, expected target 'source1', got '%s'", container.TargetID())
	}
}

func TestValueContainerCopy(t *testing.T) {
	original := core.NewValueContainerWithType("test_message")
	original.AddValue(values.NewStringValue("data", "test"))

	// Copy with values
	copyWithValues := original.Copy(true)
	if len(copyWithValues.Values()) != 1 {
		t.Errorf("Expected 1 value in copy, got %d", len(copyWithValues.Values()))
	}

	// Copy without values
	copyWithoutValues := original.Copy(false)
	if len(copyWithoutValues.Values()) != 0 {
		t.Errorf("Expected 0 values in copy, got %d", len(copyWithoutValues.Values()))
	}
	if copyWithoutValues.MessageType() != "test_message" {
		t.Errorf("Expected message type 'test_message', got '%s'", copyWithoutValues.MessageType())
	}
}

func TestJSONSerialization(t *testing.T) {
	container := core.NewValueContainerWithType("test_message")
	container.AddValue(values.NewStringValue("name", "Alice"))
	container.AddValue(values.NewInt32Value("age", 25))

	json, err := container.ToJSON()
	if err != nil {
		t.Errorf("ToJSON failed: %v", err)
	}
	if len(json) == 0 {
		t.Error("JSON string is empty")
	}
	t.Logf("JSON: %s", json)
}

func TestXMLSerialization(t *testing.T) {
	container := core.NewValueContainerWithType("test_message")
	container.AddValue(values.NewStringValue("name", "Bob"))

	xml, err := container.ToXML()
	if err != nil {
		t.Errorf("ToXML failed: %v", err)
	}
	if len(xml) == 0 {
		t.Error("XML string is empty")
	}
	t.Logf("XML: %s", xml)
}

func TestMessagePackSerialization(t *testing.T) {
	// Create container with values
	container := core.NewValueContainerFull(
		"test_source", "sub1",
		"test_target", "sub2",
		"msgpack_test",
	)
	container.AddValue(values.NewStringValue("name", "Charlie"))
	container.AddValue(values.NewInt32Value("age", 35))

	// Serialize to MessagePack
	mpData, err := container.ToMessagePack()
	if err != nil {
		t.Errorf("ToMessagePack failed: %v", err)
	}
	if len(mpData) == 0 {
		t.Error("MessagePack data is empty")
	}
	t.Logf("MessagePack size: %d bytes", len(mpData))

	// Deserialize from MessagePack
	newContainer := core.NewValueContainer()
	err = newContainer.FromMessagePack(mpData)
	if err != nil {
		t.Errorf("FromMessagePack failed: %v", err)
	}

	// Verify header fields
	if newContainer.SourceID() != "test_source" {
		t.Errorf("Expected source 'test_source', got '%s'", newContainer.SourceID())
	}
	if newContainer.TargetID() != "test_target" {
		t.Errorf("Expected target 'test_target', got '%s'", newContainer.TargetID())
	}
	if newContainer.MessageType() != "msgpack_test" {
		t.Errorf("Expected message type 'msgpack_test', got '%s'", newContainer.MessageType())
	}
}
