// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"testing"

	"github.com/hyperscale/hyperpush/pkg/hyperpush/mqtt/packets"

	"github.com/stretchr/testify/assert"
)

func TestAutenticationEvent(t *testing.T) {
	packet := packets.NewConnectPacket()

	c := NewClient(nil, nil)
	event := NewAuthenticationEvent(packet, c)

	assert.Equal(t, packet, event.Packet)
	assert.Equal(t, c, event.Client)
}
