// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTopicIDWithDollarPrefix(t *testing.T) {
	topic := TopicID("$SYS/cpu/load")

	assert.NoError(t, topic.IsValid())

	assert.True(t, topic.IsSystemTopic())

	assert.False(t, topic.Match("$SYS/cpu/info"))
	assert.True(t, topic.Match("$SYS/cpu/load"))
	assert.True(t, topic.Match("$SYS/cpu/+"))
	assert.True(t, topic.Match("$SYS/#"))
	assert.False(t, topic.Match("$SYS/+"))
	assert.True(t, topic.Match("$SYS/cpu/#"))
	assert.False(t, topic.Match("+/cpu/load"))
	assert.False(t, topic.Match("+"))
	assert.False(t, topic.Match("#/cpu/+"))

	assert.Equal(t, "$SYS/cpu/load", topic.String())
}

func TestTopicIDWithSimpleTopic(t *testing.T) {
	topic := TopicID("SYS")

	assert.NoError(t, topic.IsValid())

	assert.False(t, topic.IsSystemTopic())

	assert.True(t, topic.Match("+"))
	assert.True(t, topic.Match("#"))

	assert.Equal(t, "SYS", topic.String())
}

func TestTopicID(t *testing.T) {
	topic := TopicID("SYS/cpu/load")

	assert.NoError(t, topic.IsValid())

	assert.False(t, topic.IsSystemTopic())

	assert.False(t, topic.Match("SYS/cpu/info"))
	assert.True(t, topic.Match("SYS/cpu/load"))
	assert.True(t, topic.Match("SYS/cpu/+"))
	assert.True(t, topic.Match("SYS/#"))
	assert.False(t, topic.Match("SYS/+"))
	assert.True(t, topic.Match("SYS/cpu/#"))
	assert.True(t, topic.Match("+/cpu/#"))
	assert.True(t, topic.Match("+/cpu/load"))
	assert.True(t, topic.Match("#"))
	assert.False(t, topic.Match("#/cpu/+"))

	assert.Equal(t, "SYS/cpu/load", topic.String())
}

func TestTopicIDWithBadName(t *testing.T) {
	topic := TopicID("")

	assert.Error(t, topic.IsValid())

	topic = TopicID("#/foo")

	assert.Error(t, topic.IsValid())
}

func BenchmarkTopicIDMatch(b *testing.B) {
	b.SetParallelism(50)
	b.ReportAllocs()

	topic := TopicID("SYS/cpu/load")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			topic.Match("+/cpu/#")
		}
	})
}
