FROM scratch

HEALTHCHECK --retries=10 CMD https://localhost:1080/health

ENTRYPOINT [ "/bin/sh" ]
EXPOSE 1080

COPY bin/api /bin/sh
COPY doc/api.html /doc/api.html
COPY cacert.pem /etc/ssl/certs/ca-certificates.crt
