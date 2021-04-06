package models

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

const (
	jsonPayload = `[
	{"id":1,"trainer_id":1,"starts_at":"2020-01-24T09:00:00-08:00","ends_at":"2020-01-24T09:30:00-08:00"},
	{"id":2,"trainer_id":1,"starts_at":"2020-01-24T10:00:00-08:00","ends_at":"2020-01-24T10:30:00-08:00"},
	{"id":3,"trainer_id":1,"starts_at":"2020-01-25T10:00:00-08:00","ends_at":"2020-01-25T10:30:00-08:00"},
	{"id":4,"trainer_id":1,"starts_at":"2020-01-25T10:30:00-08:00","ends_at":"2020-01-25T11:00:00-08:00"},
	{"id":5,"trainer_id":1,"starts_at":"2020-01-26T10:00:00-08:00","ends_at":"2020-01-26T10:30:00-08:00"},
	{"id":6,"trainer_id":2,"starts_at":"2020-01-24T09:00:00-08:00","ends_at":"2020-01-24T09:30:00-08:00"},
	{"id":7,"trainer_id":2,"starts_at":"2020-01-26T10:00:00-08:00","ends_at":"2020-01-26T10:30:00-08:00"},
	{"id":8,"trainer_id":3,"starts_at":"2020-01-26T12:00:00-08:00","ends_at":"2020-01-26T12:30:00-08:00"},
	{"id":9,"trainer_id":3,"starts_at":"2020-01-26T13:00:00-08:00","ends_at":"2020-01-26T13:30:00-08:00"},
	{"id":10,"trainer_id":3,"starts_at":"2020-01-26T14:00:00-08:00","ends_at":"2020-01-26T14:30:00-08:00"}
]`
)

func must(t time.Time, err error) time.Time {
	if err != nil {
		panic(err)
	}
	return t
}

// Verifies I've modelled the JSON correctly.
func TestParse(t *testing.T) {
	appointments := []Appointment{}
	if err := json.Unmarshal([]byte(jsonPayload), &appointments); err != nil {
		t.Fatal(err)
	}
	if got, want := len(appointments), 10; got != want {
		t.Fatalf("JSON appointmnets length, got %v, want %v", got, want)
	}
	if got, want := appointments[0], (Appointment{
		ID:        1,
		TrainerID: 1,
		StartsAt:  must(time.Parse(time.RFC3339, "2020-01-24T09:00:00-08:00")),
		EndsAt:    must(time.Parse(time.RFC3339, "2020-01-24T09:30:00-08:00")),
	}); !reflect.DeepEqual(got, want) {
		t.Fatalf("expected to be the same: got %v, want %v", got, want)
	}
}
