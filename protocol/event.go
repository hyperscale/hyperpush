// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package protocol

import (
	"encoding/json"
)

// Event struct
type Event struct {
	Type    string          `json:"type"`
	Channel string          `json:"channel,omitempty"`
	User    string          `json:"user,omitempty"`
	Name    string          `json:"name,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
	Code    int             `json:"code,omitempty"`
	Message string          `json:"message,omitempty"`
}
