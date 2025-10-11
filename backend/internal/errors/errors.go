package errors

import (
	"fmt"

	"mylittleprice/internal/constants"
)

// AppError represents a structured application error
type AppError struct {
	Code       string // Error code for client
	Message    string // Human-readable message
	HTTPStatus int    // HTTP status code
	Cause      error  // Original error (for logging)
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// WithCause adds the underlying cause to the error
func (e *AppError) WithCause(cause error) *AppError {
	return &AppError{
		Code:       e.Code,
		Message:    e.Message,
		HTTPStatus: e.HTTPStatus,
		Cause:      cause,
	}
}

// WithMessage overrides the default message
func (e *AppError) WithMessage(message string) *AppError {
	return &AppError{
		Code:       e.Code,
		Message:    message,
		HTTPStatus: e.HTTPStatus,
		Cause:      e.Cause,
	}
}

// ═══════════════════════════════════════════════════════════
// PREDEFINED ERRORS
// ═══════════════════════════════════════════════════════════

// Session Errors
var (
	ErrSessionNotFound = &AppError{
		Code:       constants.ErrCodeSessionNotFound,
		Message:    "Session not found or expired",
		HTTPStatus: 404,
	}

	ErrMaxSearchesReached = &AppError{
		Code:       constants.ErrCodeMaxSearchesReached,
		Message:    constants.MsgMaxSearchesReached,
		HTTPStatus: 429,
	}

	ErrMaxMessagesReached = &AppError{
		Code:       constants.ErrCodeMaxMessagesReached,
		Message:    "Maximum messages per session reached",
		HTTPStatus: 429,
	}
)

// Validation Errors
var (
	ErrInvalidRequest = &AppError{
		Code:       constants.ErrCodeInvalidRequest,
		Message:    "Invalid request format",
		HTTPStatus: 400,
	}

	ErrValidationError = &AppError{
		Code:       constants.ErrCodeValidationError,
		Message:    "Validation failed",
		HTTPStatus: 400,
	}

	ErrQueryTooShort = &AppError{
		Code:       constants.ErrCodeValidationError,
		Message:    "Query too short",
		HTTPStatus: 400,
	}

	ErrQueryTooLong = &AppError{
		Code:       constants.ErrCodeValidationError,
		Message:    "Query too long",
		HTTPStatus: 400,
	}
)

// AI Service Errors
var (
	ErrAIProcessing = &AppError{
		Code:       constants.ErrCodeAIError,
		Message:    "Failed to process message with AI",
		HTTPStatus: 500,
	}

	ErrEmptyAIResponse = &AppError{
		Code:       constants.ErrCodeAIError,
		Message:    "AI returned empty response",
		HTTPStatus: 500,
	}

	ErrInvalidAIResponse = &AppError{
		Code:       constants.ErrCodeAIError,
		Message:    "Invalid AI response format",
		HTTPStatus: 500,
	}
)

// Search Service Errors
var (
	ErrSearchFailed = &AppError{
		Code:       constants.ErrCodeSearchError,
		Message:    "Failed to search for products",
		HTTPStatus: 500,
	}

	ErrNoProductsFound = &AppError{
		Code:       constants.ErrCodeSearchError,
		Message:    constants.MsgNoProductsFound,
		HTTPStatus: 404,
	}

	ErrProductDetailsNotFound = &AppError{
		Code:       constants.ErrCodeSearchError,
		Message:    constants.MsgProductDetailsNotFound,
		HTTPStatus: 404,
	}
)

// Cache Errors
var (
	ErrCacheMiss = &AppError{
		Code:       constants.ErrCodeCacheError,
		Message:    "Cache miss",
		HTTPStatus: 404,
	}

	ErrCacheWrite = &AppError{
		Code:       constants.ErrCodeCacheError,
		Message:    "Failed to write to cache",
		HTTPStatus: 500,
	}
)

// Internal Errors
var (
	ErrInternal = &AppError{
		Code:       constants.ErrCodeInternalError,
		Message:    "Internal server error",
		HTTPStatus: 500,
	}

	ErrDependencyFailure = &AppError{
		Code:       constants.ErrCodeInternalError,
		Message:    "Dependency service failure",
		HTTPStatus: 503,
	}
)

// ═══════════════════════════════════════════════════════════
// ERROR RESPONSE HELPERS
// ═══════════════════════════════════════════════════════════

// ErrorResponse represents JSON error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// ToResponse converts AppError to ErrorResponse
func (e *AppError) ToResponse() ErrorResponse {
	return ErrorResponse{
		Error:   "true",
		Message: e.Message,
		Code:    e.Code,
	}
}

// NewErrorResponse creates error response from any error
func NewErrorResponse(err error) ErrorResponse {
	if appErr, ok := err.(*AppError); ok {
		return appErr.ToResponse()
	}

	return ErrorResponse{
		Error:   "true",
		Message: err.Error(),
		Code:    constants.ErrCodeInternalError,
	}
}
