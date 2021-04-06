package api

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) getAvailableAppointments(w http.ResponseWriter, r *http.Request) {
	params := parseParams(r.URL.Query())
	if !params.hasAll() {
		statusError(w, http.StatusBadRequest)
		return
	}

	appointments := h.ds.AvailableAppointments(params.TrainerID, params.StartsAt, params.EndsAt)

	json.NewEncoder(w).Encode(appointments)
}

func (h *Handler) getScheduledAppointments(w http.ResponseWriter, r *http.Request) {
	params := parseParams(r.URL.Query())
	if !params.hasTrainerID() {
		statusError(w, http.StatusBadRequest)
		return
	}

	appointments := h.ds.GetAppointmentsByTrainer(params.TrainerID)

	json.NewEncoder(w).Encode(appointments)
}

func (h *Handler) scheduleAppointment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	request := createRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		statusError(w, http.StatusBadRequest)
		return
	}

	if err := request.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdAppt, err := h.ds.CreateAppointment(request.TrainerID, request.UserID, request.StartsAt, request.EndsAt)
	if err != nil {
		statusError(w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdAppt)
}
