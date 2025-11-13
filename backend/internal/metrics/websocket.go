package metrics

import (
	"log"

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
)

// init registers metrics once
func init() {
	// WebSocket connection metrics
	WebSocketConnectionsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "websocket_connections_total",
			Help: "Total number of WebSocket connection attempts",
		},
	)
	if err := prometheus.Register(WebSocketConnectionsTotal); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			WebSocketConnectionsTotal = are.ExistingCollector.(prometheus.Counter)
			log.Printf("⚠️ websocket_connections_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register websocket_connections_total: %v", err)
		}
	}

	WebSocketConnectionsActive = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "websocket_connections_active",
			Help: "Current number of active WebSocket connections",
		},
	)
	if err := prometheus.Register(WebSocketConnectionsActive); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			WebSocketConnectionsActive = are.ExistingCollector.(prometheus.Gauge)
			log.Printf("⚠️ websocket_connections_active already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register websocket_connections_active: %v", err)
		}
	}

	WebSocketConnectionsFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "websocket_connections_failed_total",
			Help: "Total number of failed WebSocket connections",
		},
	)
	if err := prometheus.Register(WebSocketConnectionsFailed); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			WebSocketConnectionsFailed = are.ExistingCollector.(prometheus.Counter)
			log.Printf("⚠️ websocket_connections_failed_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register websocket_connections_failed_total: %v", err)
		}
	}

	WebSocketConnectionDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "websocket_connection_duration_seconds",
			Help:    "WebSocket connection duration in seconds",
			Buckets: []float64{1, 5, 10, 30, 60, 300, 600, 1800, 3600},
		},
	)
	if err := prometheus.Register(WebSocketConnectionDuration); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			WebSocketConnectionDuration = are.ExistingCollector.(prometheus.Histogram)
			log.Printf("⚠️ websocket_connection_duration_seconds already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register websocket_connection_duration_seconds: %v", err)
		}
	}

	// WebSocket message metrics
	WebSocketMessagesSent = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_messages_sent_total",
			Help: "Total number of messages sent via WebSocket",
		},
		[]string{"type"},
	)
	if err := prometheus.Register(WebSocketMessagesSent); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			WebSocketMessagesSent = are.ExistingCollector.(*prometheus.CounterVec)
			log.Printf("⚠️ websocket_messages_sent_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register websocket_messages_sent_total: %v", err)
		}
	}

	WebSocketMessagesReceived = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_messages_received_total",
			Help: "Total number of messages received via WebSocket",
		},
		[]string{"type"},
	)
	if err := prometheus.Register(WebSocketMessagesReceived); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			WebSocketMessagesReceived = are.ExistingCollector.(*prometheus.CounterVec)
			log.Printf("⚠️ websocket_messages_received_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register websocket_messages_received_total: %v", err)
		}
	}

	WebSocketMessagesSentFailed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_messages_sent_failed_total",
			Help: "Total number of failed message sends",
		},
		[]string{"type", "reason"},
	)
	if err := prometheus.Register(WebSocketMessagesSentFailed); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			WebSocketMessagesSentFailed = are.ExistingCollector.(*prometheus.CounterVec)
			log.Printf("⚠️ websocket_messages_sent_failed_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register websocket_messages_sent_failed_total: %v", err)
		}
	}

	WebSocketMessageDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "websocket_message_duration_seconds",
			Help:    "WebSocket message processing duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
	if err := prometheus.Register(WebSocketMessageDuration); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			WebSocketMessageDuration = are.ExistingCollector.(prometheus.Histogram)
			log.Printf("⚠️ websocket_message_duration_seconds already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register websocket_message_duration_seconds: %v", err)
		}
	}

	WebSocketMessageQueueSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "websocket_message_queue_size",
			Help: "Current size of WebSocket message queue",
		},
	)
	if err := prometheus.Register(WebSocketMessageQueueSize); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			WebSocketMessageQueueSize = are.ExistingCollector.(prometheus.Gauge)
			log.Printf("⚠️ websocket_message_queue_size already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register websocket_message_queue_size: %v", err)
		}
	}

	// WebSocket rate limiting metrics
	WebSocketRateLimitExceeded = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_rate_limit_exceeded_total",
			Help: "Total number of WebSocket rate limit violations",
		},
		[]string{"level"}, // "connection" or "user"
	)
	if err := prometheus.Register(WebSocketRateLimitExceeded); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			WebSocketRateLimitExceeded = are.ExistingCollector.(*prometheus.CounterVec)
			log.Printf("⚠️ websocket_rate_limit_exceeded_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register websocket_rate_limit_exceeded_total: %v", err)
		}
	}

	// WebSocket broadcast metrics
	WebSocketBroadcastsSent = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "websocket_broadcasts_sent_total",
			Help: "Total number of broadcasts sent via Redis Pub/Sub",
		},
	)
	if err := prometheus.Register(WebSocketBroadcastsSent); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			WebSocketBroadcastsSent = are.ExistingCollector.(prometheus.Counter)
			log.Printf("⚠️ websocket_broadcasts_sent_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register websocket_broadcasts_sent_total: %v", err)
		}
	}

	WebSocketBroadcastsReceived = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "websocket_broadcasts_received_total",
			Help: "Total number of broadcasts received via Redis Pub/Sub",
		},
	)
	if err := prometheus.Register(WebSocketBroadcastsReceived); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			WebSocketBroadcastsReceived = are.ExistingCollector.(prometheus.Counter)
			log.Printf("⚠️ websocket_broadcasts_received_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register websocket_broadcasts_received_total: %v", err)
		}
	}
}
