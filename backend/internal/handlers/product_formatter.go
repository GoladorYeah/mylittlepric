package handlers

import (
	"fmt"

	"mylittleprice/internal/models"
)

// FormatProductDetails extracts and formats product details from SerpAPI response
func FormatProductDetails(productData map[string]interface{}) (*models.ProductDetailsResponse, error) {
	productResults, ok := productData["product_results"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid product data structure")
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

	return response, nil
}
