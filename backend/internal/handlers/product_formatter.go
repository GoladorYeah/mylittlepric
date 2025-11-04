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

	// Try "stores" field first (used by google_immersive_product with more_stores=true)
	if stores, ok := productResults["stores"].([]interface{}); ok {
		for _, store := range stores {
			if storeMap, ok := store.(map[string]interface{}); ok {
				offer := models.Offer{
					Merchant:          getStringValue(storeMap, "name"),
					Logo:              getStringValue(storeMap, "logo"),
					Price:             getStringValue(storeMap, "price"),
					ExtractedPrice:    getFloatValue(storeMap, "extracted_price"),
					Currency:          getStringValue(storeMap, "currency"),
					Link:              getStringValue(storeMap, "link"),
					Title:             getStringValue(storeMap, "title"),
					Availability:      getStringValue(storeMap, "availability"),
					Shipping:          getStringValue(storeMap, "shipping"),
					ShippingExtracted: getFloatValue(storeMap, "shipping_extracted"),
					Total:             getStringValue(storeMap, "total"),
					ExtractedTotal:    getFloatValue(storeMap, "extracted_total"),
					Rating:            float32(getFloatValue(storeMap, "rating")),
					Reviews:           getIntValue(storeMap, "reviews"),
					PaymentMethods:    getStringValue(storeMap, "payment_methods"),
					Tag:               getStringValue(storeMap, "tag"),
					MonthlyPaymentDur: getIntValue(storeMap, "monthly_payment_duration"),
					DownPayment:       getStringValue(storeMap, "down_payment"),
				}

				// Parse details_and_offers array
				if details, ok := storeMap["details_and_offers"].([]interface{}); ok {
					for _, detail := range details {
						if detailStr, ok := detail.(string); ok {
							offer.DetailsAndOffers = append(offer.DetailsAndOffers, detailStr)
						}
					}
				}

				response.Offers = append(response.Offers, offer)
			}
		}
	} else if sellers, ok := productResults["sellers"].([]interface{}); ok {
		// Fallback to "sellers" field for compatibility
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
