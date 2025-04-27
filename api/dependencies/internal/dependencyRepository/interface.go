package dependencyRepository

import (
	"context"
	"errors"
)

type Dependency struct {
	Id      string `json:"id"`
	Version string `json:"version,omitempty"`
	Name    string `json:"name,omitempty"`
}

func (d *Dependency) Validate() error {
	if d.Id == "" {
		return errors.New("dependency id is required")
	}
	return nil
}

type DependencyRepository interface {
	AddDependency(ctx context.Context, id string, dependency *Dependency) error
	GetDependencies(ctx context.Context, id string) ([]*Dependency, error)
}
