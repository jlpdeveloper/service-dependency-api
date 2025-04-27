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

func (repo mockDependencyRepository) GetDependencies(_ context.Context, _ string) ([]*dependencyRepository.Dependency, error) {
	if repo.Err != nil {
		return nil, repo.Err
	}

	// Convert the mock data to the expected return type
	data := repo.Data()
	dependencies := make([]*dependencyRepository.Dependency, 0, len(data))

	for _, item := range data {
		dep := &dependencyRepository.Dependency{}

		if id, ok := item["id"].(string); ok {
			dep.Id = id
		}
		if name, ok := item["name"].(string); ok {
			dep.Name = name
		}
		if version, ok := item["version"].(string); ok {
			dep.Version = version
		}

		dependencies = append(dependencies, dep)
	}

	return dependencies, nil
}
