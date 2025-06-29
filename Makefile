.PHONY: build run test help

# Build production image
build:
	docker build -t danielapatin/mcp-project-helper .

# Run production container
run: build
	docker run --rm -p 8080:8080 danielapatin/mcp-project-helper

# Run tests in container
test:
	docker run --rm -v $(PWD):/app -w /app golang:1.24-alpine go test -v

# Build binary locally
build-local:
	go build -o mcp-project-helper .

# Run locally (STDIO mode - default)
run-local:
	./mcp-project-helper

# Run in stdio mode (explicit)
run-stdio:
	./mcp-project-helper -transport stdio

# Run in SSE mode
run-sse:
	./mcp-project-helper -transport sse -port 8080

# Run in HTTP mode  
run-http:
	./mcp-project-helper -transport http -port 8080

# Show help
help:
	@echo "Available commands:"
	@echo "  build       - Build production Docker image"
	@echo "  run         - Run production container"
	@echo "  test        - Run tests in container"
	@echo "  build-local - Build binary locally"
	@echo "  run-local   - Run locally (STDIO mode)"
	@echo "  run-stdio   - Run in STDIO mode"
	@echo "  run-sse     - Run in SSE mode on port 8080"
	@echo "  run-http    - Run in HTTP mode on port 8080"
	@echo "  help        - Show this help"
