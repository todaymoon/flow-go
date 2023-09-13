package internal_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	p2ptest "github.com/onflow/flow-go/network/p2p/test"
	"github.com/onflow/flow-go/network/p2p/unicast/internal"
	"github.com/onflow/flow-go/network/p2p/unicast/model"
)

// TestDialConfigEntity tests the DialConfigEntity struct and its methods.
func TestDialConfigEntity(t *testing.T) {
	peerID := p2ptest.PeerIdFixture(t)

	d := &internal.DialConfigEntity{
		PeerId: peerID,
		DialConfig: model.DialConfig{
			DialBackoff:        10,
			StreamBackoff:      20,
			LastSuccessfulDial: 30,
		},
	}

	t.Run("Test ID and Checksum", func(t *testing.T) {
		// id and checksum methods must return the same value as expected.
		expectedID := internal.PeerIdToFlowId(peerID)
		require.Equal(t, expectedID, d.ID())
		require.Equal(t, expectedID, d.Checksum())
	})
}
