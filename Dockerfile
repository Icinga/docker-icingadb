# Icinga DB Docker image | (c) 2020 Icinga GmbH | GPLv2+

FROM golang:alpine as entrypoint
RUN ["sh", "-exo", "pipefail", "-c", "apk add upx; rm -vf /var/cache/apk/*"]
COPY entrypoint /entrypoint

WORKDIR /entrypoint
ENV CGO_ENABLED 0

RUN ["go", "build", "-ldflags", "-s -w", "."]
RUN ["upx", "entrypoint"]


FROM alpine as base
RUN ["mkdir", "/empty"]
COPY rootfs /rootfs
RUN ["chmod", "-R", "u=rwX,go=rX", "/rootfs"]


FROM scratch

COPY --from=base /rootfs/ /
COPY --from=base --chown=icingadb:icingadb /empty /etc/icingadb
COPY --from=entrypoint /entrypoint/entrypoint /entrypoint
COPY icingadb mysql.schema.sql.bz2 /

USER icingadb
CMD ["/entrypoint"]
