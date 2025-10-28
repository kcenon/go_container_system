/****************************************************************************
BSD 3-Clause License

Copyright (c) 2021, 🍀☀🌕🌥 🌊
All rights reserved.
****************************************************************************/

package values

import (
	"errors"

	"github.com/kcenon/go_container_system/container/core"
)

// BoolValue represents a boolean value
type BoolValue struct {
	*core.BaseValue
	value bool
}

// NewBoolValue creates a new boolean value
func NewBoolValue(name string, value bool) *BoolValue {
	data := make([]byte, 1)
	if value {
		data[0] = 1
	}
	return &BoolValue{
		BaseValue: core.NewBaseValue(name, core.BoolValue, data),
		value:     value,
	}
}

// NewBoolValueFromBytes creates a boolean value from bytes
func NewBoolValueFromBytes(name string, data []byte) (*BoolValue, error) {
	if len(data) < 1 {
		return nil, errors.New("insufficient data for bool value")
	}
	value := data[0] != 0
	return NewBoolValue(name, value), nil
}

// ToBool returns the boolean value
func (v *BoolValue) ToBool() (bool, error) {
	return v.value, nil
}

// ToInt32 converts to int32 (0 or 1)
func (v *BoolValue) ToInt32() (int32, error) {
	if v.value {
		return 1, nil
	}
	return 0, nil
}

// ToInt64 converts to int64 (0 or 1)
func (v *BoolValue) ToInt64() (int64, error) {
	if v.value {
		return 1, nil
	}
	return 0, nil
}

// ToString converts to string ("true" or "false")
func (v *BoolValue) ToString() (string, error) {
	if v.value {
		return "true", nil
	}
	return "false", nil
}

// Value returns the underlying boolean value
func (v *BoolValue) Value() bool {
	return v.value
}

// ToBytes implements complete binary format with header
// Format: [type:1][name_len:4][name][value_size:4][value:1]
func (v *BoolValue) ToBytes() ([]byte, error) {
	name := v.Name()
	nameBytes := []byte(name)
	nameLen := uint32(len(nameBytes))
	valueSize := uint32(1) // bool = 1 byte

	// Total: type(1) + name_len(4) + name + value_size(4) + value(1)
	totalSize := 1 + 4 + len(nameBytes) + 4 + 1
	result := make([]byte, 0, totalSize)

	// Type (1 byte)
	result = append(result, byte(core.BoolValue))

	// Name length (4 bytes, little-endian)
	result = append(result,
		byte(nameLen&0xFF),
		byte((nameLen>>8)&0xFF),
		byte((nameLen>>16)&0xFF),
		byte((nameLen>>24)&0xFF),
	)

	// Name
	result = append(result, nameBytes...)

	// Value size (4 bytes, little-endian)
	result = append(result,
		byte(valueSize&0xFF),
		byte((valueSize>>8)&0xFF),
		byte((valueSize>>16)&0xFF),
		byte((valueSize>>24)&0xFF),
	)

	// Value (1 byte: 0=false, 1=true)
	if v.value {
		result = append(result, 1)
	} else {
		result = append(result, 0)
	}

	return result, nil
}
