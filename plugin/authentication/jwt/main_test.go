// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPluginMetadata(t *testing.T) {
	meta := PluginMetadata()

	assert.Equal(t, "jwt", meta.Name)
	assert.Equal(t, "1.0.0", meta.Version)
}
