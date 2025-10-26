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

// ToBytes returns the byte representation
func (v *StringValue) ToBytes() ([]byte, error) {
	return []byte(v.value), nil
}

// Value returns the underlying string value
func (v *StringValue) Value() string {
	return v.value
}
