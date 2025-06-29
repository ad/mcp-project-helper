package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ad/mcp-project-helper/prompts"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	var transport = flag.String("transport", "stdio", "Transport type: stdio, sse, or http")
	var port = flag.String("port", "8080", "Port for SSE/HTTP servers")
	var promptsPath = flag.String("prompts-path", "", "Path to custom prompts directory (optional)")
	flag.Parse()

	usageDefinition := "Usage: %s [-transport stdio|sse|http] [-port PORT] [-prompts-path PATH]\n"

	mcpServer := server.NewMCPServer(
		"helper",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	prompts.RegisterTools(mcpServer, *promptsPath)

	switch *transport {
	case "stdio":
		log.Println("Starting MCP server with STDIO transport...")
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatal("STDIO server error:", err)
		}

	case "sse":
		log.Printf("Starting MCP server with SSE transport on port %s...", *port)
		sseServer := server.NewSSEServer(mcpServer)

		http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
			sseServer.ServeHTTP(w, r)
		})

		if err := http.ListenAndServe(":"+*port, nil); err != nil {
			log.Fatal("SSE server error:", err)
		}

	case "http":
		log.Printf("Starting MCP server with streamable HTTP transport on port %s...", *port)
		httpServer := server.NewStreamableHTTPServer(mcpServer)

		log.Printf("HTTP server listening on :%s/mcp", *port)
		if err := httpServer.Start(":" + *port); err != nil {
			log.Fatal("HTTP server error:", err)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown transport: %s\n", *transport)
		fmt.Fprintf(os.Stderr, usageDefinition, os.Args[0])
		os.Exit(1)
	}
}
