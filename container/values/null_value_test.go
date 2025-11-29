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

package values

import (
	"strings"
	"testing"

	"github.com/kcenon/go_container_system/container/core"
)

func TestNewNullValue(t *testing.T) {
	nv := NewNullValue("test_null")

	if nv.Name() != "test_null" {
		t.Errorf("Expected name 'test_null', got '%s'", nv.Name())
	}

	if nv.Type() != core.NullValue {
		t.Errorf("Expected type NullValue, got %v", nv.Type())
	}

	if !nv.IsNull() {
		t.Error("Expected IsNull() to return true")
	}

	if nv.Size() != 0 {
		t.Errorf("Expected size 0, got %d", nv.Size())
	}
}

func TestNullValueToString(t *testing.T) {
	nv := NewNullValue("test")
	str, err := nv.ToString()
	if err != nil {
		t.Errorf("ToString failed: %v", err)
	}
	if str != "null" {
		t.Errorf("Expected 'null', got '%s'", str)
	}
}

func TestNullValueToJSON(t *testing.T) {
	nv := NewNullValue("test_field")
	json, err := nv.ToJSON()
	if err != nil {
		t.Errorf("ToJSON failed: %v", err)
	}

	if !strings.Contains(json, `"name": "test_field"`) {
		t.Errorf("JSON should contain name field, got: %s", json)
	}
	if !strings.Contains(json, `"type": "null"`) {
		t.Errorf("JSON should contain type field, got: %s", json)
	}
	if !strings.Contains(json, `"value": null`) {
		t.Errorf("JSON should contain value field as null, got: %s", json)
	}
}

func TestNullValueToXML(t *testing.T) {
	nv := NewNullValue("test_field")
	xml, err := nv.ToXML()
	if err != nil {
		t.Errorf("ToXML failed: %v", err)
	}

	if !strings.Contains(xml, `name="test_field"`) {
		t.Errorf("XML should contain name attribute, got: %s", xml)
	}
	if !strings.Contains(xml, `type="null"`) {
		t.Errorf("XML should contain type attribute, got: %s", xml)
	}
}

func TestNullValueSerialize(t *testing.T) {
	nv := NewNullValue("optional")
	ser, err := nv.Serialize()
	if err != nil {
		t.Errorf("Serialize failed: %v", err)
	}

	// Format: name|type|size
	if ser != "optional|0|0" {
		t.Errorf("Expected 'optional|0|0', got '%s'", ser)
	}
}

func TestNullValueTypeConversions(t *testing.T) {
	nv := NewNullValue("test")

	// All type conversions should fail for null values
	if _, err := nv.ToBool(); err == nil {
		t.Error("ToBool should return error for null value")
	}

	if _, err := nv.ToInt32(); err == nil {
		t.Error("ToInt32 should return error for null value")
	}

	if _, err := nv.ToInt64(); err == nil {
		t.Error("ToInt64 should return error for null value")
	}

	if _, err := nv.ToFloat64(); err == nil {
		t.Error("ToFloat64 should return error for null value")
	}
}

func TestNullValueData(t *testing.T) {
	nv := NewNullValue("empty")
	data := nv.Data()

	if data != nil && len(data) != 0 {
		t.Errorf("Expected nil or empty data, got %v", data)
	}
}
