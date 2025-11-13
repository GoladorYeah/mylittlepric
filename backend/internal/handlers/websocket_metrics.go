package handlers

import (
	"time"

	"mylittleprice/internal/metrics"
)

// recordConnectionStart records WebSocket connection metrics when a client connects
func (h *WSHandler) recordConnectionStart() func() {
	start := time.Now()
	metrics.WebSocketConnectionsTotal.Inc()
	metrics.WebSocketConnectionsActive.Inc()

	// Return a cleanup function to be called on disconnect
	return func() {
		duration := time.Since(start).Seconds()
		metrics.WebSocketConnectionsActive.Dec()
		metrics.WebSocketConnectionDuration.Observe(duration)
	}
}

// recordConnectionFailed records failed WebSocket connection attempt
func (h *WSHandler) recordConnectionFailed() {
	metrics.WebSocketConnectionsFailed.Inc()
}

// recordMessageReceived records received WebSocket message
func (h *WSHandler) recordMessageReceived(msgType string) {
	// Prevent empty label values which cause Prometheus errors
	if msgType == "" {
		msgType = "unknown"
	}
	metrics.WebSocketMessagesReceived.WithLabelValues(msgType).Inc()
}

// recordMessageSent records sent WebSocket message
func (h *WSHandler) recordMessageSent(msgType string) {
	// Prevent empty label values which cause Prometheus errors
	if msgType == "" {
		msgType = "unknown"
	}
	metrics.WebSocketMessagesSent.WithLabelValues(msgType).Inc()
}

// recordMessageSendFailed records failed message send
func (h *WSHandler) recordMessageSendFailed(msgType string, reason string) {
	// Prevent empty label values which cause Prometheus errors
	if msgType == "" {
		msgType = "unknown"
	}
	if reason == "" {
		reason = "unknown"
	}
	metrics.WebSocketMessagesSentFailed.WithLabelValues(msgType, reason).Inc()
}

// recordMessageProcessing returns a function to record message processing time
func (h *WSHandler) recordMessageProcessing() func() {
	start := time.Now()
	return func() {
		duration := time.Since(start).Seconds()
		metrics.WebSocketMessageDuration.Observe(duration)
	}
}

// recordRateLimitViolation records WebSocket rate limit violation
func (h *WSHandler) recordRateLimitViolation(level string) {
	// Prevent empty label values which cause Prometheus errors
	if level == "" {
		level = "unknown"
	}
	metrics.WebSocketRateLimitExceeded.WithLabelValues(level).Inc()
}

// recordBroadcastSent records broadcast sent via Redis Pub/Sub
func (h *WSHandler) recordBroadcastSent() {
	metrics.WebSocketBroadcastsSent.Inc()
}

// recordBroadcastReceived records broadcast received via Redis Pub/Sub
func (h *WSHandler) recordBroadcastReceived() {
	metrics.WebSocketBroadcastsReceived.Inc()
}
