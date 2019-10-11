# eponae-api

[![Build Status](https://travis-ci.org/ViBiOh/eponae-api.svg?branch=master)](https://travis-ci.org/ViBiOh/eponae-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/eponae-api)](https://goreportcard.com/report/github.com/ViBiOh/eponae-api)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=ViBiOh/eponae-api)](https://dependabot.com)

## Usage

```bash
Usage of eponae-api:
  -address string
        [http] Listen address {EPONAE_API_ADDRESS}
  -authDisable
        [auth] Disable auth {EPONAE_API_AUTH_DISABLE}
  -authUrl string
        [auth] Auth URL, if remote {EPONAE_API_AUTH_URL}
  -authUsers string
        [auth] Allowed users and profiles (e.g. user:profile1|profile2,user2:profile3). Empty allow any identified user {EPONAE_API_AUTH_USERS}
  -basicUsers id:username:password,id2:username2:password2
        [basic] Users in the form id:username:password,id2:username2:password2 {EPONAE_API_BASIC_USERS}
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
  -dbPort string
        [db] Port {EPONAE_API_DB_PORT} (default "5432")
  -dbUser string
        [db] User {EPONAE_API_DB_USER}
  -frameOptions string
        [owasp] X-Frame-Options {EPONAE_API_FRAME_OPTIONS} (default "deny")
  -hsts
        [owasp] Indicate Strict Transport Security {EPONAE_API_HSTS} (default true)
  -key string
        [http] Key file {EPONAE_API_KEY}
  -port int
        [http] Listen port {EPONAE_API_PORT} (default 1080)
  -prometheusPath string
        [prometheus] Path for exposing metrics {EPONAE_API_PROMETHEUS_PATH} (default "/metrics")
  -readingsDefaultPage uint
        [readings] Default page {EPONAE_API_READINGS_DEFAULT_PAGE} (default 1)
  -readingsDefaultPageSize uint
        [readings] Default page size {EPONAE_API_READINGS_DEFAULT_PAGE_SIZE} (default 20)
  -readingsMaxPageSize uint
        [readings] Max page size {EPONAE_API_READINGS_MAX_PAGE_SIZE} (default 500)
  -tagsDefaultPage uint
        [tags] Default page {EPONAE_API_TAGS_DEFAULT_PAGE} (default 1)
  -tagsDefaultPageSize uint
        [tags] Default page size {EPONAE_API_TAGS_DEFAULT_PAGE_SIZE} (default 20)
  -tagsMaxPageSize uint
        [tags] Max page size {EPONAE_API_TAGS_MAX_PAGE_SIZE} (default 500)
  -tracingAgent string
        [tracing] Jaeger Agent (e.g. host:port) {EPONAE_API_TRACING_AGENT} (default "jaeger:6831")
  -tracingName string
        [tracing] Service name {EPONAE_API_TRACING_NAME}
  -url string
        [alcotest] URL to check {EPONAE_API_URL}
  -userAgent string
        [alcotest] User-Agent for check {EPONAE_API_USER_AGENT} (default "Golang alcotest")
```
