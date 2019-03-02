FROM golang:1.12 as builder

ENV APP_NAME eponae-api
ENV WORKDIR ${GOPATH}/src/github.com/ViBiOh/eponae-api

WORKDIR ${WORKDIR}
COPY ./ ${WORKDIR}/

RUN make ${APP_NAME} \
 && mkdir -p /app \
 && curl -s -o /app/cacert.pem https://curl.haxx.se/ca/cacert.pem \
 && cp bin/${APP_NAME} /app/

FROM scratch

ENV APP_NAME eponae-api
EXPOSE 1080

HEALTHCHECK --retries=10 CMD [ "/api", "-url", "https://localhost:1080/health" ]
ENTRYPOINT [ "/api" ]

COPY --from=builder /app/cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/${APP_NAME} /api
