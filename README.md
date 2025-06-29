# mcp-project-helper

A lightweight, extensible MCP (Model Context Protocol) server for running prompt-based tools and file utilities. Designed for easy integration, testing, and extension with custom prompts.

## Features
- **Prompt-based tools**: Easily add new tools by writing simple JSON prompt files.
- **File utilities**: Includes tools for reading, writing, moving, and deleting files and directories.
- **Custom prompts**: Place your own prompt definitions in the `custom_prompts/` directory.
- **Multiple transports**: Supports STDIO, SSE, and HTTP for flexible integration.
- **Extensive tests**: Includes a test script to verify all tool endpoints.

## Getting Started

### Build Locally
```sh
make build-local
```

## ⚡ Quick Start

### Install via go install

To quickly install the latest version from the repository:

```fish
go install github.com/ad/mcp-project-helper@latest
```

The binary will appear in `$GOBIN` or `$HOME/go/bin` (make sure this path is in your `$PATH`).

### 1. Build from source
```fish
# Clone the repository
git clone https://github.com/ad/mcp-project-helper.git
cd mcp-project-helper

go mod tidy

# Local build
make build-local

# Or manually
go build -o mcp-project-helper main.go

# Local build
make build-local

# Or manually
go build -o mcp-project-helper main.go

# Docker build
make build
```

### Run the Server
- STDIO (default):
```sh
./mcp-project-helper
```

- SSE:
  ```sh
  ./mcp-project-helper -transport sse -port 8080
  ```
- HTTP:
  ```sh
  ./mcp-project-helper -transport http -port 8080
  ```

### Run Tests
```sh
./test.sh
```


## 🔌 Integration

### VS Code

```
go install github.com/ad/mcp-project-helper@latest
````

Добавьте в `settings.json`:
```json
{
  "mcp": {
    "servers": {
      "helper": {
        "type": "stdio",
        "command": "/absolute/path/to/project-helper",
        "args": ["-transport", "stdio"]
      }
    }
  }
}
```

### Docker (VS Code)
```json
{
  "mcp": {
    "servers": {
      "helper": {
        "type": "stdio",
        "command": "docker",
        "args": [
          "run", "--rm", "-i",
          "danielapatin/mcp-project-helper:latest",
          "-transport", "stdio"
        ]
      }
    }
  }
}
```

### Claude Desktop
```json
{
  "mcpServers": {
    "helper": {
      "command": "/absolute/path/to/mcp-project-helper",
      "args": ["-transport", "stdio"]
    }
  }
}
```

## Adding Custom Tools
1. Create a JSON file in `custom_prompts/` (see `palette.json` for an example).
2. Each tool must define a `description` and a `prompt` field.
3. The tool will be automatically registered and available via the MCP protocol.

## Example Tools
- **tool-generator**: Generates a tool description and prompt template based on a user query.
- **palette**: Suggests a harmonious color palette for a given color.

## Project Structure
- `main.go` — Main server entry point
- `prompts/` — Built-in prompt tools
- `custom_prompts/` — User-defined prompt tools
- `test.sh` — End-to-end test script
- `Makefile` — Build and run commands

## License
MIT