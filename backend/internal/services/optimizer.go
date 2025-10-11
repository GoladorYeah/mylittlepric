package services

import (
	"regexp"
	"strings"
)

type QueryOptimizer struct {
	storagePattern   *regexp.Regexp
	colorPattern     *regexp.Regexp
	sizePattern      *regexp.Regexp
	conditionPattern *regexp.Regexp
	carrierPattern   *regexp.Regexp
	prefixPattern    *regexp.Regexp
}

func NewQueryOptimizer() *QueryOptimizer {
	return &QueryOptimizer{
		storagePattern:   regexp.MustCompile(`(?i)\b\d+\s*(gb|tb|mb|gigabyte|terabyte)\b`),
		colorPattern:     regexp.MustCompile(`(?i)\b(black|white|silver|gold|rose\s*gold|space\s*gray|grey|blue|red|green|yellow|purple|pink|orange|midnight|starlight|cosmic|deep|titanium|natural)\b`),
		sizePattern:      regexp.MustCompile(`(?i)\b(small|medium|large|x+l|[0-9]+(\.[0-9]+)?\s*(inch|cm|mm|meter|foot|feet)|size\s*[0-9]+)\b`),
		conditionPattern: regexp.MustCompile(`(?i)\b(new|refurbished|used|like\s*new|open\s*box|factory\s*sealed|unlocked)\b`),
		carrierPattern:   regexp.MustCompile(`(?i)\b(unlocked|verizon|at&t|t-mobile|sprint|gsm|cdma|dual\s*sim)\b`),
		prefixPattern:    regexp.MustCompile(`^\*+\s*`),
	}
}

func (o *QueryOptimizer) OptimizeQuery(query, searchType string) string {
	optimized := o.prefixPattern.ReplaceAllString(query, "")
	optimized = strings.TrimSpace(optimized)

	switch searchType {
	case "exact":
		optimized = o.removeAllParameters(optimized)
	case "parameters":
		optimized = o.removeSpecifications(optimized)
	case "category":
		optimized = o.broadenQuery(optimized)
	default:
		optimized = o.removeSpecifications(optimized)
	}

	optimized = o.cleanSpaces(optimized)
	return optimized
}

func (o *QueryOptimizer) removeAllParameters(query string) string {
	result := query
	result = o.storagePattern.ReplaceAllString(result, "")
	result = o.colorPattern.ReplaceAllString(result, "")
	result = o.sizePattern.ReplaceAllString(result, "")
	result = o.conditionPattern.ReplaceAllString(result, "")
	result = o.carrierPattern.ReplaceAllString(result, "")
	return result
}

func (o *QueryOptimizer) removeSpecifications(query string) string {
	result := query
	result = o.storagePattern.ReplaceAllString(result, "")
	result = o.colorPattern.ReplaceAllString(result, "")
	result = o.carrierPattern.ReplaceAllString(result, "")
	return result
}

func (o *QueryOptimizer) broadenQuery(query string) string {
	query = strings.TrimSpace(query)
	words := strings.Fields(query)

	// Remove common filler words
	fillerWords := []string{"need", "want", "looking", "for", "a", "an", "the"}
	filteredWords := []string{}

	for _, word := range words {
		isFilter := false
		for _, filler := range fillerWords {
			if strings.ToLower(word) == filler {
				isFilter = true
				break
			}
		}
		if !isFilter {
			filteredWords = append(filteredWords, word)
		}
	}

	// Keep ALL meaningful words for textiles/furniture/decor
	if len(filteredWords) <= 2 {
		return strings.Join(filteredWords, " ")
	}

	// For longer queries, keep up to 5 most relevant words
	if len(filteredWords) > 5 {
		return strings.Join(filteredWords[:5], " ")
	}

	return strings.Join(filteredWords, " ")
}

func (o *QueryOptimizer) cleanSpaces(query string) string {
	spacePattern := regexp.MustCompile(`\s+`)
	result := spacePattern.ReplaceAllString(query, " ")
	return strings.TrimSpace(result)
}

func (o *QueryOptimizer) ValidateQuery(query string) (bool, string) {
	query = strings.TrimSpace(query)

	if len(query) < 2 {
		return false, "Query too short"
	}

	if len(query) > 200 {
		return false, "Query too long"
	}

	alphanumericPattern := regexp.MustCompile(`[a-zA-Z0-9]`)
	if !alphanumericPattern.MatchString(query) {
		return false, "Query must contain letters or numbers"
	}

	return true, ""
}

func (o *QueryOptimizer) IsProductQuery(query string) bool {
	query = strings.ToLower(strings.TrimSpace(query))

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

	productTypes := []string{
		"phone", "laptop", "computer", "tablet", "watch",
		"headphones", "camera", "tv", "monitor", "keyboard",
		"mouse", "speaker", "iphone", "ipad", "macbook",
		"samsung", "sony", "lg", "dell", "hp", "lenovo",
		"pillow", "blanket", "towel", "carpet", "curtain",
		"lamp", "mirror", "vase", "sofa", "table", "chair",
	}

	for _, productType := range productTypes {
		if strings.Contains(query, productType) {
			return true
		}
	}

	return len(query) >= 3
}

func (o *QueryOptimizer) ExtractBrand(query string) string {
	query = strings.ToLower(query)

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

func (o *QueryOptimizer) ExtractProductType(query string) string {
	query = strings.ToLower(query)

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
		"pillow":     "Pillow",
		"blanket":    "Blanket",
		"towel":      "Towel",
	}

	for keyword, productType := range productTypes {
		if strings.Contains(query, keyword) {
			return productType
		}
	}

	return "Product"
}

func (o *QueryOptimizer) SuggestQueryImprovements(query string) []string {
	suggestions := []string{}

	if len(strings.Fields(query)) == 1 {
		suggestions = append(suggestions, "Try being more specific about the product model")
	}

	brand := o.ExtractBrand(query)
	if brand == "" {
		suggestions = append(suggestions, "Consider specifying a brand")
	}

	if o.storagePattern.MatchString(query) {
		suggestions = append(suggestions, "Storage size will be filtered automatically")
	}

	return suggestions
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
