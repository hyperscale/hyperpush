// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientEvent(t *testing.T) {
	c := NewClient(context.Background(), nil, nil)
	event := NewClientEvent(ClientEventTypeJoin, c)

	assert.Equal(t, ClientEventTypeJoin, event.Type)
	assert.Equal(t, c, event.Client)
}
