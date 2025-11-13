package metrics

import (
	"log"

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
)

// init registers metrics once
func init() {
	// Session metrics
	SessionsCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sessions_created_total",
			Help: "Total number of sessions created",
		},
	)
	if err := prometheus.Register(SessionsCreated); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			SessionsCreated = are.ExistingCollector.(prometheus.Counter)
			log.Printf("⚠️ sessions_created_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register sessions_created_total: %v", err)
		}
	}

	ActiveSessionsTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_sessions_total",
			Help: "Current number of active sessions",
		},
	)
	if err := prometheus.Register(ActiveSessionsTotal); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			ActiveSessionsTotal = are.ExistingCollector.(prometheus.Gauge)
			log.Printf("⚠️ active_sessions_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register active_sessions_total: %v", err)
		}
	}

	SessionLookups = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "session_lookups_total",
			Help: "Total number of session lookups",
		},
	)
	if err := prometheus.Register(SessionLookups); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			SessionLookups = are.ExistingCollector.(prometheus.Counter)
			log.Printf("⚠️ session_lookups_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register session_lookups_total: %v", err)
		}
	}

	SessionCacheMiss = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "session_cache_miss_total",
			Help: "Total number of session cache misses",
		},
	)
	if err := prometheus.Register(SessionCacheMiss); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			SessionCacheMiss = are.ExistingCollector.(prometheus.Counter)
			log.Printf("⚠️ session_cache_miss_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register session_cache_miss_total: %v", err)
		}
	}

	SessionCacheHit = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "session_cache_hit_total",
			Help: "Total number of session cache hits",
		},
	)
	if err := prometheus.Register(SessionCacheHit); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			SessionCacheHit = are.ExistingCollector.(prometheus.Counter)
			log.Printf("⚠️ session_cache_hit_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register session_cache_hit_total: %v", err)
		}
	}

	// Session cleanup metrics
	SessionCleanupLastSuccessTimestamp = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "session_cleanup_last_success_timestamp",
			Help: "Timestamp of last successful session cleanup",
		},
	)
	if err := prometheus.Register(SessionCleanupLastSuccessTimestamp); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			SessionCleanupLastSuccessTimestamp = are.ExistingCollector.(prometheus.Gauge)
			log.Printf("⚠️ session_cleanup_last_success_timestamp already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register session_cleanup_last_success_timestamp: %v", err)
		}
	}

	SessionCleanupDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "session_cleanup_duration_seconds",
			Help:    "Duration of session cleanup job in seconds",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30, 60},
		},
	)
	if err := prometheus.Register(SessionCleanupDuration); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			SessionCleanupDuration = are.ExistingCollector.(prometheus.Histogram)
			log.Printf("⚠️ session_cleanup_duration_seconds already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register session_cleanup_duration_seconds: %v", err)
		}
	}

	SessionsDeleted = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sessions_deleted_total",
			Help: "Total number of sessions deleted by cleanup job",
		},
	)
	if err := prometheus.Register(SessionsDeleted); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			SessionsDeleted = are.ExistingCollector.(prometheus.Counter)
			log.Printf("⚠️ sessions_deleted_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register sessions_deleted_total: %v", err)
		}
	}

	// Session sync metrics
	SessionSyncFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "session_sync_failed_total",
			Help: "Total number of failed Redis-PostgreSQL session syncs",
		},
	)
	if err := prometheus.Register(SessionSyncFailed); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			SessionSyncFailed = are.ExistingCollector.(prometheus.Counter)
			log.Printf("⚠️ session_sync_failed_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register session_sync_failed_total: %v", err)
		}
	}

	// Message persistence metrics
	MessagesPersisted = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "messages_persisted_total",
			Help: "Total number of messages persisted to database",
		},
	)
	if err := prometheus.Register(MessagesPersisted); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			MessagesPersisted = are.ExistingCollector.(prometheus.Counter)
			log.Printf("⚠️ messages_persisted_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register messages_persisted_total: %v", err)
		}
	}

	MessagePersistenceFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "message_persistence_failed_total",
			Help: "Total number of failed message persistence operations",
		},
	)
	if err := prometheus.Register(MessagePersistenceFailed); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			MessagePersistenceFailed = are.ExistingCollector.(prometheus.Counter)
			log.Printf("⚠️ message_persistence_failed_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register message_persistence_failed_total: %v", err)
		}
	}

	MessageCacheMiss = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "message_cache_miss_total",
			Help: "Total number of message cache misses",
		},
	)
	if err := prometheus.Register(MessageCacheMiss); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			MessageCacheMiss = are.ExistingCollector.(prometheus.Counter)
			log.Printf("⚠️ message_cache_miss_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register message_cache_miss_total: %v", err)
		}
	}

	MessageCacheHit = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "message_cache_hit_total",
			Help: "Total number of message cache hits",
		},
	)
	if err := prometheus.Register(MessageCacheHit); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			MessageCacheHit = are.ExistingCollector.(prometheus.Counter)
			log.Printf("⚠️ message_cache_hit_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register message_cache_hit_total: %v", err)
		}
	}

	MessagesDeleted = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "messages_deleted_total",
			Help: "Total number of messages deleted by cleanup job",
		},
	)
	if err := prometheus.Register(MessagesDeleted); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			MessagesDeleted = are.ExistingCollector.(prometheus.Counter)
			log.Printf("⚠️ messages_deleted_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register messages_deleted_total: %v", err)
		}
	}
}
