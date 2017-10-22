package charts

import (
	"database/sql"
	"net/http"
	"strings"
)

const conservatoriesPath = `/conservatories`

var chartsDB *sql.DB
var conservatoriesStrippedHandler = http.StripPrefix(conservatoriesPath, conservatoriesHandler())

// Init charts handler
func Init(db *sql.DB) error {
	chartsDB = db

	return nil
}

// Handler for CRUD request. Should be use with net/http
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, conservatoriesPath) {
			conservatoriesStrippedHandler.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
}
