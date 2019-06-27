// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/go-yaml/yaml"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/authentication"
	"github.com/pkg/errors"
)

// PluginMetadata entrypoint
func PluginMetadata() authentication.PluginMetadata {
	return authentication.PluginMetadata{
		Name:    "jwt",
		Version: "1.0.0",
	}
}

// PluginInit entrypoint
func PluginInit(cfg string) (authentication.Provider, error) {
	cfg = filepath.Clean(cfg)

	if cfg == "" {
		return nil, fmt.Errorf(`invalid "%s" config file`, cfg)
	}

	data, err := ioutil.ReadFile(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, `cannot read "%s" config file`, cfg)
	}

	config := &Configuration{}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, errors.Wrapf(err, `cannot parse "%s" config file`, cfg)
	}

	return &provider{
		config: config,
	}, nil
}

func main() {}
