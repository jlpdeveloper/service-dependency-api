package dependencyRepository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-dependency-api/databaseAdapter"
)

type Neo4jDependencyRepository struct {
	manager databaseAdapter.DriverManager
}

func New(driver neo4j.DriverWithContext) *Neo4jDependencyRepository {
	return &Neo4jDependencyRepository{manager: databaseAdapter.NewDriverManager(driver)}
}
