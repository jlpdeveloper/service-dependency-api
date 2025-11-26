package tools

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"service-dependency-api/repositories"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TeamIdInput struct {
	TeamId string `json:"teamId"`
}

type ServicesOutput struct {
	Services []repositories.Service `json:"services"`
}

func getServicesByTeam(ctx context.Context, team string) ([]repositories.Service, error) {
	base := ctx.Value("API_URL")
	if base == nil {
		return nil, errors.New("API_URL environment variable not set")
	}
	url, ok := base.(string)
	if !ok {
		return nil, errors.New("API_URL environment variable is not a string")
	}
	url += "teams/" + team + "/services"
	res, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	services := make([]repositories.Service, 0)
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&services)
	return services, err
}

func GetServicesByTeamTool(ctx context.Context, _ *mcp.CallToolRequest, input TeamIdInput) (
	*mcp.CallToolResult,
	ServicesOutput,
	error,
) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	services, err := getServicesByTeam(ctxWithTimeout, input.TeamId)
	if err != nil {
		return nil, ServicesOutput{}, err
	}
	return nil, ServicesOutput{Services: services}, nil
}

func GetServicesByTeamResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	uri := req.Params.URI
	uri = strings.ReplaceAll(uri, "servicemap://teams/", "")
	teamId := strings.ReplaceAll(uri, "/services", "")
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	services, err := getServicesByTeam(ctxWithTimeout, teamId)
	if err != nil {
		return nil, err
	}
	servicesJson, err := json.Marshal(services)
	if err != nil {
		return nil, err
	}
	return &mcp.ReadResourceResult{Contents: []*mcp.ResourceContents{
		{
			MIMEType: "application/json",
			Text:     string(servicesJson),
		},
	}}, nil

}
