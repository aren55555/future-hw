package api

import (
	"errors"
	"net/url"
	"strconv"
	"time"
)

const (
	queryParamKeyTrainerID = "trainer_id"
	queryParamKeyStartsAt  = "starts_at"
	queryParamKeyEndsAt    = "ends_at"
)

func parseParams(q url.Values) requestParams {
	trainerID, _ := strconv.ParseInt(q.Get(queryParamKeyTrainerID), 10, 64)
	startsAt, _ := time.Parse(time.RFC3339, q.Get(queryParamKeyStartsAt))
	endsAt, _ := time.Parse(time.RFC3339, q.Get(queryParamKeyEndsAt))

	return requestParams{
		TrainerID: trainerID,
		StartsAt:  startsAt,
		EndsAt:    endsAt,
	}
}

type requestParams struct {
	TrainerID int64
	StartsAt  time.Time
	EndsAt    time.Time
}

func (rp requestParams) hasAll() bool {
	return rp.hasTrainerID() && rp.StartsAt != time.Time{} && rp.EndsAt != time.Time{}
}

func (rp requestParams) hasTrainerID() bool {
	return rp.TrainerID != 0
}

type createRequest struct {
	TrainerID int64     `json:"trainer_id"`
	UserID    int64     `json:"user_id"`
	StartsAt  time.Time `json:"starts_at"`
	EndsAt    time.Time `json:"ends_at"`
}

func (cr createRequest) validate() error {
	if cr.TrainerID == 0 || cr.UserID == 0 || cr.StartsAt == (time.Time{}) || cr.EndsAt == (time.Time{}) {
		return errors.New("invalid request: missing values")
	}
	if cr.EndsAt.Before(cr.StartsAt) {
		return errors.New("ends_at before starts_at")
	}
	return nil
}
