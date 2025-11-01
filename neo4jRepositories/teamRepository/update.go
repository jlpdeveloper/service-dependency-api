package teamRepository

import (
	"context"
	"net/http"
	"service-dependency-api/internal/customErrors"

	"service-dependency-api/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (r Neo4jTeamRepository) UpdateTeam(ctx context.Context, team repositories.Team) error {
	_, err := r.GetTeam(ctx, team.Id)
	// Error should be a custom http error already
	if err != nil {
		return err
	}

	updateTeamTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		updateResult, updateErr := tx.Run(ctx, `
			MATCH (s:Team)
			WHERE s.id = $id
			SET s.name = $name,
				s.updated = datetime()
			RETURN s
		`, map[string]any{
			"id":   team.Id,
			"name": team.Name,
		})

		if updateErr != nil {
			return nil, customErrors.HTTPError{
				Status: http.StatusInternalServerError,
				Msg:    err.Error(),
			}
		}

		// Confirm update was successful
		if !updateResult.Next(ctx) {
			err = customErrors.HTTPError{
				Status: http.StatusInternalServerError,
				Msg:    "update Team failed",
			}
		}
		return nil, nil
	}

	_, err = r.manager.ExecuteWrite(ctx, updateTeamTransaction)
	if err != nil {
		return err
	}
	return nil
}
