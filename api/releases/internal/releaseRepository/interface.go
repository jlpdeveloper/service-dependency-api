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
	ServiceId   string    `json:"service_id"`
	ReleaseDate time.Time `json:"release_date"`
	Url         string    `json:"url"`
}

func (r *Release) Validate() error {
	if _, ok := internal.IsValidGuid(r.ServiceId); !ok {
		return errors.New("invalid Service Id")
	}
	if r.ReleaseDate.IsZero() {
		r.ReleaseDate = time.Now().UTC()
	}
	return nil
}
