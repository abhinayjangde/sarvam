package sarvam

import "context"

type ChatService struct {
	Completions *ChatCompletionService
}

type ChatCompletionService struct {
	client *Client
}

func (s *ChatCompletionService) Create(
	ctx context.Context,
	request ChatCompletionRequest,
) (*ChatCompletionResponse, error) {

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
