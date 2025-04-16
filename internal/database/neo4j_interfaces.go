package database

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ManagedTransactionWork func(tx Neo4jTransaction) (any, error)

type Neo4jSession interface {
	ExecuteWrite(ctx context.Context, work ManagedTransactionWork) (any, error)
	Close(ctx context.Context) error
}

type Neo4jTransaction interface {
	Run(ctx context.Context, cypher string, params map[string]any) (Neo4jResult, error)
}

type Neo4jResult interface {
	Single(ctx context.Context) (Neo4jRecord, error)
}

type Neo4jRecord interface {
	AsMap() map[string]any
}

type Neo4jDriver interface {
	NewSession(ctx context.Context, config neo4j.SessionConfig) Neo4jSession
}
