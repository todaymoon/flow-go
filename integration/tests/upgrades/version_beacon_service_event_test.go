package upgrades

import (
	"context"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
	"time"

	"github.com/coreos/go-semver/semver"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-core-contracts/lib/go/templates"

	sdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go/model/flow"

	"github.com/stretchr/testify/suite"
)

type TestServiceEventVersionControl struct {
	Suite
}

func (s *TestServiceEventVersionControl) TestEmittingVersionBeaconServiceEvent() {

	ctx := context.Background()

	serviceAddress := s.net.Root().Header.ChainID.Chain().ServiceAddress()
	env := templates.Environment{
		NodeVersionBeaconAddress: serviceAddress.String(),
	}

	freezePeriod := s.getFreezePeriod(ctx, env)

	s.Run("should fail adding version boundary inside the freeze period", func() {
		latestFinalized, err := s.AccessClient().GetLatestFinalizedBlockHeader(ctx)
		require.NoError(s.T(), err)

		height := latestFinalized.Height + freezePeriod - 5
		major := uint8(0)
		minor := uint8(0)
		patch := uint8(1)
		preRelease := ""

		txResult := s.sendSetVersionBoundaryTransaction(
			ctx,
			env,
			versionBoundary{
				Major:       major,
				Minor:       minor,
				Patch:       patch,
				PreRelease:  preRelease,
				BlockHeight: height,
			})
		s.Require().Error(txResult.Error)

		sealed := s.ReceiptState.WaitForReceiptFromAny(
			s.T(),
			flow.Identifier(txResult.BlockID))
		s.Require().Len(sealed.ExecutionResult.ServiceEvents, 0)
	})

	s.Run("set freeze period script should work", func() {
		// we also want to do this for the next test to conclude faster
		// this should not be too short, otherwise we might execute to many blocks
		// before the version beacon takes effect
		// if the test is flaky increase this value
		// if the test is too slow decrease this value
		newFreezePeriod := uint64(300)

		s.Require().NotEqual(newFreezePeriod, freezePeriod,
			"the test is pointless, "+
				"please change the freeze period in the test")

		setFreezePeriodScript := templates.GenerateChangeVersionFreezePeriodScript(env)
		sdkServiceAddress := sdk.Address(serviceAddress)

		latestBlockID, err := s.AccessClient().GetLatestBlockID(ctx)
		require.NoError(s.T(), err)

		tx := sdk.NewTransaction().
			SetScript(setFreezePeriodScript).
			SetReferenceBlockID(sdk.Identifier(latestBlockID)).
			SetProposalKey(sdkServiceAddress,
				0, s.AccessClient().GetSeqNumber()). // todo track sequence number
			AddAuthorizer(sdkServiceAddress).
			SetPayer(sdkServiceAddress)

		err = tx.AddArgument(cadence.NewUInt64(newFreezePeriod))
		s.Require().NoError(err)

		err = s.AccessClient().SignAndSendTransaction(ctx, tx)
		s.Require().NoError(err)

		result, err := s.AccessClient().WaitForSealed(ctx, tx.ID())
		require.NoError(s.T(), err)

		s.Require().NoError(result.Error)

		// This changes the freeze  period to something shorter,
		// so we can test the next step
		freezePeriod = s.getFreezePeriod(ctx, env)
		s.Require().Equal(newFreezePeriod, freezePeriod)
	})

	s.Run("should add version boundary after the freeze period", func() {
		latestFinalized, err := s.AccessClient().GetLatestFinalizedBlockHeader(ctx)
		require.NoError(s.T(), err)

		// make sure target height is correct
		// the height at which the version change will take effect should be after
		// the current height + the freeze period
		height := latestFinalized.Height + freezePeriod + 100

		// version 0.0.1
		// low version to not interfere with other tests
		major := uint8(0)
		minor := uint8(0)
		patch := uint8(1)
		preRelease := ""

		txResult := s.sendSetVersionBoundaryTransaction(
			ctx,
			env,
			versionBoundary{
				Major:       major,
				Minor:       minor,
				Patch:       patch,
				PreRelease:  preRelease,
				BlockHeight: height,
			})
		s.Require().NoError(txResult.Error)

		sealed := s.ReceiptState.WaitForReceiptFromAny(
			s.T(),
			flow.Identifier(txResult.BlockID))

		s.Require().Len(sealed.ExecutionResult.ServiceEvents, 1)
		s.Require().IsType(
			&flow.VersionBeacon{},
			sealed.ExecutionResult.ServiceEvents[0].Event)

		versionTable := sealed.ExecutionResult.ServiceEvents[0].Event.(*flow.VersionBeacon)
		// this should be the second ever emitted
		// the first was emitted at bootstrap
		s.Require().Equal(uint64(1), versionTable.Sequence)
		s.Require().Len(versionTable.VersionBoundaries, 2)

		// zeroth boundary should be present, as it is the one we should be on
		s.Require().Equal(uint64(0), versionTable.VersionBoundaries[0].BlockHeight)

		version, err := semver.NewVersion(versionTable.VersionBoundaries[0].Version)
		s.Require().NoError(err)
		s.Require().Equal(uint8(0), uint8(version.Major))
		s.Require().Equal(uint8(0), uint8(version.Minor))
		s.Require().Equal(uint8(0), uint8(version.Patch))

		s.Require().Equal(height, versionTable.VersionBoundaries[1].BlockHeight)

		version, err = semver.NewVersion(versionTable.VersionBoundaries[1].Version)
		s.Require().NoError(err)
		s.Require().Equal(major, uint8(version.Major))
		s.Require().Equal(minor, uint8(version.Minor))
		s.Require().Equal(patch, uint8(version.Patch))
	})

	s.Run("stop with version beacon", func() {
		latestFinalized, err := s.AccessClient().GetLatestFinalizedBlockHeader(ctx)
		require.NoError(s.T(), err)

		// make sure target height is correct
		// the height at which the version change will take effect should be after
		// the current height + the freeze period
		height := latestFinalized.Height + freezePeriod + 100

		// version 0.0.1
		// max version to be sure that the node version is lower
		major := uint8(math.MaxUint8)
		minor := uint8(math.MaxUint8)
		patch := uint8(math.MaxUint8)
		preRelease := ""

		txResult := s.sendSetVersionBoundaryTransaction(
			ctx,
			env,
			versionBoundary{
				Major:       major,
				Minor:       minor,
				Patch:       patch,
				PreRelease:  preRelease,
				BlockHeight: height,
			})
		s.Require().NoError(txResult.Error)

		sealed := s.ReceiptState.WaitForReceiptFromAny(
			s.T(),
			flow.Identifier(txResult.BlockID))

		s.Require().Len(sealed.ExecutionResult.ServiceEvents, 1)
		s.Require().IsType(
			&flow.VersionBeacon{},
			sealed.ExecutionResult.ServiceEvents[0].Event)

		versionTable := sealed.ExecutionResult.ServiceEvents[0].Event.(*flow.VersionBeacon)

		s.Require().Equal(height, versionTable.VersionBoundaries[len(versionTable.VersionBoundaries)-1].BlockHeight)
		version, err := semver.NewVersion(versionTable.VersionBoundaries[len(versionTable.VersionBoundaries)-1].Version)
		s.Require().NoError(err)
		s.Require().Equal(major, uint8(version.Major))
		s.Require().Equal(minor, uint8(version.Minor))
		s.Require().Equal(patch, uint8(version.Patch))

		shouldExecute := s.BlockState.WaitForBlocksByHeight(s.T(), height-1)
		shouldNotExecute := s.BlockState.WaitForBlocksByHeight(s.T(), height)

		s.ReceiptState.WaitForReceiptFrom(s.T(), shouldExecute[0].Header.ID(), s.exe1ID)
		s.ReceiptState.WaitForNoReceiptFrom(
			s.T(),
			5*time.Second,
			shouldNotExecute[0].Header.ID(),
			s.exe1ID,
		)

		enContainer := s.net.ContainerByID(s.exe1ID)
		err = enContainer.WaitForContainerStopped(30 * time.Second)
		s.NoError(err)
	})
}

func (s *TestServiceEventVersionControl) getFreezePeriod(
	ctx context.Context,
	env templates.Environment,
) uint64 {

	freezePeriodScript := templates.GenerateGetVersionBoundaryFreezePeriodScript(env)

	freezePeriodRaw, err := s.AccessClient().
		ExecuteScriptBytes(ctx, freezePeriodScript, nil)
	s.Require().NoError(err)

	cadenceBuffer, is := freezePeriodRaw.(cadence.UInt64)

	s.Require().True(is, "version freezePeriod script returned unknown type")

	return cadenceBuffer.ToGoValue().(uint64)
}

type versionBoundary struct {
	BlockHeight uint64
	Major       uint8
	Minor       uint8
	Patch       uint8
	PreRelease  string
}

func (s *TestServiceEventVersionControl) sendSetVersionBoundaryTransaction(
	ctx context.Context,
	env templates.Environment,
	boundary versionBoundary,
) *sdk.TransactionResult {
	serviceAddress := s.net.Root().Header.ChainID.Chain().ServiceAddress()

	versionTableChangeScript := templates.GenerateSetVersionBoundaryScript(env)

	latestBlockId, err := s.AccessClient().GetLatestBlockID(ctx)
	s.Require().NoError(err)
	seq := s.AccessClient().GetSeqNumber()

	tx := sdk.NewTransaction().
		SetScript(versionTableChangeScript).
		SetReferenceBlockID(sdk.Identifier(latestBlockId)).
		SetProposalKey(sdk.Address(serviceAddress), 0, seq).
		SetPayer(sdk.Address(serviceAddress)).
		AddAuthorizer(sdk.Address(serviceAddress))

	// args
	// newMajor: UInt8,
	// newMinor: UInt8,
	// newPatch: UInt8,
	// newPreRelease: String?,
	// targetBlockHeight: UInt64

	err = tx.AddArgument(cadence.NewUInt8(boundary.Major))
	s.Require().NoError(err)

	err = tx.AddArgument(cadence.NewUInt8(boundary.Minor))
	s.Require().NoError(err)

	err = tx.AddArgument(cadence.NewUInt8(boundary.Patch))
	s.Require().NoError(err)

	preReleaseCadenceString, err := cadence.NewString(boundary.PreRelease)
	s.Require().NoError(err)
	err = tx.AddArgument(preReleaseCadenceString)
	s.Require().NoError(err)

	err = tx.AddArgument(cadence.NewUInt64(boundary.BlockHeight))
	s.Require().NoError(err)

	err = s.AccessClient().SignAndSendTransaction(ctx, tx)
	s.Require().NoError(err)

	txResult, err := s.AccessClient().WaitForSealed(ctx, tx.ID())
	s.Require().NoError(err)
	return txResult
}

func TestVersionControlServiceEvent(t *testing.T) {
	suite.Run(t, new(TestServiceEventVersionControl))
}
