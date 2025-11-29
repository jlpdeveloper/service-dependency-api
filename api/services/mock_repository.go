package services

import (
	"context"
	"service-atlas/internal/customerrors"
	"service-atlas/repositories"
	"time"
)

type mockServiceRepository struct {
	Data func() []map[string]any
	Err  error
}

func (repo mockServiceRepository) CreateService(_ context.Context, _ repositories.Service) (string, error) {
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

func (repo mockServiceRepository) GetServiceById(_ context.Context, id string) (repositories.Service, error) {
	if repo.Err != nil {
		return repositories.Service{}, repo.Err
	}

	d := repo.Data()

	// Search for service with matching ID
	for _, v := range d {
		// Check if ID matches
		if serviceId, ok := v["id"]; ok {
			if idStr, ok := serviceId.(string); ok && idStr == id {
				service := repositories.Service{}

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
	return repositories.Service{}, nil
}

func (repo mockServiceRepository) UpdateService(_ context.Context, service repositories.Service) error {
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
	return &customerrors.HTTPError{
		Status: 404,
		Msg:    "Service not found",
	}
}

func (repo mockServiceRepository) GetAllServices(_ context.Context, page int, pageSize int) ([]repositories.Service, error) {
	if repo.Err != nil {
		return []repositories.Service{}, repo.Err
	}
	d := repo.Data()

	// Convert all data to Service objects
	var allServices []repositories.Service
	for _, v := range d {
		service := repositories.Service{}

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
		return []repositories.Service{}, nil
	}
	if endIndex > len(allServices) {
		endIndex = len(allServices)
	}

	return allServices[startIndex:endIndex], nil
}

func (repo mockServiceRepository) DeleteService(_ context.Context, id string) error {
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
