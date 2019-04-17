// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package container

import (
	"github.com/euskadi31/go-eventemitter"
	"github.com/euskadi31/go-service"
)

// Services keys
const (
	EventEmitterKey = "service.eventemitter"
)

func init() {
	service.Set(EventEmitterKey, func(c service.Container) interface{} {
		return eventemitter.New() // eventemitter.EventEmitter
	})
}
