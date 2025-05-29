package repositories

import (
	"context"
	"time"
)

// DebtRepository defines the methods for interacting with debt items.
type DebtRepository interface {
	// CreateDebtItem creates a new debt item.
	CreateDebtItem(ctx context.Context, debt Debt) error
	// UpdateStatus updates the status of an existing debt item.
	UpdateStatus(ctx context.Context, id, status string) error
	// GetDebtByServiceId retrieves debt items for a given service.
	GetDebtByServiceId(ctx context.Context, id string, page, pageSize int, onlyResolved bool) ([]Debt, error)
}

// ServiceRepository defines the methods for interacting with services.
type ServiceRepository interface {
	// GetAllServices retrieves all services.
	GetAllServices(ctx context.Context, page int, pageSize int) ([]Service, error)
	// CreateService creates a new service.
	CreateService(ctx context.Context, service Service) (string, error)
	// UpdateService updates an existing service.
	UpdateService(ctx context.Context, service Service) error
	// DeleteService deletes a service.
	DeleteService(ctx context.Context, id string) error
	// GetServiceById retrieves a service by its ID.
	GetServiceById(ctx context.Context, id string) (Service, error)
}

// DependencyRepository defines the methods for interacting with dependencies.
type DependencyRepository interface {
	// AddDependency adds a dependency to a resource.
	AddDependency(ctx context.Context, id string, dependency Dependency) error
	// GetDependencies retrieves all dependencies of a resource.
	GetDependencies(ctx context.Context, id string) ([]*Dependency, error)
	// GetDependents retrieves all resources that depend on a given resource.
	GetDependents(ctx context.Context, id string) ([]*Dependency, error)
	// DeleteDependency deletes a dependency between two resources.
	DeleteDependency(ctx context.Context, id string, dependsOnID string) error
}

// ReleaseRepository defines the methods for interacting with releases.
type ReleaseRepository interface {
	// CreateRelease creates a new release.
	CreateRelease(ctx context.Context, release Release) error
	// GetReleasesByServiceId retrieves all releases associated with a service.
	GetReleasesByServiceId(ctx context.Context, serviceId string, page, pageSize int) ([]*Release, error)
	// GetReleasesInDateRange retrieves all releases within a specified date range.
	GetReleasesInDateRange(ctx context.Context, startDate, endDate time.Time, page, pageSize int) ([]*ServiceReleaseInfo, error)
}

// ReportRepository defines the methods for gathering reports.
type ReportRepository interface {
	// GetServiceRiskReport retrieves the risk report for a service.
	GetServiceRiskReport(ctx context.Context, serviceId string) (*ServiceRiskReport, error)
}
