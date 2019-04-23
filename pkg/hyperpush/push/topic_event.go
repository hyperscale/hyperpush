// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import "github.com/hyperscale/hyperpush/pkg/hyperpush/mqtt/packets"

// TopicEventType type def
type TopicEventType uint8

/**
 * Type of TopicEvent
 */
const (
	TopicEventTypeSubscribe TopicEventType = iota + 1
	TopicEventTypeUnsubscribe
	TopicEventTypeUnsubscribeAll
)

// TopicEvent struct
type TopicEvent struct {
	Type    TopicEventType
	Name    string
	Client  *Client
	Details packets.Details
}

// NewTopicEvent func
func NewTopicEvent(t TopicEventType, name string, client *Client, details packets.Details) *TopicEvent {
	return &TopicEvent{
		Type:    t,
		Name:    name,
		Client:  client,
		Details: details,
	}
}
