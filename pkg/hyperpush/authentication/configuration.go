// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package authentication

import (
	"fmt"
	"plugin"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Configuration struct
type Configuration struct {
	Plugin string
	Config string
}

// Load plugin
func (c *Configuration) Load() (*Plugin, error) {
	pluginFile := fmt.Sprintf("authentication.%s.so", c.Plugin)

	p, err := plugin.Open(pluginFile)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot load %s plugin", pluginFile)
	}

	metaSymbol, err := p.Lookup("PluginMetadata")
	if err != nil {
		return nil, errors.Wrapf(err, "authentication plugin %s is not valid", pluginFile)
	}

	initSymbol, err := p.Lookup("PluginInit")
	if err != nil {
		return nil, errors.Wrapf(err, "authentication plugin %s is not valid", pluginFile)
	}

	meta := metaSymbol.(PluginMetaFunc)()

	log.Debug().Msgf("Load %s %s authentication plugin", meta.Name, meta.Version)

	provider, err := initSymbol.(PluginInitFunc)(c.Config)
	if err != nil {
		return nil, errors.Wrapf(err, "init %s plugin failed", c.Plugin)
	}

	return &Plugin{
		Metadata: meta,
		Provider: provider,
	}, nil
}
