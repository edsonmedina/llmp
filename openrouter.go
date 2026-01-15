package main

import (
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

// SendPrompt sends a prompt to the OpenRouter API and returns the response
func SendPrompt(apiKey, model, systemPrompt, prompt string, enableWebSearch bool) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("API key is not configured. Please add your OpenRouter API key to the config file")
	}

	// Enable web search by appending :online to the model name
	// This is the recommended OpenRouter approach that works with all models
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
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", openRouterAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var apiResp OpenRouterResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w\nResponse body: %s", err, string(body))
	}

	// Check for API error
	if apiResp.Error != nil {
		return "", fmt.Errorf("API error: %s (code: %s)\nFull response: %s", apiResp.Error.Message, string(apiResp.Error.Code), string(body))
	}

	// Extract response content
	if len(apiResp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return apiResp.Choices[0].Message.Content, nil
}
