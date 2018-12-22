// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"sync"
)

// ChannelPool struct
type ChannelPool struct {
	channels    map[string]*Channel
	channelsMtx *sync.RWMutex
}

// NewChannelPool constructor
func NewChannelPool() *ChannelPool {
	return &ChannelPool{
		channels:    make(map[string]*Channel),
		channelsMtx: &sync.RWMutex{},
	}
}

// Add channel to manager
func (c *ChannelPool) Add(channel *Channel) {
	c.channelsMtx.Lock()
	defer c.channelsMtx.Unlock()

	c.channels[channel.ID] = channel
}

// Get channel by id
func (c *ChannelPool) Get(ID string) (*Channel, bool) {
	c.channelsMtx.RLock()
	defer c.channelsMtx.RUnlock()

	channel, ok := c.channels[ID]

	return channel, ok
}

// Del channel to manager
func (c *ChannelPool) Del(ID string) {
	c.channelsMtx.Lock()
	defer c.channelsMtx.Unlock()

	delete(c.channels, ID)
}

// Size of channels
func (c *ChannelPool) Size() int {
	return len(c.channels)
}

// Channels list
func (c *ChannelPool) Channels() map[string]*Channel {
	channels := make(map[string]*Channel)

	c.channelsMtx.RLock()
	defer c.channelsMtx.RUnlock()

	for id, channel := range c.channels {
		channels[id] = channel
	}

	return channels
}
