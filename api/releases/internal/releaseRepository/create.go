package releaseRepository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-dependency-api/internal/customErrors"
)

func (r *Neo4jReleaseRepository) CreateRelease(ctx context.Context, release Release) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		_ = session.Close(ctx)
	}()

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
			return nil, &customErrors.HTTPError{
				Status: 404,
				Msg:    fmt.Sprintf("Service not found: %s", release.ServiceId),
			}
		}

		// Create the release node and connect it to the service
		query := `
			MATCH (s:Service {id: $serviceId})
			CREATE (r:Release {
				serviceId: $serviceId,
				tag: $tag,
				releaseDate: datetime($releaseDate),
				githubUrl: $githubUrl,
				notes: $notes
			})
			CREATE (s)-[rel:RELEASED]->(r)
			RETURN r
		`
		params := map[string]any{
			"serviceId":   release.ServiceId,
			"releaseDate": release.ReleaseDate.Format("2006-01-02T15:04:05Z"),
			"githubUrl":   release.Url,
		}

		_, err = tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	_, err := session.ExecuteWrite(ctx, createReleaseTransaction)
	return err
}
