package sarvam

import "fmt"

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf(
			"sarvam: API error (%d): %s",
			e.StatusCode,
			e.Message,
		)
	}

	return fmt.Sprintf(
		"sarvam: API error (%d)",
		e.StatusCode,
	)
}
