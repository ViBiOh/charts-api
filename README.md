# charts-api

[![Build Status](https://travis-ci.org/ViBiOh/charts-api.svg?branch=master)](https://travis-ci.org/ViBiOh/charts-api)
[![codecov](https://codecov.io/gh/ViBiOh/charts-api/branch/master/graph/badge.svg)](https://codecov.io/gh/ViBiOh/charts-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/charts-api)](https://goreportcard.com/report/github.com/ViBiOh/charts-api)

## Usage

```
Usage of api:
  -c string
    	URL to check
  -corsCredentials
    	Access-Control-Allow-Credentials
  -corsExpose string
    	Access-Control-Expose-Headers
  -corsHeaders string
    	Access-Control-Allow-Headers (default "Content-Type")
  -corsMethods string
    	Access-Control-Allow-Methods (default "GET")
  -corsOrigin string
    	Access-Control-Allow-Origin (default "*")
  -csp string
    	Content-Security-Policy (default "default-src 'self'")
  -dbHost string
    	Database Host
  -dbName string
    	Database Name
  -dbPass string
    	Database Pass
  -dbPort string
    	Database Port (default "5432")
  -dbUser string
    	Database User
  -hsts
    	Indicate Strict Transport Security (default true)
  -port string
    	Listen port (default "1080")
  -prometheusMetricsHost string
    	Prometheus - Allowed hostname to call metrics endpoint (default "localhost")
  -prometheusMetricsPath string
    	Prometheus - Metrics endpoint path (default "/metrics")
  -rateCount int
    	Rate IP limit (default 5000)
  -tls
    	Serve TLS content
  -tlscert string
    	TLS PEM Certificate file
  -tlshosts string
    	TLS Self-signed certificate hosts, comma separated (default "localhost")
  -tlskey string
    	TLS PEM Key file
```
