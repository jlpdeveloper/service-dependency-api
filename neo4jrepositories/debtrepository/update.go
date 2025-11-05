package debtrepository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-dependency-api/internal/customerrors"
)

func (n Neo4jDebtRepository) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := n.manager.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (d:Debt {id: $id})
			SET d.status = $status
			RETURN count(d) as updatedCount
		`, map[string]any{
			"id":     id,
			"status": status,
		})
		if err != nil {
			return nil, err
		}
		if result.Next(ctx) {
			ctr, _ := result.Record().Get("updatedCount")
			if ct, ok := ctr.(int64); ok && ct == 0 {
				return nil, &customerrors.HTTPError{
					Status: 404,
					Msg:    "Debt not found",
				}
			}
		} else if result.Err() != nil {
			return nil, result.Err()
		}
		return nil, nil
	})
	return err
}
