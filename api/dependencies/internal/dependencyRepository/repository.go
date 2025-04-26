package dependencyRepository

import "github.com/neo4j/neo4j-go-driver/v5/neo4j"

type Neo4jDependencyRepository struct {
	driver neo4j.DriverWithContext
}

func New(driver neo4j.DriverWithContext) *Neo4jDependencyRepository {
	return &Neo4jDependencyRepository{driver: driver}
}
