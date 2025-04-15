package services

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-dependency-api/internal/database"
)

type Service struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name"`
	ServiceType string `json:"type"`
	Description string `json:"description"`
}

type ServiceRepository interface {
	//GetAllServices(page int) ([]Service, error)
	CreateService(service Service) (string, error)
	//UpdateService(service Service) error
	//DeleteService(service Service) error
	//GetServiceById(id int) (Service, error)
}

// the type implements the interface
type ServiceNeo4jService struct {
	Driver database.Neo4jDriver
	ctx    context.Context
}

//https://github.com/neo4j-examples/golang-neo4j-realworld-example/blob/main/pkg/users/login.go
//https://neo4j.com/docs/go-manual/current/
//func (d *ServiceNeo4jService) GetAllServices(page int) (services []Service, err error) {
//	session := d.Driver.NewSession(d.ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
//	defer func() {
//		err = session.Close(d.ctx)
//	}()
//}

func (d *ServiceNeo4jService) CreateService(service Service) (id string, err error) {
	session := d.Driver.NewSession(d.ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		err = session.Close(d.ctx)
	}()
	createServiceTransaction := func(tx database.Neo4jTransaction) (any, error) {
		result, err := tx.Run(
			d.ctx, `
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
		svc, err := result.Single(d.ctx)
		if err != nil {
			return "", err
		}
		svcId, _ := svc.AsMap()["id"]
		return svcId.(string), err
	}
	newId, insertErr := session.ExecuteWrite(d.ctx, createServiceTransaction)
	if insertErr != nil {
		return "", insertErr
	}
	return newId.(string), nil
}
