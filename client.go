package sarvam

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	defaultBaseURL = "https://api.sarvam.ai"

	// AuthModeSubscriptionKey authenticates requests with the
	// api-subscription-key header.
	AuthModeSubscriptionKey AuthMode = "subscription_key"

	// AuthModeBearer authenticates requests with the Authorization header.
	AuthModeBearer AuthMode = "bearer"
)

type errorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// Client communicates with the Sarvam AI API.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	authMode   AuthMode

	// Chat provides access to chat completion endpoints.
	Chat *ChatService
}

// AuthMode identifies the authentication scheme used for API requests.
type AuthMode string

// ClientOption configures a Client.
type ClientOption func(*Client)

// WithBaseURL configures the client to use baseURL for API requests.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = strings.Trim(baseURL, "/")
	}
}

// WithHTTPClient configures the client to use httpClient for API requests.
// A nil client is ignored.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// WithBearerAuth configures the client to send the API key as a bearer token.
func WithBearerAuth() ClientOption {
	return func(c *Client) {
		c.authMode = AuthModeBearer
	}
}

// WithSubscriptionKeyAuth configures the client to send the API key using
// Sarvam's api-subscription-key header.
func WithSubscriptionKeyAuth() ClientOption {
	return func(c *Client) {
		c.authMode = AuthModeSubscriptionKey
	}
}

// ClientOptions contains client configuration values.
//
// ClientOptions is provided as a configuration structure for callers that
// prefer struct-based configuration. NewClient accepts ClientOption values.
type ClientOptions struct {
	// APIKey is the Sarvam API key.
	APIKey string
	// BaseURL is the base URL used for API requests.
	BaseURL string
	// HTTPClient is the HTTP client used for API requests.
	HTTPClient *http.Client
	// AuthMode selects the authentication scheme.
	AuthMode AuthMode
}

// NewClient creates a Sarvam API client using apiKey and the supplied options.
// By default, it uses Sarvam's production API URL, http.DefaultClient, and
// subscription-key authentication.
func NewClient(apiKey string, options ...ClientOption) *Client {
	client := &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: http.DefaultClient,
		authMode:   AuthModeSubscriptionKey,
	}

	for _, option := range options {
		option(client)
	}

	client.Chat = &ChatService{
		Completions: &ChatCompletionService{
			client: client,
		},
	}

	return client
}

func (c *Client) do(
	ctx context.Context,
	method string,
	path string,
	body any,
	result any,
) error {
	var reader io.Reader

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("sarvam: marshal request: %w", err)
		}

		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		c.baseURL+path,
		reader,
	)
	if err != nil {
		return fmt.Errorf("sarvam: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	switch c.authMode {
	case AuthModeBearer:
		req.Header.Set(
			"Authorization",
			"Bearer "+c.apiKey,
		)

	default:
		req.Header.Set(
			"api-subscription-key",
			c.apiKey,
		)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sarvam: request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {
		return newAPIError(resp)
	}

	if result == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf(
			"sarvam: decode response: %w",
			err,
		)
	}

	return nil
}
