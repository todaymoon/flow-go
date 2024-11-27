// Package pusher implements an engine for providing access to resources held
// by the collection node, including collections, collection guarantees, and
// transactions.
package pusher

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/engine"
	"github.com/onflow/flow-go/engine/common/fifoqueue"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/model/flow/filter"
	"github.com/onflow/flow-go/model/messages"
	"github.com/onflow/flow-go/module"
	"github.com/onflow/flow-go/module/component"
	"github.com/onflow/flow-go/module/irrecoverable"
	"github.com/onflow/flow-go/module/metrics"
	"github.com/onflow/flow-go/network"
	"github.com/onflow/flow-go/network/channels"
	"github.com/onflow/flow-go/state/protocol"
	"github.com/onflow/flow-go/storage"
	"github.com/onflow/flow-go/utils/logging"
)

// Engine is the collection pusher engine, which provides access to resources
// held by the collection node.
type Engine struct {
	log          zerolog.Logger
	engMetrics   module.EngineMetrics
	colMetrics   module.CollectionMetrics
	conduit      network.Conduit
	me           module.Local
	state        protocol.State
	collections  storage.Collections
	transactions storage.Transactions

	messageHandler *engine.MessageHandler
	notifier       engine.Notifier
	inbound        *fifoqueue.FifoQueue

	component.Component
	cm *component.ComponentManager
}

// TODO convert to network.MessageProcessor
var _ network.Engine = (*Engine)(nil)
var _ component.Component = (*Engine)(nil)

func New(log zerolog.Logger, net network.EngineRegistry, state protocol.State, engMetrics module.EngineMetrics, colMetrics module.CollectionMetrics, me module.Local, collections storage.Collections, transactions storage.Transactions) (*Engine, error) {
	// TODO length observer metrics
	inbound, err := fifoqueue.NewFifoQueue(1000)
	if err != nil {
		return nil, fmt.Errorf("could not create inbound fifoqueue: %w", err)
	}

	notifier := engine.NewNotifier()
	messageHandler := engine.NewMessageHandler(log, notifier, engine.Pattern{
		Match: engine.MatchType[*messages.SubmitCollectionGuarantee],
		Store: &engine.FifoMessageStore{FifoQueue: inbound},
	})

	e := &Engine{
		log:          log.With().Str("engine", "pusher").Logger(),
		engMetrics:   engMetrics,
		colMetrics:   colMetrics,
		me:           me,
		state:        state,
		collections:  collections,
		transactions: transactions,

		messageHandler: messageHandler,
		notifier:       notifier,
		inbound:        inbound,
	}

	conduit, err := net.Register(channels.PushGuarantees, e)
	if err != nil {
		return nil, fmt.Errorf("could not register for push protocol: %w", err)
	}
	e.conduit = conduit

	e.cm = component.NewComponentManagerBuilder().
		AddWorker(e.outboundQueueWorker).
		Build()
	e.Component = e.cm

	return e, nil
}

// Worker to process SubmitCollectionGuarantee messages coming from the Finalizer.
func (e *Engine) outboundQueueWorker(ctx irrecoverable.SignalerContext, ready component.ReadyFunc) {
	ready()

	done := ctx.Done()
	wake := e.notifier.Channel()
	for {
		select {
		case <-done:
			return
		case <-wake:
			err := e.processOutboundMessages(ctx)
			if err != nil {
				ctx.Throw(err)
			}
		}
	}
}

// processOutboundMessages processes any available messages from the queue.
// Only returns when the queue is empty (or the engine is terminated).
func (e *Engine) processOutboundMessages(ctx context.Context) error {
	for {
		nextMessage, ok := e.inbound.Pop()
		if !ok {
			return nil
		}

		asEngineWrapper := nextMessage.(*engine.Message)
		asSCGMsg := asEngineWrapper.Payload.(*messages.SubmitCollectionGuarantee)
		originID := asEngineWrapper.OriginID

		err := e.process(originID, asSCGMsg)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}
}

// SubmitLocal submits an event originating on the local node.
func (e *Engine) SubmitLocal(event interface{}) {
	err := e.messageHandler.Process(e.me.NodeID(), event)
	if err != nil {
		engine.LogError(e.log, err)
	}
}

// Submit submits the given event from the node with the given origin ID
// for processing in a non-blocking manner. It returns instantly and logs
// a potential processing error internally when done.
func (e *Engine) Submit(channel channels.Channel, originID flow.Identifier, event interface{}) {
	engine.LogError(e.log, fmt.Errorf("pusher engine should only receive local messages on the same node"))
}

// ProcessLocal processes an event originating on the local node.
func (e *Engine) ProcessLocal(event interface{}) error {
	return e.messageHandler.Process(e.me.NodeID(), event)
}

// Process processes the given event from the node with the given origin ID in
// a blocking manner. It returns the potential processing error when done.
func (e *Engine) Process(channel channels.Channel, originID flow.Identifier, message any) error {
	return fmt.Errorf("pusher engine should only receive local messages on the same node")
}

// process processes events for the pusher engine on the collection node.
func (e *Engine) process(originID flow.Identifier, event interface{}) error {
	switch ev := event.(type) {
	case *messages.SubmitCollectionGuarantee:
		e.engMetrics.MessageReceived(metrics.EngineCollectionProvider, metrics.MessageSubmitGuarantee)
		defer e.engMetrics.MessageHandled(metrics.EngineCollectionProvider, metrics.MessageSubmitGuarantee)
		if originID != e.me.NodeID() {
			return fmt.Errorf("invalid remote request to submit collection guarantee (from=%x)", originID)
		}
		return e.SubmitCollectionGuarantee(&ev.Guarantee)
	default:
		return fmt.Errorf("invalid event type (%T)", event)
	}
}

// SubmitCollectionGuarantee submits the collection guarantee to all consensus nodes.
func (e *Engine) SubmitCollectionGuarantee(guarantee *flow.CollectionGuarantee) error {
	consensusNodes, err := e.state.Final().Identities(filter.HasRole[flow.Identity](flow.RoleConsensus))
	if err != nil {
		return fmt.Errorf("could not get consensus nodes: %w", err)
	}

	// NOTE: Consensus nodes do not broadcast guarantees among themselves, so it needs that
	// at least one collection node make a publish to all of them.
	err = e.conduit.Publish(guarantee, consensusNodes.NodeIDs()...)
	if err != nil {
		return fmt.Errorf("could not submit collection guarantee: %w", err)
	}

	e.engMetrics.MessageSent(metrics.EngineCollectionProvider, metrics.MessageCollectionGuarantee)

	e.log.Debug().
		Hex("guarantee_id", logging.ID(guarantee.ID())).
		Hex("ref_block_id", logging.ID(guarantee.ReferenceBlockID)).
		Msg("submitting collection guarantee")

	return nil
}
