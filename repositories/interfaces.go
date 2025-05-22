package repositories

import "context"

type DebtRepository interface {
	CreateDebtItem(ctx context.Context, debt Debt) error
	UpdateStatus(ctx context.Context, id, status string) error
}

type ServiceRepository interface {
	GetAllServices(ctx context.Context, page int, pageSize int) ([]Service, error)
	CreateService(ctx context.Context, service Service) (string, error)
	UpdateService(ctx context.Context, service Service) error
	DeleteService(ctx context.Context, id string) error
	GetServiceById(ctx context.Context, id string) (Service, error)
}

type DependencyRepository interface {
	AddDependency(ctx context.Context, id string, dependency *Dependency) error
	GetDependencies(ctx context.Context, id string) ([]*Dependency, error)
	GetDependents(ctx context.Context, id string) ([]*Dependency, error)
	DeleteDependency(ctx context.Context, id string, dependsOnID string) error
}
