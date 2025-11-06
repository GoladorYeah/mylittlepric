package models

// ═══════════════════════════════════════════════════════════
// PRODUCT MODELS
// ═══════════════════════════════════════════════════════════

type ProductInfo struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductCard struct {
	Name        string `json:"name"`
	Price       string `json:"price"`
	OldPrice    string `json:"old_price,omitempty"`
	Link        string `json:"link"`
	Image       string `json:"image"`
	Description string `json:"description,omitempty"`
	Badge       string `json:"badge,omitempty"`
	PageToken   string `json:"page_token"`
}

type ProductDetailsRequest struct {
	PageToken string `json:"page_token"`
	Country   string `json:"country"`
}

type ProductDetailsResponse struct {
	Type            string                `json:"type"`
	Title           string                `json:"title"`
	Price           string                `json:"price"`
	Rating          float32               `json:"rating,omitempty"`
	Reviews         int                   `json:"reviews,omitempty"`
	Description     string                `json:"description,omitempty"`
	Images          []string              `json:"images,omitempty"`
	Specifications  []Specification       `json:"specifications,omitempty"`
	Variants        []Variant             `json:"variants,omitempty"`
	Offers          []Offer               `json:"offers"`
	Videos          []interface{}         `json:"videos,omitempty"`
	MoreOptions     []interface{}         `json:"more_options,omitempty"`
	RatingBreakdown []RatingBreakdownItem `json:"rating_breakdown,omitempty"`
}

type Specification struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

type Variant struct {
	Title string        `json:"title"`
	Items []interface{} `json:"items"`
}

type Offer struct {
	Merchant          string   `json:"merchant"`
	Logo              string   `json:"logo,omitempty"`
	Price             string   `json:"price"`
	ExtractedPrice    float64  `json:"extracted_price,omitempty"`
	Currency          string   `json:"currency,omitempty"`
	Link              string   `json:"link"`
	Title             string   `json:"title,omitempty"`
	Availability      string   `json:"availability,omitempty"`
	Shipping          string   `json:"shipping,omitempty"`
	ShippingExtracted float64  `json:"shipping_extracted,omitempty"`
	Total             string   `json:"total,omitempty"`
	ExtractedTotal    float64  `json:"extracted_total,omitempty"`
	Rating            float32  `json:"rating,omitempty"`
	Reviews           int      `json:"reviews,omitempty"`
	PaymentMethods    string   `json:"payment_methods,omitempty"`
	Tag               string   `json:"tag,omitempty"`
	DetailsAndOffers  []string `json:"details_and_offers,omitempty"`
	MonthlyPaymentDur int      `json:"monthly_payment_duration,omitempty"`
	DownPayment       string   `json:"down_payment,omitempty"`
}

type RatingBreakdownItem struct {
	Stars  int `json:"stars"`
	Amount int `json:"amount"`
}
