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
	defaultBaseURL                   = "https://api.sarvam.ai"
	AuthModeSubscriptionKey AuthMode = "subscription_key"
	AuthModeBearer          AuthMode = "bearer"
)

type errorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	authMode   AuthMode

	Chat *ChatService
}

type AuthMode string

type ClientOption func(*Client)

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = strings.Trim(baseURL, "/")
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

func WithBearerAuth() ClientOption {
	return func(c *Client) {
		c.authMode = AuthModeBearer
	}
}

func WithSubscriptionKeyAuth() ClientOption {
	return func(c *Client) {
		c.authMode = AuthModeSubscriptionKey
	}
}

type ClientOptions struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
	AuthMode   AuthMode
}

// new client function
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
