# eponae-api

[![Build Status](https://travis-ci.org/ViBiOh/eponae-api.svg?branch=master)](https://travis-ci.org/ViBiOh/eponae-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/eponae-api)](https://goreportcard.com/report/github.com/ViBiOh/eponae-api)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=ViBiOh/eponae-api)](https://dependabot.com)

## Usage

```bash
Usage of eponae-api:
  -address string
        [http] Listen address {EPONAE-API_ADDRESS}
  -authDisable
        [auth] Disable auth {EPONAE-API_AUTH_DISABLE}
  -authUrl string
        [auth] Auth URL, if remote {EPONAE-API_AUTH_URL}
  -authUsers string
        [auth] Allowed users and profiles (e.g. user:profile1|profile2,user2:profile3). Empty allow any identified user {EPONAE-API_AUTH_USERS}
  -basicUsers id:username:password,id2:username2:password2
        [basic] Users in the form id:username:password,id2:username2:password2 {EPONAE-API_BASIC_USERS}
  -cert string
        [http] Certificate file {EPONAE-API_CERT}
  -corsCredentials
        [cors] Access-Control-Allow-Credentials {EPONAE-API_CORS_CREDENTIALS}
  -corsExpose string
        [cors] Access-Control-Expose-Headers {EPONAE-API_CORS_EXPOSE}
  -corsHeaders string
        [cors] Access-Control-Allow-Headers {EPONAE-API_CORS_HEADERS} (default "Content-Type")
  -corsMethods string
        [cors] Access-Control-Allow-Methods {EPONAE-API_CORS_METHODS} (default "GET")
  -corsOrigin string
        [cors] Access-Control-Allow-Origin {EPONAE-API_CORS_ORIGIN} (default "*")
  -csp string
        [owasp] Content-Security-Policy {EPONAE-API_CSP} (default "default-src 'self'; base-uri 'self'")
  -dbHost string
        [db] Host {EPONAE-API_DB_HOST}
  -dbName string
        [db] Name {EPONAE-API_DB_NAME}
  -dbPass string
        [db] Pass {EPONAE-API_DB_PASS}
  -dbPort string
        [db] Port {EPONAE-API_DB_PORT} (default "5432")
  -dbSslmode string
        [db] SSL Mode {EPONAE-API_DB_SSLMODE} (default "disable")
  -dbUser string
        [db] User {EPONAE-API_DB_USER}
  -frameOptions string
        [owasp] X-Frame-Options {EPONAE-API_FRAME_OPTIONS} (default "deny")
  -hsts
        [owasp] Indicate Strict Transport Security {EPONAE-API_HSTS} (default true)
  -key string
        [http] Key file {EPONAE-API_KEY}
  -port int
        [http] Listen port {EPONAE-API_PORT} (default 1080)
  -prometheusPath string
        [prometheus] Path for exposing metrics {EPONAE-API_PROMETHEUS_PATH} (default "/metrics")
  -readingsDefaultPage uint
        [readings] Default page {EPONAE-API_READINGS_DEFAULT_PAGE} (default 1)
  -readingsDefaultPageSize uint
        [readings] Default page size {EPONAE-API_READINGS_DEFAULT_PAGE_SIZE} (default 20)
  -readingsMaxPageSize uint
        [readings] Max page size {EPONAE-API_READINGS_MAX_PAGE_SIZE} (default 100)
  -tagsDefaultPage uint
        [tags] Default page {EPONAE-API_TAGS_DEFAULT_PAGE} (default 1)
  -tagsDefaultPageSize uint
        [tags] Default page size {EPONAE-API_TAGS_DEFAULT_PAGE_SIZE} (default 20)
  -tagsMaxPageSize uint
        [tags] Max page size {EPONAE-API_TAGS_MAX_PAGE_SIZE} (default 100)
  -url string
        [alcotest] URL to check {EPONAE-API_URL}
  -userAgent string
        [alcotest] User-Agent for check {EPONAE-API_USER_AGENT} (default "Golang alcotest")
```
