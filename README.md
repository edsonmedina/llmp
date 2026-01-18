# llmp - OpenRouter CLI Tool

A lightweight command-line tool for interacting with OpenRouter's LLM API. Built with Go using only the standard library.

## Features

- 🔑 Simple configuration management with auto-generated config file
- 🤖 Support for any OpenRouter model
- ⚡ Real-time streaming response
- 🌐 Optional web search capability
- 📝 Output to stdout or file
- 🚀 Zero external dependencies

## Installation

```bash
cd /home/edson/lab/llmp
go build -o llmp
```

Optionally, move the binary to your PATH:
```bash
sudo mv llmp /usr/local/bin/
```

## Configuration

On first run, `llmp` will create a configuration file at `~/.config/llmp/config.json`:

```json
{
  "api_key": "",
  "default_model": "openai/gpt-3.5-turbo"
}
```

Add your OpenRouter API key to this file:

```json
{
  "api_key": "sk-or-v1-YOUR-API-KEY-HERE",
  "default_model": "openai/gpt-3.5-turbo"
}
```

You can also change the default model to any model supported by OpenRouter.

## Usage

### Basic usage (stdout)

```bash
echo "Explain quantum computing in one sentence" | llmp
```

### Using a specific model

```bash
echo "What is the capital of France?" | llmp -model anthropic/claude-3-sonnet
```

### Enable web search

```bash
echo "What's the latest news about AI?" | llmp -ws
```

### Output to file

```bash
echo "Write a haiku about programming" | llmp -o output.txt
```

### Combining options

```bash
echo "What's the weather in Tokyo today?" | llmp -model openai/gpt-4 -ws -o weather.txt
```

### Using heredoc for multi-line prompts

```bash
llmp << EOF
Write a short poem about:
- The ocean
- Sunset
- Peace
EOF
```

### Using input files

Read user prompt from a file:
```bash
llmp -u prompt.txt
```

Provide a system prompt from a file:
```bash
echo "Hola" | llmp -s system_prompt.txt
```

Combine both:
```bash
llmp -s system.txt -u prompt.txt
```

## Command-line Flags

- `-model <string>` - Specify the OpenRouter model to use (defaults to config value)
- `-o <file>` - Write output to a file instead of stdout
- `-ws` - Enable web search capability
- `-u <file>` - Read user prompt from a file
- `-s <file>` - Read system prompt from a file
- `-h` - Show help message

## Examples

**Quick question:**
```bash
echo "What is 2+2?" | llmp
```

**Code generation:**
```bash
echo "Write a Python function to calculate fibonacci numbers" | llmp -model anthropic/claude-3-opus
```

**Research with web search:**
```bash
echo "What are the latest developments in quantum computing?" | llmp -ws -o research.txt
```

## Supported Models

`llmp` supports any model available on OpenRouter. Popular options include:

- `openai/gpt-4`
- `openai/gpt-3.5-turbo`
- `anthropic/claude-3-opus`
- `anthropic/claude-3-sonnet`
- `google/gemini-pro`
- `meta-llama/llama-3-70b-instruct`

See [OpenRouter's model list](https://openrouter.ai/models) for all available models.

## Error Handling

- If the API key is missing, the tool will prompt you to add it to the config file
- API errors are displayed with descriptive messages
- Empty prompts are rejected with an error message

## License

MIT
