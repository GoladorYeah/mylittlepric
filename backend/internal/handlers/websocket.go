package handlers

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"

	"mylittleprice/internal/container"
	"mylittleprice/internal/models"
)

type Client struct {
	Conn   *websocket.Conn
	UserID *uuid.UUID // nil for anonymous users
}

type WSHandler struct {
	container *container.Container
	processor *ChatProcessor
	clients   map[string]*Client // clientID -> Client
	userConns map[uuid.UUID]map[string]bool // userID -> set of clientIDs
	mu        sync.RWMutex
}

func NewWSHandler(c *container.Container) *WSHandler {
	return &WSHandler{
		container: c,
		processor: NewChatProcessor(c),
		clients:   make(map[string]*Client),
		userConns: make(map[uuid.UUID]map[string]bool),
	}
}

type WSMessage struct {
	Type            string                 `json:"type"`
	SessionID       string                 `json:"session_id"`
	Message         string                 `json:"message"`
	Country         string                 `json:"country"`
	Language        string `json:"language"`
	Currency        string                 `json:"currency"`
	NewSearch       bool                   `json:"new_search"`
	PageToken       string                 `json:"page_token"`
	CurrentCategory string                 `json:"current_category"`
	AccessToken     string                 `json:"access_token,omitempty"` // Optional JWT token for authentication
	Preferences     map[string]interface{} `json:"preferences,omitempty"`  // For preferences sync
	SavedSearch     map[string]interface{} `json:"saved_search,omitempty"` // For saved search sync
}

type WSResponse struct {
	Type           string                         `json:"type"`
	Output         string                         `json:"output,omitempty"`
	QuickReplies   []string                       `json:"quick_replies,omitempty"`
	Products       []models.ProductCard           `json:"products,omitempty"`
	SearchType     string                         `json:"search_type,omitempty"`
	SessionID      string                         `json:"session_id"`
	MessageCount   int                            `json:"message_count,omitempty"`
	SearchState    *models.SearchStateResponse    `json:"search_state,omitempty"`
	ProductDetails *models.ProductDetailsResponse `json:"product_details,omitempty"`
	Error          string                         `json:"error,omitempty"`
	Message        string                         `json:"message,omitempty"`
}

func (h *WSHandler) HandleWebSocket(c *websocket.Conn) {
	clientID := uuid.New().String()
	var userID *uuid.UUID

	log.Printf("🔌 Client connected: %s", clientID)

	// First message should contain access_token if user is authenticated
	// We'll update userID as messages come in with access_token
	client := &Client{
		Conn:   c,
		UserID: nil,
	}
	h.addClient(clientID, client)
	defer h.removeClient(clientID)

	for {
		var msg WSMessage
		err := c.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("❌ WebSocket error: %v", err)
			}
			break
		}

		// Update userID if access_token is provided
		if msg.AccessToken != "" {
			claims, err := h.container.JWTService.ValidateAccessToken(msg.AccessToken)
			if err == nil {
				if userID == nil || *userID != claims.UserID {
					// First time or changed user - update mapping
					h.updateClientUser(clientID, &claims.UserID)
					userID = &claims.UserID
					client.UserID = userID
					log.Printf("🔐 Client %s authenticated as user %s", clientID, userID.String())
				}
			}
		}

		h.handleMessage(c, &msg, clientID)
	}

	log.Printf("🔌 Client disconnected: %s", clientID)
}

func (h *WSHandler) handleMessage(c *websocket.Conn, msg *WSMessage, clientID string) {
	switch msg.Type {
	case "chat":
		h.handleChat(c, msg, clientID)
	case "product_details":
		h.handleProductDetails(c, msg)
	case "ping":
		h.sendResponse(c, &WSResponse{Type: "pong"})
	case "sync_preferences":
		h.handleSyncPreferences(c, msg, clientID)
	case "sync_saved_search":
		h.handleSyncSavedSearch(c, msg, clientID)
	case "sync_session":
		h.handleSyncSession(c, msg, clientID)
	default:
		h.sendError(c, "unknown_message_type", "Unknown message type")
	}
}

func (h *WSHandler) handleChat(c *websocket.Conn, msg *WSMessage, clientID string) {
	// Extract user ID from access token if provided
	var userID *uuid.UUID
	if msg.AccessToken != "" {
		claims, err := h.container.JWTService.ValidateAccessToken(msg.AccessToken)
		if err == nil {
			userID = &claims.UserID
		}
	}

	// Broadcast user message to other devices BEFORE processing
	if userID != nil {
		userMsgSync := &WSResponse{
			Type:      "user_message_sync",
			Output:    msg.Message,
			SessionID: msg.SessionID,
		}
		h.broadcastToUser(*userID, userMsgSync, clientID)
	}

	// Process chat using shared processor
	processorReq := &ChatRequest{
		SessionID:       msg.SessionID,
		UserID:          userID,
		Message:         msg.Message,
		Country:         msg.Country,
		Language:        msg.Language,
		Currency:        msg.Currency,
		NewSearch:       msg.NewSearch,
		CurrentCategory: msg.CurrentCategory,
	}

	result := h.processor.ProcessChat(processorReq)

	// Handle errors
	if result.Error != nil {
		h.sendError(c, result.Error.Code, result.Error.Message)
		return
	}

	// Build response
	response := &WSResponse{
		Type:         result.Type,
		Output:       result.Output,
		QuickReplies: result.QuickReplies,
		Products:     result.Products,
		SearchType:   result.SearchType,
		SessionID:    result.SessionID,
		MessageCount: result.MessageCount,
		SearchState:  result.SearchState,
	}

	// Send response to the sender
	h.sendResponse(c, response)

	// Broadcast assistant message to other devices of the same user
	if userID != nil {
		// Create sync message for other devices
		syncMsg := &WSResponse{
			Type:         "assistant_message_sync",
			Output:       result.Output,
			QuickReplies: result.QuickReplies,
			Products:     result.Products,
			SearchType:   result.SearchType,
			SessionID:    result.SessionID,
			MessageCount: result.MessageCount,
			SearchState:  result.SearchState,
		}
		h.broadcastToUser(*userID, syncMsg, clientID)
	}
}

func (h *WSHandler) handleProductDetails(c *websocket.Conn, msg *WSMessage) {
	if msg.PageToken == "" {
		h.sendError(c, "validation_error", "Page token is required")
		return
	}

	if msg.Country == "" {
		msg.Country = h.container.Config.DefaultCountry
	}

	cachedProduct, err := h.container.CacheService.GetProductByToken(msg.PageToken)
	if err == nil && cachedProduct != nil {
		h.sendProductDetailsResponse(c, cachedProduct, msg.SessionID)
		return
	}

	startTime := time.Now()
	productDetails, keyIndex, err := h.container.SerpService.GetProductDetailsByToken(msg.PageToken)
	responseTime := time.Since(startTime)

	h.container.SerpRotator.RecordUsage(keyIndex, err == nil, responseTime)

	if err != nil {
		h.sendError(c, "fetch_error", "Failed to fetch product details")
		return
	}

	if err := h.container.CacheService.SetProductByToken(msg.PageToken, productDetails, h.container.Config.CacheImmersiveTTL); err != nil {
		fmt.Printf("⚠️ Failed to cache product details: %v\n", err)
	}

	h.sendProductDetailsResponse(c, productDetails, msg.SessionID)
}

func (h *WSHandler) sendProductDetailsResponse(c *websocket.Conn, productData map[string]interface{}, sessionID string) {
	details, err := FormatProductDetails(productData)
	if err != nil {
		h.sendError(c, "parse_error", err.Error())
		return
	}

	h.sendResponse(c, &WSResponse{
		Type:           "product_details",
		ProductDetails: details,
		SessionID:      sessionID,
	})
}

func (h *WSHandler) addClient(id string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[id] = client
}

func (h *WSHandler) removeClient(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	client, exists := h.clients[id]
	if !exists {
		return
	}

	// Remove from userConns if user was authenticated
	if client.UserID != nil {
		if connSet, ok := h.userConns[*client.UserID]; ok {
			delete(connSet, id)
			// Remove user entry if no more connections
			if len(connSet) == 0 {
				delete(h.userConns, *client.UserID)
			}
		}
	}

	delete(h.clients, id)
}

func (h *WSHandler) updateClientUser(clientID string, userID *uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	client, exists := h.clients[clientID]
	if !exists {
		return
	}

	// Remove from old user's connection set
	if client.UserID != nil {
		if connSet, ok := h.userConns[*client.UserID]; ok {
			delete(connSet, clientID)
			if len(connSet) == 0 {
				delete(h.userConns, *client.UserID)
			}
		}
	}

	// Add to new user's connection set
	if userID != nil {
		if _, ok := h.userConns[*userID]; !ok {
			h.userConns[*userID] = make(map[string]bool)
		}
		h.userConns[*userID][clientID] = true
	}

	client.UserID = userID
}

// broadcastToUser sends a message to all connections of a user except the sender
func (h *WSHandler) broadcastToUser(userID uuid.UUID, response *WSResponse, excludeClientID string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clientIDs, ok := h.userConns[userID]
	if !ok {
		return
	}

	for cid := range clientIDs {
		if cid == excludeClientID {
			continue
		}

		client, exists := h.clients[cid]
		if !exists {
			continue
		}

		if err := client.Conn.WriteJSON(response); err != nil {
			log.Printf("❌ Failed to broadcast to client %s: %v", cid, err)
		}
	}
}

// handleSyncPreferences handles preference synchronization across devices
func (h *WSHandler) handleSyncPreferences(c *websocket.Conn, msg *WSMessage, clientID string) {
	// Extract user ID from access token
	if msg.AccessToken == "" {
		h.sendError(c, "auth_required", "Authentication required for preferences sync")
		return
	}

	claims, err := h.container.JWTService.ValidateAccessToken(msg.AccessToken)
	if err != nil {
		h.sendError(c, "invalid_token", "Invalid access token")
		return
	}

	// Broadcast preferences update to other devices
	syncMsg := &WSResponse{
		Type:      "preferences_updated",
		SessionID: msg.SessionID,
		Message:   "Preferences updated",
	}

	h.sendResponse(c, &WSResponse{Type: "sync_ack"})
	h.broadcastToUser(claims.UserID, syncMsg, clientID)
}

// handleSyncSavedSearch handles saved search synchronization across devices
func (h *WSHandler) handleSyncSavedSearch(c *websocket.Conn, msg *WSMessage, clientID string) {
	// Extract user ID from access token
	if msg.AccessToken == "" {
		// Anonymous users can't sync across devices
		return
	}

	claims, err := h.container.JWTService.ValidateAccessToken(msg.AccessToken)
	if err != nil {
		return
	}

	// Broadcast saved search update to other devices
	syncMsg := &WSResponse{
		Type:      "saved_search_updated",
		SessionID: msg.SessionID,
	}

	h.sendResponse(c, &WSResponse{Type: "sync_ack"})
	h.broadcastToUser(claims.UserID, syncMsg, clientID)
}

// handleSyncSession handles session change synchronization across devices
func (h *WSHandler) handleSyncSession(c *websocket.Conn, msg *WSMessage, clientID string) {
	// Extract user ID from access token
	if msg.AccessToken == "" {
		return
	}

	claims, err := h.container.JWTService.ValidateAccessToken(msg.AccessToken)
	if err != nil {
		return
	}

	// Broadcast session change to other devices
	syncMsg := &WSResponse{
		Type:      "session_changed",
		SessionID: msg.SessionID,
	}

	h.sendResponse(c, &WSResponse{Type: "sync_ack"})
	h.broadcastToUser(claims.UserID, syncMsg, clientID)
}

func (h *WSHandler) sendResponse(c *websocket.Conn, response *WSResponse) {
	if err := c.WriteJSON(response); err != nil {
		log.Printf("❌ Failed to send response: %v", err)
	}
}

func (h *WSHandler) sendError(c *websocket.Conn, errorCode, message string) {
	h.sendResponse(c, &WSResponse{
		Type:    "error",
		Error:   errorCode,
		Message: message,
	})
}
