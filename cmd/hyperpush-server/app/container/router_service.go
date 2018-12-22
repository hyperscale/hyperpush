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
	//"github.com/hyperscale/hyperpush/cmd/hyperpush-server/app/config"
	hlogger "github.com/hyperscale/hyperpush/pkg/hyperpush/logger"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/push"
)


// Services keys
const (
	RouterKey = "service.http.router"
)

func init() {
	service.Set(RouterKey, func(c service.Container) interface{} {
		logger := c.Get(LoggerKey).(zerolog.Logger)
		//cfg := c.Get(ConfigKey).(*config.Configuration)

		pushHandler := c.Get(PushServerKey).(push.Server)

		router := server.New(&server.Configuration{
			HTTP: &server.HTTPConfiguration{
				Port: 12456,
			},
			HTTPS: &server.HTTPSConfiguration{
				Port:     12457,
				CertFile: "./testdata/server.crt",
				KeyFile:  "./testdata/server.key",
			},
			Profiling:   true,
			Metrics:     true,
			HealthCheck: true,
		})

		router.Use(hlog.NewHandler(logger))
		router.Use(hlog.AccessHandler(hlogger.Handler))
		router.Use(hlog.RemoteAddrHandler("ip"))
		router.Use(hlog.UserAgentHandler("user_agent"))
		router.Use(hlog.RefererHandler("referer"))
		router.Use(hlog.RequestIDHandler("req_id", "Request-Id"))

		router.EnableCors()

		router.SetNotFoundFunc(func(w http.ResponseWriter, r *http.Request) {
			response.Encode(w, r, http.StatusNotFound, map[string]interface{}{
				"error": map[string]interface{}{
					"message": fmt.Sprintf("%s %s not found", r.Method, r.URL.Path),
				},
			})
		})

		router.Handle("/ws", pushHandler).Methods(http.MethodGet)

		return router
	})
}