// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package container

import (
	"github.com/euskadi31/go-eventemitter"
	"github.com/euskadi31/go-service"
	"github.com/hyperscale/hyperpush/cmd/hyperpush-server/app/config"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/push"
)

// Services keys
const (
	PushServerKey = "service.push.server"
)

func init() {
	service.Set(PushServerKey, func(c service.Container) interface{} {
		cfg := c.Get(ConfigKey).(*config.Configuration)
		emitter := c.Get(EventEmitterKey).(eventemitter.EventEmitter)

		return push.NewServer(cfg.Push, emitter)
	})
}
