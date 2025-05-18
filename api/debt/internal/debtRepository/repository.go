package debtRepository

import "github.com/neo4j/neo4j-go-driver/v5/neo4j"

type Neo4jDebtRepository struct {
	driver neo4j.DriverWithContext
}

func New(driver neo4j.DriverWithContext) *Neo4jDebtRepository {
	return &Neo4jDebtRepository{driver: driver}
}
