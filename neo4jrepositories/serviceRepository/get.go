package serviceRepository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	nRepo "service-dependency-api/neo4jrepositories"
	"service-dependency-api/repositories"
)

func (d *Neo4jServiceRepository) GetAllServices(ctx context.Context, page int, pageSize int) (services []repositories.Service, err error) {
	services = []repositories.Service{}
	getPagedData := func(tx neo4j.ManagedTransaction) (any, error) {
		skip := (page - 1) * pageSize

		result, err := tx.Run(ctx, `
		    MATCH (s:Service)
			RETURN s
			ORDER BY s.createdDate DESC
			SKIP $skip
			LIMIT $limit
		`, map[string]any{
			"skip":  skip,
			"limit": pageSize,
		})

		if err != nil {
			return nil, err
		}

		for result.Next(ctx) {
			record := result.Record()
			node, ok := record.Get("s")
			if !ok {
				continue
			}

			n, ok := node.(neo4j.Node)
			if !ok {
				continue
			}

			svc := nRepo.MapNodeToService(n)
			services = append(services, svc)
		}
		return services, nil
	}
	_, readErr := d.manager.ExecuteRead(ctx, getPagedData)
	if readErr != nil {
		return nil, readErr
	}
	return services, nil
}

func (d *Neo4jServiceRepository) GetServiceById(ctx context.Context, id string) (svc repositories.Service, err error) {
	getServiceById := func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (s:Service)
			WHERE s.id = $id
			RETURN s
		`, map[string]any{
			"id": id,
		})

		if err != nil {
			return repositories.Service{}, err
		}

		if !result.Next(ctx) {
			return repositories.Service{}, nil // No service found with this ID
		}

		record := result.Record()
		node, ok := record.Get("s")
		if !ok {
			return repositories.Service{}, nil
		}

		n, ok := node.(neo4j.Node)
		if !ok {
			return repositories.Service{}, nil
		}

		return nRepo.MapNodeToService(n), nil
	}

	service, readErr := d.manager.ExecuteRead(ctx, getServiceById)
	if readErr != nil {
		return repositories.Service{}, readErr
	}

	if service == nil {
		return repositories.Service{}, nil
	}

	return service.(repositories.Service), nil
}
