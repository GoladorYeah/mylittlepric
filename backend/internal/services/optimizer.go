package services

import (
	"regexp"
	"strings"
)

// QueryOptimizer optimizes search queries for Google Shopping
type QueryOptimizer struct {
	// Patterns to remove from queries
	storagePattern   *regexp.Regexp
	colorPattern     *regexp.Regexp
	sizePattern      *regexp.Regexp
	conditionPattern *regexp.Regexp
	carrierPattern   *regexp.Regexp
	prefixPattern    *regexp.Regexp
}

// NewQueryOptimizer creates a new query optimizer
func NewQueryOptimizer() *QueryOptimizer {
	return &QueryOptimizer{
		// Remove storage specifications (256GB, 512 GB, 1TB, etc.)
		storagePattern: regexp.MustCompile(`(?i)\b\d+\s*(gb|tb|mb|gigabyte|terabyte)\b`),

		// Remove color specifications
		colorPattern: regexp.MustCompile(`(?i)\b(black|white|silver|gold|rose\s*gold|space\s*gray|grey|blue|red|green|yellow|purple|pink|orange|midnight|starlight|cosmic|deep|titanium|natural)\b`),

		// Remove size specifications (for clothes, screens, etc.)
		sizePattern: regexp.MustCompile(`(?i)\b(small|medium|large|x+l|[0-9]+(\.[0-9]+)?\s*(inch|cm|mm|meter|foot|feet)|size\s*[0-9]+)\b`),

		// Remove condition specifications
		conditionPattern: regexp.MustCompile(`(?i)\b(new|refurbished|used|like\s*new|open\s*box|factory\s*sealed|unlocked)\b`),

		// Remove carrier/network specifications
		carrierPattern: regexp.MustCompile(`(?i)\b(unlocked|verizon|at&t|t-mobile|sprint|gsm|cdma|dual\s*sim)\b`),

		// Remove search type prefixes (*, **, ***)
		prefixPattern: regexp.MustCompile(`^\*+\s*`),
	}
}

// OptimizeQuery optimizes a search query based on search type
func (o *QueryOptimizer) OptimizeQuery(query, searchType string) string {
	// Remove prefix markers
	optimized := o.prefixPattern.ReplaceAllString(query, "")

	// Trim whitespace
	optimized = strings.TrimSpace(optimized)

	// Apply different optimization levels based on search type
	switch searchType {
	case "exact":
		// For exact search, remove all parameters
		optimized = o.removeAllParameters(optimized)

	case "parameters":
		// For parameter search, keep brand but remove specific specs
		optimized = o.removeSpecifications(optimized)

	case "category":
		// For category search, keep it broad
		optimized = o.removeAllParameters(optimized)
		optimized = o.broadenQuery(optimized)

	default:
		// Default: moderate optimization
		optimized = o.removeSpecifications(optimized)
	}

	// Clean up extra spaces
	optimized = o.cleanSpaces(optimized)

	return optimized
}

// removeAllParameters removes all specific parameters from query
func (o *QueryOptimizer) removeAllParameters(query string) string {
	result := query

	// Remove storage
	result = o.storagePattern.ReplaceAllString(result, "")

	// Remove colors
	result = o.colorPattern.ReplaceAllString(result, "")

	// Remove sizes
	result = o.sizePattern.ReplaceAllString(result, "")

	// Remove conditions
	result = o.conditionPattern.ReplaceAllString(result, "")

	// Remove carrier info
	result = o.carrierPattern.ReplaceAllString(result, "")

	return result
}

// removeSpecifications removes technical specifications but keeps brand/model
func (o *QueryOptimizer) removeSpecifications(query string) string {
	result := query

	// Remove storage (most important for relevance)
	result = o.storagePattern.ReplaceAllString(result, "")

	// Remove colors (can affect results significantly)
	result = o.colorPattern.ReplaceAllString(result, "")

	// Remove carrier info
	result = o.carrierPattern.ReplaceAllString(result, "")

	return result
}

// broadenQuery makes query more general for category search
func (o *QueryOptimizer) broadenQuery(query string) string {
	// Remove model numbers that are too specific
	// Example: "iPhone 15 Pro Max" -> "iPhone"

	query = strings.TrimSpace(query)
	words := strings.Fields(query)

	// If query has more than 2 words, try to simplify
	if len(words) > 2 {
		// Keep first 1-2 words (usually brand and product type)
		return strings.Join(words[:min(2, len(words))], " ")
	}

	return query
}

// cleanSpaces removes extra whitespace
func (o *QueryOptimizer) cleanSpaces(query string) string {
	// Replace multiple spaces with single space
	spacePattern := regexp.MustCompile(`\s+`)
	result := spacePattern.ReplaceAllString(query, " ")

	return strings.TrimSpace(result)
}

// ValidateQuery checks if query is valid for searching
func (o *QueryOptimizer) ValidateQuery(query string) (bool, string) {
	query = strings.TrimSpace(query)

	// Check minimum length
	if len(query) < 2 {
		return false, "Query too short"
	}

	// Check maximum length
	if len(query) > 200 {
		return false, "Query too long"
	}

	// Check if query contains only special characters
	alphanumericPattern := regexp.MustCompile(`[a-zA-Z0-9]`)
	if !alphanumericPattern.MatchString(query) {
		return false, "Query must contain letters or numbers"
	}

	return true, ""
}

// IsProductQuery checks if query is related to product search
func (o *QueryOptimizer) IsProductQuery(query string) bool {
	query = strings.ToLower(strings.TrimSpace(query))

	// Negative patterns (NOT product queries)
	negativePatterns := []string{
		"hello", "hi", "hey", "how are you",
		"what is your name", "who are you",
		"help", "thanks", "thank you",
		"bye", "goodbye",
	}

	for _, pattern := range negativePatterns {
		if strings.Contains(query, pattern) {
			return false
		}
	}

	// Positive patterns (likely product queries)
	positivePatterns := []string{
		"buy", "purchase", "looking for", "need", "want",
		"find", "search", "show me", "price", "cheap",
		"best", "recommend", "where to buy",
	}

	for _, pattern := range positivePatterns {
		if strings.Contains(query, pattern) {
			return true
		}
	}

	// If query contains common product types, it's likely a product query
	productTypes := []string{
		"phone", "laptop", "computer", "tablet", "watch",
		"headphones", "camera", "tv", "monitor", "keyboard",
		"mouse", "speaker", "iphone", "ipad", "macbook",
		"samsung", "sony", "lg", "dell", "hp", "lenovo",
	}

	for _, productType := range productTypes {
		if strings.Contains(query, productType) {
			return true
		}
	}

	// Default: assume it's a product query if it passes basic validation
	return len(query) >= 3
}

// ExtractBrand tries to extract brand name from query
func (o *QueryOptimizer) ExtractBrand(query string) string {
	query = strings.ToLower(query)

	// Common brands
	brands := []string{
		"apple", "samsung", "google", "sony", "lg", "huawei",
		"xiaomi", "oppo", "vivo", "oneplus", "motorola", "nokia",
		"dell", "hp", "lenovo", "asus", "acer", "msi",
		"microsoft", "surface", "bose", "jbl", "beats",
		"canon", "nikon", "panasonic", "gopro",
	}

	for _, brand := range brands {
		if strings.Contains(query, brand) {
			return capitalizeFirst(brand)
		}
	}

	return ""
}

// ExtractProductType tries to extract product type from query
func (o *QueryOptimizer) ExtractProductType(query string) string {
	query = strings.ToLower(query)

	// Common product types
	productTypes := map[string]string{
		"phone":      "Smartphone",
		"smartphone": "Smartphone",
		"iphone":     "Smartphone",
		"laptop":     "Laptop",
		"notebook":   "Laptop",
		"macbook":    "Laptop",
		"tablet":     "Tablet",
		"ipad":       "Tablet",
		"watch":      "Smartwatch",
		"smartwatch": "Smartwatch",
		"headphones": "Headphones",
		"earbuds":    "Earbuds",
		"camera":     "Camera",
		"tv":         "Television",
		"television": "Television",
		"monitor":    "Monitor",
		"keyboard":   "Keyboard",
		"mouse":      "Mouse",
		"speaker":    "Speaker",
	}

	for keyword, productType := range productTypes {
		if strings.Contains(query, keyword) {
			return productType
		}
	}

	return "Product"
}

// SuggestQueryImprovements suggests ways to improve the query
func (o *QueryOptimizer) SuggestQueryImprovements(query string) []string {
	suggestions := []string{}

	// Check if query is too vague
	if len(strings.Fields(query)) == 1 {
		suggestions = append(suggestions, "Try being more specific about the product model")
	}

	// Check if brand is missing
	brand := o.ExtractBrand(query)
	if brand == "" {
		suggestions = append(suggestions, "Consider specifying a brand")
	}

	// Check if query has too many parameters
	if o.storagePattern.MatchString(query) {
		suggestions = append(suggestions, "Storage size will be filtered automatically")
	}

	return suggestions
}

// Helper functions

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// Examples of optimization:
// Input: "iPhone 15 Pro Max 256GB Space Gray"
// Output: "iPhone 15 Pro Max"
//
// Input: "Samsung Galaxy S24 Ultra 512GB Unlocked"
// Output: "Samsung Galaxy S24 Ultra"
//
// Input: "MacBook Pro 16 inch M3 1TB Silver"
// Output: "MacBook Pro"
