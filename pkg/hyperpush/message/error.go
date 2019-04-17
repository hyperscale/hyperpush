// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import (
	"fmt"
)

// ErrorCode type
type ErrorCode int

// ErrorCode constant
const (
	ErrorCodeUnauthorized            ErrorCode = 400
	ErrorCodeBadCredentials          ErrorCode = 401
	ErrorCodeCredentialsExpired      ErrorCode = 402
	ErrorCodeUnexpectedSigningMethod ErrorCode = 403
	ErrorCodeInvalidCredentials      ErrorCode = 404
	ErrorCodeUnknown                 ErrorCode = 999
)

var errorCodeMessage = map[ErrorCode]string{
	ErrorCodeUnknown:                 "Unknown",
	ErrorCodeUnauthorized:            "Unauthorized",
	ErrorCodeBadCredentials:          "Bad Credentials",
	ErrorCodeCredentialsExpired:      "Credentials Expired",
	ErrorCodeUnexpectedSigningMethod: "Unexpected signing method: %v",
}

// NewEventFromError func
func NewEventFromError(err ErrorInterface) *Event {
	return NewEventFromErrorCode(err.Code())
}

// NewEventFromErrorCode from ErrorCode
func NewEventFromErrorCode(code ErrorCode) *Event {
	return &Event{
		Type:    EventTypeError,
		Code:    int(code),
		Message: errorCodeMessage[code],
	}
}

// ErrorInterface for error
//go:generate mockery -case=underscore -inpkg -name=ErrorInterface
type ErrorInterface interface {
	error
	Code() ErrorCode
}

// Error type
type Error struct {
	code    ErrorCode
	message string
}

// NewError from ErrorCode
func NewError(code ErrorCode) Error {
	return Error{
		code:    code,
		message: errorCodeMessage[code],
	}
}

// NewErrorf from ErrorCode with arguments
func NewErrorf(code ErrorCode, a ...interface{}) Error {
	return Error{
		code:    code,
		message: fmt.Sprintf(errorCodeMessage[code], a...),
	}
}

// FromError std
func FromError(err error) Error {
	return Error{
		code:    ErrorCodeUnknown,
		message: err.Error(),
	}
}

func (e Error) Error() string {
	return e.message
}

// Code of error
func (e Error) Code() ErrorCode {
	return e.code
}
