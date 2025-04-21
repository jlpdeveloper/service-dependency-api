package services

import (
	"context"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"service-dependency-api/internal"
	"sync"
	"time"
)

type Service struct {
	Id          string    `json:"id,omitempty"`
	Name        string    `json:"name"`
	ServiceType string    `json:"type"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated,omitempty"`
}

type ServiceRepository interface {
	GetAllServices(page int, pageSize int) ([]Service, error)
	CreateService(service Service) (string, error)
	UpdateService(service Service) (bool, error)
	DeleteService(id string) error
	GetServiceById(id string) (Service, error)
}

// ServiceNeo4jRepository type implements the interface for service repository above
type ServiceNeo4jRepository struct {
	Driver neo4j.DriverWithContext
	Ctx    context.Context
}

func (d *ServiceNeo4jRepository) DeleteService(id string) (err error) {
	session := d.Driver.NewSession(d.Ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		closeErr := session.Close(d.Ctx)
		if err == nil {
			err = closeErr
		}
	}()
	deleteServiceTransaction := func(tx neo4j.ManagedTransaction) (any, error) {

		result, err := tx.Run(d.Ctx, `
    		MATCH (s:Service { id: $id })
    		RETURN count(s) as count
		`, map[string]interface{}{"id": id})

		if err != nil {
			return nil, err
		}

		if record, err := result.Single(d.Ctx); err == nil {
			count, _ := record.Get("count")
			if count.(int64) == 0 {
				return nil, &internal.HTTPError{
					Status: 404,
					Msg:    "Service not found",
				}
			}
		}
		result, err = tx.Run(d.Ctx, `
		MATCH(s:Service { id: $id})
		DELETE s;`, map[string]interface{}{"id": id})
		if err != nil {
			log.Println("Error deleting service: " + id)
			return nil, err
		}

		summary, err := result.Consume(d.Ctx)
		if err != nil {
			return nil, &internal.HTTPError{Status: 500, Msg: "Error deleting service: " + id}
		}

		if summary.Counters().NodesDeleted() == 0 {
			log.Println("Error deleting service: " + id + ". Database transaction not successful")
			return nil, &internal.HTTPError{Status: 500, Msg: "Error deleting service: " + id}
		}
		return nil, nil
	}

	_, err = session.ExecuteWrite(d.Ctx, deleteServiceTransaction)
	return err
}

func (d *ServiceNeo4jRepository) UpdateService(service Service) (found bool, err error) {
	session := d.Driver.NewSession(d.Ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		closeErr := session.Close(d.Ctx)
		if err == nil {
			err = closeErr
		}
	}()
	updateServiceTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		// First check if the service exists
		result, err := tx.Run(d.Ctx, `
			MATCH (s:Service)
			WHERE s.id = $id
			RETURN s
		`, map[string]any{
			"id": service.Id,
		})

		if err != nil {
			return found, err
		}

		found = result.Next(d.Ctx)

		if found {

			// Service exists, update it
			updateResult, updateErr := tx.Run(d.Ctx, `
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
			if !updateResult.Next(d.Ctx) {
				err = errors.New("update Service failed")
			}
		}
		return found, err
	}

	result, execErr := session.ExecuteWrite(d.Ctx, updateServiceTransaction)
	if execErr != nil {
		return false, execErr
	}

	return result.(bool), nil
}

func (d *ServiceNeo4jRepository) GetServiceById(id string) (svc Service, err error) {
	session := d.Driver.NewSession(d.Ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() {
		closeErr := session.Close(d.Ctx)
		if err == nil {
			err = closeErr
		}
	}()

	getServiceById := func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(d.Ctx, `
			MATCH (s:Service)
			WHERE s.id = $id
			RETURN s
		`, map[string]any{
			"id": id,
		})

		if err != nil {
			return Service{}, err
		}

		if !result.Next(d.Ctx) {
			return Service{}, nil // No service found with this ID
		}

		record := result.Record()
		node, ok := record.Get("s")
		if !ok {
			return Service{}, nil
		}

		n, ok := node.(neo4j.Node)
		if !ok {
			return Service{}, nil
		}

		return d.mapNodeToService(n), nil
	}

	service, readErr := session.ExecuteRead(d.Ctx, getServiceById)
	if readErr != nil {
		return Service{}, readErr
	}

	if service == nil {
		return Service{}, nil
	}

	return service.(Service), nil
}

// mapNodeToService converts a Neo4j node to a Service object
func (d *ServiceNeo4jRepository) mapNodeToService(n neo4j.Node) Service {
	svc := Service{}

	// Safely extract name with validation
	if name, ok := n.Props["name"]; ok {
		if nameStr, ok := name.(string); ok {
			svc.Name = nameStr
		}
	}

	// Safely extract description with validation
	if desc, ok := n.Props["description"]; ok {
		if descStr, ok := desc.(string); ok {
			svc.Description = descStr
		}
	}

	// Safely extract service type with validation
	if svcType, ok := n.Props["type"]; ok {
		if typeStr, ok := svcType.(string); ok {
			svc.ServiceType = typeStr
		}
	}

	// Safely extract ID with validation
	if id, ok := n.Props["id"]; ok {
		if idStr, ok := id.(string); ok {
			svc.Id = idStr
		}
	}

	// Safely extract created date with validation
	if date, ok := n.Props["created"]; ok {
		if dateStr, ok := date.(time.Time); ok {
			svc.Created = dateStr
		}
	}

	if date, ok := n.Props["updated"]; ok {
		if dateStr, ok := date.(time.Time); ok {
			svc.Updated = dateStr
		}
	}
	return svc
}

func (d *ServiceNeo4jRepository) GetAllServices(page int, pageSize int) (services []Service, err error) {
	session := d.Driver.NewSession(d.Ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() {
		closeErr := session.Close(d.Ctx)
		if err == nil {
			err = closeErr
		}
	}()
	services = []Service{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	getPagedData := func(tx neo4j.ManagedTransaction) (any, error) {
		defer wg.Done()
		skip := (page - 1) * pageSize

		result, err := tx.Run(d.Ctx, `
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

		for result.Next(d.Ctx) {
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
	_, readErr := session.ExecuteRead(d.Ctx, getPagedData)
	wg.Wait()
	if readErr != nil {
		return nil, readErr
	}
	return services, nil
}

func (d *ServiceNeo4jRepository) CreateService(service Service) (id string, err error) {
	session := d.Driver.NewSession(d.Ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		closeErr := session.Close(d.Ctx)
		if err == nil {
			err = closeErr
		}
	}()

	createServiceTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(
			d.Ctx, `
        CREATE (n: Service {id: randomuuid(), created: datetime(), name: $name, type: $type, description: $description})
        RETURN n.id AS id
        `, map[string]any{
				"name":        service.Name,
				"type":        service.ServiceType,
				"description": service.Description,
			})
		if err != nil {
			return "", err
		}
		svc, err := result.Single(d.Ctx)
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
	newId, insertErr := session.ExecuteWrite(d.Ctx, createServiceTransaction)
	if insertErr != nil {
		return "", insertErr
	}
	return newId.(string), nil
}
