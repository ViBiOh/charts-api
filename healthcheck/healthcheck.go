package healthcheck

import (
	"database/sql"
	"net/http"

	"github.com/ViBiOh/httputils/db"
)

var chartsDB *sql.DB

// Init charts handler
func Init(db *sql.DB) error {
	chartsDB = db

	return nil
}

// Handler for Health request. Should be use with net/http
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if db.Ping(chartsDB) {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusServiceUnavailable)
			}
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
