package metrics

import (
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// WebSocket connection metrics
	WebSocketConnectionsTotal prometheus.Counter
	WebSocketConnectionsActive prometheus.Gauge
	WebSocketConnectionsFailed prometheus.Counter
	WebSocketConnectionDuration prometheus.Histogram

	// WebSocket message metrics
	WebSocketMessagesSent *prometheus.CounterVec
	WebSocketMessagesReceived *prometheus.CounterVec
	WebSocketMessagesSentFailed *prometheus.CounterVec
	WebSocketMessageDuration prometheus.Histogram
	WebSocketMessageQueueSize prometheus.Gauge

	// WebSocket rate limiting metrics
	WebSocketRateLimitExceeded *prometheus.CounterVec

	// WebSocket broadcast metrics
	WebSocketBroadcastsSent prometheus.Counter
	WebSocketBroadcastsReceived prometheus.Counter

	// Ensure metrics are registered only once
	wsMetricsOnce sync.Once
)

// RegisterWebSocketMetrics registers all WebSocket metrics to default registry
func RegisterWebSocketMetrics() {
	wsMetricsOnce.Do(func() {
		log.Printf("🔧 Registering WebSocket metrics")
		// WebSocket connection metrics
		WebSocketConnectionsTotal = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "websocket_connections_total",
				Help: "Total number of WebSocket connection attempts",
			},
		)
		prometheus.MustRegister(WebSocketConnectionsTotal)

		WebSocketConnectionsActive = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "websocket_connections_active",
				Help: "Current number of active WebSocket connections",
			},
		)
		prometheus.MustRegister(WebSocketConnectionsActive)

		WebSocketConnectionsFailed = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "websocket_connections_failed_total",
				Help: "Total number of failed WebSocket connections",
			},
		)
		prometheus.MustRegister(WebSocketConnectionsFailed)

		WebSocketConnectionDuration = prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "websocket_connection_duration_seconds",
				Help:    "WebSocket connection duration in seconds",
				Buckets: []float64{1, 5, 10, 30, 60, 300, 600, 1800, 3600},
			},
		)
		prometheus.MustRegister(WebSocketConnectionDuration)

		// WebSocket message metrics
		WebSocketMessagesSent = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "websocket_messages_sent_total",
				Help: "Total number of messages sent via WebSocket",
			},
			[]string{"type"},
		)
		prometheus.MustRegister(WebSocketMessagesSent)

		WebSocketMessagesReceived = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "websocket_messages_received_total",
				Help: "Total number of messages received via WebSocket",
			},
			[]string{"type"},
		)
		prometheus.MustRegister(WebSocketMessagesReceived)

		WebSocketMessagesSentFailed = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "websocket_messages_sent_failed_total",
				Help: "Total number of failed message sends",
			},
			[]string{"type", "reason"},
		)
		prometheus.MustRegister(WebSocketMessagesSentFailed)

		WebSocketMessageDuration = prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "websocket_message_duration_seconds",
				Help:    "WebSocket message processing duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
		)
		prometheus.MustRegister(WebSocketMessageDuration)

		WebSocketMessageQueueSize = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "websocket_message_queue_size",
				Help: "Current size of WebSocket message queue",
			},
		)
		prometheus.MustRegister(WebSocketMessageQueueSize)

		// WebSocket rate limiting metrics
		WebSocketRateLimitExceeded = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "websocket_rate_limit_exceeded_total",
				Help: "Total number of WebSocket rate limit violations",
			},
			[]string{"level"}, // "connection" or "user"
		)
		prometheus.MustRegister(WebSocketRateLimitExceeded)

		// WebSocket broadcast metrics
		WebSocketBroadcastsSent = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "websocket_broadcasts_sent_total",
				Help: "Total number of broadcasts sent via Redis Pub/Sub",
			},
		)
		prometheus.MustRegister(WebSocketBroadcastsSent)

		WebSocketBroadcastsReceived = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "websocket_broadcasts_received_total",
				Help: "Total number of broadcasts received via Redis Pub/Sub",
			},
		)
		prometheus.MustRegister(WebSocketBroadcastsReceived)

		log.Printf("✅ WebSocket metrics registered successfully")
	})
}
