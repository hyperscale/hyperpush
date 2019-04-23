// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"net"
	"net/http"
	"time"

	"github.com/euskadi31/go-eventemitter"
	"github.com/gobwas/ws"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/authentication"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/metrics"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/mqtt/packets"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/pool"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/transport"
	"github.com/mailru/easygo/netpoll"
	"github.com/rs/zerolog/log"
)

// ErrorHandler func
type ErrorHandler func(w http.ResponseWriter, code int, err error)

// Server interface
//go:generate mockery -case=underscore -inpkg -name=Server
type Server interface {
	Authenticate(packet *packets.ConnectPacket, client *Client)
	SetAuthenticationProvider(provider authentication.Provider)
	JoinTopic(packet *packets.SubscribePacket, client *Client)
	LeaveTopic(packet *packets.UnsubscribePacket, client *Client)
	Leave(client *Client)
	Publish(message *packets.PublishPacket)
	ListenAndServe() error
}

type server struct {
	config              *Configuration
	ErrorHandler        ErrorHandler
	clients             *ClientPool
	clientsEventCh      chan *ClientEvent
	users               UserPool
	topics              *TopicPool
	topicsEventCh       chan *TopicEvent
	authentication      authentication.Provider
	authenticateEventCh chan *AuthenticationEvent
	messagesCh          chan *packets.PublishPacket
	emitter             eventemitter.EventEmitter
	poller              netpoll.Poller
	pool                *pool.Pool
}

// NewServer server
func NewServer(config *Configuration, emitter eventemitter.EventEmitter) (Server, error) {
	poller, err := netpoll.New(nil)
	if err != nil {
		return nil, err
	}

	p, err := pool.NewPool(config.ConnectionWorkerSize, config.ConnectionQueueSize, 1)
	if err != nil {
		return nil, err
	}

	return &server{
		config: config,
		ErrorHandler: func(w http.ResponseWriter, code int, err error) {
			log.Error().Err(err).Msg("error")

			http.Error(w, err.Error(), code)
		},
		clients:             NewClientPool(),
		users:               NewUserPool(),
		clientsEventCh:      make(chan *ClientEvent, config.ClientQueueSize),
		topics:              NewTopicPool(),
		topicsEventCh:       make(chan *TopicEvent, config.TopicQueueSize),
		authenticateEventCh: make(chan *AuthenticationEvent, config.AuthenticationQueueSize),
		messagesCh:          make(chan *packets.PublishPacket, config.MessageQueueSize),
		emitter:             emitter,
		poller:              poller,
		pool:                p,
	}, nil
}

// SetAuthenticationProvider to server
func (p *server) SetAuthenticationProvider(provider authentication.Provider) {
	p.authentication = provider
}

// Authenticate user
func (p *server) Authenticate(packet *packets.ConnectPacket, client *Client) {
	event := NewAuthenticationEvent(packet, client)

	p.authenticateEventCh <- event
}

// Publish message to server
func (p *server) Publish(message *packets.PublishPacket) {
	p.messagesCh <- message
}

// JoinTopic added client to topic
func (p *server) JoinTopic(packet *packets.SubscribePacket, client *Client) {
	for _, topic := range packet.Topics {
		event := NewTopicEvent(TopicEventTypeSubscribe, topic, client, packet.Details())

		p.topicsEventCh <- event
	}
}

// LeaveTopic removed client to topic
func (p *server) LeaveTopic(packet *packets.UnsubscribePacket, client *Client) {
	for _, topic := range packet.Topics {
		event := NewTopicEvent(TopicEventTypeUnsubscribe, topic, client, packet.Details())

		p.topicsEventCh <- event
	}
}

// Leave removed client to all topic
func (p *server) Leave(client *Client) {
	event := NewClientEvent(ClientEventTypeLeave, client)

	p.clientsEventCh <- event
}

func (p *server) dispatchMessage(event *packets.PublishPacket) {
	/*if event.User != "" {
		// send message to user on all client
		if clients, ok := p.users.Get(event.User); ok {
			for _, client := range clients {
				p.pool.Schedule(func() {
					client.Write(event)
				})
			}
		}
	} else if event.Topic != "" {*/
	// broadcast message to topic
	if topic, ok := p.topics.Get(event.TopicName); ok {
		topic.Publish(event)
	}
	//}
}

func (p *server) processMessage() {
	for {
		select {
		case e := <-p.messagesCh:
			p.dispatchMessage(e)
		}
	}
}

func (p *server) processAuthenticateEvent() {
	for {
		select {
		case e := <-p.authenticateEventCh:
			// authenticate user here
			user, err := p.authentication.Authenticate(e.Packet)
			if err != nil {
				e.Client.Write(packets.NewConnackPacketFromErr(err))
			} else {
				e.Client.UserID = user.ID

				if p.users.HasClient(user.ID, e.Client.ID) {
					log.Debug().Msgf("Client %s already autenticated", e.Client.ID)

					// @see:
					//e.Client.Close()

					break
				}

				p.users.Add(user.ID, e.Client)

				metrics.ClientAuthenticate.WithLabelValues().Set(float64(p.users.Size()))

				e.Client.Write(packets.NewConnackPacket())

				log.Debug().Msgf("Authenticate user %s for client %s", user.ID, e.Client.ID)
			}
		}
	}
}

func (p *server) processClientEvent() {
	for {
		select {
		case e := <-p.clientsEventCh:
			switch e.Type {
			case ClientEventTypeJoin:
				if p.clients.Has(e.Client.ID) {
					log.Debug().Msgf("Client %s already connected.", e.Client.ID)

					break
				}

				p.clients.Add(e.Client)

				metrics.ClientLive.WithLabelValues().Set(float64(p.clients.Size()))

				log.Debug().Msgf("New Client %s connected.", e.Client.ID)
				log.Debug().Msgf("Now %d clients connected.", p.clients.Size())

			case ClientEventTypeLeave:
				event := NewTopicEvent(TopicEventTypeUnsubscribeAll, "all", e.Client, packets.Details{})
				p.topicsEventCh <- event

				if e.Client.IsAuthenticated() {
					p.users.DelClient(e.Client.UserID, e.Client.ID)

					metrics.ClientAuthenticate.WithLabelValues().Set(float64(p.users.Size()))

					log.Debug().Msgf("Client %s unauthenticated.", e.Client.ID)
				}

				log.Debug().Msgf("Client %s disconnected.", e.Client.ID)

				p.clients.Del(e.Client.ID)

				metrics.ClientLive.WithLabelValues().Set(float64(p.clients.Size()))

				log.Debug().Msgf("Now %d clients connected.", p.clients.Size())
			}
		}
	}
}

func (p *server) cleanTopic(topic *Topic) {
	if topic.Size() == 0 {
		log.Info().Msgf(`Remove "%s" topic`, topic.ID)

		p.topics.Del(topic.ID.String())

		topic.Close()
	}
}

func (p *server) processTopicEvent() {
	for {
		select {
		case e := <-p.topicsEventCh:
			switch e.Type {
			case TopicEventTypeSubscribe:
				// check if client is authenticated
				/*
					if strings.HasPrefix(e.Name, PrivateTopicPrefix) && !e.Client.IsAuthenticated() {
						e.Client.Write(message.NewEventFromErrorCode(message.ErrorCodeUnauthorized))

						break
					}
				*/
				topic, ok := p.topics.Get(e.Name)
				if !ok {
					log.Info().Msgf(`Create new "%s" topic`, e.Name)

					// Create new topic
					topic = NewTopic(e.Name)

					p.topics.Add(topic)

					// this goroutine is release by topic.Close() in cleanTopic()
					go topic.Listen()
				}

				// Add client to topic
				topic.Join(e.Client)

			case TopicEventTypeUnsubscribe:
				if topic, ok := p.topics.Get(e.Name); ok {
					topic.Leave(e.Client)

					p.cleanTopic(topic)
				}

			case TopicEventTypeUnsubscribeAll:
				topics := p.topics.Topics()

				for _, topic := range topics {
					topic.Leave(e.Client)

					p.cleanTopic(topic)
				}
			}
		}
	}
}

// ListenAndServe push server
func (p *server) ListenAndServe() error {
	go p.processTopicEvent()
	go p.processClientEvent()
	go p.processAuthenticateEvent()
	go p.processMessage()

	exit := make(chan struct{})

	// Create incoming connections listener.
	ln, err := net.Listen("tcp", p.config.Addr())
	if err != nil {
		return err
	}

	log.Info().Msgf("Push Server is listening on %s", ln.Addr().String())

	// Create netpoll descriptor for the listener.
	// We use OneShot here to manually resume events stream when we want to.
	acceptDesc := netpoll.Must(netpoll.HandleListener(
		ln,
		netpoll.EventRead|netpoll.EventOneShot,
	))

	// accept is a topic to signal about next incoming connection Accept()
	// results.
	accept := make(chan error, 1)

	// Subscribe to events about listener.
	p.poller.Start(acceptDesc, func(e netpoll.Event) {
		// We do not want to accept incoming connection when goroutine pool is
		// busy. So if there are no free goroutines during 1ms we want to
		// cooldown the server and do not receive connection for some short
		// time.
		err := p.pool.ScheduleTimeout(time.Millisecond, func() {
			conn, err := ln.Accept()
			if err != nil {
				accept <- err

				return
			}

			accept <- nil

			p.ServeTCP(conn)
		})
		if err == nil {
			err = <-accept
		}

		if err != nil {
			if err != pool.ErrScheduleTimeout {
				goto cooldown
			}

			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				goto cooldown
			}

			log.Fatal().Err(err).Msg("accept error")

		cooldown:
			delay := 5 * time.Millisecond

			log.Warn().Err(err).Msgf("accept error: retrying in %v", delay)

			time.Sleep(delay)
		}

		p.poller.Resume(acceptDesc)
	})

	<-exit

	return nil
}

func (p *server) ServeTCP(conn net.Conn) {
	if p.clients.Size() >= p.config.MaxConnections {
		//http.Error(conn, "Too Many Requests", http.StatusTooManyRequests)
		//@TODO: write http response
		return
	}

	sock := transport.NewConnection(conn, p.config.IOTimeoutDuration)

	// Zero-copy upgrade to WebSocket connection.
	hs, err := ws.Upgrade(sock)
	if err != nil {
		log.Error().Err(err).Msg("WebSocket.Upgrade")

		conn.Close()

		return
	}

	ts := transport.NewWebSocket(sock)

	log.Debug().Msgf("%s > %s: established websocket connection: %+v", conn.LocalAddr().String(), conn.RemoteAddr().String(), hs)

	//log.Printf("%s: established websocket connection: %+v", nameConn(conn), hs)

	client := NewClient(ts, p)

	event := NewClientEvent(ClientEventTypeJoin, client)
	p.clientsEventCh <- event

	metrics.ClientConnection.WithLabelValues().Add(1)

	// Create netpoll event descriptor for conn.
	// We want to handle only read events of it.
	desc := netpoll.Must(netpoll.HandleRead(conn))

	// Subscribe to events about conn.
	p.poller.Start(desc, func(ev netpoll.Event) {
		if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
			// When ReadHup or Hup received, this mean that client has
			// closed at least write end of the connection or connections
			// itself. So we want to stop receive events about such conn
			// and remove it from the chat registry.
			p.poller.Stop(desc)
			client.Close()
			//chat.Remove(user)

			return
		}
		// Here we can read some new message from connection.
		// We can not read it right here in callback, because then we will
		// block the poller's inner loop.
		// We do not want to spawn a new goroutine to read single message.
		// But we want to reuse previously spawned goroutine.
		p.pool.Schedule(func() {
			if err := client.ReadEvent(); err != nil {
				// When receive failed, we can only disconnect broken
				// connection and stop to receive events about it.
				p.poller.Stop(desc)
				client.Close()
				//chat.Remove(user)
			}
		})
	})
}
