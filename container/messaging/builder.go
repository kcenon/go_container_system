/****************************************************************************
BSD 3-Clause License

Copyright (c) 2021, 🍀☀🌕🌥 🌊
All rights reserved.
****************************************************************************/

// Package messaging provides a fluent builder API for constructing ValueContainer instances.
package messaging

import (
	"github.com/kcenon/go_container_system/container/core"
)

// ContainerBuilder provides a fluent API for constructing ValueContainer instances.
// It allows chaining method calls to configure various container properties
// before building the final container.
//
// Example usage:
//
//	container, err := messaging.NewContainerBuilder().
//	    WithSource("client", "1").
//	    WithTarget("server", "main").
//	    WithType("request").
//	    Build()
type ContainerBuilder struct {
	sourceID    string
	sourceSubID string
	targetID    string
	targetSubID string
	messageType string
	values      []core.Value
	threadSafe  bool
}

// NewContainerBuilder creates a new ContainerBuilder instance.
func NewContainerBuilder() *ContainerBuilder {
	return &ContainerBuilder{
		values: make([]core.Value, 0),
	}
}

// WithSource sets the source ID and sub ID for the container.
// Returns the builder for method chaining.
func (b *ContainerBuilder) WithSource(id, subID string) *ContainerBuilder {
	b.sourceID = id
	b.sourceSubID = subID
	return b
}

// WithTarget sets the target ID and sub ID for the container.
// Returns the builder for method chaining.
func (b *ContainerBuilder) WithTarget(id, subID string) *ContainerBuilder {
	b.targetID = id
	b.targetSubID = subID
	return b
}

// WithType sets the message type for the container.
// Returns the builder for method chaining.
func (b *ContainerBuilder) WithType(msgType string) *ContainerBuilder {
	b.messageType = msgType
	return b
}

// WithValues adds values to the container.
// Returns the builder for method chaining.
func (b *ContainerBuilder) WithValues(values ...core.Value) *ContainerBuilder {
	b.values = append(b.values, values...)
	return b
}

// WithThreadSafe enables thread-safe mode for the container.
// Returns the builder for method chaining.
func (b *ContainerBuilder) WithThreadSafe(enabled bool) *ContainerBuilder {
	b.threadSafe = enabled
	return b
}

// Build creates a new ValueContainer with the configured properties.
// Returns the constructed container and any error encountered.
func (b *ContainerBuilder) Build() (*core.ValueContainer, error) {
	container := core.NewValueContainerFull(
		b.sourceID,
		b.sourceSubID,
		b.targetID,
		b.targetSubID,
		b.messageType,
		b.values...,
	)

	if b.threadSafe {
		container.EnableThreadSafe()
	}

	return container, nil
}
