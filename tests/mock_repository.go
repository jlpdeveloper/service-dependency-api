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
	var serviceList []services.Service
	d := repo.Data()
	for _, v := range d {
		serviceList = append(serviceList, services.Service{
			Id:          v["id"].(string),
			Name:        v["name"].(string),
			Description: v["description"].(string),
			ServiceType: v["service_type"].(string),
		})
	}

	return serviceList, nil
}
