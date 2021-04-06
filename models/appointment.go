package models

import (
	"errors"
	"time"
)

// Models the JSON representation of a Client's Appointment with a Trainer.
type Appointment struct {
	ID        int64     `json:"id,omitempty"`
	TrainerID int64     `json:"trainer_id,omitempty"`
	StartsAt  time.Time `json:"starts_at,omitempty"`
	EndsAt    time.Time `json:"ends_at,omitempty"`
}

// Validate the presence of the fields (are non zero)
func (a *Appointment) Validate() error {
	if a.ID == 0 || a.TrainerID == 0 || a.StartsAt == (time.Time{}) || a.EndsAt == (time.Time{}) {
		return errors.New("invalid appointment: missing values")
	}
	if a.EndsAt.Before(a.StartsAt) {
		return errors.New("ends_at before starts_at")
	}
	return nil
}
