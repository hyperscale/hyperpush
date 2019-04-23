// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/mqtt/packets"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/transport"
)

// Client struct
type Client struct {
	ID          string
	UserID      string
	ts          transport.Transport
	server      Server
	isConnected bool
}

// NewClient constructor
func NewClient(ts transport.Transport, server Server) *Client {
	return &Client{
		ID:          uuid.New().String(),
		ts:          ts,
		server:      server,
		isConnected: true,
	}
}

// Write event to client
func (c *Client) Write(event packets.ControlPacket) error {
	return c.ts.Write(event)
}

// IsAuthenticated client
func (c *Client) IsAuthenticated() bool {
	return c.UserID != ""
}

// Close client
func (c *Client) Close() error {
	c.isConnected = false

	c.server.Leave(c)

	if err := c.Write(packets.NewControlPacket(packets.Disconnect)); err != nil {
		return err
	}

	return c.ts.Close()
}

func (c *Client) process(event packets.ControlPacket) error {
	switch evt := event.(type) {
	case *packets.PublishPacket:
		/*
			if c.UserID == "" {
				c.Write(message.NewEventFromErrorCode(message.ErrorCodeUnauthorized))

				break
			}
		*/
		c.server.Publish(evt)

	case *packets.SubscribePacket:
		c.server.JoinTopic(evt, c)

	case *packets.UnsubscribePacket:
		c.server.LeaveTopic(evt, c)

	case *packets.PingreqPacket:
		resp := packets.NewControlPacket(packets.Pingresp)

		c.Write(resp)

	case *packets.ConnectPacket:
		c.ID = evt.ClientIdentifier
		c.server.Authenticate(evt, c)

	default:
		return fmt.Errorf(`Unsupported "%v" event type`, evt)
	}

	return nil
}

// ReadEvent events
func (c *Client) ReadEvent() error {
	event, err := c.ts.Read()
	if err != nil {
		//@TODO: check error type, if bad request return error event
		/*
			log.Error().Err(err).Msg("Bad request")

			c.Write(&message.Event{
				Type: message.EventTypeError,
				Data: json.RawMessage("Bad request"),
			})
		*/
		c.ts.Close()

		return err
	}

	return c.process(event)
}
