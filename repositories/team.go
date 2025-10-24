package repositories

import (
	"errors"
	"time"
)

type Team struct {
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated,omitempty"`
	Id      string    `json:"id"`
}

func (t *Team) Validate() error {
	if t.Name == "" {
		return errors.New("team name is required")
	}
	return nil
}
