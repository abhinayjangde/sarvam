# Sarvam Go SDK

A lightweight, idiomatic Go client for the [Sarvam AI](https://sarvam.ai) API.

The SDK provides an OpenAI-compatible interface for Sarvam Chat Completions.

## Installation

```bash
go get github.com/abhinayjangde/sarvam
```

## Authentication

Create a Sarvam API key and set it as an environment variable:

```bash
export SARVAM_API_KEY="your-api-key"
```

PowerShell:

```powershell
$env:SARVAM_API_KEY="your-api-key"
```

Never commit your API key to source control.

## Basic Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/abhinayjangde/sarvam"
)

func main() {
	client := sarvam.NewClient(
		os.Getenv("SARVAM_API_KEY"),
	)

	response, err := client.Chat.Completions.Create(
		context.Background(),
		sarvam.ChatCompletionRequest{
			Model: sarvam.ModelSarvam105B,
			Messages: []sarvam.Message{
				{
					Role:    "user",
					Content: "Explain goroutines in simple Hindi.",
				},
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.Choices[0].Message.Content)
}
```

## Custom HTTP Client

You can provide your own `http.Client`:

```go
httpClient := &http.Client{
	Timeout: 60 * time.Second,
}

client := sarvam.NewClient(
	os.Getenv("SARVAM_API_KEY"),
	sarvam.WithHTTPClient(httpClient),
)
```

## Custom Base URL

A custom base URL can be useful for testing or advanced integrations:

```go
client := sarvam.NewClient(
	apiKey,
	sarvam.WithBaseURL("https://api.example.com"),
)
```

## Bearer Authentication

Sarvam also supports Bearer authentication:

```go
client := sarvam.NewClient(
	apiKey,
	sarvam.WithBearerAuth(),
)
```

By default, the SDK uses Sarvam's `api-subscription-key` authentication.

## Supported Models

```go
sarvam.ModelSarvam105B
sarvam.ModelSarvam105BConversations
```

You can also provide a model ID directly:

```go
sarvam.ChatCompletionRequest{
	Model: "sarvam-105b",
	// ...
}
```

## Error Handling

API errors can be inspected using `errors.As`:

```go
response, err := client.Chat.Completions.Create(
	context.Background(),
	request,
)

if err != nil {
	var apiErr *sarvam.APIError

	if errors.As(err, &apiErr) {
		fmt.Println(apiErr.StatusCode)
		fmt.Println(apiErr.Message)
	}

	return
}
```

## Status

This project is currently in early development.

Current support:

* Chat Completions
* OpenAI-compatible request/response structures
* API-key authentication
* Bearer authentication
* Custom HTTP clients
* Custom base URLs
* Context cancellation
* API error handling
* Request validation

Streaming support is planned for a future release.

## License

MIT
