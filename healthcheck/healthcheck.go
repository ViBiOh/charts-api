package healthcheck

import (
	"fmt"
	"net/http"

	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/writer"
)

// App stores informations and secret of API
type App struct {
	handlers map[string]http.Handler
}

// NewApp creates new App from Handlers list
func NewApp(handlers map[string]http.Handler) *App {
	return &App{
		handlers: handlers,
	}
}

// Handler for Health request. Should be use with net/http
func (a *App) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			for url, handler := range a.handlers {
				fakeWriter := writer.ResponseWriter{}
				request, err := http.NewRequest(http.MethodGet, fmt.Sprintf(`%s/health`, url), nil)
				if err != nil {
					httputils.InternalServerError(w, fmt.Errorf(`Error while creating health request: %v`, err))
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
