// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import (
	crypto "github.com/onflow/crypto"
	mock "github.com/stretchr/testify/mock"
)

// EpochRecoveryMyBeaconKey is an autogenerated mock type for the EpochRecoveryMyBeaconKey type
type EpochRecoveryMyBeaconKey struct {
	mock.Mock
}

// OverwriteMyBeaconPrivateKey provides a mock function with given fields: epochCounter, key
func (_m *EpochRecoveryMyBeaconKey) OverwriteMyBeaconPrivateKey(epochCounter uint64, key crypto.PrivateKey) error {
	ret := _m.Called(epochCounter, key)

	if len(ret) == 0 {
		panic("no return value specified for OverwriteMyBeaconPrivateKey")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64, crypto.PrivateKey) error); ok {
		r0 = rf(epochCounter, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RetrieveMyBeaconPrivateKey provides a mock function with given fields: epochCounter
func (_m *EpochRecoveryMyBeaconKey) RetrieveMyBeaconPrivateKey(epochCounter uint64) (crypto.PrivateKey, bool, error) {
	ret := _m.Called(epochCounter)

	if len(ret) == 0 {
		panic("no return value specified for RetrieveMyBeaconPrivateKey")
	}

	var r0 crypto.PrivateKey
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(uint64) (crypto.PrivateKey, bool, error)); ok {
		return rf(epochCounter)
	}
	if rf, ok := ret.Get(0).(func(uint64) crypto.PrivateKey); ok {
		r0 = rf(epochCounter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(crypto.PrivateKey)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64) bool); ok {
		r1 = rf(epochCounter)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(uint64) error); ok {
		r2 = rf(epochCounter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewEpochRecoveryMyBeaconKey creates a new instance of EpochRecoveryMyBeaconKey. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEpochRecoveryMyBeaconKey(t interface {
	mock.TestingT
	Cleanup(func())
}) *EpochRecoveryMyBeaconKey {
	mock := &EpochRecoveryMyBeaconKey{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
