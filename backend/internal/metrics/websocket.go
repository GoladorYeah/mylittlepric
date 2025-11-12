package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// WebSocket connection metrics
	WebSocketConnectionsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "websocket_connections_total",
			Help: "Total number of WebSocket connection attempts",
		},
	)

	WebSocketConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "websocket_connections_active",
			Help: "Current number of active WebSocket connections",
		},
	)

	WebSocketConnectionsFailed = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "websocket_connections_failed_total",
			Help: "Total number of failed WebSocket connections",
		},
	)

	WebSocketConnectionDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "websocket_connection_duration_seconds",
			Help:    "WebSocket connection duration in seconds",
			Buckets: []float64{1, 5, 10, 30, 60, 300, 600, 1800, 3600},
		},
	)

	// WebSocket message metrics
	WebSocketMessagesSent = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_messages_sent_total",
			Help: "Total number of messages sent via WebSocket",
		},
		[]string{"type"},
	)

	WebSocketMessagesReceived = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_messages_received_total",
			Help: "Total number of messages received via WebSocket",
		},
		[]string{"type"},
	)

	WebSocketMessagesSentFailed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_messages_sent_failed_total",
			Help: "Total number of failed message sends",
		},
		[]string{"type", "reason"},
	)

	WebSocketMessageDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "websocket_message_duration_seconds",
			Help:    "WebSocket message processing duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	WebSocketMessageQueueSize = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "websocket_message_queue_size",
			Help: "Current size of WebSocket message queue",
		},
	)

	// WebSocket rate limiting metrics
	WebSocketRateLimitExceeded = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_rate_limit_exceeded_total",
			Help: "Total number of WebSocket rate limit violations",
		},
		[]string{"level"}, // "connection" or "user"
	)

	// WebSocket broadcast metrics
	WebSocketBroadcastsSent = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "websocket_broadcasts_sent_total",
			Help: "Total number of broadcasts sent via Redis Pub/Sub",
		},
	)

	WebSocketBroadcastsReceived = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "websocket_broadcasts_received_total",
			Help: "Total number of broadcasts received via Redis Pub/Sub",
		},
	)
)
