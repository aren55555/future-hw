package compute

import (
	"errors"
	"time"

	"github.com/aren55555/future-hw/models"
)

func AvailableTimes(trainerID int64, start time.Time, end time.Time, existing []*models.Appointment) ([]*models.Appointment, error) {
	taken := map[time.Time]interface{}{}
	for _, appt := range existing {
		taken[appt.StartsAt] = true
	}
	return availableTimes(trainerID, start, end, taken)
}

func availableTimes(trainerID int64, start time.Time, end time.Time, taken map[time.Time]interface{}) ([]*models.Appointment, error) {
	if err := validateInputs(start, end); err != nil {
		return nil, err
	}

	availables := []*models.Appointment{}

	currentTime := start
	for currentTime.Before(end) {
		if isBizHours(currentTime) {
			if _, ok := taken[currentTime]; !ok {
				availables = append(availables, &models.Appointment{
					TrainerID: trainerID,
					StartsAt:  currentTime,
					EndsAt:    currentTime.Add(30 * time.Minute),
				})
			}
		}
		currentTime = currentTime.Add(30 * time.Minute)
	}

	return availables, nil
}

func validateInputs(start, end time.Time) error {
	if end.Before(start) {
		return errors.New("end before start")
	}
	if start.Second() != 0 {
		return errors.New("start must be :00 secs")
	}
	if start.Minute() != 0 && start.Minute() != 30 {
		return errors.New("start must be :00 or :30 mins")
	}
	if end.Second() != 0 {
		return errors.New("end must be :00 secs")
	}
	if end.Minute() != 0 && end.Minute() != 30 {
		return errors.New("end must be :00 or :30 mins")
	}
	return nil
}

func isBizHours(t time.Time) bool {
	// TODO: holidays?
	// TODO: DST?
	if t.Weekday() == time.Sunday || t.Weekday() == time.Saturday {
		return false
	}
	if t.Hour() < 8 || t.Hour() > 17 {
		return false
	}
	return true
}
