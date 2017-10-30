package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/ViBiOh/alcotest/alcotest"
	"github.com/ViBiOh/eponae-api/charts"
	"github.com/ViBiOh/eponae-api/healthcheck"
	"github.com/ViBiOh/eponae-api/readings"
	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/cert"
	"github.com/ViBiOh/httputils/cors"
	"github.com/ViBiOh/httputils/db"
	"github.com/ViBiOh/httputils/owasp"
	"github.com/ViBiOh/httputils/prometheus"
	"github.com/ViBiOh/httputils/rate"
)

const healthcheckPath = `/health`
const conservatoriesPath = `/conservatories`
const readingsPath = `/readings`

var healthcheckHandler = http.StripPrefix(healthcheckPath, healthcheck.Handler())
var chartsHandler = http.StripPrefix(conservatoriesPath, charts.Handler())
var readingsHandler = http.StripPrefix(readingsPath, readings.Handler())
var restHandler http.Handler

func handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, healthcheckPath) {
			healthcheckHandler.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, conservatoriesPath) {
			chartsHandler.ServeHTTP(w, r)
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
	url := flag.String(`c`, ``, `URL to check`)
	port := flag.String(`port`, `1080`, `Listen port`)
	tls := flag.Bool(`tls`, false, `Serve TLS content`)
	prometheusConfig := prometheus.Flags(`prometheus`)
	rateConfig := rate.Flags(`rate`)
	owaspConfig := owasp.Flags(``)
	corsConfig := cors.Flags(`cors`)
	dbConfig := db.Flags(`db`)
	flag.Parse()

	if *url != `` {
		alcotest.Do(url)
		return
	}

	chartsDB, err := db.GetDB(dbConfig)
	if err != nil {
		log.Printf(`Error while initializing database: %v`, err)
	} else if chartsDB != nil {
		log.Print(`Database ready`)
	}

	log.Printf(`Starting server on port %s`, *port)

	if err := healthcheck.Init(chartsDB); err != nil {
		log.Printf(`Error while initializing healthcheck: %v`, err)
	}
	if err := charts.Init(chartsDB); err != nil {
		log.Printf(`Error while initializing charts: %v`, err)
	}

	restHandler = prometheus.Handler(prometheusConfig, rate.Handler(rateConfig, gziphandler.GzipHandler(owasp.Handler(owaspConfig, cors.Handler(corsConfig, handler())))))
	server := &http.Server{
		Addr:    `:` + *port,
		Handler: restHandler,
	}

	var serveError = make(chan error)
	go func() {
		defer close(serveError)
		if *tls {
			log.Print(`Listening with TLS enabled`)
			serveError <- cert.ListenAndServeTLS(server)
		} else {
			serveError <- server.ListenAndServe()
		}
	}()

	httputils.ServerGracefulClose(server, serveError, nil)
}
