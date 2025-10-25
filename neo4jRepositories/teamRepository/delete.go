package teamRepository

import (
	"context"
	"errors"
)

func (r Neo4jTeamRepository) DeleteTeam(ctx context.Context, teamId string) error {
	return errors.New("Not implemented")
}
