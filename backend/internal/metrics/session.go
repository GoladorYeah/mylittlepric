package metrics

import (
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Session metrics
	SessionsCreated prometheus.Counter
	ActiveSessionsTotal prometheus.Gauge
	SessionLookups prometheus.Counter
	SessionCacheMiss prometheus.Counter
	SessionCacheHit prometheus.Counter

	// Session cleanup metrics
	SessionCleanupLastSuccessTimestamp prometheus.Gauge
	SessionCleanupDuration prometheus.Histogram
	SessionsDeleted prometheus.Counter

	// Session sync metrics
	SessionSyncFailed prometheus.Counter

	// Message persistence metrics
	MessagesPersisted prometheus.Counter
	MessagePersistenceFailed prometheus.Counter
	MessageCacheMiss prometheus.Counter
	MessageCacheHit prometheus.Counter
	MessagesDeleted prometheus.Counter

	// Ensure metrics are registered only once
	sessionMetricsOnce sync.Once
)

// RegisterSessionMetrics registers all session metrics to default registry
func RegisterSessionMetrics() {
	sessionMetricsOnce.Do(func() {
		log.Printf("🔧 Registering Session metrics")

		// Session metrics
		SessionsCreated = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "sessions_created_total",
				Help: "Total number of sessions created",
			},
		)
		prometheus.MustRegister(SessionsCreated)

		ActiveSessionsTotal = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_sessions_total",
				Help: "Current number of active sessions",
			},
		)
		prometheus.MustRegister(ActiveSessionsTotal)

		SessionLookups = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "session_lookups_total",
				Help: "Total number of session lookups",
			},
		)
		prometheus.MustRegister(SessionLookups)

		SessionCacheMiss = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "session_cache_miss_total",
				Help: "Total number of session cache misses",
			},
		)
		prometheus.MustRegister(SessionCacheMiss)

		SessionCacheHit = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "session_cache_hit_total",
				Help: "Total number of session cache hits",
			},
		)
		prometheus.MustRegister(SessionCacheHit)

		// Session cleanup metrics
		SessionCleanupLastSuccessTimestamp = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "session_cleanup_last_success_timestamp",
				Help: "Timestamp of last successful session cleanup",
			},
		)
		prometheus.MustRegister(SessionCleanupLastSuccessTimestamp)

		SessionCleanupDuration = prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "session_cleanup_duration_seconds",
				Help:    "Duration of session cleanup job in seconds",
				Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30, 60},
			},
		)
		prometheus.MustRegister(SessionCleanupDuration)

		SessionsDeleted = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "sessions_deleted_total",
				Help: "Total number of sessions deleted by cleanup job",
			},
		)
		prometheus.MustRegister(SessionsDeleted)

		// Session sync metrics
		SessionSyncFailed = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "session_sync_failed_total",
				Help: "Total number of failed Redis-PostgreSQL session syncs",
			},
		)
		prometheus.MustRegister(SessionSyncFailed)

		// Message persistence metrics
		MessagesPersisted = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "messages_persisted_total",
				Help: "Total number of messages persisted to database",
			},
		)
		prometheus.MustRegister(MessagesPersisted)

		MessagePersistenceFailed = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "message_persistence_failed_total",
				Help: "Total number of failed message persistence operations",
			},
		)
		prometheus.MustRegister(MessagePersistenceFailed)

		MessageCacheMiss = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "message_cache_miss_total",
				Help: "Total number of message cache misses",
			},
		)
		prometheus.MustRegister(MessageCacheMiss)

		MessageCacheHit = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "message_cache_hit_total",
				Help: "Total number of message cache hits",
			},
		)
		prometheus.MustRegister(MessageCacheHit)

		MessagesDeleted = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "messages_deleted_total",
				Help: "Total number of messages deleted by cleanup job",
			},
		)
		prometheus.MustRegister(MessagesDeleted)

		log.Printf("✅ Session metrics registered successfully")
	})
}
