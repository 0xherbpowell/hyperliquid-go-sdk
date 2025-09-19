package errors

import (
	"fmt"
	"net/http"
)

// APIError represents an error returned by the Hyperliquid API
type APIError struct {
	StatusCode int    `json:"status_code"`
	Code       string `json:"code,omitempty"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("API error %d: %s - %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// NewAPIError creates a new API error
func NewAPIError(statusCode int, code, message string, data any) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Data:       data,
	}
}

// IsClientError returns true if the error is a 4xx client error
func (e *APIError) IsClientError() bool {
	return e.StatusCode >= 400 && e.StatusCode < 500
}

// IsServerError returns true if the error is a 5xx server error
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500
}

// Common error types
var (
	ErrInvalidAddress     = fmt.Errorf("invalid address format")
	ErrInvalidPrivateKey  = fmt.Errorf("invalid private key")
	ErrInvalidSignature   = fmt.Errorf("invalid signature")
	ErrInvalidCloid       = fmt.Errorf("invalid cloid format")
	ErrWebSocketNotReady  = fmt.Errorf("websocket connection not ready")
	ErrSubscriptionExists = fmt.Errorf("subscription already exists")
	ErrFloatPrecision     = fmt.Errorf("float precision error in conversion")
)

// HTTP status code errors
func NewBadRequestError(message string) *APIError {
	return NewAPIError(http.StatusBadRequest, "BAD_REQUEST", message, nil)
}

func NewUnauthorizedError(message string) *APIError {
	return NewAPIError(http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

func NewForbiddenError(message string) *APIError {
	return NewAPIError(http.StatusForbidden, "FORBIDDEN", message, nil)
}

func NewNotFoundError(message string) *APIError {
	return NewAPIError(http.StatusNotFound, "NOT_FOUND", message, nil)
}

func NewInternalServerError(message string) *APIError {
	return NewAPIError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, nil)
}
