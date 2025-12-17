/****************************************************************************
BSD 3-Clause License

Copyright (c) 2021, 🍀☀🌕🌥 🌊
All rights reserved.
****************************************************************************/

package di

import (
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for container system dependencies.
// Include this set in your wire.Build() call to automatically wire
// all container-related dependencies.
//
// Example:
//
//	func InitializeService() (*Service, error) {
//	    wire.Build(
//	        di.ProviderSet,
//	        NewService,
//	    )
//	    return nil, nil
//	}
var ProviderSet = wire.NewSet(
	NewContainerFactory,
	wire.Bind(new(ContainerFactory), new(*DefaultContainerFactory)),
)
