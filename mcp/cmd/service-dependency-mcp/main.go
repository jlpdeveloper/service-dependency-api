package main

import (
	"context"
	"log"
	"service-dependency-api/mcp/prompts"
	"service-dependency-api/mcp/tools"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "Service Dependency API", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "hello_world", Description: "say hi"}, tools.HelloWorld)
	prompts.SetupPrompts(server)
	// Run the server over stdin/stdout, until the client disconnects.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
