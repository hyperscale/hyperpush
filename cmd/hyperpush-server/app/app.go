// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/euskadi31/go-service"
	"github.com/hyperscale/hyperpush/cmd/hyperpush-server/app/container"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/push"
	"github.com/rs/zerolog/log"
)

// Run push server
func Run() error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	_ = service.Get(container.LoggerKey)
	//router := service.Get(container.RouterKey).(*server.Server)
	server := service.Get(container.PushServerKey).(push.Server)

	log.Info().Msg("Rinning")
	/*
		go func() {
			if err := router.Run(); err != nil {
				log.Error().Err(err).Msg("router.Run")
			}
		}()
	*/
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("server.ListenAndServe")
		}
	}()

	<-sig

	log.Info().Msg("Shutdown")

	//return router.Shutdown()
	return nil
}
