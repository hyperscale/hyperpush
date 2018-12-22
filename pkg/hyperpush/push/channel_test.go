// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	ws "github.com/gorilla/websocket"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/message"
	"github.com/stretchr/testify/assert"
)

func TestChannel(t *testing.T) {
	server := &ServerMock{}

	var wgj sync.WaitGroup
	var wgm sync.WaitGroup

	wgj.Add(1)
	wgm.Add(1)

	conn := &WebSocketConnection{
		writeMessageAssertType: "join",
		joinAssertFunc: func(messageType int, data []byte) error {
			defer wgj.Done()
			assert.Equal(t, ws.TextMessage, messageType)

			event, err := message.Decode(data)
			assert.NoError(t, err)

			assert.Equal(t, message.EventTypeSubscribed, event.Type)
			assert.Equal(t, "test", event.Channel)

			return nil
		},
		leaveAssertFunc: func(messageType int, data []byte) error {
			assert.Equal(t, ws.TextMessage, messageType)

			event, err := message.Decode(data)
			assert.NoError(t, err)

			assert.Equal(t, message.EventTypeUnsubscribed, event.Type)
			assert.Equal(t, "test", event.Channel)

			return nil
		},
		writeAssertFunc: func(messageType int, data []byte) error {
			defer wgm.Done()
			assert.Equal(t, ws.TextMessage, messageType)

			event, err := message.Decode(data)
			assert.NoError(t, err)

			assert.Equal(t, message.EventTypeMessage, event.Type)
			assert.Equal(t, "test", event.Channel)
			assert.Equal(t, json.RawMessage(`"bar"`), event.Data)

			return nil
		},
	}

	client := NewClient(context.Background(), conn, server)
	go client.processMessages()

	c := NewChannel("test")
	go c.Listen()

	c.Join(client)
	c.Join(client)

	wgj.Wait()

	assert.Equal(t, 1, c.clients.Size())

	conn.writeMessageAssertType = "write"

	c.Publish(&message.Event{
		Type:    message.EventTypeMessage,
		Channel: "test",
		Name:    "foo",
		Data:    json.RawMessage(`"bar"`),
	})

	wgm.Wait()

	conn.writeMessageAssertType = "leave"

	c.Leave(client)
	c.Leave(client)

	assert.Equal(t, 0, c.Size())
}
