package models

import "time"

// Models the JSON representation of a Client's Appointment with a Trainer.
type Appointment struct {
	ID        int64     `json:"id"`
	StartsAt  time.Time `json:"starts_at"`
	EndsAt    time.Time `json:"ends_at"`
	TrainerID int64     `json:"trainer_id"`
}
