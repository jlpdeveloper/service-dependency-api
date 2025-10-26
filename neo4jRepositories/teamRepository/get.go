package teamRepository

import (
	"context"
	"errors"
	"service-dependency-api/repositories"
)

func (r Neo4jTeamRepository) GetTeam(ctx context.Context, teamId string) (*repositories.Team, error) {
	return nil, errors.New("Not implemented")
}

func (r Neo4jTeamRepository) GetTeams(ctx context.Context, page, pageSize int) ([]repositories.Team, error) {
	return nil, errors.New("Not implemented")
}
