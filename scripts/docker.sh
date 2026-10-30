#!/bin/sh
# A project-local context keeps this stack attached to the same Docker engine.
set -eu
cd "$(dirname "$0")/.."
if [ -z "${DOCKER_CONTEXT:-}" ] && [ -z "${DOCKER_HOST:-}" ] && [ -f deploy/docker-context ]; then
    IFS= read -r project_context < deploy/docker-context || true
    if [ -z "${project_context:-}" ]; then
        printf '%s\n' 'Пустой deploy/docker-context. Укажите имя Docker-контекста или удалите этот локальный файл.' >&2
        exit 1
    fi
    export DOCKER_CONTEXT="$project_context"
fi
exec docker "$@"
