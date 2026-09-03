package sarvam

import (
	"context"
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
