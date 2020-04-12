package main

import (
	"flag"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/ViBiOh/auth/v2/pkg/handler"
	"github.com/ViBiOh/auth/v2/pkg/ident/basic"
	basicDb "github.com/ViBiOh/auth/v2/pkg/provider/db"
	"github.com/ViBiOh/eponae-api/pkg/reading"
	"github.com/ViBiOh/eponae-api/pkg/readingtag"
	"github.com/ViBiOh/eponae-api/pkg/tag"
	"github.com/ViBiOh/httputils/v3/pkg/alcotest"
	"github.com/ViBiOh/httputils/v3/pkg/cors"
	"github.com/ViBiOh/httputils/v3/pkg/crud"
	"github.com/ViBiOh/httputils/v3/pkg/db"
	"github.com/ViBiOh/httputils/v3/pkg/httputils"
	"github.com/ViBiOh/httputils/v3/pkg/logger"
	"github.com/ViBiOh/httputils/v3/pkg/owasp"
	"github.com/ViBiOh/httputils/v3/pkg/prometheus"
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
	owaspConfig := owasp.Flags(fs, "")
	corsConfig := cors.Flags(fs, "cors")

	dbConfig := db.Flags(fs, "db")

	readingsConfig := crud.GetConfiguredFlags(readingsPath, "Readings")(fs, "readings")
	tagsConfig := crud.GetConfiguredFlags(tagsPath, "Tags")(fs, "tags")

	logger.Fatal(fs.Parse(os.Args[1:]))

	alcotest.DoAndExit(alcotestConfig)

	apiDB, err := db.New(dbConfig)
	logger.Fatal(err)

	basicApp := basicDb.New(apiDB)
	basicProvider := basic.New(basicApp)
	authApp := basicDb.New(apiDB)

	tagService := tag.New(apiDB)
	readingTagService := readingtag.New(apiDB, tagService)
	readingService := reading.New(apiDB, readingTagService, tagService)

	readingsApp, err := crud.New(readingsConfig, readingService)
	logger.Fatal(err)
	tagsApp, err := crud.New(tagsConfig, tagService)
	logger.Fatal(err)

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

	server := httputils.New(serverConfig)
	server.Health(apiDB.Ping)
	server.Middleware(prometheus.New(prometheusConfig).Middleware)
	server.Middleware(owasp.New(owaspConfig).Middleware)
	server.Middleware(cors.New(corsConfig).Middleware)
	server.Middleware(handler.New(authApp, basicProvider).Middleware)
	server.ListenServeWait(apihandler)
}
