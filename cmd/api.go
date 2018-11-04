package main

import (
	"flag"
	"net/http"
	"strings"

	"github.com/ViBiOh/auth/pkg/auth"
	"github.com/ViBiOh/auth/pkg/provider/basic"
	authService "github.com/ViBiOh/auth/pkg/service"
	apiHealthcheck "github.com/ViBiOh/eponae-api/pkg/healthcheck"
	"github.com/ViBiOh/eponae-api/pkg/readings"
	"github.com/ViBiOh/httputils/pkg"
	"github.com/ViBiOh/httputils/pkg/alcotest"
	"github.com/ViBiOh/httputils/pkg/cors"
	"github.com/ViBiOh/httputils/pkg/db"
	"github.com/ViBiOh/httputils/pkg/gzip"
	"github.com/ViBiOh/httputils/pkg/healthcheck"
	"github.com/ViBiOh/httputils/pkg/logger"
	"github.com/ViBiOh/httputils/pkg/opentracing"
	"github.com/ViBiOh/httputils/pkg/owasp"
	"github.com/ViBiOh/httputils/pkg/prometheus"
	"github.com/ViBiOh/httputils/pkg/server"
)

const (
	healthcheckPath = `/health`
	readingsPath    = `/readings`
)

func main() {
	serverConfig := httputils.Flags(``)
	alcotestConfig := alcotest.Flags(``)
	prometheusConfig := prometheus.Flags(`prometheus`)
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
	prometheusApp := prometheus.NewApp(prometheusConfig)
	opentracingApp := opentracing.NewApp(opentracingConfig)
	gzipApp := gzip.NewApp()
	owaspApp := owasp.NewApp(owaspConfig)
	corsApp := cors.NewApp(corsConfig)

	eponaeDB, err := db.GetDB(eponaeDbConfig)
	if err != nil {
		logger.Fatal(`%+v`, err)
	}
	readingsAuthApp := auth.NewApp(readingsAuthConfig, authService.NewBasicApp(readingsAuthBasicConfig))
	readingsApp := readings.NewApp(eponaeDB)

	readingsHandler := server.ChainMiddlewares(http.StripPrefix(readingsPath, readingsApp.Handler()), readingsAuthApp)

	healthcheckApp.NextHealthcheck(http.StripPrefix(healthcheckPath, apiHealthcheck.NewApp(map[string]http.Handler{
		readingsPath: readingsHandler,
	}).Handler()))

	apihandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, readingsPath) {
			readingsHandler.ServeHTTP(w, r)
			return
		}

		if r.Method == http.MethodGet && (r.URL.Path == `/` || r.URL.Path == ``) {
			http.ServeFile(w, r, `doc/api.html`)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	})

	handler := server.ChainMiddlewares(apihandler, prometheusApp, opentracingApp, gzipApp, owaspApp, corsApp)

	serverApp.ListenAndServe(handler, nil, healthcheckApp)
}
