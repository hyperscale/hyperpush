// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"

	ws "github.com/gorilla/websocket"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/message"
	"github.com/stretchr/testify/assert"
)

func TestClientWrtie(t *testing.T) {
	server := &ServerMock{}

	conn := &WebSocketConnection{
		writeMessageAssertType: "write",
		writeAssertFunc: func(messageType int, data []byte) error {

			assert.Equal(t, ws.TextMessage, messageType)

			event, err := message.Decode(data)
			assert.NoError(t, err)

			assert.Equal(t, message.EventTypeMessage, event.Type)
			assert.Equal(t, "test", event.Channel)
			assert.Equal(t, "foo", event.Name)
			assert.Equal(t, json.RawMessage(`"bar"`), event.Data)

			return nil
		},
	}

	client := NewClient(context.Background(), conn, server)

	go client.processMessages()

	client.Write(&message.Event{
		Type:    message.EventTypeMessage,
		Channel: "test",
		Name:    "foo",
		Data:    json.RawMessage(`"bar"`),
	})

	client.Close()
}

func TestClientWriteFail(t *testing.T) {
	conn := &WebSocketConnection{
		writeMessageAssertType: "write",
		writeAssertFunc: func(messageType int, data []byte) error {
			assert.Equal(t, ws.TextMessage, messageType)

			event, err := message.Decode(data)
			assert.NoError(t, err)

			assert.Equal(t, message.EventTypeMessage, event.Type)
			assert.Equal(t, "test", event.Channel)
			assert.Equal(t, json.RawMessage(`"bar"`), event.Data)

			return errors.New("test error")
		},
	}

	client := NewClient(context.Background(), conn, nil)

	client.Write(&message.Event{
		Type:    message.EventTypeMessage,
		Channel: "test",
		Name:    "foo",
		Data:    json.RawMessage(`"bar"`),
	})
}

func TestClientPing(t *testing.T) {
	server := &ServerMock{}

	readCount := 0

	conn := &WebSocketConnection{
		writeMessageAssertType: "write",
		writeAssertFunc: func(messageType int, data []byte) error {
			assert.Equal(t, ws.TextMessage, messageType)

			event, err := message.Decode(data)
			assert.NoError(t, err)

			assert.Equal(t, "pong", event.Type)

			return nil
		},
		readMessageAssertType: "ping",
		readMessageAssertFunc: func() (messageType int, p []byte, err error) {
			if readCount == 1 {
				return 1, []byte(""), &ws.CloseError{Code: ws.CloseUnsupportedData}
			}

			readCount++

			return 1, []byte(`{"type":"ping"}`), nil
		},
	}

	client := NewClient(context.Background(), conn, server)

	client.Listen()

	client.Close()
}

func TestClientSubscribe(t *testing.T) {
	readCount := 0

	server := &ServerMock{}

	conn := &WebSocketConnection{
		writeMessageAssertType: "write",
		writeAssertFunc: func(messageType int, data []byte) error {
			assert.Equal(t, ws.TextMessage, messageType)

			event, err := message.Decode(data)
			assert.NoError(t, err)

			switch readCount {
			case 1:
				assert.Equal(t, "subscribed", event.Type)
				assert.Equal(t, "test", event.Channel)
			case 2:
				assert.Equal(t, "unsubscribed", event.Type)
				assert.Equal(t, "test", event.Channel)
			}

			return nil
		},
		readMessageAssertFunc: func() (messageType int, p []byte, err error) {
			defer func() {
				readCount++
			}()

			switch readCount {
			case 0:
				return 1, []byte(`{"type":"subscribe","channel":"test"}`), nil
			case 1:
				return 1, []byte(`{"type":"unsubscribe","channel":"test"}`), nil
			default:
				return 1, []byte(""), &ws.CloseError{Code: ws.CloseUnsupportedData}
			}
		},
	}

	client := NewClient(context.Background(), conn, server)

	client.Listen()

	client.Close()
}

func TestClientMessageWithoutAuthenticate(t *testing.T) {
	server := &ServerMock{}

	readCount := 0

	conn := &WebSocketConnection{
		writeMessageAssertType: "write",
		writeAssertFunc: func(messageType int, data []byte) error {
			assert.Equal(t, ws.TextMessage, messageType)

			event, err := message.Decode(data)
			assert.NoError(t, err)

			assert.Equal(t, "error", event.Type)
			assert.Equal(t, int(message.ErrorCodeUnauthorized), event.Code)
			assert.Equal(t, "Unauthorized", event.Message)

			return nil
		},
		readMessageAssertFunc: func() (messageType int, p []byte, err error) {
			if readCount == 1 {
				return 1, []byte(""), &ws.CloseError{Code: ws.CloseUnsupportedData}
			}

			readCount++

			return 1, []byte(`{"type":"message","channel":"test","name":"foo","data":"bar"}`), nil
		},
	}

	client := NewClient(context.Background(), conn, server)

	assert.False(t, client.IsAuthenticated())

	client.Listen()

	client.Close()
}

func TestClientBadEventType(t *testing.T) {
	readCount := 0

	server := &ServerMock{}

	conn := &WebSocketConnection{
		writeMessageAssertType: "write",
		writeAssertFunc: func(messageType int, data []byte) error {
			if messageType == ws.CloseMessage {
				return nil
			}

			assert.Equal(t, ws.TextMessage, messageType)

			event, err := message.Decode(data)
			assert.NoError(t, err)

			assert.Equal(t, "error", event.Type)
			assert.Equal(t, json.RawMessage(`"Unsupported \"bad\" event type"`), event.Data)

			return nil
		},
		readMessageAssertType: "bad",
		readMessageAssertFunc: func() (messageType int, p []byte, err error) {
			if readCount == 1 {
				return 1, []byte(""), &ws.CloseError{Code: ws.CloseUnsupportedData}
			}

			readCount++

			return 1, []byte(`{"type":"bad","channel":"test"}`), nil
		},
	}

	client := NewClient(context.Background(), conn, server)

	client.Listen()

	client.Close()
}

func TestProcessMessageWriteFail(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)

	server := &ServerMock{}

	conn := &WebSocketConnection{
		writeMessageAssertType: "write",
		writeAssertFunc: func(messageType int, data []byte) error {
			defer wg.Done()

			return errors.New("test")
		},
	}

	client := NewClient(context.Background(), conn, server)

	go client.processMessages()

	client.Write(&message.Event{
		Type:    message.EventTypeMessage,
		Channel: "test",
		Name:    "foo",
		Data:    json.RawMessage(`"bar"`),
	})

	wg.Wait()
}

func TestProcessMessageClosedChannel(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)

	server := &ServerMock{}

	conn := &WebSocketConnection{
		writeMessageAssertType: "write",
		writeAssertFunc: func(messageType int, data []byte) error {
			defer wg.Done()

			assert.Equal(t, ws.CloseMessage, messageType)

			return nil
		},
	}

	client := NewClient(context.Background(), conn, server)

	go client.processMessages()

	close(client.messagesCh)

	wg.Wait()
}
