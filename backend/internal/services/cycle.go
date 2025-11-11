package services

import (
	"mylittleprice/internal/models"
)

// CycleService handles cycle-related operations for prompt management
// Separated from SessionService for better SRP (Single Responsibility Principle)
type CycleService struct {
	universalPromptMgr *UniversalPromptManager
}

// NewCycleService creates a new CycleService instance
func NewCycleService() *CycleService {
	return &CycleService{
		universalPromptMgr: NewUniversalPromptManager(),
	}
}

// GetUniversalPromptManager returns the universal prompt manager
func (s *CycleService) GetUniversalPromptManager() *UniversalPromptManager {
	return s.universalPromptMgr
}

// IncrementCycleIterationInMemory increments the iteration using an in-memory session (avoids N+1)
// Returns true if we should start a new cycle
func (s *CycleService) IncrementCycleIterationInMemory(session *models.ChatSession) bool {
	shouldStartNew := !s.universalPromptMgr.IncrementIteration(&session.CycleState)
	return shouldStartNew
}

// StartNewCycleInMemory starts a new cycle using an in-memory session (avoids N+1)
func (s *CycleService) StartNewCycleInMemory(session *models.ChatSession, lastRequest string, products []models.ProductInfo) {
	s.universalPromptMgr.StartNewCycle(&session.CycleState, lastRequest, products)
}

// AddToCycleHistoryInMemory adds a message to cycle history using an in-memory session (avoids N+1)
func (s *CycleService) AddToCycleHistoryInMemory(session *models.ChatSession, role, content string) {
	s.universalPromptMgr.AddToCycleHistory(&session.CycleState, role, content)
}

// InitializeCycleState initializes a new cycle state for a session
func (s *CycleService) InitializeCycleState() models.CycleState {
	return s.universalPromptMgr.InitializeCycleState()
}
