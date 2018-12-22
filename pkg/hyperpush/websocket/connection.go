// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package websocket

import (
	"time"
)

// Connection interface
type Connection interface {
	SetReadLimit(limit int64)
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
	SetWriteDeadline(t time.Time) error
}
