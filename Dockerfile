FROM alpine as fetcher

WORKDIR /app

RUN apk --update add curl \
 && curl -q -sSL --max-time 10 -o /app/cacert.pem https://curl.haxx.se/ca/cacert.pem

FROM scratch

EXPOSE 1080

HEALTHCHECK --retries=10 CMD [ "/eponae", "-url", "http://localhost:1080/health" ]
ENTRYPOINT [ "/eponae" ]

ARG APP_VERSION
ENV VERSION=${APP_VERSION}

ARG OS
ARG ARCH

COPY ./doc /doc
COPY --from=fetcher /app/cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY release/eponae_${OS}_${ARCH} /eponae
