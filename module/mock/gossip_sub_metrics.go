// Code generated by mockery v2.21.4. DO NOT EDIT.

package mock

import (
	channels "github.com/onflow/flow-go/network/channels"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// GossipSubMetrics is an autogenerated mock type for the GossipSubMetrics type
type GossipSubMetrics struct {
	mock.Mock
}

// AsyncProcessingFinished provides a mock function with given fields: msgType, duration
func (_m *GossipSubMetrics) AsyncProcessingFinished(msgType string, duration time.Duration) {
	_m.Called(msgType, duration)
}

// AsyncProcessingStarted provides a mock function with given fields: msgType
func (_m *GossipSubMetrics) AsyncProcessingStarted(msgType string) {
	_m.Called(msgType)
}

// OnAppSpecificScoreUpdated provides a mock function with given fields: _a0
func (_m *GossipSubMetrics) OnAppSpecificScoreUpdated(_a0 float64) {
	_m.Called(_a0)
}

// OnBehaviourPenaltyUpdated provides a mock function with given fields: _a0
func (_m *GossipSubMetrics) OnBehaviourPenaltyUpdated(_a0 float64) {
	_m.Called(_a0)
}

// OnFirstMessageDeliveredUpdated provides a mock function with given fields: _a0, _a1
func (_m *GossipSubMetrics) OnFirstMessageDeliveredUpdated(_a0 channels.Topic, _a1 float64) {
	_m.Called(_a0, _a1)
}

// OnGraftReceived provides a mock function with given fields: count
func (_m *GossipSubMetrics) OnGraftReceived(count int) {
	_m.Called(count)
}

// OnIHaveReceived provides a mock function with given fields: count
func (_m *GossipSubMetrics) OnIHaveReceived(count int) {
	_m.Called(count)
}

// OnIPColocationFactorUpdated provides a mock function with given fields: _a0
func (_m *GossipSubMetrics) OnIPColocationFactorUpdated(_a0 float64) {
	_m.Called(_a0)
}

// OnIWantReceived provides a mock function with given fields: count
func (_m *GossipSubMetrics) OnIWantReceived(count int) {
	_m.Called(count)
}

// OnIncomingRpcAcceptedFully provides a mock function with given fields:
func (_m *GossipSubMetrics) OnIncomingRpcAcceptedFully() {
	_m.Called()
}

// OnIncomingRpcAcceptedOnlyForControlMessages provides a mock function with given fields:
func (_m *GossipSubMetrics) OnIncomingRpcAcceptedOnlyForControlMessages() {
	_m.Called()
}

// OnIncomingRpcRejected provides a mock function with given fields:
func (_m *GossipSubMetrics) OnIncomingRpcRejected() {
	_m.Called()
}

// OnInvalidMessageDeliveredUpdated provides a mock function with given fields: _a0, _a1
func (_m *GossipSubMetrics) OnInvalidMessageDeliveredUpdated(_a0 channels.Topic, _a1 float64) {
	_m.Called(_a0, _a1)
}

// OnLocalMeshSizeUpdated provides a mock function with given fields: topic, size
func (_m *GossipSubMetrics) OnLocalMeshSizeUpdated(topic string, size int) {
	_m.Called(topic, size)
}

// OnMeshMessageDeliveredUpdated provides a mock function with given fields: _a0, _a1
func (_m *GossipSubMetrics) OnMeshMessageDeliveredUpdated(_a0 channels.Topic, _a1 float64) {
	_m.Called(_a0, _a1)
}

// OnOverallPeerScoreUpdated provides a mock function with given fields: _a0
func (_m *GossipSubMetrics) OnOverallPeerScoreUpdated(_a0 float64) {
	_m.Called(_a0)
}

// OnPruneReceived provides a mock function with given fields: count
func (_m *GossipSubMetrics) OnPruneReceived(count int) {
	_m.Called(count)
}

// OnPublishedGossipMessagesReceived provides a mock function with given fields: count
func (_m *GossipSubMetrics) OnPublishedGossipMessagesReceived(count int) {
	_m.Called(count)
}

// OnTimeInMeshUpdated provides a mock function with given fields: _a0, _a1
func (_m *GossipSubMetrics) OnTimeInMeshUpdated(_a0 channels.Topic, _a1 time.Duration) {
	_m.Called(_a0, _a1)
}

// PreProcessingFinished provides a mock function with given fields: msgType, sampleSize, duration
func (_m *GossipSubMetrics) PreProcessingFinished(msgType string, sampleSize uint, duration time.Duration) {
	_m.Called(msgType, sampleSize, duration)
}

// PreProcessingStarted provides a mock function with given fields: msgType, sampleSize
func (_m *GossipSubMetrics) PreProcessingStarted(msgType string, sampleSize uint) {
	_m.Called(msgType, sampleSize)
}

// SetWarningStateCount provides a mock function with given fields: _a0
func (_m *GossipSubMetrics) SetWarningStateCount(_a0 uint) {
	_m.Called(_a0)
}

type mockConstructorTestingTNewGossipSubMetrics interface {
	mock.TestingT
	Cleanup(func())
}

// NewGossipSubMetrics creates a new instance of GossipSubMetrics. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGossipSubMetrics(t mockConstructorTestingTNewGossipSubMetrics) *GossipSubMetrics {
	mock := &GossipSubMetrics{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
