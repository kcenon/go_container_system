/****************************************************************************
BSD 3-Clause License

Copyright (c) 2021, 🍀☀🌕🌥 🌊
All rights reserved.
****************************************************************************/

package values

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/kcenon/go_container_system/container/core"
)

// ArrayValue represents an array/list that can hold multiple elements
// ArrayValue (type 15) is an extension to support homogeneous or heterogeneous
// collections of values, similar to JSON arrays.
//
// Wire format (binary):
// [type:1=15][name_len:4 LE][name:UTF-8][value_size:4 LE][count:4 LE][values...]
//
// Text format:
// [name,15,count];[element1][element2]...
type ArrayValue struct {
	*core.BaseValue
	elements []core.Value
}

// NewArrayValue creates a new array value
func NewArrayValue(name string, elements ...core.Value) *ArrayValue {
	av := &ArrayValue{
		BaseValue: core.NewBaseValue(name, core.ArrayValue, nil),
		elements:  make([]core.Value, 0),
	}
	for _, element := range elements {
		av.elements = append(av.elements, element)
	}
	return av
}

// Elements returns all elements
func (v *ArrayValue) Elements() []core.Value {
	return v.elements
}

// Count returns the number of elements
func (v *ArrayValue) Count() int {
	return len(v.elements)
}

// IsEmpty checks if the array is empty
func (v *ArrayValue) IsEmpty() bool {
	return len(v.elements) == 0
}

// At gets element at index
func (v *ArrayValue) At(index int) (core.Value, error) {
	if index < 0 || index >= len(v.elements) {
		return nil, fmt.Errorf("ArrayValue index %d out of range (size: %d)", index, len(v.elements))
	}
	return v.elements[index], nil
}

// Append adds an element to the end of the array
func (v *ArrayValue) Append(element core.Value) error {
	v.elements = append(v.elements, element)
	return nil
}

// PushBack adds an element (C++ compatibility name)
func (v *ArrayValue) PushBack(element core.Value) error {
	return v.Append(element)
}

// Clear removes all elements
func (v *ArrayValue) Clear() {
	v.elements = make([]core.Value, 0)
}

// Serialize serializes the array and all its elements
func (v *ArrayValue) Serialize() (string, error) {
	result := fmt.Sprintf("[%s,%s,%d];", v.Name(), v.Type().String(), len(v.elements))
	for _, element := range v.elements {
		elemSer, err := element.Serialize()
		if err != nil {
			return "", err
		}
		result += elemSer
	}
	return result, nil
}

// ToXML converts to XML representation
func (v *ArrayValue) ToXML() (string, error) {
	type XMLArray struct {
		XMLName  xml.Name `xml:"array"`
		Name     string   `xml:"name,attr"`
		Type     string   `xml:"type,attr"`
		Count    int      `xml:"count,attr"`
		Elements []string `xml:"element"`
	}

	xmlArr := XMLArray{
		Name:     v.Name(),
		Type:     v.Type().TypeName(),
		Count:    len(v.elements),
		Elements: make([]string, 0),
	}

	for _, element := range v.elements {
		elemXML, err := element.ToXML()
		if err != nil {
			return "", err
		}
		xmlArr.Elements = append(xmlArr.Elements, elemXML)
	}

	data, err := xml.Marshal(xmlArr)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ToJSON converts to JSON representation
func (v *ArrayValue) ToJSON() (string, error) {
	type JSONArray struct {
		Name     string        `json:"name"`
		Type     string        `json:"type"`
		Elements []interface{} `json:"elements"`
	}

	jsonArr := JSONArray{
		Name:     v.Name(),
		Type:     "array",
		Elements: make([]interface{}, 0),
	}

	for _, element := range v.elements {
		elemJSON, err := element.ToJSON()
		if err != nil {
			return "", err
		}
		var elemData interface{}
		if err := json.Unmarshal([]byte(elemJSON), &elemData); err != nil {
			return "", err
		}
		jsonArr.Elements = append(jsonArr.Elements, elemData)
	}

	data, err := json.MarshalIndent(jsonArr, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Data returns a human-readable description as bytes
func (v *ArrayValue) Data() []byte {
	return []byte(fmt.Sprintf("Array(%d elements)", len(v.elements)))
}

// Size returns the size in bytes (for serialization)
func (v *ArrayValue) Size() int {
	size := 4 // count (4 bytes)
	for _, element := range v.elements {
		size += element.Size()
	}
	return size
}
