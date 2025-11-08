package databaseadapter

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// DriverManager interface abstracts Neo4j driver operations
type DriverManager interface {
	ExecuteWrite(ctx context.Context, work func(tx neo4j.ManagedTransaction) (any, error)) (any, error)
	ExecuteRead(ctx context.Context, work func(tx neo4j.ManagedTransaction) (any, error)) (any, error)
}
