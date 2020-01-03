// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import flow "github.com/dapperlabs/flow-go/model/flow"
import mock "github.com/stretchr/testify/mock"
import state "github.com/dapperlabs/flow-go/engine/execution/execution/state"

// ExecutionState is an autogenerated mock type for the ExecutionState type
type ExecutionState struct {
	mock.Mock
}

// CommitDelta provides a mock function with given fields: _a0
func (_m *ExecutionState) CommitDelta(_a0 state.Delta) (flow.StateCommitment, error) {
	ret := _m.Called(_a0)

	var r0 flow.StateCommitment
	if rf, ok := ret.Get(0).(func(state.Delta) flow.StateCommitment); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(flow.StateCommitment)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(state.Delta) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewView provides a mock function with given fields: _a0
func (_m *ExecutionState) NewView(_a0 flow.StateCommitment) *state.View {
	ret := _m.Called(_a0)

	var r0 *state.View
	if rf, ok := ret.Get(0).(func(flow.StateCommitment) *state.View); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*state.View)
		}
	}

	return r0
}
