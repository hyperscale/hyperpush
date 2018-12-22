// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	encodeEventExpected = []byte(`{"type":"publish","channel":"bitcoin","name":"tick","data":{"price":0.5678}}`)
)

func TestEncode(t *testing.T) {
	event, err := Encode(&Event{
		Type:    "publish",
		Channel: "bitcoin",
		Name:    "tick",
		Data:    json.RawMessage(`{"price":0.5678}`),
	})
	assert.NoError(t, err)

	assert.Equal(t, encodeEventExpected, event)
}

func TestEncodeWithRaw(t *testing.T) {
	actual, err := Encode(&Event{
		Type:    "publish",
		Channel: "bitcoin",
		Name:    "tick",
		Data:    json.RawMessage(`{"price": 1.456}`),
		Raw:     encodeEventExpected,
	})
	assert.NoError(t, err)

	assert.Equal(t, encodeEventExpected, actual)
}

func BenchmarkEncode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Encode(&Event{
			Type:    "publish",
			Channel: "bitcoin",
			Name:    "tick",
			Data:    json.RawMessage(`{"foo": 1.3578}`),
		})
	}
}
