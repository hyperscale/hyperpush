// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"net"
	"sync"
	"sync/atomic"

	"github.com/hyperscale/hyperpush/pkg/hyperpush/mqtt/packets"
)

const (
	disconnected uint32 = iota
	connecting
	reconnecting
	connected
)

// Client interface
type Client interface {
	// Connect will create a connection to the message broker, by default
	// it will attempt to connect at v3.1.1 and auto retry at v3.1 if that
	// fails
	Connect() Token
	// Disconnect will end the connection with the server, but not before waiting
	// the specified number of milliseconds to wait for existing work to be
	// completed.
	Disconnect(quiesce uint)
	// Publish will publish a message with the specified QoS and content
	// to the specified topic.
	// Returns a token to track delivery of the message to the broker
	Publish(topic string, qos byte, retained bool, payload interface{}) Token
	// Subscribe starts a new subscription. Provide a MessageHandler to be executed when
	// a message is published on the topic provided, or nil for the default handler
	Subscribe(topic string, qos byte, callback MessageHandler) Token
	// SubscribeMultiple starts a new subscription for multiple topics. Provide a MessageHandler to
	// be executed when a message is published on one of the topics provided, or nil for the
	// default handler
	SubscribeMultiple(filters map[string]byte, callback MessageHandler) Token
	// Unsubscribe will end the subscription from each of the topics provided.
	// Messages published to those topics from other clients will no longer be
	// received.
	Unsubscribe(topics ...string) Token
}

// client implements the Client interface
type client struct {
	lastSent        atomic.Value
	lastReceived    atomic.Value
	pingOutstanding int32
	status          uint32
	sync.RWMutex
	messageIds
	conn            net.Conn
	ibound          chan packets.ControlPacket
	obound          chan *PacketAndToken
	oboundP         chan *PacketAndToken
	msgRouter       *router
	stopRouter      chan bool
	incomingPubChan chan *packets.PublishPacket
	errors          chan error
	stop            chan struct{}
	persist         Store
	options         ClientOptions
	workers         sync.WaitGroup
}

// NewClient will create an MQTT v3.1.1 client with all of the options specified
// in the provided ClientOptions. The client must have the Connect method called
// on it before it may be used. This is to make sure resources (such as a net
// connection) are created before the application is actually ready.
func NewClient(o *ClientOptions) Client {
	c := &client{}
	c.options = *o

	if c.options.Store == nil {
		c.options.Store = NewMemoryStore()
	}
	switch c.options.ProtocolVersion {
	case 3, 4:
		c.options.protocolVersionExplicit = true
	case 0x83, 0x84:
		c.options.protocolVersionExplicit = true
	default:
		c.options.ProtocolVersion = 4
		c.options.protocolVersionExplicit = false
	}
	c.persist = c.options.Store
	c.status = disconnected
	c.messageIds = messageIds{index: make(map[uint16]tokenCompletor)}
	c.msgRouter, c.stopRouter = newRouter()
	c.msgRouter.setDefaultHandler(c.options.DefaultPublishHandler)
	if !c.options.AutoReconnect {
		c.options.MessageChannelDepth = 0
	}
	return c
}
