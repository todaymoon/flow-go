// Code generated by mockery v2.13.0. DO NOT EDIT.

package mock

import (
	context "context"

	flow "github.com/onflow/flow-go/model/flow"
	execution_data "github.com/onflow/flow-go/module/executiondatasync/execution_data"

	mock "github.com/stretchr/testify/mock"
)

// Downloader is an autogenerated mock type for the Downloader type
type Downloader struct {
	mock.Mock
}

// Done provides a mock function with given fields:
func (_m *Downloader) Done() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// Download provides a mock function with given fields: ctx, executionDataID
func (_m *Downloader) Download(ctx context.Context, executionDataID flow.Identifier) (*execution_data.BlockExecutionData, error) {
	ret := _m.Called(ctx, executionDataID)

	var r0 *execution_data.BlockExecutionData
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) *execution_data.BlockExecutionData); ok {
		r0 = rf(ctx, executionDataID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*execution_data.BlockExecutionData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) error); ok {
		r1 = rf(ctx, executionDataID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Ready provides a mock function with given fields:
func (_m *Downloader) Ready() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

type NewDownloaderT interface {
	mock.TestingT
	Cleanup(func())
}

// NewDownloader creates a new instance of Downloader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDownloader(t NewDownloaderT) *Downloader {
	mock := &Downloader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
