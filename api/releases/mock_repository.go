package releases

import (
	"context"
	"service-dependency-api/api/releases/internal/releaseRepository"
	"time"
)

type mockReleaseRepository struct {
	Err         error
	Releases    []*releaseRepository.Release
	ServiceInfo []*releaseRepository.ServiceReleaseInfo
}

func (repo mockReleaseRepository) CreateRelease(_ context.Context, _ releaseRepository.Release) error {
	if repo.Err != nil {
		return repo.Err
	}

	// If no error, we consider the operation successful
	return nil
}

func (repo mockReleaseRepository) GetReleasesByServiceId(_ context.Context, _ string, page, pageSize int) ([]*releaseRepository.Release, error) {
	if repo.Err != nil {
		return nil, repo.Err
	}

	// Calculate start and end indices for pagination
	start := (page - 1) * pageSize
	end := start + pageSize

	// Check if start is beyond the array length
	if start >= len(repo.Releases) {
		return []*releaseRepository.Release{}, nil
	}

	// Ensure end doesn't exceed array length
	if end > len(repo.Releases) {
		end = len(repo.Releases)
	}

	// Return the paginated mock releases
	return repo.Releases[start:end], nil
}

func (repo mockReleaseRepository) GetReleasesInDateRange(_ context.Context, _ time.Time, _ time.Time, page, pageSize int) ([]*releaseRepository.ServiceReleaseInfo, error) {
	if repo.Err != nil {
		return nil, repo.Err
	}
	// Calculate start and end indices for pagination
	start := (page - 1) * pageSize
	end := start + pageSize

	// Check if start is beyond the array length
	if start >= len(repo.ServiceInfo) {
		return []*releaseRepository.ServiceReleaseInfo{}, nil
	}

	// Ensure end doesn't exceed array length
	if end > len(repo.ServiceInfo) {
		end = len(repo.ServiceInfo)
	}

	// Return the paginated mock releases
	return repo.ServiceInfo[start:end], nil
}
