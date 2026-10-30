#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
if ./scripts/docker.sh compose version >/dev/null 2>&1; then
    exec ./scripts/docker.sh compose "$@"
fi
# Export the same context for the standalone Compose v2 executable.
if [ -z "${DOCKER_CONTEXT:-}" ] && [ -z "${DOCKER_HOST:-}" ] && [ -f deploy/docker-context ]; then
    IFS= read -r DOCKER_CONTEXT < deploy/docker-context
    export DOCKER_CONTEXT
fi
exec docker-compose "$@"
