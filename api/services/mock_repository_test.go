package services

import (
	"context"
	"service-dependency-api/api/services/service_repository"
	"service-dependency-api/internal"
	"time"
)

type MockServiceRepository struct {
	Data func() []map[string]any
	Err  error
}

func (repo MockServiceRepository) CreateService(_ context.Context, _ service_repository.Service) (string, error) {
	if repo.Err != nil {
		return "", repo.Err
	}

	data := repo.Data()
	if len(data) == 0 {
		return "", nil
	}

	if id, ok := data[0]["id"]; ok {
		if idStr, ok := id.(string); ok {
			return idStr, nil
		}
	}

	return "", nil
}

func (repo MockServiceRepository) GetServiceById(_ context.Context, id string) (service_repository.Service, error) {
	if repo.Err != nil {
		return service_repository.Service{}, repo.Err
	}

	d := repo.Data()

	// Search for service with matching ID
	for _, v := range d {
		// Check if ID matches
		if serviceId, ok := v["id"]; ok {
			if idStr, ok := serviceId.(string); ok && idStr == id {
				service := service_repository.Service{}

				// Set the ID
				service.Id = idStr

				// Safely extract name with validation
				if name, ok := v["name"]; ok {
					if nameStr, ok := name.(string); ok {
						service.Name = nameStr
					}
				}

				// Safely extract description with validation
				if desc, ok := v["description"]; ok {
					if descStr, ok := desc.(string); ok {
						service.Description = descStr
					}
				}

				// Safely extract service type with validation
				if svcType, ok := v["type"]; ok {
					if typeStr, ok := svcType.(string); ok {
						service.ServiceType = typeStr
					}
				}

				// Safely extract created date with validation
				if date, ok := v["created"]; ok {
					if dateTime, ok := date.(time.Time); ok {
						service.Created = dateTime
					}
				}

				return service, nil
			}
		}
	}

	// Service not found
	return service_repository.Service{}, nil
}

func (repo MockServiceRepository) UpdateService(_ context.Context, service service_repository.Service) error {
	if repo.Err != nil {
		return repo.Err
	}

	d := repo.Data()

	// Search for service with matching ID
	for _, v := range d {
		// Check if ID matches
		if serviceId, ok := v["id"]; ok {
			if idStr, ok := serviceId.(string); ok && idStr == service.Id {
				// Service found, update would be successful
				return nil
			}
		}
	}

	// Service not found
	return &internal.HTTPError{
		Status: 404,
		Msg:    "Service not found",
	}
}

func (repo MockServiceRepository) GetAllServices(_ context.Context, page int, pageSize int) ([]service_repository.Service, error) {
	if repo.Err != nil {
		return []service_repository.Service{}, repo.Err
	}
	d := repo.Data()

	// Convert all data to Service objects
	var allServices []service_repository.Service
	for _, v := range d {
		service := service_repository.Service{}

		// Safely extract ID with validation
		if id, ok := v["id"]; ok {
			if idStr, ok := id.(string); ok {
				service.Id = idStr
			}
		}

		// Safely extract name with validation
		if name, ok := v["name"]; ok {
			if nameStr, ok := name.(string); ok {
				service.Name = nameStr
			}
		}

		// Safely extract description with validation
		if desc, ok := v["description"]; ok {
			if descStr, ok := desc.(string); ok {
				service.Description = descStr
			}
		}

		// Safely extract service type with validation
		if svcType, ok := v["type"]; ok {
			if typeStr, ok := svcType.(string); ok {
				service.ServiceType = typeStr
			}
		}

		// Safely extract created date with validation
		if date, ok := v["created"]; ok {
			if dateTime, ok := date.(time.Time); ok {
				service.Created = dateTime
			}
		}

		allServices = append(allServices, service)
	}

	// Apply pagination
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	// Handle edge cases
	if startIndex >= len(allServices) {
		return []service_repository.Service{}, nil
	}
	if endIndex > len(allServices) {
		endIndex = len(allServices)
	}

	return allServices[startIndex:endIndex], nil
}

func (repo MockServiceRepository) DeleteService(_ context.Context, id string) error {
	if repo.Err != nil {
		return repo.Err
	}

	d := repo.Data()

	// Search for service with matching ID
	for _, v := range d {
		// Check if ID matches
		if serviceId, ok := v["id"]; ok {
			if idStr, ok := serviceId.(string); ok && idStr == id {
				// Service found, delete would be successful
				return nil
			}
		}
	}

	// Service not found, but not returning an error as the delete operation
	// is idempotent - deleting a non-existent resource is not an error
	return nil
}
