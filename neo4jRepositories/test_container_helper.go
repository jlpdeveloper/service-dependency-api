package neo4jRepositories

import (
	"context"
	"errors"

	"github.com/testcontainers/testcontainers-go/modules/neo4j"
)

type TestContainerHelper struct {
	Container *neo4j.Neo4jContainer
	Endpoint  string
}

// NewTestContainerHelper creates a new Neo4j container for testing purposes
// returns a TestContainerHelper struct that contains the container and the endpoint
func NewTestContainerHelper(ctx context.Context) (*TestContainerHelper, error) {
	neo4jContainer, err := neo4j.Run(ctx,
		"neo4j:latest",
		neo4j.WithAdminPassword("letmein!"),
	)
	if err != nil {
		return nil, err
	}
	if neo4jContainer == nil {
		return nil, errors.New("neo4jContainer is nil")
	}
	dbPort, err := neo4jContainer.MappedPort(ctx, "7687/tcp")
	if err != nil {
		return nil, err
	}
	host, err := neo4jContainer.Host(ctx)
	if err != nil {
		return nil, err
	}
	endpoint := "neo4j://" + host + ":" + dbPort.Port()
	return &TestContainerHelper{
		Container: neo4jContainer,
		Endpoint:  endpoint,
	}, nil

}
