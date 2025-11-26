package teamrepository

import (
	"context"
	"net/http"
	"service-atlas/internal/customerrors"
	"service-atlas/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (r Neo4jTeamRepository) CreateTeam(ctx context.Context, team repositories.Team) (string, error) {
	createTeamTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(
			ctx, `
        CREATE (n: Team {id: randomuuid(), created: datetime(), updated: datetime(), name: $name})
        RETURN n.id AS id
        `, map[string]any{
				"name": team.Name,
			})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			id, ok := result.Record().Get("id")
			if !ok {
				return nil, &customerrors.HTTPError{
					Status: http.StatusInternalServerError,
					Msg:    "Id not returned when creating team",
				}
			}
			return id, nil
		}
		return nil, &customerrors.HTTPError{
			Status: http.StatusInternalServerError,
			Msg:    "No id returned from creating team",
		}

	}
	result, err := r.manager.ExecuteWrite(ctx, createTeamTransaction)
	if err != nil {
		return "", &customerrors.HTTPError{
			Status: http.StatusInternalServerError,
			Msg:    "Error creating team",
		}
	}
	id, ok := result.(string)
	if !ok {
		return "", &customerrors.HTTPError{
			Status: http.StatusInternalServerError,
			Msg:    "Error creating team",
		}
	}
	return id, nil
}
