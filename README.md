# eponae-api

[![Build Status](https://travis-ci.com/ViBiOh/eponae-api.svg?branch=master)](https://travis-ci.com/ViBiOh/eponae-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/eponae-api)](https://goreportcard.com/report/github.com/ViBiOh/eponae-api)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=ViBiOh/eponae-api)](https://dependabot.com)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ViBiOh_eponae-api&metric=alert_status)](https://sonarcloud.io/dashboard?id=ViBiOh_eponae-api)

## Usage

```bash
Usage of eponae-api:
  -address string
        [http] Listen address {EPONAE_API_ADDRESS}
  -cert string
        [http] Certificate file {EPONAE_API_CERT}
  -corsCredentials
        [cors] Access-Control-Allow-Credentials {EPONAE_API_CORS_CREDENTIALS}
  -corsExpose string
        [cors] Access-Control-Expose-Headers {EPONAE_API_CORS_EXPOSE}
  -corsHeaders string
        [cors] Access-Control-Allow-Headers {EPONAE_API_CORS_HEADERS} (default "Content-Type")
  -corsMethods string
        [cors] Access-Control-Allow-Methods {EPONAE_API_CORS_METHODS} (default "GET")
  -corsOrigin string
        [cors] Access-Control-Allow-Origin {EPONAE_API_CORS_ORIGIN} (default "*")
  -csp string
        [owasp] Content-Security-Policy {EPONAE_API_CSP} (default "default-src 'self'; base-uri 'self'")
  -dbHost string
        [db] Host {EPONAE_API_DB_HOST}
  -dbName string
        [db] Name {EPONAE_API_DB_NAME}
  -dbPass string
        [db] Pass {EPONAE_API_DB_PASS}
  -dbPort uint
        [db] Port {EPONAE_API_DB_PORT} (default 5432)
  -dbSslmode string
        [db] SSL Mode {EPONAE_API_DB_SSLMODE} (default "disable")
  -dbUser string
        [db] User {EPONAE_API_DB_USER}
  -frameOptions string
        [owasp] X-Frame-Options {EPONAE_API_FRAME_OPTIONS} (default "deny")
  -hsts
        [owasp] Indicate Strict Transport Security {EPONAE_API_HSTS} (default true)
  -key string
        [http] Key file {EPONAE_API_KEY}
  -okStatus int
        [http] Healthy HTTP Status code {EPONAE_API_OK_STATUS} (default 204)
  -port uint
        [http] Listen port {EPONAE_API_PORT} (default 1080)
  -prometheusPath string
        [prometheus] Path for exposing metrics {EPONAE_API_PROMETHEUS_PATH} (default "/metrics")
  -readingsDefaultPage uint
        [readings] Default page {EPONAE_API_READINGS_DEFAULT_PAGE} (default 1)
  -readingsDefaultPageSize uint
        [readings] Default page size {EPONAE_API_READINGS_DEFAULT_PAGE_SIZE} (default 20)
  -readingsMaxPageSize uint
        [readings] Max page size {EPONAE_API_READINGS_MAX_PAGE_SIZE} (default 100)
  -tagsDefaultPage uint
        [tags] Default page {EPONAE_API_TAGS_DEFAULT_PAGE} (default 1)
  -tagsDefaultPageSize uint
        [tags] Default page size {EPONAE_API_TAGS_DEFAULT_PAGE_SIZE} (default 20)
  -tagsMaxPageSize uint
        [tags] Max page size {EPONAE_API_TAGS_MAX_PAGE_SIZE} (default 100)
  -url string
        [alcotest] URL to check {EPONAE_API_URL}
  -userAgent string
        [alcotest] User-Agent for check {EPONAE_API_USER_AGENT} (default "Alcotest")
```
