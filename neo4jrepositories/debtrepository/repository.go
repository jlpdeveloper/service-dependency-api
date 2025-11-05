package debtrepository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-dependency-api/databaseadapter"
)

type Neo4jDebtRepository struct {
	manager databaseadapter.DriverManager
}

func New(driver neo4j.DriverWithContext) *Neo4jDebtRepository {
	return &Neo4jDebtRepository{manager: databaseadapter.NewDriverManager(driver)}
}

const (
	// DefaultStatus The default status for any newly created debt items
	DefaultStatus = "pending"
)
