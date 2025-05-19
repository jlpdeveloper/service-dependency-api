package debt

import (
	"context"
	"service-dependency-api/api/debt/internal/debtRepository"
)

type mockDebtRepository struct {
	Err   error
	Debts []*debtRepository.Debt
}

func (repo mockDebtRepository) CreateDebtItem(_ context.Context, _ debtRepository.Debt) error {
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
