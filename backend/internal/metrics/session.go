package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Session metrics
	SessionsCreated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "sessions_created_total",
			Help: "Total number of sessions created",
		},
	)

	ActiveSessionsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_sessions_total",
			Help: "Current number of active sessions",
		},
	)

	SessionLookups = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "session_lookups_total",
			Help: "Total number of session lookups",
		},
	)

	SessionCacheMiss = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "session_cache_miss_total",
			Help: "Total number of session cache misses",
		},
	)

	SessionCacheHit = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "session_cache_hit_total",
			Help: "Total number of session cache hits",
		},
	)

	// Session cleanup metrics
	SessionCleanupLastSuccessTimestamp = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "session_cleanup_last_success_timestamp",
			Help: "Timestamp of last successful session cleanup",
		},
	)

	SessionCleanupDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "session_cleanup_duration_seconds",
			Help:    "Duration of session cleanup job in seconds",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30, 60},
		},
	)

	SessionsDeleted = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "sessions_deleted_total",
			Help: "Total number of sessions deleted by cleanup job",
		},
	)

	// Session sync metrics
	SessionSyncFailed = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "session_sync_failed_total",
			Help: "Total number of failed Redis-PostgreSQL session syncs",
		},
	)

	// Message persistence metrics
	MessagesPersisted = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "messages_persisted_total",
			Help: "Total number of messages persisted to database",
		},
	)

	MessagePersistenceFailed = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "message_persistence_failed_total",
			Help: "Total number of failed message persistence operations",
		},
	)

	MessageCacheMiss = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "message_cache_miss_total",
			Help: "Total number of message cache misses",
		},
	)

	MessageCacheHit = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "message_cache_hit_total",
			Help: "Total number of message cache hits",
		},
	)

	MessagesDeleted = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "messages_deleted_total",
			Help: "Total number of messages deleted by cleanup job",
		},
	)
)
