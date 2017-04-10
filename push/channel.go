// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"github.com/hyperscale/hyperpush/protocol"
)

// Channel struct
type Channel struct {
	ID       string
	clients  *ClientManager
	eventsCh chan *protocol.Event
}

// NewChannel constructor
func NewChannel(ID string) *Channel {
	return &Channel{
		ID:       ID,
		clients:  NewClientManager(),
		eventsCh: make(chan *protocol.Event),
	}
}

// Join client to channel
func (c *Channel) Join(client *Client) {
	c.clients.Add(client)

	client.Write(&protocol.Event{
		Type:    "subscribed",
		Channel: c.ID,
	})
}

// Leave client to channel
func (c *Channel) Leave(client *Client) {
	c.clients.Remove(client)

	client.Write(&protocol.Event{
		Type:    "unsubscribed",
		Channel: c.ID,
	})
}

// Write to channel
func (c *Channel) Write(event *protocol.Event) {
	c.eventsCh <- event
}

// Listen channel
func (c *Channel) Listen() {
	for {
		select {
		case e := <-c.eventsCh:
			for _, client := range c.clients.Clients() {
				client.Write(e)
			}
		}
	}
}
