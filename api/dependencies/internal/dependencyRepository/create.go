package dependencyRepository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-dependency-api/internal/customErrors"
)

func (d *Neo4jDependencyRepository) AddDependency(ctx context.Context, id string, dependency *Dependency) error {
	session := d.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		_ = session.Close(ctx)
	}()

	createDependencyTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		// Check if both services exist
		checkQuery := `
			MATCH (s1:Service {id: $serviceId})
			MATCH (s2:Service {id: $dependencyId})
			RETURN s1, s2
		`
		result, err := tx.Run(ctx, checkQuery, map[string]any{
			"serviceId":    id,
			"dependencyId": dependency.Id,
		})
		if err != nil {
			return nil, err
		}

		// If no records are returned, one or both services don't exist
		records, err := result.Collect(ctx)
		if err != nil {
			return nil, err
		}
		if len(records) == 0 {
			return nil, &customErrors.HTTPError{
				Status: 404,
				Msg:    fmt.Sprintf("One or both services not found: %s, %s", id, dependency.Id),
			}
		}

		// Create the dependency relationship
		var query string
		var params map[string]any

		if dependency.Version != "" {
			query = `
				MATCH (s1:Service {id: $serviceId})
				MATCH (s2:Service {id: $dependencyId})
				MERGE (s1)-[r:DEPENDS_ON {version: $version}]->(s2)
				RETURN r
			`
			params = map[string]any{
				"serviceId":    id,
				"dependencyId": dependency.Id,
				"version":      dependency.Version,
			}
		} else {
			query = `
				MATCH (s1:Service {id: $serviceId})
				MATCH (s2:Service {id: $dependencyId})
				MERGE (s1)-[r:DEPENDS_ON]->(s2)
				RETURN r
			`
			params = map[string]any{
				"serviceId":    id,
				"dependencyId": dependency.Id,
			}
		}

		_, err = tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	_, err := session.ExecuteWrite(ctx, createDependencyTransaction)
	return err
}
