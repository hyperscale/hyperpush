// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/authentication"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/message"
)

func convertError(err error) message.Error {
	switch err.Error() {
	case jwt.ErrSignatureInvalid.Error(), "token contains an invalid number of segments":
		return message.NewError(message.ErrorCodeBadCredentials)
	default:
		return message.FromError(err)
	}
}

type provider struct {
	config *Configuration
}

// Authenticate user
func (p provider) Authenticate(accessToken string) (*authentication.User, message.ErrorInterface) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, message.NewErrorf(message.ErrorCodeUnexpectedSigningMethod, token.Header["alg"])
		}

		return []byte(p.config.Key), nil
	})
	if err != nil {
		return nil, convertError(err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok == false || token.Valid == false {
		return nil, message.NewError(message.ErrorCodeBadCredentials)
	}

	if _, ok := claims["sub"]; !ok {
		return nil, message.NewError(message.ErrorCodeInvalidCredentials)
	}

	if _, ok := claims["exp"]; !ok {
		return nil, message.NewError(message.ErrorCodeInvalidCredentials)
	}

	now := time.Now().UTC().Unix()
	if int64(claims["exp"].(float64)) < now {
		return nil, message.NewError(message.ErrorCodeCredentialsExpired)
	}

	return &authentication.User{
		ID: claims["sub"].(string),
	}, nil
}
