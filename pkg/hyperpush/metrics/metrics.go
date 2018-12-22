// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(ClientLive)
	prometheus.MustRegister(ClientConnection)
	prometheus.MustRegister(ClientAuthenticate)
	prometheus.MustRegister(ChannelClient)
	prometheus.MustRegister(MessageReceivedTotal)
	prometheus.MustRegister(ChannelMessageReceivedTotal)
	prometheus.MustRegister(ChannelMessageReceivedBytes)
	prometheus.MustRegister(ChannelMessageDeliveredTotal)
	prometheus.MustRegister(ChannelMessageDeliveredBytes)
	prometheus.MustRegister(ChannelMessageDeliveredSeconds)
}

// ClientLive metrics
var ClientLive = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "client_live_total",
		Help: "The count of client connected.",
	},
	[]string{},
)

// ClientConnection metrics
var ClientConnection = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "client_connection_total",
		Help: "The count of connection.",
	},
	[]string{},
)

// ClientAuthenticate metrics
var ClientAuthenticate = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "client_authenticate_total",
		Help: "The count of client authenticated.",
	},
	[]string{},
)

// MessageReceivedTotal metrics
var MessageReceivedTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "message_received_total",
		Help: "The count of message received.",
	},
	[]string{"type"},
)

// ChannelClient metrics
var ChannelClient = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "channel_client_total",
		Help: "The count of client connected on channel.",
	},
	[]string{"channel"},
)

// ChannelMessageReceivedTotal metrics
var ChannelMessageReceivedTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "channel_message_received_total",
		Help: "The count of message received.",
	},
	[]string{"channel"},
)

// ChannelMessageReceivedBytes metrics
var ChannelMessageReceivedBytes = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "channel_message_received_bytes",
		Help: "The bytes of message received.",
	},
	[]string{"channel"},
)

// ChannelMessageDeliveredTotal metrics
var ChannelMessageDeliveredTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "channel_message_delivered_total",
		Help: "The count of message delivered.",
	},
	[]string{"channel"},
)

// ChannelMessageDeliveredBytes metrics
var ChannelMessageDeliveredBytes = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "channel_message_delivered_bytes",
		Help: "The bytes of message delivered.",
	},
	[]string{"channel"},
)

// ChannelMessageDeliveredSeconds metrics
var ChannelMessageDeliveredSeconds = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "channel_message_delivered_seconds",
		Help:    "The latencies of message delivered in seconds.",
		Buckets: []float64{0.01, 0.1, 0.3, 0.5, 1., 2., 5.},
	},
	[]string{"channel"},
)
