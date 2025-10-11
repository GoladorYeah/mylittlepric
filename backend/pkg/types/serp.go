package types

// GoogleShoppingResponse represents SERP API Google Shopping response
type GoogleShoppingResponse struct {
	SearchMetadata SearchMetadata `json:"search_metadata"`
	Shopping       []ShoppingItem `json:"shopping"`
}

type SearchMetadata struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// ShoppingItem represents a single product from Google Shopping
type ShoppingItem struct {
	Position    int     `json:"position"`
	Title       string  `json:"title"`
	Link        string  `json:"link"`
	ProductLink string  `json:"product_link"`
	ProductID   string  `json:"product_id"`
	Thumbnail   string  `json:"thumbnail"`
	Price       string  `json:"price"`
	OldPrice    string  `json:"old_price,omitempty"`
	Merchant    string  `json:"merchant"`
	Rating      float32 `json:"rating,omitempty"`
	Reviews     int     `json:"reviews,omitempty"`
	Delivery    string  `json:"delivery,omitempty"`
	// Google Immersive Product token for detailed view
	PageToken string `json:"page_token,omitempty"`
	// Alternative fields from different SERP responses
	SerpAPILink string `json:"serpapi_link,omitempty"`
}

// GoogleImmersiveProductResponse represents detailed product info
type GoogleImmersiveProductResponse struct {
	SearchMetadata   SearchMetadata      `json:"search_metadata"`
	SearchParameters ImmersiveParameters `json:"search_parameters"`
	ProductResults   ProductResults      `json:"product_results"`
}

type ImmersiveParameters struct {
	Engine    string `json:"engine"`
	PageToken string `json:"page_token"`
}

type ProductResults struct {
	Title           string            `json:"title"`
	Description     string            `json:"description,omitempty"`
	Price           string            `json:"price"`
	ExtractedPrice  float64           `json:"extracted_price,omitempty"`
	Rating          float32           `json:"rating,omitempty"`
	Reviews         int               `json:"reviews,omitempty"`
	Thumbnails      []string          `json:"thumbnails"`
	Offers          []Offer           `json:"offers"`
	Variants        []Variant         `json:"variants,omitempty"`
	MoreOptions     []MoreOption      `json:"more_options,omitempty"`
	Videos          []Video           `json:"videos,omitempty"`
	AboutTheProduct AboutProduct      `json:"about_the_product,omitempty"`
	Specifications  []Spec            `json:"specifications,omitempty"`
	Ratings         []RatingBreakdown `json:"ratings,omitempty"`
	ReviewsImages   []string          `json:"reviews_images,omitempty"`
	StoresNextToken string            `json:"stores_next_page_token,omitempty"`
}

type AboutProduct struct {
	Description string   `json:"description"`
	Highlights  []string `json:"highlights,omitempty"`
}

type Spec struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

type RatingBreakdown struct {
	Stars  int `json:"stars"`
	Amount int `json:"amount"`
}

type Offer struct {
	Merchant       string  `json:"merchant"`
	MerchantLink   string  `json:"merchant_link"`
	Price          string  `json:"price"`
	ExtractedPrice float64 `json:"extracted_price,omitempty"`
	Availability   string  `json:"availability,omitempty"`
	Shipping       string  `json:"shipping,omitempty"`
	Rating         float32 `json:"rating,omitempty"`
	Reviews        int     `json:"reviews,omitempty"`
	Link           string  `json:"link"`
}

type Variant struct {
	Title string        `json:"title"`
	Items []VariantItem `json:"items"`
}

type VariantItem struct {
	Name        string `json:"name"`
	Selected    bool   `json:"selected,omitempty"`
	Available   bool   `json:"available"`
	SerpAPILink string `json:"serpapi_link"`
}

type MoreOption struct {
	Title       string  `json:"title"`
	Thumbnail   string  `json:"thumbnail"`
	Price       string  `json:"price"`
	Rating      float32 `json:"rating,omitempty"`
	Reviews     int     `json:"reviews,omitempty"`
	SerpAPILink string  `json:"serpapi_link"`
}

type Video struct {
	Title     string `json:"title"`
	Link      string `json:"link"`
	Source    string `json:"source"`
	Channel   string `json:"channel,omitempty"`
	Duration  string `json:"duration,omitempty"`
	Thumbnail string `json:"thumbnail"`
}

// ExtractPageToken extracts page_token from serpapi_link
func ExtractPageToken(serpAPILink string) string {
	// serpapi_link format: https://serpapi.com/search.json?engine=google_immersive_product&page_token=XXX
	// We need to extract the page_token value
	// For simplicity, we'll return the full serpapi_link and extract token in service
	return serpAPILink
}
