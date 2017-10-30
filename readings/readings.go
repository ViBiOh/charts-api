package readings

import (
	"database/sql"
	"log"
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
		log.Printf(`[readings] Error while initializing database: %v`, err)
	} else if readingsDB != nil {
		log.Print(`[readings] Database ready`)
	}

	return nil
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

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(user.Username))
	})
}
