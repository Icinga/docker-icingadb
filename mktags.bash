#!/bin/bash
# Icinga DB Docker image | (c) 2023 Icinga GmbH | GPLv2+
set -exo pipefail

if [[ "$1" =~ ^v((([0-9]+).([0-9]+)).[0-9]+)$ ]]; then
  XYZ="${BASH_REMATCH[1]}"
  XY="${BASH_REMATCH[2]}"
  X="${BASH_REMATCH[3]}"
  Y="${BASH_REMATCH[4]}"

  BUILDX=(docker buildx build --platform "$(cat "$(realpath "$(dirname "$0")")/platforms.txt")" --push)
  cd "$(mktemp -d)"

  echo "FROM icinga/icingadb:$XYZ" >Dockerfile

  "${BUILDX[@]}" -t "icinga/icingadb:$XY" .

  NEXT="${X}.$(($Y+1))"

  case "$(curl --head -sSLo /dev/null -w '%{http_code}' "https://hub.docker.com/v2/namespaces/icinga/repositories/icingadb/tags/$NEXT")" in
    200)
      ;;
    404)
      "${BUILDX[@]}" -t "icinga/icingadb:$X" .
      "${BUILDX[@]}" -t icinga/icingadb .
      ;;
    *)
      echo "Can't check for icinga/icingadb:$NEXT"
      false
      ;;
  esac
fi
