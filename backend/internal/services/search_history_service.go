package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"mylittleprice/ent"
	"mylittleprice/ent/searchhistory"
	"mylittleprice/internal/models"
)

type SearchHistoryService struct {
	client *ent.Client
}

func NewSearchHistoryService(client *ent.Client) *SearchHistoryService {
	return &SearchHistoryService{
		client: client,
	}
}

// SaveSearchHistory saves a search to history
func (s *SearchHistoryService) SaveSearchHistory(ctx context.Context, history *models.SearchHistory) error {
	// Set expiration for anonymous users (30 days)
	if history.UserID == nil && history.ExpiresAt == nil {
		expiresAt := time.Now().Add(30 * 24 * time.Hour)
		history.ExpiresAt = &expiresAt
	}

	if history.ID == uuid.Nil {
		history.ID = uuid.New()
	}

	if history.CreatedAt.IsZero() {
		history.CreatedAt = time.Now()
	}

	// Create using Ent
	builder := s.client.SearchHistory.Create().
		SetID(history.ID).
		SetSearchQuery(history.SearchQuery).
		SetSearchType(history.SearchType).
		SetCountryCode(history.CountryCode).
		SetLanguageCode(history.LanguageCode).
		SetCurrency(history.Currency).
		SetResultCount(history.ResultCount).
		SetCreatedAt(history.CreatedAt)

	// Set optional fields
	if history.UserID != nil {
		builder.SetUserID(*history.UserID)
	}

	if history.SessionID != nil {
		builder.SetSessionID(*history.SessionID)
	}

	if history.OptimizedQuery != nil {
		builder.SetOptimizedQuery(*history.OptimizedQuery)
	}

	if history.Category != nil {
		builder.SetCategory(*history.Category)
	}

	if history.ProductsFound != nil && len(history.ProductsFound) > 0 {
		// Convert []models.ProductCard to []map[string]interface{}
		products := make([]map[string]interface{}, len(history.ProductsFound))
		for i, p := range history.ProductsFound {
			products[i] = map[string]interface{}{
				"name":        p.Name,
				"price":       p.Price,
				"old_price":   p.OldPrice,
				"link":        p.Link,
				"image":       p.Image,
				"description": p.Description,
				"badge":       p.Badge,
				"page_token":  p.PageToken,
			}
		}
		builder.SetProductsFound(products)
	}

	if history.ExpiresAt != nil {
		builder.SetExpiresAt(*history.ExpiresAt)
	}

	_, err := builder.Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to save search history: %w", err)
	}

	return nil
}

// GetUserSearchHistory retrieves search history for a user or session
func (s *SearchHistoryService) GetUserSearchHistory(ctx context.Context, userID *uuid.UUID, sessionID *string, limit, offset int) (*models.SearchHistoryListResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Build query
	query := s.client.SearchHistory.Query()

	if userID != nil {
		// For authenticated users, show all their history
		query = query.Where(searchhistory.UserIDEQ(*userID))
	} else if sessionID != nil {
		// For anonymous users with session_id, only show their session's searches
		query = query.Where(
			searchhistory.And(
				searchhistory.UserIDIsNil(),
				searchhistory.SessionIDEQ(*sessionID),
				searchhistory.Or(
					searchhistory.ExpiresAtIsNil(),
					searchhistory.ExpiresAtGT(time.Now()),
				),
			),
		)
	} else {
		// No user_id or session_id - return empty result
		return &models.SearchHistoryListResponse{
			Items:   []models.SearchHistory{},
			Total:   0,
			Limit:   limit,
			Offset:  offset,
			HasMore: false,
		}, nil
	}

	// Get total count
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count search history: %w", err)
	}

	// Get items
	items, err := query.
		Order(ent.Desc(searchhistory.FieldCreatedAt)).
		Limit(limit).
		Offset(offset).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get search history: %w", err)
	}

	// Convert Ent entities to response models
	responseItems := make([]models.SearchHistory, len(items))
	for i, item := range items {
		// Convert UUID fields to pointers
		var userID *uuid.UUID
		if item.UserID != uuid.Nil {
			userID = &item.UserID
		}

		// Convert string fields to pointers
		var sessionID, optimizedQuery, category, clickedProductID *string
		if item.SessionID != "" {
			sessionID = &item.SessionID
		}
		if item.OptimizedQuery != "" {
			optimizedQuery = &item.OptimizedQuery
		}
		if item.Category != "" {
			category = &item.Category
		}
		if item.ClickedProductID != "" {
			clickedProductID = &item.ClickedProductID
		}

		// Convert time fields to pointers
		var expiresAt *time.Time
		if !item.ExpiresAt.IsZero() {
			expiresAt = &item.ExpiresAt
		}

		responseItems[i] = models.SearchHistory{
			ID:               item.ID,
			UserID:           userID,
			SessionID:        sessionID,
			SearchQuery:      item.SearchQuery,
			OptimizedQuery:   optimizedQuery,
			SearchType:       item.SearchType,
			Category:         category,
			CountryCode:      item.CountryCode,
			LanguageCode:     item.LanguageCode,
			Currency:         item.Currency,
			ResultCount:      item.ResultCount,
			ClickedProductID: clickedProductID,
			CreatedAt:        item.CreatedAt,
			ExpiresAt:        expiresAt,
		}

		// Convert products from []map[string]interface{} to []models.ProductCard
		if item.ProductsFound != nil && len(item.ProductsFound) > 0 {
			products := make([]models.ProductCard, len(item.ProductsFound))
			for j, p := range item.ProductsFound {
				products[j] = models.ProductCard{
					Name:        getStringFromMap(p, "name"),
					Price:       getStringFromMap(p, "price"),
					OldPrice:    getStringFromMap(p, "old_price"),
					Link:        getStringFromMap(p, "link"),
					Image:       getStringFromMap(p, "image"),
					Description: getStringFromMap(p, "description"),
					Badge:       getStringFromMap(p, "badge"),
					PageToken:   getStringFromMap(p, "page_token"),
				}
			}
			responseItems[i].ProductsFound = products
		}
	}

	return &models.SearchHistoryListResponse{
		Items:   responseItems,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: offset+len(items) < total,
	}, nil
}

// Helper function to safely extract string from map
func getStringFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// DeleteSearchHistory deletes a search history entry
func (s *SearchHistoryService) DeleteSearchHistory(ctx context.Context, id uuid.UUID, userID *uuid.UUID) error {
	query := s.client.SearchHistory.Delete().Where(searchhistory.IDEQ(id))

	if userID != nil {
		// For authenticated users, ensure they own the history entry
		query = query.Where(searchhistory.UserIDEQ(*userID))
	} else {
		// For anonymous users, ensure it's not owned by anyone
		query = query.Where(searchhistory.UserIDIsNil())
	}

	rowsAffected, err := query.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete search history: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("search history not found or permission denied")
	}

	return nil
}

// DeleteAllUserSearchHistory deletes all search history for a user
func (s *SearchHistoryService) DeleteAllUserSearchHistory(ctx context.Context, userID uuid.UUID) error {
	_, err := s.client.SearchHistory.Delete().
		Where(searchhistory.UserIDEQ(userID)).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete all search history: %w", err)
	}
	return nil
}

// UpdateClickedProduct updates which product was clicked in a search
func (s *SearchHistoryService) UpdateClickedProduct(ctx context.Context, historyID uuid.UUID, productID string) error {
	_, err := s.client.SearchHistory.UpdateOneID(historyID).
		SetClickedProductID(productID).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to update clicked product: %w", err)
	}
	return nil
}

// CleanupExpiredAnonymousHistory removes expired anonymous search history
func (s *SearchHistoryService) CleanupExpiredAnonymousHistory(ctx context.Context) (int64, error) {
	rowsAffected, err := s.client.SearchHistory.Delete().
		Where(
			searchhistory.And(
				searchhistory.UserIDIsNil(),
				searchhistory.ExpiresAtLT(time.Now()),
			),
		).
		Exec(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired history: %w", err)
	}

	return int64(rowsAffected), nil
}
