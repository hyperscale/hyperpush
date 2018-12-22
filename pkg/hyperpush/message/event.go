// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import (
	"encoding/json"
)

// EventType type
type EventType string

// EventType enums
const (
	EventTypeConnect        EventType = "connect"
	EventTypeConnected      EventType = "connected"
	EventTypePublish        EventType = "publish"
	EventTypeError          EventType = "error"
	EventTypeMessage        EventType = "message"
	EventTypeSubscribe      EventType = "subscribe"
	EventTypeSubscribed     EventType = "subscribed"
	EventTypeUnsubscribe    EventType = "unsubscribe"
	EventTypeUnsubscribed   EventType = "unsubscribed"
	EventTypePing           EventType = "ping"
	EventTypePong           EventType = "pong"
	EventTypeAuthentication EventType = "authentication"
	EventTypeAuthenticated  EventType = "authenticated"
)

// Event struct
type Event struct {
	Type    EventType       `json:"type"`
	Channel string          `json:"channel,omitempty"`
	User    int             `json:"user,omitempty"`
	Name    string          `json:"name,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
	Token   string          `json:"token,omitempty"`
	Code    int             `json:"code,omitempty"`
	Message string          `json:"message,omitempty"`
	Raw     []byte          `json:"-"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *Event) UnmarshalJSON(data []byte) error {
	d := struct {
		Type    EventType       `json:"type"`
		Channel string          `json:"channel,omitempty"`
		User    int             `json:"user,omitempty"`
		Name    string          `json:"name,omitempty"`
		Data    json.RawMessage `json:"data,omitempty"`
		Token   string          `json:"token,omitempty"`
		Code    int             `json:"code,omitempty"`
		Message string          `json:"message,omitempty"`
	}{}

	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	e.Type = d.Type
	e.Channel = d.Channel
	e.User = d.User
	e.Name = d.Name
	e.Data = d.Data
	e.Token = d.Token
	e.Code = d.Code
	e.Message = d.Message
	e.Raw = data

	return nil
}
