package healthcheck

import (
	"net/http"
)

const healthcheckPath = `/health`

var handlersToCheck []http.Handler

// Init charts handler
func Init(handlers []http.Handler) error {
	handlersToCheck = handlers

	return nil
}

// Handler for Health request. Should be use with net/http
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
