package prompts

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func SetupPrompts(server *mcp.Server) {
	server.AddPrompt(&mcp.Prompt{
		Name:  "hello_world",
		Title: "Hello World",
	}, helloWorld)
}

func helloWorld(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "Hi prompt",
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{Text: "To say hello to a user, use the `hello_world` tool, passing in a required 'name' parameter." +
					"It will return a greeting message."},
			},
		},
	}, nil
}
