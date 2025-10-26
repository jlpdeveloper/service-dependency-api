package teamRepository

import (
	"context"
	"errors"

	"service-dependency-api/repositories"
)

func (r Neo4jTeamRepository) UpdateTeam(ctx context.Context, team repositories.Team) error {
	return errors.New("Not implemented")
}
