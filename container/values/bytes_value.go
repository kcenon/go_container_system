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

// ToBytes returns the byte value
func (v *BytesValue) ToBytes() ([]byte, error) {
	// Return a copy to prevent external modification
	result := make([]byte, len(v.value))
	copy(result, v.value)
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
