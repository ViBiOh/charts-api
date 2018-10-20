package healthcheck

import (
	"fmt"
	"net/http"

	"github.com/ViBiOh/httputils/pkg/httperror"
	"github.com/ViBiOh/httputils/pkg/writer"
	"github.com/pkg/errors"
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
func (a App) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		for url, handler := range a.handlers {
			fakeWriter := writer.ResponseWriter{}
			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf(`%s/health`, url), nil)
			if err != nil {
				httperror.InternalServerError(w, errors.WithStack(err))
				return
			}

			handler.ServeHTTP(&fakeWriter, request)

			if status := fakeWriter.Status(); status != http.StatusOK {
				http.Error(w, fmt.Sprintf(`Bad status while pinging endpoint %s`, url), status)
				return
			}
		}

		w.WriteHeader(http.StatusOK)

	})
}
