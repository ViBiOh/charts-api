package readings

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ViBiOh/auth/auth"
	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/db"
)

const healthcheckPath = `/health`

var authURL string
var authUsers map[string]*auth.User

var authConfig = auth.Flags(`readingsAuth`)
var dbConfig = db.Flags(`readingsDb`)
var readingsDB *sql.DB

// Init readings API
func Init() (err error) {
	authURL = *authConfig[`url`]
	authUsers = auth.LoadUsersProfiles(*authConfig[`users`])

	readingsDB, err = db.GetDB(dbConfig)
	if err != nil {
		err = fmt.Errorf(`Error while initializing database: %v`, err)
	}

	return
}

func listReadings(w http.ResponseWriter, r *http.Request, user *auth.User) {
	if list, err := listReadingsOfUser(user); err == nil {
		httputils.ResponseArrayJSON(w, http.StatusOK, list, httputils.IsPretty(r.URL.RawQuery))
	} else {
		httputils.InternalServerError(w, err)
	}
}

// Handler for Readings request. Should be use with net/http
func Handler() http.Handler {
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

		user, err := auth.IsAuthenticated(authURL, authUsers, r)
		if err != nil {
			httputils.Unauthorized(w, err)
			return
		}

		if r.Method == http.MethodGet && (r.URL.Path == `/` || r.URL.Path == ``) {
			listReadings(w, r, user)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
