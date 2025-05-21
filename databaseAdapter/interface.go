package databaseAdapter

import "context"

// DriverManager interface abstracts Neo4j driver operations
type DriverManager interface {
	NewSession(ctx context.Context, config SessionConfig) SessionManager
}

// SessionManager interface abstracts Neo4j session operations
type SessionManager interface {
	ExecuteWrite(ctx context.Context, work func(tx TransactionManager) (any, error)) (any, error)
	Close(ctx context.Context) error
}

// TransactionManager interface abstracts Neo4j transaction operations
type TransactionManager interface {
	Run(ctx context.Context, query string, params map[string]any) (ResultManager, error)
}

// ResultManager interface abstracts Neo4j result operations
type ResultManager interface {
	Collect(ctx context.Context) ([]Record, error)
}

// Record interface abstracts Neo4j record operations
type Record interface {
	// Add methods as needed
}

// SessionConfig struct to hold session configuration
type SessionConfig struct {
	AccessMode string
}
