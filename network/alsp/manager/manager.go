package alspmgr

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/engine/common/worker"
	"github.com/onflow/flow-go/module"
	"github.com/onflow/flow-go/module/component"
	"github.com/onflow/flow-go/module/mempool/queue"
	"github.com/onflow/flow-go/module/metrics"
	"github.com/onflow/flow-go/network"
	"github.com/onflow/flow-go/network/alsp"
	"github.com/onflow/flow-go/network/alsp/internal"
	"github.com/onflow/flow-go/network/alsp/model"
	"github.com/onflow/flow-go/network/channels"
	"github.com/onflow/flow-go/utils/logging"
)

const (
	defaultMisbehaviorReportManagerWorkers = 2
	// DefaultSpamRecordQueueSize is the default size of the queue that stores the spam records to be processed by the
	// worker pool. The queue size should be large enough to handle the spam records during attacks. The recommended
	// size is 100 * number of nodes in the network. By default, the ALSP module will disallow-list the misbehaving
	// node after 100 spam reports are received (if no penalty value are amplified). Therefore, the queue size should
	// be at least 100 * number of nodes in the network.
	DefaultSpamRecordQueueSize = 100 * 1000
)

// MisbehaviorReportManager is responsible for handling misbehavior reports.
// The current version is at the minimum viable product stage and only logs the reports.
// TODO: the mature version should be able to handle the reports and take actions accordingly, i.e., penalize the misbehaving node
//
//	and report the node to be disallow-listed if the overall penalty of the misbehaving node drops below the disallow-listing threshold.
type MisbehaviorReportManager struct {
	component.Component
	logger  zerolog.Logger
	metrics module.AlspMetrics
	cache   alsp.SpamRecordCache
	// disablePenalty indicates whether applying the penalty to the misbehaving node is disabled.
	// When disabled, the ALSP module logs the misbehavior reports and updates the metrics, but does not apply the penalty.
	// This is useful for managing production incidents.
	// Note: under normal circumstances, the ALSP module should not be disabled.
	disablePenalty bool

	// workerPool is the worker pool for handling the misbehavior reports in a thread-safe and non-blocking manner.
	workerPool *worker.Pool[*internal.ReportedMisbehaviorWork]
}

var _ network.MisbehaviorReportManager = (*MisbehaviorReportManager)(nil)

type MisbehaviorReportManagerConfig struct {
	Logger zerolog.Logger
	// SpamRecordsCacheSize is the size of the spam record cache that stores the spam records for the authorized nodes.
	// It should be as big as the number of authorized nodes in Flow network.
	// Recommendation: for small network sizes 10 * number of authorized nodes to ensure that the cache can hold all the spam records of the authorized nodes.
	SpamRecordsCacheSize uint32
	// SpamReportQueueSize is the size of the queue that stores the spam records to be processed by the worker pool.
	SpamReportQueueSize uint32
	// AlspMetrics is the metrics instance for the alsp module (collecting spam reports).
	AlspMetrics module.AlspMetrics
	// HeroCacheMetricsFactory is the metrics factory for the HeroCache-related metrics.
	// Having factory as part of the config allows to create the metrics locally in the module.
	HeroCacheMetricsFactory metrics.HeroCacheMetricsFactory
	// DisablePenalty indicates whether applying the penalty to the misbehaving node is disabled.
	// When disabled, the ALSP module logs the misbehavior reports and updates the metrics, but does not apply the penalty.
	// This is useful for managing production incidents.
	// Note: under normal circumstances, the ALSP module should not be disabled.
	DisablePenalty bool
}

// validate validates the MisbehaviorReportManagerConfig instance. It returns an error if the config is invalid.
// It only validates the numeric fields of the config that may yield a stealth error in the production.
// It does not validate the struct fields of the config against a nil value.
// Args:
//
//	None.
//
// Returns:
//
//	An error if the config is invalid.
func (c MisbehaviorReportManagerConfig) validate() error {
	if c.SpamRecordsCacheSize == 0 {
		return fmt.Errorf("spam record cache size is not set")
	}
	if c.SpamReportQueueSize == 0 {
		return fmt.Errorf("spam report queue size is not set")
	}
	return nil
}

type MisbehaviorReportManagerOption func(*MisbehaviorReportManager)

// WithSpamRecordsCache sets the spam record cache for the MisbehaviorReportManager.
// Args:
//
//	cache: the spam record cache instance.
//
// Returns:
//
//	a MisbehaviorReportManagerOption that sets the spam record cache for the MisbehaviorReportManager.
//
// Note: this option is used for testing purposes. The production version of the MisbehaviorReportManager should use the
//
//	NewSpamRecordCache function to create the spam record cache.
func WithSpamRecordsCache(cache alsp.SpamRecordCache) MisbehaviorReportManagerOption {
	return func(m *MisbehaviorReportManager) {
		m.cache = cache
	}
}

// NewMisbehaviorReportManager creates a new instance of the MisbehaviorReportManager.
// Args:
//
//	logger: the logger instance.
//	metrics: the metrics instance.
//	cache: the spam record cache instance.
//
// Returns:
//
//		A new instance of the MisbehaviorReportManager.
//	 An error if the config is invalid. The error is considered irrecoverable.
func NewMisbehaviorReportManager(cfg *MisbehaviorReportManagerConfig, opts ...MisbehaviorReportManagerOption) (*MisbehaviorReportManager, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration for MisbehaviorReportManager: %w", err)
	}

	lg := cfg.Logger.With().Str("module", "misbehavior_report_manager").Logger()
	m := &MisbehaviorReportManager{
		logger:         lg,
		metrics:        cfg.AlspMetrics,
		disablePenalty: cfg.DisablePenalty,
	}

	m.cache = internal.NewSpamRecordCache(
		cfg.SpamRecordsCacheSize,
		lg.With().Str("component", "spam_record_cache").Logger(),
		metrics.ApplicationLayerSpamRecordCacheMetricFactory(cfg.HeroCacheMetricsFactory),
		model.SpamRecordFactory())

	store := queue.NewHeroStore(
		cfg.SpamReportQueueSize,
		lg.With().Str("component", "spam_record_queue").Logger(),
		metrics.ApplicationLayerSpamRecordQueueMetricsFactory(cfg.HeroCacheMetricsFactory))

	m.workerPool = worker.NewWorkerPoolBuilder[*internal.ReportedMisbehaviorWork](
		cfg.Logger,
		store,
		m.processMisbehaviorReport).Build()

	for _, opt := range opts {
		opt(m)
	}

	builder := component.NewComponentManagerBuilder()
	for i := 0; i < defaultMisbehaviorReportManagerWorkers; i++ {
		builder.AddWorker(m.workerPool.WorkerLogic())
	}

	m.Component = builder.Build()

	if m.disablePenalty {
		m.logger.Warn().Msg("penalty mechanism of alsp is disabled")
	}
	return m, nil
}

// HandleMisbehaviorReport is called upon a new misbehavior is reported.
// The implementation of this function should be thread-safe and non-blocking.
// Args:
//
//	channel: the channel on which the misbehavior is reported.
//	report: the misbehavior report.
//
// Returns:
//
//	none.
func (m *MisbehaviorReportManager) HandleMisbehaviorReport(channel channels.Channel, report network.MisbehaviorReport) {
	lg := m.logger.With().
		Str("channel", channel.String()).
		Hex("misbehaving_id", logging.ID(report.OriginId())).
		Str("reason", report.Reason().String()).
		Float64("penalty", report.Penalty()).Logger()
	m.metrics.OnMisbehaviorReported(channel.String(), report.Reason().String())

	if ok := m.workerPool.Submit(&internal.ReportedMisbehaviorWork{
		Channel:  channel,
		OriginId: report.OriginId(),
		Reason:   report.Reason(),
		Penalty:  report.Penalty(),
	}); !ok {
		lg.Warn().Msg("discarding misbehavior report because either the queue is full or the misbehavior report is duplicate")
	}
}

// processMisbehaviorReport is the worker function that processes the misbehavior reports.
// It is called by the worker pool.
// It applies the penalty to the misbehaving node and updates the spam record cache.
// Implementation must be thread-safe so that it can be called concurrently.
// Args:
//
//	report: the misbehavior report to be processed.
//
// Returns:
//
//		error: the error that occurred during the processing of the misbehavior report. The returned error is
//	 irrecoverable and the node should crash if it occurs (indicating a bug in the ALSP module).
func (m *MisbehaviorReportManager) processMisbehaviorReport(report *internal.ReportedMisbehaviorWork) error {
	lg := m.logger.With().
		Str("channel", report.Channel.String()).
		Hex("misbehaving_id", logging.ID(report.OriginId)).
		Str("reason", report.Reason.String()).
		Float64("penalty", report.Penalty).Logger()

	if m.disablePenalty {
		// when penalty mechanism disabled, the misbehavior is logged and metrics are updated,
		// but no further actions are taken.
		lg.Trace().Msg("discarding misbehavior report because alsp penalty is disabled")
		return nil
	}

	applyPenalty := func() (float64, error) {
		return m.cache.Adjust(report.OriginId, func(record model.ProtocolSpamRecord) (model.ProtocolSpamRecord, error) {
			if report.Penalty > 0 {
				// this should never happen, unless there is a bug in the misbehavior report handling logic.
				// we should crash the node in this case to prevent further misbehavior reports from being lost and fix the bug.
				return record, fmt.Errorf("penalty value is positive: %f", report.Penalty)
			}
			record.Penalty += report.Penalty // penalty value is negative. We add it to the current penalty.
			return record, nil
		})
	}

	init := func() {
		initialized := m.cache.Init(report.OriginId)
		lg.Trace().Bool("initialized", initialized).Msg("initialized spam record")
	}

	// we first try to apply the penalty to the spam record, if it does not exist, cache returns ErrSpamRecordNotFound.
	// in this case, we initialize the spam record and try to apply the penalty again. We use an optimistic update by
	// first assuming that the spam record exists and then initializing it if it does not exist. In this way, we avoid
	// acquiring the lock twice per misbehavior report, reducing the contention on the lock and improving the performance.
	updatedPenalty, err := internal.TryWithRecoveryIfHitError(internal.ErrSpamRecordNotFound, applyPenalty, init)
	if err != nil {
		// this should never happen, unless there is a bug in the spam record cache implementation.
		// we should crash the node in this case to prevent further misbehavior reports from being lost and fix the bug.
		return fmt.Errorf("failed to apply penalty to the spam record: %w", err)
	}

	lg.Debug().Float64("updated_penalty", updatedPenalty).Msg("misbehavior report handled")
	return nil
}
