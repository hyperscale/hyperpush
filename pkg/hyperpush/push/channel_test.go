// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/hyperscale/hyperpush/pkg/hyperpush/message"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/websocket"
	"github.com/stretchr/testify/assert"
)

func TestChannel(t *testing.T) {
	serverMock := &MockServer{}

	var wgj sync.WaitGroup
	var wgm sync.WaitGroup
	var wgl sync.WaitGroup

	wgj.Add(1)
	wgm.Add(1)
	wgl.Add(1)

	connMock := &websocket.MockConnection{}

	connMock.On("SetWriteDeadline", mock.Anything).Return(nil)
	connMock.On("WriteMessage", 1, []byte(`{"type":"subscribed","channel":"test"}`)).Return(nil).Run(func(fn mock.Arguments) {
		defer wgj.Done()
	})

	connMock.On("WriteMessage", 1, []byte(`{"type":"message","channel":"test","name":"foo","data":"bar"}`)).Return(nil).Run(func(fn mock.Arguments) {
		defer wgm.Done()
	})

	connMock.On("WriteMessage", 1, []byte(`{"type":"unsubscribed","channel":"test"}`)).Return(nil).Run(func(fn mock.Arguments) {
		defer wgl.Done()
	})

	client := NewClient(context.Background(), connMock, serverMock)
	go client.processMessages()

	c := NewChannel("test")
	go c.Listen()

	c.Join(client)
	c.Join(client)

	wgj.Wait()

	assert.Equal(t, 1, c.clients.Size())

	c.Publish(&message.Event{
		Type:    message.EventTypeMessage,
		Channel: "test",
		Name:    "foo",
		Data:    json.RawMessage(`"bar"`),
	})

	wgm.Wait()

	c.Leave(client)
	c.Leave(client)

	wgl.Wait()

	assert.Equal(t, 0, c.Size())
}
