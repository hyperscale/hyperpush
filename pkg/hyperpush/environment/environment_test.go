// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package environment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentTypeString(t *testing.T) {
	env := Prod

	assert.Equal(t, "prod", env.String())
}

func TestEnvironmentFromString(t *testing.T) {
	assert.Equal(t, Prod, FromString("prod"))
	assert.Equal(t, PreProd, FromString("preprod"))
	assert.Equal(t, Dev, FromString("dev"))
	assert.Equal(t, Dev, FromString("bad"))
}
