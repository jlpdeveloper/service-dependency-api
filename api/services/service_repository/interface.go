package service_repository

import (
	"context"
	"time"
)

type Service struct {
	Id          string    `json:"id,omitempty"`
	Name        string    `json:"name"`
	ServiceType string    `json:"type"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated,omitempty"`
}

type ServiceRepository interface {
	GetAllServices(ctx context.Context, page int, pageSize int) ([]Service, error)
	CreateService(ctx context.Context, service Service) (string, error)
	UpdateService(ctx context.Context, service Service) (bool, error)
	DeleteService(ctx context.Context, id string) error
	GetServiceById(ctx context.Context, id string) (Service, error)
}
