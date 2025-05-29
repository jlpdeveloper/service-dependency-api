package repositories

import (
	"context"
	"time"
)

type DebtRepository interface {
	CreateDebtItem(ctx context.Context, debt Debt) error
	UpdateStatus(ctx context.Context, id, status string) error
	GetDebtByServiceId(ctx context.Context, id string, page, pageSize int, onlyResolved bool) ([]Debt, error)
}

type ServiceRepository interface {
	GetAllServices(ctx context.Context, page int, pageSize int) ([]Service, error)
	CreateService(ctx context.Context, service Service) (string, error)
	UpdateService(ctx context.Context, service Service) error
	DeleteService(ctx context.Context, id string) error
	GetServiceById(ctx context.Context, id string) (Service, error)
}

type DependencyRepository interface {
	AddDependency(ctx context.Context, id string, dependency Dependency) error
	GetDependencies(ctx context.Context, id string) ([]*Dependency, error)
	GetDependents(ctx context.Context, id string) ([]*Dependency, error)
	DeleteDependency(ctx context.Context, id string, dependsOnID string) error
}

type ReleaseRepository interface {
	CreateRelease(ctx context.Context, release Release) error
	GetReleasesByServiceId(ctx context.Context, serviceId string, page, pageSize int) ([]*Release, error)
	GetReleasesInDateRange(ctx context.Context, startDate, endDate time.Time, page, pageSize int) ([]*ServiceReleaseInfo, error)
}

type ReportRepository interface {
	GetServiceRiskReport(ctx context.Context, serviceId string) (*ServiceRiskReport, error)
}
