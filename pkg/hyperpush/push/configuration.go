// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"net"
	"strconv"
	"time"
)

// Configuration struct
type Configuration struct {
	Host                    string
	Port                    int
	ClientQueueSize         int           `mapstructure:"client_queue_size"`
	TopicQueueSize          int           `mapstructure:"topic_queue_size"`
	AuthenticationQueueSize int           `mapstructure:"authentication_queue_size"`
	MessageQueueSize        int           `mapstructure:"message_queue_size"`
	MaxConnections          int           `mapstructure:"max_connections"`
	ConnectionWorkerSize    int           `mapstructure:"connection_worker_size"`
	ConnectionQueueSize     int           `mapstructure:"connection_queue_size"`
	IOTimeoutDuration       time.Duration `mapstructure:"io_timeout_duration"`
}

// Addr return host and port string
func (c Configuration) Addr() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}
