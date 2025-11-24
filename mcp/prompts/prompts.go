package prompts

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// SetupPrompts registers the prompts with the server.
func SetupPrompts(server *mcp.Server) {
	server.AddPrompt(&mcp.Prompt{
		Name:  "hello_world",
		Title: "Hello World",
	}, helloWorld)

	server.AddPrompt(&mcp.Prompt{
		Name:  "get_services_by_team",
		Title: "Get Services By Team",
		Arguments: []*mcp.PromptArgument{
			{Name: "team_id", Description: "Team ID", Required: true},
		},
	}, getServicesByTeam)
}

// helloWorld constructs a prompt that guides the LLM to say hello to a user.
func helloWorld(_ context.Context, _ *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
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

// getServicesByTeam constructs a prompt that guides the LLM to retrieve services for a specific team.
//
// It expects a 'team_id' argument in the request.
// The prompt instructs the LLM to either use the `servicemap` resource or the `get_services_by_team` tool.
func getServicesByTeam(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	teamId := req.Params.Arguments["team_id"]
	return &mcp.GetPromptResult{
		Description: "Get services by team",
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{Text: fmt.Sprintf("To get a list of services for a team, use the resource Use the resource: servicemap://teams/%s/services` or the `get_services_by_team` tool, passing in a required 'team' parameter."+
					"It will return a list of services.", teamId)},
			},
		},
	}, nil
}
