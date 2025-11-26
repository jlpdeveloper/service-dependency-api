package servicerepository

import (
	"service-atlas/databaseadapter"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Neo4jServiceRepository type implements the interface for service repository above
type Neo4jServiceRepository struct {
	manager databaseadapter.DriverManager
}

func New(driver neo4j.DriverWithContext) *Neo4jServiceRepository {
	return &Neo4jServiceRepository{manager: databaseadapter.NewDriverManager(driver)}
}
