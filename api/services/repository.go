package services

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"sync"
)

type Service struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name"`
	ServiceType string `json:"type"`
	Description string `json:"description"`
}

type ServiceRepository interface {
	GetAllServices(page int, pageSize int) ([]Service, error)
	CreateService(service Service) (string, error)
	//UpdateService(service Service) error
	//DeleteService(service Service) error
	//GetServiceById(id int) (Service, error)
}

// the type implements the interface
type ServiceNeo4jRepository struct {
	Driver neo4j.DriverWithContext
	Ctx    context.Context
}

func (d *ServiceNeo4jRepository) GetAllServices(page int, pageSize int) (services []Service, err error) {
	session := d.Driver.NewSession(d.Ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() {
		err = session.Close(d.Ctx)
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

			svc := Service{
				Name:        n.Props["name"].(string),
				Description: n.Props["description"].(string),
				ServiceType: n.Props["type"].(string),
				Id:          n.Props["id"].(string),
			}

			services = append(services, svc)
		}
		return services, nil
	}
	_, readErr := session.ExecuteRead(d.Ctx, getPagedData)
	if readErr != nil {
		return nil, readErr
	}
	return services, nil
}

func (d *ServiceNeo4jRepository) CreateService(service Service) (id string, err error) {
	session := d.Driver.NewSession(d.Ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		err = session.Close(d.Ctx)
	}()

	createServiceTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(
			d.Ctx, `
        CREATE (n: Service {id: randomuuid(), createdDate: datetime(), name: $name, type: $type, description: $description})
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
		svcId, _ := svc.AsMap()["id"]
		return svcId.(string), err

	}
	newId, insertErr := session.ExecuteWrite(d.Ctx, createServiceTransaction)
	if insertErr != nil {
		return "", insertErr
	}
	return newId.(string), nil
}
