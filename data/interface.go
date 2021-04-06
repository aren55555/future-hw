package data

import (
	"time"

	"github.com/aren55555/future-hw/models"
)

type Store interface {
	AvailableAppointments(trainerID int64, startsAt, endsAt time.Time) []*models.Appointment
	CreateAppointment(trainerID, userID int64, startsAt, endsAt time.Time) (*models.Appointment, error)
	GetAppointmentsByTrainer(trainerID int64) []*models.Appointment
}
