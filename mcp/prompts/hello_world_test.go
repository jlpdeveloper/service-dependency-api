package prompts

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestHelloWorldPrompt(t *testing.T) {
	ctx := context.Background()
	result, err := helloWorld(ctx, new(mcp.GetPromptRequest))

	if err != nil {
		t.Error(err)
	}
	if result == nil {
		t.Fatal("Expected result")
	}
	if result.Description != "Hi prompt" {
		t.Errorf("Expected description to be 'Hi prompt', got '%s'", result.Description)
	}
	if len(result.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(result.Messages))
	}
	msg := result.Messages[0]
	if msg.Role != "user" {
		t.Errorf("Expected role 'user', got '%s'", msg.Role)
	}

}
