package serviceRepository

import (
	"context"
	"service-dependency-api/databaseAdapter"
	"service-dependency-api/repositories"
)

func (d *Neo4jServiceRepository) CreateService(ctx context.Context, service repositories.Service) (id string, err error) {
	session := d.manager.NewSession(ctx, databaseAdapter.SessionConfig{AccessMode: "write"})
	defer func() {
		closeErr := session.Close(ctx)
		if err == nil {
			err = closeErr
		}
	}()

	createServiceTransaction := func(tx databaseAdapter.TransactionManager) (any, error) {
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
		svcId, err := result.GetProperty(ctx, "id")
		if svcId != nil {
			if idStr, ok := svcId.(string); ok {
				return idStr, err
			}
		}
		return "", err

	}
	newId, insertErr := session.ExecuteWrite(ctx, createServiceTransaction)
	if insertErr != nil {
		return "", insertErr
	}
	return newId.(string), nil
}
