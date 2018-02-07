package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/ViBiOh/alcotest/alcotest"
	"github.com/ViBiOh/auth/auth"
	"github.com/ViBiOh/auth/provider/basic"
	authService "github.com/ViBiOh/auth/service"
	"github.com/ViBiOh/eponae-api/conservatories"
	"github.com/ViBiOh/eponae-api/healthcheck"
	"github.com/ViBiOh/eponae-api/readings"
	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/cert"
	"github.com/ViBiOh/httputils/cors"
	"github.com/ViBiOh/httputils/db"
	"github.com/ViBiOh/httputils/owasp"
)

const healthcheckPath = `/health`
const conservatoriesPath = `/conservatories`
const readingsPath = `/readings`

var (
	conservatoriesHandler http.Handler
	readingsHandler       http.Handler
	healthcheckHandler    http.Handler
)

func handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
}

func main() {
	port := flag.Int(`port`, 1080, `Listen port`)
	tls := flag.Bool(`tls`, true, `Serve TLS content`)
	alcotestConfig := alcotest.Flags(``)
	tlsConfig := cert.Flags(`tls`)
	owaspConfig := owasp.Flags(``)
	corsConfig := cors.Flags(`cors`)

	conservatoriesDbConfig := db.Flags(`conservatoriesDb`)

	readingsDbConfig := db.Flags(`readingsDb`)
	readingsAuthConfig := auth.Flags(`readingsAuth`)
	readingsAuthBasicConfig := basic.Flags(`readingsBasic`)

	flag.Parse()

	alcotest.DoAndExit(alcotestConfig)

	log.Printf(`Starting server on port %d`, *port)

	conservatoriesDB, err := db.GetDB(conservatoriesDbConfig)
	if err != nil {
		err = fmt.Errorf(`Error while initializing conservatories database: %v`, err)
	}
	conservatoriesApp := conservatories.NewApp(conservatoriesDB)
	conservatoriesHandler = http.StripPrefix(conservatoriesPath, conservatoriesApp.Handler())

	readingsDB, err := db.GetDB(readingsDbConfig)
	if err != nil {
		err = fmt.Errorf(`Error while initializing readings database: %v`, err)
	}
	readingsAuthApp := auth.NewApp(readingsAuthConfig, authService.NewBasicApp(readingsAuthBasicConfig))
	readingsApp := readings.NewApp(readingsDB, readingsAuthApp)
	readingsHandler = http.StripPrefix(readingsPath, readingsApp.Handler())

	healthcheckApp := healthcheck.NewApp(map[string]http.Handler{
		conservatoriesPath: conservatoriesHandler,
		readingsPath:       readingsHandler,
	})
	healthcheckHandler = http.StripPrefix(healthcheckPath, healthcheckApp.Handler())

	server := &http.Server{
		Addr:    fmt.Sprintf(`:%d`, *port),
		Handler: gziphandler.GzipHandler(owasp.Handler(owaspConfig, cors.Handler(corsConfig, handler()))),
	}

	var serveError = make(chan error)
	go func() {
		defer close(serveError)
		if *tls {
			log.Print(`Listening with TLS enabled`)
			serveError <- cert.ListenAndServeTLS(tlsConfig, server)
		} else {
			log.Print(`⚠ api is running without secure connection ⚠`)
			serveError <- server.ListenAndServe()
		}
	}()

	httputils.ServerGracefulClose(server, serveError, nil)
}
