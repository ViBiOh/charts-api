FROM scratch

EXPOSE 1080
COPY ./doc /doc

HEALTHCHECK --retries=10 CMD [ "/eponae", "-url", "http://localhost:1080/health" ]
ENTRYPOINT [ "/eponae" ]

ARG VERSION
ENV VERSION=${VERSION}

ARG TARGETOS
ARG TARGETARCH

COPY cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY release/eponae_${TARGETOS}_${TARGETARCH} /eponae
