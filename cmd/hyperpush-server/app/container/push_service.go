// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package container

import (
	"fmt"
	"net/http"

	"github.com/euskadi31/go-server"
	"github.com/euskadi31/go-server/response"
	"github.com/euskadi31/go-service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/hyperscale/hyperpush/pkg/hyperpush/push"
)


// Services keys
const (
	PushServerKey = "service.push.server"
)

func init() {
	service.Set(PushServerKey, func(c service.Container) interface{} {
		cfg := c.Get(ConfigKey).(*config.Configuration)

		return push.NewServer(cfg.Push)
	})
}