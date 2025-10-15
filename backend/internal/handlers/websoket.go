// backend/internal/handlers/websoket.go
package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"

	"mylittleprice/internal/container"
	"mylittleprice/internal/models"
)

type WSHandler struct {
	container *container.Container
	clients   map[string]*websocket.Conn
	mu        sync.RWMutex
}

func NewWSHandler(c *container.Container) *WSHandler {
	return &WSHandler{
		container: c,
		clients:   make(map[string]*websocket.Conn),
	}
}

type WSMessage struct {
	Type      string `json:"type"`
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
	Country   string `json:"country"`
	Language  string `json:"language"`
	Currency  string `json:"currency"`
	NewSearch bool   `json:"new_search"`
	PageToken string `json:"page_token"`
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
	if msg.Message == "" {
		h.sendError(c, "validation_error", "Message is required")
		return
	}

	if msg.Country == "" {
		msg.Country = "CH"
	}

	if msg.Language == "" {
		msg.Language = "en"
	}

	var session *models.ChatSession
	var err error

	if msg.SessionID != "" {
		session, err = h.container.SessionService.GetSession(msg.SessionID)
		if err != nil {
			msg.SessionID = uuid.New().String()
			session, err = h.container.SessionService.CreateSession(msg.SessionID, msg.Country, msg.Language)
			if err != nil {
				h.sendError(c, "session_error", "Failed to create session")
				return
			}
		}
	} else {
		msg.SessionID = uuid.New().String()
		session, err = h.container.SessionService.CreateSession(msg.SessionID, msg.Country, msg.Language)
		if err != nil {
			h.sendError(c, "session_error", "Failed to create session")
			return
		}
	}

	if msg.NewSearch {
		fmt.Printf("🔄 New search for session %s\n", msg.SessionID)
		if err := h.container.SessionService.StartNewSearch(msg.SessionID); err != nil {
			fmt.Printf("⚠️ Failed to start new search: %v\n", err)
		}
		session, _ = h.container.SessionService.GetSession(msg.SessionID)
	}

	if session.SearchState.SearchCount >= h.container.SessionService.GetMaxSearches() {
		h.sendResponse(c, &WSResponse{
			Type:         "text",
			Output:       "You have reached the maximum number of searches. Please start a new search.",
			SessionID:    msg.SessionID,
			MessageCount: session.MessageCount,
			SearchState: &models.SearchStateResponse{
				Status:      string(session.SearchState.Status),
				CanContinue: false,
				SearchCount: session.SearchState.SearchCount,
				MaxSearches: h.container.SessionService.GetMaxSearches(),
				Message:     "Search limit reached",
			},
		})
		return
	}

	userMessage := &models.Message{
		ID:        uuid.New(),
		SessionID: session.ID,
		Role:      "user",
		Content:   msg.Message,
		CreatedAt: time.Now(),
	}

	if err := h.container.SessionService.AddMessage(msg.SessionID, userMessage); err != nil {
		h.sendError(c, "storage_error", "Failed to store message")
		return
	}

	if err := h.container.SessionService.IncrementMessageCount(msg.SessionID); err != nil {
		fmt.Printf("⚠️ Failed to increment message count: %v\n", err)
	}

	conversationHistory, err := h.container.SessionService.GetConversationHistory(msg.SessionID)
	if err != nil {
		conversationHistory = []map[string]string{}
	}

	geminiResponse, _, err := h.container.GeminiService.ProcessMessageWithContext(
		msg.Message,
		conversationHistory,
		msg.Country,
		msg.Language,
		session.SearchState.Category,
		session.SearchState.LastProduct,
	)

	if err != nil {
		log.Printf("❌ Gemini processing error: %v", err)
		h.sendError(c, "processing_error", fmt.Sprintf("AI processing failed: %v", err))
		return
	}

	if geminiResponse == nil {
		log.Printf("❌ Gemini returned nil response")
		h.sendError(c, "processing_error", "AI returned empty response")
		return
	}

	session.SearchState.Category = geminiResponse.Category

	assistantMessage := &models.Message{
		ID:           uuid.New(),
		SessionID:    session.ID,
		Role:         "assistant",
		Content:      geminiResponse.Output,
		ResponseType: geminiResponse.ResponseType,
		QuickReplies: geminiResponse.QuickReplies,
		CreatedAt:    time.Now(),
	}

	if err := h.container.SessionService.AddMessage(msg.SessionID, assistantMessage); err != nil {
		fmt.Printf("⚠️ Failed to store assistant message: %v\n", err)
	}

	response := &WSResponse{
		Type:         geminiResponse.ResponseType,
		Output:       geminiResponse.Output,
		QuickReplies: geminiResponse.QuickReplies,
		SessionID:    msg.SessionID,
		MessageCount: session.MessageCount + 1,
	}

	if geminiResponse.ResponseType == "search" {
		products, searchErr := h.performSearch(geminiResponse, msg.Country, msg.Language)
		if searchErr != nil {
			log.Printf("⚠️ Search failed: %v", searchErr)
			response.Output = "Sorry, I couldn't find any products. Please try different keywords."
			response.Type = "text"
		} else if len(products) > 0 {
			response.Products = products
			response.SearchType = geminiResponse.SearchType

			if len(products) > 0 {
				priceStr := products[0].Price
				priceStr = strings.ReplaceAll(priceStr, "$", "")
				priceStr = strings.ReplaceAll(priceStr, "€", "")
				priceStr = strings.ReplaceAll(priceStr, "£", "")
				priceStr = strings.ReplaceAll(priceStr, "CHF", "")
				priceStr = strings.TrimSpace(priceStr)
				priceStr = strings.ReplaceAll(priceStr, ",", "")

				price, _ := strconv.ParseFloat(priceStr, 64)

				session.SearchState.LastProduct = &models.ProductInfo{
					Name:  products[0].Name,
					Price: price,
				}
			}

			session.SearchState.SearchCount++
			assistantMessage.Products = products
		}
	}

	session.SearchState.Status = models.SearchStatusIdle
	if err := h.container.SessionService.UpdateSession(session); err != nil {
		fmt.Printf("⚠️ Failed to update session: %v\n", err)
	}

	response.SearchState = &models.SearchStateResponse{
		Status:      string(session.SearchState.Status),
		Category:    session.SearchState.Category,
		CanContinue: session.SearchState.SearchCount < h.container.SessionService.GetMaxSearches(),
		SearchCount: session.SearchState.SearchCount,
		MaxSearches: h.container.SessionService.GetMaxSearches(),
	}

	h.sendResponse(c, response)
}

func (h *WSHandler) handleProductDetails(c *websocket.Conn, msg *WSMessage) {
	if msg.PageToken == "" {
		h.sendError(c, "validation_error", "Page token is required")
		return
	}

	if msg.Country == "" {
		msg.Country = "CH"
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
	productResults, ok := productData["product_results"].(map[string]interface{})
	if !ok {
		h.sendError(c, "parse_error", "Invalid product data structure")
		return
	}

	details := &models.ProductDetailsResponse{
		Type:    "product_details",
		Title:   getStringValue(productResults, "title"),
		Price:   getStringValue(productResults, "price"),
		Rating:  float32(getFloatValue(productResults, "rating")),
		Reviews: getIntValue(productResults, "reviews"),
	}

	if aboutProduct, ok := productResults["about_the_product"].(map[string]interface{}); ok {
		details.Description = getStringValue(aboutProduct, "description")
	}

	if thumbnails, ok := productResults["thumbnails"].([]interface{}); ok {
		for _, thumb := range thumbnails {
			if thumbStr, ok := thumb.(string); ok {
				details.Images = append(details.Images, thumbStr)
			}
		}
	}

	if specs, ok := productResults["specifications"].([]interface{}); ok {
		for _, spec := range specs {
			if specMap, ok := spec.(map[string]interface{}); ok {
				details.Specifications = append(details.Specifications, models.Specification{
					Title: getStringValue(specMap, "title"),
					Value: getStringValue(specMap, "value"),
				})
			}
		}
	}

	if sellers, ok := productResults["sellers"].([]interface{}); ok {
		for _, seller := range sellers {
			if sellerMap, ok := seller.(map[string]interface{}); ok {
				offer := models.Offer{
					Merchant: getStringValue(sellerMap, "name"),
					Price:    getStringValue(sellerMap, "price"),
					Link:     getStringValue(sellerMap, "link"),
					Rating:   float32(getFloatValue(sellerMap, "rating")),
				}
				details.Offers = append(details.Offers, offer)
			}
		}
	}

	h.sendResponse(c, &WSResponse{
		Type:           "product_details",
		ProductDetails: details,
		SessionID:      sessionID,
	})
}

func (h *WSHandler) performSearch(geminiResp *models.GeminiResponse, country, language string) ([]models.ProductCard, error) {
	products, _, err := h.container.SerpService.SearchProducts(
		geminiResp.SearchPhrase,
		geminiResp.SearchType,
		country,
	)

	if err != nil {
		return nil, err
	}

	return products, nil
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
