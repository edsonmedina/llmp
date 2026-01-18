package openrouter

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const openRouterAPIURL = "https://openrouter.ai/api/v1/chat/completions"

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenRouterRequest represents the request payload for OpenRouter API
type OpenRouterRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream,omitempty"`
}

// OpenRouterResponse represents the response from OpenRouter API
type OpenRouterResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string          `json:"message"`
		Code    json.RawMessage `json:"code"`
	} `json:"error,omitempty"`
}

// SendPrompt sends a prompt to the OpenRouter API and returns the response as a stream
func SendPrompt(apiKey, model, systemPrompt, prompt string, enableWebSearch bool) (io.ReadCloser, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is not configured. Please add your OpenRouter API key to the config file")
	}

	// Enable web search by appending :online to the model name
	if enableWebSearch && !strings.HasSuffix(model, ":online") {
		model = model + ":online"
	}

	// Build messages
	var messages []Message
	if systemPrompt != "" {
		messages = append(messages, Message{
			Role:    "system",
			Content: systemPrompt,
		})
	}
	messages = append(messages, Message{
		Role:    "user",
		Content: prompt,
	})

	// Build request payload
	reqPayload := OpenRouterRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", openRouterAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return &streamReader{
		response: resp,
		reader:   bufio.NewReader(resp.Body),
	}, nil
}
