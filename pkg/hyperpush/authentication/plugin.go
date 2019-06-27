// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package authentication

// PluginMetadata struct
type PluginMetadata struct {
	Version string
	Name    string
}

// PluginMetaFunc type
type PluginMetaFunc func() PluginMetadata

// PluginInitFunc type
type PluginInitFunc func(cfg string) (Provider, error)

// Plugin struct
type Plugin struct {
	Metadata PluginMetadata
	Provider Provider
}
