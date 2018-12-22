// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package container

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/euskadi31/go-service"
	"github.com/hyperscale/hyperpush/cmd/hyperpush-server/app/config"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/environment"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Services keys
const (
	ConfigKey = "service.config"
)

const name = "hyperpush"

func init() {
	service.Set(ConfigKey, func(c service.Container) interface{} {
		cfg := config.NewConfiguration()

		var cfgFile string
		var env string

		cmd := c.Get(FlagKey).(*flag.FlagSet)

		cmd.StringVar(&cfgFile, "config", "", "config file (default is $HOME/config.yaml)")
		cmd.StringVar(&env, "env", "dev", "Environment {dev, preprod, prod}")

		// Ignore errors; cmd is set for ExitOnError.
		// nolint:gosec
		_ = cmd.Parse(os.Args[1:])

		options := viper.New()

		options.SetDefault("environment", "dev")
		options.SetDefault("logger.level", "debug")
		options.SetDefault("logger.prefix", name)
		options.SetDefault("server.host", "")
		options.SetDefault("server.port", 8080)
		options.SetDefault("server.read_header_timeout", 10*time.Millisecond)
		options.SetDefault("server.read_timeout", 10*time.Second)
		options.SetDefault("server.shutdown_timeout", 10*time.Second)
		options.SetDefault("server.write_timeout", 10*time.Second)

		options.SetConfigName("config") // name of config file (without extension)

		options.AddConfigPath("/etc/" + name + "/")   // path to look for the config file in
		options.AddConfigPath("$HOME/." + name + "/") // call multiple times to add many search paths
		options.AddConfigPath(".")

		if cfgFile != "" { // enable ability to specify config file via flag
			options.SetConfigFile(cfgFile)
		}

		if environment := os.Getenv("ENV"); environment != "" {
			if err := os.Setenv("HYPERPUSH_ENVIRONMENT", environment); err != nil {
				log.Error().Err(err).Msg("Setenv failed")
			}
		}

		options.SetEnvPrefix("HYPERPUSH")
		options.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		options.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := options.ReadInConfig(); err == nil {
			log.Info().Msgf("Using config file: %s", options.ConfigFileUsed())
		}

		if err := options.Unmarshal(cfg); err != nil {
			log.Fatal().Err(err).Msg(ConfigKey)
		}

		cfg.Environment = environment.FromString(env)

		return cfg
	})
}
