/****************************************************************************
BSD 3-Clause License

Copyright (c) 2021, 🍀☀🌕🌥 🌊
All rights reserved.
****************************************************************************/

package values

import (
	"encoding/base64"

	"github.com/kcenon/go_container_system/container/core"
)

// BytesValue represents a binary data value
type BytesValue struct {
	*core.BaseValue
	value []byte
}

// NewBytesValue creates a new bytes value
func NewBytesValue(name string, value []byte) *BytesValue {
	// Make a copy to prevent external modification
	data := make([]byte, len(value))
	copy(data, value)
	return &BytesValue{
		BaseValue: core.NewBaseValue(name, core.BytesValue, data),
		value:     data,
	}
}

// ToBytes implements complete binary format with header
// Format: [type:1][name_len:4][name][value_size:4][bytes]
func (v *BytesValue) ToBytes() ([]byte, error) {
	name := v.Name()
	nameBytes := []byte(name)
	nameLen := uint32(len(nameBytes))
	valueSize := uint32(len(v.value))

	// Total: type(1) + name_len(4) + name + value_size(4) + value
	totalSize := 1 + 4 + len(nameBytes) + 4 + len(v.value)
	result := make([]byte, 0, totalSize)

	// Type (1 byte)
	result = append(result, byte(core.BytesValue))

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

	// Raw bytes
	result = append(result, v.value...)

	return result, nil
}

// ToString returns base64 encoded string
func (v *BytesValue) ToString() (string, error) {
	return base64.StdEncoding.EncodeToString(v.value), nil
}

// Value returns the underlying byte array
func (v *BytesValue) Value() []byte {
	result := make([]byte, len(v.value))
	copy(result, v.value)
	return result
}
