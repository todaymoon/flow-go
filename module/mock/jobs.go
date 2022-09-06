// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import (
	module "github.com/onflow/flow-go/module"
	mock "github.com/stretchr/testify/mock"
)

// Jobs is an autogenerated mock type for the Jobs type
type Jobs struct {
	mock.Mock
}

// AtIndex provides a mock function with given fields: index
func (_m *Jobs) AtIndex(index uint64) (module.Job, error) {
	ret := _m.Called(index)

	var r0 module.Job
	if rf, ok := ret.Get(0).(func(uint64) module.Job); ok {
		r0 = rf(index)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(module.Job)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(index)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Head provides a mock function with given fields:
func (_m *Jobs) Head() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewJobs interface {
	mock.TestingT
	Cleanup(func())
}

// NewJobs creates a new instance of Jobs. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewJobs(t mockConstructorTestingTNewJobs) *Jobs {
	mock := &Jobs{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
