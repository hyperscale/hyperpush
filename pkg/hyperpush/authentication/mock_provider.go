// Code generated by mockery v1.0.0. DO NOT EDIT.

package authentication

import mock "github.com/stretchr/testify/mock"
import packets "github.com/hyperscale/hyperpush/pkg/hyperpush/mqtt/packets"

// MockProvider is an autogenerated mock type for the Provider type
type MockProvider struct {
	mock.Mock
}

// Authenticate provides a mock function with given fields: auth
func (_m *MockProvider) Authenticate(auth *packets.ConnectPacket) (*User, error) {
	ret := _m.Called(auth)

	var r0 *User
	if rf, ok := ret.Get(0).(func(*packets.ConnectPacket) *User); ok {
		r0 = rf(auth)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*packets.ConnectPacket) error); ok {
		r1 = rf(auth)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
