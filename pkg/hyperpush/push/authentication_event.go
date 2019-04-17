// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

// AuthenticationEvent struct
type AuthenticationEvent struct {
	Token  string
	Client *Client
}

// NewAuthenticationEvent func
func NewAuthenticationEvent(token string, client *Client) *AuthenticationEvent {
	return &AuthenticationEvent{
		Token:  token,
		Client: client,
	}
}
