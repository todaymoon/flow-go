// Code generated by mockery v2.21.4. DO NOT EDIT.

package mock

import (
	context "context"

	access "github.com/onflow/flow-go/access"

	execution "github.com/onflow/flow/protobuf/go/flow/execution"

	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"
)

// API is an autogenerated mock type for the API type
type API struct {
	mock.Mock
}

// ExecuteScriptAtBlockHeight provides a mock function with given fields: ctx, blockHeight, script, arguments
func (_m *API) ExecuteScriptAtBlockHeight(ctx context.Context, blockHeight uint64, script []byte, arguments [][]byte) ([]byte, error) {
	ret := _m.Called(ctx, blockHeight, script, arguments)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, []byte, [][]byte) ([]byte, error)); ok {
		return rf(ctx, blockHeight, script, arguments)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, []byte, [][]byte) []byte); ok {
		r0 = rf(ctx, blockHeight, script, arguments)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, []byte, [][]byte) error); ok {
		r1 = rf(ctx, blockHeight, script, arguments)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExecuteScriptAtBlockID provides a mock function with given fields: ctx, blockID, script, arguments
func (_m *API) ExecuteScriptAtBlockID(ctx context.Context, blockID flow.Identifier, script []byte, arguments [][]byte) ([]byte, error) {
	ret := _m.Called(ctx, blockID, script, arguments)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier, []byte, [][]byte) ([]byte, error)); ok {
		return rf(ctx, blockID, script, arguments)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier, []byte, [][]byte) []byte); ok {
		r0 = rf(ctx, blockID, script, arguments)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier, []byte, [][]byte) error); ok {
		r1 = rf(ctx, blockID, script, arguments)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExecuteScriptAtLatestBlock provides a mock function with given fields: ctx, script, arguments
func (_m *API) ExecuteScriptAtLatestBlock(ctx context.Context, script []byte, arguments [][]byte) ([]byte, error) {
	ret := _m.Called(ctx, script, arguments)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte, [][]byte) ([]byte, error)); ok {
		return rf(ctx, script, arguments)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []byte, [][]byte) []byte); ok {
		r0 = rf(ctx, script, arguments)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []byte, [][]byte) error); ok {
		r1 = rf(ctx, script, arguments)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccount provides a mock function with given fields: ctx, address
func (_m *API) GetAccount(ctx context.Context, address flow.Address) (*flow.Account, error) {
	ret := _m.Called(ctx, address)

	var r0 *flow.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Address) (*flow.Account, error)); ok {
		return rf(ctx, address)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Address) *flow.Account); ok {
		r0 = rf(ctx, address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Address) error); ok {
		r1 = rf(ctx, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccountAtBlockHeight provides a mock function with given fields: ctx, address, height
func (_m *API) GetAccountAtBlockHeight(ctx context.Context, address flow.Address, height uint64) (*flow.Account, error) {
	ret := _m.Called(ctx, address, height)

	var r0 *flow.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Address, uint64) (*flow.Account, error)); ok {
		return rf(ctx, address, height)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Address, uint64) *flow.Account); ok {
		r0 = rf(ctx, address, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Address, uint64) error); ok {
		r1 = rf(ctx, address, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccountAtLatestBlock provides a mock function with given fields: ctx, address
func (_m *API) GetAccountAtLatestBlock(ctx context.Context, address flow.Address) (*flow.Account, error) {
	ret := _m.Called(ctx, address)

	var r0 *flow.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Address) (*flow.Account, error)); ok {
		return rf(ctx, address)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Address) *flow.Account); ok {
		r0 = rf(ctx, address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Address) error); ok {
		r1 = rf(ctx, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockByHeight provides a mock function with given fields: ctx, height
func (_m *API) GetBlockByHeight(ctx context.Context, height uint64) (*flow.Block, flow.BlockStatus, error) {
	ret := _m.Called(ctx, height)

	var r0 *flow.Block
	var r1 flow.BlockStatus
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) (*flow.Block, flow.BlockStatus, error)); ok {
		return rf(ctx, height)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *flow.Block); ok {
		r0 = rf(ctx, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Block)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64) flow.BlockStatus); ok {
		r1 = rf(ctx, height)
	} else {
		r1 = ret.Get(1).(flow.BlockStatus)
	}

	if rf, ok := ret.Get(2).(func(context.Context, uint64) error); ok {
		r2 = rf(ctx, height)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetBlockByID provides a mock function with given fields: ctx, id
func (_m *API) GetBlockByID(ctx context.Context, id flow.Identifier) (*flow.Block, flow.BlockStatus, error) {
	ret := _m.Called(ctx, id)

	var r0 *flow.Block
	var r1 flow.BlockStatus
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) (*flow.Block, flow.BlockStatus, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) *flow.Block); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Block)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) flow.BlockStatus); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Get(1).(flow.BlockStatus)
	}

	if rf, ok := ret.Get(2).(func(context.Context, flow.Identifier) error); ok {
		r2 = rf(ctx, id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetBlockHeaderByHeight provides a mock function with given fields: ctx, height
func (_m *API) GetBlockHeaderByHeight(ctx context.Context, height uint64) (*flow.Header, flow.BlockStatus, error) {
	ret := _m.Called(ctx, height)

	var r0 *flow.Header
	var r1 flow.BlockStatus
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) (*flow.Header, flow.BlockStatus, error)); ok {
		return rf(ctx, height)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *flow.Header); ok {
		r0 = rf(ctx, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Header)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64) flow.BlockStatus); ok {
		r1 = rf(ctx, height)
	} else {
		r1 = ret.Get(1).(flow.BlockStatus)
	}

	if rf, ok := ret.Get(2).(func(context.Context, uint64) error); ok {
		r2 = rf(ctx, height)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetBlockHeaderByID provides a mock function with given fields: ctx, id
func (_m *API) GetBlockHeaderByID(ctx context.Context, id flow.Identifier) (*flow.Header, flow.BlockStatus, error) {
	ret := _m.Called(ctx, id)

	var r0 *flow.Header
	var r1 flow.BlockStatus
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) (*flow.Header, flow.BlockStatus, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) *flow.Header); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Header)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) flow.BlockStatus); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Get(1).(flow.BlockStatus)
	}

	if rf, ok := ret.Get(2).(func(context.Context, flow.Identifier) error); ok {
		r2 = rf(ctx, id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetCollectionByID provides a mock function with given fields: ctx, id
func (_m *API) GetCollectionByID(ctx context.Context, id flow.Identifier) (*flow.LightCollection, error) {
	ret := _m.Called(ctx, id)

	var r0 *flow.LightCollection
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) (*flow.LightCollection, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) *flow.LightCollection); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.LightCollection)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEventsForBlockIDs provides a mock function with given fields: ctx, eventType, blockIDs, eventEncodingVersion
func (_m *API) GetEventsForBlockIDs(ctx context.Context, eventType string, blockIDs []flow.Identifier, eventEncodingVersion execution.EventEncodingVersion) ([]flow.BlockEvents, error) {
	ret := _m.Called(ctx, eventType, blockIDs, eventEncodingVersion)

	var r0 []flow.BlockEvents
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []flow.Identifier, execution.EventEncodingVersion) ([]flow.BlockEvents, error)); ok {
		return rf(ctx, eventType, blockIDs, eventEncodingVersion)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, []flow.Identifier, execution.EventEncodingVersion) []flow.BlockEvents); ok {
		r0 = rf(ctx, eventType, blockIDs, eventEncodingVersion)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]flow.BlockEvents)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, []flow.Identifier, execution.EventEncodingVersion) error); ok {
		r1 = rf(ctx, eventType, blockIDs, eventEncodingVersion)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEventsForHeightRange provides a mock function with given fields: ctx, eventType, startHeight, endHeight, eventEncodingVersion
func (_m *API) GetEventsForHeightRange(ctx context.Context, eventType string, startHeight uint64, endHeight uint64, eventEncodingVersion execution.EventEncodingVersion) ([]flow.BlockEvents, error) {
	ret := _m.Called(ctx, eventType, startHeight, endHeight, eventEncodingVersion)

	var r0 []flow.BlockEvents
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, uint64, uint64, execution.EventEncodingVersion) ([]flow.BlockEvents, error)); ok {
		return rf(ctx, eventType, startHeight, endHeight, eventEncodingVersion)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, uint64, uint64, execution.EventEncodingVersion) []flow.BlockEvents); ok {
		r0 = rf(ctx, eventType, startHeight, endHeight, eventEncodingVersion)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]flow.BlockEvents)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, uint64, uint64, execution.EventEncodingVersion) error); ok {
		r1 = rf(ctx, eventType, startHeight, endHeight, eventEncodingVersion)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetExecutionResultByID provides a mock function with given fields: ctx, id
func (_m *API) GetExecutionResultByID(ctx context.Context, id flow.Identifier) (*flow.ExecutionResult, error) {
	ret := _m.Called(ctx, id)

	var r0 *flow.ExecutionResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) (*flow.ExecutionResult, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) *flow.ExecutionResult); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.ExecutionResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetExecutionResultForBlockID provides a mock function with given fields: ctx, blockID
func (_m *API) GetExecutionResultForBlockID(ctx context.Context, blockID flow.Identifier) (*flow.ExecutionResult, error) {
	ret := _m.Called(ctx, blockID)

	var r0 *flow.ExecutionResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) (*flow.ExecutionResult, error)); ok {
		return rf(ctx, blockID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) *flow.ExecutionResult); ok {
		r0 = rf(ctx, blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.ExecutionResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) error); ok {
		r1 = rf(ctx, blockID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLatestBlock provides a mock function with given fields: ctx, isSealed
func (_m *API) GetLatestBlock(ctx context.Context, isSealed bool) (*flow.Block, flow.BlockStatus, error) {
	ret := _m.Called(ctx, isSealed)

	var r0 *flow.Block
	var r1 flow.BlockStatus
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, bool) (*flow.Block, flow.BlockStatus, error)); ok {
		return rf(ctx, isSealed)
	}
	if rf, ok := ret.Get(0).(func(context.Context, bool) *flow.Block); ok {
		r0 = rf(ctx, isSealed)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Block)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, bool) flow.BlockStatus); ok {
		r1 = rf(ctx, isSealed)
	} else {
		r1 = ret.Get(1).(flow.BlockStatus)
	}

	if rf, ok := ret.Get(2).(func(context.Context, bool) error); ok {
		r2 = rf(ctx, isSealed)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetLatestBlockHeader provides a mock function with given fields: ctx, isSealed
func (_m *API) GetLatestBlockHeader(ctx context.Context, isSealed bool) (*flow.Header, flow.BlockStatus, error) {
	ret := _m.Called(ctx, isSealed)

	var r0 *flow.Header
	var r1 flow.BlockStatus
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, bool) (*flow.Header, flow.BlockStatus, error)); ok {
		return rf(ctx, isSealed)
	}
	if rf, ok := ret.Get(0).(func(context.Context, bool) *flow.Header); ok {
		r0 = rf(ctx, isSealed)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Header)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, bool) flow.BlockStatus); ok {
		r1 = rf(ctx, isSealed)
	} else {
		r1 = ret.Get(1).(flow.BlockStatus)
	}

	if rf, ok := ret.Get(2).(func(context.Context, bool) error); ok {
		r2 = rf(ctx, isSealed)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetLatestProtocolStateSnapshot provides a mock function with given fields: ctx
func (_m *API) GetLatestProtocolStateSnapshot(ctx context.Context) ([]byte, error) {
	ret := _m.Called(ctx)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]byte, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []byte); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNetworkParameters provides a mock function with given fields: ctx
func (_m *API) GetNetworkParameters(ctx context.Context) access.NetworkParameters {
	ret := _m.Called(ctx)

	var r0 access.NetworkParameters
	if rf, ok := ret.Get(0).(func(context.Context) access.NetworkParameters); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(access.NetworkParameters)
	}

	return r0
}

// GetNodeVersionInfo provides a mock function with given fields: ctx
func (_m *API) GetNodeVersionInfo(ctx context.Context) (*access.NodeVersionInfo, error) {
	ret := _m.Called(ctx)

	var r0 *access.NodeVersionInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*access.NodeVersionInfo, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *access.NodeVersionInfo); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*access.NodeVersionInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransaction provides a mock function with given fields: ctx, id
func (_m *API) GetTransaction(ctx context.Context, id flow.Identifier) (*flow.TransactionBody, error) {
	ret := _m.Called(ctx, id)

	var r0 *flow.TransactionBody
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) (*flow.TransactionBody, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) *flow.TransactionBody); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.TransactionBody)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionResult provides a mock function with given fields: ctx, id, blockID, collectionID, eventEncodingVersion
func (_m *API) GetTransactionResult(ctx context.Context, id flow.Identifier, blockID flow.Identifier, collectionID flow.Identifier, eventEncodingVersion execution.EventEncodingVersion) (*access.TransactionResult, error) {
	ret := _m.Called(ctx, id, blockID, collectionID, eventEncodingVersion)

	var r0 *access.TransactionResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier, flow.Identifier, flow.Identifier, execution.EventEncodingVersion) (*access.TransactionResult, error)); ok {
		return rf(ctx, id, blockID, collectionID, eventEncodingVersion)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier, flow.Identifier, flow.Identifier, execution.EventEncodingVersion) *access.TransactionResult); ok {
		r0 = rf(ctx, id, blockID, collectionID, eventEncodingVersion)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*access.TransactionResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier, flow.Identifier, flow.Identifier, execution.EventEncodingVersion) error); ok {
		r1 = rf(ctx, id, blockID, collectionID, eventEncodingVersion)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionResultByIndex provides a mock function with given fields: ctx, blockID, index, eventEncodingVersion
func (_m *API) GetTransactionResultByIndex(ctx context.Context, blockID flow.Identifier, index uint32, eventEncodingVersion execution.EventEncodingVersion) (*access.TransactionResult, error) {
	ret := _m.Called(ctx, blockID, index, eventEncodingVersion)

	var r0 *access.TransactionResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier, uint32, execution.EventEncodingVersion) (*access.TransactionResult, error)); ok {
		return rf(ctx, blockID, index, eventEncodingVersion)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier, uint32, execution.EventEncodingVersion) *access.TransactionResult); ok {
		r0 = rf(ctx, blockID, index, eventEncodingVersion)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*access.TransactionResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier, uint32, execution.EventEncodingVersion) error); ok {
		r1 = rf(ctx, blockID, index, eventEncodingVersion)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionResultsByBlockID provides a mock function with given fields: ctx, blockID, eventEncodingVersion
func (_m *API) GetTransactionResultsByBlockID(ctx context.Context, blockID flow.Identifier, eventEncodingVersion execution.EventEncodingVersion) ([]*access.TransactionResult, error) {
	ret := _m.Called(ctx, blockID, eventEncodingVersion)

	var r0 []*access.TransactionResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier, execution.EventEncodingVersion) ([]*access.TransactionResult, error)); ok {
		return rf(ctx, blockID, eventEncodingVersion)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier, execution.EventEncodingVersion) []*access.TransactionResult); ok {
		r0 = rf(ctx, blockID, eventEncodingVersion)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*access.TransactionResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier, execution.EventEncodingVersion) error); ok {
		r1 = rf(ctx, blockID, eventEncodingVersion)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionsByBlockID provides a mock function with given fields: ctx, blockID
func (_m *API) GetTransactionsByBlockID(ctx context.Context, blockID flow.Identifier) ([]*flow.TransactionBody, error) {
	ret := _m.Called(ctx, blockID)

	var r0 []*flow.TransactionBody
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) ([]*flow.TransactionBody, error)); ok {
		return rf(ctx, blockID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) []*flow.TransactionBody); ok {
		r0 = rf(ctx, blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*flow.TransactionBody)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) error); ok {
		r1 = rf(ctx, blockID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Ping provides a mock function with given fields: ctx
func (_m *API) Ping(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendTransaction provides a mock function with given fields: ctx, tx
func (_m *API) SendTransaction(ctx context.Context, tx *flow.TransactionBody) error {
	ret := _m.Called(ctx, tx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *flow.TransactionBody) error); ok {
		r0 = rf(ctx, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewAPI interface {
	mock.TestingT
	Cleanup(func())
}

// NewAPI creates a new instance of API. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAPI(t mockConstructorTestingTNewAPI) *API {
	mock := &API{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
