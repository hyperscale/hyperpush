// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"github.com/hyperscale/hyperpush/pkg/hyperpush/mqtt/packets"
)

// AuthenticationEvent struct
type AuthenticationEvent struct {
	Packet *packets.ConnectPacket
	Client *Client
}

// NewAuthenticationEvent func
func NewAuthenticationEvent(packet *packets.ConnectPacket, client *Client) *AuthenticationEvent {
	return &AuthenticationEvent{
		Packet: packet,
		Client: client,
	}
}
