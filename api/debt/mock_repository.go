package debt

import (
	"context"
	"service-atlas/repositories"
)

type mockDebtRepository struct {
	Err   error
	Debts []repositories.Debt
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

func (repo mockDebtRepository) GetDebtByServiceId(_ context.Context, id string, _, _ int, onlyResolved bool) ([]repositories.Debt, error) {
	if repo.Err != nil {
		return nil, repo.Err
	}
	var debts []repositories.Debt
	for _, d := range repo.Debts {
		if d.ServiceId == id {
			if onlyResolved {
				if d.Status == "remediated" {
					debts = append(debts, d)
				}
			} else {
				debts = append(debts, d)
			}
		}
	}
	return debts, nil

}
