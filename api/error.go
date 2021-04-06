package api

import "net/http"

// TODO: make this return a formatted JSON error rather than just status text!
func statusError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
