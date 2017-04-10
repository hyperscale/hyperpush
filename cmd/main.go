// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/hyperscale/hyperpush/version"
	"github.com/rs/xlog"
	"log"
)

func main() {
	lc := xlog.Config{
		Fields: xlog.F{
			"app":      "hyperpush",
			"version":  version.Version,
			"revision": version.Revision,
		},
		Level:  xlog.LevelInfo,
		Output: xlog.NewConsoleOutput(),
	}

	// Enable debug level
	/*if Config.Debug {
		lc.Level = xlog.LevelDebug
	}*/

	logger := xlog.New(lc)

	log.SetOutput(logger)
	xlog.SetLogger(logger)

	xlog.Info("Hyperpush")
}
