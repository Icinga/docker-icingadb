FROM golang:alpine as entrypoint
RUN ["sh", "-exo", "pipefail", "-c", "apk add upx; rm -vf /var/cache/apk/*"]
COPY entrypoint /entrypoint

WORKDIR /entrypoint
ENV CGO_ENABLED 0

RUN ["go", "build", "-ldflags", "-s -w", "."]
RUN ["upx", "entrypoint"]


FROM alpine as empty
RUN ["mkdir", "/empty"]


FROM scratch
COPY rootfs/ /
COPY --from=empty --chown=icingadb:icingadb /empty /etc/icingadb
COPY --from=entrypoint /entrypoint/entrypoint /entrypoint
