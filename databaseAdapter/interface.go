package databaseAdapter

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// DriverManager interface abstracts Neo4j driver operations
type DriverManager interface {
	ExecuteWrite(ctx context.Context, work func(tx neo4j.ManagedTransaction) (any, error)) (any, error)
	ExecuteRead(ctx context.Context, work func(tx neo4j.ManagedTransaction) (any, error)) (any, error)
}

func NewDriverManager(driver neo4j.DriverWithContext) DriverManager {
	return &Neo4jDriverAdapter{driver}
}

// Neo4jDriverAdapter adapts Neo4j driver to DriverManager
type Neo4jDriverAdapter struct {
	driver neo4j.DriverWithContext
}

func (n Neo4jDriverAdapter) executeInSession(ctx context.Context, work func(tx neo4j.ManagedTransaction) (any, error), mode neo4j.AccessMode) (any, error) {
	session := n.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: mode})
	defer func() {
		_ = session.Close(ctx)
	}()
	result, err := session.ExecuteWrite(ctx, work)

	return result, err
}

func (n Neo4jDriverAdapter) ExecuteWrite(ctx context.Context, work func(tx neo4j.ManagedTransaction) (any, error)) (any, error) {
	return n.executeInSession(ctx, work, neo4j.AccessModeWrite)
}

func (n Neo4jDriverAdapter) ExecuteRead(ctx context.Context, work func(tx neo4j.ManagedTransaction) (any, error)) (any, error) {
	return n.executeInSession(ctx, work, neo4j.AccessModeRead)
}
