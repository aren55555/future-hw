package mem

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/aren55555/future-hw/data"
	"github.com/aren55555/future-hw/models"
)

var _ data.Store = &Client{}

type Client struct {
	client
}

func New() *Client {
	return &Client{
		client{
			nextID:    1,
			byTrainer: map[int64][]*models.Appointment{},
		},
	}
}

func (c *Client) Seed(fileLocation string) error {
	if fileLocation == "" {
		return errors.New("no file location specified")
	}

	jsonData, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		return err
	}

	seedAppointments := []*models.Appointment{}
	if err := json.Unmarshal(jsonData, &seedAppointments); err != nil {
		return err
	}

	c.seed(seedAppointments)
	return nil
}

type client struct {
	lock      sync.RWMutex
	nextID    int64
	byTrainer map[int64][]*models.Appointment
}

func (c *client) seed(appointments []*models.Appointment) {
	c.lock.Lock()
	defer c.lock.Unlock()

	maxID := int64(0)
	for _, a := range appointments {
		if a.ID >= maxID {
			maxID = a.ID
		}

		var trainerAppointments []*models.Appointment
		var ok bool
		if appointments, ok = c.byTrainer[a.TrainerID]; !ok {
			trainerAppointments = []*models.Appointment{}
		}

		trainerAppointments = append(trainerAppointments, a)
		c.byTrainer[a.TrainerID] = trainerAppointments
	}

	c.nextID = maxID + 1
}

func (c *Client) AvailableAppointments(trainerID int64, startsAt, endsAt time.Time) []*models.Appointment {
	scheduled := c.GetAppointmentsByTrainer(trainerID)
	// TODO: compute availability by walking and computing the leftover aviailable ranges
	fmt.Println(scheduled)
	return []*models.Appointment{}
}

func (c *Client) GetAppointmentsByTrainer(trainerID int64) []*models.Appointment {
	c.lock.RLock()
	defer c.lock.RUnlock()

	var appointments []*models.Appointment
	var ok bool
	if appointments, ok = c.byTrainer[trainerID]; !ok {
		appointments = []*models.Appointment{}
	}

	return appointments
}

func (c *Client) CreateAppointment(trainerID, _ int64, startsAt, endsAt time.Time) (*models.Appointment, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var appointments []*models.Appointment
	var ok bool
	if appointments, ok = c.byTrainer[trainerID]; !ok {
		appointments = []*models.Appointment{}
	}

	newAppointment := &models.Appointment{
		ID:        c.nextID,
		TrainerID: trainerID,
		StartsAt:  startsAt,
		EndsAt:    endsAt,
	}
	c.nextID++

	appointments = append(appointments, newAppointment)
	c.byTrainer[trainerID] = appointments

	return newAppointment, nil
}
