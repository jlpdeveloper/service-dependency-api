package debtRepository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-dependency-api/databaseAdapter"
	"service-dependency-api/neo4jRepositories/neo4jAdapter"
)

type Neo4jDebtRepository struct {
	manager databaseAdapter.DriverManager
}

func New(driver neo4j.DriverWithContext) *Neo4jDebtRepository {
	manager := neo4jAdapter.NewDriverManager(driver)
	return &Neo4jDebtRepository{manager: manager}
}

const (
	DefaultStatus = "pending"
)
