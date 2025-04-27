package dependencyRepository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-dependency-api/internal/customErrors"
)

func (d *Neo4jDependencyRepository) GetDependencies(ctx context.Context, id string) ([]*Dependency, error) {
	session := d.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() {
		_ = session.Close(ctx)
	}()

	query := `
			MATCH (s1:Service {id: $serviceId})-[r:DEPENDS_ON]->(s2:Service)
			RETURN s2.id as id, s2.name as name, r.version as version
		`
	result, err := session.ExecuteRead(ctx, makeGetTransaction(ctx, id, query))
	if err != nil {
		return nil, err
	}

	return result.([]*Dependency), nil
}

func (d *Neo4jDependencyRepository) GetDependents(ctx context.Context, id string) ([]*Dependency, error) {
	session := d.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() {
		_ = session.Close(ctx)
	}()

	query := `
			MATCH (s1:Service)-[r:DEPENDS_ON]->(s2:Service {id: $serviceId})
			RETURN s1.id as id, s1.name as name, r.version as version
		`
	result, err := session.ExecuteRead(ctx, makeGetTransaction(ctx, id, query))
	if err != nil {
		return nil, err
	}

	return result.([]*Dependency), nil
}

func makeGetTransaction(ctx context.Context, id string, query string) func(tx neo4j.ManagedTransaction) (any, error) {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		// First check if the service exists
		checkQuery := `
			MATCH (s:Service {id: $serviceId})
			RETURN s
		`
		result, err := tx.Run(ctx, checkQuery, map[string]any{
			"serviceId": id,
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
				Msg:    fmt.Sprintf("Service not found: %s", id),
			}
		}

		// Find all services that depend on the service with the given ID

		result, err = tx.Run(ctx, query, map[string]any{
			"serviceId": id,
		})
		if err != nil {
			return nil, err
		}

		dependencies := []*Dependency{}
		records, err = result.Collect(ctx)
		if err != nil {
			return nil, err
		}

		for _, record := range records {
			id, _ := record.Get("id")
			name, _ := record.Get("name")
			version, _ := record.Get("version")

			dependency := &Dependency{
				Id: id.(string),
			}

			// Only set name and version if they exist
			if name != nil {
				dependency.Name = name.(string)
			}
			if version != nil {
				dependency.Version = version.(string)
			}

			dependencies = append(dependencies, dependency)
		}

		return dependencies, nil
	}
}
