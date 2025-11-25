package tools

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"service-dependency-api/repositories"
)

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
