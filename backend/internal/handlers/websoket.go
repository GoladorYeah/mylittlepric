package handlers

import (
	"encoding/json"
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

type WSMessage struct {
	Type      string                 `json:"type"`
	SessionID string                 `json:"session_id,omitempty"`
	Message   string                 `json:"message,omitempty"`
	Country   string                 `json:"country,omitempty"`
	Language  string                 `json:"language,omitempty"`
	NewSearch bool                   `json:"new_search,omitempty"`
	PageToken string                 `json:"page_token,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

type WSResponse struct {
	Type           string                         `json:"type"`
	Output         string                         `json:"output,omitempty"`
	QuickReplies   []string                       `json:"quick_replies,omitempty"`
	Products       []models.ProductCard           `json:"products,omitempty"`
	SearchType     string                         `json:"search_type,omitempty"`
	SessionID      string                         `json:"session_id"`
	MessageCount   int                            `json:"message_count"`
	SearchState    *models.SearchStateInfo        `json:"search_state,omitempty"`
	ProductDetails *models.ProductDetailsResponse `json:"product_details,omitempty"`
	Error          string                         `json:"error,omitempty"`
}

func NewWSHandler(c *container.Container) *WSHandler {
	return &WSHandler{
		container: c,
		clients:   make(map[string]*websocket.Conn),
	}
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
		msg.Language = h.container.SessionService.GetLanguageForCountry(msg.Country)
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
			h.sendError(c, "session_error", "Failed to start new search")
			return
		}
		session, _ = h.container.SessionService.GetSession(msg.SessionID)
		fmt.Printf("✅ New search started (%d)\n", session.SearchState.SearchCount)
	} else if session.SearchState.Status == models.SearchStatusIdle {
		session.SearchState.Status = models.SearchStatusInProgress
		h.container.SessionService.UpdateSession(session)
	}

	if h.container.SessionService.IsSearchCompleted(msg.SessionID) {
		messageLower := strings.ToLower(msg.Message)
		refinementWords := []string{"cheap", "expensive", "premium", "budget", "price", "cost", "less", "more"}
		isRefinement := false
		for _, word := range refinementWords {
			if strings.Contains(messageLower, word) {
				isRefinement = true
				break
			}
		}

		if isRefinement {
			h.container.SessionService.ResetSearchStatus(msg.SessionID)
			session.SearchState.Status = models.SearchStatusInProgress
			h.container.SessionService.UpdateSession(session)
		}
	}

	userMessage := &models.Message{
		ID:        uuid.New(),
		SessionID: session.ID,
		Role:      "user",
		Content:   msg.Message,
		Category:  session.SearchState.Category,
		CreatedAt: time.Now(),
	}

	if err := h.container.SessionService.AddMessage(msg.SessionID, userMessage); err != nil {
		h.sendError(c, "message_error", "Failed to save message")
		return
	}

	if err := h.container.SessionService.IncrementMessageCount(msg.SessionID); err != nil {
		h.sendError(c, "session_error", "Failed to update session")
		return
	}

	history, err := h.container.SessionService.GetConversationHistory(msg.SessionID)
	if err != nil {
		h.sendError(c, "session_error", "Failed to get history")
		return
	}

	fmt.Printf("🤖 Processing (category: %s)\n", session.SearchState.Category)

	startTime := time.Now()
	geminiResponse, keyIndex, err := h.container.GeminiService.ProcessMessageWithContext(
		msg.Message,
		history,
		session.CountryCode,
		session.LanguageCode,
		session.SearchState.Category,
		session.SearchState.LastProduct,
	)
	responseTime := time.Since(startTime)

	h.container.GeminiRotator.RecordUsage(keyIndex, err == nil, responseTime)

	if err != nil {
		fmt.Printf("❌ Gemini error: %v\n", err)
		h.sendError(c, "ai_error", fmt.Sprintf("AI error: %v", err))
		return
	}

	if geminiResponse.Category != "" && session.SearchState.Category == "" {
		h.container.SessionService.SetCategory(msg.SessionID, geminiResponse.Category)
		session.SearchState.Category = geminiResponse.Category
		fmt.Printf("   📂 Category set: %s\n\n", geminiResponse.Category)
	}

	var response *WSResponse

	switch geminiResponse.ResponseType {
	case "dialogue":
		response = h.createDialogueResponse(msg, session, geminiResponse)

	case "search":
		response = h.createSearchResponse(msg, session, geminiResponse)
		if response.Error != "" {
			h.sendResponse(c, response)
			return
		}

	default:
		h.sendError(c, "invalid_response", "Invalid response type")
		return
	}

	assistantMessage := &models.Message{
		ID:           uuid.New(),
		SessionID:    session.ID,
		Role:         "assistant",
		Content:      geminiResponse.Output,
		ResponseType: geminiResponse.ResponseType,
		QuickReplies: geminiResponse.QuickReplies,
		Category:     session.SearchState.Category,
		CreatedAt:    time.Now(),
	}

	if err := h.container.SessionService.AddMessage(msg.SessionID, assistantMessage); err != nil {
		fmt.Printf("⚠️  Warning: Failed to save assistant message: %v\n", err)
	}

	h.sendResponse(c, response)
}

func (h *WSHandler) handleProductDetails(c *websocket.Conn, msg *WSMessage) {
	if msg.PageToken == "" {
		h.sendError(c, "validation_error", "Page token is required")
		return
	}

	country := msg.Country
	if country == "" {
		country = "CH"
	}

	cachedProduct, err := h.container.CacheService.GetProductByToken(msg.PageToken)
	if err == nil && cachedProduct != nil {
		details := h.formatProductDetails(cachedProduct)
		h.sendResponse(c, &WSResponse{
			Type:           "product_details",
			ProductDetails: details,
		})
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
		log.Printf("Warning: Failed to cache product details: %v", err)
	}

	details := h.formatProductDetails(productDetails)
	h.sendResponse(c, &WSResponse{
		Type:           "product_details",
		ProductDetails: details,
	})
}

func (h *WSHandler) createDialogueResponse(msg *WSMessage, session *models.ChatSession, geminiResponse *models.GeminiResponse) *WSResponse {
	return &WSResponse{
		Type:         "text",
		Output:       geminiResponse.Output,
		QuickReplies: geminiResponse.QuickReplies,
		SessionID:    msg.SessionID,
		MessageCount: session.MessageCount + 1,
		SearchState:  h.container.SessionService.GetSearchStateInfo(msg.SessionID),
	}
}

func (h *WSHandler) createSearchResponse(msg *WSMessage, session *models.ChatSession, geminiResponse *models.GeminiResponse) *WSResponse {
	isValid, errMsg := h.container.Optimizer.ValidateQuery(geminiResponse.SearchPhrase)
	if !isValid {
		return &WSResponse{Type: "error", Error: errMsg}
	}

	if !h.container.Optimizer.IsProductQuery(geminiResponse.SearchPhrase) {
		return &WSResponse{
			Type:         "text",
			Output:       "Please ask about a specific product.",
			SessionID:    msg.SessionID,
			MessageCount: session.MessageCount + 1,
			SearchState:  h.container.SessionService.GetSearchStateInfo(msg.SessionID),
		}
	}

	optimizedQuery := h.container.Optimizer.OptimizeQuery(geminiResponse.SearchPhrase, geminiResponse.SearchType)

	if geminiResponse.PriceFilter != "" && session.SearchState.LastProduct != nil {
		optimizedQuery = session.SearchState.LastProduct.Name
	}

	fmt.Printf("   🔍 Query: '%s' → '%s'\n", geminiResponse.SearchPhrase, optimizedQuery)
	if geminiResponse.PriceFilter != "" {
		fmt.Printf("   💰 Price filter: %s (ref: %.2f)\n", geminiResponse.PriceFilter, session.SearchState.LastProduct.Price)
	}

	startTime := time.Now()
	products, keyIndex, err := h.container.SerpService.SearchWithCache(
		optimizedQuery,
		geminiResponse.SearchType,
		session.CountryCode,
		h.container.CacheService,
	)
	responseTime := time.Since(startTime)

	if keyIndex != -1 {
		h.container.SerpRotator.RecordUsage(keyIndex, err == nil, responseTime)
		fmt.Printf("   📊 SERP API (key %d, %dms)\n", keyIndex, responseTime.Milliseconds())
	} else {
		fmt.Printf("   💾 Cache hit (%dms)\n", responseTime.Milliseconds())
	}

	if err != nil {
		return &WSResponse{Type: "error", Error: "Search failed"}
	}

	if len(products) == 0 {
		return &WSResponse{
			Type:         "text",
			Output:       "No products found. Try different description.",
			SessionID:    msg.SessionID,
			MessageCount: session.MessageCount + 1,
			SearchState:  h.container.SessionService.GetSearchStateInfo(msg.SessionID),
		}
	}

	if geminiResponse.PriceFilter != "" && session.SearchState.LastProduct != nil {
		products = h.filterByPrice(products, session.SearchState.LastProduct.Price, geminiResponse.PriceFilter, session.Currency)
	}

	if len(products) > 0 {
		firstProduct := products[0]
		price := h.extractPrice(firstProduct.Price)
		h.container.SessionService.SetLastProduct(msg.SessionID, firstProduct.Name, price)
	}

	h.container.SessionService.MarkSearchCompleted(msg.SessionID)
	fmt.Printf("   ✅ Found %d products.\n", len(products))

	output := "Here are the best options:"
	if geminiResponse.SearchType == "exact" {
		output = "Here is the product:"
	}
	if geminiResponse.PriceFilter == "cheaper" {
		output = "Here are cheaper alternatives:"
	} else if geminiResponse.PriceFilter == "expensive" {
		output = "Here are more premium options:"
	}

	return &WSResponse{
		Type:         "product_card",
		Output:       output,
		Products:     products,
		SearchType:   geminiResponse.SearchType,
		SessionID:    msg.SessionID,
		MessageCount: session.MessageCount + 1,
		SearchState:  h.container.SessionService.GetSearchStateInfo(msg.SessionID),
	}
}

func (h *WSHandler) formatProductDetails(productData map[string]interface{}) *models.ProductDetailsResponse {
	productResults, ok := productData["product_results"].(map[string]interface{})
	if !ok {
		return nil
	}

	response := &models.ProductDetailsResponse{
		Type:    "product_details",
		Title:   getStringValue(productResults, "title"),
		Price:   getStringValue(productResults, "price"),
		Rating:  float32(getFloatValue(productResults, "rating")),
		Reviews: getIntValue(productResults, "reviews"),
	}

	if aboutProduct, ok := productResults["about_the_product"].(map[string]interface{}); ok {
		response.Description = getStringValue(aboutProduct, "description")
	}

	if thumbnails, ok := productResults["thumbnails"].([]interface{}); ok {
		for _, thumb := range thumbnails {
			if thumbStr, ok := thumb.(string); ok {
				response.Images = append(response.Images, thumbStr)
			}
		}
	}

	if specs, ok := productResults["specifications"].([]interface{}); ok {
		for _, specData := range specs {
			if spec, ok := specData.(map[string]interface{}); ok {
				response.Specifications = append(response.Specifications, models.ProductSpec{
					Title: getStringValue(spec, "title"),
					Value: getStringValue(spec, "value"),
				})
			}
		}
	}

	if variants, ok := productResults["variants"].([]interface{}); ok {
		for _, variantData := range variants {
			if variant, ok := variantData.(map[string]interface{}); ok {
				productVariant := models.ProductVariant{
					Title: getStringValue(variant, "title"),
					Items: []models.VariantItem{},
				}

				if items, ok := variant["items"].([]interface{}); ok {
					for _, itemData := range items {
						if item, ok := itemData.(map[string]interface{}); ok {
							pageToken := extractPageTokenFromLink(getStringValue(item, "serpapi_link"))

							productVariant.Items = append(productVariant.Items, models.VariantItem{
								Name:      getStringValue(item, "name"),
								Selected:  getBoolValue(item, "selected"),
								Available: getBoolValue(item, "available"),
								PageToken: pageToken,
							})
						}
					}
				}

				response.Variants = append(response.Variants, productVariant)
			}
		}
	}

	if offers, ok := productResults["offers"].([]interface{}); ok {
		for _, offerData := range offers {
			if offer, ok := offerData.(map[string]interface{}); ok {
				response.Offers = append(response.Offers, models.ProductOffer{
					Merchant:     getStringValue(offer, "name"),
					Price:        getStringValue(offer, "price"),
					Currency:     extractCurrency(getStringValue(offer, "price")),
					Link:         getStringValue(offer, "link"),
					Availability: getStringValue(offer, "details_and_offers"),
					Shipping:     getStringValue(offer, "shipping"),
					Rating:       float32(getFloatValue(offer, "rating")),
				})
			}
		}
	}

	if videos, ok := productResults["videos"].([]interface{}); ok {
		for _, videoData := range videos {
			if video, ok := videoData.(map[string]interface{}); ok {
				response.Videos = append(response.Videos, models.ProductVideo{
					Title:     getStringValue(video, "title"),
					Link:      getStringValue(video, "link"),
					Source:    getStringValue(video, "source"),
					Channel:   getStringValue(video, "channel"),
					Duration:  getStringValue(video, "duration"),
					Thumbnail: getStringValue(video, "thumbnail"),
				})
			}
		}
	}

	if moreOptions, ok := productResults["more_options"].([]interface{}); ok {
		for _, optionData := range moreOptions {
			if option, ok := optionData.(map[string]interface{}); ok {
				pageToken := extractPageTokenFromLink(getStringValue(option, "serpapi_link"))

				response.MoreOptions = append(response.MoreOptions, models.AlternativeProduct{
					Title:     getStringValue(option, "title"),
					Thumbnail: getStringValue(option, "thumbnail"),
					Price:     getStringValue(option, "price"),
					Rating:    float32(getFloatValue(option, "rating")),
					Reviews:   getIntValue(option, "reviews"),
					PageToken: pageToken,
				})
			}
		}
	}

	if ratings, ok := productResults["ratings"].([]interface{}); ok {
		for _, ratingData := range ratings {
			if rating, ok := ratingData.(map[string]interface{}); ok {
				response.RatingBreakdown = append(response.RatingBreakdown, models.RatingBreakdown{
					Stars:  getIntValue(rating, "stars"),
					Amount: getIntValue(rating, "amount"),
				})
			}
		}
	}

	return response
}

func (h *WSHandler) filterByPrice(products []models.ProductCard, refPrice float64, filter string, currency string) []models.ProductCard {
	filtered := []models.ProductCard{}

	for _, product := range products {
		price := h.extractPrice(product.Price)

		if filter == "cheaper" && price < refPrice && price > refPrice*0.3 {
			filtered = append(filtered, product)
		} else if filter == "expensive" && price > refPrice && price < refPrice*3 {
			filtered = append(filtered, product)
		}
	}

	if len(filtered) == 0 {
		return products
	}

	return filtered
}

func (h *WSHandler) extractPrice(priceStr string) float64 {
	cleaned := strings.TrimSpace(priceStr)
	cleaned = strings.ReplaceAll(cleaned, "CHF", "")
	cleaned = strings.ReplaceAll(cleaned, "EUR", "")
	cleaned = strings.ReplaceAll(cleaned, "USD", "")
	cleaned = strings.ReplaceAll(cleaned, "GBP", "")
	cleaned = strings.ReplaceAll(cleaned, "$", "")
	cleaned = strings.ReplaceAll(cleaned, "€", "")
	cleaned = strings.ReplaceAll(cleaned, "£", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.TrimSpace(cleaned)

	price, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0
	}

	return price
}

func (h *WSHandler) sendResponse(c *websocket.Conn, response *WSResponse) {
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("❌ Failed to marshal response: %v", err)
		return
	}

	if err := c.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("❌ Failed to send response: %v", err)
	}
}

func (h *WSHandler) sendError(c *websocket.Conn, code, message string) {
	h.sendResponse(c, &WSResponse{
		Type:  "error",
		Error: fmt.Sprintf("%s: %s", code, message),
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
