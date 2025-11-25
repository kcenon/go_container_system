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

// ContainerValue represents a container that can hold multiple child values
type ContainerValue struct {
	*core.BaseValue
	children []core.Value
}

// NewContainerValue creates a new container value
func NewContainerValue(name string, children ...core.Value) *ContainerValue {
	cv := &ContainerValue{
		BaseValue: core.NewBaseValue(name, core.ContainerValue, nil),
		children:  make([]core.Value, 0),
	}
	for _, child := range children {
		cv.children = append(cv.children, child)
	}
	return cv
}

// Children returns all child values
func (v *ContainerValue) Children() []core.Value {
	return v.children
}

// ChildCount returns the number of children
func (v *ContainerValue) ChildCount() int {
	return len(v.children)
}

// GetChild gets a child by name and index
func (v *ContainerValue) GetChild(name string, index int) core.Value {
	count := 0
	for _, child := range v.children {
		if child.Name() == name {
			if count == index {
				return child
			}
			count++
		}
	}
	return core.NewBaseValue("", core.NullValue, nil)
}

// AddChild adds a child value
func (v *ContainerValue) AddChild(child core.Value) error {
	v.children = append(v.children, child)
	return nil
}

// RemoveChild removes all children with the given name
func (v *ContainerValue) RemoveChild(name string) error {
	newChildren := make([]core.Value, 0)
	for _, child := range v.children {
		if child.Name() != name {
			newChildren = append(newChildren, child)
		}
	}
	v.children = newChildren
	return nil
}

// Serialize serializes the container and all its children to C++ compatible format.
//
// Format: [name,type_code,child_count];[child1][child2]...
// This matches the Python/C++ wire protocol format for cross-language compatibility.
func (v *ContainerValue) Serialize() (string, error) {
	// Container header with child count (type 14 = container_value)
	result := fmt.Sprintf("[%s,%s,%d];", v.Name(), v.Type().String(), len(v.children))

	// Append all child serializations (recursive)
	for _, child := range v.children {
		childSer, err := child.Serialize()
		if err != nil {
			return "", err
		}
		result += childSer
	}
	return result, nil
}

// ToXML converts to XML representation
func (v *ContainerValue) ToXML() (string, error) {
	type XMLContainer struct {
		XMLName  xml.Name `xml:"container"`
		Name     string   `xml:"name,attr"`
		Type     string   `xml:"type,attr"`
		Children []string `xml:"child"`
	}

	xmlCont := XMLContainer{
		Name:     v.Name(),
		Type:     v.Type().TypeName(),
		Children: make([]string, 0),
	}

	for _, child := range v.children {
		childXML, err := child.ToXML()
		if err != nil {
			return "", err
		}
		xmlCont.Children = append(xmlCont.Children, childXML)
	}

	data, err := xml.MarshalIndent(xmlCont, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToJSON converts to JSON representation
func (v *ContainerValue) ToJSON() (string, error) {
	jsonCont := map[string]interface{}{
		"name":     v.Name(),
		"type":     v.Type().TypeName(),
		"children": make([]map[string]interface{}, 0),
	}

	children := make([]map[string]interface{}, 0)
	for _, child := range v.children {
		childJSON, err := child.ToJSON()
		if err != nil {
			return "", err
		}
		var childMap map[string]interface{}
		if err := json.Unmarshal([]byte(childJSON), &childMap); err != nil {
			return "", err
		}
		children = append(children, childMap)
	}
	jsonCont["children"] = children

	data, err := json.MarshalIndent(jsonCont, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
