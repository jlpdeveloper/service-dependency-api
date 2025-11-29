package dependencyrepository

import (
	"service-atlas/databaseadapter"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jDependencyRepository struct {
	manager databaseadapter.DriverManager
}

func New(driver neo4j.DriverWithContext) *Neo4jDependencyRepository {
	return &Neo4jDependencyRepository{manager: databaseadapter.NewDriverManager(driver)}
}
