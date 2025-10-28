/****************************************************************************
BSD 3-Clause License

Copyright (c) 2021, 🍀☀🌕🌥 🌊
All rights reserved.
****************************************************************************/

package values

import (
	"github.com/kcenon/go_container_system/container/core"
)

// StringValue represents a string value
type StringValue struct {
	*core.BaseValue
	value string
}

// NewStringValue creates a new string value
func NewStringValue(name string, value string) *StringValue {
	return &StringValue{
		BaseValue: core.NewBaseValue(name, core.StringValue, []byte(value)),
		value:     value,
	}
}

// ToString returns the string value
func (v *StringValue) ToString() (string, error) {
	return v.value, nil
}

// ToBytes implements complete binary format with header
// Format: [type:1][name_len:4][name][value_size:4][string_bytes]
func (v *StringValue) ToBytes() ([]byte, error) {
	name := v.Name()
	nameBytes := []byte(name)
	nameLen := uint32(len(nameBytes))

	valueBytes := []byte(v.value)
	valueSize := uint32(len(valueBytes))

	// Total: type(1) + name_len(4) + name + value_size(4) + value
	totalSize := 1 + 4 + len(nameBytes) + 4 + len(valueBytes)
	result := make([]byte, 0, totalSize)

	// Type (1 byte)
	result = append(result, byte(core.StringValue))

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

	// String bytes (UTF-8)
	result = append(result, valueBytes...)

	return result, nil
}

// Value returns the underlying string value
func (v *StringValue) Value() string {
	return v.value
}
