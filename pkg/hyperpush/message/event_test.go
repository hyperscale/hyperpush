// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventUnmarshalJSON(t *testing.T) {
	event := &Event{}

	err := event.UnmarshalJSON(msg)
	assert.NoError(t, err)
	assert.Equal(t, EventTypePublish, event.Type)
	assert.Equal(t, "bitcoin", event.Channel)
	assert.Equal(t, "tick", event.Name)
}

func TestEventUnmarshalJSONWithoutData(t *testing.T) {
	event := &Event{}

	err := event.UnmarshalJSON(nil)
	assert.Error(t, err)
}
