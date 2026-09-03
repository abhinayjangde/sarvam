package sarvam

import (
	"context"
	"strconv"
	"strings"
)

type ChatService struct {
	// Completions provides access to chat completion requests.
	Completions *ChatCompletionService
}

// ChatCompletionService provides operations for the chat completions API.
type ChatCompletionService struct {
	client *Client
}

// Create sends a chat completion request and returns the generated response.
// It returns a ValidationError when request fields are invalid and an APIError
// when Sarvam rejects the request.
func (s *ChatCompletionService) Create(
	ctx context.Context,
	request ChatCompletionRequest,
) (*ChatCompletionResponse, error) {
	if err := request.validate(); err != nil {
		return nil, err
	}

	var response ChatCompletionResponse

	err := s.client.do(
		ctx,
		"POST",
		"/v1/chat/completions",
		request,
		&response,
	)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (r ChatCompletionRequest) validate() error {
	if strings.TrimSpace(r.Model) == "" {
		return &ValidationError{
			Field:   "model",
			Message: "model is required",
		}
	}

	if len(r.Messages) == 0 {
		return &ValidationError{
			Field:   "messages",
			Message: "at least one message is required",
		}
	}

	for i, message := range r.Messages {
		if strings.TrimSpace(message.Role) == "" {
			return &ValidationError{
				Field: "messages[" +
					strconv.Itoa(i) +
					"].role",
				Message: "role is required",
			}
		}
	}

	return nil
}
