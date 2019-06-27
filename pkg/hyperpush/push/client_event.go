// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

// ClientEventType type def
type ClientEventType uint8

/**
 * Type of ClientMessage
 */
const (
	ClientEventTypeJoin ClientEventType = iota + 1
	ClientEventTypeLeave
)

// ClientEvent struct
type ClientEvent struct {
	Type   ClientEventType
	Client *Client
}

// NewClientEvent func
func NewClientEvent(t ClientEventType, client *Client) *ClientEvent {
	return &ClientEvent{
		Type:   t,
		Client: client,
	}
}
