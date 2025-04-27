package dependencies

import (
	"context"
	"service-dependency-api/api/dependencies/internal/dependencyRepository"
)

type mockDependencyRepository struct {
	Data func() []map[string]any
	Err  error
}

func (repo mockDependencyRepository) AddDependency(_ context.Context, _ string, _ *dependencyRepository.Dependency) error {
	if repo.Err != nil {
		return repo.Err
	}

	// If no error, we consider the operation successful
	// In a real implementation, we might want to check if the service exists, etc.
	return nil
}
