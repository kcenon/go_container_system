/****************************************************************************
BSD 3-Clause License

Copyright (c) 2021, 🍀☀🌕🌥 🌊
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
****************************************************************************/

package core

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
)

// Value represents the interface for all value types in the container system
type Value interface {
	// Basic accessors
	Name() string
	Type() ValueType
	Data() []byte
	Size() int

	// Type checking
	IsNull() bool
	IsBytes() bool
	IsBoolean() bool
	IsNumeric() bool
	IsString() bool
	IsContainer() bool

	// Type conversions
	ToBool() (bool, error)
	ToInt16() (int16, error)
	ToUInt16() (uint16, error)
	ToInt32() (int32, error)
	ToUInt32() (uint32, error)
	ToInt64() (int64, error)
	ToUInt64() (uint64, error)
	ToFloat32() (float32, error)
	ToFloat64() (float64, error)
	ToString() (string, error)
	ToBytes() ([]byte, error)

	// Serialization
	Serialize() (string, error)
	ToXML() (string, error)
	ToJSON() (string, error)

	// For container values
	Children() []Value
	ChildCount() int
	GetChild(name string, index int) Value
	AddChild(child Value) error
	RemoveChild(name string) error
}

// BaseValue provides the base implementation for all value types
type BaseValue struct {
	name   string
	vtype  ValueType
	data   []byte
	parent Value
	units  []Value
}

// NewBaseValue creates a new base value
func NewBaseValue(name string, vtype ValueType, data []byte) *BaseValue {
	return &BaseValue{
		name:  name,
		vtype: vtype,
		data:  data,
		units: make([]Value, 0),
	}
}

// Name returns the name of the value
func (v *BaseValue) Name() string {
	return v.name
}

// Type returns the type of the value
func (v *BaseValue) Type() ValueType {
	return v.vtype
}

// Data returns the raw data
func (v *BaseValue) Data() []byte {
	return v.data
}

// Size returns the size of the data
func (v *BaseValue) Size() int {
	return len(v.data)
}

// IsNull checks if the value is null
func (v *BaseValue) IsNull() bool {
	return v.vtype == NullValue
}

// IsBytes checks if the value is bytes
func (v *BaseValue) IsBytes() bool {
	return v.vtype == BytesValue
}

// IsBoolean checks if the value is boolean
func (v *BaseValue) IsBoolean() bool {
	return v.vtype == BoolValue
}

// IsNumeric checks if the value is numeric
func (v *BaseValue) IsNumeric() bool {
	return v.vtype >= ShortValue && v.vtype <= DoubleValue
}

// IsString checks if the value is string
func (v *BaseValue) IsString() bool {
	return v.vtype == StringValue
}

// IsContainer checks if the value is container
func (v *BaseValue) IsContainer() bool {
	return v.vtype == ContainerValue
}

// ToBool converts to boolean (default implementation)
func (v *BaseValue) ToBool() (bool, error) {
	if v.IsNull() {
		return false, errors.New("cannot convert null_value to bool")
	}
	return false, errors.New("type conversion not supported")
}

// ToInt16 converts to int16 (default implementation)
func (v *BaseValue) ToInt16() (int16, error) {
	if v.IsNull() {
		return 0, errors.New("cannot convert null_value to int16")
	}
	return 0, errors.New("type conversion not supported")
}

// ToUInt16 converts to uint16 (default implementation)
func (v *BaseValue) ToUInt16() (uint16, error) {
	if v.IsNull() {
		return 0, errors.New("cannot convert null_value to uint16")
	}
	return 0, errors.New("type conversion not supported")
}

// ToInt32 converts to int32 (default implementation)
func (v *BaseValue) ToInt32() (int32, error) {
	if v.IsNull() {
		return 0, errors.New("cannot convert null_value to int32")
	}
	return 0, errors.New("type conversion not supported")
}

// ToUInt32 converts to uint32 (default implementation)
func (v *BaseValue) ToUInt32() (uint32, error) {
	if v.IsNull() {
		return 0, errors.New("cannot convert null_value to uint32")
	}
	return 0, errors.New("type conversion not supported")
}

// ToInt64 converts to int64 (default implementation)
func (v *BaseValue) ToInt64() (int64, error) {
	if v.IsNull() {
		return 0, errors.New("cannot convert null_value to int64")
	}
	return 0, errors.New("type conversion not supported")
}

// ToUInt64 converts to uint64 (default implementation)
func (v *BaseValue) ToUInt64() (uint64, error) {
	if v.IsNull() {
		return 0, errors.New("cannot convert null_value to uint64")
	}
	return 0, errors.New("type conversion not supported")
}

// ToFloat32 converts to float32 (default implementation)
func (v *BaseValue) ToFloat32() (float32, error) {
	if v.IsNull() {
		return 0, errors.New("cannot convert null_value to float32")
	}
	return 0, errors.New("type conversion not supported")
}

// ToFloat64 converts to float64 (default implementation)
func (v *BaseValue) ToFloat64() (float64, error) {
	if v.IsNull() {
		return 0, errors.New("cannot convert null_value to float64")
	}
	return 0, errors.New("type conversion not supported")
}

// ToString converts to string (default implementation)
func (v *BaseValue) ToString() (string, error) {
	if v.IsNull() {
		return "", nil
	}
	return "", errors.New("type conversion not supported")
}

// ToBytes converts to bytes
func (v *BaseValue) ToBytes() ([]byte, error) {
	return v.data, nil
}

// Serialize serializes the value to string
func (v *BaseValue) Serialize() (string, error) {
	return fmt.Sprintf("%s|%s|%d", v.name, v.vtype.String(), len(v.data)), nil
}

// ToXML converts to XML representation
func (v *BaseValue) ToXML() (string, error) {
	type XMLValue struct {
		XMLName xml.Name `xml:"value"`
		Name    string   `xml:"name,attr"`
		Type    string   `xml:"type,attr"`
		Data    string   `xml:",chardata"`
	}

	xmlVal := XMLValue{
		Name: v.name,
		Type: v.vtype.TypeName(),
		Data: string(v.data),
	}

	data, err := xml.MarshalIndent(xmlVal, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToJSON converts to JSON representation
func (v *BaseValue) ToJSON() (string, error) {
	jsonVal := map[string]interface{}{
		"name": v.name,
		"type": v.vtype.TypeName(),
		"data": string(v.data),
	}

	data, err := json.MarshalIndent(jsonVal, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Children returns child values (for container values)
func (v *BaseValue) Children() []Value {
	return v.units
}

// ChildCount returns the number of children
func (v *BaseValue) ChildCount() int {
	return len(v.units)
}

// GetChild gets a child by name and index
func (v *BaseValue) GetChild(name string, index int) Value {
	count := 0
	for _, child := range v.units {
		if child.Name() == name {
			if count == index {
				return child
			}
			count++
		}
	}
	return NewBaseValue("", NullValue, nil)
}

// AddChild adds a child value (default implementation throws error)
func (v *BaseValue) AddChild(child Value) error {
	return errors.New("cannot add child to non-container value")
}

// RemoveChild removes a child by name (default implementation throws error)
func (v *BaseValue) RemoveChild(name string) error {
	return errors.New("cannot remove child from non-container value")
}
