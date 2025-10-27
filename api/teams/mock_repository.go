package teams

import (
	"context"
	"service-dependency-api/repositories"
)

type mockTeamRepository struct {
	Err   error
	team  repositories.Team
	teams []repositories.Team
}

func (repo mockTeamRepository) GetTeam(_ context.Context, _ string) (*repositories.Team, error) {
	if repo.Err != nil {
		return nil, repo.Err
	}
	return &repo.team, nil
}

func (repo mockTeamRepository) GetTeams(_ context.Context, _ int, _ int) ([]repositories.Team, error) {
	if repo.Err != nil {
		return nil, repo.Err
	}
	return repo.teams, nil
}
func (repo mockTeamRepository) CreateTeam(_ context.Context, _ repositories.Team) (string, error) {
	if repo.Err != nil {
		return "", repo.Err
	}
	return "", nil
}
func (repo mockTeamRepository) UpdateTeam(_ context.Context, _ repositories.Team) error {
	if repo.Err != nil {
		return repo.Err
	}
	return nil
}

func (repo mockTeamRepository) DeleteTeam(_ context.Context, _ string) error {
	if repo.Err != nil {
		return repo.Err
	}
	return nil
}
