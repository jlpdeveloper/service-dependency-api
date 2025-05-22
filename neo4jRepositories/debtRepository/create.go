package debtRepository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-dependency-api/internal/customErrors"
	"service-dependency-api/repositories"
)

func (n Neo4jDebtRepository) CreateDebtItem(ctx context.Context, debt repositories.Debt) error {
	createDebtTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		// Check if the service exists
		checkQuery := `
			MATCH (s:Service {id: $serviceId})
			RETURN s
		`
		result, err := tx.Run(ctx, checkQuery, map[string]any{
			"serviceId": debt.ServiceId,
		})
		if err != nil {
			return nil, err
		}

		// If no records are returned, the service doesn't exist
		records, err := result.Collect(ctx)
		if err != nil {
			return nil, err
		}
		if len(records) == 0 {
			return nil, &customErrors.HTTPError{
				Status: 404,
				Msg:    fmt.Sprintf("Service not found: %s", debt.ServiceId),
			}
		}
		_, err = tx.Run(ctx, `
				MATCH (s:Service {id: $serviceId})
				CREATE (n:Debt {id: randomuuid(), created: datetime(), title: $title, type: $type, description: $description, status: $status})
				CREATE (s)-[r:OWNS]->(n)
        `, map[string]any{
			"title":       debt.Title,
			"type":        debt.Type,
			"description": debt.Description,
			"status":      DefaultStatus,
			"serviceId":   debt.ServiceId,
		})
		return nil, err
	}
	_, err := n.manager.ExecuteWrite(ctx, createDebtTransaction)
	return err
}
