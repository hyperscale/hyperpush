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
	msg = []byte(`{"type":"publish","channel":"bitcoin","name":"tick","data":{"price": 1.5788}}`)
)

func TestDecode(t *testing.T) {
	message, err := Decode(msg)
	if err != nil {
		t.Errorf(err.Error())
	}

	var data map[string]interface{}

	err = json.Unmarshal(message.Data, &data)
	assert.NoError(t, err)

	assert.Equal(t, EventTypePublish, message.Type)
	assert.Equal(t, "bitcoin", message.Channel)
	assert.Equal(t, "tick", message.Name)
	assert.Equal(t, map[string]interface{}{"price": 1.5788}, data)
}

func TestDecodeBadJson(t *testing.T) {
	_, err := Decode([]byte(`{"foo":"bar"`))

	assert.Error(t, err)
}

func BenchmarkDecode(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		Decode(msg)
	}
}
