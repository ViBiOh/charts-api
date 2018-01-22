package readings

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/ViBiOh/auth/auth"
	authProvider "github.com/ViBiOh/auth/provider"
	"github.com/ViBiOh/auth/provider/basic"
	authService "github.com/ViBiOh/auth/service"
	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/db"
)

const healthcheckPath = `/health`

var (
	authConfig        = auth.Flags(`readingsAuth`)
	authServiceConfig = authService.Flags(`readings`)
	authBasicConfig   = basic.Flags(`readingsBasic`)

	dbConfig   = db.Flags(`readingsDb`)
	readingsDB *sql.DB
)

// Init readings API
func Init() (err error) {
	readingsDB, err = db.GetDB(dbConfig)
	if err != nil {
		err = fmt.Errorf(`Error while initializing database: %v`, err)
	}

	return
}

func listReadings(w http.ResponseWriter, r *http.Request, user *authProvider.User) {
	if list, err := listReadingsOfUser(user); err == nil {
		httputils.ResponseArrayJSON(w, http.StatusOK, list, httputils.IsPretty(r.URL.RawQuery))
	} else {
		httputils.InternalServerError(w, err)
	}
}

// Handler for Readings request. Should be use with net/http
func Handler() http.Handler {
	authApp := auth.NewApp(authConfig, authService.NewApp(authServiceConfig, authBasicConfig, nil))

	authHandler := authApp.Handler(func(w http.ResponseWriter, r *http.Request, user *authProvider.User) {
		if strings.HasPrefix(r.URL.Path, tagsPath) {
			tagsHandler(w, r, user, strings.TrimPrefix(r.URL.Path, tagsPath))
		} else if r.Method == http.MethodGet && (r.URL.Path == `/` || r.URL.Path == ``) {
			listReadings(w, r, user)
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
			if db.Ping(readingsDB) {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusServiceUnavailable)
			}
			return
		}

		authHandler.ServeHTTP(w, r)
	})
}
