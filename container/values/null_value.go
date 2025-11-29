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

// Package values provides concrete implementations of the Value interface.
package values

import (
	"encoding/json"
	"encoding/xml"

	"github.com/kcenon/go_container_system/container/core"
)

// NullValue represents an absent or undefined value.
// This is useful for optional fields or placeholder values.
//
// Wire Protocol:
//   - Type code: 0
//   - Payload size: 0 bytes (no payload)
//
// C++ Compatibility:
// This type corresponds to `null_value` in the C++ container system,
// which uses `std::monostate` as the underlying type.
type NullValue struct {
	*core.BaseValue
}

// NewNullValue creates a new null value with the specified name.
//
// Example:
//
//	nullVal := values.NewNullValue("optional_field")
//	fmt.Println(nullVal.IsNull()) // true
func NewNullValue(name string) *NullValue {
	return &NullValue{
		BaseValue: core.NewBaseValue(name, core.NullValue, nil),
	}
}

// ToString returns "null" for null values.
func (v *NullValue) ToString() (string, error) {
	return "null", nil
}

// ToJSON returns the JSON representation of the null value.
func (v *NullValue) ToJSON() (string, error) {
	jsonVal := map[string]interface{}{
		"name":  v.Name(),
		"type":  "null",
		"value": nil,
	}

	data, err := json.MarshalIndent(jsonVal, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToXML returns the XML representation of the null value.
func (v *NullValue) ToXML() (string, error) {
	type XMLNullValue struct {
		XMLName xml.Name `xml:"value"`
		Name    string   `xml:"name,attr"`
		Type    string   `xml:"type,attr"`
	}

	xmlVal := XMLNullValue{
		Name: v.Name(),
		Type: "null",
	}

	data, err := xml.MarshalIndent(xmlVal, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Serialize returns the string representation for the null value.
func (v *NullValue) Serialize() (string, error) {
	// Format: name|type|size (size is always 0 for null)
	return v.Name() + "|0|0", nil
}
