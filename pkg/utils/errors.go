package utils

import (
	"fmt"
	"net/http"
)

// APIError represents errors returned by the API
type APIError struct {
	StatusCode int               `json:"status_code"`
	Code       *string           `json:"code,omitempty"`
	Message    string            `json:"msg"`
	Headers    http.Header       `json:"-"`
	Data       interface{}       `json:"data,omitempty"`
}

func (e *APIError) Error() string {
	if e.Code != nil {
		return fmt.Sprintf("API error %d (%s): %s", e.StatusCode, *e.Code, e.Message)
	}
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// ClientError represents 4xx errors
type ClientError struct {
	*APIError
}

// ServerError represents 5xx errors
type ServerError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("Server error %d: %s", e.StatusCode, e.Message)
}

// NewClientError creates a new client error
func NewClientError(statusCode int, code *string, message string, headers http.Header, data interface{}) *ClientError {
	return &ClientError{
		APIError: &APIError{
			StatusCode: statusCode,
			Code:       code,
			Message:    message,
			Headers:    headers,
			Data:       data,
		},
	}
}

// NewServerError creates a new server error
func NewServerError(statusCode int, message string) *ServerError {
	return &ServerError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}