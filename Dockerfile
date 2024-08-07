# Icinga DB Docker image | (c) 2020 Icinga GmbH | GPLv2+

FROM golang:alpine as icingadb
RUN ["sh", "-exo", "pipefail", "-c", "apk add git; rm -vf /var/cache/apk/*"]
ENV CGO_ENABLED 0

COPY --from=icingadb-git . /icingadb-src/.git
WORKDIR /icingadb-src
RUN ["git", "checkout", "."]

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    ["go", "build", "-ldflags", "-s -w", "./cmd/icingadb"]

RUN ["bzip2", "-k", "schema/mysql/schema.sql"]
RUN ["bzip2", "-k", "schema/pgsql/schema.sql"]


FROM golang:alpine as entrypoint
ENV CGO_ENABLED 0

COPY entrypoint /entrypoint
WORKDIR /entrypoint

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    ["go", "build", "-ldflags", "-s -w", "."]


FROM alpine as base
RUN ["mkdir", "/empty"]
COPY rootfs /rootfs
RUN ["chmod", "-R", "u=rwX,go=rX", "/rootfs"]


FROM scratch

LABEL org.opencontainers.image.documentation='https://icinga.com/docs/icinga-db/latest/doc/01-About/' \
      org.opencontainers.image.source='https://github.com/Icinga/icingadb' \
      org.opencontainers.image.licenses='GPL-2.0+'

COPY --from=base /rootfs/ /
COPY --from=base --chown=icingadb:icingadb /empty /etc/icingadb
COPY --from=entrypoint /entrypoint/entrypoint /entrypoint
COPY --from=icingadb /icingadb-src/icingadb /
COPY --from=icingadb /icingadb-src/schema/mysql/schema.sql.bz2 /mysql.schema.sql.bz2
COPY --from=icingadb /icingadb-src/schema/pgsql/schema.sql.bz2 /pgsql.schema.sql.bz2

USER icingadb
CMD ["/entrypoint"]
