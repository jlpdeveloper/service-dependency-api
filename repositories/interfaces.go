package repositories

import "context"

type DebtRepository interface {
	CreateDebtItem(ctx context.Context, debt Debt) error
	UpdateStatus(ctx context.Context, id, status string) error
}
