package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the application configuration
type Config struct {
	APIKey       string `json:"api_key"`
	DefaultModel string `json:"default_model"`
}

// LoadConfig loads the configuration from the config file
// If the config file doesn't exist, it creates one with default values
func LoadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create config directory if it doesn't exist
		configDir := filepath.Dir(configPath)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Config file not found. Let's create one!\n\n")

		// Prompt for API key (required)
		apiKey := promptForAPIKey()

		// Prompt for default model (optional)
		defaultModel := promptForDefaultModel()

		// Create config with user input
		newConfig := &Config{
			APIKey:       apiKey,
			DefaultModel: defaultModel,
		}

		// Write config to file
		if err := saveConfig(configPath, newConfig); err != nil {
			return nil, err
		}

		fmt.Fprintf(os.Stderr, "\n✓ Config file created at: %s\n", configPath)
		return newConfig, nil
	}

	// Read existing config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "llmp", "config.json"), nil
}

// saveConfig saves the config to the specified path
func saveConfig(path string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// promptForAPIKey prompts the user for an API key and validates it's non-empty
func promptForAPIKey() string {
	// Open /dev/tty to read from terminal directly (not stdin)
	tty, err := os.Open("/dev/tty")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not open terminal for input: %v\n", err)
		fmt.Fprintf(os.Stderr, "Please manually edit the config file at ~/.config/llmp/config.json\n")
		return ""
	}
	defer tty.Close()

	reader := bufio.NewReader(tty)

	for {
		fmt.Fprint(os.Stderr, "Enter your OpenRouter API key: ")
		apiKey, _ := reader.ReadString('\n')
		apiKey = strings.TrimSpace(apiKey)

		if apiKey != "" {
			return apiKey
		}

		fmt.Fprintf(os.Stderr, "⚠ API key cannot be empty. Please try again.\n")
	}
}

// promptForDefaultModel prompts the user for a default model with a sane fallback
func promptForDefaultModel() string {
	const defaultModel = "openai/gpt-3.5-turbo"

	// Open /dev/tty to read from terminal directly (not stdin)
	tty, err := os.Open("/dev/tty")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Using default model: %s\n", defaultModel)
		return defaultModel
	}
	defer tty.Close()

	reader := bufio.NewReader(tty)
	fmt.Fprintf(os.Stderr, "Enter default model (press Enter for '%s'): ", defaultModel)

	model, _ := reader.ReadString('\n')
	model = strings.TrimSpace(model)

	if model == "" {
		fmt.Fprintf(os.Stderr, "Using default model: %s\n", defaultModel)
		return defaultModel
	}

	return model
}
