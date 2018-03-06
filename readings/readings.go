package readings

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/ViBiOh/auth/auth"
	authProvider "github.com/ViBiOh/auth/provider"
	"github.com/ViBiOh/httputils/db"
	"github.com/ViBiOh/httputils/httperror"
	"github.com/ViBiOh/httputils/httpjson"
)

const healthcheckPath = `/health`

// App stores informations and secret of API
type App struct {
	db      *sql.DB
	authApp *auth.App
}

// NewApp creates new App from Flags' config
func NewApp(db *sql.DB, authApp *auth.App) *App {
	return &App{
		db:      db,
		authApp: authApp,
	}
}

// Flags add flags for given prefix
func Flags(prefix string) map[string]*string {
	return nil
}

func (a *App) listReadings(w http.ResponseWriter, r *http.Request, user *authProvider.User) {
	if list, err := a.listReadingsOfUser(user); err == nil {
		if err := httpjson.ResponseArrayJSON(w, http.StatusOK, list, httpjson.IsPretty(r.URL.RawQuery)); err != nil {
			httperror.InternalServerError(w, err)
		}
	} else {
		httperror.InternalServerError(w, err)
	}
}

// Handler for Readings request. Should be use with net/http
func (a *App) Handler() http.Handler {
	authHandler := a.authApp.Handler(func(w http.ResponseWriter, r *http.Request, user *authProvider.User) {
		if strings.HasPrefix(r.URL.Path, tagsPath) {
			a.tagsHandler(w, r, user, strings.TrimPrefix(r.URL.Path, tagsPath))
		} else if r.Method == http.MethodGet && (r.URL.Path == `/` || r.URL.Path == ``) {
			a.listReadings(w, r, user)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method == http.MethodGet && r.URL.Path == healthcheckPath {
			if db.Ping(a.db) {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusServiceUnavailable)
			}
			return
		}

		authHandler.ServeHTTP(w, r)
	})
}
