// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"errors"
	"testing"

	"github.com/hyperscale/hyperpush/pkg/hyperpush/mqtt/packets"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientWrtie(t *testing.T) {
	serverMock := &MockServer{}

	serverMock.On("Leave", mock.AnythingOfType("*push.Client"))

	event := packets.NewPublishPacket("test", []byte("bar"))

	connMock := &transport.MockTransport{}

	connMock.On("Write", event).Return(nil).Once()

	connMock.On("Write", packets.NewDisconnectPacket()).Return(nil).Once()

	connMock.On("Close").Return(nil)

	client := NewClient(connMock, serverMock)

	err := client.Write(event)
	assert.NoError(t, err)

	client.Close()

	connMock.AssertExpectations(t)

	serverMock.AssertExpectations(t)
}

func TestClientIsAuthenticated(t *testing.T) {
	client := NewClient(nil, nil)

	assert.False(t, client.IsAuthenticated())
}

func TestClientWriteFail(t *testing.T) {
	connMock := &transport.MockTransport{}

	event := packets.NewPublishPacket("test", []byte("bar"))

	connMock.On("Write", event).Return(errors.New("fail"))

	client := NewClient(connMock, nil)

	err := client.Write(event)
	assert.Error(t, err)

	connMock.AssertExpectations(t)
}

func TestClientPing(t *testing.T) {
	event := packets.NewControlPacket(packets.Pingreq).(*packets.PingreqPacket)

	serverMock := &MockServer{}

	serverMock.On("Leave", mock.AnythingOfType("*push.Client"))

	connMock := &transport.MockTransport{}

	connMock.On("Write", packets.NewDisconnectPacket()).Return(nil).Once()

	connMock.On("Write", packets.NewPingrespPacket()).Return(nil).Once()

	connMock.On("Read").Return(event, nil).Once()

	connMock.On("Close").Return(nil)

	client := NewClient(connMock, serverMock)

	err := client.ReadEvent()
	assert.NoError(t, err)

	client.Close()

	connMock.AssertExpectations(t)
	serverMock.AssertExpectations(t)
}

func TestClientSubscribe(t *testing.T) {
	event := packets.NewControlPacket(packets.Subscribe).(*packets.SubscribePacket)
	event.Topics = []string{"test"}

	serverMock := &MockServer{}

	serverMock.On("Leave", mock.AnythingOfType("*push.Client"))

	serverMock.On("JoinTopic", event, mock.AnythingOfType("*push.Client"))

	connMock := &transport.MockTransport{}

	connMock.On("Write", packets.NewDisconnectPacket()).Return(nil).Once()

	connMock.On("Read").Return(event, nil).Once()

	connMock.On("Close").Return(nil)

	client := NewClient(connMock, serverMock)

	err := client.ReadEvent()
	assert.NoError(t, err)

	client.Close()

	connMock.AssertExpectations(t)
	serverMock.AssertExpectations(t)
}

/*
func TestClientMessageWithoutAuthenticate(t *testing.T) {
	serverMock := &MockServer{}

	connMock := &transport.MockTransport{}

	client := NewClient(context.Background(), connMock, serverMock)

	assert.False(t, client.IsAuthenticated())

	client.Listen()

	client.Close()
}
*/
func TestClientBadEventType(t *testing.T) {
	serverMock := &MockServer{}

	serverMock.On("Leave", mock.AnythingOfType("*push.Client")).Once()

	connMock := &transport.MockTransport{}

	connMock.On("Write", packets.NewDisconnectPacket()).Return(nil).Once()

	connMock.On("Read").Return(packets.NewControlPacket(packets.Connack), nil).Once()

	connMock.On("Close").Return(nil).Once()

	client := NewClient(connMock, serverMock)

	err := client.ReadEvent()
	assert.Error(t, err)

	client.Close()

	connMock.AssertExpectations(t)
	serverMock.AssertExpectations(t)
}

func TestProcessMessageWriteFail(t *testing.T) {
	serverMock := &MockServer{}

	event := packets.NewPublishPacket("test", []byte("bar"))

	connMock := &transport.MockTransport{}

	connMock.On("Write", event).Return(errors.New("fail")).Once()

	client := NewClient(connMock, serverMock)

	client.Write(event)

	connMock.AssertExpectations(t)
	serverMock.AssertExpectations(t)
}
