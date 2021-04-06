package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/aren55555/future-hw/data"
	"github.com/aren55555/future-hw/data/mock"
	"github.com/aren55555/future-hw/helpers"
	"github.com/aren55555/future-hw/models"
)

var (
	startsAt = helpers.TimeMust(time.Parse(time.RFC3339, "2020-01-24T09:00:00-08:00"))
	endsAt   = helpers.TimeMust(time.Parse(time.RFC3339, "2020-01-24T09:30:00-08:00"))
)

type responseValidator func(*http.Response) error

func appointmentListResponseValidationFactory(expectedLength int) responseValidator {
	return responseValidator(func(resp *http.Response) error {
		defer resp.Body.Close()
		appointments := []models.Appointment{}
		if err := json.NewDecoder(resp.Body).Decode(&appointments); err != nil {
			return err
		}
		if got, want := len(appointments), expectedLength; got != want {
			return fmt.Errorf("appoinments JSON length: got %v, want %v", got, want)
		}
		return nil
	})
}

func TestGetAvailableAppointments(t *testing.T) {
	for _, tc := range []struct {
		desc                 string
		query                url.Values
		mock                 data.Store
		expectedStatus       int
		responseVerification func(*http.Response) error
	}{
		{
			desc:  "invalid_params",
			query: url.Values{ /* there are no params */ },
			mock: &mock.Client{
				AppointmentsReply: []*models.Appointment{},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			desc: "no_available_appointments",
			query: url.Values{
				queryParamKeyTrainerID: []string{"1"},
				queryParamKeyStartsAt:  []string{"2020-01-24T09:00:00-08:00"},
				queryParamKeyEndsAt:    []string{"2020-01-24T09:30:00-08:00"},
			},
			mock: &mock.Client{
				AppointmentsReply: []*models.Appointment{},
			},
			expectedStatus:       http.StatusOK,
			responseVerification: appointmentListResponseValidationFactory(0),
		},
		{
			desc: "has_available_appointments",
			query: url.Values{
				queryParamKeyTrainerID: []string{"1"},
				queryParamKeyStartsAt:  []string{"2020-01-24T09:00:00-08:00"},
				queryParamKeyEndsAt:    []string{"2020-01-24T09:30:00-08:00"},
			},
			mock: &mock.Client{
				AppointmentsReply: []*models.Appointment{
					{
						ID:        1,
						StartsAt:  startsAt,
						EndsAt:    endsAt,
						TrainerID: 1,
					},
					{
						ID:        2,
						StartsAt:  startsAt,
						EndsAt:    endsAt,
						TrainerID: 2,
					},
				},
			},
			expectedStatus:       http.StatusOK,
			responseVerification: appointmentListResponseValidationFactory(2),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			h := New(tc.mock)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?%s", tc.query.Encode()), nil)

			rr := httptest.NewRecorder()

			h.ServeHTTP(rr, req)

			// Verify HTTP Status Code
			if got, want := rr.Result().StatusCode, tc.expectedStatus; got != want {
				t.Fatalf("http status got %v, want %v", got, want)
			}
			if tc.responseVerification != nil {
				if err := tc.responseVerification(rr.Result()); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestGetScheduledAppointments(t *testing.T) {
	for _, tc := range []struct {
		desc                 string
		query                url.Values
		mock                 data.Store
		expectedStatus       int
		responseVerification func(*http.Response) error
	}{
		{
			desc: "invalid_params",
			query: url.Values{
				queryParamKeyFilter: []string{queryParamValueScheduled},
				/* there are no params */
			},
			mock: &mock.Client{
				AppointmentsReply: []*models.Appointment{},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			desc: "no_scheduled_appointments",
			query: url.Values{
				queryParamKeyFilter:    []string{queryParamValueScheduled},
				queryParamKeyTrainerID: []string{"1"},
			},
			mock: &mock.Client{
				AppointmentsReply: []*models.Appointment{},
			},
			expectedStatus:       http.StatusOK,
			responseVerification: appointmentListResponseValidationFactory(0),
		},
		{
			desc: "has_scheduled_appointments",
			query: url.Values{
				queryParamKeyFilter:    []string{queryParamValueScheduled},
				queryParamKeyTrainerID: []string{"1"},
			},
			mock: &mock.Client{
				AppointmentsReply: []*models.Appointment{
					{
						ID:        1,
						StartsAt:  startsAt,
						EndsAt:    endsAt,
						TrainerID: 1,
					}, {
						ID:        2,
						StartsAt:  startsAt,
						EndsAt:    endsAt,
						TrainerID: 2,
					},
				},
			},
			expectedStatus:       http.StatusOK,
			responseVerification: appointmentListResponseValidationFactory(2),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			h := New(tc.mock)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?%s", tc.query.Encode()), nil)

			rr := httptest.NewRecorder()

			h.ServeHTTP(rr, req)

			// Verify HTTP Status Code
			if got, want := rr.Result().StatusCode, tc.expectedStatus; got != want {
				t.Fatalf("http status got %v, want %v", got, want)
			}
			if tc.responseVerification != nil {
				if err := tc.responseVerification(rr.Result()); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestScheduleAppointment(t *testing.T) {
	for _, tc := range []struct {
		desc                 string
		request              createRequest
		mock                 data.Store
		expectedStatus       int
		responseVerification func(*http.Response) error
	}{
		{
			desc:           "invalid_request",
			request:        createRequest{ /* no fields are present! */ },
			expectedStatus: http.StatusBadRequest,
		},
		{
			desc: "ends_at_before_starts_at",
			request: createRequest{
				TrainerID: 1,
				UserID:    1,
				StartsAt:  endsAt,
				EndsAt:    startsAt,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			desc: "data_store_returns_error",
			request: createRequest{
				TrainerID: 1,
				UserID:    1,
				StartsAt:  startsAt,
				EndsAt:    endsAt,
			},
			mock: &mock.Client{
				CreateReply: struct {
					Appointment *models.Appointment
					Error       error
				}{nil, errors.New("timeout")},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			desc: "success",
			request: createRequest{
				TrainerID: 1,
				UserID:    1,
				StartsAt:  startsAt,
				EndsAt:    endsAt,
			},
			mock: &mock.Client{
				CreateReply: struct {
					Appointment *models.Appointment
					Error       error
				}{&models.Appointment{
					ID:        1,
					TrainerID: 1,
					StartsAt:  startsAt,
					EndsAt:    endsAt,
				}, nil},
			},
			expectedStatus: http.StatusCreated,
			responseVerification: func(resp *http.Response) error {
				defer resp.Body.Close()
				appointment := models.Appointment{}
				if err := json.NewDecoder(resp.Body).Decode(&appointment); err != nil {
					return err
				}
				if got, want := appointment.ID, int64(1); got != want {
					return fmt.Errorf("appointment id: got %v, want %v", got, want)
				}
				return nil
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			jsonData, _ := json.Marshal(tc.request)

			h := New(tc.mock)
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))

			rr := httptest.NewRecorder()

			h.ServeHTTP(rr, req)

			// Verify HTTP Status Code
			if got, want := rr.Result().StatusCode, tc.expectedStatus; got != want {
				t.Fatalf("http status got %v, want %v", got, want)
			}
			if tc.responseVerification != nil {
				if err := tc.responseVerification(rr.Result()); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
