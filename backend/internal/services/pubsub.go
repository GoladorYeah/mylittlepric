package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"mylittleprice/internal/models"
)

// PubSubService handles Redis Pub/Sub for cross-server WebSocket communication
type PubSubService struct {
	redis      *redis.Client
	ctx        context.Context
	serverID   string // Unique ID for this server instance
	mu         sync.RWMutex
	handlers   map[string]BroadcastHandler // Channel -> handler mapping
	cancelFunc context.CancelFunc
}

// BroadcastHandler is a callback function for handling broadcast messages
type BroadcastHandler func(msg *BroadcastMessage)

// BroadcastMessage represents a message to be broadcast across servers
type BroadcastMessage struct {
	ServerID  string       `json:"server_id"`  // ID of server that sent the message
	UserID    uuid.UUID    `json:"user_id"`    // Target user ID
	SessionID string       `json:"session_id"` // Session ID
	Type      string       `json:"type"`       // Message type
	Payload   interface{}  `json:"payload"`    // Message payload
}

// NewPubSubService creates a new PubSubService
func NewPubSubService(redis *redis.Client) *PubSubService {
	ctx, cancel := context.WithCancel(context.Background())

	return &PubSubService{
		redis:      redis,
		ctx:        ctx,
		serverID:   uuid.New().String(), // Generate unique server ID
		handlers:   make(map[string]BroadcastHandler),
		cancelFunc: cancel,
	}
}

// GetServerID returns this server's unique ID
func (s *PubSubService) GetServerID() string {
	return s.serverID
}

// Subscribe subscribes to a channel and registers a handler
func (s *PubSubService) Subscribe(channel string, handler BroadcastHandler) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Register handler
	s.handlers[channel] = handler

	// Subscribe to Redis channel
	pubsub := s.redis.Subscribe(s.ctx, channel)

	// Start listening in a goroutine
	go s.listen(channel, pubsub)

	log.Printf("🔔 Server %s subscribed to channel: %s", s.serverID[:8], channel)
	return nil
}

// listen listens for messages on a Redis channel
func (s *PubSubService) listen(channel string, pubsub *redis.PubSub) {
	defer pubsub.Close()

	ch := pubsub.Channel()

	for {
		select {
		case <-s.ctx.Done():
			log.Printf("🔕 Server %s stopped listening to channel: %s", s.serverID[:8], channel)
			return
		case msg := <-ch:
			if msg == nil {
				continue
			}

			// Parse message
			var broadcastMsg BroadcastMessage
			if err := json.Unmarshal([]byte(msg.Payload), &broadcastMsg); err != nil {
				log.Printf("❌ Failed to unmarshal broadcast message: %v", err)
				continue
			}

			// Skip messages from this server (avoid echo)
			if broadcastMsg.ServerID == s.serverID {
				continue
			}

			// Get handler for this channel
			s.mu.RLock()
			handler, exists := s.handlers[channel]
			s.mu.RUnlock()

			if !exists {
				log.Printf("⚠️ No handler registered for channel: %s", channel)
				continue
			}

			// Call handler
			handler(&broadcastMsg)
		}
	}
}

// Publish publishes a message to a channel
func (s *PubSubService) Publish(channel string, msg *BroadcastMessage) error {
	// Set server ID
	msg.ServerID = s.serverID

	// Marshal message
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %w", err)
	}

	// Publish to Redis
	if err := s.redis.Publish(s.ctx, channel, data).Err(); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// BroadcastToUser broadcasts a message to all servers for a specific user
func (s *PubSubService) BroadcastToUser(userID uuid.UUID, msgType string, payload interface{}) error {
	msg := &BroadcastMessage{
		UserID:  userID,
		Type:    msgType,
		Payload: payload,
	}

	// Use a channel per user for efficient targeting
	channel := fmt.Sprintf("user:%s", userID.String())
	return s.Publish(channel, msg)
}

// BroadcastToSession broadcasts a message to all servers for a specific session
func (s *PubSubService) BroadcastToSession(sessionID string, msgType string, payload interface{}) error {
	msg := &BroadcastMessage{
		SessionID: sessionID,
		Type:      msgType,
		Payload:   payload,
	}

	channel := fmt.Sprintf("session:%s", sessionID)
	return s.Publish(channel, msg)
}

// SubscribeToUser subscribes to all messages for a specific user
func (s *PubSubService) SubscribeToUser(userID uuid.UUID, handler BroadcastHandler) error {
	channel := fmt.Sprintf("user:%s", userID.String())
	return s.Subscribe(channel, handler)
}

// SubscribeToSession subscribes to all messages for a specific session
func (s *PubSubService) SubscribeToSession(sessionID string, handler BroadcastHandler) error {
	channel := fmt.Sprintf("session:%s", sessionID)
	return s.Subscribe(channel, handler)
}

// SubscribeToAllUsers subscribes to a global user broadcast channel
// This is used for WebSocket handlers to receive all user messages
func (s *PubSubService) SubscribeToAllUsers(handler BroadcastHandler) error {
	return s.Subscribe("users:broadcast", handler)
}

// BroadcastToAllUsers broadcasts a message to all users (all servers)
func (s *PubSubService) BroadcastToAllUsers(userID uuid.UUID, msgType string, payload interface{}) error {
	msg := &BroadcastMessage{
		UserID:  userID,
		Type:    msgType,
		Payload: payload,
	}

	return s.Publish("users:broadcast", msg)
}

// Close closes the PubSubService and cancels all subscriptions
func (s *PubSubService) Close() error {
	log.Printf("🔕 Closing PubSubService for server %s", s.serverID[:8])
	s.cancelFunc()
	return nil
}

// ConvertToWSResponse converts a broadcast payload to WSResponse
func ConvertToWSResponse(payload interface{}) (*models.Message, error) {
	// Marshal and unmarshal to convert map to struct
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	var msg models.Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
