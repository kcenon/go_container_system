/****************************************************************************
BSD 3-Clause License

Copyright (c) 2021, 🍀☀🌕🌥 🌊
All rights reserved.
****************************************************************************/

// Package di provides dependency injection support for the container system.
// It defines standard interfaces and providers for integration with Go DI frameworks
// such as Google Wire and Uber Dig.
//
// Example usage with Google Wire:
//
//	// wire.go
//	//go:build wireinject
//	// +build wireinject
//
//	package main
//
//	import (
//	    "github.com/google/wire"
//	    "github.com/kcenon/go_container_system/container/di"
//	)
//
//	func InitializeApp() (*App, error) {
//	    wire.Build(di.ProviderSet, NewApp)
//	    return nil, nil
//	}
//
// Example usage with Uber Dig:
//
//	container := dig.New()
//	container.Provide(di.NewContainerFactory)
package di

import (
	"github.com/kcenon/go_container_system/container/core"
	"github.com/kcenon/go_container_system/container/messaging"
)

// ContainerFactory defines the interface for creating ValueContainer instances.
// This interface allows for easy mocking in tests and provides a standard
// abstraction for container creation across the application.
type ContainerFactory interface {
	// NewContainer creates a new empty ValueContainer.
	NewContainer() *core.ValueContainer

	// NewContainerWithType creates a ValueContainer with the specified message type.
	NewContainerWithType(messageType string) *core.ValueContainer

	// NewContainerWithTarget creates a ValueContainer with target information.
	NewContainerWithTarget(targetID, targetSubID, messageType string) *core.ValueContainer

	// NewContainerFull creates a ValueContainer with full header information.
	NewContainerFull(sourceID, sourceSubID, targetID, targetSubID, messageType string) *core.ValueContainer

	// NewBuilder creates a new ContainerBuilder for fluent container construction.
	NewBuilder() *messaging.ContainerBuilder
}

// DefaultContainerFactory is the default implementation of ContainerFactory.
// It creates ValueContainer instances using the standard constructors from the core package.
type DefaultContainerFactory struct{}

// NewContainerFactory creates a new ContainerFactory instance.
// This is the provider function for dependency injection frameworks.
func NewContainerFactory() ContainerFactory {
	return &DefaultContainerFactory{}
}

// NewContainer creates a new empty ValueContainer.
func (f *DefaultContainerFactory) NewContainer() *core.ValueContainer {
	return core.NewValueContainer()
}

// NewContainerWithType creates a ValueContainer with the specified message type.
func (f *DefaultContainerFactory) NewContainerWithType(messageType string) *core.ValueContainer {
	return core.NewValueContainerWithType(messageType)
}

// NewContainerWithTarget creates a ValueContainer with target information.
func (f *DefaultContainerFactory) NewContainerWithTarget(targetID, targetSubID, messageType string) *core.ValueContainer {
	return core.NewValueContainerWithTarget(targetID, targetSubID, messageType)
}

// NewContainerFull creates a ValueContainer with full header information.
func (f *DefaultContainerFactory) NewContainerFull(sourceID, sourceSubID, targetID, targetSubID, messageType string) *core.ValueContainer {
	return core.NewValueContainerFull(sourceID, sourceSubID, targetID, targetSubID, messageType)
}

// NewBuilder creates a new ContainerBuilder for fluent container construction.
func (f *DefaultContainerFactory) NewBuilder() *messaging.ContainerBuilder {
	return messaging.NewContainerBuilder()
}
