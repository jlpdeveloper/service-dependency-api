package teamRepository

import (
	"context"
	"net/http"
	"service-dependency-api/internal/customErrors"
	"service-dependency-api/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (r Neo4jTeamRepository) CreateTeam(ctx context.Context, team repositories.Team) error {
	createTeamTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(
			ctx, `
        CREATE (n: Team {id: randomuuid(), created: datetime(), updated: datetime(), name: $name})
        RETURN n.id AS id
        `, map[string]any{
				"name": team.Name,
			})
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
	_, err := r.manager.ExecuteWrite(ctx, createTeamTransaction)

	if err != nil {
		return &customErrors.HTTPError{
			Status: http.StatusInternalServerError,
			Msg:    "Error creating team",
		}
	}
	return nil
}
