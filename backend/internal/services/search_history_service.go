package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"mylittleprice/internal/models"
)

type SearchHistoryService struct {
	db *sqlx.DB
}

func NewSearchHistoryService(db *sqlx.DB) *SearchHistoryService {
	return &SearchHistoryService{
		db: db,
	}
}

// SaveSearchHistory saves a search to history
func (s *SearchHistoryService) SaveSearchHistory(ctx context.Context, history *models.SearchHistory) error {
	// Convert products to JSONB
	var productsJSON []byte
	var err error
	if history.ProductsFound != nil && len(history.ProductsFound) > 0 {
		// Save all products (no limit)
		productsJSON, err = json.Marshal(history.ProductsFound)
		if err != nil {
			return fmt.Errorf("failed to marshal products: %w", err)
		}
	}

	// Set expiration for anonymous users (30 days)
	if history.UserID == nil && history.ExpiresAt == nil {
		expiresAt := time.Now().Add(30 * 24 * time.Hour)
		history.ExpiresAt = &expiresAt
	}

	query := `
		INSERT INTO search_history (
			id, user_id, session_id, search_query, optimized_query,
			search_type, category, country_code, language_code,
			currency, result_count, products_found, created_at, expires_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
		RETURNING id, created_at
	`

	if history.ID == uuid.Nil {
		history.ID = uuid.New()
	}

	if history.CreatedAt.IsZero() {
		history.CreatedAt = time.Now()
	}

	err = s.db.QueryRowContext(
		ctx,
		query,
		history.ID,
		history.UserID,
		history.SessionID,
		history.SearchQuery,
		history.OptimizedQuery,
		history.SearchType,
		history.Category,
		history.CountryCode,
		history.LanguageCode,
		history.Currency,
		history.ResultCount,
		productsJSON,
		history.CreatedAt,
		history.ExpiresAt,
	).Scan(&history.ID, &history.CreatedAt)

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

	// Build query based on whether userID is provided
	var countQuery string
	var selectQuery string
	var args []interface{}

	if userID != nil {
		// For authenticated users, show all their history
		countQuery = `SELECT COUNT(*) FROM search_history WHERE user_id = $1`
		selectQuery = `
			SELECT id, user_id, session_id, search_query, optimized_query,
				   search_type, category, country_code, language_code,
				   currency, result_count, products_found, clicked_product_id,
				   created_at, expires_at
			FROM search_history
			WHERE user_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{userID, limit, offset}
	} else if sessionID != nil {
		// For anonymous users with session_id, only show their session's searches
		countQuery = `SELECT COUNT(*) FROM search_history WHERE user_id IS NULL AND session_id = $1 AND (expires_at IS NULL OR expires_at > NOW())`
		selectQuery = `
			SELECT id, user_id, session_id, search_query, optimized_query,
				   search_type, category, country_code, language_code,
				   currency, result_count, products_found, clicked_product_id,
				   created_at, expires_at
			FROM search_history
			WHERE user_id IS NULL AND session_id = $1 AND (expires_at IS NULL OR expires_at > NOW())
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{sessionID, limit, offset}
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
	var total int
	if userID != nil {
		err := s.db.GetContext(ctx, &total, countQuery, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to count search history: %w", err)
		}
	} else if sessionID != nil {
		err := s.db.GetContext(ctx, &total, countQuery, sessionID)
		if err != nil {
			return nil, fmt.Errorf("failed to count search history: %w", err)
		}
	}

	// Get items
	type SearchHistoryRow struct {
		ID               uuid.UUID      `db:"id"`
		UserID           *uuid.UUID     `db:"user_id"`
		SessionID        *string        `db:"session_id"`
		SearchQuery      string         `db:"search_query"`
		OptimizedQuery   *string        `db:"optimized_query"`
		SearchType       string         `db:"search_type"`
		Category         *string        `db:"category"`
		CountryCode      string         `db:"country_code"`
		LanguageCode     string         `db:"language_code"`
		Currency         string         `db:"currency"`
		ResultCount      int            `db:"result_count"`
		ProductsFound    sql.NullString `db:"products_found"`
		ClickedProductID *string        `db:"clicked_product_id"`
		CreatedAt        time.Time      `db:"created_at"`
		ExpiresAt        *time.Time     `db:"expires_at"`
	}

	var rows []SearchHistoryRow
	err := s.db.SelectContext(ctx, &rows, selectQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get search history: %w", err)
	}

	// Convert to response models
	items := make([]models.SearchHistory, len(rows))
	for i, row := range rows {
		items[i] = models.SearchHistory{
			ID:               row.ID,
			UserID:           row.UserID,
			SessionID:        row.SessionID,
			SearchQuery:      row.SearchQuery,
			OptimizedQuery:   row.OptimizedQuery,
			SearchType:       row.SearchType,
			Category:         row.Category,
			CountryCode:      row.CountryCode,
			LanguageCode:     row.LanguageCode,
			Currency:         row.Currency,
			ResultCount:      row.ResultCount,
			ClickedProductID: row.ClickedProductID,
			CreatedAt:        row.CreatedAt,
			ExpiresAt:        row.ExpiresAt,
		}

		// Unmarshal products
		if row.ProductsFound.Valid {
			var products []models.ProductCard
			if err := json.Unmarshal([]byte(row.ProductsFound.String), &products); err == nil {
				items[i].ProductsFound = products
			}
		}
	}

	return &models.SearchHistoryListResponse{
		Items:   items,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: offset+len(items) < total,
	}, nil
}

// DeleteSearchHistory deletes a search history entry
func (s *SearchHistoryService) DeleteSearchHistory(ctx context.Context, id uuid.UUID, userID *uuid.UUID) error {
	var query string
	var args []interface{}

	if userID != nil {
		// For authenticated users, ensure they own the history entry
		query = `DELETE FROM search_history WHERE id = $1 AND user_id = $2`
		args = []interface{}{id, userID}
	} else {
		// For anonymous users, allow deletion without user check
		query = `DELETE FROM search_history WHERE id = $1 AND user_id IS NULL`
		args = []interface{}{id}
	}

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete search history: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("search history not found or permission denied")
	}

	return nil
}

// DeleteAllUserSearchHistory deletes all search history for a user
func (s *SearchHistoryService) DeleteAllUserSearchHistory(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM search_history WHERE user_id = $1`
	_, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete all search history: %w", err)
	}
	return nil
}

// UpdateClickedProduct updates which product was clicked in a search
func (s *SearchHistoryService) UpdateClickedProduct(ctx context.Context, historyID uuid.UUID, productID string) error {
	query := `UPDATE search_history SET clicked_product_id = $1 WHERE id = $2`
	_, err := s.db.ExecContext(ctx, query, productID, historyID)
	if err != nil {
		return fmt.Errorf("failed to update clicked product: %w", err)
	}
	return nil
}

// CleanupExpiredAnonymousHistory removes expired anonymous search history
func (s *SearchHistoryService) CleanupExpiredAnonymousHistory(ctx context.Context) (int64, error) {
	query := `DELETE FROM search_history WHERE user_id IS NULL AND expires_at < NOW()`
	result, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired history: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}
