// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserPool(t *testing.T) {
	up := NewUserPool()

	assert.Equal(t, 0, up.Size())

	c1 := NewClient(context.Background(), nil, nil)
	up.Add(1, c1)

	clients, ok := up.Get(1)
	assert.True(t, ok)
	assert.Equal(t, 1, len(clients))

	ok = up.HasClient(1, c1.ID)

	assert.True(t, ok)

	ok = up.HasClient(2, c1.ID)

	assert.False(t, ok)

	assert.Equal(t, 1, up.Size())
	c2 := NewClient(context.Background(), nil, nil)

	up.Add(1, c2)
	up.Add(1, c2)

	ok = up.Has(1)
	assert.True(t, ok)

	clients, ok = up.Get(1)
	assert.True(t, ok)
	assert.Equal(t, 2, len(clients))

	assert.Equal(t, 1, up.Size())

	up.DelClient(1, c2.ID)

	clients, ok = up.Get(1)
	assert.True(t, ok)
	assert.Equal(t, 1, len(clients))

	up.Add(2, NewClient(context.Background(), nil, nil))

	assert.Equal(t, 2, up.Size())

	c3 := NewClient(context.Background(), nil, nil)
	up.Add(3, c3)

	assert.Equal(t, 3, up.Size())

	_, ok = up.Get(3)

	assert.True(t, ok)

	up.DelClient(1, c1.ID)

	assert.Equal(t, 2, up.Size())
}

func TestUserPoolConcurrency(t *testing.T) {
	var wg sync.WaitGroup

	up := NewUserPool()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < 100000; i++ {
			up.Add(i, NewClient(context.Background(), nil, nil))
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 100000; i < 200000; i++ {
			up.Add(i, NewClient(context.Background(), nil, nil))
		}
	}()

	wg.Wait()

	assert.Equal(t, 200000, up.Size())
}

func BenchmarkUserPool(b *testing.B) {
	b.SetParallelism(50)
	b.ReportAllocs()

	up := NewUserPool()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++

			up.Add(i, NewClient(context.Background(), nil, nil))
		}
	})
}
