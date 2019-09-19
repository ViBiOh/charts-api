package main

import (
	"flag"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/ViBiOh/auth/pkg/auth"
	"github.com/ViBiOh/auth/pkg/ident/basic"
	identService "github.com/ViBiOh/auth/pkg/ident/service"
	"github.com/ViBiOh/eponae-api/pkg/reading"
	"github.com/ViBiOh/eponae-api/pkg/readingtag"
	"github.com/ViBiOh/eponae-api/pkg/tag"
	httputils "github.com/ViBiOh/httputils/v2/pkg"
	"github.com/ViBiOh/httputils/v2/pkg/alcotest"
	"github.com/ViBiOh/httputils/v2/pkg/cors"
	"github.com/ViBiOh/httputils/v2/pkg/crud"
	"github.com/ViBiOh/httputils/v2/pkg/db"
	"github.com/ViBiOh/httputils/v2/pkg/logger"
	"github.com/ViBiOh/httputils/v2/pkg/opentracing"
	"github.com/ViBiOh/httputils/v2/pkg/owasp"
	"github.com/ViBiOh/httputils/v2/pkg/prometheus"
)

const (
	readingsPath = "/readings"
	tagsPath     = "/tags"

	docPath = "doc/"
)

func main() {
	fs := flag.NewFlagSet("eponae-api", flag.ExitOnError)

	serverConfig := httputils.Flags(fs, "")
	alcotestConfig := alcotest.Flags(fs, "")
	prometheusConfig := prometheus.Flags(fs, "prometheus")
	opentracingConfig := opentracing.Flags(fs, "tracing")
	owaspConfig := owasp.Flags(fs, "")
	corsConfig := cors.Flags(fs, "cors")

	dbConfig := db.Flags(fs, "db")
	authConfig := auth.Flags(fs, "auth")
	basicConfig := basic.Flags(fs, "basic")

	readingsConfig := crud.Flags(fs, "readings")
	tagsConfig := crud.Flags(fs, "tags")

	logger.Fatal(fs.Parse(os.Args[1:]))

	alcotest.DoAndExit(alcotestConfig)

	prometheusApp := prometheus.New(prometheusConfig)
	opentracingApp := opentracing.New(opentracingConfig)
	owaspApp := owasp.New(owaspConfig)
	corsApp := cors.New(corsConfig)

	apiDB, err := db.New(dbConfig)
	logger.Fatal(err)
	authApp := auth.NewService(authConfig, identService.NewBasic(basicConfig, apiDB))

	tagService := tag.New(apiDB)
	readingTagService := readingtag.New(apiDB, tagService)
	readingService := reading.New(apiDB, readingTagService, tagService)

	readingsApp := crud.New(readingsConfig, readingService)
	tagsApp := crud.New(tagsConfig, tagService)

	readingsHandler := http.StripPrefix(readingsPath, readingsApp.Handler())
	tagsHandler := http.StripPrefix(tagsPath, tagsApp.Handler())

	apihandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, readingsPath) {
			readingsHandler.ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(r.URL.Path, tagsPath) {
			tagsHandler.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(w, r, path.Join(docPath, r.URL.Path))
	})

	handler := httputils.ChainMiddlewares(apihandler, prometheusApp, opentracingApp, owaspApp, corsApp, authApp)

	httputils.New(serverConfig).ListenAndServe(handler, httputils.HealthHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if db.Ping(apiDB) {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})), nil)
}
