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

if ! docker version; then
	echo 'Docker not found' >&2
	false
fi

if ! docker buildx version; then
	echo '"docker buildx" not found (see https://docs.docker.com/buildx/working-with-buildx/ )' >&2
	false
fi

docker buildx build --load -t icinga/icingadb --build-context "icingadb-git=$(realpath "$IDBSRC")/.git" "$(realpath "$(dirname "$0")")"
