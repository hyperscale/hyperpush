// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"testing"

	"github.com/hyperscale/hyperpush/pkg/hyperpush/mqtt/packets"

	"github.com/stretchr/testify/assert"
)

func TestTopicEvent(t *testing.T) {
	c := NewClient(nil, nil)
	event := NewTopicEvent(TopicEventTypeSubscribe, "test", c, packets.Details{})

	assert.Equal(t, TopicEventTypeSubscribe, event.Type)
	assert.Equal(t, "test", event.Name)
	assert.Equal(t, c, event.Client)
}
