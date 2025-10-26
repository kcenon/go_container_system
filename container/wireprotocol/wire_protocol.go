package wireprotocol

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/kcenon/go_container_system/container/core"
	"github.com/kcenon/go_container_system/container/values"
)

// C++ header field IDs (matching container.cpp constants)
const (
	targetIDField      = 1
	targetSubIDField   = 2
	sourceIDField      = 3
	sourceSubIDField   = 4
	messageTypeField   = 5
	messageVersionField = 6
)

// SerializeCppWire serializes a ValueContainer to C++ wire protocol format
//
// Format: @header={{[id,value];...}};@data={{[name,type,data];...}};
//
// This produces byte-for-byte compatible output with C++ container_system
// and Python container_system for cross-language data exchange.
//
// Example:
//   container := core.NewValueContainer()
//   container.SetSource("client", "session")
//   container.AddValue(values.NewInt32Value("count", 42))
//   wireData := wireprotocol.SerializeCppWire(container)
//   // Result: @header={{[3,client];[4,session];[5,data_container];[6,1.0.0.0];}};@data={{[count,int_value,42];}};
func SerializeCppWire(c *core.ValueContainer) (string, error) {
	var result strings.Builder
	result.Grow(512) // Pre-allocate buffer

	// Serialize header
	result.WriteString("@header={{")

	// Only include routing fields if message_type is not "data_container"
	messageType := c.MessageType()
	if messageType != "data_container" {
		targetID := c.TargetID()
		targetSubID := c.TargetSubID()
		if targetID != "" || targetSubID != "" {
			result.WriteString(fmt.Sprintf("[%d,%s];", targetIDField, targetID))
			result.WriteString(fmt.Sprintf("[%d,%s];", targetSubIDField, targetSubID))
		}
		sourceID := c.SourceID()
		sourceSubID := c.SourceSubID()
		if sourceID != "" || sourceSubID != "" {
			result.WriteString(fmt.Sprintf("[%d,%s];", sourceIDField, sourceID))
			result.WriteString(fmt.Sprintf("[%d,%s];", sourceSubIDField, sourceSubID))
		}
	}

	// Always include message_type and version
	result.WriteString(fmt.Sprintf("[%d,%s];", messageTypeField, messageType))
	result.WriteString(fmt.Sprintf("[%d,%s];", messageVersionField, c.Version()))
	result.WriteString("}};")

	// Serialize data
	result.WriteString("@data={{")

	// Serialize all values
	for _, value := range c.Values() {
		serialized, err := serializeValueCpp(value)
		if err != nil {
			// Skip values that fail to serialize
			continue
		}
		result.WriteString(serialized)
	}

	result.WriteString("}};")

	return result.String(), nil
}

// serializeValueCpp serializes a single value to C++ wire protocol format
//
// Format: [name,type_name,data];
func serializeValueCpp(value core.Value) (string, error) {
	name := value.Name()
	valueType := value.Type()
	typeName := valueTypeToCppName(valueType)

	// Serialize data based on type
	var dataStr string
	switch valueType {
	case core.BoolValue:
		boolVal, err := value.ToBool()
		if err != nil {
			return "", err
		}
		if boolVal {
			dataStr = "true"
		} else {
			dataStr = "false"
		}
	case core.ShortValue:
		val, err := value.ToInt16()
		if err != nil {
			return "", err
		}
		dataStr = fmt.Sprintf("%d", val)
	case core.UShortValue:
		val, err := value.ToUInt16()
		if err != nil {
			return "", err
		}
		dataStr = fmt.Sprintf("%d", val)
	case core.IntValue:
		val, err := value.ToInt32()
		if err != nil {
			return "", err
		}
		dataStr = fmt.Sprintf("%d", val)
	case core.UIntValue:
		val, err := value.ToUInt32()
		if err != nil {
			return "", err
		}
		dataStr = fmt.Sprintf("%d", val)
	case core.LongValue:
		val, err := value.ToInt64()
		if err != nil {
			return "", err
		}
		dataStr = fmt.Sprintf("%d", val)
	case core.ULongValue:
		val, err := value.ToUInt64()
		if err != nil {
			return "", err
		}
		dataStr = fmt.Sprintf("%d", val)
	case core.LLongValue:
		val, err := value.ToInt64()
		if err != nil {
			return "", err
		}
		dataStr = fmt.Sprintf("%d", val)
	case core.ULLongValue:
		val, err := value.ToUInt64()
		if err != nil {
			return "", err
		}
		dataStr = fmt.Sprintf("%d", val)
	case core.FloatValue:
		val, err := value.ToFloat32()
		if err != nil {
			return "", err
		}
		dataStr = fmt.Sprintf("%g", val)
	case core.DoubleValue:
		val, err := value.ToFloat64()
		if err != nil {
			return "", err
		}
		dataStr = fmt.Sprintf("%g", val)
	case core.StringValue:
		val, err := value.ToString()
		if err != nil {
			return "", err
		}
		dataStr = val
	case core.BytesValue:
		// Convert bytes to hex string (matching C++ hex encoding)
		bytes, err := value.ToBytes()
		if err != nil {
			return "", err
		}
		dataStr = hex.EncodeToString(bytes)
	case core.ContainerValue:
		// For containers, store child count (matching C++ behavior)
		dataStr = "0" // Placeholder - full support requires nested container work
	case core.NullValue:
		dataStr = ""
	default:
		return "", fmt.Errorf("unsupported value type: %d", valueType)
	}

	return fmt.Sprintf("[%s,%s,%s];", name, typeName, dataStr), nil
}

// valueTypeToCppName converts ValueType to C++ type name string
func valueTypeToCppName(vt core.ValueType) string {
	switch vt {
	case core.BoolValue:
		return "bool_value"
	case core.ShortValue:
		return "short_value"
	case core.UShortValue:
		return "ushort_value"
	case core.IntValue:
		return "int_value"
	case core.UIntValue:
		return "uint_value"
	case core.LongValue:
		return "long_value"
	case core.ULongValue:
		return "ulong_value"
	case core.LLongValue:
		return "llong_value"
	case core.ULLongValue:
		return "ullong_value"
	case core.FloatValue:
		return "float_value"
	case core.DoubleValue:
		return "double_value"
	case core.StringValue:
		return "string_value"
	case core.BytesValue:
		return "bytes_value"
	case core.ContainerValue:
		return "container_value"
	case core.NullValue:
		return "null_value"
	default:
		return "null_value"
	}
}

// cppNameToValueType converts C++ type name string to ValueType
func cppNameToValueType(name string) (core.ValueType, error) {
	switch name {
	case "bool_value":
		return core.BoolValue, nil
	case "short_value":
		return core.ShortValue, nil
	case "ushort_value":
		return core.UShortValue, nil
	case "int_value":
		return core.IntValue, nil
	case "uint_value":
		return core.UIntValue, nil
	case "long_value":
		return core.LongValue, nil
	case "ulong_value":
		return core.ULongValue, nil
	case "llong_value":
		return core.LLongValue, nil
	case "ullong_value":
		return core.ULLongValue, nil
	case "float_value":
		return core.FloatValue, nil
	case "double_value":
		return core.DoubleValue, nil
	case "string_value":
		return core.StringValue, nil
	case "bytes_value":
		return core.BytesValue, nil
	case "container_value":
		return core.ContainerValue, nil
	case "null_value":
		return core.NullValue, nil
	default:
		return core.NullValue, fmt.Errorf("unknown C++ type name: %s", name)
	}
}

// DeserializeCppWire deserializes a ValueContainer from C++ wire protocol format
//
// This can parse data generated by C++ container_system, Python container_system,
// Rust container_system, or any other system using the C++ wire protocol.
//
// Format: @header={{[id,value];...}};@data={{[name,type,data];...}};
func DeserializeCppWire(wireData string) (*core.ValueContainer, error) {
	// Remove newlines for easier parsing
	cleanData := strings.ReplaceAll(wireData, "\r\n", "")
	cleanData = strings.ReplaceAll(cleanData, "\n", "")

	container := core.NewValueContainer()

	// Parse header section
	headerRegex := regexp.MustCompile(`@header=\s*\{\{?\s*(.*?)\s*\}\}?;`)
	headerMatch := headerRegex.FindStringSubmatch(cleanData)

	var targetID, targetSubID, sourceID, sourceSubID, messageType string

	if len(headerMatch) > 1 {
		headerContent := headerMatch[1]

		// Parse header pairs: [id,value];
		pairRegex := regexp.MustCompile(`\[(\d+),(.*?)\];`)
		pairMatches := pairRegex.FindAllStringSubmatch(headerContent, -1)

		for _, match := range pairMatches {
			if len(match) < 3 {
				continue
			}

			id, err := strconv.Atoi(match[1])
			if err != nil {
				continue
			}
			value := strings.TrimSpace(match[2])

			switch id {
			case targetIDField:
				targetID = value
			case targetSubIDField:
				targetSubID = value
			case sourceIDField:
				sourceID = value
			case sourceSubIDField:
				sourceSubID = value
			case messageTypeField:
				messageType = value
			// version field is read-only, ignore it
			}
		}
	}

	// Apply header fields to container
	if targetID != "" || targetSubID != "" {
		container.SetTarget(targetID, targetSubID)
	}
	if sourceID != "" || sourceSubID != "" {
		container.SetSource(sourceID, sourceSubID)
	}
	if messageType != "" {
		container.SetMessageType(messageType)
	}

	// Parse data section
	dataRegex := regexp.MustCompile(`@data=\s*\{\{?\s*(.*?)\s*\}\}?;`)
	dataMatch := dataRegex.FindStringSubmatch(cleanData)

	if len(dataMatch) > 1 {
		dataContent := dataMatch[1]

		// Parse value items: [name,type,data];
		itemRegex := regexp.MustCompile(`\[(\w+),\s*(\w+),\s*(.*?)\];`)
		itemMatches := itemRegex.FindAllStringSubmatch(dataContent, -1)

		for _, match := range itemMatches {
			if len(match) < 4 {
				continue
			}

			name := match[1]
			typeName := match[2]
			dataStr := match[3]

			valueType, err := cppNameToValueType(typeName)
			if err != nil {
				continue
			}

			// Parse value based on type
			var parsedValue core.Value
			switch valueType {
			case core.BoolValue:
				val := (dataStr == "true")
				parsedValue = values.NewBoolValue(name, val)

			case core.ShortValue:
				val, err := strconv.ParseInt(dataStr, 10, 16)
				if err != nil {
					continue
				}
				parsedValue = values.NewInt16Value(name, int16(val))

			case core.UShortValue:
				val, err := strconv.ParseUint(dataStr, 10, 16)
				if err != nil {
					continue
				}
				parsedValue = values.NewUInt16Value(name, uint16(val))

			case core.IntValue:
				val, err := strconv.ParseInt(dataStr, 10, 32)
				if err != nil {
					continue
				}
				parsedValue = values.NewInt32Value(name, int32(val))

			case core.UIntValue:
				val, err := strconv.ParseUint(dataStr, 10, 32)
				if err != nil {
					continue
				}
				parsedValue = values.NewUInt32Value(name, uint32(val))

			case core.LongValue, core.LLongValue:
				val, err := strconv.ParseInt(dataStr, 10, 64)
				if err != nil {
					continue
				}
				parsedValue = values.NewInt64Value(name, val)

			case core.ULongValue, core.ULLongValue:
				val, err := strconv.ParseUint(dataStr, 10, 64)
				if err != nil {
					continue
				}
				parsedValue = values.NewUInt64Value(name, val)

			case core.FloatValue:
				val, err := strconv.ParseFloat(dataStr, 32)
				if err != nil {
					continue
				}
				parsedValue = values.NewFloat32Value(name, float32(val))

			case core.DoubleValue:
				val, err := strconv.ParseFloat(dataStr, 64)
				if err != nil {
					continue
				}
				parsedValue = values.NewFloat64Value(name, val)

			case core.StringValue:
				parsedValue = values.NewStringValue(name, dataStr)

			case core.BytesValue:
				// Decode hex string
				bytes, err := hex.DecodeString(dataStr)
				if err != nil {
					continue
				}
				parsedValue = values.NewBytesValue(name, bytes)

			case core.ContainerValue, core.NullValue:
				// TODO: Implement nested container support
				continue

			default:
				continue
			}

			if parsedValue != nil {
				container.AddValue(parsedValue)
			}
		}
	}

	return container, nil
}
