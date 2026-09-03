package sarvam

// Message is a single message in a chat completion conversation.
type Message struct {
	// Role identifies the message author, such as "user" or "assistant".
	Role string `json:"role"`
	// Content is the text of the message.
	Content string `json:"content"`
}

const (
	// ModelSarvam105B identifies the Sarvam 105B model.
	ModelSarvam105B = "sarvam-105b"
	// ModelSarvam105BConversations identifies the conversational Sarvam 105B model.
	ModelSarvam105BConversations = "sarvam-105b-conversations"
)

// ChatCompletionRequest contains the parameters for a chat completion.
type ChatCompletionRequest struct {
	// Model is the model identifier to use.
	Model string `json:"model"`
	// Messages contains the conversation history.
	Messages []Message `json:"messages"`
	// Temperature controls the randomness of the generated response.
	Temperature *float64 `json:"temperature,omitempty"`
	// TopP controls nucleus sampling for the generated response.
	TopP *float64 `json:"top_p,omitempty"`
	// MaxTokens limits the number of generated tokens.
	MaxTokens *int `json:"max_tokens,omitempty"`
	// Stream requests a streamed response when supported.
	Stream bool `json:"stream,omitempty"`
	// Stop contains sequences that stop generation.
	Stop []string `json:"stop,omitempty"`
	// N specifies the number of completions to generate.
	N *int `json:"n,omitempty"`
	// Seed makes sampling deterministic when supported.
	Seed *int `json:"seed,omitempty"`
	// ReasoningEffort controls the reasoning effort when supported by the model.
	ReasoningEffort string `json:"reasoning_effort,omitempty"`
}

// ChatCompletionResponse is the response returned by the chat completions API.
type ChatCompletionResponse struct {
	// ID is the unique identifier for the completion.
	ID string `json:"id"`
	// Object identifies the response object type.
	Object string `json:"object"`
	// Created is the Unix timestamp when the completion was created.
	Created int64 `json:"created"`
	// Model is the model that generated the response.
	Model string `json:"model"`
	// Choices contains the generated completion choices.
	Choices []Choice `json:"choices"`
	// Usage contains token usage details, when provided.
	Usage *Usage `json:"usage,omitempty"`
	// ServiceTier identifies the service tier used, when provided.
	ServiceTier *string `json:"service_tier,omitempty"`
	// SystemFingerprint identifies the backend configuration, when provided.
	SystemFingerprint *string `json:"system_fingerprint,omitempty"`
}

// Choice is one generated completion choice.
type Choice struct {
	// Index is the zero-based choice index.
	Index int `json:"index"`
	// Message is the generated assistant message.
	Message Message `json:"message"`
	// FinishReason explains why generation stopped.
	FinishReason string `json:"finish_reason"`
}

// Usage contains token usage statistics for a completion.
type Usage struct {
	// PromptTokens is the number of tokens in the prompt.
	PromptTokens int `json:"prompt_tokens"`
	// CompletionTokens is the number of tokens generated.
	CompletionTokens int `json:"completion_tokens"`
	// TotalTokens is the combined prompt and completion token count.
	TotalTokens int `json:"total_tokens"`
}
