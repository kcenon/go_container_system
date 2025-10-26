/****************************************************************************
BSD 3-Clause License

Copyright (c) 2021, 🍀☀🌕🌥 🌊
All rights reserved.
****************************************************************************/

package core

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
)

// ValueContainer represents a message container with header and values
type ValueContainer struct {
	// Header fields
	sourceID    string
	sourceSubID string
	targetID    string
	targetSubID string
	messageType string
	version     string

	// Values
	units []Value
}

// NewValueContainer creates a new empty container
func NewValueContainer() *ValueContainer {
	return &ValueContainer{
		version: "1.0.0.0",
		units:   make([]Value, 0),
	}
}

// NewValueContainerWithType creates a container with message type
func NewValueContainerWithType(messageType string, units ...Value) *ValueContainer {
	return &ValueContainer{
		messageType: messageType,
		version:     "1.0.0.0",
		units:       units,
	}
}

// NewValueContainerWithTarget creates a container with target info
func NewValueContainerWithTarget(targetID, targetSubID, messageType string, units ...Value) *ValueContainer {
	return &ValueContainer{
		targetID:    targetID,
		targetSubID: targetSubID,
		messageType: messageType,
		version:     "1.0.0.0",
		units:       units,
	}
}

// NewValueContainerFull creates a container with full header
func NewValueContainerFull(sourceID, sourceSubID, targetID, targetSubID, messageType string, units ...Value) *ValueContainer {
	return &ValueContainer{
		sourceID:    sourceID,
		sourceSubID: sourceSubID,
		targetID:    targetID,
		targetSubID: targetSubID,
		messageType: messageType,
		version:     "1.0.0.0",
		units:       units,
	}
}

// SetSource sets the source ID and sub ID
func (c *ValueContainer) SetSource(sourceID, sourceSubID string) {
	c.sourceID = sourceID
	c.sourceSubID = sourceSubID
}

// SetTarget sets the target ID and sub ID
func (c *ValueContainer) SetTarget(targetID, targetSubID string) {
	c.targetID = targetID
	c.targetSubID = targetSubID
}

// SetMessageType sets the message type
func (c *ValueContainer) SetMessageType(messageType string) {
	c.messageType = messageType
}

// SwapHeader swaps source and target
func (c *ValueContainer) SwapHeader() {
	c.sourceID, c.targetID = c.targetID, c.sourceID
	c.sourceSubID, c.targetSubID = c.targetSubID, c.sourceSubID
}

// Accessors
func (c *ValueContainer) SourceID() string       { return c.sourceID }
func (c *ValueContainer) SourceSubID() string    { return c.sourceSubID }
func (c *ValueContainer) TargetID() string       { return c.targetID }
func (c *ValueContainer) TargetSubID() string    { return c.targetSubID }
func (c *ValueContainer) MessageType() string    { return c.messageType }
func (c *ValueContainer) Version() string        { return c.version }
func (c *ValueContainer) Values() []Value        { return c.units }

// AddValue adds a value to the container
func (c *ValueContainer) AddValue(value Value) {
	c.units = append(c.units, value)
}

// RemoveValue removes all values with the given name
func (c *ValueContainer) RemoveValue(name string) {
	newUnits := make([]Value, 0)
	for _, unit := range c.units {
		if unit.Name() != name {
			newUnits = append(newUnits, unit)
		}
	}
	c.units = newUnits
}

// GetValue gets the first value with the given name
func (c *ValueContainer) GetValue(name string, index int) Value {
	count := 0
	for _, unit := range c.units {
		if unit.Name() == name {
			if count == index {
				return unit
			}
			count++
		}
	}
	return NewBaseValue("", NullValue, nil)
}

// GetValues gets all values with the given name
func (c *ValueContainer) GetValues(name string) []Value {
	result := make([]Value, 0)
	for _, unit := range c.units {
		if unit.Name() == name {
			result = append(result, unit)
		}
	}
	return result
}

// ClearValues removes all values
func (c *ValueContainer) ClearValues() {
	c.units = make([]Value, 0)
}

// Copy creates a copy of this container
func (c *ValueContainer) Copy(containingValues bool) *ValueContainer {
	newContainer := &ValueContainer{
		sourceID:    c.sourceID,
		sourceSubID: c.sourceSubID,
		targetID:    c.targetID,
		targetSubID: c.targetSubID,
		messageType: c.messageType,
		version:     c.version,
		units:       make([]Value, 0),
	}

	if containingValues {
		newContainer.units = make([]Value, len(c.units))
		copy(newContainer.units, c.units)
	}

	return newContainer
}

// Serialize serializes the container to string format
func (c *ValueContainer) Serialize() (string, error) {
	// Header: sourceID|sourceSubID|targetID|targetSubID|messageType|version
	header := fmt.Sprintf("%s|%s|%s|%s|%s|%s",
		c.sourceID, c.sourceSubID, c.targetID, c.targetSubID,
		c.messageType, c.version)

	// Values
	valueStrs := make([]string, len(c.units))
	for i, unit := range c.units {
		valStr, err := unit.Serialize()
		if err != nil {
			return "", err
		}
		valueStrs[i] = valStr
	}

	data := strings.Join(valueStrs, "|")
	return fmt.Sprintf("%s\n%s", header, data), nil
}

// SerializeArray serializes the container to byte array
func (c *ValueContainer) SerializeArray() ([]byte, error) {
	str, err := c.Serialize()
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

// Deserialize deserializes from string
func (c *ValueContainer) Deserialize(data string) error {
	lines := strings.Split(data, "\n")
	if len(lines) < 1 {
		return fmt.Errorf("invalid data format")
	}

	// Parse header
	headerParts := strings.Split(lines[0], "|")
	if len(headerParts) >= 6 {
		c.sourceID = headerParts[0]
		c.sourceSubID = headerParts[1]
		c.targetID = headerParts[2]
		c.targetSubID = headerParts[3]
		c.messageType = headerParts[4]
		c.version = headerParts[5]
	}

	// TODO: Parse values from lines[1] if present
	// This would require a value factory to create values based on type

	return nil
}

// DeserializeArray deserializes from byte array
func (c *ValueContainer) DeserializeArray(data []byte) error {
	return c.Deserialize(string(data))
}

// ToXML converts to XML representation
func (c *ValueContainer) ToXML() (string, error) {
	type XMLContainer struct {
		XMLName     xml.Name `xml:"container"`
		SourceID    string   `xml:"source_id"`
		SourceSubID string   `xml:"source_sub_id"`
		TargetID    string   `xml:"target_id"`
		TargetSubID string   `xml:"target_sub_id"`
		MessageType string   `xml:"message_type"`
		Version     string   `xml:"version"`
		Values      []string `xml:"values>value"`
	}

	xmlCont := XMLContainer{
		SourceID:    c.sourceID,
		SourceSubID: c.sourceSubID,
		TargetID:    c.targetID,
		TargetSubID: c.targetSubID,
		MessageType: c.messageType,
		Version:     c.version,
		Values:      make([]string, 0),
	}

	for _, unit := range c.units {
		unitXML, err := unit.ToXML()
		if err != nil {
			return "", err
		}
		xmlCont.Values = append(xmlCont.Values, unitXML)
	}

	data, err := xml.MarshalIndent(xmlCont, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToJSON converts to JSON representation
func (c *ValueContainer) ToJSON() (string, error) {
	jsonCont := map[string]interface{}{
		"source_id":     c.sourceID,
		"source_sub_id": c.sourceSubID,
		"target_id":     c.targetID,
		"target_sub_id": c.targetSubID,
		"message_type":  c.messageType,
		"version":       c.version,
		"values":        make([]map[string]interface{}, 0),
	}

	values := make([]map[string]interface{}, 0)
	for _, unit := range c.units {
		unitJSON, err := unit.ToJSON()
		if err != nil {
			return "", err
		}
		var unitMap map[string]interface{}
		if err := json.Unmarshal([]byte(unitJSON), &unitMap); err != nil {
			return "", err
		}
		values = append(values, unitMap)
	}
	jsonCont["values"] = values

	data, err := json.MarshalIndent(jsonCont, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
