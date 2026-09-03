package sarvam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ValidationError struct {
	// Field identifies the invalid request field.
	Field string
	// Message describes why the field is invalid.
	Message string
}

// Error returns the validation error message.
func (e *ValidationError) Error() string {
	return fmt.Sprintf(
		"sarvam: invalid %s: %s",
		e.Field,
		e.Message,
	)
}

type APIError struct {
	// StatusCode is the HTTP status code returned by the API.
	StatusCode int
	// Message is the API-provided error message.
	Message string
	// Type is the API-provided error type.
	Type string
	// Code is the API-provided error code.
	Code string
	// RequestID is the request identifier returned by the API, when available.
	RequestID string
}

// Error returns a human-readable API error message.
func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf(
			"sarvam: API error (%d)",
			e.StatusCode,
		)
	}

	return fmt.Sprintf(
		"sarvam: API error (%d): %s",
		e.StatusCode,
		e.Message,
	)
}

type apiErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

func newAPIError(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    http.StatusText(resp.StatusCode),
		}
	}

	var payload apiErrorResponse

	if err := json.Unmarshal(body, &payload); err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    payload.Error.Message,
		Type:       payload.Error.Type,
		Code:       payload.Error.Code,
		RequestID:  resp.Header.Get("x-request-id"),
	}
}

// IsUnauthorized reports whether the API rejected the credentials or access.
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized ||
		e.StatusCode == http.StatusForbidden
}

// IsRateLimited reports whether the API rejected the request due to rate limits.
func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// IsClientError reports whether the API returned a 4xx status code.
func (e *APIError) IsClientError() bool {
	return e.StatusCode >= 400 &&
		e.StatusCode < 500
}

// IsServerError reports whether the API returned a 5xx status code.
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500
}
