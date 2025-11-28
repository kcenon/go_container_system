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

// Package core provides the fundamental types and interfaces for the container system.
package core

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"sync"
	"sync/atomic"
)

// BinaryVersion is the version byte for binary serialization format
const BinaryVersion uint8 = 1

// ValueStore provides domain-agnostic value storage.
//
// Pure value storage layer without messaging-specific fields.
// Can be used as a general-purpose serialization container.
//
// Features:
//   - Type-safe value storage
//   - JSON/Binary serialization support
//   - Thread-safe operations (optional)
//   - Key-value storage interface
//   - Statistics tracking
//
// Note: This is part of the domain separation initiative.
// See ValueContainer for messaging-specific wrapper.
type ValueStore struct {
	// Key-value storage
	values map[string]Value

	// Thread safety
	mutex            sync.RWMutex
	threadSafeEnabled atomic.Bool

	// Statistics
	readCount          atomic.Uint64
	writeCount         atomic.Uint64
	serializationCount atomic.Uint64
}

// NewValueStore creates a new empty ValueStore
func NewValueStore() *ValueStore {
	return &ValueStore{
		values: make(map[string]Value),
	}
}

// =========================================================================
// Value Management
// =========================================================================

// Add adds a value with a key.
// If the key already exists, the value will be overwritten.
// Thread-safe if EnableThreadSafety was called.
func (vs *ValueStore) Add(key string, value Value) {
	if vs.threadSafeEnabled.Load() {
		vs.mutex.Lock()
		defer vs.mutex.Unlock()
	}

	vs.values[key] = value
	vs.writeCount.Add(1)
}

// Get retrieves a value by key.
// Returns nil if the key doesn't exist.
// Thread-safe if EnableThreadSafety was called.
func (vs *ValueStore) Get(key string) Value {
	if vs.threadSafeEnabled.Load() {
		vs.mutex.RLock()
		defer vs.mutex.RUnlock()
	}

	value, exists := vs.values[key]
	if exists {
		vs.readCount.Add(1)
		return value
	}
	return nil
}

// Contains checks if a key exists.
func (vs *ValueStore) Contains(key string) bool {
	if vs.threadSafeEnabled.Load() {
		vs.mutex.RLock()
		defer vs.mutex.RUnlock()
	}

	_, exists := vs.values[key]
	return exists
}

// Remove removes a value by key.
// Returns true if removed, false if not found.
func (vs *ValueStore) Remove(key string) bool {
	if vs.threadSafeEnabled.Load() {
		vs.mutex.Lock()
		defer vs.mutex.Unlock()
	}

	_, exists := vs.values[key]
	if exists {
		delete(vs.values, key)
		return true
	}
	return false
}

// Clear removes all values.
func (vs *ValueStore) Clear() {
	if vs.threadSafeEnabled.Load() {
		vs.mutex.Lock()
		defer vs.mutex.Unlock()
	}

	vs.values = make(map[string]Value)
}

// Size returns the number of stored values.
func (vs *ValueStore) Size() int {
	if vs.threadSafeEnabled.Load() {
		vs.mutex.RLock()
		defer vs.mutex.RUnlock()
	}

	return len(vs.values)
}

// Empty checks if the store is empty.
func (vs *ValueStore) Empty() bool {
	return vs.Size() == 0
}

// Keys returns all keys.
func (vs *ValueStore) Keys() []string {
	if vs.threadSafeEnabled.Load() {
		vs.mutex.RLock()
		defer vs.mutex.RUnlock()
	}

	keys := make([]string, 0, len(vs.values))
	for key := range vs.values {
		keys = append(keys, key)
	}
	return keys
}

// Values returns all values.
func (vs *ValueStore) Values() []Value {
	if vs.threadSafeEnabled.Load() {
		vs.mutex.RLock()
		defer vs.mutex.RUnlock()
	}

	values := make([]Value, 0, len(vs.values))
	for _, value := range vs.values {
		values = append(values, value)
	}
	return values
}

// =========================================================================
// Serialization
// =========================================================================

// valueJSON is used for JSON serialization
type valueJSON struct {
	Name string      `json:"name"`
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Serialize serializes to JSON string.
func (vs *ValueStore) Serialize() (string, error) {
	if vs.threadSafeEnabled.Load() {
		vs.mutex.RLock()
		defer vs.mutex.RUnlock()
	}

	vs.serializationCount.Add(1)

	result := make(map[string]valueJSON)

	for key, value := range vs.values {
		dataStr, err := value.ToString()
		if err != nil {
			dataStr = ""
		}

		result[key] = valueJSON{
			Name: value.Name(),
			Type: value.Type().TypeName(),
			Data: dataStr,
		}
	}

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// SerializeBinary serializes to binary format.
//
// Binary format:
//   - Version byte (1)
//   - Number of entries (4 bytes, uint32, little-endian)
//   - For each entry:
//   - Key length (4 bytes, uint32, little-endian)
//   - Key data (UTF-8)
//   - Value type (1 byte)
//   - Value length (4 bytes, uint32, little-endian)
//   - Value data
func (vs *ValueStore) SerializeBinary() ([]byte, error) {
	if vs.threadSafeEnabled.Load() {
		vs.mutex.RLock()
		defer vs.mutex.RUnlock()
	}

	vs.serializationCount.Add(1)

	// Pre-calculate size for efficiency
	size := 1 + 4 // version + count
	for key, value := range vs.values {
		size += 4 + len(key) + 1 + 4 + len(value.Data())
	}

	result := make([]byte, 0, size)

	// Version byte
	result = append(result, BinaryVersion)

	// Number of entries
	countBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(countBytes, uint32(len(vs.values)))
	result = append(result, countBytes...)

	// Serialize each key-value pair
	for key, value := range vs.values {
		// Key length and key
		keyBytes := []byte(key)
		keyLenBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(keyLenBytes, uint32(len(keyBytes)))
		result = append(result, keyLenBytes...)
		result = append(result, keyBytes...)

		// Value type
		result = append(result, byte(value.Type()))

		// Value data
		valueData := value.Data()
		valueLenBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(valueLenBytes, uint32(len(valueData)))
		result = append(result, valueLenBytes...)
		result = append(result, valueData...)
	}

	return result, nil
}

// DeserializeBinary deserializes from binary format.
// Note: This requires a ValueFactory function to create values from type and data.
// For now, this returns an error as it requires integration with value factories.
func DeserializeBinary(data []byte, factory func(name string, vtype ValueType, data []byte) (Value, error)) (*ValueStore, error) {
	if len(data) < 5 {
		return nil, errors.New("invalid data: too small")
	}

	offset := 0

	// Read version
	version := data[offset]
	offset++

	if version != BinaryVersion {
		return nil, errors.New("unsupported binary version")
	}

	// Read count
	count := binary.LittleEndian.Uint32(data[offset:])
	offset += 4

	store := NewValueStore()

	// Read each key-value pair
	for i := uint32(0); i < count; i++ {
		if offset+4 > len(data) {
			return nil, errors.New("truncated data at entry")
		}

		// Read key length
		keyLen := binary.LittleEndian.Uint32(data[offset:])
		offset += 4

		if offset+int(keyLen)+5 > len(data) {
			return nil, errors.New("truncated key data")
		}

		// Read key
		key := string(data[offset : offset+int(keyLen)])
		offset += int(keyLen)

		// Read value type
		vtype := ValueType(data[offset])
		offset++

		// Read value length
		valueLen := binary.LittleEndian.Uint32(data[offset:])
		offset += 4

		if offset+int(valueLen) > len(data) {
			return nil, errors.New("truncated value data")
		}

		// Read value data
		valueData := data[offset : offset+int(valueLen)]
		offset += int(valueLen)

		// Create value using factory
		if factory != nil {
			value, err := factory(key, vtype, valueData)
			if err != nil {
				return nil, err
			}
			store.values[key] = value
		}
	}

	return store, nil
}

// ToJSON converts to JSON format (alias for Serialize)
func (vs *ValueStore) ToJSON() (string, error) {
	return vs.Serialize()
}

// =========================================================================
// Thread Safety
// =========================================================================

// EnableThreadSafety enables thread-safe operations.
func (vs *ValueStore) EnableThreadSafety() {
	vs.threadSafeEnabled.Store(true)
}

// DisableThreadSafety disables thread-safe operations.
// Use only if you guarantee single-threaded access.
func (vs *ValueStore) DisableThreadSafety() {
	vs.threadSafeEnabled.Store(false)
}

// IsThreadSafe checks if thread safety is enabled.
func (vs *ValueStore) IsThreadSafe() bool {
	return vs.threadSafeEnabled.Load()
}

// =========================================================================
// Statistics
// =========================================================================

// GetReadCount returns the number of read operations.
func (vs *ValueStore) GetReadCount() uint64 {
	return vs.readCount.Load()
}

// GetWriteCount returns the number of write operations.
func (vs *ValueStore) GetWriteCount() uint64 {
	return vs.writeCount.Load()
}

// GetSerializationCount returns the number of serialization operations.
func (vs *ValueStore) GetSerializationCount() uint64 {
	return vs.serializationCount.Load()
}

// ResetStatistics resets all statistics to zero.
func (vs *ValueStore) ResetStatistics() {
	vs.readCount.Store(0)
	vs.writeCount.Store(0)
	vs.serializationCount.Store(0)
}

// =========================================================================
// Iteration Support
// =========================================================================

// Range calls fn for each key-value pair.
// If fn returns false, the iteration stops.
// Thread-safe if EnableThreadSafety was called.
func (vs *ValueStore) Range(fn func(key string, value Value) bool) {
	if vs.threadSafeEnabled.Load() {
		vs.mutex.RLock()
		defer vs.mutex.RUnlock()
	}

	for key, value := range vs.values {
		if !fn(key, value) {
			break
		}
	}
}
