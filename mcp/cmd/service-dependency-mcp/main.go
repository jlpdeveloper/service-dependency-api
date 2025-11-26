package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"service-dependency-api/mcp/prompts"
	"service-dependency-api/mcp/tools"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()
	url, err := getApiUrl()
	if err != nil {
		log.Fatal(err)
	}
	apiCtx := context.WithValue(ctx, "API_URL", url)
	server := mcp.NewServer(&mcp.Implementation{Name: "Service Map MCP", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "hello_world", Description: "say hi"}, tools.HelloWorld)
	registerGetServicesByTeam(server)
	prompts.SetupPrompts(server)
	// Run the server over stdin/stdout, until the client disconnects.
	if err := server.Run(apiCtx, &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

func registerGetServicesByTeam(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "get_services_by_team", Description: "get services by team"}, tools.GetServicesByTeamTool)
	server.AddResourceTemplate(&mcp.ResourceTemplate{
		Name:        "Get Services By Team",
		MIMEType:    "application/json",
		URITemplate: "servicemap://teams/{teamId}/services",
	}, tools.GetServicesByTeamResource)
}

func getApiUrl() (string, error) {
	envUrl, ok := os.LookupEnv("API_URL")
	if !ok {
		return "", errors.New("API_URL environment variable not set")
	}
	c := http.Client{
		Timeout: 1 * time.Second,
	}
	if envUrl[len(envUrl)-1:] != "/" {
		envUrl += "/"
	}
	_, err := c.Get(envUrl + "helloworld")
	if err != nil {
		return "", err
	}
	return envUrl, nil
}
