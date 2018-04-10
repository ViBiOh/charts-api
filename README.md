# eponae-api

[![Build Status](https://travis-ci.org/ViBiOh/eponae-api.svg?branch=master)](https://travis-ci.org/ViBiOh/eponae-api)
[![codecov](https://codecov.io/gh/ViBiOh/eponae-api/branch/master/graph/badge.svg)](https://codecov.io/gh/ViBiOh/eponae-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/eponae-api)](https://goreportcard.com/report/github.com/ViBiOh/eponae-api)

## Usage

```
Usage of api:
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
  -eponaeDbHost string
      [database] Host
  -eponaeDbName string
      [database] Name
  -eponaeDbPass string
      [database] Pass
  -eponaeDbPort string
      [database] Port (default "5432")
  -eponaeDbUser string
      [database] User
  -frameOptions string
      [owasp] X-Frame-Options (default "deny")
  -hsts
      [owasp] Indicate Strict Transport Security (default true)
  -port string
      Listen port (default "1080")
  -readingsAuthUrl string
      [auth] Auth URL, if remote
  -readingsAuthUsers string
      [auth] List of allowed users and profiles (e.g. user:profile1|profile2,user2:profile3)
  -readingsBasicUsers string
      [Basic] Users in the form "id:username:password,id2:username2:password2"
  -tls
      Serve TLS content
  -tlsCert string
      [tls] PEM Certificate file
  -tlsHosts string
      [tls] Self-signed certificate hosts, comma separated (default "localhost")
  -tlsKey string
      [tls] PEM Key file
  -url string
      [health] URL to check
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
