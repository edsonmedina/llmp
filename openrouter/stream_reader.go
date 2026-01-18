package openrouter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// streamReader adapts the OpenRouter SSE stream to io.ReadCloser.
//
// The OpenRouter API (compatible with OpenAI) uses Server-Sent Events (SSE) for streaming.
// The raw stream consists of lines starting with "data: ", followed by a JSON object
// representing the chunk, or "[DONE]" to indicate the end of the stream.
//
// Example raw stream:
//
//	data: {"choices": [{"delta": {"content": "Hello"}}]}
//	data: {"choices": [{"delta": {"content": " world"}}]}
//	data: [DONE]
//
// If we read the raw body directly, we would get this internal protocol data.
// This adapter parses these events on the fly, extracting just the "content"
// text (e.g., "Hello world") so that the consumer sees a clean stream of bytes.
type streamReader struct {
	response *http.Response
	reader   *bufio.Reader
	pending  []byte
	err      error
}

// Read implements io.Reader
func (s *streamReader) Read(p []byte) (n int, err error) {
	if s.err != nil {
		return 0, s.err
	}

	// Serve pending data first
	if len(s.pending) > 0 {
		n = copy(p, s.pending)
		s.pending = s.pending[n:]
		return n, nil
	}

	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			s.err = err
			if err == io.EOF {
				return 0, io.EOF
			}
			return 0, err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			s.err = io.EOF
			return 0, io.EOF
		}

		var streamResp struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
			Error *struct {
				Message string `json:"message"`
			} `json:"error,omitempty"`
		}

		if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
			// Skip malformed JSON
			continue
		}

		if streamResp.Error != nil {
			s.err = fmt.Errorf("API Error: %s", streamResp.Error.Message)
			return 0, s.err
		}

		if len(streamResp.Choices) > 0 {
			content := streamResp.Choices[0].Delta.Content
			if content != "" {
				// We found content. Copy it to p.
				contentBytes := []byte(content)
				n = copy(p, contentBytes)
				if n < len(contentBytes) {
					// We couldn't fit everything. Save the rest.
					s.pending = contentBytes[n:]
				}
				return n, nil
			}
		}
	}
}

// Close implements io.Closer
func (s *streamReader) Close() error {
	return s.response.Body.Close()
}
