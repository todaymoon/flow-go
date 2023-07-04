// Code generated by mockery v2.21.4. DO NOT EDIT.

package mocks

import (
	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"

	model "github.com/onflow/flow-go/consensus/hotstuff/model"
)

// Validator is an autogenerated mock type for the Validator type
type Validator struct {
	mock.Mock
}

// ValidateProposal provides a mock function with given fields: proposal
func (_m *Validator) ValidateProposal(proposal *model.Proposal) error {
	ret := _m.Called(proposal)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Proposal) error); ok {
		r0 = rf(proposal)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateQC provides a mock function with given fields: qc
func (_m *Validator) ValidateQC(qc *flow.QuorumCertificate) error {
	ret := _m.Called(qc)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.QuorumCertificate) error); ok {
		r0 = rf(qc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateTC provides a mock function with given fields: tc
func (_m *Validator) ValidateTC(tc *flow.TimeoutCertificate) error {
	ret := _m.Called(tc)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.TimeoutCertificate) error); ok {
		r0 = rf(tc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateVote provides a mock function with given fields: vote
func (_m *Validator) ValidateVote(vote *model.Vote) (*flow.IdentitySkeleton, error) {
	ret := _m.Called(vote)

	var r0 *flow.Identity
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.Vote) (*flow.Identity, error)); ok {
		return rf(vote)
	}
	if rf, ok := ret.Get(0).(func(*model.Vote) *flow.Identity); ok {
		r0 = rf(vote)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Identity)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.Vote) error); ok {
		r1 = rf(vote)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewValidator interface {
	mock.TestingT
	Cleanup(func())
}

// NewValidator creates a new instance of Validator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewValidator(t mockConstructorTestingTNewValidator) *Validator {
	mock := &Validator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
