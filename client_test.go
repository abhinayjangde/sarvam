package sarvam

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-key")

	if client == nil {
		t.Fatal("expected client")
	}

	if client.baseURL != defaultBaseURL {
		t.Fatalf(
			"expected %s, got %s",
			defaultBaseURL,
			client.baseURL,
		)
	}

	if client.httpClient == nil {
		t.Fatal("expected HTTP client")
	}

	if client.Chat == nil {
		t.Fatal("expected Chat service")
	}

	if client.Chat.Completions == nil {
		t.Fatal("expected Completions service")
	}
}

func TestWithBaseURL(t *testing.T) {
	client := NewClient(
		"test-key",
		WithBaseURL("https://example.com/"),
	)

	if client.baseURL != "https://example.com" {
		t.Fatalf(
			"unexpected base URL: %s",
			client.baseURL,
		)
	}
}

func TestBearerAuth(t *testing.T) {
	client := NewClient(
		"test-key",
		WithBearerAuth(),
	)

	if client.authMode != AuthModeBearer {
		t.Fatal("expected bearer authentication")
	}
}

func TestSubscriptionKeyAuth(t *testing.T) {
	client := NewClient(
		"test-key",
		WithSubscriptionKeyAuth(),
	)

	if client.authMode != AuthModeSubscriptionKey {
		t.Fatal("expected subscription-key authentication")
	}
}
