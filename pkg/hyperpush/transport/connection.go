// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package transport

import (
	"net"
	"time"
)

// Connection struct
type Connection struct {
	net.Conn
	t time.Duration
}

// NewConnection constructor
func NewConnection(conn net.Conn, deadline time.Duration) Connection {
	return Connection{
		Conn: conn,
		t:    deadline,
	}
}

func (c Connection) Write(p []byte) (int, error) {
	if err := c.Conn.SetWriteDeadline(time.Now().Add(c.t)); err != nil {
		return 0, err
	}

	return c.Conn.Write(p)
}

func (c Connection) Read(p []byte) (int, error) {
	if err := c.Conn.SetReadDeadline(time.Now().Add(c.t)); err != nil {
		return 0, err
	}

	return c.Conn.Read(p)
}
