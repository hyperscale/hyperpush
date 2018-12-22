// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEventFromErrorCode(t *testing.T) {
	event := NewEventFromErrorCode(ErrorCodeBadCredentials)

	assert.Equal(t, EventTypeError, event.Type)
	assert.Equal(t, 401, event.Code)
	assert.Equal(t, "Bad Credentials", event.Message)
}

func TestNewEventFromError(t *testing.T) {
	err := NewError(ErrorCodeUnauthorized)

	assert.Equal(t, "Unauthorized", err.Error())

	event := NewEventFromError(err)

	assert.Equal(t, EventTypeError, event.Type)
	assert.Equal(t, 400, event.Code)
	assert.Equal(t, "Unauthorized", event.Message)
}

func TestNewErrorf(t *testing.T) {
	err := NewErrorf(ErrorCodeUnexpectedSigningMethod, "test")

	assert.Equal(t, ErrorCodeUnexpectedSigningMethod, err.Code())
	assert.Equal(t, "Unexpected signing method: test", err.Error())
}

func TestFromError(t *testing.T) {
	err := FromError(errors.New("fail"))

	assert.Equal(t, ErrorCodeUnknown, err.Code())
	assert.Equal(t, "fail", err.Error())
}
