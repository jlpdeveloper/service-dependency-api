package serviceRepository

import (
	"context"
	"errors"
	"net/url"
	"time"
)

type Service struct {
	Id          string    `json:"id,omitempty"`
	Name        string    `json:"name"`
	ServiceType string    `json:"type"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated,omitempty"`
	Url         string    `json:"url,omitempty"`
}

func (service *Service) Validate() error {
	switch {
	case service.Name == "":
		return errors.New("service name is required")
	case service.Url == "":
		return errors.New("service url is required")
	case service.ServiceType == "":
		return errors.New("service type is required")
	}

	// Validate URL format
	parsedURL, err := url.Parse(service.Url)
	if err != nil {
		return errors.New("service url is not a valid URL format")
	}

	// Ensure URL has a scheme (http or https)
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("service url must use http or https protocol")
	}

	return nil
}

type ServiceRepository interface {
	GetAllServices(ctx context.Context, page int, pageSize int) ([]Service, error)
	CreateService(ctx context.Context, service Service) (string, error)
	UpdateService(ctx context.Context, service Service) error
	DeleteService(ctx context.Context, id string) error
	GetServiceById(ctx context.Context, id string) (Service, error)
}
