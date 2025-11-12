package services

// SessionAdapter adapts SessionService for middleware use
// This avoids circular dependencies by returning generic interface
type SessionAdapter struct {
	*SessionService
}

// GetSession wraps SessionService.GetSession to return interface{}
func (a *SessionAdapter) GetSession(sessionID string) (interface{}, error) {
	return a.SessionService.GetSession(sessionID)
}
