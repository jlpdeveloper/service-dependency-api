package releaserepository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-atlas/internal/customerrors"
	"service-atlas/repositories"
)

func (r *Neo4jReleaseRepository) CreateRelease(ctx context.Context, release repositories.Release) error {
	createReleaseTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		// Check if the service exists
		checkQuery := `
			MATCH (s:Service {id: $serviceId})
			RETURN s
		`
		result, err := tx.Run(ctx, checkQuery, map[string]any{
			"serviceId": release.ServiceId,
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
			return nil, &customerrors.HTTPError{
				Status: 404,
				Msg:    fmt.Sprintf("Service not found: %s", release.ServiceId),
			}
		}

		params := map[string]any{
			"serviceId":   release.ServiceId,
			"releaseDate": release.ReleaseDate.Format(releaseTimeFormat),
		}
		// Build the Cypher query dynamically
		propertiesString := "releaseDate: datetime($releaseDate)"
		if release.Url != "" {
			propertiesString += ", url: $url"
			params["url"] = release.Url
		}
		if release.Version != "" {
			propertiesString += ", version: $version"
			params["version"] = release.Version
		}

		query := fmt.Sprintf(`
			MATCH (s:Service {id: $serviceId})
			CREATE (r:Release {%s})
			CREATE (s)-[rel:RELEASED]->(r)
			RETURN r
		`, propertiesString)

		_, err = tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	_, err := r.manager.ExecuteWrite(ctx, createReleaseTransaction)
	return err
}
