// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	assert.True(t, prometheus.Unregister(ClientLive))
	assert.True(t, prometheus.Unregister(ClientConnection))
	assert.True(t, prometheus.Unregister(ClientAuthenticate))
	assert.True(t, prometheus.Unregister(ChannelClient))
	assert.True(t, prometheus.Unregister(ChannelMessageReceivedTotal))
	assert.True(t, prometheus.Unregister(ChannelMessageReceivedBytes))
	assert.True(t, prometheus.Unregister(ChannelMessageDeliveredTotal))
	assert.True(t, prometheus.Unregister(ChannelMessageDeliveredBytes))
	assert.True(t, prometheus.Unregister(ChannelMessageDeliveredSeconds))
}
