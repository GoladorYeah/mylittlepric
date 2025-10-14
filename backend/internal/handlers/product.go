package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"mylittleprice/internal/container"
	"mylittleprice/internal/models"
)

type ProductHandler struct {
	container *container.Container
}

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

	if req.PageToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: "Page token is required",
		})
	}

	if req.Country == "" {
		req.Country = "CH"
	}

	cachedProduct, err := h.container.CacheService.GetProductByToken(req.PageToken)
	if err == nil && cachedProduct != nil {
		return h.formatProductResponse(c, cachedProduct)
	}

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
		for _, spec := range specs {
			if specMap, ok := spec.(map[string]interface{}); ok {
				response.Specifications = append(response.Specifications, models.Specification{
					Title: getStringValue(specMap, "title"),
					Value: getStringValue(specMap, "value"),
				})
			}
		}
	}

	if variants, ok := productResults["variants"].([]interface{}); ok {
		for _, variant := range variants {
			if variantMap, ok := variant.(map[string]interface{}); ok {
				items := []interface{}{}
				if variantItems, ok := variantMap["items"].([]interface{}); ok {
					items = variantItems
				}

				response.Variants = append(response.Variants, models.Variant{
					Title: getStringValue(variantMap, "title"),
					Items: items,
				})
			}
		}
	}

	if sellers, ok := productResults["sellers"].([]interface{}); ok {
		for _, seller := range sellers {
			if sellerMap, ok := seller.(map[string]interface{}); ok {
				offer := models.Offer{
					Merchant:     getStringValue(sellerMap, "name"),
					Price:        getStringValue(sellerMap, "price"),
					Currency:     getStringValue(sellerMap, "currency"),
					Link:         getStringValue(sellerMap, "link"),
					Availability: getStringValue(sellerMap, "availability"),
					Shipping:     getStringValue(sellerMap, "shipping"),
					Rating:       float32(getFloatValue(sellerMap, "rating")),
				}
				response.Offers = append(response.Offers, offer)
			}
		}
	}

	if videos, ok := productResults["videos"].([]interface{}); ok {
		response.Videos = videos
	}

	if moreOptions, ok := productResults["more_options"].([]interface{}); ok {
		response.MoreOptions = moreOptions
	}

	if ratingBreakdown, ok := productResults["rating_breakdown"].([]interface{}); ok {
		for _, item := range ratingBreakdown {
			if itemMap, ok := item.(map[string]interface{}); ok {
				response.RatingBreakdown = append(response.RatingBreakdown, models.RatingBreakdownItem{
					Stars:  getIntValue(itemMap, "stars"),
					Amount: getIntValue(itemMap, "amount"),
				})
			}
		}
	}

	return c.JSON(response)
}
