// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/message"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/websocket"
	"github.com/rs/zerolog/log"
)

const (
	// Maximum message size allowed from peer.
	maxMessageSize = 10240
)

// Client struct
type Client struct {
	ID             string
	UserID         string
	ws             websocket.Connection
	server         Server
	isConnected    bool
	messagesCh     chan *message.Event
	messagesDoneCh chan bool
	Ctx            context.Context
}

// NewClient constructor
func NewClient(ctx context.Context, ws websocket.Connection, server Server) *Client {
	id := uuid.New().String()

	// ctx = log.CtxAddField(ctx, "client_id", id)

	return &Client{
		Ctx:            ctx,
		ID:             id,
		ws:             ws,
		server:         server,
		isConnected:    true,
		messagesCh:     make(chan *message.Event, 250),
		messagesDoneCh: make(chan bool, 1),
	}
}

// Write event to client
func (c *Client) Write(event *message.Event) {
	if c.isConnected {
		c.messagesCh <- event
	}
}

// IsAuthenticated client
func (c *Client) IsAuthenticated() bool {
	return c.UserID != ""
}

func (c *Client) processMessages() {
	defer func() {
		c.isConnected = false
	}()

	for {
		select {
		case e, ok := <-c.messagesCh:
			if err := c.ws.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
				log.Error().Err(err).Msg("")

				continue
			}

			if !ok {
				// The hub closed the channel.
				if err := c.ws.WriteMessage(ws.CloseMessage, []byte{}); err != nil {
					log.Error().Err(err).Msg("")
				}

				return
			}

			msg, err := message.Encode(e)
			if err != nil {
				log.Error().Err(err).Msg("message.Encode")

				continue
			}

			if err := c.ws.WriteMessage(ws.TextMessage, msg); err != nil {
				log.Error().Err(err).Msg("")

				return
			}
		case <-c.messagesDoneCh:
			return
		}
	}
}

// Close client
func (c *Client) Close() error {
	c.isConnected = false

	c.server.Leave(c)

	c.messagesDoneCh <- true

	return c.ws.Close()
}

func (c *Client) process(event *message.Event) {
	switch event.Type {
	case message.EventTypeMessage:
		if c.UserID == "" {
			c.Write(message.NewEventFromErrorCode(message.ErrorCodeUnauthorized))

			break
		}
		c.server.Publish(event)

	case message.EventTypeSubscribe:
		c.server.JoinChannel(event.Channel, c)

	case message.EventTypeUnsubscribe:
		c.server.LeaveChannel(event.Channel, c)

	case message.EventTypePing:
		c.Write(&message.Event{
			Type: message.EventTypePong,
		})

	case message.EventTypeAuthentication:
		c.server.Authenticate(event.Token, c)

	default:
		msg := fmt.Sprintf(`Unsupported "%s" event type`, event.Type)

		log.Error().Msg(msg)

		data, err := json.Marshal(msg)
		if err != nil {
			log.Error().Err(err).Msg("json.Marshal")
		}

		c.Write(&message.Event{
			Type: message.EventTypeError,
			Data: json.RawMessage(data),
		})
	}
}

// Listen func
func (c *Client) Listen() {
	go c.processMessages()

	c.ws.SetReadLimit(maxMessageSize)

	for {
		if !c.isConnected {
			return
		}

		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(
				err,
				ws.CloseAbnormalClosure,
				ws.CloseNoStatusReceived,
				ws.CloseNormalClosure,
				ws.CloseGoingAway,
			) {
				log.Error().Err(err).Msg("ws.ReadMessage")
			}
			return
		}

		event, err := message.Decode(msg)
		if err != nil {
			log.Error().Err(err).Msg("Bad request")

			c.Write(&message.Event{
				Type: message.EventTypeError,
				Data: json.RawMessage("Bad request"),
			})
		} else {
			c.process(event)
		}
	}
}
