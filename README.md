# eponae-api

[![Build Status](https://travis-ci.org/ViBiOh/eponae-api.svg?branch=master)](https://travis-ci.org/ViBiOh/eponae-api)
[![codecov](https://codecov.io/gh/ViBiOh/eponae-api/branch/master/graph/badge.svg)](https://codecov.io/gh/ViBiOh/eponae-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/eponae-api)](https://goreportcard.com/report/github.com/ViBiOh/eponae-api)

## Usage

```
Usage of api:
  -c string
        [health] URL to check
  -conservatoriesDbHost string
        [database] Host
  -conservatoriesDbName string
        [database] Name
  -conservatoriesDbPass string
        [database] Pass
  -conservatoriesDbPort string
        [database] Port (default "5432")
  -conservatoriesDbUser string
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
  -frameOptions string
        [owasp] X-Frame-Options (default "deny")
  -hsts
        [owasp] Indicate Strict Transport Security (default true)
  -port int
        Listen port (default 1080)
  -readingsAuthUrl string
        [auth] Auth URL, if remote
  -readingsAuthUsers string
        [auth] List of allowed users and profiles (e.g. user:profile1|profile2,user2:profile3)
  -readingsBasicUsers string
        [Basic] Users in the form "id:username:password,id2:username2:password2"
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
read -p "EPONAE_DATABASE_PASS=" EPONAE_DATABASE_PASS
read -p "EPONAE_DATABASE_PORT=" EPONAE_DATABASE_PORT

export EPONAE_DATABASE_DIR=`realpath ./data`
export EPONAE_DATABASE_PASS=${EPONAE_DATABASE_PASS}
export EPONAE_DATABASE_PORT=${EPONAE_DATABASE_PORT}

mkdir ${EPONAE_DATABASE_DIR}
sudo chown -R 70:70 ${EPONAE_DATABASE_DIR}

docker-compose -p eponae -f docker-compose.db.yml up -d
docker exec -it eponae_db_1 psql
```

### Postgres configuration

```sql
\c eponae
```
