package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"mylittleprice/internal/container"
	"mylittleprice/internal/models"
)

// ProductHandler handles product-related requests
type ProductHandler struct {
	container *container.Container
}

// NewProductHandler creates a new product handler with dependency injection
func NewProductHandler(c *container.Container) *ProductHandler {
	return &ProductHandler{
		container: c,
	}
}

func (h *ProductHandler) HandleProductDetails(c *fiber.Ctx) error {
	var req models.ProductDetailsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Failed to parse request body",
		})
	}

	// Validate request
	if req.PageToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: "Page token is required",
		})
	}

	if req.Country == "" {
		req.Country = "CH"
	}

	// Check cache first
	cachedProduct, err := h.container.CacheService.GetProductByToken(req.PageToken)
	if err == nil && cachedProduct != nil {
		return h.formatProductResponse(c, cachedProduct)
	}

	// Fetch from SERP API
	startTime := time.Now()
	productDetails, keyIndex, err := h.container.SerpService.GetProductDetailsByToken(req.PageToken)
	responseTime := time.Since(startTime)

	h.container.SerpRotator.RecordUsage(keyIndex, err == nil, responseTime)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "fetch_error",
			Message: "Failed to fetch product details",
		})
	}

	// Cache the result
	if err := h.container.CacheService.SetProductByToken(req.PageToken, productDetails, h.container.Config.CacheImmersiveTTL); err != nil {
		c.Context().Logger().Printf("Warning: Failed to cache product details: %v", err)
	}

	return h.formatProductResponse(c, productDetails)
}

func (h *ProductHandler) formatProductResponse(c *fiber.Ctx, productData map[string]interface{}) error {
	productResults, ok := productData["product_results"].(map[string]interface{})
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "parse_error",
			Message: "Invalid product data structure",
		})
	}

	response := models.ProductDetailsResponse{
		Type:    "product_details",
		Title:   getStringValue(productResults, "title"),
		Price:   getStringValue(productResults, "price"),
		Rating:  float32(getFloatValue(productResults, "rating")),
		Reviews: getIntValue(productResults, "reviews"),
	}

	// Extract description
	if aboutProduct, ok := productResults["about_the_product"].(map[string]interface{}); ok {
		response.Description = getStringValue(aboutProduct, "description")
	}

	// Extract thumbnails
	if thumbnails, ok := productResults["thumbnails"].([]interface{}); ok {
		for _, thumb := range thumbnails {
			if thumbStr, ok := thumb.(string); ok {
				response.Images = append(response.Images, thumbStr)
			}
		}
	}

	// Extract specifications
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

	// Extract variants
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

	// Extract offers
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

	// Extract videos
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

	// Extract more options
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

	// Extract rating breakdown
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

	return c.JSON(response)
}

// Helper functions

func getStringValue(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func getFloatValue(data map[string]interface{}, key string) float64 {
	switch v := data[key].(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return 0
	}
}

func getIntValue(data map[string]interface{}, key string) int {
	if val, ok := data[key].(float64); ok {
		return int(val)
	}
	if val, ok := data[key].(int); ok {
		return val
	}
	return 0
}

func getBoolValue(data map[string]interface{}, key string) bool {
	if val, ok := data[key].(bool); ok {
		return val
	}
	return false
}

func extractCurrency(price string) string {
	if len(price) >= 3 {
		return price[:3]
	}
	return "CHF"
}

func extractPageTokenFromLink(serpAPILink string) string {
	if serpAPILink == "" {
		return ""
	}

	tokenStart := strings.Index(serpAPILink, "page_token=")
	if tokenStart == -1 {
		return ""
	}

	tokenStart += len("page_token=")
	tokenEnd := strings.Index(serpAPILink[tokenStart:], "&")

	if tokenEnd == -1 {
		return serpAPILink[tokenStart:]
	}

	return serpAPILink[tokenStart : tokenStart+tokenEnd]
}
