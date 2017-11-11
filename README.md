# eponae-api

[![Build Status](https://travis-ci.org/ViBiOh/eponae-api.svg?branch=master)](https://travis-ci.org/ViBiOh/eponae-api)
[![codecov](https://codecov.io/gh/ViBiOh/eponae-api/branch/master/graph/badge.svg)](https://codecov.io/gh/ViBiOh/eponae-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/eponae-api)](https://goreportcard.com/report/github.com/ViBiOh/eponae-api)

## Usage

```
Usage of api:
  -c string
    	[health] URL to check
  -chartsDbHost string
    	[database] Host
  -chartsDbName string
    	[database] Name
  -chartsDbPass string
    	[database] Pass
  -chartsDbPort string
    	[database] Port (default "5432")
  -chartsDbUser string
    	[database] User
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
    	[owasp] Content-Security-Policy (default "default-src 'self'")
  -hsts
    	[owasp] Indicate Strict Transport Security (default true)
  -port string
    	Listen port (default "1080")
  -prometheusMetricsHost string
    	[prometheus] Allowed hostname to call metrics endpoint (default "localhost")
  -prometheusMetricsPath string
    	[prometheus] Metrics endpoint path (default "/metrics")
  -prometheusPrefix string
    	[prometheus] Prefix (default "http")
  -rateCount uint
    	[rate] IP limit (default 5000)
  -readingsAuthUrl string
    	[auth] Auth URL
  -readingsAuthUsers string
    	[auth] List of allowed users and profiles (e.g. user:profile1|profile2,user2:profile3)
  -readingsDbHost string
    	[database] Host
  -readingsDbName string
    	[database] Name
  -readingsDbPass string
    	[database] Pass
  -readingsDbPort string
    	[database] Port (default "5432")
  -readingsDbUser string
    	[database] User
  -tls
    	Serve TLS content (default true)
  -tlsCert string
    	[tls] PEM Certificate file
  -tlsHosts string
    	[tls] Self-signed certificate hosts, comma separated (default "localhost")
  -tlsKey string
    	[tls] PEM Key file
```

## Postgres installation

```bash
export CHARTS_DATABASE_DIR=`realpath ./data_charts`
export READINGS_DATABASE_DIR=`realpath ./data_readings`
mkdir ${CHARTS_DATABASE_DIR}
mkdir ${READINGS_DATABASE_DIR}
sudo chown -R 70:70 ${CHARTS_DATABASE_DIR}
sudo chown -R 70:70 ${READINGS_DATABASE_DIR}
```
