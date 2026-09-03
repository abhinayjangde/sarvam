package sarvam

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

const (
	ModelSarvam105B              = "sarvam-105b"
	ModelSarvam105BConversations = "sarvam-105b-conversations"
)

type ChatCompletionRequest struct {
	Model           string    `json:"model"`
	Messages        []Message `json:"messages"`
	Temperature     *float64  `json:"temperature,omitempty"`
	TopP            *float64  `json:"top_p,omitempty"`
	MaxTokens       *int      `json:"max_tokens,omitempty"`
	Stream          bool      `json:"stream,omitempty"`
	Stop            []string  `json:"stop,omitempty"`
	N               *int      `json:"n,omitempty"`
	Seed            *int      `json:"seed,omitempty"`
	ReasoningEffort string    `json:"reasoning_effort,omitempty"`
}

type ChatCompletionResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             *Usage   `json:"usage,omitempty"`
	ServiceTier       *string  `json:"service_tier,omitempty"`
	SystemFingerprint *string  `json:"system_fingerprint,omitempty"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
