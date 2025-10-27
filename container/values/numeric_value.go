/****************************************************************************
BSD 3-Clause License

Copyright (c) 2021, 🍀☀🌕🌥 🌊
All rights reserved.
****************************************************************************/

package values

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/kcenon/go_container_system/container/core"
)

// Int16Value represents a 16-bit signed integer
type Int16Value struct {
	*core.BaseValue
	value int16
}

// NewInt16Value creates a new int16 value
func NewInt16Value(name string, value int16) *Int16Value {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, uint16(value))
	return &Int16Value{
		BaseValue: core.NewBaseValue(name, core.ShortValue, data),
		value:     value,
	}
}

func (v *Int16Value) ToInt16() (int16, error) { return v.value, nil }
func (v *Int16Value) ToInt32() (int32, error) { return int32(v.value), nil }
func (v *Int16Value) ToInt64() (int64, error) { return int64(v.value), nil }
func (v *Int16Value) Value() int16             { return v.value }

// UInt16Value represents a 16-bit unsigned integer
type UInt16Value struct {
	*core.BaseValue
	value uint16
}

// NewUInt16Value creates a new uint16 value
func NewUInt16Value(name string, value uint16) *UInt16Value {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, value)
	return &UInt16Value{
		BaseValue: core.NewBaseValue(name, core.UShortValue, data),
		value:     value,
	}
}

func (v *UInt16Value) ToUInt16() (uint16, error) { return v.value, nil }
func (v *UInt16Value) ToUInt32() (uint32, error) { return uint32(v.value), nil }
func (v *UInt16Value) ToUInt64() (uint64, error) { return uint64(v.value), nil }
func (v *UInt16Value) Value() uint16              { return v.value }

// Int32Value represents a 32-bit signed integer
type Int32Value struct {
	*core.BaseValue
	value int32
}

// NewInt32Value creates a new int32 value
func NewInt32Value(name string, value int32) *Int32Value {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(value))
	return &Int32Value{
		BaseValue: core.NewBaseValue(name, core.IntValue, data),
		value:     value,
	}
}

func (v *Int32Value) ToInt32() (int32, error) { return v.value, nil }
func (v *Int32Value) ToInt64() (int64, error) { return int64(v.value), nil }
func (v *Int32Value) Value() int32             { return v.value }

// UInt32Value represents a 32-bit unsigned integer
type UInt32Value struct {
	*core.BaseValue
	value uint32
}

// NewUInt32Value creates a new uint32 value
func NewUInt32Value(name string, value uint32) *UInt32Value {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, value)
	return &UInt32Value{
		BaseValue: core.NewBaseValue(name, core.UIntValue, data),
		value:     value,
	}
}

func (v *UInt32Value) ToUInt32() (uint32, error) { return v.value, nil }
func (v *UInt32Value) ToUInt64() (uint64, error) { return uint64(v.value), nil }
func (v *UInt32Value) Value() uint32              { return v.value }

// Int64Value represents a 64-bit signed integer
type Int64Value struct {
	*core.BaseValue
	value int64
}

// NewInt64Value creates a new int64 value
func NewInt64Value(name string, value int64) *Int64Value {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, uint64(value))
	return &Int64Value{
		BaseValue: core.NewBaseValue(name, core.LLongValue, data),
		value:     value,
	}
}

func (v *Int64Value) ToInt64() (int64, error) { return v.value, nil }
func (v *Int64Value) Value() int64             { return v.value }

// UInt64Value represents a 64-bit unsigned integer
type UInt64Value struct {
	*core.BaseValue
	value uint64
}

// NewUInt64Value creates a new uint64 value
func NewUInt64Value(name string, value uint64) *UInt64Value {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, value)
	return &UInt64Value{
		BaseValue: core.NewBaseValue(name, core.ULLongValue, data),
		value:     value,
	}
}

func (v *UInt64Value) ToUInt64() (uint64, error) { return v.value, nil }
func (v *UInt64Value) Value() uint64              { return v.value }

// Float32Value represents a 32-bit floating point
type Float32Value struct {
	*core.BaseValue
	value float32
}

// NewFloat32Value creates a new float32 value
func NewFloat32Value(name string, value float32) *Float32Value {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, math.Float32bits(value))
	return &Float32Value{
		BaseValue: core.NewBaseValue(name, core.FloatValue, data),
		value:     value,
	}
}

func (v *Float32Value) ToFloat32() (float32, error) { return v.value, nil }
func (v *Float32Value) ToFloat64() (float64, error) { return float64(v.value), nil }
func (v *Float32Value) Value() float32               { return v.value }

// Float64Value represents a 64-bit floating point
type Float64Value struct {
	*core.BaseValue
	value float64
}

// NewFloat64Value creates a new float64 value
func NewFloat64Value(name string, value float64) *Float64Value {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, math.Float64bits(value))
	return &Float64Value{
		BaseValue: core.NewBaseValue(name, core.DoubleValue, data),
		value:     value,
	}
}

func (v *Float64Value) ToFloat64() (float64, error) { return v.value, nil }
func (v *Float64Value) Value() float64               { return v.value }

// =============================================================================
// Long/ULong Types (32-bit range policy)
// =============================================================================

// LongValue represents a 32-bit signed integer (type 6).
// Policy: Enforces 32-bit range [-2^31, 2^31-1].
// Values exceeding this range should use Int64Value.
// Always serializes as 4 bytes (int32) regardless of platform.
type LongValue struct {
	*core.BaseValue
	value int32
}

const (
	int32Min = -2147483648
	int32Max = 2147483647
)

// NewLongValue creates a new long value with range checking
func NewLongValue(name string, value int64) (*LongValue, error) {
	// Enforce strict 32-bit range policy
	if value < int32Min || value > int32Max {
		return nil, fmt.Errorf(
			"LongValue: value %d exceeds 32-bit range [%d, %d]. "+
				"Use Int64Value for 64-bit values",
			value, int32Min, int32Max,
		)
	}

	val32 := int32(value)
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(val32))

	return &LongValue{
		BaseValue: core.NewBaseValue(name, core.LongValue, data),
		value:     val32,
	}, nil
}

func (v *LongValue) ToInt32() (int32, error) { return v.value, nil }
func (v *LongValue) ToInt64() (int64, error) { return int64(v.value), nil }
func (v *LongValue) Value() int32             { return v.value }

// ULongValue represents a 32-bit unsigned integer (type 7).
// Policy: Enforces 32-bit range [0, 2^32-1].
// Values exceeding this range should use UInt64Value.
// Always serializes as 4 bytes (uint32) regardless of platform.
type ULongValue struct {
	*core.BaseValue
	value uint32
}

const uint32Max = 4294967295

// NewULongValue creates a new ulong value with range checking
func NewULongValue(name string, value uint64) (*ULongValue, error) {
	// Enforce strict 32-bit range policy
	if value > uint32Max {
		return nil, fmt.Errorf(
			"ULongValue: value %d exceeds 32-bit range [0, %d]. "+
				"Use UInt64Value for 64-bit values",
			value, uint32Max,
		)
	}

	val32 := uint32(value)
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, val32)

	return &ULongValue{
		BaseValue: core.NewBaseValue(name, core.ULongValue, data),
		value:     val32,
	}, nil
}

func (v *ULongValue) ToUInt32() (uint32, error) { return v.value, nil }
func (v *ULongValue) ToUInt64() (uint64, error) { return uint64(v.value), nil }
func (v *ULongValue) Value() uint32              { return v.value }
