package readings

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/ViBiOh/auth/pkg/auth"
	"github.com/ViBiOh/auth/pkg/model"
	"github.com/ViBiOh/httputils/pkg/db"
	"github.com/ViBiOh/httputils/pkg/errors"
	"github.com/ViBiOh/httputils/pkg/httperror"
	"github.com/ViBiOh/httputils/pkg/httpjson"
)

const healthcheckPath = `/health`

// App stores informations and secret of API
type App struct {
	db *sql.DB
}

// NewApp creates new App from Flags' config
func NewApp(db *sql.DB) *App {
	return &App{
		db: db,
	}
}

// Flags add flags for given prefix
func Flags(prefix string) map[string]*string {
	return nil
}

func (a App) listReadings(w http.ResponseWriter, r *http.Request, user *model.User) {
	if list, err := a.listReadingsOfUser(user); err == nil {
		if err := httpjson.ResponseArrayJSON(w, http.StatusOK, list, httpjson.IsPretty(r)); err != nil {
			httperror.InternalServerError(w, err)
		}
	} else {
		httperror.InternalServerError(w, err)
	}
}

// Handler for Readings request. Should be use with net/http
func (a App) Handler() http.Handler {
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

		user := auth.UserFromContext(r.Context())
		if user == nil {
			httperror.InternalServerError(w, errors.New(`no user provided`))
		}

		if strings.HasPrefix(r.URL.Path, tagsPath) {
			a.tagsHandler(w, r, user, strings.TrimPrefix(r.URL.Path, tagsPath))
		} else if r.Method == http.MethodGet && (r.URL.Path == `/` || r.URL.Path == ``) {
			a.listReadings(w, r, user)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
