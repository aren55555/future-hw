package mock

import (
	"time"

	"github.com/aren55555/future-hw/data"
	"github.com/aren55555/future-hw/models"
)

var _ data.Store = &Client{}

// Client implements the data.Store interface, should be used in testing/stubbing.
type Client struct {
	Appointments []*models.Appointment
	Appointment  *models.Appointment
	Error        error
}

func (m *Client) AvailableAppointments(_ int64, _, _ time.Time) ([]*models.Appointment, error) {
	return m.Appointments, m.Error
}

func (m *Client) GetAppointmentsByTrainer(_ int64) []*models.Appointment {
	return m.Appointments
}

func (m *Client) CreateAppointment(trainerID, userID int64, startsAt, endsAt time.Time) (*models.Appointment, error) {
	return m.Appointment, m.Error
}
