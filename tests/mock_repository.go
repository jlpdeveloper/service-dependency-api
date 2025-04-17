package tests

import "service-dependency-api/api/services"

type MockServiceRepository struct {
	Data func() []map[string]any
	Err  error
}

func (repo MockServiceRepository) CreateService(service services.Service) (string, error) {
	if repo.Err != nil {
		return "", repo.Err
	}
	return repo.Data()[0]["id"].(string), nil
}

func (repo MockServiceRepository) GetAllServices(page int, pageSize int) ([]services.Service, error) {
	if repo.Err != nil {
		return []services.Service{}, repo.Err
	}
	d := repo.Data()

	// Convert all data to Service objects
	var allServices []services.Service
	for _, v := range d {
		allServices = append(allServices, services.Service{
			Id:          v["id"].(string),
			Name:        v["name"].(string),
			Description: v["description"].(string),
			ServiceType: v["type"].(string),
		})
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
