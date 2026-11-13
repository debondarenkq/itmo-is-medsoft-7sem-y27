#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
network_name=${HIS_GATEWAY_NETWORK:-his_gateway}
if ! ./scripts/docker.sh network inspect "$network_name" >/dev/null 2>&1; then
    ./scripts/docker.sh network create "$network_name" >/dev/null
fi
