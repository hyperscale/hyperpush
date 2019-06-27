// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package transport

import (
	"io"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/mqtt/packets"
)

var _ Transport = (*WebSocket)(nil)

// WebSocket struct
type WebSocket struct {
	conn io.ReadWriteCloser
}

// NewWebSocket transport
func NewWebSocket(conn io.ReadWriteCloser) *WebSocket {
	return &WebSocket{
		conn,
	}
}

// Read message from WebSocket
func (t *WebSocket) Read() (packets.ControlPacket, error) {
	h, r, err := wsutil.NextReader(t.conn, ws.StateServerSide)
	if err != nil {
		return nil, err
	}

	if h.OpCode.IsControl() {
		return nil, wsutil.ControlFrameHandler(t.conn, ws.StateServerSide)(h, r)
	}

	return packets.ReadPacket(r)
}

// Write message to WebSocket
func (t *WebSocket) Write(event packets.ControlPacket) error {
	w := wsutil.NewWriter(t.conn, ws.StateServerSide, ws.OpText)

	if err := event.Write(w); err != nil {
		return err
	}

	return w.Flush()
}

// Close socket
func (t *WebSocket) Close() error {
	return t.conn.Close()
}
