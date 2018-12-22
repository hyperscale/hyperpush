// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"sync"
)

// UserPool interface
type UserPool interface {
	Add(id int, client *Client)
	Has(id int) bool
	HasClient(id int, clientID string) bool
	Get(id int) (map[string]*Client, bool)
	DelClient(id int, clientID string)
	Del(id int)
	Size() int
}

// UserPool struct
type userPool struct {
	users      map[int]map[string]*Client
	clientsMtx *sync.RWMutex
}

// NewUserPool constructor
func NewUserPool() UserPool {
	return &userPool{
		users:      make(map[int]map[string]*Client),
		clientsMtx: &sync.RWMutex{},
	}
}

// Add user to pool
func (u *userPool) Add(id int, client *Client) {
	u.clientsMtx.Lock()
	defer u.clientsMtx.Unlock()

	clients, ok := u.users[id]
	if ok == false {
		clients = make(map[string]*Client)

		u.users[id] = clients
	}

	clients[client.ID] = client

}

// Has user authenticated
func (u *userPool) Has(id int) bool {
	u.clientsMtx.RLock()
	defer u.clientsMtx.RUnlock()

	_, ok := u.users[id]

	return ok
}

// HasClient authenticated
func (u *userPool) HasClient(id int, clientID string) bool {
	if clients, ok := u.Get(id); ok {
		u.clientsMtx.RLock()
		defer u.clientsMtx.RUnlock()

		_, ok := clients[clientID]

		return ok
	}

	return false
}

// Get clients by user id
func (u *userPool) Get(id int) (map[string]*Client, bool) {
	u.clientsMtx.RLock()
	defer u.clientsMtx.RUnlock()

	clients, ok := u.users[id]

	return clients, ok
}

// DelClient by id
func (u *userPool) DelClient(id int, clientID string) {
	if clients, ok := u.Get(id); ok {
		u.clientsMtx.Lock()
		delete(clients, clientID)
		u.clientsMtx.Unlock()

		if len(clients) == 0 {
			u.Del(id)
		}
	}
}

// Del user to pool
func (u *userPool) Del(id int) {
	u.clientsMtx.Lock()
	defer u.clientsMtx.Unlock()

	delete(u.users, id)
}

// Size of users
func (u *userPool) Size() int {
	return len(u.users)
}
