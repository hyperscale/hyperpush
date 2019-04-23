// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package transport

import "github.com/hyperscale/hyperpush/pkg/hyperpush/mqtt/packets"

// Transport interface
//go:generate mockery -case=underscore -inpkg -name=Transport
type Transport interface {
	// Wrtie message event to interface
	Write(event packets.ControlPacket) error

	// Read message event from interface
	Read() (packets.ControlPacket, error)

	// Close socket client
	Close() error
}
