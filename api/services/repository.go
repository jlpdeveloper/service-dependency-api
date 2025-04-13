package services

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Service struct {
	id          int    `json:"id"`
	name        string `json:"name"`
	serviceType string `json:"type"`
	description string `json:"description"`
}

type ServiceRepository interface {
	//GetAllServices(page int) ([]Service, error)
	//CreateService(service Service) (int, error)
	//UpdateService(service Service) error
	//DeleteService(service Service) error
	//GetServiceById(id int) (Service, error)
}

// the type implements the interface
type ServiceNeo4jService struct {
	Driver *neo4j.DriverWithContext
	ctx    *context.Context
}

//https://github.com/neo4j-examples/golang-neo4j-realworld-example/blob/main/pkg/users/login.go
//https://neo4j.com/docs/go-manual/current/
//func (d *ServiceNeo4jService) GetAllServices(page int) (services []Service, err error) {
//	session := d.Driver.NewSession(d.ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
//	defer func() {
//		err = session.Close(d.ctx)
//	}()
//}
