package sarvam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf(
		"sarvam: invalid %s: %s",
		e.Field,
		e.Message,
	)
}

type APIError struct {
	StatusCode int
	Message    string
	Type       string
	Code       string
	RequestID  string
}

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

func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized ||
		e.StatusCode == http.StatusForbidden
}

func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

func (e *APIError) IsClientError() bool {
	return e.StatusCode >= 400 &&
		e.StatusCode < 500
}

func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500
}
