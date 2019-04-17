// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"time"

	"github.com/hyperscale/hyperpush/pkg/hyperpush/message"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/metrics"
	"github.com/rs/zerolog/log"
)

// Prefix channel
const (
	PublicChannelPrefix  = "public-"
	PrivateChannelPrefix = "private-"
)

// Channel struct
type Channel struct {
	ID              string
	clients         *ClientPool
	messagesEventCh chan *message.Event
	doneCh          chan bool
}

// NewChannel constructor
func NewChannel(ID string) *Channel {
	return &Channel{
		ID:              ID,
		clients:         NewClientPool(),
		messagesEventCh: make(chan *message.Event, 300),
		doneCh:          make(chan bool),
	}
}

// Join client to channel
func (c *Channel) Join(client *Client) {
	if c.clients.Has(client.ID) {
		log.Debug().Msgf("Client %s already joined the %s channel", client.ID, c.ID)

		return
	}

	c.clients.Add(client)

	client.Write(&message.Event{
		Type:    message.EventTypeSubscribed,
		Channel: c.ID,
	})

	metrics.ChannelClient.With(map[string]string{
		"channel": c.ID,
	}).Set(float64(c.clients.Size()))

	log.Debug().Msgf("Client %s join %s channel.", client.ID, c.ID)
	log.Debug().Msgf("Now %d clients in the %s channel.", c.clients.Size(), c.ID)
}

// Leave client to channel
func (c *Channel) Leave(client *Client) {
	if c.clients.Has(client.ID) == false {
		log.Debug().Msgf("Client %s has not joined the %s channel", client.ID, c.ID)

		return
	}

	c.clients.Del(client.ID)

	client.Write(&message.Event{
		Type:    message.EventTypeUnsubscribed,
		Channel: c.ID,
	})

	metrics.ChannelClient.With(map[string]string{
		"channel": c.ID,
	}).Set(float64(c.clients.Size()))

	log.Debug().Msgf("Client %s leave %s channel.", client.ID, c.ID)
	log.Debug().Msgf("Now %d clients in the %s channel.", c.clients.Size(), c.ID)
}

// Publish to channel
func (c *Channel) Publish(event *message.Event) {
	c.messagesEventCh <- event
}

// Close channel
func (c *Channel) Close() {
	c.doneCh <- true
}

// Size of clients in channel
func (c *Channel) Size() int {
	return c.clients.Size()
}

// Listen channel
func (c *Channel) Listen() {
	for {
		select {
		case <-c.doneCh:
			return
		case e := <-c.messagesEventCh:
			length := len(e.Raw)

			metrics.ChannelMessageReceivedTotal.With(map[string]string{
				"channel": c.ID,
			}).Add(1)

			metrics.ChannelMessageReceivedBytes.With(map[string]string{
				"channel": c.ID,
			}).Add(float64(length))

			clients := c.clients.Clients()

			ts := time.Now()

			for _, client := range clients {
				client.Write(e)

				metrics.ChannelMessageDeliveredTotal.With(map[string]string{
					"channel": c.ID,
				}).Add(1)

				metrics.ChannelMessageDeliveredBytes.With(map[string]string{
					"channel": c.ID,
				}).Add(float64(length))
			}

			metrics.ChannelMessageDeliveredSeconds.With(map[string]string{
				"channel": c.ID,
			}).Observe(time.Since(ts).Seconds())
		}
	}
}
