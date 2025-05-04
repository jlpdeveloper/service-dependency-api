package releases

import (
	"context"
	"service-dependency-api/api/releases/internal/releaseRepository"
)

type mockReleaseRepository struct {
	Err error
}

func (repo mockReleaseRepository) CreateRelease(_ context.Context, _ releaseRepository.Release) error {
	if repo.Err != nil {
		return repo.Err
	}

	// If no error, we consider the operation successful
	return nil
}
