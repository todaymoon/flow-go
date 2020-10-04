package execution_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/engine/execution/state/delta"
	"github.com/onflow/flow-go/engine/testutil"
	testmock "github.com/onflow/flow-go/engine/testutil/mock"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module/mempool/entity"
	"github.com/onflow/flow-go/network/stub"
	"github.com/onflow/flow-go/utils/unittest"
)

// Test execution node can sync execution state from another EN node.
//
// It creates 2 ENs, and 5 blocks:
// A <- B
//      ^- C (seal_A)
//         ^- D (seal_B, finalized)
//      		  ^- E <- F <- G (finalizes D, seals B)
//
// set the sync threshold as 2, meaning requires 2 sealed and unexecuted
// blocks to trigger state syncing
// set the block execution for EN1 to be fast
// set the block execution for EN2 to be slow
//
// EN1 will have A,B,C,D,E,F,G
// EN2 will receive A,B,C,D,E,F and verify no state syncing is triggered.
// Once EN2 has received G, state syncing will be triggered. verify that
// EN2 has the statecommitment for D
func TestStateSyncFlow(t *testing.T) {
	// create two EN nodes,
	// EN1 is able to execute blocks fast,
	// EN2 is slow to execute any block, it has to rely on state syncing
	// to catch up.
	withNodes(t, func(EN1, EN2 *testmock.ExecutionNode) {
		log := unittest.Logger()
		log.Debug().Msg("nodes created")

		EN2.ExecutionEngine.OnComputeBlock = func(ctx context.Context, block *entity.ExecutableBlock, view *delta.View) {
			log.Info().Msgf("EN2 is about to compute block: %v, let it be slow...", block.ID())
			time.Sleep(time.Second * 5)
		}

		genesis, err := EN1.State.AtHeight(0).Head()
		require.NoError(t, err, "could not get genesis")

		// create all the blocks needed for this tests:
		// A <- B <- C <- D <- E <- F <- G
		blockA := makeBlockWithParentAndSeal(genesis, nil)
		blockB := makeBlockWithParentAndSeal(blockA.Header, nil)
		blockC := makeBlockWithParentAndSeal(blockB.Header, blockA.Header)
		blockD := makeBlockWithParentAndSeal(blockC.Header, blockB.Header)
		blockE := makeBlockWithParentAndSeal(blockD.Header, nil)
		blockF := makeBlockWithParentAndSeal(blockE.Header, nil)
		blockG := makeBlockWithParentAndSeal(blockF.Header, nil)

		logBlocks(log, blockA, blockB, blockC, blockD, blockE, blockF, blockG)

		// EN1 receives all the blocks and we wait until EN1 has considered
		// blockB is sealed.
		sendBlockToEN(t, blockA, nil, EN1)
		sendBlockToEN(t, blockB, nil, EN1)
		sendBlockToEN(t, blockC, nil, EN1)
		sendBlockToEN(t, blockD, blockA, EN1)
		sendBlockToEN(t, blockE, blockB, EN1)
		sendBlockToEN(t, blockF, blockC, EN1)
		sendBlockToEN(t, blockG, blockD, EN1)
		waitTimeout := 5 * time.Second
		waitUntilBlockSealed(t, blockB, EN1, waitTimeout)

		// send all the blocks except G to EN2, which will not trigger
		// state syncing yet. That's because there is only 1 sealed and unexecuted block,
		// which is A, and the threshold for triggering state syncing is 2.
		sendBlockToEN(t, blockA, nil, EN2)
		sendBlockToEN(t, blockB, nil, EN2)
		sendBlockToEN(t, blockC, nil, EN2)
		sendBlockToEN(t, blockD, blockA, EN2)
		sendBlockToEN(t, blockE, blockB, EN2)
		sendBlockToEN(t, blockF, blockC, EN2) // once this function returns, it means blockF is executable.

		// verify that the state syncing is not triggered
		//
		// send G to EN2 to trigger state syncing
		sendBlockToEN(t, blockG, blockD, EN2)
		//
		// verify the state delta is called.
		// wait for a short period of time for the statecommitment for G is ready from syncing state.
		// it will timeout if state syncing didn't happen, because the block execution is too long to
		// wait for.
		waitUntilBlockIsExecuted(t, blockG, EN2, waitTimeout)
	})
}

func withNodes(t *testing.T, f func(en1, en2 *testmock.ExecutionNode)) {
	hub := stub.NewNetworkHub()

	chainID := flow.Mainnet

	colID := unittest.IdentityFixture(unittest.WithRole(flow.RoleCollection))
	conID := unittest.IdentityFixture(unittest.WithRole(flow.RoleConsensus))
	exe1ID := unittest.IdentityFixture(unittest.WithRole(flow.RoleExecution))
	exe2ID := unittest.IdentityFixture(unittest.WithRole(flow.RoleExecution))
	verID := unittest.IdentityFixture(unittest.WithRole(flow.RoleVerification))
	identities := unittest.CompleteIdentitySet(colID, conID, exe1ID, exe2ID, verID)

	syncThreshold := 2
	exeNode1 := testutil.ExecutionNode(t, hub, exe1ID, identities, syncThreshold, chainID)
	exeNode1.Ready()
	defer exeNode1.Done()
	exeNode2 := testutil.ExecutionNode(t, hub, exe2ID, identities, syncThreshold, chainID)
	exeNode2.Ready()
	defer exeNode2.Done()

	collectionNode := testutil.GenericNode(t, hub, colID, identities, chainID)
	defer collectionNode.Done()
	verificationNode := testutil.GenericNode(t, hub, verID, identities, chainID)
	defer verificationNode.Done()
	consensusNode := testutil.GenericNode(t, hub, conID, identities, chainID)
	defer consensusNode.Done()

	f(&exeNode1, &exeNode2)
}

func makeBlockWithParentAndSeal(
	parent *flow.Header, sealed *flow.Header) *flow.Block {
	block := unittest.BlockWithParentFixture(parent)
	payload := flow.Payload{
		Guarantees: nil,
	}

	if sealed != nil {
		payload.Seals = []*flow.Seal{
			unittest.SealFixture(
				unittest.SealWithBlockID(sealed.ID()),
			),
		}
	}

	block.SetPayload(payload)
	return &block
}

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func logBlocks(log zerolog.Logger, blocks ...*flow.Block) {
	for i, b := range blocks {
		name := string(alphabet[i])
		log.Debug().Msgf("creating blocks for testing, block %v's ID:%v", name, b.ID())
	}
}

func sendBlockToEN(t *testing.T, block *flow.Block, finalizes *flow.Block, en *testmock.ExecutionNode) {
	// simulating block finalization
	err := en.State.Mutate().HeaderExtend(block)
	require.NoError(t, err)

	if finalizes != nil {
		err = en.State.Mutate().Finalize(finalizes.ID())
		require.NoError(t, err)
	}

	en.IngestionEngine.BlockProcessable(block.Header)
}

func waitUntilBlockSealed(
	t *testing.T, block *flow.Block, en *testmock.ExecutionNode, timeout time.Duration) {
	require.Eventually(t, func() bool {
		sealed, err := en.GenericNode.State.Sealed().Head()
		fmt.Println("===============> sealed", sealed.Height)
		require.NoError(t, err)
		return sealed.Height >= block.Header.Height
	}, timeout, time.Millisecond*500,
		fmt.Sprintf("expect block %v to be sealed, but timeout", block.ID()))
}

func waitUntilBlockIsExecuted(
	t *testing.T, block *flow.Block, en *testmock.ExecutionNode, timeout time.Duration) {
	blockID := block.ID()
	require.Eventually(t, func() bool {
		_, err := en.ExecutionState.StateCommitmentByBlockID(context.Background(), blockID)
		return err == nil
	}, timeout, time.Millisecond*500,
		fmt.Sprintf("expect block %v to be executed, but timeout", block.ID()))
}
