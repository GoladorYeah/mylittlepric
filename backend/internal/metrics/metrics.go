package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/go-redis/v9"
)

var (
	// HTTP метрики
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// WebSocket метрики
	WebSocketConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "websocket_connections_active",
			Help: "Number of active WebSocket connections",
		},
	)

	WebSocketMessagesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_messages_total",
			Help: "Total number of WebSocket messages",
		},
		[]string{"type"}, // sent, received
	)

	// Database метрики
	DBQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"database", "operation"}, // postgresql/redis, select/insert/update/delete
	)

	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"database", "operation"},
	)

	RedisConnectionPoolActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "redis_connection_pool_active",
			Help: "Number of active Redis connections",
		},
	)

	RedisConnectionPoolIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "redis_connection_pool_idle",
			Help: "Number of idle Redis connections",
		},
	)

	// AI сервисы метрики
	AIRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ai_requests_total",
			Help: "Total number of AI API requests",
		},
		[]string{"service", "model", "status"}, // gemini/serp, gemini-2.0-flash/etc, success/error
	)

	AIRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ai_request_duration_seconds",
			Help:    "AI API request latency",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30},
		},
		[]string{"service", "model"},
	)

	AITokensUsed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ai_tokens_used_total",
			Help: "Total number of AI tokens used",
		},
		[]string{"service", "model", "type"}, // input/output
	)

	// Session метрики
	ActiveSessions = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_sessions_total",
			Help: "Number of active chat sessions",
		},
	)

	MessagesProcessedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "messages_processed_total",
			Help: "Total number of processed messages",
		},
		[]string{"status"}, // success/error
	)

	MessageProcessingDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "message_processing_duration_seconds",
			Help:    "Message processing duration",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30, 60},
		},
	)

	// Cleanup job метрики
	CleanupJobRunsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cleanup_job_runs_total",
			Help: "Total number of cleanup job runs",
		},
	)

	CleanupJobRecordsDeleted = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cleanup_job_records_deleted_total",
			Help: "Total number of records deleted by cleanup job",
		},
	)

	CleanupJobDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "cleanup_job_duration_seconds",
			Help:    "Cleanup job execution duration",
			Buckets: []float64{1, 5, 10, 30, 60, 120},
		},
	)

	// Error метрики
	ErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors",
		},
		[]string{"type", "source"}, // type: database/ai/validation, source: service name
	)
)

// UpdateRedisPoolMetrics обновляет метрики Redis connection pool
func UpdateRedisPoolMetrics(stats *redis.PoolStats) {
	RedisConnectionPoolActive.Set(float64(stats.TotalConns - stats.IdleConns))
	RedisConnectionPoolIdle.Set(float64(stats.IdleConns))
}
