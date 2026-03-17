package evilmail

import (
	"errors"
	"fmt"
	"net/http"
)

// APIError represents an error response from the EvilMail API.
// It captures the HTTP status code, the raw response body, and any
// message extracted from the API's JSON error response.
type APIError struct {
	// StatusCode is the HTTP status code returned by the API.
	StatusCode int

	// Status is the HTTP status text (e.g., "404 Not Found").
	Status string

	// Message is the error message extracted from the API response body,
	// if available.
	Message string

	// Body is the raw response body bytes.
	Body []byte
}

// Error returns a human-readable representation of the API error.
func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("evilmail: API error %d %s: %s", e.StatusCode, e.Status, e.Message)
	}
	return fmt.Sprintf("evilmail: API error %d %s", e.StatusCode, e.Status)
}

// IsNotFound reports whether the error represents a 404 Not Found response.
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

// IsUnauthorized reports whether the error represents a 401 Unauthorized response.
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized
}

// IsForbidden reports whether the error represents a 403 Forbidden response.
func (e *APIError) IsForbidden() bool {
	return e.StatusCode == http.StatusForbidden
}

// IsRateLimited reports whether the error represents a 429 Too Many Requests response.
func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// IsServerError reports whether the error represents a 5xx server error response.
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}

// AuthError represents an authentication failure, typically caused
// by a missing or invalid API key.
type AuthError struct {
	Err error
}

// Error returns a human-readable representation of the authentication error.
func (e *AuthError) Error() string {
	return fmt.Sprintf("evilmail: authentication failed: %v", e.Err)
}

// Unwrap returns the underlying error.
func (e *AuthError) Unwrap() error {
	return e.Err
}

// ValidationError represents a client-side validation failure that
// occurs before a request is sent to the API.
type ValidationError struct {
	// Field is the name of the field that failed validation.
	Field string

	// Message describes why validation failed.
	Message string
}

// Error returns a human-readable representation of the validation error.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("evilmail: validation error on field %q: %s", e.Field, e.Message)
}

// IsNotFoundError reports whether err represents a 404 Not Found API response.
// It works with both *APIError and *AuthError (which wraps *APIError).
func IsNotFoundError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.IsNotFound()
	}
	return false
}

// IsAuthError reports whether err represents an authentication failure
// (401 Unauthorized or 403 Forbidden).
func IsAuthError(err error) bool {
	var authErr *AuthError
	return errors.As(err, &authErr)
}

// IsRateLimitError reports whether err represents a 429 Too Many Requests response.
func IsRateLimitError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.IsRateLimited()
	}
	return false
}
