// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannelPool(t *testing.T) {
	cm := NewChannelPool()

	assert.Equal(t, 0, cm.Size())

	cm.Add(NewChannel("test"))

	assert.Equal(t, 1, cm.Size())

	cm.Add(NewChannel("test 2"))

	assert.Equal(t, 2, cm.Size())

	c1 := NewChannel("test")
	cm.Add(c1)
	cm.Add(c1)

	assert.Equal(t, 2, cm.Size())

	c2, ok := cm.Get("test")

	assert.True(t, ok)

	assert.Equal(t, c1, c2)

	cm.Del("test")

	assert.Equal(t, 1, cm.Size())

	_, ok = cm.Get("test")

	assert.False(t, ok)

	channels := cm.Channels()

	assert.Equal(t, cm.Size(), len(channels))
}

func TestChannelPoolConcurrency(t *testing.T) {
	var wg sync.WaitGroup

	cm := NewChannelPool()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < 10000; i++ {
			cm.Add(NewChannel(fmt.Sprintf("test-%d", i)))
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 10000; i < 20000; i++ {
			cm.Add(NewChannel(fmt.Sprintf("test-%d", i)))
		}
	}()

	wg.Wait()

	assert.Equal(t, 20000, cm.Size())
}
