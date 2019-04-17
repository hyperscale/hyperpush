// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
func TestServerAuthenticationProvider(t *testing.T) {
	authenticationMock := NewAuthenticationMock()

	server := NewServer(Configuration{})

	server.SetAuthenticationProvider(authenticationMock)

	assert.Equal(t, authenticationMock, server.authentication)
}
*/

func TestServerNotFoundUrl(t *testing.T) {
	server := NewServer(&Configuration{
		ClientQueueSize:         500,
		ChannelQueueSize:        500,
		AuthenticationQueueSize: 500,
		MessageQueueSize:        500,
		MaxConnections:          30000,
	}, nil)

	go server.Listen()

	ts := httptest.NewServer(server)
	defer ts.Close()

	_, err := http.Get(ts.URL)
	assert.NoError(t, err)
}

/*
func TestServer(t *testing.T) {
	token, err := makeAccessToken(testKey)
	assert.NoError(t, err)

	var wg sync.WaitGroup

	wg.Add(11)

	server := NewServer(Configuration{
		ClientQueueSize:         500,
		ChannelQueueSize:        500,
		AuthenticationQueueSize: 500,
		MessageQueueSize:        500,
		MaxConnections:          1,
	})
	server.SetAuthenticationProvider(NewAuthenticationJWT(authentication.Configuration{
		Key: testKey,
	}))

	go server.Listen()

	ts := httptest.NewServer(server)
	defer ts.Close()

	c, _, err := websocket.DefaultDialer.Dial(strings.Replace(ts.URL, "http://", "ws://", 1), nil)
	assert.NoError(t, err)

	go func() {
		defer c.Close()
		state := 0

		for {
			_, message, err := c.ReadMessage()
			if _, ok := err.(net.Error); ok {
				return
			}
			assert.NoError(t, err)

			event, err := message.Decode(message)
			assert.NoError(t, err)

			switch state {
			case 0:
				assert.Equal(t, message.EventTypeConnected, event.Type)
				state++
				wg.Done()
			case 1:
				assert.Equal(t, message.EventTypeSubscribed, event.Type)
				assert.Equal(t, "test", event.Channel)
				state++
				wg.Done()
			case 2:
				assert.Equal(t, message.EventTypeError, event.Type)
				assert.Equal(t, 400, event.Code)
				state++
				wg.Done()
			case 3:
				assert.Equal(t, message.EventTypeError, event.Type)
				assert.Equal(t, 400, event.Code)
				state++
				wg.Done()
			case 4:
				assert.Equal(t, message.EventTypeAuthenticated, event.Type)
				state++
				wg.Done()
			case 5:
				assert.Equal(t, message.EventTypeAuthenticated, event.Type)
				state++
				wg.Done()
			case 6:
				assert.Equal(t, message.EventTypeMessage, event.Type)
				assert.Equal(t, "test", event.Channel)
				assert.Equal(t, "foo", event.Name)
				assert.Equal(t, json.RawMessage(`"bar"`), event.Data)
				state++
				wg.Done()
			case 7:
				assert.Equal(t, message.EventTypeUnsubscribed, event.Type)
				assert.Equal(t, "test", event.Channel)
				state++
				wg.Done()
			case 8:
				assert.Equal(t, message.EventTypePong, event.Type)
				state++
				wg.Done()
			case 9:
				assert.Equal(t, message.EventTypeMessage, event.Type)
				assert.Equal(t, 351775, event.User)
				assert.Equal(t, "foo", event.Name)
				assert.Equal(t, json.RawMessage(`"bar"`), event.Data)
				state++
				wg.Done()
			case 10:
				assert.Equal(t, message.EventTypeError, event.Type)
				assert.Equal(t, 999, event.Code)
				state++
				wg.Done()
			}

			time.Sleep(time.Millisecond * 100)
		}
	}()

	time.Sleep(time.Millisecond * 500)

	assert.NoError(t, c.WriteMessage(websocket.TextMessage, []byte(`{"type":"subscribe","channel":"test"}`)))
	time.Sleep(time.Millisecond * 500)

	assert.NoError(t, c.WriteMessage(websocket.TextMessage, []byte(`{"type":"message","channel":"test","name":"foo","data":"bar"}`)))
	time.Sleep(time.Millisecond * 500)

	assert.NoError(t, c.WriteMessage(websocket.TextMessage, []byte(`{"type":"subscribe","channel":"private-test"}`)))
	time.Sleep(time.Millisecond * 500)

	assert.NoError(t, c.WriteMessage(websocket.TextMessage, []byte(`{"type":"authentication","token":"`+token+`"}`)))
	time.Sleep(time.Millisecond * 500)

	assert.NoError(t, c.WriteMessage(websocket.TextMessage, []byte(`{"type":"authentication","token":"`+token+`"}`)))
	time.Sleep(time.Millisecond * 500)

	assert.NoError(t, c.WriteMessage(websocket.TextMessage, []byte(`{"type":"message","channel":"test","name":"foo","data":"bar"}`)))
	time.Sleep(time.Millisecond * 500)

	assert.NoError(t, c.WriteMessage(websocket.TextMessage, []byte(`{"type":"unsubscribe","channel":"test"}`)))
	time.Sleep(time.Millisecond * 500)

	assert.NoError(t, c.WriteMessage(websocket.TextMessage, []byte(`{"type":"ping"}`)))
	time.Sleep(time.Millisecond * 500)

	assert.NoError(t, c.WriteMessage(websocket.TextMessage, []byte(`{"type":"message","user":351775,"name":"foo","data":"bar"}`)))
	time.Sleep(time.Millisecond * 500)

	assert.NoError(t, c.WriteMessage(websocket.TextMessage, []byte(`{"type":"authentication","token":"`+token+`55"}`)))

	resp, err := http.Get(ts.URL)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)

	wg.Wait()

	c.Close()

	time.Sleep(time.Second * 1)
}
*/

/*
func BenchmarkTest(b *testing.B) {
	b.SetParallelism(50)
	b.ReportAllocs()

	token, _ := makeAccessToken(testKey)

	server := NewServer(&Configuration{
		ClientQueueSize:         500,
		ChannelQueueSize:        500,
		AuthenticationQueueSize: 500,
		MessageQueueSize:        500,
		MaxConnections:          10000,
	}, nil)
	server.SetAuthenticationProvider(NewAuthenticationJWT(auth.Configuration{
		Key: testKey,
	}))

	go server.Listen()

	ts := httptest.NewServer(server)
	defer ts.Close()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c, _, err := websocket.DefaultDialer.Dial(strings.Replace(ts.URL, "http://", "ws://", 1), nil)
			if err != nil {
				//@TODO log
				continue
			}

			time.Sleep(time.Millisecond * 200)

			c.WriteMessage(websocket.TextMessage, []byte(`{"type":"subscribe","channel":"test"}`))
			time.Sleep(time.Millisecond * 500)

			c.WriteMessage(websocket.TextMessage, []byte(`{"type":"authentication","token":"`+token+`"}`))
			time.Sleep(time.Millisecond * 500)

			c.WriteMessage(websocket.TextMessage, []byte(`{"type":"message","channel":"test","name":"foo","data":"bar"}`))
			time.Sleep(time.Millisecond * 500)

			c.WriteMessage(websocket.TextMessage, []byte(`{"type":"unsubscribe","channel":"test"}`))
			time.Sleep(time.Millisecond * 500)

			c.WriteMessage(websocket.TextMessage, []byte(`{"type":"ping"}`))

			time.Sleep(time.Millisecond * 500)

			c.Close()
		}
	})
}
*/
