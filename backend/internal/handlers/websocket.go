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

type WSHandler struct {
	container *container.Container
	processor *ChatProcessor
	clients   map[string]*websocket.Conn
	mu        sync.RWMutex
}

func NewWSHandler(c *container.Container) *WSHandler {
	return &WSHandler{
		container: c,
		processor: NewChatProcessor(c),
		clients:   make(map[string]*websocket.Conn),
	}
}

type WSMessage struct {
	Type            string `json:"type"`
	SessionID       string `json:"session_id"`
	Message         string `json:"message"`
	Country         string `json:"country"`
	Language        string `json:"language"`
	Currency        string `json:"currency"`
	NewSearch       bool   `json:"new_search"`
	PageToken       string `json:"page_token"`
	CurrentCategory string `json:"current_category"`
	AccessToken     string `json:"access_token,omitempty"` // Optional JWT token for authentication
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
	h.addClient(clientID, c)
	defer h.removeClient(clientID)

	log.Printf("🔌 Client connected: %s", clientID)

	for {
		var msg WSMessage
		err := c.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("❌ WebSocket error: %v", err)
			}
			break
		}

		h.handleMessage(c, &msg)
	}

	log.Printf("🔌 Client disconnected: %s", clientID)
}

func (h *WSHandler) handleMessage(c *websocket.Conn, msg *WSMessage) {
	switch msg.Type {
	case "chat":
		h.handleChat(c, msg)
	case "product_details":
		h.handleProductDetails(c, msg)
	case "ping":
		h.sendResponse(c, &WSResponse{Type: "pong"})
	default:
		h.sendError(c, "unknown_message_type", "Unknown message type")
	}
}

func (h *WSHandler) handleChat(c *websocket.Conn, msg *WSMessage) {
	// Extract user ID from access token if provided
	var userID *uuid.UUID
	if msg.AccessToken != "" {
		claims, err := h.container.JWTService.ValidateAccessToken(msg.AccessToken)
		if err == nil {
			userID = &claims.UserID
		}
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

	h.sendResponse(c, response)
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

func (h *WSHandler) addClient(id string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[id] = conn
}

func (h *WSHandler) removeClient(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, id)
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
