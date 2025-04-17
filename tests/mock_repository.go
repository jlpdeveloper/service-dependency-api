package tests

import "service-dependency-api/api/services"

type MockServiceRepository struct {
	Data map[string]any
	Err  error
}

func (repo MockServiceRepository) CreateService(service services.Service) (string, error) {
	if repo.Err != nil {
		return "", repo.Err
	}
	return repo.Data["id"].(string), nil
}
