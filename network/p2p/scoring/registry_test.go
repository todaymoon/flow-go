package scoring_test

import (
	"math"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/module/metrics"
	"github.com/onflow/flow-go/module/mock"
	"github.com/onflow/flow-go/network/p2p"
	netcache "github.com/onflow/flow-go/network/p2p/cache"
	mockp2p "github.com/onflow/flow-go/network/p2p/mock"
	"github.com/onflow/flow-go/network/p2p/scoring"
	"github.com/onflow/flow-go/utils/unittest"
)

// TestDefaultDecayFunction tests the default decay function used by the peer scorer.
// The default decay function is used when no custom decay function is provided.
// The test evaluates the following cases:
// 1. score is non-negative and should not be decayed.
// 2. score is negative and above the skipDecayThreshold and lastUpdated is too recent. In this case, the score should not be decayed.
// 3. score is negative and above the skipDecayThreshold and lastUpdated is too old. In this case, the score should not be decayed.
// 4. score is negative and below the skipDecayThreshold and lastUpdated is too recent. In this case, the score should not be decayed.
// 5. score is negative and below the skipDecayThreshold and lastUpdated is too old. In this case, the score should be decayed.
func TestDefaultDecayFunction(t *testing.T) {
	type args struct {
		record      p2p.GossipSubSpamRecord
		lastUpdated time.Time
	}

	type want struct {
		record p2p.GossipSubSpamRecord
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			// 1. score is non-negative and should not be decayed.
			name: "score is non-negative",
			args: args{
				record: p2p.GossipSubSpamRecord{
					Penalty: 5,
					Decay:   0.8,
				},
				lastUpdated: time.Now(),
			},
			want: want{
				record: p2p.GossipSubSpamRecord{
					Penalty: 5,
					Decay:   0.8,
				},
			},
		},
		{ // 2. score is negative and above the skipDecayThreshold and lastUpdated is too recent. In this case, the score should not be decayed,
			// since less than a second has passed since last update.
			name: "score is negative and but above skipDecayThreshold and lastUpdated is too recent",
			args: args{
				record: p2p.GossipSubSpamRecord{
					Penalty: -0.09, // -0.09 is above skipDecayThreshold of -0.1
					Decay:   0.8,
				},
				lastUpdated: time.Now(),
			},
			want: want{
				record: p2p.GossipSubSpamRecord{
					Penalty: 0, // score is set to 0
					Decay:   0.8,
				},
			},
		},
		{
			// 3. score is negative and above the skipDecayThreshold and lastUpdated is too old. In this case, the score should not be decayed,
			// since score is between [skipDecayThreshold, 0] and more than a second has passed since last update.
			name: "score is negative and but above skipDecayThreshold and lastUpdated is too old",
			args: args{
				record: p2p.GossipSubSpamRecord{
					Penalty: -0.09, // -0.09 is above skipDecayThreshold of -0.1
					Decay:   0.8,
				},
				lastUpdated: time.Now().Add(-10 * time.Second),
			},
			want: want{
				record: p2p.GossipSubSpamRecord{
					Penalty: 0, // score is set to 0
					Decay:   0.8,
				},
			},
		},
		{
			// 4. score is negative and below the skipDecayThreshold and lastUpdated is too recent. In this case, the score should not be decayed,
			// since less than a second has passed since last update.
			name: "score is negative and below skipDecayThreshold but lastUpdated is too recent",
			args: args{
				record: p2p.GossipSubSpamRecord{
					Penalty: -5,
					Decay:   0.8,
				},
				lastUpdated: time.Now(),
			},
			want: want{
				record: p2p.GossipSubSpamRecord{
					Penalty: -5,
					Decay:   0.8,
				},
			},
		},
		{
			// 5. score is negative and below the skipDecayThreshold and lastUpdated is too old. In this case, the score should be decayed.
			name: "score is negative and below skipDecayThreshold but lastUpdated is too old",
			args: args{
				record: p2p.GossipSubSpamRecord{
					Penalty: -15,
					Decay:   0.8,
				},
				lastUpdated: time.Now().Add(-10 * time.Second),
			},
			want: want{
				record: p2p.GossipSubSpamRecord{
					Penalty: -15 * math.Pow(0.8, 10),
					Decay:   0.8,
				},
			},
		},
	}

	decayFunc := scoring.DefaultDecayFunction()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decayFunc(tt.args.record, tt.args.lastUpdated)
			assert.NoError(t, err)
			assert.Less(t, math.Abs(got.Penalty-tt.want.record.Penalty), 10e-3)
			assert.Equal(t, got.Decay, tt.want.record.Decay)
		})
	}
}

// TestInit tests when a peer id is queried for the first time by the
// app specific score function, the score is initialized to the initial state.
func TestInitSpamRecords(t *testing.T) {
	reg, cache := newGossipSubAppSpecificScoreRegistry(t)
	peerID := peer.ID("peer-1")

	// initially, the cache should not have the peer id.
	assert.False(t, cache.Has(peerID))

	// when the app specific score function is called for the first time, the score should be initialized to the initial state.
	score := reg.AppSpecificScoreFunc()(peerID)
	assert.Equal(t, score, scoring.InitAppScoreRecordState().Penalty) // score should be initialized to the initial state.

	// the cache should now have the peer id.
	assert.True(t, cache.Has(peerID))
	record, err, ok := cache.Get(peerID) // get the record from the cache.
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.Equal(t, record.Penalty, scoring.InitAppScoreRecordState().Penalty) // score should be initialized to the initial state.
	assert.Equal(t, record.Decay, scoring.InitAppScoreRecordState().Decay)     // decay should be initialized to the initial state.
}

func TestInitWhenGetGoesFirst(t *testing.T) {
	t.Run("graft", func(t *testing.T) {
		testInitWhenGetFirst(t, p2p.CtrlMsgGraft, penaltyValueFixtures().Graft)
	})
	t.Run("prune", func(t *testing.T) {
		testInitWhenGetFirst(t, p2p.CtrlMsgPrune, penaltyValueFixtures().Prune)
	})
	t.Run("ihave", func(t *testing.T) {
		testInitWhenGetFirst(t, p2p.CtrlMsgIHave, penaltyValueFixtures().IHave)
	})
	t.Run("iwant", func(t *testing.T) {
		testInitWhenGetFirst(t, p2p.CtrlMsgIWant, penaltyValueFixtures().IWant)
	})
}

// testInitWhenGetFirst tests when a peer id is queried for the first time by the
// app specific score function, the score is initialized to the initial state. Then, the score is reported and the
// score is updated in the cache. The next time the app specific score function is called, the score should be the
// updated score.
func testInitWhenGetFirst(t *testing.T, messageType p2p.ControlMessageType, expectedPenalty float64) {
	reg, cache := newGossipSubAppSpecificScoreRegistry(t)
	peerID := peer.ID("peer-1")

	// initially, the cache should not have the peer id.
	assert.False(t, cache.Has(peerID))

	// when the app specific score function is called for the first time, the score should be initialized to the initial state.
	score := reg.AppSpecificScoreFunc()(peerID)
	assert.Equal(t, score, scoring.InitAppScoreRecordState().Penalty) // score should be initialized to the initial state.

	// the cache should now have the peer id.
	assert.True(t, cache.Has(peerID))
	record, err, ok := cache.Get(peerID) // get the record from the cache.
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.Equal(t, record.Penalty, scoring.InitAppScoreRecordState().Penalty) // score should be initialized to the initial state.
	assert.Equal(t, record.Decay, scoring.InitAppScoreRecordState().Decay)     // decay should be initialized to the initial state.

	// report a misbehavior for the peer id.
	reg.OnInvalidControlMessageNotification(&p2p.InvalidControlMessageNotification{
		PeerID:  peerID,
		MsgType: messageType,
		Count:   1,
	})

	// the score should now be updated.
	record, err, ok = cache.Get(peerID) // get the record from the cache.
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.Less(t, math.Abs(expectedPenalty-record.Penalty), 10e-3)        // score should be updated to -10.
	assert.Equal(t, scoring.InitAppScoreRecordState().Decay, record.Decay) // decay should be initialized to the initial state.

	// when the app specific score function is called again, the score should be updated.
	score = reg.AppSpecificScoreFunc()(peerID)
	assert.Less(t, math.Abs(expectedPenalty-score), 10e-3) // score should be updated to -10.
}

func TestInitWhenReportGoesFirst(t *testing.T) {
	t.Run("graft", func(t *testing.T) {
		testInitWhenReportGoesFirst(t, p2p.CtrlMsgGraft, penaltyValueFixtures().Graft)
	})
	t.Run("prune", func(t *testing.T) {
		testInitWhenReportGoesFirst(t, p2p.CtrlMsgPrune, penaltyValueFixtures().Prune)
	})
	t.Run("ihave", func(t *testing.T) {
		testInitWhenReportGoesFirst(t, p2p.CtrlMsgIHave, penaltyValueFixtures().IHave)
	})
	t.Run("iwant", func(t *testing.T) {
		testInitWhenReportGoesFirst(t, p2p.CtrlMsgIWant, penaltyValueFixtures().IWant)
	})
}

// testInitWhenReportGoesFirst tests situation where a peer id is reported for the first time
// before the app specific score function is called for the first time on it.
// The test expects the score to be initialized to the initial state and then updated by the penalty value.
// Subsequent calls to the app specific score function should return the updated score.
func testInitWhenReportGoesFirst(t *testing.T, messageType p2p.ControlMessageType, expectedPenalty float64) {
	reg, cache := newGossipSubAppSpecificScoreRegistry(t)
	peerID := peer.ID("peer-1")

	// report a misbehavior for the peer id.
	reg.OnInvalidControlMessageNotification(&p2p.InvalidControlMessageNotification{
		PeerID:  peerID,
		MsgType: p2p.CtrlMsgGraft,
		Count:   1,
	})

	// the score should now be updated.
	record, err, ok := cache.Get(peerID) // get the record from the cache.
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.Less(t, math.Abs(scoring.DefaultGossipSubCtrlMsgPenaltyValue().Graft-record.Penalty), 10e-3) // score should be updated to -10, we account for decay.
	assert.Equal(t, scoring.InitAppScoreRecordState().Decay, record.Decay)                              // decay should be initialized to the initial state.

	// when the app specific score function is called for the first time, the score should be updated.
	score := reg.AppSpecificScoreFunc()(peerID)
	assert.Less(t, math.Abs(scoring.DefaultGossipSubCtrlMsgPenaltyValue().Graft-score), 10e-3) // score should be updated to -10, we account for decay.
}

// TestSpamPenaltyDecaysInCache tests that the spam penalty records decay over time in the cache.
func TestSpamPenaltyDecaysInCache(t *testing.T) {
	peerID := peer.ID("peer-1")
	reg, _ := newGossipSubAppSpecificScoreRegistry(t,
		withStakedIdentity(peerID),
		withValidSubscriptions(peerID))

	// report a misbehavior for the peer id.
	reg.OnInvalidControlMessageNotification(&p2p.InvalidControlMessageNotification{
		PeerID:  peerID,
		MsgType: p2p.CtrlMsgPrune,
		Count:   1,
	})

	time.Sleep(1 * time.Second) // wait for the penalty to decay.

	reg.OnInvalidControlMessageNotification(&p2p.InvalidControlMessageNotification{
		PeerID:  peerID,
		MsgType: p2p.CtrlMsgGraft,
		Count:   1,
	})

	time.Sleep(1 * time.Second) // wait for the penalty to decay.

	reg.OnInvalidControlMessageNotification(&p2p.InvalidControlMessageNotification{
		PeerID:  peerID,
		MsgType: p2p.CtrlMsgIHave,
		Count:   1,
	})

	time.Sleep(1 * time.Second) // wait for the penalty to decay.

	reg.OnInvalidControlMessageNotification(&p2p.InvalidControlMessageNotification{
		PeerID:  peerID,
		MsgType: p2p.CtrlMsgIWant,
		Count:   1,
	})

	time.Sleep(1 * time.Second) // wait for the penalty to decay.

	// when the app specific score function is called for the first time, the decay functionality should be kicked in
	// the cache, and the score should be updated. Note that since the penalty values are negative, the default staked identity
	// reward is not applied. Hence, the score is only comprised of the penalties.
	score := reg.AppSpecificScoreFunc()(peerID)
	// the upper bound is the sum of the penalties without decay.
	scoreUpperBound := penaltyValueFixtures().Prune +
		penaltyValueFixtures().Graft +
		penaltyValueFixtures().IHave +
		penaltyValueFixtures().IWant
	// the lower bound is the sum of the penalties with decay assuming the decay is applied 4 times to the sum of the penalties.
	// in reality, the decay is applied 4 times to the first penalty, then 3 times to the second penalty, and so on.
	scoreLowerBound := scoreUpperBound * math.Pow(scoring.InitAppScoreRecordState().Decay, 4)

	// with decay, the score should be between the upper and lower bounds.
	assert.Greater(t, score, scoreUpperBound)
	assert.Less(t, score, scoreLowerBound)
}

// TestSpamPenaltyDecayToZero tests that the spam penalty decays to zero over time, and when the spam penalty of
// a peer is set back to zero, its app specific score is also reset to the initial state.
func TestSpamPenaltyDecayToZero(t *testing.T) {
	cache := netcache.NewGossipSubSpamRecordCache(100, unittest.Logger(), metrics.NewNoopCollector(), scoring.DefaultDecayFunction())

	// mocks peer has an staked identity and is subscribed to the allowed topics.
	idProvider := mock.NewIdentityProvider(t)
	peerID := peer.ID("peer-1")
	idProvider.On("ByPeerID", peerID).Return(unittest.IdentityFixture(), true).Maybe()

	validator := mockp2p.NewSubscriptionValidator(t)
	validator.On("CheckSubscribedToAllowedTopics", peerID, testifymock.Anything).Return(nil).Maybe()

	reg := scoring.NewGossipSubAppSpecificScoreRegistry(&scoring.GossipSubAppSpecificScoreRegistryConfig{
		Logger:        unittest.Logger(),
		DecayFunction: scoring.DefaultDecayFunction(),
		Penalty:       penaltyValueFixtures(),
		Validator:     validator,
		IdProvider:    idProvider,
		CacheFactory: func() p2p.GossipSubSpamRecordCache {
			return cache
		},
		Init: func() p2p.GossipSubSpamRecord {
			return p2p.GossipSubSpamRecord{
				Decay:   0.02, // we choose a small decay value to speed up the test.
				Penalty: 0,
			}
		},
	})

	// report a misbehavior for the peer id.
	reg.OnInvalidControlMessageNotification(&p2p.InvalidControlMessageNotification{
		PeerID:  peerID,
		MsgType: p2p.CtrlMsgGraft,
		Count:   1,
	})

	// decays happen every second, so we wait for 1 second to make sure the score is updated.
	time.Sleep(1 * time.Second)
	// the score should now be updated, it should be still negative but greater than the penalty value (due to decay).
	score := reg.AppSpecificScoreFunc()(peerID)
	require.Less(t, score, float64(0))                      // the score should be less than zero.
	require.Greater(t, score, penaltyValueFixtures().Graft) // the score should be less than the penalty value due to decay.

	require.Eventually(t, func() bool {
		// the spam penalty should eventually decay to zero.
		r, err, ok := cache.Get(peerID)
		return ok && err == nil && r.Penalty == 0.0
	}, 5*time.Second, 100*time.Millisecond)

	require.Eventually(t, func() bool {
		// when the spam penalty is decayed to zero, the app specific score of the node should reset back to its initial state (i.e., max reward).
		return reg.AppSpecificScoreFunc()(peerID) == scoring.MaxAppSpecificReward
	}, 5*time.Second, 100*time.Millisecond)

	// the score should now be zero.
	record, err, ok := cache.Get(peerID) // get the record from the cache.
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, record.Penalty) // score should be zero.
}

// withStakedIdentity returns a function that sets the identity provider to return an staked identity for the given peer id.
// It is used for testing purposes, and causes the given peer id to benefit from the staked identity reward in GossipSub.
func withStakedIdentity(peerId peer.ID) func(cfg *scoring.GossipSubAppSpecificScoreRegistryConfig) {
	return func(cfg *scoring.GossipSubAppSpecificScoreRegistryConfig) {
		cfg.IdProvider.(*mock.IdentityProvider).On("ByPeerID", peerId).Return(unittest.IdentityFixture(), true).Maybe()
	}
}

// withValidSubscriptions returns a function that sets the subscription validator to return nil for the given peer id.
// It is used for testing purposes and causes the given peer id to never be penalized for subscribing to invalid topics.
func withValidSubscriptions(peer peer.ID) func(cfg *scoring.GossipSubAppSpecificScoreRegistryConfig) {
	return func(cfg *scoring.GossipSubAppSpecificScoreRegistryConfig) {
		cfg.Validator.(*mockp2p.SubscriptionValidator).On("CheckSubscribedToAllowedTopics", peer, testifymock.Anything).Return(nil).Maybe()
	}
}

// newGossipSubAppSpecificScoreRegistry returns a new instance of GossipSubAppSpecificScoreRegistry with default values
// for the testing purposes.
func newGossipSubAppSpecificScoreRegistry(t *testing.T, opts ...func(*scoring.GossipSubAppSpecificScoreRegistryConfig)) (*scoring.GossipSubAppSpecificScoreRegistry, *netcache.GossipSubSpamRecordCache) {
	cache := netcache.NewGossipSubSpamRecordCache(100, unittest.Logger(), metrics.NewNoopCollector(), scoring.DefaultDecayFunction())
	cfg := &scoring.GossipSubAppSpecificScoreRegistryConfig{
		Logger:        unittest.Logger(),
		DecayFunction: scoring.DefaultDecayFunction(),
		Init:          scoring.InitAppScoreRecordState,
		Penalty:       penaltyValueFixtures(),
		IdProvider:    mock.NewIdentityProvider(t),
		Validator:     mockp2p.NewSubscriptionValidator(t),
		CacheFactory: func() p2p.GossipSubSpamRecordCache {
			return cache
		},
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return scoring.NewGossipSubAppSpecificScoreRegistry(cfg), cache
}

// penaltyValueFixtures returns a set of penalty values for testing purposes.
// The values are not realistic. The important thing is that they are different from each other. This is to make sure
// that the tests are not passing because of the default values.
func penaltyValueFixtures() scoring.GossipSubCtrlMsgPenaltyValue {
	return scoring.GossipSubCtrlMsgPenaltyValue{
		Graft: -100,
		Prune: -50,
		IHave: -20,
		IWant: -10,
	}
}
