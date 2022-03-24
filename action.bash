#!/bin/bash
# Icinga DB Docker image | (c) 2020 Icinga GmbH | GPLv2+
set -exo pipefail

TARGET=icinga/icingadb

mkimg () {
	test -n "$TAG"

	node /actions/checkout/dist/index.js |grep -vFe ::add-matcher::

	CGO_ENABLED=0 go build -ldflags '-s -w' ./cmd/icingadb
	upx icingadb
	bzip2 <schema/mysql/schema.sql >mysql.schema.sql.bz2
	bzip2 <schema/pgsql/schema.sql >pgsql.schema.sql.bz2

	cp -r /entrypoint .
	cp -r /rootfs .

	docker build -f /Dockerfile -t "${TARGET}:$TAG" .

	STATE_isPost=1 node /actions/checkout/dist/index.js
}

push () {
	test -n "$TAG"

	if [ "$(tr -d '\n' <<<"$DOCKER_HUB_PASSWORD" |wc -c)" -gt 0 ]; then
		docker login -u icingaadmin --password-stdin <<<"$DOCKER_HUB_PASSWORD"
		docker push "${TARGET}:$TAG"
		docker logout
	fi
}

case "$GITHUB_EVENT_NAME" in
	pull_request)
		grep -qEe '^refs/pull/[0-9]+' <<<"$GITHUB_REF"
		TAG="pr$(grep -oEe '[0-9]+' <<<"$GITHUB_REF")"
		mkimg
		;;
	push)
		grep -qEe '^refs/heads/.' <<<"$GITHUB_REF"
		TAG="$(cut -d / -f 3- <<<"$GITHUB_REF")"
		mkimg
		push
		;;
	release)
		grep -qEe '^refs/tags/v[0-9]' <<<"$GITHUB_REF"
		TAG="$(cut -d v -f 2- <<<"$GITHUB_REF")"
		mkimg
		push
		;;
	*)
		echo "Unknown event: $GITHUB_EVENT_NAME" >&2
		false
		;;
esac
