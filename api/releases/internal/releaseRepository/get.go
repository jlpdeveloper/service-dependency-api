package releaseRepository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
	"service-dependency-api/internal/customErrors"
	"time"
)

func (r *Neo4jReleaseRepository) GetReleasesByServiceId(ctx context.Context, serviceId string, page, pageSize int) ([]*Release, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() {
		_ = session.Close(ctx)
	}()

	if page <= 0 || pageSize <= 0 {
		return nil, &customErrors.HTTPError{
			Status: http.StatusBadRequest,
			Msg:    "page and page size must be positive integers",
		}
	}

	getReleasesByServiceIdTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		// First check if the service exists
		checkQuery := `
			MATCH (s:Service {id: $serviceId})
			RETURN s
		`
		result, err := tx.Run(ctx, checkQuery, map[string]any{
			"serviceId": serviceId,
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
				Msg:    fmt.Sprintf("Service not found: %s", serviceId),
			}
		}

		// Find all releases for the service with the given ID, ordered by release date descending with pagination
		query := `
			MATCH (s:Service {id: $serviceId})-[rel:RELEASED]->(r:Release)
			RETURN r.releaseDate as releaseDate, r.url as url, r.version as version
			ORDER BY r.releaseDate DESC
			SKIP $skip
			LIMIT $limit
		`
		result, err = tx.Run(ctx, query, map[string]any{
			"serviceId": serviceId,
			"skip":      (page - 1) * pageSize,
			"limit":     pageSize,
		})
		if err != nil {
			return nil, err
		}

		releases := []*Release{}
		records, err = result.Collect(ctx)
		if err != nil {
			return nil, err
		}

		for _, record := range records {
			releaseDate, _ := record.Get("releaseDate")
			url, _ := record.Get("url")
			version, _ := record.Get("version")

			release := &Release{
				ServiceId:   serviceId,
				ReleaseDate: releaseDate.(time.Time),
			}

			// Only set url and version if they exist
			if url != nil {
				release.Url = url.(string)
			}
			if version != nil {
				release.Version = version.(string)
			}

			releases = append(releases, release)
		}

		return releases, nil
	}

	result, err := session.ExecuteRead(ctx, getReleasesByServiceIdTransaction)
	if err != nil {
		return nil, err
	}

	return result.([]*Release), nil
}
