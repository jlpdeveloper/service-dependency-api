package service_repository

import (
	"context"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (d *ServiceNeo4jRepository) UpdateService(ctx context.Context, service Service) (found bool, err error) {
	session := d.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		closeErr := session.Close(ctx)
		if err == nil {
			err = closeErr
		}
	}()
	updateServiceTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		// First check if the service exists
		result, err := tx.Run(ctx, `
			MATCH (s:Service)
			WHERE s.id = $id
			RETURN s
		`, map[string]any{
			"id": service.Id,
		})

		if err != nil {
			return found, err
		}

		found = result.Next(ctx)

		if found {

			// Service exists, update it
			updateResult, updateErr := tx.Run(ctx, `
			MATCH (s:Service)
			WHERE s.id = $id
			SET s.name = $name, 
				s.type = $type, 
				s.description = $description,
				s.updated = datetime()
			RETURN s
		`, map[string]any{
				"id":          service.Id,
				"name":        service.Name,
				"type":        service.ServiceType,
				"description": service.Description,
			})

			if updateErr != nil {
				err = updateErr
			}

			// Confirm update was successful
			if !updateResult.Next(ctx) {
				err = errors.New("update Service failed")
			}
		}
		return found, err
	}

	result, execErr := session.ExecuteWrite(ctx, updateServiceTransaction)
	if execErr != nil {
		return false, execErr
	}

	return result.(bool), nil
}
