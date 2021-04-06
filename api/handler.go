package api

import (
	"net/http"

	"github.com/aren55555/future-hw/data"
)

const (
	headerContentType = "Content-Type"
	mimeJSON          = "application/json"

	queryParamKeyFilter      = "filter"
	queryParamValueScheduled = "scheduled"
)

type Handler struct {
	ds data.Store
}

func New(ds data.Store) *Handler {
	return &Handler{
		ds: ds,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(headerContentType, mimeJSON) // TODO: the errors are currently textual, but will still have a JSON Content Type header.

	switch r.Method {
	case http.MethodGet:
		// TODO: consider making this a POST so it can take a JSON body of arguments
		//       rather than a GET with query params. Times will encode in a much
		//       more readable fashion in JSON. I only used GET since the task's
		//       requirements explicitly used the word "get".
		if r.URL.Query().Get(queryParamKeyFilter) == queryParamValueScheduled {
			h.getScheduledAppointments(w, r)
			return
		}
		// Any other value for the filter param will simply return the available Appointments.
		h.getAvailableAppointments(w, r)
	case http.MethodPost:
		h.scheduleAppointment(w, r)
	default:
		statusError(w, http.StatusMethodNotAllowed)
	}
}
