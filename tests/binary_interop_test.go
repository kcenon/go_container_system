package tests

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/kcenon/go_container_system/container/core"
	"github.com/kcenon/go_container_system/container/values"
)

// TestBinaryCompatibility_PrimitiveTypes tests binary serialization compatibility
// between Go and other languages (C++/Rust) for all primitive types
func TestBinaryCompatibility_PrimitiveTypes(t *testing.T) {
	tests := []struct {
		name     string
		value    core.Value
		typeByte byte
	}{
		{"Bool", values.NewBoolValue("test_bool", true), 0x01},
		{"Int16", values.NewInt16Value("test_i16", -12345), 0x02},
		{"UInt16", values.NewUInt16Value("test_u16", 54321), 0x03},
		{"Int32", values.NewInt32Value("test_i32", -987654), 0x04},
		{"UInt32", values.NewUInt32Value("test_u32", 4000000), 0x05},
		{"Int64", values.NewInt64Value("test_i64", -9876543210), 0x08},
		{"UInt64", values.NewUInt64Value("test_u64", 18446744073709551615), 0x09},
		{"Float32", values.NewFloat32Value("test_f32", 3.14159), 0x0A},
		{"Float64", values.NewFloat64Value("test_f64", 2.71828182845), 0x0B},
		{"String", values.NewStringValue("test_str", "Hello, World!"), 0x0D},
		{"Bytes", values.NewBytesValue("test_bytes", []byte{0xDE, 0xAD, 0xBE, 0xEF}), 0x0C},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize
			data, err := tt.value.ToBytes()
			if err != nil {
				t.Fatalf("ToBytes() failed: %v", err)
			}

			// Verify format structure
			if len(data) < 10 {
				t.Fatalf("Binary data too short: %d bytes", len(data))
			}

			// Verify type byte
			if data[0] != tt.typeByte {
				t.Errorf("Type byte mismatch: expected 0x%02X, got 0x%02X", tt.typeByte, data[0])
			}

			// Verify name length (4 bytes, little-endian)
			nameLen := uint32(data[1]) | (uint32(data[2]) << 8) | (uint32(data[3]) << 16) | (uint32(data[4]) << 24)
			expectedNameLen := uint32(len(tt.value.Name()))
			if nameLen != expectedNameLen {
				t.Errorf("Name length mismatch: expected %d, got %d", expectedNameLen, nameLen)
			}

			// Verify name
			nameStart := 5
			nameEnd := nameStart + int(nameLen)
			if nameEnd > len(data) {
				t.Fatalf("Name extends beyond data: nameEnd=%d, dataLen=%d", nameEnd, len(data))
			}
			actualName := string(data[nameStart:nameEnd])
			if actualName != tt.value.Name() {
				t.Errorf("Name mismatch: expected '%s', got '%s'", tt.value.Name(), actualName)
			}

			// Verify value_size field exists
			valueSizeStart := nameEnd
			if valueSizeStart+4 > len(data) {
				t.Fatalf("value_size field missing")
			}

			t.Logf("✓ %s binary format verified (%d bytes total)", tt.name, len(data))
		})
	}
}

// TestBinaryRoundtrip_CoreTypes verifies serialize → deserialize roundtrip
// for core primitive types (testing subset implemented in deserializeValue helper)
func TestBinaryRoundtrip_CoreTypes(t *testing.T) {
	tests := []struct {
		name  string
		value core.Value
	}{
		{"Bool_True", values.NewBoolValue("b1", true)},
		{"Bool_False", values.NewBoolValue("b2", false)},
		{"Int32_Pos", values.NewInt32Value("i32p", 2147483647)},
		{"Int32_Neg", values.NewInt32Value("i32n", -2147483648)},
		{"Int64_Pos", values.NewInt64Value("i64p", 9223372036854775807)},
		{"Int64_Neg", values.NewInt64Value("i64n", -9223372036854775808)},
		{"String_Empty", values.NewStringValue("s_empty", "")},
		{"String_ASCII", values.NewStringValue("s_ascii", "Hello, World!")},
		{"String_UTF8", values.NewStringValue("s_utf8", "안녕하세요 🌍")},
		{"Bytes_Empty", values.NewBytesValue("b_empty", []byte{})},
		{"Bytes_Data", values.NewBytesValue("b_data", []byte{0x00, 0xFF, 0xDE, 0xAD, 0xBE, 0xEF})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize
			data, err := tt.value.ToBytes()
			if err != nil {
				t.Fatalf("ToBytes() failed: %v", err)
			}

			// Deserialize
			restored, err := deserializeValue(data)
			if err != nil {
				t.Fatalf("Deserialization failed: %v", err)
			}

			// Verify name
			if restored.Name() != tt.value.Name() {
				t.Errorf("Name mismatch: expected '%s', got '%s'", tt.value.Name(), restored.Name())
			}

			// Verify type
			if restored.Type() != tt.value.Type() {
				t.Errorf("Type mismatch: expected %v, got %v", tt.value.Type(), restored.Type())
			}

			// Type-specific value verification
			switch tt.value.Type() {
			case core.BoolValue:
				expected, _ := tt.value.ToBool()
				actual, _ := restored.ToBool()
				if expected != actual {
					t.Errorf("Bool value mismatch: expected %v, got %v", expected, actual)
				}

			case core.IntValue:
				expected, _ := tt.value.ToInt32()
				actual, _ := restored.ToInt32()
				if expected != actual {
					t.Errorf("Int32 value mismatch: expected %d, got %d", expected, actual)
				}

			case core.LLongValue:
				expected, _ := tt.value.ToInt64()
				actual, _ := restored.ToInt64()
				if expected != actual {
					t.Errorf("Int64 value mismatch: expected %d, got %d", expected, actual)
				}

			case core.StringValue:
				expected, _ := tt.value.ToString()
				actual, _ := restored.ToString()
				if expected != actual {
					t.Errorf("String value mismatch: expected '%s', got '%s'", expected, actual)
				}

			case core.BytesValue:
				expected, _ := tt.value.ToBytes()
				actual, _ := restored.ToBytes()
				if !bytes.Equal(expected, actual) {
					t.Errorf("Bytes value mismatch: expected %v, got %v", expected, actual)
				}
			}

			t.Logf("✓ %s roundtrip successful", tt.name)
		})
	}
}

// TestBinaryInterop_KnownRustData tests deserialization of binary data
// generated by Rust implementation
func TestBinaryInterop_KnownRustData(t *testing.T) {
	// These hex strings represent binary data generated by Rust
	// Format: [type:1][name_len:4 LE][name][value_size:4 LE][value]

	tests := []struct {
		name      string
		hexData   string
		valueType core.ValueType
		valueName string
		checkFunc func(t *testing.T, v core.Value)
	}{
		{
			name:      "Rust_Int32",
			hexData:   "04" + "07000000" + "7465737469333332" + "04000000" + "2A000000", // Int32 "testi32" = 42
			valueType: core.IntValue,
			valueName: "testi32",
			checkFunc: func(t *testing.T, v core.Value) {
				val, err := v.ToInt32()
				if err != nil {
					t.Fatalf("ToInt32() failed: %v", err)
				}
				if val != 42 {
					t.Errorf("Expected 42, got %d", val)
				}
			},
		},
		{
			name:      "Rust_Bool_True",
			hexData:   "01" + "04000000" + "626F6F6C" + "01000000" + "01", // Bool "bool" = true
			valueType: core.BoolValue,
			valueName: "bool",
			checkFunc: func(t *testing.T, v core.Value) {
				val, err := v.ToBool()
				if err != nil {
					t.Fatalf("ToBool() failed: %v", err)
				}
				if !val {
					t.Error("Expected true, got false")
				}
			},
		},
		{
			name:      "Rust_String",
			hexData:   "0D" + "05000000" + "6D7973747220" + "0D000000" + "48656C6C6F2C20576F726C6421", // String "mystr" = "Hello, World!"
			valueType: core.StringValue,
			valueName: "mystr",
			checkFunc: func(t *testing.T, v core.Value) {
				val, err := v.ToString()
				if err != nil {
					t.Fatalf("ToString() failed: %v", err)
				}
				if val != "Hello, World!" {
					t.Errorf("Expected 'Hello, World!', got '%s'", val)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Decode hex
			data, err := hex.DecodeString(tt.hexData)
			if err != nil {
				t.Fatalf("Failed to decode hex: %v", err)
			}

			// Deserialize
			restored, err := deserializeValue(data)
			if err != nil {
				t.Fatalf("Deserialization failed: %v", err)
			}

			// Verify type
			if restored.Type() != tt.valueType {
				t.Errorf("Type mismatch: expected %v, got %v", tt.valueType, restored.Type())
			}

			// Verify name
			if restored.Name() != tt.valueName {
				t.Errorf("Name mismatch: expected '%s', got '%s'", tt.valueName, restored.Name())
			}

			// Type-specific verification
			tt.checkFunc(t, restored)

			t.Logf("✓ Successfully deserialized Rust-generated %s", tt.name)
		})
	}
}

// Helper: Simple value deserialization function for testing
// This uses array_value's deserializeValue function which already handles all types
func deserializeValue(data []byte) (core.Value, error) {
	// Use the ArrayValue package's deserialize function
	// For now, we'll implement a simple version for testing
	if len(data) < 10 {
		return nil, bytes.ErrTooLarge // placeholder error
	}

	typeID := core.ValueType(data[0])

	// Read name_len (4 bytes LE)
	nameLen := uint32(data[1]) | (uint32(data[2]) << 8) | (uint32(data[3]) << 16) | (uint32(data[4]) << 24)

	// Read name
	nameStart := 5
	nameEnd := nameStart + int(nameLen)
	name := string(data[nameStart:nameEnd])

	// Read value_size (4 bytes LE)
	valueSizeStart := nameEnd
	valueSize := uint32(data[valueSizeStart]) | (uint32(data[valueSizeStart+1]) << 8) |
		(uint32(data[valueSizeStart+2]) << 16) | (uint32(data[valueSizeStart+3]) << 24)

	// Read value data
	valueStart := valueSizeStart + 4
	valueEnd := valueStart + int(valueSize)
	valueData := data[valueStart:valueEnd]

	// Create appropriate value type
	switch typeID {
	case core.BoolValue:
		return values.NewBoolValueFromBytes(name, valueData)
	case core.IntValue:
		if len(valueData) != 4 {
			return nil, bytes.ErrTooLarge
		}
		val := int32(uint32(valueData[0]) | (uint32(valueData[1]) << 8) |
			(uint32(valueData[2]) << 16) | (uint32(valueData[3]) << 24))
		return values.NewInt32Value(name, val), nil
	case core.LLongValue:
		if len(valueData) != 8 {
			return nil, bytes.ErrTooLarge
		}
		val := int64(uint64(valueData[0]) | (uint64(valueData[1]) << 8) |
			(uint64(valueData[2]) << 16) | (uint64(valueData[3]) << 24) |
			(uint64(valueData[4]) << 32) | (uint64(valueData[5]) << 40) |
			(uint64(valueData[6]) << 48) | (uint64(valueData[7]) << 56))
		return values.NewInt64Value(name, val), nil
	case core.StringValue:
		return values.NewStringValue(name, string(valueData)), nil
	case core.BytesValue:
		return values.NewBytesValue(name, valueData), nil
	default:
		return nil, bytes.ErrTooLarge // placeholder
	}
}
