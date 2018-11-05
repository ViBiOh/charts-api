package main

import (
	"flag"
	"net/http"
	"strings"

	"github.com/ViBiOh/auth/pkg/auth"
	"github.com/ViBiOh/auth/pkg/ident/basic"
	identService "github.com/ViBiOh/auth/pkg/ident/service"
	"github.com/ViBiOh/eponae-api/pkg/reading"
	"github.com/ViBiOh/eponae-api/pkg/readingtag"
	"github.com/ViBiOh/eponae-api/pkg/tag"
	"github.com/ViBiOh/httputils/pkg"
	"github.com/ViBiOh/httputils/pkg/alcotest"
	"github.com/ViBiOh/httputils/pkg/cors"
	"github.com/ViBiOh/httputils/pkg/crud"
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
	readingsPath = `/readings`
	tagsPath     = `/tags`
)

func main() {
	serverConfig := httputils.Flags(``)
	alcotestConfig := alcotest.Flags(``)
	prometheusConfig := prometheus.Flags(`prometheus`)
	opentracingConfig := opentracing.Flags(`tracing`)
	owaspConfig := owasp.Flags(``)
	corsConfig := cors.Flags(`cors`)

	dbConfig := db.Flags(`db`)
	authConfig := auth.Flags(`auth`)
	basicConfig := basic.Flags(`basic`)

	readingsConfig := crud.Flags(`readings`)
	tagsConfig := crud.Flags(`tags`)

	flag.Parse()

	alcotest.DoAndExit(alcotestConfig)

	serverApp := httputils.NewApp(serverConfig)
	healthcheckApp := healthcheck.NewApp()
	prometheusApp := prometheus.NewApp(prometheusConfig)
	opentracingApp := opentracing.NewApp(opentracingConfig)
	gzipApp := gzip.NewApp()
	owaspApp := owasp.NewApp(owaspConfig)
	corsApp := cors.NewApp(corsConfig)

	apiDB, err := db.GetDB(dbConfig)
	if err != nil {
		logger.Fatal(`%+v`, err)
	}
	authApp := auth.NewServiceApp(authConfig, identService.NewBasicApp(basicConfig, apiDB))

	tagService := tag.NewService(apiDB)
	readingTagService := readingtag.NewService(apiDB, tagService)
	readingService := reading.NewService(apiDB, readingTagService)

	readingsApp := crud.NewApp(readingsConfig, readingService)
	tagsApp := crud.NewApp(tagsConfig, tagService)

	readingsHandler := readingsApp.Handler()
	tagsHandler := tagsApp.Handler()

	apihandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, readingsPath) {
			readingsHandler.ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(r.URL.Path, tagsPath) {
			tagsHandler.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	})

	handler := server.ChainMiddlewares(apihandler, prometheusApp, opentracingApp, gzipApp, owaspApp, corsApp, authApp)
	healthcheckApp.NextHealthcheck(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if db.Ping(apiDB) {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}))

	serverApp.ListenAndServe(handler, nil, healthcheckApp)
}
