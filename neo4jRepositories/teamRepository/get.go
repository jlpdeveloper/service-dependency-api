package teamRepository

import (
	"context"
	"errors"
	"net/http"
	"service-dependency-api/internal/customErrors"
	"service-dependency-api/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (r Neo4jTeamRepository) GetTeam(ctx context.Context, teamId string) (*repositories.Team, error) {
	getTeamTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (s:Team)
			WHERE s.id = $id
			RETURN s
		`, map[string]any{
			"id": teamId,
		})

		if err != nil {
			return nil, customErrors.HTTPError{
				Status: http.StatusInternalServerError,
				Msg:    err.Error(),
			}
		}

		if !result.Next(ctx) {
			return nil, customErrors.HTTPError{
				Status: http.StatusNotFound,
				Msg:    "Team not found",
			} // No team found with this ID
		}

		record := result.Record()
		node, ok := record.Get("s")
		if !ok {
			return nil, customErrors.HTTPError{
				Status: http.StatusInternalServerError,
				Msg:    "Error converting result to team",
			}
		}

		n, ok := node.(neo4j.Node)
		if !ok {
			return nil, customErrors.HTTPError{
				Status: http.StatusInternalServerError,
				Msg:    "Error converting result to team",
			}
		}

		return r.mapNodeToTeam(n), nil
	}
	result, err := r.manager.ExecuteRead(ctx, getTeamTransaction)
	if err != nil {
		return nil, err
	}
	team, ok := result.(repositories.Team)
	if !ok {
		return nil, errors.New("error converting result to team")
	}
	return &team, nil

}

func (r Neo4jTeamRepository) GetTeams(ctx context.Context, page, pageSize int) ([]repositories.Team, error) {
	return nil, errors.New("not implemented")
}
