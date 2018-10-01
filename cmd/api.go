package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/ViBiOh/auth/pkg/auth"
	"github.com/ViBiOh/auth/pkg/provider/basic"
	authService "github.com/ViBiOh/auth/pkg/service"
	"github.com/ViBiOh/eponae-api/pkg/conservatories"
	apiHealthcheck "github.com/ViBiOh/eponae-api/pkg/healthcheck"
	"github.com/ViBiOh/eponae-api/pkg/readings"
	"github.com/ViBiOh/httputils/pkg"
	"github.com/ViBiOh/httputils/pkg/alcotest"
	"github.com/ViBiOh/httputils/pkg/cors"
	"github.com/ViBiOh/httputils/pkg/db"
	"github.com/ViBiOh/httputils/pkg/gzip"
	"github.com/ViBiOh/httputils/pkg/healthcheck"
	"github.com/ViBiOh/httputils/pkg/opentracing"
	"github.com/ViBiOh/httputils/pkg/owasp"
	"github.com/ViBiOh/httputils/pkg/server"
)

const (
	healthcheckPath    = `/health`
	conservatoriesPath = `/conservatories`
	readingsPath       = `/readings`
)

func main() {
	serverConfig := httputils.Flags(``)
	alcotestConfig := alcotest.Flags(``)
	opentracingConfig := opentracing.Flags(`tracing`)
	owaspConfig := owasp.Flags(``)
	corsConfig := cors.Flags(`cors`)

	eponaeDbConfig := db.Flags(`eponaeDb`)
	readingsAuthConfig := auth.Flags(`readingsAuth`)
	readingsAuthBasicConfig := basic.Flags(`readingsBasic`)

	flag.Parse()

	alcotest.DoAndExit(alcotestConfig)

	serverApp := httputils.NewApp(serverConfig)
	healthcheckApp := healthcheck.NewApp()
	opentracingApp := opentracing.NewApp(opentracingConfig)
	owaspApp := owasp.NewApp(owaspConfig)
	corsApp := cors.NewApp(corsConfig)
	gzipApp := gzip.NewApp()

	eponaeDB, err := db.GetDB(eponaeDbConfig)
	if err != nil {
		err = fmt.Errorf(`error while initializing database: %v`, err)
	}
	conservatoriesApp := conservatories.NewApp(eponaeDB)
	readingsAuthApp := auth.NewApp(readingsAuthConfig, authService.NewBasicApp(readingsAuthBasicConfig))
	readingsApp := readings.NewApp(eponaeDB, readingsAuthApp)

	conservatoriesHandler := http.StripPrefix(conservatoriesPath, conservatoriesApp.Handler())
	readingsHandler := http.StripPrefix(readingsPath, readingsApp.Handler())

	healthcheckApp.NextHealthcheck(http.StripPrefix(healthcheckPath, apiHealthcheck.NewApp(map[string]http.Handler{
		conservatoriesPath: conservatoriesHandler,
		readingsPath:       readingsHandler,
	}).Handler()))

	apihandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, conservatoriesPath) {
			conservatoriesHandler.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, readingsPath) {
			readingsHandler.ServeHTTP(w, r)
		} else if r.Method == http.MethodGet && (r.URL.Path == `/` || r.URL.Path == ``) {
			http.ServeFile(w, r, `doc/api.html`)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	handler := server.ChainMiddlewares(apihandler, opentracingApp, gzipApp, owaspApp, corsApp)

	serverApp.ListenAndServe(handler, nil, healthcheckApp)
}
