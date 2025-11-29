package dependencyrepository

import (
	"context"
	"fmt"
	"service-atlas/internal/customerrors"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (d *Neo4jDependencyRepository) DeleteDependency(ctx context.Context, id string, dependsOnID string) error {
	deleteDependencyTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		// Check if both services exist and the dependency relationship exists
		checkQuery := `
			MATCH (s1:Service {id: $serviceId})-[r:DEPENDS_ON]->(s2:Service {id: $dependsOnID})
			RETURN r
		`
		result, err := tx.Run(ctx, checkQuery, map[string]any{
			"serviceId":   id,
			"dependsOnID": dependsOnID,
		})
		if err != nil {
			return nil, err
		}

		// If no records are returned, the dependency relationship doesn't exist
		records, err := result.Collect(ctx)
		if err != nil {
			return nil, err
		}
		if len(records) == 0 {
			return nil, &customerrors.HTTPError{
				Status: 404,
				Msg:    fmt.Sprintf("Dependency relationship not found between services: %s -> %s", id, dependsOnID),
			}
		}

		// Delete the dependency relationship
		deleteQuery := `
			MATCH (s1:Service {id: $serviceId})-[r:DEPENDS_ON]->(s2:Service {id: $dependsOnID})
			DELETE r
		`
		_, err = tx.Run(ctx, deleteQuery, map[string]any{
			"serviceId":   id,
			"dependsOnID": dependsOnID,
		})
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	_, err := d.manager.ExecuteWrite(ctx, deleteDependencyTransaction)
	return err
}
