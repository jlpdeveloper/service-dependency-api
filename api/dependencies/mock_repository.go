package dependencies

import (
	"context"
	"fmt"
	"service-dependency-api/api/dependencies/internal/dependencyRepository"
	"service-dependency-api/internal/customErrors"
)

type mockDependencyRepository struct {
	Data func() []map[string]any
	Err  error
	// DependencyExists is used to determine if a dependency exists in the mock repository
	DependencyExists bool
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

func (repo mockDependencyRepository) GetDependents(_ context.Context, _ string) ([]*dependencyRepository.Dependency, error) {
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

func (repo mockDependencyRepository) DeleteDependency(_ context.Context, id string, dependsOnID string) error {
	if repo.Err != nil {
		return repo.Err
	}

	// If DependencyExists is false, return a 404 error
	if !repo.DependencyExists {
		return &customErrors.HTTPError{
			Status: 404,
			Msg:    fmt.Sprintf("Dependency relationship not found between services: %s -> %s", id, dependsOnID),
		}
	}

	// If no error and dependency exists, we consider the operation successful
	return nil
}
