// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"sync"
)

// TopicPool struct
type TopicPool struct {
	topics    map[string]*Topic
	topicsMtx *sync.RWMutex
}

// NewTopicPool constructor
func NewTopicPool() *TopicPool {
	return &TopicPool{
		topics:    make(map[string]*Topic),
		topicsMtx: &sync.RWMutex{},
	}
}

// Add channel to manager
func (c *TopicPool) Add(topic *Topic) {
	c.topicsMtx.Lock()
	defer c.topicsMtx.Unlock()

	c.topics[topic.ID.String()] = topic
}

// Get channel by id
func (c *TopicPool) Get(ID string) (*Topic, bool) {
	c.topicsMtx.RLock()
	defer c.topicsMtx.RUnlock()

	channel, ok := c.topics[ID]

	return channel, ok
}

// Del channel to manager
func (c *TopicPool) Del(ID string) {
	c.topicsMtx.Lock()
	defer c.topicsMtx.Unlock()

	delete(c.topics, ID)
}

// Size of topics
func (c *TopicPool) Size() int {
	return len(c.topics)
}

// Topics list
func (c *TopicPool) Topics() map[string]*Topic {
	topics := make(map[string]*Topic)

	c.topicsMtx.RLock()
	defer c.topicsMtx.RUnlock()

	for id, channel := range c.topics {
		topics[id] = channel
	}

	return topics
}
