package mem

import (
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	c := New()

	c.CreateAppointment(1, 1, time.Now(), time.Now())
	c.CreateAppointment(1, 1, time.Now(), time.Now())
	c.CreateAppointment(1, 1, time.Now(), time.Now())
	c.CreateAppointment(2, 1, time.Now(), time.Now())

	if got, want := len(c.GetAppointmentsByTrainer(1)), 3; got != want {
		t.Fatalf("appointments len: got %v, want %v", got, want)
	}
	if got, want := len(c.GetAppointmentsByTrainer(2)), 1; got != want {
		t.Fatalf("appointments len: got %v, want %v", got, want)
	}
	if got, want := len(c.GetAppointmentsByTrainer(3)), 0; got != want {
		t.Fatalf("appointments len: got %v, want %v", got, want)
	}
}
