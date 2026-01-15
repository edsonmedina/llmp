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
	inputFlag := flag.String("u", "", "File to read user prompt from (optional alternative to stdin)")
	systemFlag := flag.String("s", "", "File to read system prompt from")
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

	// Read system prompt from file if provided
	var systemPrompt string
	if *systemFlag != "" {
		data, err := os.ReadFile(*systemFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading system prompt file: %v\n", err)
			os.Exit(1)
		}
		systemPrompt = string(data)
	}

	// Read user prompt
	var prompt string
	if *inputFlag != "" {
		// Read from file
		data, err := os.ReadFile(*inputFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading user prompt file: %v\n", err)
			os.Exit(1)
		}
		prompt = string(data)
	} else {
		// Check if data is being piped to stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// Read from stdin (pipe)
			promptBytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
				os.Exit(1)
			}
			prompt = string(promptBytes)
		}
		// If it is a terminal (CharDevice), we do nothing and prompt stays empty
	}

	if prompt == "" {
		flag.Usage()
		os.Exit(0)
	}

	// Send prompt to OpenRouter API
	response, err := SendPrompt(config.APIKey, model, systemPrompt, prompt, *webSearchFlag)
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
