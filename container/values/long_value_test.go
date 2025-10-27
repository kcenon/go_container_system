/****************************************************************************
BSD 3-Clause License

Copyright (c) 2021, 🍀☀🌕🌥 🌊
All rights reserved.
****************************************************************************/

package values

import (
	"testing"
)

// =============================================================================
// LongValue (type 6) Tests - Signed 32-bit Range
// =============================================================================

func TestLongValue_AcceptsValidPositiveValue(t *testing.T) {
	lv, err := NewLongValue("test", 1000000)
	if err != nil {
		t.Fatalf("NewLongValue failed: %v", err)
	}
	if lv.Value() != 1000000 {
		t.Errorf("Expected 1000000, got %d", lv.Value())
	}
}

func TestLongValue_AcceptsValidNegativeValue(t *testing.T) {
	lv, err := NewLongValue("test", -1000000)
	if err != nil {
		t.Fatalf("NewLongValue failed: %v", err)
	}
	if lv.Value() != -1000000 {
		t.Errorf("Expected -1000000, got %d", lv.Value())
	}
}

func TestLongValue_AcceptsZero(t *testing.T) {
	lv, err := NewLongValue("test", 0)
	if err != nil {
		t.Fatalf("NewLongValue failed: %v", err)
	}
	if lv.Value() != 0 {
		t.Errorf("Expected 0, got %d", lv.Value())
	}
}

func TestLongValue_AcceptsInt32Max(t *testing.T) {
	lv, err := NewLongValue("test", int32Max)
	if err != nil {
		t.Fatalf("NewLongValue failed: %v", err)
	}
	if lv.Value() != int32Max {
		t.Errorf("Expected %d, got %d", int32Max, lv.Value())
	}
}

func TestLongValue_AcceptsInt32Min(t *testing.T) {
	lv, err := NewLongValue("test", int32Min)
	if err != nil {
		t.Fatalf("NewLongValue failed: %v", err)
	}
	if lv.Value() != int32Min {
		t.Errorf("Expected %d, got %d", int32Min, lv.Value())
	}
}

func TestLongValue_RejectsInt32MaxPlusOne(t *testing.T) {
	_, err := NewLongValue("test", int64(int32Max)+1)
	if err == nil {
		t.Error("Expected error for value exceeding int32 max, got nil")
	}
}

func TestLongValue_RejectsInt32MinMinusOne(t *testing.T) {
	_, err := NewLongValue("test", int64(int32Min)-1)
	if err == nil {
		t.Error("Expected error for value below int32 min, got nil")
	}
}

func TestLongValue_RejectsLargePositiveValue(t *testing.T) {
	_, err := NewLongValue("test", 5000000000)
	if err == nil {
		t.Error("Expected error for large positive value, got nil")
	}
}

func TestLongValue_RejectsLargeNegativeValue(t *testing.T) {
	_, err := NewLongValue("test", -5000000000)
	if err == nil {
		t.Error("Expected error for large negative value, got nil")
	}
}

// =============================================================================
// ULongValue (type 7) Tests - Unsigned 32-bit Range
// =============================================================================

func TestULongValue_AcceptsValidValue(t *testing.T) {
	ulv, err := NewULongValue("test", 1000000)
	if err != nil {
		t.Fatalf("NewULongValue failed: %v", err)
	}
	if ulv.Value() != 1000000 {
		t.Errorf("Expected 1000000, got %d", ulv.Value())
	}
}

func TestULongValue_AcceptsZero(t *testing.T) {
	ulv, err := NewULongValue("test", 0)
	if err != nil {
		t.Fatalf("NewULongValue failed: %v", err)
	}
	if ulv.Value() != 0 {
		t.Errorf("Expected 0, got %d", ulv.Value())
	}
}

func TestULongValue_AcceptsUInt32Max(t *testing.T) {
	ulv, err := NewULongValue("test", uint32Max)
	if err != nil {
		t.Fatalf("NewULongValue failed: %v", err)
	}
	if ulv.Value() != uint32Max {
		t.Errorf("Expected %d, got %d", uint32Max, ulv.Value())
	}
}

func TestULongValue_RejectsUInt32MaxPlusOne(t *testing.T) {
	_, err := NewULongValue("test", uint64(uint32Max)+1)
	if err == nil {
		t.Error("Expected error for value exceeding uint32 max, got nil")
	}
}

func TestULongValue_RejectsLargeValue(t *testing.T) {
	_, err := NewULongValue("test", 10000000000)
	if err == nil {
		t.Error("Expected error for large value, got nil")
	}
}

// =============================================================================
// Serialization Tests - Data Size Verification
// =============================================================================

func TestLongValue_SerializesAs4Bytes(t *testing.T) {
	lv, err := NewLongValue("test", 12345)
	if err != nil {
		t.Fatalf("NewLongValue failed: %v", err)
	}
	data := lv.BaseValue.Data()
	if len(data) != 4 {
		t.Errorf("Expected 4 bytes, got %d", len(data))
	}
}

func TestULongValue_SerializesAs4Bytes(t *testing.T) {
	ulv, err := NewULongValue("test", 12345)
	if err != nil {
		t.Fatalf("NewULongValue failed: %v", err)
	}
	data := ulv.BaseValue.Data()
	if len(data) != 4 {
		t.Errorf("Expected 4 bytes, got %d", len(data))
	}
}

// =============================================================================
// Type Conversion Tests
// =============================================================================

func TestLongValue_ConvertsToInt32(t *testing.T) {
	lv, err := NewLongValue("test", 12345)
	if err != nil {
		t.Fatalf("NewLongValue failed: %v", err)
	}
	val32, err := lv.ToInt32()
	if err != nil {
		t.Fatalf("ToInt32 failed: %v", err)
	}
	if val32 != 12345 {
		t.Errorf("Expected 12345, got %d", val32)
	}
}

func TestLongValue_ConvertsToInt64(t *testing.T) {
	lv, err := NewLongValue("test", 12345)
	if err != nil {
		t.Fatalf("NewLongValue failed: %v", err)
	}
	val64, err := lv.ToInt64()
	if err != nil {
		t.Fatalf("ToInt64 failed: %v", err)
	}
	if val64 != 12345 {
		t.Errorf("Expected 12345, got %d", val64)
	}
}

func TestULongValue_ConvertsToUInt32(t *testing.T) {
	ulv, err := NewULongValue("test", 12345)
	if err != nil {
		t.Fatalf("NewULongValue failed: %v", err)
	}
	val32, err := ulv.ToUInt32()
	if err != nil {
		t.Fatalf("ToUInt32 failed: %v", err)
	}
	if val32 != 12345 {
		t.Errorf("Expected 12345, got %d", val32)
	}
}

func TestULongValue_ConvertsToUInt64(t *testing.T) {
	ulv, err := NewULongValue("test", 12345)
	if err != nil {
		t.Fatalf("NewULongValue failed: %v", err)
	}
	val64, err := ulv.ToUInt64()
	if err != nil {
		t.Fatalf("ToUInt64 failed: %v", err)
	}
	if val64 != 12345 {
		t.Errorf("Expected 12345, got %d", val64)
	}
}

// =============================================================================
// Error Message Validation Tests
// =============================================================================

func TestLongValue_ErrorMessageIsDescriptive(t *testing.T) {
	_, err := NewLongValue("test", 5000000000)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	errMsg := err.Error()
	if !contains(errMsg, "LongValue") {
		t.Errorf("Error message should contain 'LongValue': %s", errMsg)
	}
	if !contains(errMsg, "32-bit") {
		t.Errorf("Error message should contain '32-bit': %s", errMsg)
	}
	if !contains(errMsg, "Int64Value") {
		t.Errorf("Error message should contain 'Int64Value': %s", errMsg)
	}
}

func TestULongValue_ErrorMessageIsDescriptive(t *testing.T) {
	_, err := NewULongValue("test", 10000000000)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	errMsg := err.Error()
	if !contains(errMsg, "ULongValue") {
		t.Errorf("Error message should contain 'ULongValue': %s", errMsg)
	}
	if !contains(errMsg, "32-bit") {
		t.Errorf("Error message should contain '32-bit': %s", errMsg)
	}
	if !contains(errMsg, "UInt64Value") {
		t.Errorf("Error message should contain 'UInt64Value': %s", errMsg)
	}
}

// =============================================================================
// Boundary Value Table Tests
// =============================================================================

func TestLongValue_BoundaryValues(t *testing.T) {
	tests := []struct {
		name    string
		value   int64
		wantErr bool
	}{
		{"INT32_MIN", int32Min, false},
		{"Negative million", -1000000, false},
		{"Negative one", -1, false},
		{"Zero", 0, false},
		{"Positive one", 1, false},
		{"Positive million", 1000000, false},
		{"INT32_MAX", int32Max, false},
		{"INT32_MIN - 1", int64(int32Min) - 1, true},
		{"INT32_MAX + 1", int64(int32Max) + 1, true},
		{"Large negative", -5000000000, true},
		{"Large positive", 5000000000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewLongValue("test", tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLongValue(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestULongValue_BoundaryValues(t *testing.T) {
	tests := []struct {
		name    string
		value   uint64
		wantErr bool
	}{
		{"Zero", 0, false},
		{"One", 1, false},
		{"One million", 1000000, false},
		{"UINT32_MAX", uint32Max, false},
		{"UINT32_MAX + 1", uint64(uint32Max) + 1, true},
		{"Large value", 10000000000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewULongValue("test", tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewULongValue(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && indexString(s, substr) >= 0)
}

func indexString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
