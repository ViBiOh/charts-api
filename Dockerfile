FROM scratch

HEALTHCHECK --retries=10 CMD [ "/api", "-url", "https://localhost:1080/health" ]

ENTRYPOINT [ "/api" ]
EXPOSE 1080

COPY cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY doc/api.html /doc/api.html
COPY bin/api /api
