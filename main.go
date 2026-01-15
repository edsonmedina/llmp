package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// Define command-line flags
	modelFlag := flag.String("model", "", "OpenRouter model to use (uses config default if not set)")
	outputFlag := flag.String("o", "", "Output file path (writes to stdout if not set)")
	webSearchFlag := flag.Bool("ws", false, "Enable web search")
	flag.Parse()

	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Determine which model to use
	model := *modelFlag
	if model == "" {
		model = config.DefaultModel
	}

	// Read prompt from stdin
	promptBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
		os.Exit(1)
	}

	prompt := string(promptBytes)
	if prompt == "" {
		fmt.Fprintf(os.Stderr, "Error: no prompt provided via stdin\n")
		os.Exit(1)
	}

	// Send prompt to OpenRouter API
	response, err := SendPrompt(config.APIKey, model, prompt, *webSearchFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Write response to output
	if *outputFlag != "" {
		// Write to file
		if err := os.WriteFile(*outputFlag, []byte(response), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Response written to: %s\n", *outputFlag)
	} else {
		// Write to stdout
		fmt.Println(response)
	}
}
