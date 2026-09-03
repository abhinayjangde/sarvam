package sarvam

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChatCompletion(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Method != http.MethodPost {
				t.Fatalf("expected POST, got %s", r.Method)
			}

			if r.URL.Path != "/v1/chat/completions" {
				t.Fatalf(
					"unexpected path: %s",
					r.URL.Path,
				)
			}

			if r.Header.Get("api-subscription-key") != "test-key" {
				t.Fatal("missing API key")
			}

			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			w.WriteHeader(http.StatusOK)

			w.Write([]byte(`{
				"id": "test-id",
				"object": "chat.completion",
				"created": 1735689600,
				"model": "sarvam-105b",
				"choices": [
					{
						"index": 0,
						"message": {
							"role": "assistant",
							"content": "Hello from Sarvam!"
						},
						"finish_reason": "stop"
					}
				],
				"usage": {
					"prompt_tokens": 5,
					"completion_tokens": 4,
					"total_tokens": 9
				}
			}`))
		}),
	)

	defer server.Close()

	client := NewClient("test-key")

	client.baseURL = server.URL

	response, err := client.Chat.Completions.Create(
		context.Background(),
		ChatCompletionRequest{
			Model: ModelSarvam105B,
			Messages: []Message{
				{
					Role:    "user",
					Content: "Hello",
				},
			},
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.ID != "test-id" {
		t.Fatalf(
			"expected test-id, got %s",
			response.ID,
		)
	}

	if response.Choices[0].Message.Content != "Hello from Sarvam!" {
		t.Fatalf("unexpected response")
	}
}

func TestChatCompletionValidation(t *testing.T) {
	client := NewClient("test-key")

	_, err := client.Chat.Completions.Create(
		context.Background(),
		ChatCompletionRequest{},
	)

	if err == nil {
		t.Fatal("expected validation error")
	}

	var validationErr *ValidationError

	if !errors.As(err, &validationErr) {
		t.Fatalf(
			"expected ValidationError, got %T",
			err,
		)
	}

	if validationErr.Field != "model" {
		t.Fatalf(
			"expected model validation error, got %s",
			validationErr.Field,
		)
	}
}

func TestChatCompletionAPIError(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			w.WriteHeader(http.StatusForbidden)

			_, _ = w.Write([]byte(`{
				"error": {
					"message": "Invalid API key",
					"type": "invalid_api_key_error",
					"code": "invalid_api_key"
				}
			}`))
		}),
	)

	defer server.Close()

	client := NewClient(
		"bad-key",
		WithBaseURL(server.URL),
	)

	_, err := client.Chat.Completions.Create(
		context.Background(),
		ChatCompletionRequest{
			Model: ModelSarvam105B,
			Messages: []Message{
				{
					Role:    "user",
					Content: "Hello",
				},
			},
		},
	)

	if err == nil {
		t.Fatal("expected API error")
	}

	var apiErr *APIError

	if !errors.As(err, &apiErr) {
		t.Fatalf(
			"expected APIError, got %T",
			err,
		)
	}

	if apiErr.StatusCode != http.StatusForbidden {
		t.Fatalf(
			"expected 403, got %d",
			apiErr.StatusCode,
		)
	}

	if apiErr.Code != "invalid_api_key" {
		t.Fatalf(
			"unexpected error code: %s",
			apiErr.Code,
		)
	}
}

func TestChatCompletionRequest(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			if r.Method != http.MethodPost {
				t.Fatalf(
					"expected POST, got %s",
					r.Method,
				)
			}

			if r.URL.Path != "/v1/chat/completions" {
				t.Fatalf(
					"unexpected path: %s",
					r.URL.Path,
				)
			}

			if got := r.Header.Get(
				"api-subscription-key",
			); got != "test-key" {
				t.Fatalf(
					"unexpected API key: %s",
					got,
				)
			}

			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			_, _ = w.Write([]byte(`{
				"id": "test-id",
				"object": "chat.completion",
				"created": 1735689600,
				"model": "sarvam-105b",
				"choices": [
					{
						"index": 0,
						"message": {
							"role": "assistant",
							"content": "Namaste!"
						},
						"finish_reason": "stop"
					}
				],
				"usage": {
					"prompt_tokens": 5,
					"completion_tokens": 2,
					"total_tokens": 7
				}
			}`))
		}),
	)

	defer server.Close()

	client := NewClient(
		"test-key",
		WithBaseURL(server.URL),
	)

	response, err := client.Chat.Completions.Create(
		context.Background(),
		ChatCompletionRequest{
			Model: ModelSarvam105B,
			Messages: []Message{
				{
					Role:    "user",
					Content: "Say hello",
				},
			},
		},
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if response.ID != "test-id" {
		t.Fatalf(
			"unexpected ID: %s",
			response.ID,
		)
	}

	if response.Choices[0].Message.Content != "Namaste!" {
		t.Fatalf("unexpected response content")
	}
}

func TestContextCancellation(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			<-r.Context().Done()
		}),
	)

	defer server.Close()

	client := NewClient(
		"test-key",
		WithBaseURL(server.URL),
	)

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	cancel()

	_, err := client.Chat.Completions.Create(
		ctx,
		ChatCompletionRequest{
			Model: ModelSarvam105B,
			Messages: []Message{
				{
					Role:    "user",
					Content: "Hello",
				},
			},
		},
	)

	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}
