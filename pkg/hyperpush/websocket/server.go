// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package websocket

// Server interface
type Server interface {
	ListenAndServe() error
}

type server struct {
}

// New Websocket server
func New() Server {
	return &server{}
}

func (s *server) ListenAndServe() error {

	return nil
}
