package healthcheck

import (
	"fmt"
	"net/http"

	"github.com/ViBiOh/httputils/writer"
)

var handlersToCheck []http.Handler
var healthRequest *http.Request

// Init charts handler
func Init(handlers []http.Handler) (err error) {
	handlersToCheck = handlers

	healthRequest, err = http.NewRequest(http.MethodGet, `/health`, nil)
	if err != nil {
		err = fmt.Errorf(`Error while creating health request: %v`, err)
	}

	return
}

// Handler for Health request. Should be use with net/http
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			for index, handler := range handlersToCheck {
				fakeWriter := writer.ResponseWriter{}

				handler.ServeHTTP(&fakeWriter, healthRequest)

				if status := fakeWriter.Status(); status != http.StatusOK {
					w.WriteHeader(status)
					return
				}
			}

			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
