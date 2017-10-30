package healthcheck

import (
	"fmt"
	"net/http"

	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/writer"
)

var handlers map[string]http.Handler

// Init charts handler
func Init(handlersToCheck map[string]http.Handler) (err error) {
	handlers = handlersToCheck

	return
}

// Handler for Health request. Should be use with net/http
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			for url, handler := range handlers {
				fakeWriter := writer.ResponseWriter{}
				request, err := http.NewRequest(http.MethodGet, url+`/health`, nil)
				if err != nil {
					httputils.InternalServer(w, fmt.Errorf(`Error while creating health request: %v`, err))
					return
				}

				handler.ServeHTTP(&fakeWriter, request)

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
