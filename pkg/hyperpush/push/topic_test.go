// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"sync"
	"testing"

	"github.com/hyperscale/hyperpush/pkg/hyperpush/mqtt/packets"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTopic(t *testing.T) {
	serverMock := &MockServer{}

	var wgj sync.WaitGroup
	var wgm sync.WaitGroup
	var wgl sync.WaitGroup

	wgj.Add(1)
	wgm.Add(1)
	wgl.Add(1)

	connMock := &transport.MockTransport{}

	connMock.On("Write", packets.NewSubackPacket()).Return(nil).Run(func(fn mock.Arguments) {
		defer wgj.Done()
	}).Once()

	event := packets.NewPublishPacket("test", []byte("bar"))

	connMock.On("Write", event).Return(nil).Run(func(fn mock.Arguments) {
		defer wgm.Done()
	}).Once()

	connMock.On("Write", packets.NewUnsubackPacket()).Return(nil).Run(func(fn mock.Arguments) {
		defer wgl.Done()
	}).Once()

	client := NewClient(connMock, serverMock)

	c := NewTopic("test")
	go c.Listen()

	c.Join(client)
	c.Join(client)

	wgj.Wait()

	assert.Equal(t, 1, c.clients.Size())

	c.Publish(event)

	wgm.Wait()

	c.Leave(client)
	c.Leave(client)

	wgl.Wait()

	assert.Equal(t, 0, c.Size())

	c.Close()
}
