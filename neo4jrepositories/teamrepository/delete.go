package teamrepository

import (
	"context"
	"service-atlas/internal/customerrors"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (r Neo4jTeamRepository) DeleteTeam(ctx context.Context, id string) error {
	deleteTeamTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			OPTIONAL MATCH (t:Team { id: $id })
			DETACH DELETE t
			RETURN count(t) AS deletedCount
		`, map[string]any{"id": id})
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}

		deletedVal, ok := record.Get("deletedCount")
		if !ok {
			return nil, &customerrors.HTTPError{Status: 500, Msg: "Error deleting team: " + id}
		}

		var deletedCount int64
		switch v := deletedVal.(type) {
		case int64:
			deletedCount = v
		case int:
			deletedCount = int64(v)
		default:
			return nil, &customerrors.HTTPError{Status: 500, Msg: "Error deleting team: " + id}
		}

		if deletedCount == 0 {
			return nil, &customerrors.HTTPError{Status: 404, Msg: "Team not found"}
		}
		return nil, nil
	}
	_, err := r.manager.ExecuteWrite(ctx, deleteTeamTransaction)
	return err
}
