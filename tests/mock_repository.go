package tests

import (
	"service-dependency-api/api/services"
	"time"
)

type MockServiceRepository struct {
	Data func() []map[string]any
	Err  error
}

func (repo MockServiceRepository) CreateService(_ services.Service) (string, error) {
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

func (repo MockServiceRepository) GetServiceById(id string) (services.Service, error) {
	if repo.Err != nil {
		return services.Service{}, repo.Err
	}

	d := repo.Data()

	// Search for service with matching ID
	for _, v := range d {
		// Check if ID matches
		if serviceId, ok := v["id"]; ok {
			if idStr, ok := serviceId.(string); ok && idStr == id {
				service := services.Service{}

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

				// Safely extract createdDate with validation
				if date, ok := v["createdDate"]; ok {
					if dateTime, ok := date.(time.Time); ok {
						service.CreatedDate = dateTime
					}
				}

				return service, nil
			}
		}
	}

	// Service not found
	return services.Service{}, nil
}

func (repo MockServiceRepository) GetAllServices(page int, pageSize int) ([]services.Service, error) {
	if repo.Err != nil {
		return []services.Service{}, repo.Err
	}
	d := repo.Data()

	// Convert all data to Service objects
	var allServices []services.Service
	for _, v := range d {
		service := services.Service{}

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

		// Safely extract createdDate with validation
		if date, ok := v["createdDate"]; ok {
			if dateTime, ok := date.(time.Time); ok {
				service.CreatedDate = dateTime
			}
		}

		allServices = append(allServices, service)
	}

	// Apply pagination
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	// Handle edge cases
	if startIndex >= len(allServices) {
		return []services.Service{}, nil
	}
	if endIndex > len(allServices) {
		endIndex = len(allServices)
	}

	return allServices[startIndex:endIndex], nil
}
