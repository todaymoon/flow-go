package flow

import (
	"fmt"

	clone "github.com/huandu/go-clone/generic"
	"golang.org/x/exp/slices"
)

// DynamicIdentityEntry encapsulates nodeID and dynamic portion of identity.
type DynamicIdentityEntry struct {
	NodeID  Identifier
	Ejected bool
}

type DynamicIdentityEntryList []*DynamicIdentityEntry

// EpochMinStateEntry represents a snapshot of the identity table (incl. the set of all notes authorized to
// be part of the network) at some point in time. It allows to reconstruct the state of identity table using
// epoch setup events and dynamic identities. It tracks attempts of invalid state transitions.
// It also holds information about the next epoch, if it has been already committed.
// This structure is used to persist protocol state in the database.
//
// Note that the current implementation does not store the identity table directly. Instead, we store
// the original events that constituted the _initial_ identity table at the beginning of the epoch
// plus some modifiers. We intend to restructure this code soon.
type EpochMinStateEntry struct {
	PreviousEpoch *EpochStateContainer // minimal dynamic properties for previous epoch [optional, nil for first epoch after spork, genesis]
	CurrentEpoch  EpochStateContainer  // minimal dynamic properties for current epoch
	NextEpoch     *EpochStateContainer // minimal dynamic properties for next epoch [optional, nil iff we are in staking phase]

	// EpochFallbackTriggered encodes whether an invalid epoch transition
	// has been detected in this fork. Under normal operations, this value is false.
	// Node-internally, the EpochFallback notification is emitted when a block is
	// finalized that changes this flag from false to true.
	// A state transition from true -> false is possible only when protocol undergoes epoch recovery.
	EpochFallbackTriggered bool
}

// EpochStateContainer holds the data pertaining to a _single_ epoch but no information about
// any adjacent epochs. To perform a transition from epoch N to N+1, EpochStateContainers for
// both epochs are necessary.
type EpochStateContainer struct {
	// ID of setup event for this epoch, never nil.
	SetupID Identifier
	// ID of commit event for this epoch. Could be ZeroID if epoch was not committed.
	CommitID Identifier
	// ActiveIdentities contains the dynamic identity properties for the nodes that
	// are active in this epoch. Active means that these nodes are authorized to contribute to
	// extending the chain. Nodes are listed in `ActiveIdentities` if and only if
	// they are part of the EpochSetup event for the respective epoch.
	// The dynamic identity properties can change from block to block. Each non-deferred
	// identity-mutating operation is applied independently to the `ActiveIdentities`
	// of the relevant epoch's EpochStateContainer separately.
	// Identities are always sorted in canonical order.
	//
	// Context: In comparison, nodes that are joining in the next epoch or left as of this
	// epoch are only allowed to listen to the network but not actively contribute. Such
	// nodes are _not_ part of `Identities`.
	ActiveIdentities DynamicIdentityEntryList

	// EpochExtensions contains potential EFM-extensions of this epoch. In the happy path
	// it is nil or empty. An Epoch in which Epoch-Fallback-Mode [EFM] is triggered, will
	// have at least one extension. By convention, the initial extension must satisfy
	//   EpochSetup.FinalView + 1 = EpochExtensions[0].FirstView
	// and each consecutive pair of slice elements must obey
	//   EpochExtensions[i].FinalView+1 = EpochExtensions[i+1].FirstView
	EpochExtensions []EpochExtension
}

// EpochExtension represents a range of views, which contiguously extends this epoch.
type EpochExtension struct {
	FirstView     uint64
	FinalView     uint64
	TargetEndTime uint64
}

// ID returns an identifier for this EpochStateContainer by hashing internal fields.
// Per convention, the ID of a `nil` EpochStateContainer is `flow.ZeroID`.
func (c *EpochStateContainer) ID() Identifier {
	if c == nil {
		return ZeroID
	}
	return MakeID(c)
}

// EventIDs returns the `flow.EventIDs` with the hashes of the EpochSetup and EpochCommit events.
// Per convention, for a `nil` EpochStateContainer, we return `flow.ZeroID` for both events.
func (c *EpochStateContainer) EventIDs() EventIDs {
	if c == nil {
		return EventIDs{ZeroID, ZeroID}
	}
	return EventIDs{c.SetupID, c.CommitID}
}

// Copy returns a full copy of the entry.
// Embedded Identities are deep-copied, _except_ for their keys, which are copied by reference.
// Per convention, the ID of a `nil` EpochStateContainer is `flow.ZeroID`.
func (c *EpochStateContainer) Copy() *EpochStateContainer {
	if c == nil {
		return nil
	}
	return &EpochStateContainer{
		SetupID:          c.SetupID,
		CommitID:         c.CommitID,
		ActiveIdentities: c.ActiveIdentities.Copy(),
		EpochExtensions:  clone.Clone(c.EpochExtensions),
	}
}

type EpochStateEntry struct {
	*EpochMinStateEntry

	PreviousEpochSetup  *EpochSetup
	PreviousEpochCommit *EpochCommit
	CurrentEpochSetup   *EpochSetup
	CurrentEpochCommit  *EpochCommit
	NextEpochSetup      *EpochSetup
	NextEpochCommit     *EpochCommit
}

func NewEpochStateEntry(
	protocolState *EpochMinStateEntry,
	previousEpochSetup *EpochSetup,
	previousEpochCommit *EpochCommit,
	currentEpochSetup *EpochSetup,
	currentEpochCommit *EpochCommit,
	nextEpochSetup *EpochSetup,
	nextEpochCommit *EpochCommit,
) (*EpochStateEntry, error) {
	result := &EpochStateEntry{
		EpochMinStateEntry:  protocolState,
		PreviousEpochSetup:  previousEpochSetup,
		PreviousEpochCommit: previousEpochCommit,
		CurrentEpochSetup:   currentEpochSetup,
		CurrentEpochCommit:  currentEpochCommit,
		NextEpochSetup:      nextEpochSetup,
		NextEpochCommit:     nextEpochCommit,
	}

	// If previous epoch is specified: ensure respective epoch service events are not nil and consistent with commitments in `EpochMinStateEntry.PreviousEpoch`
	if protocolState.PreviousEpoch != nil {
		if protocolState.PreviousEpoch.SetupID != previousEpochSetup.ID() { // calling ID() will panic is EpochSetup event is nil
			return nil, fmt.Errorf("supplied previous epoch's setup event (%x) does not match commitment (%x) in EpochMinStateEntry", previousEpochSetup.ID(), protocolState.PreviousEpoch.SetupID)
		}
		if protocolState.PreviousEpoch.CommitID != previousEpochCommit.ID() { // calling ID() will panic is EpochCommit event is nil
			return nil, fmt.Errorf("supplied previous epoch's commit event (%x) does not match commitment (%x) in EpochMinStateEntry", previousEpochCommit.ID(), protocolState.PreviousEpoch.CommitID)
		}
	} else {
		if previousEpochSetup != nil {
			return nil, fmt.Errorf("no previous epoch but gotten non-nil EpochSetup event")
		}
		if previousEpochCommit != nil {
			return nil, fmt.Errorf("no previous epoch but gotten non-nil EpochCommit event")
		}
	}

	// For current epoch: ensure respective epoch service events are not nil and consistent with commitments in `EpochMinStateEntry.CurrentEpoch`
	if protocolState.CurrentEpoch.SetupID != currentEpochSetup.ID() { // calling ID() will panic is EpochSetup event is nil
		return nil, fmt.Errorf("supplied current epoch's setup event (%x) does not match commitment (%x) in EpochMinStateEntry", currentEpochSetup.ID(), protocolState.CurrentEpoch.SetupID)
	}
	if protocolState.CurrentEpoch.CommitID != currentEpochCommit.ID() { // calling ID() will panic is EpochCommit event is nil
		return nil, fmt.Errorf("supplied current epoch's commit event (%x) does not match commitment (%x) in EpochMinStateEntry", currentEpochCommit.ID(), protocolState.CurrentEpoch.CommitID)
	}

	// If we are in staking phase (i.e. protocolState.NextEpoch == nil):
	//  (1) Full identity table contains active identities from current epoch.
	//      If previous epoch exists, we add nodes from previous epoch that are leaving in the current epoch with `EpochParticipationStatusLeaving` status.
	// Otherwise, we are in epoch setup or epoch commit phase (i.e. protocolState.NextEpoch ≠ nil):
	//  (2a) Full identity table contains active identities from current epoch + nodes joining in next epoch with `EpochParticipationStatusJoining` status.
	//  (2b) Furthermore, we also build the full identity table for the next epoch's staking phase:
	//       active identities from next epoch + nodes from current epoch that are leaving at the end of the current epoch with `flow.EpochParticipationStatusLeaving` status.
	nextEpoch := protocolState.NextEpoch
	if nextEpoch == nil { // in staking phase: build full identity table for current epoch according to (1)
		if nextEpochSetup != nil {
			return nil, fmt.Errorf("no next epoch but gotten non-nil EpochSetup event")
		}
		if nextEpochCommit != nil {
			return nil, fmt.Errorf("no next epoch but gotten non-nil EpochCommit event")
		}
	} else { // protocolState.NextEpoch ≠ nil, i.e. we are in epoch setup or epoch commit phase
		// ensure respective epoch service events are not nil and consistent with commitments in `EpochMinStateEntry.NextEpoch`
		if nextEpoch.SetupID != nextEpochSetup.ID() {
			return nil, fmt.Errorf("supplied next epoch's setup event (%x) does not match commitment (%x) in EpochMinStateEntry", nextEpoch.SetupID, nextEpochSetup.ID())
		}
		if nextEpoch.CommitID != ZeroID {
			if nextEpoch.CommitID != nextEpochCommit.ID() {
				return nil, fmt.Errorf("supplied next epoch's commit event (%x) does not match commitment (%x) in EpochMinStateEntry", nextEpoch.CommitID, nextEpochCommit.ID())
			}
		} else {
			if nextEpochCommit != nil {
				return nil, fmt.Errorf("next epoch not yet committed but got EpochCommit event")
			}
		}
	}
	return result, nil
}

// EpochRichStateEntry is a EpochMinStateEntry which has additional fields that are cached
// from storage layer for convenience.
// Using this structure instead of EpochMinStateEntry allows us to avoid querying
// the database for epoch setups and commits and full identity table.
// It holds several invariants, such as:
//   - CurrentEpochSetup and CurrentEpochCommit are for the same epoch. Never nil.
//   - PreviousEpochSetup and PreviousEpochCommit are for the same epoch. Can be nil.
//   - CurrentEpochIdentityTable is the full (dynamic) identity table for the current epoch.
//     Identities are sorted in canonical order. Without duplicates. Never nil.
//   - NextEpochIdentityTable is the full (dynamic) identity table for the next epoch. Can be nil.
//
// NOTE regarding `CurrentEpochIdentityTable` and `NextEpochIdentityTable`:
// The Identity Table is generally a super-set of the identities listed in the Epoch
// Service Events for the respective epoch. This is because the service events only list
// nodes that are authorized to _actively_ contribute to extending the chain. In contrast,
// the Identity Table additionally contains nodes (with weight zero) from the previous or
// upcoming epoch, which are transitioning into / out of the network and are only allowed
// to listen but not to actively contribute.
type EpochRichStateEntry struct {
	*EpochStateEntry

	CurrentEpochIdentityTable IdentityList
	NextEpochIdentityTable    IdentityList
}

// NewEpochRichStateEntry constructs a EpochRichStateEntry from an EpochMinStateEntry and additional data.
// No errors are expected during normal operation. All errors indicate inconsistent or invalid inputs.
func NewEpochRichStateEntry(
	protocolState *EpochStateEntry,
) (*EpochRichStateEntry, error) {
	result := &EpochRichStateEntry{
		EpochStateEntry:           protocolState,
		CurrentEpochIdentityTable: IdentityList{},
		NextEpochIdentityTable:    IdentityList{},
	}
	// If we are in staking phase (i.e. protocolState.NextEpoch == nil):
	//  (1) Full identity table contains active identities from current epoch.
	//      If previous epoch exists, we add nodes from previous epoch that are leaving in the current epoch with `EpochParticipationStatusLeaving` status.
	// Otherwise, we are in epoch setup or epoch commit phase (i.e. protocolState.NextEpoch ≠ nil):
	//  (2a) Full identity table contains active identities from current epoch + nodes joining in next epoch with `EpochParticipationStatusJoining` status.
	//  (2b) Furthermore, we also build the full identity table for the next epoch's staking phase:
	//       active identities from next epoch + nodes from current epoch that are leaving at the end of the current epoch with `flow.EpochParticipationStatusLeaving` status.
	var err error
	nextEpoch := protocolState.NextEpoch
	if nextEpoch == nil { // in staking phase: build full identity table for current epoch according to (1)
		var previousEpochIdentitySkeletons IdentitySkeletonList
		var previousEpochDynamicIdentities DynamicIdentityEntryList
		if previousEpochSetup := protocolState.PreviousEpochSetup; previousEpochSetup != nil {
			previousEpochIdentitySkeletons = previousEpochSetup.Participants
			previousEpochDynamicIdentities = protocolState.PreviousEpoch.ActiveIdentities
		}
		result.CurrentEpochIdentityTable, err = BuildIdentityTable(
			protocolState.CurrentEpochSetup.Participants,
			protocolState.CurrentEpoch.ActiveIdentities,
			previousEpochIdentitySkeletons,
			previousEpochDynamicIdentities,
			EpochParticipationStatusLeaving,
		)
		if err != nil {
			return nil, fmt.Errorf("could not build identity table for staking phase: %w", err)
		}
	} else { // protocolState.NextEpoch ≠ nil, i.e. we are in epoch setup or epoch commit phase
		result.CurrentEpochIdentityTable, err = BuildIdentityTable(
			protocolState.CurrentEpochSetup.Participants,
			protocolState.CurrentEpoch.ActiveIdentities,
			protocolState.NextEpochSetup.Participants,
			nextEpoch.ActiveIdentities,
			EpochParticipationStatusJoining,
		)
		if err != nil {
			return nil, fmt.Errorf("could not build identity table for setup/commit phase: %w", err)
		}

		result.NextEpochIdentityTable, err = BuildIdentityTable(
			protocolState.NextEpochSetup.Participants,
			nextEpoch.ActiveIdentities,
			protocolState.CurrentEpochSetup.Participants,
			protocolState.CurrentEpoch.ActiveIdentities,
			EpochParticipationStatusLeaving,
		)
		if err != nil {
			return nil, fmt.Errorf("could not build next epoch identity table: %w", err)
		}
	}
	return result, nil
}

// ID returns hash of entry by hashing all fields.
func (e *EpochMinStateEntry) ID() Identifier {
	if e == nil {
		return ZeroID
	}
	body := struct {
		PreviousEpochID        Identifier
		CurrentEpochID         Identifier
		NextEpochID            Identifier
		EpochFallbackTriggered bool
	}{
		PreviousEpochID:        e.PreviousEpoch.ID(),
		CurrentEpochID:         e.CurrentEpoch.ID(),
		NextEpochID:            e.NextEpoch.ID(),
		EpochFallbackTriggered: e.EpochFallbackTriggered,
	}
	return MakeID(body)
}

// Copy returns a full copy of the entry.
// Embedded Identities are deep-copied, _except_ for their keys, which are copied by reference.
func (e *EpochMinStateEntry) Copy() *EpochMinStateEntry {
	if e == nil {
		return nil
	}
	return &EpochMinStateEntry{
		PreviousEpoch:          e.PreviousEpoch.Copy(),
		CurrentEpoch:           *e.CurrentEpoch.Copy(),
		NextEpoch:              e.NextEpoch.Copy(),
		EpochFallbackTriggered: e.EpochFallbackTriggered,
	}
}

// Copy returns a full copy of the EpochStateEntry.
//   - Embedded service events are copied by reference (not deep-copied).
func (e *EpochStateEntry) Copy() *EpochStateEntry {
	if e == nil {
		return nil
	}
	return &EpochStateEntry{
		EpochMinStateEntry:  e.EpochMinStateEntry.Copy(),
		PreviousEpochSetup:  e.PreviousEpochSetup,
		PreviousEpochCommit: e.PreviousEpochCommit,
		CurrentEpochSetup:   e.CurrentEpochSetup,
		CurrentEpochCommit:  e.CurrentEpochCommit,
		NextEpochSetup:      e.NextEpochSetup,
		NextEpochCommit:     e.NextEpochCommit,
	}
}

// Copy returns a full copy of the EpochRichStateEntry.
//   - Embedded service events are copied by reference (not deep-copied).
//   - CurrentEpochIdentityTable and NextEpochIdentityTable are deep-copied, _except_ for their keys, which are copied by reference.
func (e *EpochRichStateEntry) Copy() *EpochRichStateEntry {
	if e == nil {
		return nil
	}
	return &EpochRichStateEntry{
		EpochStateEntry:           e.EpochStateEntry.Copy(),
		CurrentEpochIdentityTable: e.CurrentEpochIdentityTable.Copy(),
		NextEpochIdentityTable:    e.NextEpochIdentityTable.Copy(),
	}
}

// CurrentEpochFinalView returns the final view of the current epoch, taking into account possible epoch extensions.
// If there are no epoch extensions, the final view is the final view of the current epoch setup,
// otherwise it is the final view of the last epoch extension.
func (e *EpochRichStateEntry) CurrentEpochFinalView() uint64 {
	l := len(e.CurrentEpoch.EpochExtensions)
	if l > 0 {
		return e.CurrentEpoch.EpochExtensions[l-1].FinalView
	}
	return e.CurrentEpochSetup.FinalView
}

// EpochPhase returns the current epoch phase.
// The receiver EpochMinStateEntry must be properly constructed.
// See flow.EpochPhase for detailed documentation.
func (e *EpochMinStateEntry) EpochPhase() EpochPhase {
	// CAUTION: the logic below that deduces the EpochPhase must be consistent with `epochs.FallbackStateMachine`,
	// which sets the fields we are using here. Specifically, we require that the FallbackStateMachine clears out
	// any tentative values for a subsequent epoch _unless_ that epoch is already committed.
	if e.EpochFallbackTriggered {
		// If the next epoch has been committed, we are in EpochPhaseCommitted regardless of EFM status.
		// We will enter EpochPhaseFallback after completing the transition into the committed next epoch.
		if e.NextEpoch != nil && e.NextEpoch.CommitID != ZeroID {
			return EpochPhaseCommitted
		}
		// If the next epoch has not been committed and EFM is triggered, we immediately enter EpochPhaseFallback.
		return EpochPhaseFallback
	}

	// The epoch phase is determined by how much information we have about the next epoch
	if e.NextEpoch == nil {
		return EpochPhaseStaking // if no information about the next epoch is known, we are in the Staking Phase
	}
	// Per convention, NextEpoch ≠ nil if and only if NextEpoch.SetupID is specified.
	if e.NextEpoch.CommitID == ZeroID {
		return EpochPhaseSetup // if only the Setup event is known for the next epoch but not the Commit event, we are in the Setup Phase
	}
	return EpochPhaseCommitted // if the Setup and Commit events are known for the next epoch, we are in the Committed Phase
}

// EpochCounter returns the current epoch counter.
// The receiver EpochRichStateEntry must be properly constructed.
func (e *EpochRichStateEntry) EpochCounter() uint64 {
	return e.CurrentEpochSetup.Counter
}

func (ll DynamicIdentityEntryList) Lookup() map[Identifier]*DynamicIdentityEntry {
	result := make(map[Identifier]*DynamicIdentityEntry, len(ll))
	for _, entry := range ll {
		result[entry.NodeID] = entry
	}
	return result
}

// Sorted returns whether the list is sorted by the input ordering.
func (ll DynamicIdentityEntryList) Sorted(less IdentifierOrder) bool {
	return slices.IsSortedFunc(ll, func(lhs, rhs *DynamicIdentityEntry) int {
		return less(lhs.NodeID, rhs.NodeID)
	})
}

// ByNodeID gets a node from the list by node ID.
func (ll DynamicIdentityEntryList) ByNodeID(nodeID Identifier) (*DynamicIdentityEntry, bool) {
	for _, identity := range ll {
		if identity.NodeID == nodeID {
			return identity, true
		}
	}
	return nil, false
}

// Copy returns a copy of the DynamicIdentityEntryList. The resulting slice uses
// a different backing array, meaning appends and insert operations on either slice
// are guaranteed to only affect that slice.
//
// Copy should be used when modifying an existing identity list by either
// appending new elements, re-ordering, or inserting new elements in an
// existing index.
//
// CAUTION:
// All Identity fields are deep-copied, _except_ for their keys, which
// are copied by reference.
func (ll DynamicIdentityEntryList) Copy() DynamicIdentityEntryList {
	lenList := len(ll)
	dup := make(DynamicIdentityEntryList, 0, lenList)
	for i := 0; i < lenList; i++ {
		// copy the object
		next := *(ll[i])
		dup = append(dup, &next)
	}
	return dup
}

// Sort sorts the list by the input ordering. Returns a new, sorted list without modifying the input.
// CAUTION:
// All Identity fields are deep-copied, _except_ for their keys, which are copied by reference.
func (ll DynamicIdentityEntryList) Sort(less IdentifierOrder) DynamicIdentityEntryList {
	dup := ll.Copy()
	slices.SortFunc(dup, func(lhs, rhs *DynamicIdentityEntry) int {
		return less(lhs.NodeID, rhs.NodeID)
	})
	return dup
}

// BuildIdentityTable constructs the full identity table for the target epoch by combining data from:
//  1. The IdentitySkeletons for the nodes that are _active_ in the target epoch
//     (recorded in EpochSetup event and immutable throughout the epoch).
//  2. The Dynamic Identities for the nodes that are _active_ in the target epoch (i.e. the dynamic identity
//     fields for the IdentitySkeletons contained in the EpochSetup event for the respective epoch).
//
// Optionally, identity information for an adjacent epoch is given if and only if an adjacent epoch exists. For
// a target epoch N, the epochs N-1 and N+1 are defined to be adjacent. Adjacent epochs do not necessarily exist
// (e.g. consider a spork comprising only a single epoch), in which case the respective inputs are nil or empty.
//  3. [optional] An adjacent epoch's IdentitySkeletons as recorded in the adjacent epoch's setup event.
//  4. [optional] An adjacent epoch's Dynamic Identities.
//  5. An adjacent epoch's identities participation status, this could be joining or leaving depending on epoch phase.
//
// The function enforces that the input slices pertaining to the same epoch contain the same identities
// (compared by nodeID) in the same order. Otherwise, an exception is returned.
// No errors are expected during normal operation. All errors indicate inconsistent or invalid inputs.
func BuildIdentityTable(
	targetEpochIdentitySkeletons IdentitySkeletonList,
	targetEpochDynamicIdentities DynamicIdentityEntryList,
	adjacentEpochIdentitySkeletons IdentitySkeletonList,
	adjacentEpochDynamicIdentities DynamicIdentityEntryList,
	adjacentIdentitiesStatus EpochParticipationStatus,
) (IdentityList, error) {
	if adjacentIdentitiesStatus != EpochParticipationStatusLeaving &&
		adjacentIdentitiesStatus != EpochParticipationStatusJoining {
		return nil, fmt.Errorf("invalid adjacent identity status, expect %s or %s, got %s",
			EpochParticipationStatusLeaving.String(),
			EpochParticipationStatusJoining.String(),
			adjacentIdentitiesStatus)
	}
	targetEpochParticipants, err := ComposeFullIdentities(targetEpochIdentitySkeletons, targetEpochDynamicIdentities, EpochParticipationStatusActive)
	if err != nil {
		return nil, fmt.Errorf("could not reconstruct participants for target epoch: %w", err)
	}
	adjacentEpochParticipants, err := ComposeFullIdentities(adjacentEpochIdentitySkeletons, adjacentEpochDynamicIdentities, adjacentIdentitiesStatus)
	if err != nil {
		return nil, fmt.Errorf("could not reconstruct participants for adjacent epoch: %w", err)
	}

	// Combine the participants of the current and adjacent epoch. The method `GenericIdentityList.Union`
	// already implements the following required conventions:
	//  1. Preference for IdentitySkeleton of the target epoch:
	//     In case an IdentitySkeleton with the same NodeID exists in the target epoch as well as
	//     in the adjacent epoch, we use the IdentitySkeleton for the target epoch (for example,
	//     to account for changes of keys, address, initial weight, etc).
	//  2. Canonical ordering
	return targetEpochParticipants.Union(adjacentEpochParticipants), nil
}

// DynamicIdentityEntryListFromIdentities converts IdentityList to DynamicIdentityEntryList.
func DynamicIdentityEntryListFromIdentities(identities IdentityList) DynamicIdentityEntryList {
	dynamicIdentities := make(DynamicIdentityEntryList, 0, len(identities))
	for _, identity := range identities {
		dynamicIdentities = append(dynamicIdentities, &DynamicIdentityEntry{
			NodeID:  identity.NodeID,
			Ejected: identity.IsEjected(),
		})
	}
	return dynamicIdentities
}

// ComposeFullIdentities combines identity skeletons and dynamic identities to produce a flow.IdentityList.
// It enforces that the input slices `skeletons` and `dynamics` list the same identities (compared by nodeID)
// in the same order. Otherwise, an exception is returned. For each identity i, we set
// `i.EpochParticipationStatus` to the `defaultEpochParticipationStatus` _unless_ i is ejected.
// No errors are expected during normal operations.
func ComposeFullIdentities(
	skeletons IdentitySkeletonList,
	dynamics DynamicIdentityEntryList,
	defaultEpochParticipationStatus EpochParticipationStatus,
) (IdentityList, error) {
	// sanity check: list of skeletons and dynamic should be the same
	if len(skeletons) != len(dynamics) {
		return nil, fmt.Errorf("invalid number of identities to reconstruct: expected %d, got %d", len(skeletons), len(dynamics))
	}

	// reconstruct identities from skeleton and dynamic parts
	var result IdentityList
	for i := range dynamics {
		// sanity check: identities should be sorted in the same order
		if dynamics[i].NodeID != skeletons[i].NodeID {
			return nil, fmt.Errorf("identites in protocol state are not consistently ordered: expected %s, got %s", skeletons[i].NodeID, dynamics[i].NodeID)
		}
		status := defaultEpochParticipationStatus
		if dynamics[i].Ejected {
			status = EpochParticipationStatusEjected
		}
		result = append(result, &Identity{
			IdentitySkeleton: *skeletons[i],
			DynamicIdentity: DynamicIdentity{
				EpochParticipationStatus: status,
			},
		})
	}
	return result, nil
}

// PSKeyValueStoreData is a binary blob with a version attached, specifying the format
// of the marshaled data. In a nutshell, it serves as a binary snapshot of a ProtocolKVStore.
// This structure is useful for version-agnostic storage, where snapshots with different versions
// can co-exist. The PSKeyValueStoreData is a generic format that can be later decoded to
// potentially different strongly typed structures based on version. When reading from the store,
// callers must know how to deal with the binary representation.
type PSKeyValueStoreData struct {
	Version uint64
	Data    []byte
}
