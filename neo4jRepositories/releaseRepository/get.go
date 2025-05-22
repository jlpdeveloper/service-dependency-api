package releaseRepository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
	"service-dependency-api/internal/customErrors"
	"service-dependency-api/repositories"
	"time"
)

func (r *Neo4jReleaseRepository) GetReleasesByServiceId(ctx context.Context, serviceId string, page, pageSize int) ([]*repositories.Release, error) {
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

		releases := []*repositories.Release{}
		records, err = result.Collect(ctx)
		if err != nil {
			return nil, err
		}

		for _, record := range records {
			releaseDate, _ := record.Get("releaseDate")
			release := &repositories.Release{
				ServiceId:   serviceId,
				ReleaseDate: releaseDate.(time.Time),
			}

			// Only set url and version if they exist
			if url, ok := record.Get("url"); ok {
				release.Url = url.(string)
			}
			if version, ok := record.Get("version"); ok {
				release.Version = version.(string)
			}

			releases = append(releases, release)
		}

		return releases, nil
	}

	result, err := r.manager.ExecuteRead(ctx, getReleasesByServiceIdTransaction)
	if err != nil {
		return nil, err
	}

	return result.([]*repositories.Release), nil
}

func (r *Neo4jReleaseRepository) GetReleasesInDateRange(ctx context.Context, startDate, endDate time.Time, page, pageSize int) ([]*repositories.ServiceReleaseInfo, error) {
	getReleasesInRangeTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (s:Service)-[rel:RELEASED]->(r:Release)
			WHERE r.releaseDate >= datetime($startDate) AND r.releaseDate <= datetime($endDate)
			RETURN r.releaseDate as releaseDate, r.url as url, r.version as version, s.id as serviceId, s.name as serviceName, s.type as serviceType
			ORDER BY r.releaseDate DESC
			SKIP $skip
			LIMIT $limit
		`

		result, err := tx.Run(ctx, query, map[string]any{
			"startDate": startDate.Format("2006-01-02"),
			"endDate":   endDate.Format("2006-01-02"),
			"skip":      (page - 1) * pageSize,
			"limit":     pageSize,
		})
		if err != nil {
			return nil, err
		}
		records, err := result.Collect(ctx)
		if err != nil {
			return nil, err
		}

		releases := []*repositories.ServiceReleaseInfo{}
		for _, record := range records {
			releaseDate, _ := record.Get("releaseDate")
			serviceId, _ := record.Get("serviceId")
			serviceName, _ := record.Get("serviceName")
			serviceType, _ := record.Get("serviceType")
			release := &repositories.ServiceReleaseInfo{
				ServiceName: serviceName.(string),
				ServiceType: serviceType.(string),
				Release: repositories.Release{
					ReleaseDate: releaseDate.(time.Time),
					ServiceId:   serviceId.(string),
				},
			}
			if url, ok := record.Get("url"); ok && url != nil {
				release.Url = url.(string)
			}
			if version, ok := record.Get("version"); ok && version != nil {
				release.Version = version.(string)
			}

			releases = append(releases, release)
		}

		return releases, nil

	}
	result, err := r.manager.ExecuteRead(ctx, getReleasesInRangeTransaction)
	if err != nil {
		return nil, err
	}
	return result.([]*repositories.ServiceReleaseInfo), nil
}
