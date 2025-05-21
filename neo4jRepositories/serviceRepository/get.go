package serviceRepository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-dependency-api/repositories"
	"sync"
)

func (d *Neo4jServiceRepository) GetAllServices(ctx context.Context, page int, pageSize int) (services []repositories.Service, err error) {
	session := d.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() {
		closeErr := session.Close(ctx)
		if err == nil {
			err = closeErr
		}
	}()
	services = []repositories.Service{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	getPagedData := func(tx neo4j.ManagedTransaction) (any, error) {
		defer wg.Done()
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

			svc := d.mapNodeToService(n)
			services = append(services, svc)
		}
		return services, nil
	}
	_, readErr := session.ExecuteRead(ctx, getPagedData)
	wg.Wait()
	if readErr != nil {
		return nil, readErr
	}
	return services, nil
}

func (d *Neo4jServiceRepository) GetServiceById(ctx context.Context, id string) (svc repositories.Service, err error) {
	session := d.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() {
		closeErr := session.Close(ctx)
		if err == nil {
			err = closeErr
		}
	}()

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

		return d.mapNodeToService(n), nil
	}

	service, readErr := session.ExecuteRead(ctx, getServiceById)
	if readErr != nil {
		return repositories.Service{}, readErr
	}

	if service == nil {
		return repositories.Service{}, nil
	}

	return service.(repositories.Service), nil
}
