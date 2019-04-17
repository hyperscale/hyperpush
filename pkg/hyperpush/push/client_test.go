// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"

	ws "github.com/gorilla/websocket"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/message"
	"github.com/hyperscale/hyperpush/pkg/hyperpush/websocket"
	"github.com/stretchr/testify/mock"
)

func TestClientWrtie(t *testing.T) {
	var wgm sync.WaitGroup

	serverMock := &MockServer{}

	serverMock.On("Leave", mock.AnythingOfType("*push.Client"))

	connMock := &websocket.MockConnection{}

	wgm.Add(1)

	connMock.On("SetWriteDeadline", mock.Anything).Return(nil)

	connMock.On("WriteMessage", 1, []byte(`{"type":"message","channel":"test","name":"foo","data":"bar"}`)).Return(nil).Run(func(fn mock.Arguments) {
		defer wgm.Done()
	})

	connMock.On("Close").Return(nil)

	client := NewClient(context.Background(), connMock, serverMock)

	go client.processMessages()

	client.Write(&message.Event{
		Type:    message.EventTypeMessage,
		Channel: "test",
		Name:    "foo",
		Data:    json.RawMessage(`"bar"`),
	})

	wgm.Wait()

	client.Close()

	connMock.AssertExpectations(t)

	serverMock.AssertExpectations(t)
}

func TestClientWriteFail(t *testing.T) {
	connMock := &websocket.MockConnection{}

	client := NewClient(context.Background(), connMock, nil)

	client.Write(&message.Event{
		Type:    message.EventTypeMessage,
		Channel: "test",
		Name:    "foo",
		Data:    json.RawMessage(`"bar"`),
	})

	connMock.AssertExpectations(t)
}

func TestClientPing(t *testing.T) {
	var wgw sync.WaitGroup

	serverMock := &MockServer{}

	serverMock.On("Leave", mock.AnythingOfType("*push.Client"))

	connMock := &websocket.MockConnection{}

	connMock.On("SetWriteDeadline", mock.Anything).Return(nil)
	connMock.On("SetReadLimit", mock.Anything).Return(nil)

	connMock.On("ReadMessage").Return(1, []byte(`{"type":"ping"}`), nil).Once()

	wgw.Add(1)

	connMock.On("WriteMessage", 1, []byte(`{"type":"pong"}`)).Return(nil).Run(func(args mock.Arguments) {
		defer wgw.Done()
	})

	connMock.On("ReadMessage").Return(1, []byte(""), &ws.CloseError{Code: ws.CloseUnsupportedData}).Once()

	connMock.On("Close").Return(nil)

	client := NewClient(context.Background(), connMock, serverMock)

	client.Listen()

	wgw.Wait()

	client.Close()

	connMock.AssertExpectations(t)
	serverMock.AssertExpectations(t)
}

func TestClientSubscribe(t *testing.T) {
	serverMock := &MockServer{}

	serverMock.On("Leave", mock.AnythingOfType("*push.Client"))

	serverMock.On("JoinChannel", "test", mock.AnythingOfType("*push.Client"))

	connMock := &websocket.MockConnection{}

	//connMock.On("SetWriteDeadline", mock.Anything).Return(nil)
	connMock.On("SetReadLimit", mock.Anything).Return(nil)

	connMock.On("ReadMessage").Return(1, []byte(`{"type":"subscribe","channel":"test"}`), nil).Once()

	connMock.On("ReadMessage").Return(1, []byte(""), &ws.CloseError{Code: ws.CloseUnsupportedData}).Once()

	connMock.On("Close").Return(nil)

	client := NewClient(context.Background(), connMock, serverMock)

	client.Listen()

	client.Close()

	connMock.AssertExpectations(t)
	serverMock.AssertExpectations(t)
}

/*
func TestClientMessageWithoutAuthenticate(t *testing.T) {
	serverMock := &MockServer{}

	connMock := &websocket.MockConnection{}

	client := NewClient(context.Background(), connMock, serverMock)

	assert.False(t, client.IsAuthenticated())

	client.Listen()

	client.Close()
}
*/
func TestClientBadEventType(t *testing.T) {
	var wgw sync.WaitGroup

	serverMock := &MockServer{}

	serverMock.On("Leave", mock.AnythingOfType("*push.Client")).Once()

	connMock := &websocket.MockConnection{}

	connMock.On("SetWriteDeadline", mock.Anything).Return(nil).Once()

	connMock.On("SetReadLimit", mock.Anything).Return(nil)

	wgw.Add(1)

	connMock.On("WriteMessage", 1, []byte(`{"type":"error","data":"Unsupported \"bad\" event type"}`)).Return(nil).Run(func(args mock.Arguments) {
		defer wgw.Done()
	}).Once()

	connMock.On("ReadMessage").Return(1, []byte(`{"type":"bad","channel":"test"}`), nil).Once()

	connMock.On("ReadMessage").Return(1, []byte(""), &ws.CloseError{Code: ws.CloseUnsupportedData}).Once()

	connMock.On("Close").Return(nil).Once()

	client := NewClient(context.Background(), connMock, serverMock)

	client.Listen()

	wgw.Wait()

	client.Close()

	connMock.AssertExpectations(t)
	serverMock.AssertExpectations(t)
}

//@TODO: Bug
/*
func TestClientBadMessageFormat(t *testing.T) {
	var wgw sync.WaitGroup

	serverMock := &MockServer{}

	serverMock.On("Leave", mock.AnythingOfType("*push.Client")).Once()

	connMock := &websocket.MockConnection{}

	connMock.On("SetWriteDeadline", mock.Anything).Return(nil).Once()

	connMock.On("SetReadLimit", mock.Anything).Return(nil)

	wgw.Add(1)

	connMock.On("WriteMessage", 1, []byte(`{"type":"error","data":"Bad request"}`)).Return(nil).Run(func(args mock.Arguments) {
		defer wgw.Done()
	}).Once()

	connMock.On("ReadMessage").Return(1, []byte(`{"type":"bad","channel":"test"`), nil).Once()

	connMock.On("ReadMessage").Return(1, []byte(""), &ws.CloseError{Code: ws.CloseUnsupportedData}).Once().WaitUntil(time.After(1 * time.Second))

	connMock.On("Close").Return(nil).Once()

	client := NewClient(context.Background(), connMock, serverMock)

	client.Listen()

	wgw.Wait()

	client.Close()

	connMock.AssertExpectations(t)
	serverMock.AssertExpectations(t)
}
*/

func TestProcessMessageWriteFail(t *testing.T) {
	var wg sync.WaitGroup

	serverMock := &MockServer{}

	connMock := &websocket.MockConnection{}

	connMock.On("SetWriteDeadline", mock.Anything).Return(nil)

	wg.Add(1)
	connMock.On("WriteMessage", 1, []byte(`{"type":"message","channel":"test","name":"foo","data":"bar"}`)).Return(errors.New("fail")).Run(func(args mock.Arguments) {
		defer wg.Done()
	})

	client := NewClient(context.Background(), connMock, serverMock)

	go client.processMessages()

	client.Write(&message.Event{
		Type:    message.EventTypeMessage,
		Channel: "test",
		Name:    "foo",
		Data:    json.RawMessage(`"bar"`),
	})

	wg.Wait()

	connMock.AssertExpectations(t)
	serverMock.AssertExpectations(t)
}

func TestProcessMessageClosedChannel(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)

	serverMock := &MockServer{}

	connMock := &websocket.MockConnection{}

	connMock.On("SetWriteDeadline", mock.Anything).Return(nil)

	connMock.On("WriteMessage", 8, []byte("")).Return(nil).Run(func(args mock.Arguments) {
		defer wg.Done()
	})

	client := NewClient(context.Background(), connMock, serverMock)

	go client.processMessages()

	close(client.messagesCh)

	wg.Wait()

	connMock.AssertExpectations(t)
	serverMock.AssertExpectations(t)
}
