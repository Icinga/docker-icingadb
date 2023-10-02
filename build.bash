#!/bin/bash
# Icinga DB Docker image | (c) 2020 Icinga GmbH | GPLv2+
set -exo pipefail

IDBSRC="$1"
ACTION="${2:-local}"
TAG="${3:-test}"

if [ -z "$IDBSRC" ]; then
	cat <<EOF >&2
Usage: ${0} /icingadb/source/dir [local|all|push [TAG]]
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

OUR_DIR="$(realpath "$(dirname "$0")")"
COMMON_ARGS=(-t "icinga/icingadb:$TAG" --build-context "icingadb-git=$(realpath "$IDBSRC")/.git" "$OUR_DIR")
BUILDX=(docker buildx build --platform "$(cat "${OUR_DIR}/platforms.txt")")

case "$ACTION" in
	all)
		"${BUILDX[@]}" "${COMMON_ARGS[@]}"
		;;
	push)
		"${BUILDX[@]}" --push "${COMMON_ARGS[@]}"
		;;
	*)
		docker buildx build --load "${COMMON_ARGS[@]}"
		;;
esac
