package backend

import (
	"context"
	"errors"
	"fmt"

	"github.com/onflow/flow-go/state"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/onflow/flow-go/access"
	"github.com/onflow/flow-go/cmd/build"
	"github.com/onflow/flow-go/engine/common/rpc/convert"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/state/protocol"
)

var SnapshotHistoryLimitErr = fmt.Errorf("reached the snapshot history limit")

type backendNetwork struct {
	state                protocol.State
	chainID              flow.ChainID
	snapshotHistoryLimit int
}

/*
NetworkAPI func

The observer and access nodes need to be able to handle GetNetworkParameters
and GetLatestProtocolStateSnapshot RPCs so this logic was split into
the backendNetwork so that we can ignore the rest of the backend logic
*/
func NewNetworkAPI(state protocol.State, chainID flow.ChainID, snapshotHistoryLimit int) *backendNetwork {
	return &backendNetwork{
		state:                state,
		chainID:              chainID,
		snapshotHistoryLimit: snapshotHistoryLimit,
	}
}

func (b *backendNetwork) GetNetworkParameters(_ context.Context) access.NetworkParameters {
	return access.NetworkParameters{
		ChainID: b.chainID,
	}
}

func (b *backendNetwork) GetNodeVersionInfo(ctx context.Context) (*access.NodeVersionInfo, error) {
	stateParams := b.state.Params()
	sporkId, err := stateParams.SporkID()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read spork ID: %v", err)
	}

	protocolVersion, err := stateParams.ProtocolVersion()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read protocol version: %v", err)
	}

	return &access.NodeVersionInfo{
		Semver:          build.Version(),
		Commit:          build.Commit(),
		SporkId:         sporkId,
		ProtocolVersion: uint64(protocolVersion),
	}, nil
}

// GetLatestProtocolStateSnapshot returns the latest finalized snapshot
func (b *backendNetwork) GetLatestProtocolStateSnapshot(_ context.Context) ([]byte, error) {
	snapshot := b.state.Final()

	validSnapshot, err := b.getValidSnapshot(snapshot, 0, true)
	if err != nil {
		return nil, err
	}

	data, err := convert.SnapshotToBytes(validSnapshot)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert snapshot to bytes: %v", err)
	}

	return data, nil
}

// GetProtocolStateSnapshotByBlockID returns serializable Snapshot for a block, by blockID.
// The requested block must be finalized, otherwise an error is returned.
// Expected errors during normal operation:
//   - status.Error[codes.NotFound] - No block with the given ID was found
//   - status.Error[codes.InvalidArgument] - We will never return a snapshot for this block ID:
//     1. A block was found, but it is not finalized and is below the finalized height, so it will never be finalized.
//     2. A block was found, however its sealing segment spans an epoch phase transition, yielding an invalid snapshot.
//   - status.Error[codes.FailedPrecondition] - A block was found, but it is not finalized and is above the finalized height.
//     The block may or may not be finalized in the future; the client can retry later.
func (b *backendNetwork) GetProtocolStateSnapshotByBlockID(_ context.Context, blockID flow.Identifier) ([]byte, error) {
	snapshot := b.state.AtBlockID(blockID)
	snapshotHeadByBlockId, err := snapshot.Head()
	if err != nil {
		if errors.Is(err, state.ErrUnknownSnapshotReference) {
			return nil, status.Errorf(codes.NotFound, "failed to get a valid snapshot: block not found")
		}

		return nil, status.Errorf(codes.Internal, "could not get header by blockID: %v", err)
	}

	snapshotByHeight := b.state.AtHeight(snapshotHeadByBlockId.Height)
	snapshotHeadByHeight, err := snapshotByHeight.Head()

	if err != nil {
		if errors.Is(err, state.ErrUnknownSnapshotReference) {
			return nil, status.Errorf(codes.InvalidArgument,
				"failed to retrieve snapshot for block by height %d: block not finalized", snapshotHeadByBlockId.Height)
		}

		return nil, status.Errorf(codes.Internal, "failed to find snapshot: %v", err)
	}

	if snapshotHeadByHeight.ID() != blockID {
		return nil, status.Errorf(codes.InvalidArgument, "failed to retrieve snapshot for block: block not finalized")
	}

	validSnapshot, err := b.getValidSnapshot(snapshotByHeight, 0, false)
	if err != nil {
		if errors.Is(err, ErrSnapshotPhaseMismatch) {
			return nil, status.Errorf(codes.InvalidArgument, "failed to retrieve snapshot for block, try again with different block: "+
				"%v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get a valid snapshot: %v", err)
	}

	data, err := convert.SnapshotToBytes(validSnapshot)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert snapshot to bytes: %v", err)
	}

	return data, nil
}

// GetProtocolStateSnapshotByHeight returns serializable Snapshot by block height.
// The block must be finalized (otherwise the by-height query is ambiguous).
// Expected errors during normal operation:
//   - status.Error[codes.NotFound] - No block with the given height was found.
//     The block height may or may not be finalized in the future; the client can retry later.
//   - status.Error[codes.InvalidArgument] - A block was found, however its sealing segment spans an epoch phase transition,
//     yielding an invalid snapshot. Therefore we will never return a snapshot for this block height.
func (b *backendNetwork) GetProtocolStateSnapshotByHeight(_ context.Context, blockHeight uint64) ([]byte, error) {
	snapshot := b.state.AtHeight(blockHeight)

	_, err := snapshot.Head()
	if err != nil {
		if errors.Is(err, state.ErrUnknownSnapshotReference) {
			return nil, status.Errorf(codes.NotFound, "failed to find snapshot: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to get a valid snapshot: %v", err)
	}

	validSnapshot, err := b.getValidSnapshot(snapshot, 0, false)
	if err != nil {
		if errors.Is(err, ErrSnapshotPhaseMismatch) {
			return nil, status.Errorf(codes.InvalidArgument, "failed to retrieve snapshot for block, try again with different block: "+
				"%v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get a valid snapshot: %v", err)
	}

	data, err := convert.SnapshotToBytes(validSnapshot)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert snapshot to bytes: %v", err)
	}

	return data, nil
}

func (b *backendNetwork) isEpochOrPhaseDifferent(counter1, counter2 uint64, phase1, phase2 flow.EpochPhase) bool {
	return counter1 != counter2 || phase1 != phase2
}

// getValidSnapshot will return a valid snapshot that has a sealing segment which
// 1. does not contain any blocks that span an epoch transition
// 2. does not contain any blocks that span an epoch phase transition
// If a snapshot does contain an invalid sealing segment query the state
// by height of each block in the segment and return a snapshot at the point
// where the transition happens.
// Expected error returns during normal operations:
// * ErrSnapshotPhaseMismatch - snapshot does not contain a valid sealing segment
// All other errors should be treated as exceptions.
func (b *backendNetwork) getValidSnapshot(snapshot protocol.Snapshot, blocksVisited int, findNextValidSnapshot bool) (protocol.Snapshot, error) {
	segment, err := snapshot.SealingSegment()
	if err != nil {
		return nil, fmt.Errorf("failed to get sealing segment: %w", err)
	}

	counterAtHighest, phaseAtHighest, err := b.getCounterAndPhase(segment.Highest().Header.Height)
	if err != nil {
		return nil, fmt.Errorf("failed to get counter and phase at highest block in the segment: %w", err)
	}

	counterAtLowest, phaseAtLowest, err := b.getCounterAndPhase(segment.Sealed().Header.Height)
	if err != nil {
		return nil, fmt.Errorf("failed to get counter and phase at lowest block in the segment: %w", err)
	}

	// Check if the counters and phase are different this indicates that the sealing segment
	// of the snapshot requested spans either an epoch transition or phase transition.
	if b.isEpochOrPhaseDifferent(counterAtHighest, counterAtLowest, phaseAtHighest, phaseAtLowest) {
		if !findNextValidSnapshot {
			return nil, ErrSnapshotPhaseMismatch
		}

		// Visit each node in strict order of decreasing height starting at head
		// to find the block that straddles the transition boundary.
		for i := len(segment.Blocks) - 1; i >= 0; i-- {
			blocksVisited++

			// NOTE: Check if we have reached our history limit, in edge cases
			// where the sealing segment is abnormally long we want to short circuit
			// the recursive calls and return an error. The API caller can retry.
			if blocksVisited > b.snapshotHistoryLimit {
				return nil, fmt.Errorf("%w: (%d)", SnapshotHistoryLimitErr, b.snapshotHistoryLimit)
			}

			counterAtBlock, phaseAtBlock, err := b.getCounterAndPhase(segment.Blocks[i].Header.Height)
			if err != nil {
				return nil, fmt.Errorf("failed to get epoch counter and phase for snapshot at block %s: %w", segment.Blocks[i].ID(), err)
			}

			// Check if this block straddles the transition boundary, if it does return the snapshot
			// at that block height.
			if b.isEpochOrPhaseDifferent(counterAtHighest, counterAtBlock, phaseAtHighest, phaseAtBlock) {
				return b.getValidSnapshot(b.state.AtHeight(segment.Blocks[i].Header.Height), blocksVisited, true)
			}
		}
	}

	return snapshot, nil
}

// getCounterAndPhase will return the epoch counter and phase at the specified height in state
func (b *backendNetwork) getCounterAndPhase(height uint64) (uint64, flow.EpochPhase, error) {
	snapshot := b.state.AtHeight(height)

	counter, err := snapshot.Epochs().Current().Counter()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get counter for block (height=%d): %w", height, err)
	}

	phase, err := snapshot.Phase()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get phase for block (height=%d): %w", height, err)
	}

	return counter, phase, nil
}
