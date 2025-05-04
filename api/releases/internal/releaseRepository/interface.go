package releaseRepository

import (
	"context"
	"errors"
	"service-dependency-api/internal"
	"time"
)

type ReleaseRepository interface {
	CreateRelease(ctx context.Context, release Release) error
}

type Release struct {
	ServiceId   string    `json:"serviceId"`
	Url         string    `json:"url"`
	ReleaseDate time.Time `json:"releaseDate"`
}

func (r *Release) Validate() error {
	if _, ok := internal.IsValidGuid(r.ServiceId); !ok {
		return errors.New("Invalid Service Id")
	}
	r.ReleaseDate = time.Now().UTC()
	return nil
}
