// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"time"

	"github.com/hyperscale/hyperpush/pkg/hyperpush/metrics"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/mqtt/packets"
	"github.com/rs/zerolog/log"
)

// Topic struct
type Topic struct {
	ID              TopicID
	clients         *ClientPool
	messagesEventCh chan *packets.PublishPacket
	doneCh          chan bool
}

// NewTopic constructor
func NewTopic(ID string) *Topic {
	return &Topic{
		ID:              TopicID(ID),
		clients:         NewClientPool(),
		messagesEventCh: make(chan *packets.PublishPacket, 300),
		doneCh:          make(chan bool),
	}
}

// Join client to topic
func (c *Topic) Join(client *Client) {
	if c.clients.Has(client.ID) {
		log.Debug().Msgf("Client %s already joined the %s topic", client.ID, c.ID)

		return
	}

	c.clients.Add(client)

	client.Write(packets.NewControlPacket(packets.Suback))

	metrics.TopicClient.With(map[string]string{
		"topic": c.ID.String(),
	}).Set(float64(c.clients.Size()))

	log.Debug().Msgf("Client %s join %s topic.", client.ID, c.ID)
	log.Debug().Msgf("Now %d clients in the %s topic.", c.clients.Size(), c.ID)
}

// Leave client to topic
func (c *Topic) Leave(client *Client) {
	if c.clients.Has(client.ID) == false {
		log.Debug().Msgf("Client %s has not joined the %s topic", client.ID, c.ID)

		return
	}

	c.clients.Del(client.ID)

	client.Write(packets.NewControlPacket(packets.Unsuback))

	metrics.TopicClient.With(map[string]string{
		"topic": c.ID.String(),
	}).Set(float64(c.clients.Size()))

	log.Debug().Msgf("Client %s leave %s topic.", client.ID, c.ID)
	log.Debug().Msgf("Now %d clients in the %s topic.", c.clients.Size(), c.ID)
}

// Publish to topic
func (c *Topic) Publish(event *packets.PublishPacket) {
	c.messagesEventCh <- event
}

// Close topic
func (c *Topic) Close() {
	c.doneCh <- true
}

// Size of clients in topic
func (c *Topic) Size() int {
	return c.clients.Size()
}

// Listen topic
func (c *Topic) Listen() {
	for {
		select {
		case <-c.doneCh:
			return
		case e := <-c.messagesEventCh:
			//length := len(e.Raw)

			metrics.TopicMessageReceivedTotal.With(map[string]string{
				"topic": c.ID.String(),
			}).Add(1)
			/*
				metrics.TopicMessageReceivedBytes.With(map[string]string{
					"topic": c.ID,
				}).Add(float64(length))
			*/
			clients := c.clients.Clients()

			ts := time.Now()

			for _, client := range clients {
				client.Write(e)

				metrics.TopicMessageDeliveredTotal.With(map[string]string{
					"topic": c.ID.String(),
				}).Add(1)
				/*
					metrics.TopicMessageDeliveredBytes.With(map[string]string{
						"topic": c.ID,
					}).Add(float64(length))
				*/
			}

			metrics.TopicMessageDeliveredSeconds.With(map[string]string{
				"topic": c.ID.String(),
			}).Observe(time.Since(ts).Seconds())
		}
	}
}
