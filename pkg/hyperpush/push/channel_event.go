// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

// ChannelEventType type def
type ChannelEventType uint8

/**
 * Type of ChannelMessage
 */
const (
	ChannelEventTypeSubscribe ChannelEventType = iota + 1
	ChannelEventTypeUnsubscribe
	ChannelEventTypeUnsubscribeAll
)

// ChannelEvent struct
type ChannelEvent struct {
	Type   ChannelEventType
	Name   string
	Client *Client
}

// NewChannelEvent func
func NewChannelEvent(t ChannelEventType, name string, client *Client) *ChannelEvent {
	return &ChannelEvent{
		Type:   t,
		Name:   name,
		Client: client,
	}
}
