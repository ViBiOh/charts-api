# eponae-api

[![Build Status](https://travis-ci.org/ViBiOh/eponae-api.svg?branch=master)](https://travis-ci.org/ViBiOh/eponae-api)
[![codecov](https://codecov.io/gh/ViBiOh/eponae-api/branch/master/graph/badge.svg)](https://codecov.io/gh/ViBiOh/eponae-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/eponae-api)](https://goreportcard.com/report/github.com/ViBiOh/eponae-api)

## Usage

```bash
Usage of api:
  -authDisable
      [auth] Disable auth
  -authUrl string
      [auth] Auth URL, if remote
  -authUsers string
      [auth] Allowed users and profiles (e.g. user:profile1|profile2,user2:profile3). Empty allow any identified user
  -basicUsers string
      [Basic] Users in the form "id:username:password,id2:username2:password2"
  -corsCredentials
      [cors] Access-Control-Allow-Credentials
  -corsExpose string
      [cors] Access-Control-Expose-Headers
  -corsHeaders string
      [cors] Access-Control-Allow-Headers (default "Content-Type")
  -corsMethods string
      [cors] Access-Control-Allow-Methods (default "GET")
  -corsOrigin string
      [cors] Access-Control-Allow-Origin (default "*")
  -csp string
      [owasp] Content-Security-Policy (default "default-src 'self'; base-uri 'self'")
  -dbHost string
      [database] Host
  -dbName string
      [database] Name
  -dbPass string
      [database] Pass
  -dbPort string
      [database] Port (default "5432")
  -dbUser string
      [database] User
  -frameOptions string
      [owasp] X-Frame-Options (default "deny")
  -hsts
      [owasp] Indicate Strict Transport Security (default true)
  -port int
      Listen port (default 1080)
  -prometheusPath string
      [prometheus] Path for exposing metrics (default "/metrics")
  -readingsDefaultPage uint
      [readings] Default page (default 1)
  -readingsDefaultPageSize uint
      [readings] Default page size (default 20)
  -readingsMaxPageSize uint
      [readings] Max page size (default 500)
  -readingsPath string
      [readings] HTTP Path prefix (default "/readings")
  -tagsDefaultPage uint
      [tags] Default page (default 1)
  -tagsDefaultPageSize uint
      [tags] Default page size (default 20)
  -tagsMaxPageSize uint
      [tags] Max page size (default 500)
  -tagsPath string
      [tags] HTTP Path prefix (default "/tags")
  -tls
      Serve TLS content (default true)
  -tlsCert string
      [tls] PEM Certificate file
  -tlsHosts string
      [tls] Self-signed certificate hosts, comma separated (default "localhost")
  -tlsKey string
      [tls] PEM Key file
  -tlsOrganization string
      [tls] Self-signed certificate organization (default "ViBiOh")
  -tracingAgent string
      [opentracing] Jaeger Agent (e.g. host:port) (default "jaeger:6831")
  -tracingName string
      [opentracing] Service name
  -url string
      [health] URL to check
  -userAgent string
      [health] User-Agent for check (default "Golang alcotest")
```

## Postgres installation

```bash
export EPONAE_DATABASE_DIR=`realpath ./data`
export EPONAE_DATABASE_PASS=password
export EPONAE_DATABASE_PORT=5432

mkdir ${EPONAE_DATABASE_DIR}
sudo chown -R 70:70 ${EPONAE_DATABASE_DIR}

docker-compose -p eponae -f docker-compose.db.yml up -d
```

### Postgres configuration

```sql
\c eponae
```
