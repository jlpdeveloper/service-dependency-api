package dependencyRepository

import (
	"errors"
)

type Dependency struct {
	Id      string `json:"id"`
	Version string `json:"version,omitempty"`
}

func (d *Dependency) Validate() error {
	if d.Id == "" {
		return errors.New("dependency id is required")
	}
	return nil
}

type DependencyRepository interface {
	AddDependency(id string, dependency *Dependency) error
}
