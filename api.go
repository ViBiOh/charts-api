package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/ViBiOh/auth/auth"
	"github.com/ViBiOh/auth/provider/basic"
	authService "github.com/ViBiOh/auth/service"
	"github.com/ViBiOh/eponae-api/conservatories"
	"github.com/ViBiOh/eponae-api/healthcheck"
	"github.com/ViBiOh/eponae-api/readings"
	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/cors"
	"github.com/ViBiOh/httputils/db"
	"github.com/ViBiOh/httputils/owasp"
)

const (
	healthcheckPath    = `/health`
	conservatoriesPath = `/conservatories`
	readingsPath       = `/readings`
)

func main() {
	owaspConfig := owasp.Flags(``)
	corsConfig := cors.Flags(`cors`)
	eponaeDbConfig := db.Flags(`eponaeDb`)
	readingsAuthConfig := auth.Flags(`readingsAuth`)
	readingsAuthBasicConfig := basic.Flags(`readingsBasic`)

	httputils.NewApp(httputils.Flags(``), func() http.Handler {
		eponaeDB, err := db.GetDB(eponaeDbConfig)
		if err != nil {
			err = fmt.Errorf(`Error while initializing database: %v`, err)
		}

		conservatoriesApp := conservatories.NewApp(eponaeDB)
		conservatoriesHandler := http.StripPrefix(conservatoriesPath, conservatoriesApp.Handler())

		readingsAuthApp := auth.NewApp(readingsAuthConfig, authService.NewBasicApp(readingsAuthBasicConfig))
		readingsApp := readings.NewApp(eponaeDB, readingsAuthApp)
		readingsHandler := http.StripPrefix(readingsPath, readingsApp.Handler())

		healthcheckApp := healthcheck.NewApp(map[string]http.Handler{
			conservatoriesPath: conservatoriesHandler,
			readingsPath:       readingsHandler,
		})
		healthcheckHandler := http.StripPrefix(healthcheckPath, healthcheckApp.Handler())

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, healthcheckPath) {
				healthcheckHandler.ServeHTTP(w, r)
			} else if strings.HasPrefix(r.URL.Path, conservatoriesPath) {
				conservatoriesHandler.ServeHTTP(w, r)
			} else if strings.HasPrefix(r.URL.Path, readingsPath) {
				readingsHandler.ServeHTTP(w, r)
			} else if r.Method == http.MethodGet && (r.URL.Path == `/` || r.URL.Path == ``) {
				http.ServeFile(w, r, `doc/api.html`)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		})

		return gziphandler.GzipHandler(owasp.Handler(owaspConfig, cors.Handler(corsConfig, handler)))
	}, nil).ListenAndServe()
}
