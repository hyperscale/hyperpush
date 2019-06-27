// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientPool(t *testing.T) {
	cp := NewClientPool()

	assert.Equal(t, 0, cp.Size())

	c0 := NewClient(nil, nil)
	cp.Add(c0)

	assert.Equal(t, 1, cp.Size())

	c1 := NewClient(nil, nil)
	cp.Add(c1)

	cp.Add(c1)

	assert.Equal(t, 2, cp.Size())

	ok := cp.Has(c1.ID)

	assert.True(t, ok)

	c2, ok := cp.Get(c1.ID)

	assert.True(t, ok)

	assert.Equal(t, c1, c2)

	cp.Del(c1.ID)

	assert.Equal(t, 1, cp.Size())

	_, ok = cp.Get(c1.ID)

	assert.False(t, ok)

	clients := cp.Clients()

	assert.Equal(t, cp.Size(), len(clients))
}

func TestClientPoolConcurrency(t *testing.T) {
	var wg sync.WaitGroup

	cp := NewClientPool()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < 10000; i++ {
			cp.Add(NewClient(nil, nil))
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 10000; i < 20000; i++ {
			cp.Add(NewClient(nil, nil))
		}
	}()

	wg.Wait()

	assert.Equal(t, 20000, cp.Size())
}
