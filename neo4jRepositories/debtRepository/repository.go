package debtRepository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-dependency-api/databaseAdapter"
)

type Neo4jDebtRepository struct {
	manager databaseAdapter.DriverManager
}

func New(driver neo4j.DriverWithContext) *Neo4jDebtRepository {
	manager := databaseAdapter.NewDriverManager(driver)
	return &Neo4jDebtRepository{manager: manager}
}

const (
	DefaultStatus = "pending"
)
