# eponae-api

[![Build Status](https://travis-ci.org/ViBiOh/eponae-api.svg?branch=master)](https://travis-ci.org/ViBiOh/eponae-api)
[![codecov](https://codecov.io/gh/ViBiOh/eponae-api/branch/master/graph/badge.svg)](https://codecov.io/gh/ViBiOh/eponae-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/eponae-api)](https://goreportcard.com/report/github.com/ViBiOh/eponae-api)

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
    	Prometheus allowed hostname to call metrics endpoint (default "localhost")
  -prometheusMetricsPath string
    	Prometheus metrics endpoint path (default "/metrics")
  -prometheusPrefix string
    	Prometheus prefix (default "http")
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
