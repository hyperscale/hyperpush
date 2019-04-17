// Code generated by mockery v1.0.0. DO NOT EDIT.

package message

import mock "github.com/stretchr/testify/mock"

// MockErrorInterface is an autogenerated mock type for the ErrorInterface type
type MockErrorInterface struct {
	mock.Mock
}

// Code provides a mock function with given fields:
func (_m *MockErrorInterface) Code() ErrorCode {
	ret := _m.Called()

	var r0 ErrorCode
	if rf, ok := ret.Get(0).(func() ErrorCode); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(ErrorCode)
	}

	return r0
}

// Error provides a mock function with given fields:
func (_m *MockErrorInterface) Error() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
