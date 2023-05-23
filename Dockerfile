# Icinga DB Docker image | (c) 2020 Icinga GmbH | GPLv2+

FROM golang:alpine as builder
RUN ["sh", "-exo", "pipefail", "-c", "apk add git; rm -vf /var/cache/apk/*"]
ENV CGO_ENABLED 0


FROM builder as download-once
COPY --from=icingadb-git . /icingadb-src/.git
WORKDIR /icingadb-src
RUN ["git", "checkout", "."]
RUN --mount=type=cache,target=/go/pkg/mod ["go", "mod", "download"]


FROM --platform=$BUILDPLATFORM download-once as download


FROM builder as icingadb

COPY --from=download /icingadb-src /icingadb-src
COPY --from=download /go/pkg/mod /go/pkg/mod
WORKDIR /icingadb-src

RUN --mount=type=cache,target=/root/.cache/go-build \
    ["go", "build", "-ldflags", "-s -w", "./cmd/icingadb"]

RUN ["bzip2", "-k", "schema/mysql/schema.sql"]
RUN ["bzip2", "-k", "schema/pgsql/schema.sql"]


FROM builder as download-entrypoint-once
COPY entrypoint /entrypoint
WORKDIR /entrypoint
RUN --mount=type=cache,target=/go/pkg/mod ["go", "mod", "download"]


FROM --platform=$BUILDPLATFORM download-entrypoint-once as download-entrypoint


FROM builder as entrypoint

COPY --from=download-entrypoint /entrypoint /entrypoint
COPY --from=download-entrypoint /go/pkg/mod /go/pkg/mod
WORKDIR /entrypoint

RUN --mount=type=cache,target=/root/.cache/go-build \
    ["go", "build", "-ldflags", "-s -w", "."]


FROM golang:alpine as base
RUN ["mkdir", "/empty"]
COPY rootfs /rootfs
RUN ["chmod", "-R", "u=rwX,go=rX", "/rootfs"]


FROM scratch

COPY --from=base /rootfs/ /
COPY --from=base --chown=icingadb:icingadb /empty /etc/icingadb
COPY --from=entrypoint /entrypoint/entrypoint /entrypoint
COPY --from=icingadb /icingadb-src/icingadb /
COPY --from=icingadb /icingadb-src/schema/mysql/schema.sql.bz2 /mysql.schema.sql.bz2
COPY --from=icingadb /icingadb-src/schema/pgsql/schema.sql.bz2 /pgsql.schema.sql.bz2

USER icingadb
CMD ["/entrypoint"]
