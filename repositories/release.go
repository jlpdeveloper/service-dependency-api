package repositories

import (
	"errors"
	"service-dependency-api/internal"
	"time"
)

type Release struct {
	ServiceId   string    `json:"service_id"`
	ReleaseDate time.Time `json:"release_date"`
	Url         string    `json:"url"`
	Version     string    `json:"version"`
}

func (r *Release) Validate() error {
	if _, ok := internal.IsValidGuid(r.ServiceId); !ok {
		return errors.New("invalid Service Id")
	}
	if r.ReleaseDate.IsZero() {
		r.ReleaseDate = time.Now().UTC()
	}
	if r.Url == "" && r.Version == "" {
		return errors.New("missing URL or Version")
	}
	return nil
}
