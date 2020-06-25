FROM golang:alpine as entrypoint
COPY entrypoint /entrypoint

WORKDIR /entrypoint
ENV CGO_ENABLED 0
RUN ["go", "build", "."]


FROM alpine as empty
RUN ["mkdir", "/empty"]


FROM scratch
COPY rootfs/ /
COPY --from=empty --chown=icingadb:icingadb /empty /etc/icingadb
COPY --from=entrypoint /entrypoint/entrypoint /entrypoint
