package teamRepository

import (
	"context"
	"service-dependency-api/internal/customErrors"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (r Neo4jTeamRepository) DeleteTeam(ctx context.Context, id string) error {
	deleteTeamTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
    		MATCH (s:Team { id: $id })
    		RETURN count(s) as count
		`, map[string]interface{}{"id": id})

		if err != nil {
			return nil, err
		}

		if record, err := result.Single(ctx); err == nil {
			count, _ := record.Get("count")
			if count.(int64) == 0 {
				return nil, &customErrors.HTTPError{
					Status: 404,
					Msg:    "Team not found",
				}
			}
		}
		result, err = tx.Run(ctx, `
		MATCH(s:Team { id: $id})
		DETACH DELETE s;`, map[string]interface{}{"id": id})
		if err != nil {
			return nil, err
		}

		summary, err := result.Consume(ctx)
		if err != nil {
			return nil, &customErrors.HTTPError{Status: 500, Msg: "Error deleting team: " + id}
		}

		if summary.Counters().NodesDeleted() == 0 {
			return nil, &customErrors.HTTPError{Status: 500, Msg: "Error deleting team: " + id}
		}
		return nil, nil
	}
	_, err := r.manager.ExecuteWrite(ctx, deleteTeamTransaction)
	return err
}
