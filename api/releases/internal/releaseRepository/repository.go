package releaseRepository

import "github.com/neo4j/neo4j-go-driver/v5/neo4j"

type Neo4jReleaseRepository struct {
	driver neo4j.DriverWithContext
}

func New(driver neo4j.DriverWithContext) *Neo4jReleaseRepository {
	return &Neo4jReleaseRepository{
		driver: driver,
	}
}
