package serviceRepository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-dependency-api/repositories"
)

func (d *Neo4jServiceRepository) CreateService(ctx context.Context, service repositories.Service) (id string, err error) {
	createServiceTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(
			ctx, `
        CREATE (n: Service {id: randomuuid(), created: datetime(), name: $name, type: $type, description: $description, url: $url})
        RETURN n.id AS id
        `, map[string]any{
				"name":        service.Name,
				"type":        service.ServiceType,
				"description": service.Description,
				"url":         service.Url,
			})
		if err != nil {
			return "", err
		}
		svc, err := result.Single(ctx)
		if err != nil {
			return "", err
		}
		svcMap := svc.AsMap()
		if svcId, ok := svcMap["id"]; ok {
			if idStr, ok := svcId.(string); ok {
				return idStr, err
			}
		}
		return "", err

	}
	newId, insertErr := d.manager.ExecuteWrite(ctx, createServiceTransaction)
	if insertErr != nil {
		return "", insertErr
	}
	return newId.(string), nil
}
