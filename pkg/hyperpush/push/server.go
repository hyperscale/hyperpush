// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"net/http"
	"strings"

	"github.com/euskadi31/go-eventemitter"
	"github.com/gorilla/websocket"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/authentication"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/message"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/metrics"
	"github.com/rs/zerolog/log"
)

// ErrorHandler func
type ErrorHandler func(w http.ResponseWriter, code int, err error)

// Server interface
//go:generate mockery -case=underscore -inpkg -name=Server
type Server interface {
	Authenticate(token string, client *Client)
	SetAuthenticationProvider(provider authentication.Provider)
	JoinChannel(ID string, client *Client)
	LeaveChannel(ID string, client *Client)
	Leave(client *Client)
	Publish(message *message.Event)
	Listen()
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type server struct {
	config              *Configuration
	upgrader            websocket.Upgrader
	ErrorHandler        ErrorHandler
	clients             *ClientPool
	clientsEventCh      chan *ClientEvent
	users               UserPool
	channels            *ChannelPool
	channelsEventCh     chan *ChannelEvent
	authentication      authentication.Provider
	authenticateEventCh chan *AuthenticationEvent
	messagesCh          chan *message.Event
	emitter             eventemitter.EventEmitter
}

// NewServer server
func NewServer(config *Configuration, emitter eventemitter.EventEmitter) Server {
	return &server{
		config: config,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		ErrorHandler: func(w http.ResponseWriter, code int, err error) {
			log.Error().Err(err).Msg("error")

			http.Error(w, err.Error(), code)
		},
		clients:             NewClientPool(),
		users:               NewUserPool(),
		clientsEventCh:      make(chan *ClientEvent, config.ClientQueueSize),
		channels:            NewChannelPool(),
		channelsEventCh:     make(chan *ChannelEvent, config.ChannelQueueSize),
		authenticateEventCh: make(chan *AuthenticationEvent, config.AuthenticationQueueSize),
		messagesCh:          make(chan *message.Event, config.MessageQueueSize),
		emitter:             emitter,
	}
}

// SetAuthenticationProvider to server
func (p *server) SetAuthenticationProvider(provider authentication.Provider) {
	p.authentication = provider
}

// Authenticate user
func (p *server) Authenticate(token string, client *Client) {
	event := NewAuthenticationEvent(token, client)

	p.authenticateEventCh <- event
}

// Publish message to server
func (p *server) Publish(message *message.Event) {
	p.messagesCh <- message
}

// JoinChannel added client to channel
func (p *server) JoinChannel(ID string, client *Client) {
	event := NewChannelEvent(ChannelEventTypeSubscribe, ID, client)

	p.channelsEventCh <- event
}

// LeaveChannel removed client to channel
func (p *server) LeaveChannel(ID string, client *Client) {
	event := NewChannelEvent(ChannelEventTypeUnsubscribe, ID, client)

	p.channelsEventCh <- event
}

// Leave removed client to all channel
func (p *server) Leave(client *Client) {
	event := NewClientEvent(ClientEventTypeLeave, client)

	p.clientsEventCh <- event
}

func (p *server) dispatchMessage(event *message.Event) {
	if event.User != "" {
		// send message to user on all client
		if clients, ok := p.users.Get(event.User); ok {
			for _, client := range clients {
				client.Write(event)
			}
		}
	} else if event.Channel != "" {
		// broadcast message to channel
		if channel, ok := p.channels.Get(event.Channel); ok {
			channel.Publish(event)
		}
	}
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
			user, err := p.authentication.Authenticate(e.Token)
			if err != nil {
				e.Client.Write(message.NewEventFromError(err))
			} else {
				e.Client.UserID = user.ID

				if p.users.HasClient(user.ID, e.Client.ID) {
					log.Debug().Msgf("Client %s already autenticated", e.Client.ID)

					e.Client.Write(&message.Event{
						Type: message.EventTypeAuthenticated,
						User: user.ID,
					})

					break
				}

				p.users.Add(user.ID, e.Client)

				metrics.ClientAuthenticate.WithLabelValues().Set(float64(p.users.Size()))

				e.Client.Write(&message.Event{
					Type: message.EventTypeAuthenticated,
					User: user.ID,
				})

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

				e.Client.Write(&message.Event{
					Type: message.EventTypeConnected,
				})

				metrics.ClientLive.WithLabelValues().Set(float64(p.clients.Size()))

				log.Debug().Msgf("New Client %s connected.", e.Client.ID)
				log.Debug().Msgf("Now %d clients connected.", p.clients.Size())

			case ClientEventTypeLeave:
				event := NewChannelEvent(ChannelEventTypeUnsubscribeAll, "all", e.Client)
				p.channelsEventCh <- event

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

func (p *server) cleanChannel(channel *Channel) {
	if channel.Size() == 0 {
		log.Info().Msgf(`Remove "%s" channel`, channel.ID)

		p.channels.Del(channel.ID)

		channel.Close()
	}
}

func (p *server) processChannelEvent() {
	for {
		select {
		case e := <-p.channelsEventCh:
			switch e.Type {
			case ChannelEventTypeSubscribe:
				// check if client is authenticated
				if strings.HasPrefix(e.Name, PrivateChannelPrefix) && !e.Client.IsAuthenticated() {
					e.Client.Write(message.NewEventFromErrorCode(message.ErrorCodeUnauthorized))

					break
				}

				channel, ok := p.channels.Get(e.Name)
				if !ok {
					log.Info().Msgf(`Create new "%s" channel`, e.Name)

					// Create new channel
					channel = NewChannel(e.Name)

					p.channels.Add(channel)

					// this goroutine is release by channel.Close() in cleanChannel()
					go channel.Listen()
				}

				// Add client to channel
				channel.Join(e.Client)

			case ChannelEventTypeUnsubscribe:
				if channel, ok := p.channels.Get(e.Name); ok {
					channel.Leave(e.Client)
					p.cleanChannel(channel)
				}

			case ChannelEventTypeUnsubscribeAll:
				channels := p.channels.Channels()

				for _, channel := range channels {
					channel.Leave(e.Client)
					p.cleanChannel(channel)
				}
			}
		}
	}
}

// Listen server
func (p *server) Listen() {
	go p.processChannelEvent()
	go p.processClientEvent()
	go p.processAuthenticateEvent()
	go p.processMessage()

}

// ListenAndServe push server
func (p *server) ListenAndServe() error {

	return nil
}

// ServeHTTP handler
func (p *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.clients.Size() >= p.config.MaxConnections {
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)

		return
	}

	ws, err := p.upgrader.Upgrade(w, r, nil)
	if err != nil {
		p.ErrorHandler(w, http.StatusInternalServerError, err)

		return
	}

	ctx := r.Context()

	client := NewClient(ctx, ws, p)
	defer func() {
		log.Debug().Msgf("Closing client %s", client.ID)

		if err := client.Close(); err != nil {
			log.Error().Err(err).Msg("Client.Close")
		}
	}()

	event := NewClientEvent(ClientEventTypeJoin, client)
	p.clientsEventCh <- event

	metrics.ClientConnection.WithLabelValues().Add(1)

	client.Listen()
}
