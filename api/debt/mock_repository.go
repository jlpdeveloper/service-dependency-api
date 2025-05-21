package debt

import (
	"context"
	"service-dependency-api/repositories"
)

type mockDebtRepository struct {
	Err   error
	Debts []*repositories.Debt
}

func (repo mockDebtRepository) CreateDebtItem(_ context.Context, _ repositories.Debt) error {
	if repo.Err != nil {
		return repo.Err
	}

	// If no error, we consider the operation successful
	return nil
}

func (repo mockDebtRepository) UpdateStatus(_ context.Context, _, _ string) error {
	if repo.Err != nil {
		return repo.Err
	}

	// If no error, we consider the operation successful
	return nil
}
