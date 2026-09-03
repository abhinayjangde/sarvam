package sarvam

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const defaultBaseURL = "https://api.sarvam.ai"

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

	Chat *ChatService
}

type ClientOptions struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(apiKey string) *Client {
	client := &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{},
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
	var bodyReader *bytes.Reader

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("sarvam: marshal request: %w", err)
		}

		bodyReader = bytes.NewReader(data)
	} else {
		bodyReader = bytes.NewReader(nil)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		c.baseURL+path,
		bodyReader,
	)
	if err != nil {
		return fmt.Errorf("sarvam: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-subscription-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sarvam: request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr errorResponse

		_ = json.NewDecoder(resp.Body).Decode(&apiErr)

		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    apiErr.Error.Message,
		}
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
