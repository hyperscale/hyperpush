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
	prometheus.MustRegister(TopicClient)
	prometheus.MustRegister(MessageReceivedTotal)
	prometheus.MustRegister(TopicMessageReceivedTotal)
	prometheus.MustRegister(TopicMessageReceivedBytes)
	prometheus.MustRegister(TopicMessageDeliveredTotal)
	prometheus.MustRegister(TopicMessageDeliveredBytes)
	prometheus.MustRegister(TopicMessageDeliveredSeconds)
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

// TopicClient metrics
var TopicClient = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "topic_client_total",
		Help: "The count of client connected on topic.",
	},
	[]string{"topic"},
)

// TopicMessageReceivedTotal metrics
var TopicMessageReceivedTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "topic_message_received_total",
		Help: "The count of message received.",
	},
	[]string{"topic"},
)

// TopicMessageReceivedBytes metrics
var TopicMessageReceivedBytes = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "topic_message_received_bytes",
		Help: "The bytes of message received.",
	},
	[]string{"topic"},
)

// TopicMessageDeliveredTotal metrics
var TopicMessageDeliveredTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "topic_message_delivered_total",
		Help: "The count of message delivered.",
	},
	[]string{"topic"},
)

// TopicMessageDeliveredBytes metrics
var TopicMessageDeliveredBytes = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "topic_message_delivered_bytes",
		Help: "The bytes of message delivered.",
	},
	[]string{"topic"},
)

// TopicMessageDeliveredSeconds metrics
var TopicMessageDeliveredSeconds = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "topic_message_delivered_seconds",
		Help:    "The latencies of message delivered in seconds.",
		Buckets: []float64{0.01, 0.1, 0.3, 0.5, 1., 2., 5.},
	},
	[]string{"topic"},
)
