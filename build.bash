#!/bin/bash
# Icinga DB Docker image | (c) 2020 Icinga GmbH | GPLv2+
set -exo pipefail

IDBSRC="$1"

if [ -z "$IDBSRC" ]; then
	cat <<EOF >&2
Usage: ${0} /icingadb/source/dir
EOF

	false
fi

IDBSRC="$(realpath "$IDBSRC")"
BLDCTX="$(realpath "$(dirname "$0")")"

docker build -f "${BLDCTX}/action.Dockerfile" -t icinga/icingadb-builder "$BLDCTX"

docker run --rm -i \
	-v "${IDBSRC}:/idbsrc:ro" \
	-v "${BLDCTX}:/bldctx:ro" \
	-v /var/run/docker.sock:/var/run/docker.sock \
	icinga/icingadb-builder bash <<EOF
set -exo pipefail

git -C /idbsrc archive --prefix=idbcp/ HEAD |tar -xC /
cd /idbcp

CGO_ENABLED=0 go build -ldflags '-s -w' ./cmd/icingadb
upx icingadb
bzip2 <schema/mysql/schema.sql >mysql.schema.sql.bz2

cp -r /entrypoint .
cp -r /rootfs .

docker build -f /bldctx/Dockerfile -t icinga/icingadb .
EOF
