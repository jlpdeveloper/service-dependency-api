package teamrepository

import (
	"context"
	"net/http"
	"service-atlas/internal/customerrors"

	"service-atlas/repositories"

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
			return nil, customerrors.HTTPError{
				Status: http.StatusInternalServerError,
				Msg:    updateErr.Error(),
			}
		}

		// Confirm update was successful
		if !updateResult.Next(ctx) {
			if resultErr := updateResult.Err(); resultErr != nil {
				return nil, customerrors.HTTPError{
					Status: http.StatusInternalServerError,
					Msg:    resultErr.Error(),
				}
			}
			return nil, customerrors.HTTPError{
				Status: http.StatusInternalServerError,
				Msg:    "Failed to confirm update",
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
