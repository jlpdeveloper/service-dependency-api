package teamRepository

import (
	"context"
	"net/http"
	"service-dependency-api/internal/customErrors"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (r Neo4jTeamRepository) CreateTeamAssociation(ctx context.Context, teamId, serviceId string) error {
	createTeamAssociationTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (s:Service {id: $serviceId}), (t:Team {id: $teamId})
			CREATE (t) -[r:OWNS]-> (s)
			RETURN r
		`, map[string]any{
			"serviceId": serviceId,
			"teamId":    teamId,
		})
		if err != nil {
			return nil, customErrors.HTTPError{
				Status: http.StatusInternalServerError,
				Msg:    err.Error(),
			}
		}
		if rErr := result.Err(); rErr != nil {
			return nil, customErrors.HTTPError{
				Status: http.StatusInternalServerError,
				Msg:    rErr.Error(),
			}
		}
		if !result.Next(ctx) {
			return nil, customErrors.HTTPError{
				Status: http.StatusNotFound,
				Msg:    "Failed to create team association",
			}
		}

		return nil, nil
	}
	_, err := r.manager.ExecuteWrite(ctx, createTeamAssociationTransaction)
	if err != nil {
		return err
	}
	return nil
}

func (r Neo4jTeamRepository) DeleteTeamAssociation(ctx context.Context, teamId, serviceId string) error {
	deleteTeamAssociationTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (s:Service {id: $serviceId}), (t:Team {id: $teamId})
			MATCH (t) -[r:OWNS]-> (s)
			DELETE r
			RETURN count(r) as deleted
		`, map[string]any{
			"serviceId": serviceId,
			"teamId":    teamId,
		})
		if err != nil {
			return nil, customErrors.HTTPError{
				Status: http.StatusInternalServerError,
				Msg:    err.Error(),
			}
		}
		if rErr := result.Err(); rErr != nil {
			return nil, customErrors.HTTPError{
				Status: http.StatusInternalServerError,
				Msg:    rErr.Error(),
			}
		}
		if !result.Next(ctx) {
			return nil, customErrors.HTTPError{
				Status: http.StatusNotFound,
				Msg:    "Failed to delete team association",
			}
		}
		return nil, nil
	}
	_, err := r.manager.ExecuteWrite(ctx, deleteTeamAssociationTransaction)
	if err != nil {
		return err
	}
	return nil
}
