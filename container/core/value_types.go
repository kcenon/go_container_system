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

// ValueType represents the type of value stored in the container
type ValueType int

const (
	NullValue ValueType = iota
	BoolValue
	ShortValue
	UShortValue
	IntValue
	UIntValue
	LongValue
	ULongValue
	LLongValue
	ULLongValue
	FloatValue
	DoubleValue
	BytesValue
	StringValue
	ContainerValue
	ArrayValue
)

// String returns the string representation of the value type
func (vt ValueType) String() string {
	switch vt {
	case NullValue:
		return "0"
	case BoolValue:
		return "1"
	case ShortValue:
		return "2"
	case UShortValue:
		return "3"
	case IntValue:
		return "4"
	case UIntValue:
		return "5"
	case LongValue:
		return "6"
	case ULongValue:
		return "7"
	case LLongValue:
		return "8"
	case ULLongValue:
		return "9"
	case FloatValue:
		return "10"
	case DoubleValue:
		return "11"
	case BytesValue:
		return "12"
	case StringValue:
		return "13"
	case ContainerValue:
		return "14"
	case ArrayValue:
		return "15"
	default:
		return "0"
	}
}

// ParseValueType converts a string to a ValueType
func ParseValueType(s string) ValueType {
	switch s {
	case "0":
		return NullValue
	case "1":
		return BoolValue
	case "2":
		return ShortValue
	case "3":
		return UShortValue
	case "4":
		return IntValue
	case "5":
		return UIntValue
	case "6":
		return LongValue
	case "7":
		return ULongValue
	case "8":
		return LLongValue
	case "9":
		return ULLongValue
	case "10":
		return FloatValue
	case "11":
		return DoubleValue
	case "12":
		return BytesValue
	case "13":
		return StringValue
	case "14":
		return ContainerValue
	case "15":
		return ArrayValue
	default:
		return NullValue
	}
}

// TypeName returns a human-readable name for the value type
func (vt ValueType) TypeName() string {
	switch vt {
	case NullValue:
		return "null"
	case BoolValue:
		return "bool"
	case ShortValue:
		return "short"
	case UShortValue:
		return "ushort"
	case IntValue:
		return "int"
	case UIntValue:
		return "uint"
	case LongValue:
		return "long"
	case ULongValue:
		return "ulong"
	case LLongValue:
		return "llong"
	case ULLongValue:
		return "ullong"
	case FloatValue:
		return "float"
	case DoubleValue:
		return "double"
	case BytesValue:
		return "bytes"
	case StringValue:
		return "string"
	case ContainerValue:
		return "container"
	case ArrayValue:
		return "array"
	default:
		return "unknown"
	}
}
