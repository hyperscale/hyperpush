// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"sync"
)

// ClientPool struct
type ClientPool struct {
	clients    map[string]*Client
	clientsMtx *sync.RWMutex
}

// NewClientPool constructor
func NewClientPool() *ClientPool {
	return &ClientPool{
		clients:    make(map[string]*Client),
		clientsMtx: &sync.RWMutex{},
	}
}

// Add client to manager
func (c *ClientPool) Add(client *Client) {
	c.clientsMtx.Lock()
	defer c.clientsMtx.Unlock()

	c.clients[client.ID] = client
}

// Has client exists by id
func (c *ClientPool) Has(ID string) bool {
	c.clientsMtx.RLock()
	defer c.clientsMtx.RUnlock()

	_, ok := c.clients[ID]

	return ok
}

// Get client by id
func (c *ClientPool) Get(ID string) (*Client, bool) {
	c.clientsMtx.RLock()
	defer c.clientsMtx.RUnlock()

	client, ok := c.clients[ID]

	return client, ok
}

// Del client to manager
func (c *ClientPool) Del(ID string) {
	c.clientsMtx.Lock()
	defer c.clientsMtx.Unlock()

	delete(c.clients, ID)
}

// Size of clients
func (c *ClientPool) Size() int {
	return len(c.clients)
}

// Clients list
func (c *ClientPool) Clients() []*Client {
	clients := []*Client{}

	c.clientsMtx.RLock()
	defer c.clientsMtx.RUnlock()

	for _, client := range c.clients {
		clients = append(clients, client)
	}

	return clients
}
