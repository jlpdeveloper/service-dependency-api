package repositories

import (
	"errors"
	"service-dependency-api/internal"
)

type Debt struct {
	ServiceId   string `json:"serviceId"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Id          string `json:"id"`
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
