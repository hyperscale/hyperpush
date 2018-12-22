// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"github.com/euskadi31/go-server"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/authentication"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/environment"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/logger"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/push"
)

// Configuration struct
type Configuration struct {
	Environment    environment.Env
	Logger         *logger.Configuration
	Server         *server.Configuration
	Push           *push.Configuration
	Authentication *authentication.Configuration
}

// NewConfiguration constructor
func NewConfiguration() *Configuration {
	return &Configuration{
		Environment:    environment.Dev,
		Logger:         &logger.Configuration{},
		Server:         &server.Configuration{},
		Push:           &push.Configuration{},
		Authentication: &authentication.Configuration{},
	}
}
