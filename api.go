package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/ViBiOh/alcotest/alcotest"
	"github.com/ViBiOh/eponae-api/conservatories"
	"github.com/ViBiOh/eponae-api/healthcheck"
	"github.com/ViBiOh/eponae-api/readings"
	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/cert"
	"github.com/ViBiOh/httputils/cors"
	"github.com/ViBiOh/httputils/owasp"
	"github.com/ViBiOh/httputils/prometheus"
	"github.com/ViBiOh/httputils/rate"
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
	prometheusConfig := prometheus.Flags(`prometheus`)
	rateConfig := rate.Flags(`rate`)
	owaspConfig := owasp.Flags(``)
	corsConfig := cors.Flags(`cors`)

	flag.Parse()

	alcotest.DoAndExit(alcotestConfig)

	log.Printf(`Starting server on port %d`, *port)

	if err := conservatories.Init(); err != nil {
		log.Printf(`[conservatories] Error while initializing: %v`, err)
	}
	if err := readings.Init(); err != nil {
		log.Printf(`[readings] Error while initializing: %v`, err)
	}
	if err := healthcheck.Init(map[string]http.Handler{`/conservatories`: conservatoriesHandler, `/readings`: readingsHandler}); err != nil {
		log.Printf(`[healthcheck] Error while initializing: %v`, err)
	}

	conservatoriesHandler = http.StripPrefix(conservatoriesPath, conservatories.Handler())
	readingsHandler = http.StripPrefix(readingsPath, readings.Handler())
	healthcheckHandler = http.StripPrefix(healthcheckPath, healthcheck.Handler())

	server := &http.Server{
		Addr:    fmt.Sprintf(`:%d`, *port),
		Handler: prometheus.Handler(prometheusConfig, rate.Handler(rateConfig, gziphandler.GzipHandler(owasp.Handler(owaspConfig, cors.Handler(corsConfig, handler()))))),
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
