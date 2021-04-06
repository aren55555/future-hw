package compute

import (
	"testing"
	"time"

	"github.com/aren55555/future-hw/helpers"
)

const (
	trainerID = int64(1)
)

func TestAvailableTimes(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		start    string
		end      string
		takenSet map[time.Time]interface{}
		wantErr  bool
		wantLen  int
	}{
		{
			desc:     "invalid_mins",
			start:    "2020-01-24T07:25:55-08:00",
			end:      "2020-01-24T10:30:00-08:00",
			takenSet: map[time.Time]interface{}{},
			wantErr:  true,
		},
		{
			desc:     "invalid_secs",
			start:    "2020-01-24T07:30:55-08:00",
			end:      "2020-01-24T10:30:00-08:00",
			takenSet: map[time.Time]interface{}{},
			wantErr:  true,
		},
		{
			desc:     "none_taken",
			start:    "2020-01-24T07:00:00-08:00",
			end:      "2020-01-24T10:30:00-08:00",
			takenSet: map[time.Time]interface{}{},
			wantLen:  5,
		},
		{
			desc:  "two_slots_taken",
			start: "2020-01-24T07:00:00-08:00",
			end:   "2020-01-24T10:30:00-08:00",
			takenSet: map[time.Time]interface{}{
				helpers.TimeMust(time.Parse(time.RFC3339, "2020-01-24T08:30:00-08:00")): true,
				helpers.TimeMust(time.Parse(time.RFC3339, "2020-01-24T09:00:00-08:00")): true,
			},
			wantLen: 3,
		},
		{
			desc:     "outside_biz",
			start:    "2020-01-24T00:00:00-08:00",
			end:      "2020-01-24T07:30:00-08:00",
			takenSet: map[time.Time]interface{}{},
			wantLen:  0,
		},
		{
			desc:     "weekend",
			start:    "2020-01-25T08:00:00-08:00",
			end:      "2020-01-25T10:30:00-08:00",
			takenSet: map[time.Time]interface{}{},
			wantLen:  0,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			// I know I'm testing the private function; but I don't want to rewrite this since I'm done with my allocated time box.
			available, err := availableTimes(
				trainerID,
				helpers.TimeMust(time.Parse(time.RFC3339, tc.start)),
				helpers.TimeMust(time.Parse(time.RFC3339, tc.end)),
				tc.takenSet,
			)
			if got, want := (err != nil), tc.wantErr; got != want {
				t.Fatalf("AvailableTimes error: got %v, want %v", got, want)
			}
			if got, want := len(available), tc.wantLen; got != want {
				t.Fatalf("len: got %v, want %v", got, want)
			}
		})
	}
}
