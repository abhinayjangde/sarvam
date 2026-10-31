package sarvam

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ChatCompletionStream reads server-sent events from a streamed chat
// completion. Call Close when the stream is no longer needed.
type ChatCompletionStream struct {
	body    io.ReadCloser
	scanner *bufio.Scanner
	closed  bool
	done    bool
}

func newChatCompletionStream(resp *http.Response) *ChatCompletionStream {
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 4096), 1024*1024)

	return &ChatCompletionStream{
		body:    resp.Body,
		scanner: scanner,
	}
}

// Recv returns the next streamed completion chunk. It returns io.EOF after the
// server sends the SSE [DONE] marker or closes the response.
func (s *ChatCompletionStream) Recv() (*ChatCompletionChunk, error) {
	if s.closed || s.done {
		return nil, io.EOF
	}

	data, err := s.nextEvent()
	if err != nil {
		if err == io.EOF {
			s.done = true
		}
		return nil, err
	}

	if data == "[DONE]" {
		s.done = true
		return nil, io.EOF
	}

	var chunk ChatCompletionChunk
	if err := json.Unmarshal([]byte(data), &chunk); err != nil {
		return nil, fmt.Errorf("sarvam: decode stream event: %w", err)
	}

	return &chunk, nil
}

func (s *ChatCompletionStream) nextEvent() (string, error) {
	var data []string

	for s.scanner.Scan() {
		line := strings.TrimSuffix(s.scanner.Text(), "\r")
		if line == "" {
			if len(data) > 0 {
				return strings.Join(data, "\n"), nil
			}
			continue
		}

		if strings.HasPrefix(line, ":") || strings.HasPrefix(line, "event:") {
			continue
		}

		if strings.HasPrefix(line, "data:") {
			data = append(data, strings.TrimSpace(strings.TrimPrefix(line, "data:")))
		}
	}

	if len(data) > 0 {
		return strings.Join(data, "\n"), nil
	}

	if err := s.scanner.Err(); err != nil {
		return "", fmt.Errorf("sarvam: read stream: %w", err)
	}

	return "", io.EOF
}

// Close closes the underlying HTTP response body.
func (s *ChatCompletionStream) Close() error {
	if s.closed {
		return nil
	}

	s.closed = true
	s.done = true
	return s.body.Close()
}

var _ io.Closer = (*ChatCompletionStream)(nil)
