// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import mock "github.com/stretchr/testify/mock"

// GossipSubLocalMeshMetrics is an autogenerated mock type for the GossipSubLocalMeshMetrics type
type GossipSubLocalMeshMetrics struct {
	mock.Mock
}

// OnLocalMeshSizeUpdated provides a mock function with given fields: topic, size
func (_m *GossipSubLocalMeshMetrics) OnLocalMeshSizeUpdated(topic string, size int) {
	_m.Called(topic, size)
}

type mockConstructorTestingTNewGossipSubLocalMeshMetrics interface {
	mock.TestingT
	Cleanup(func())
}

// NewGossipSubLocalMeshMetrics creates a new instance of GossipSubLocalMeshMetrics. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGossipSubLocalMeshMetrics(t mockConstructorTestingTNewGossipSubLocalMeshMetrics) *GossipSubLocalMeshMetrics {
	mock := &GossipSubLocalMeshMetrics{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
