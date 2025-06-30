package prompts

import (
	"context"
	"embed"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

//go:embed *.json
var promptFiles embed.FS

type ToolPrompt struct {
	Description string `json:"description"`
	Prompt      string `json:"prompt"`
}

type ToolResult map[string]interface{}

func RegisterTools(mcpServer *server.MCPServer, promptsPath string) {
	toolFiles := make(map[string][]byte)

	embeddedEntries, err := promptFiles.ReadDir(".")
	if err != nil {
		log.Printf("[MCP][ERROR] reading embedded prompts: %v", err)
	} else {
		for _, entry := range embeddedEntries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
				continue
			}
			data, err := promptFiles.ReadFile(entry.Name())
			if err != nil {
				log.Printf("[MCP][ERROR] reading embedded prompt file %s: %v", entry.Name(), err)
				continue
			}
			toolFiles[entry.Name()] = data
		}
	}

	if promptsPath != "" {
		userEntries, err := os.ReadDir(promptsPath)
		if err != nil {
			log.Printf("[MCP][ERROR] reading prompts from %s: %v", promptsPath, err)
		} else {
			for _, entry := range userEntries {
				if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
					continue
				}
				data, err := os.ReadFile(promptsPath + "/" + entry.Name())
				if err != nil {
					log.Printf("[MCP][ERROR] reading prompt file %s: %v", entry.Name(), err)
					continue
				}
				toolFiles[entry.Name()] = data
			}
		}
	}

	for fileName, data := range toolFiles {
		toolName := strings.TrimSuffix(fileName, ".json")
		var tp ToolPrompt
		if err := json.Unmarshal(data, &tp); err != nil {
			log.Printf("[MCP][ERROR] parsing prompt file %s: %v", fileName, err)
			continue
		}
		mcpServer.AddTool(
			mcp.NewTool(toolName,
				mcp.WithDescription(tp.Description),
				mcp.WithString("query", mcp.Description("User query for the tool"), mcp.Required()),
			),
			makeHandlePromptTool(tp.Prompt),
		)
	}
}

func makeHandlePromptTool(prompt string) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var params struct {
			Query string `json:"query"`
		}
		if err := decodeParams(request.Params.Arguments, &params); err != nil {
			return nil, err
		}
		finalPrompt := prompt + "\n---\n" + params.Query
		return wrapResult(ToolResult{"prompt": finalPrompt}), nil
	}
}

func decodeParams(args interface{}, out interface{}) error {
	b, err := json.Marshal(args)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}

func wrapResult(res ToolResult) *mcp.CallToolResult {
	isError := false
	if v, ok := res["isError"]; ok {
		if b, ok := v.(bool); ok && b {
			isError = true
		}
	}
	b, _ := json.Marshal(res)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(b),
			},
		},
		IsError: isError,
	}
}
