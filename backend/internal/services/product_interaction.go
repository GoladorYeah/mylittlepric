// backend/internal/services/product_interaction.go
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"mylittleprice/ent"
	"mylittleprice/ent/productinteraction"
	"mylittleprice/internal/models"
)

// ProductInteractionService tracks user interactions with products
type ProductInteractionService struct {
	db  *ent.Client
	ctx context.Context
}

// NewProductInteractionService creates a new product interaction service
func NewProductInteractionService(db *ent.Client) *ProductInteractionService {
	return &ProductInteractionService{
		db:  db,
		ctx: context.Background(),
	}
}

// TrackProductView tracks when a user views a product
func (s *ProductInteractionService) TrackProductView(
	userID *uuid.UUID,
	sessionID string,
	product *models.ProductInfo,
	messagePosition int,
	positionInResults int,
	searchQuery string,
	searchType string,
) error {
	builder := s.db.ProductInteraction.Create().
		SetSessionID(sessionID).
		SetProductID(product.ID).
		SetProductName(product.Name).
		SetInteractionType("viewed").
		SetMessagePosition(messagePosition).
		SetPositionInResults(positionInResults).
		SetInteractedAt(time.Now())

	if userID != nil {
		builder.SetUserID(*userID)
	}

	if product.Price > 0 {
		builder.SetProductPrice(product.Price)
	}

	if product.Currency != "" {
		builder.SetProductCurrency(product.Currency)
	}

	if product.Source != "" {
		builder.SetProductBrand(product.Source)
	}

	if product.URL != "" {
		builder.SetProductURL(product.URL)
	}

	if searchQuery != "" {
		builder.SetSearchQuery(searchQuery)
	}

	if searchType != "" {
		builder.SetSearchType(searchType)
	}

	// Store full product data for later analysis
	productData := map[string]interface{}{
		"id":       product.ID,
		"name":     product.Name,
		"price":    product.Price,
		"currency": product.Currency,
		"rating":   product.Rating,
		"source":   product.Source,
		"url":      product.URL,
	}
	builder.SetProductData(productData)

	_, err := builder.Save(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to track product view: %w", err)
	}

	return nil
}

// TrackProductClick tracks when a user clicks on a product
func (s *ProductInteractionService) TrackProductClick(
	userID *uuid.UUID,
	sessionID string,
	productID string,
) error {
	// Find existing interaction or create new one
	interaction, err := s.db.ProductInteraction.
		Query().
		Where(
			productinteraction.SessionID(sessionID),
			productinteraction.ProductID(productID),
		).
		First(s.ctx)

	if err != nil {
		// No existing interaction, create new click interaction
		builder := s.db.ProductInteraction.Create().
			SetSessionID(sessionID).
			SetProductID(productID).
			SetProductName("Unknown"). // We don't have full product info here
			SetInteractionType("clicked").
			SetClickCount(1).
			SetImplicitScore(0.7). // Clicks are strong positive signals
			SetInteractedAt(time.Now())

		if userID != nil {
			builder.SetUserID(*userID)
		}

		_, err = builder.Save(s.ctx)
		return err
	}

	// Update existing interaction
	_, err = s.db.ProductInteraction.UpdateOne(interaction).
		SetInteractionType("clicked").
		AddClickCount(1).
		SetImplicitScore(0.7).
		Save(s.ctx)

	return err
}

// TrackProductComparison tracks when a user adds product to comparison
func (s *ProductInteractionService) TrackProductComparison(
	userID *uuid.UUID,
	sessionID string,
	productID string,
) error {
	interaction, err := s.db.ProductInteraction.
		Query().
		Where(
			productinteraction.SessionID(sessionID),
			productinteraction.ProductID(productID),
		).
		First(s.ctx)

	if err != nil {
		// Create new comparison interaction
		builder := s.db.ProductInteraction.Create().
			SetSessionID(sessionID).
			SetProductID(productID).
			SetProductName("Unknown").
			SetInteractionType("compared").
			SetAddedToComparison(true).
			SetImplicitScore(0.5). // Comparisons indicate interest
			SetInteractedAt(time.Now())

		if userID != nil {
			builder.SetUserID(*userID)
		}

		_, err = builder.Save(s.ctx)
		return err
	}

	// Update existing interaction
	_, err = s.db.ProductInteraction.UpdateOne(interaction).
		SetInteractionType("compared").
		SetAddedToComparison(true).
		SetImplicitScore(0.5).
		Save(s.ctx)

	return err
}

// TrackProductDismissal tracks when a user dismisses/skips a product
func (s *ProductInteractionService) TrackProductDismissal(
	userID *uuid.UUID,
	sessionID string,
	productID string,
	feedback string,
) error {
	interaction, err := s.db.ProductInteraction.
		Query().
		Where(
			productinteraction.SessionID(sessionID),
			productinteraction.ProductID(productID),
		).
		First(s.ctx)

	if err != nil {
		// Create new dismissal interaction
		builder := s.db.ProductInteraction.Create().
			SetSessionID(sessionID).
			SetProductID(productID).
			SetProductName("Unknown").
			SetInteractionType("dismissed").
			SetImplicitScore(-0.3). // Negative signal
			SetInteractedAt(time.Now())

		if userID != nil {
			builder.SetUserID(*userID)
		}

		if feedback != "" {
			builder.SetFeedback(feedback)
		}

		_, err = builder.Save(s.ctx)
		return err
	}

	// Update existing interaction
	updater := s.db.ProductInteraction.UpdateOne(interaction).
		SetInteractionType("dismissed").
		SetImplicitScore(-0.3)

	if feedback != "" {
		updater.SetFeedback(feedback)
	}

	_, err = updater.Save(s.ctx)
	return err
}

// UpdateViewDuration updates how long a user viewed a product
func (s *ProductInteractionService) UpdateViewDuration(
	sessionID string,
	productID string,
	durationSeconds int,
) error {
	interaction, err := s.db.ProductInteraction.
		Query().
		Where(
			productinteraction.SessionID(sessionID),
			productinteraction.ProductID(productID),
		).
		First(s.ctx)

	if err != nil {
		return fmt.Errorf("interaction not found: %w", err)
	}

	// Longer view duration = higher interest
	implicitScore := 0.1
	if durationSeconds > 30 {
		implicitScore = 0.4
	} else if durationSeconds > 10 {
		implicitScore = 0.2
	}

	_, err = s.db.ProductInteraction.UpdateOne(interaction).
		SetViewDurationSeconds(durationSeconds).
		SetImplicitScore(implicitScore).
		Save(s.ctx)

	return err
}

// GetUserInteractions retrieves all interactions for a user
func (s *ProductInteractionService) GetUserInteractions(userID uuid.UUID, limit int) ([]*ent.ProductInteraction, error) {
	return s.db.ProductInteraction.
		Query().
		Where(productinteraction.UserID(userID)).
		Order(ent.Desc(productinteraction.FieldInteractedAt)).
		Limit(limit).
		All(s.ctx)
}

// GetSessionInteractions retrieves all interactions for a session
func (s *ProductInteractionService) GetSessionInteractions(sessionID string) ([]*ent.ProductInteraction, error) {
	return s.db.ProductInteraction.
		Query().
		Where(productinteraction.SessionID(sessionID)).
		Order(ent.Asc(productinteraction.FieldInteractionSequence)).
		All(s.ctx)
}

// GetProductPopularity gets popularity metrics for a product
func (s *ProductInteractionService) GetProductPopularity(productID string) (map[string]int, error) {
	interactions, err := s.db.ProductInteraction.
		Query().
		Where(productinteraction.ProductID(productID)).
		All(s.ctx)

	if err != nil {
		return nil, err
	}

	stats := map[string]int{
		"views":       0,
		"clicks":      0,
		"comparisons": 0,
		"dismissals":  0,
	}

	for _, interaction := range interactions {
		switch interaction.InteractionType {
		case "viewed":
			stats["views"]++
		case "clicked":
			stats["clicks"]++
		case "compared":
			stats["comparisons"]++
		case "dismissed":
			stats["dismissals"]++
		}
	}

	return stats, nil
}

// GetCategoryInteractions gets user's interactions within a category
func (s *ProductInteractionService) GetCategoryInteractions(userID uuid.UUID, category string) ([]*ent.ProductInteraction, error) {
	return s.db.ProductInteraction.
		Query().
		Where(
			productinteraction.UserID(userID),
			productinteraction.ProductCategory(category),
		).
		Order(ent.Desc(productinteraction.FieldImplicitScore)).
		All(s.ctx)
}

// GetTopProductsForUser gets top products based on implicit scores
func (s *ProductInteractionService) GetTopProductsForUser(userID uuid.UUID, limit int) ([]*ent.ProductInteraction, error) {
	return s.db.ProductInteraction.
		Query().
		Where(productinteraction.UserID(userID)).
		Where(productinteraction.ImplicitScoreGT(0)). // Only positive interactions
		Order(ent.Desc(productinteraction.FieldImplicitScore)).
		Limit(limit).
		All(s.ctx)
}

// GetBrandInteractions gets user's interactions with a specific brand
func (s *ProductInteractionService) GetBrandInteractions(userID uuid.UUID, brand string) ([]*ent.ProductInteraction, error) {
	return s.db.ProductInteraction.
		Query().
		Where(
			productinteraction.UserID(userID),
			productinteraction.ProductBrand(brand),
		).
		Order(ent.Desc(productinteraction.FieldInteractedAt)).
		All(s.ctx)
}

// GetInteractionStats gets summary statistics for a user
func (s *ProductInteractionService) GetInteractionStats(userID uuid.UUID) (map[string]interface{}, error) {
	interactions, err := s.db.ProductInteraction.
		Query().
		Where(productinteraction.UserID(userID)).
		All(s.ctx)

	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_interactions": len(interactions),
		"views":              0,
		"clicks":             0,
		"comparisons":        0,
		"dismissals":         0,
		"avg_implicit_score": 0.0,
		"engagement_rate":    0.0,
	}

	totalScore := 0.0
	engagedCount := 0

	for _, interaction := range interactions {
		switch interaction.InteractionType {
		case "viewed":
			stats["views"] = stats["views"].(int) + 1
		case "clicked":
			stats["clicks"] = stats["clicks"].(int) + 1
			engagedCount++
		case "compared":
			stats["comparisons"] = stats["comparisons"].(int) + 1
			engagedCount++
		case "dismissed":
			stats["dismissals"] = stats["dismissals"].(int) + 1
		}
		totalScore += interaction.ImplicitScore
	}

	if len(interactions) > 0 {
		stats["avg_implicit_score"] = totalScore / float64(len(interactions))
		stats["engagement_rate"] = float64(engagedCount) / float64(len(interactions))
	}

	return stats, nil
}
