// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

// Configuration struct
type Configuration struct {
	ClientQueueSize         int `mapstructure:"client_queue_size"`
	ChannelQueueSize        int `mapstructure:"channel_queue_size"`
	AuthenticationQueueSize int `mapstructure:"authentication_queue_size"`
	MessageQueueSize        int `mapstructure:"message_queue_size"`
	MaxConnections          int `mapstructure:"max_connections"`
}
