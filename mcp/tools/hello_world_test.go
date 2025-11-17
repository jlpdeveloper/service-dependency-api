package tools

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestHelloWorld(t *testing.T) {
	ctx := context.Background()

	result, output, err := HelloWorld(ctx, new(mcp.CallToolRequest), Input{Name: "Bob"})
	if err != nil {
		t.Error(err)
	}
	if result != nil {
		t.Error("Expected nil result")
	}
	if output.Greeting != "Hi Bob" {
		t.Errorf("Expected greeting to be 'Hi Bob', got '%s'", output.Greeting)
	}
}
