package domain

import "mylittleprice/internal/constants"

// ═══════════════════════════════════════════════════════════
// SEARCH TYPES
// ═══════════════════════════════════════════════════════════

// SearchType represents type of product search
type SearchType string

const (
	SearchTypeExact      SearchType = constants.SearchTypeExact
	SearchTypeParameters SearchType = constants.SearchTypeParameters
	SearchTypeCategory   SearchType = constants.SearchTypeCategory
)

// IsValid checks if search type is valid
func (st SearchType) IsValid() bool {
	switch st {
	case SearchTypeExact, SearchTypeParameters, SearchTypeCategory:
		return true
	default:
		return false
	}
}

// MaxProducts returns maximum products for this search type
func (st SearchType) MaxProducts() int {
	switch st {
	case SearchTypeExact:
		return constants.MaxProductsExactSearch
	case SearchTypeParameters:
		return constants.MaxProductsParameterSearch
	case SearchTypeCategory:
		return constants.MaxProductsCategorySearch
	default:
		return constants.MaxProductsParameterSearch
	}
}

// RelevanceThreshold returns minimum relevance score
func (st SearchType) RelevanceThreshold() float32 {
	switch st {
	case SearchTypeExact:
		return constants.MinRelevanceThresholdExact
	case SearchTypeParameters:
		return constants.MinRelevanceThresholdParameters
	case SearchTypeCategory:
		return constants.MinRelevanceThresholdCategory
	default:
		return constants.MinRelevanceThresholdParameters
	}
}

// SearchContext encapsulates search parameters
type SearchContext struct {
	Query     string
	Type      SearchType
	Locale    Locale
	Category  string
	Optimized bool // Whether query has been optimized
}

// NewSearchContext creates a new search context
func NewSearchContext(query string, searchType SearchType, locale Locale) SearchContext {
	return SearchContext{
		Query:     query,
		Type:      searchType,
		Locale:    locale,
		Optimized: false,
	}
}

// ═══════════════════════════════════════════════════════════
// PRODUCT CATEGORY
// ═══════════════════════════════════════════════════════════

// Category represents product category
type Category string

const (
	CategoryElectronics Category = "electronics"
	CategoryClothing    Category = "clothing"
	CategoryFurniture   Category = "furniture"
	CategoryKitchen     Category = "kitchen"
	CategorySports      Category = "sports"
	CategoryTools       Category = "tools"
	CategoryDecor       Category = "decor"
	CategoryTextiles    Category = "textiles"
)

// IsValid checks if category is valid
func (c Category) IsValid() bool {
	switch c {
	case CategoryElectronics, CategoryClothing, CategoryFurniture,
		CategoryKitchen, CategorySports, CategoryTools,
		CategoryDecor, CategoryTextiles:
		return true
	default:
		return false
	}
}

// String returns string representation
func (c Category) String() string {
	return string(c)
}
