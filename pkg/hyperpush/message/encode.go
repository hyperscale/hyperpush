// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import (
	"encoding/json"
)

// Encode event to json string
func Encode(msg *Event) ([]byte, error) {
	if len(msg.Raw) > 0 {
		return msg.Raw, nil
	}

	return json.Marshal(msg)
}
