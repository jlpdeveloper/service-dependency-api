package prompts

import (
	"context"
	"fmt"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestGetServicesByTeamPrompt(t *testing.T) {
	ctx := context.Background()
	r := &mcp.GetPromptRequest{
		Params: &mcp.GetPromptParams{
			Arguments: map[string]string{
				"team_id": "test-team-id",
			},
		},
	}
	result, err := getServicesByTeam(ctx, r)

	if err != nil {
		t.Error(err)
	}
	if result == nil {
		t.Fatal("Expected result")
	}
	if result.Description != "Get services by team" {
		t.Errorf("Expected description to be 'Get services by team', got '%s'", result.Description)
	}
	if len(result.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(result.Messages))
	}
	msg := result.Messages[0]
	if msg.Role != "user" {
		t.Errorf("Expected role 'user', got '%s'", msg.Role)
	}
	c := msg.Content.(*mcp.TextContent)
	expectedContent := fmt.Sprintf("To get a list of services for a team, use the resource Use the resource: servicemap://teams/%s/services` or the `get_services_by_team` tool, passing in a required 'team' parameter."+
		"It will return a list of services. If you haven't previously called the 'get_teams' tool, you will need to do so you can get team Ids based on team names. This call expects a guid id passed in", "test-team-id")
	if c.Text != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, c.Text)
	}

}
