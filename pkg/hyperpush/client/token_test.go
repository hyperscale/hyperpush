// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"testing"
	"time"
)

func TestWaitTimeout(t *testing.T) {
	b := baseToken{}

	if b.WaitTimeout(time.Second) {
		t.Fatal("Should have failed")
	}
}
