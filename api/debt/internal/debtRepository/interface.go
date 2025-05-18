package debtRepository

import (
	"context"
	"errors"
	"service-dependency-api/internal"
)

type Repository interface {
	CreateDebtItem(ctx context.Context, debt Debt) error
	UpdateStatus(ctx context.Context, id, status string) error
}

type Debt struct {
	ServiceId   string `json:"serviceId"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (d *Debt) Validate() error {
	if d.Status == "" {
		d.Status = "pending"
	}
	if !internal.DebtTypes.IsMember(d.Type) {
		return errors.New("invalid debt type")
	}
	if !internal.DebtStatus.IsMember(d.Status) {
		return errors.New("invalid debt status")
	}
	if d.Title == "" {
		return errors.New("title is required")
	}
	return nil
}
